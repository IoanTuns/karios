package datebaza

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"sort"
	"strconv"
	"strings"
	"time"

	"pontix/internal/ai"
	"pontix/internal/bd/acces"
	"pontix/internal/calcul"
	"pontix/internal/i18n"
)

type PersoanaPlanAI struct {
	ID    string `json:"id"`
	Marca string `json:"marca"`
	Nume  string `json:"nume"`
}

type RandPlanAI struct {
	AngajatID string `json:"angajat_id"`
	Marca     string `json:"marca"`
	Nume      string `json:"nume"`
	Data      string `json:"data"`
	Inainte   string `json:"inainte"`
	Dupa      string `json:"dupa"`
	Revizie   string `json:"revizie"`
	Stare     string `json:"stare"`
	Motiv     string `json:"motiv"`
}

type PlanAI struct {
	Anulat            bool                   `json:"anulat"`
	OperatieID        string                 `json:"operatie_id"`
	SelectieExplicata []ExplicatieSelectieAI `json:"selectie_explicata"`
	SumarLuni         []SumarPlanAI          `json:"sumar_luni"`
	SumarPersoane     []SumarPlanAI          `json:"sumar_persoane"`
	Intentie          ai.IntentiePlan        `json:"intentie"`
	Persoane          []PersoanaPlanAI       `json:"persoane"`
	Luni              []string               `json:"luni"`
	Randuri           []RandPlanAI           `json:"randuri"`
	Token             string                 `json:"token"`
	Noi               int                    `json:"noi"`
	Inlocuiri         int                    `json:"inlocuiri"`
	Identice          int                    `json:"identice"`
	Ocupate           int                    `json:"ocupate"`
	Blocate           int                    `json:"blocate"`
	Aplicate          int                    `json:"aplicate"`
}

func normalPlan(s string) string {
	return strings.ToLower(strings.Join(strings.Fields(strings.NewReplacer("ă", "a", "â", "a", "î", "i", "ș", "s", "ş", "s", "ț", "t", "ţ", "t").Replace(s)), " "))
}

func numePlan(s string) string {
	var text i18n.Text
	if json.Unmarshal([]byte(s), &text) == nil {
		return text.In(i18n.Romana)
	}
	return s
}

func (s *Serviciu) PregatestePlanAI(ctx context.Context, intentie ai.IntentiePlan) (PlanAI, error) {
	tx, err := s.db.BeginTx(ctx, &sql.TxOptions{ReadOnly: true})
	if err != nil {
		return PlanAI{}, traduceEroare(err)
	}
	defer tx.Rollback()
	p, _, err := construiestePlanAI(ctx, tx, acces.New(tx), intentie)
	if err == nil {
		err = tx.Commit()
	}
	return p, traduceEroare(err)
}

// construiestePlanAI folosește aceeași tranzacție pentru date, selecție și token.
// Limitele de un an și 10000 de celule protejează memoria și durata confirmării.
func construiestePlanAI(ctx context.Context, tx *sql.Tx, q *acces.Queries, i ai.IntentiePlan) (PlanAI, sursePlanificare, error) {
	p := PlanAI{Intentie: i, Persoane: []PersoanaPlanAI{}, Luni: []string{}, Randuri: []RandPlanAI{}}
	var surse sursePlanificare
	prima, e1 := time.Parse(time.DateOnly, i.DeLa)
	ultima, e2 := time.Parse(time.DateOnly, i.PanaLa)
	if i.AnDeConfirmat || e1 != nil || e2 != nil || prima.Year() < 2000 || ultima.Year() > 2100 || ultima.Before(prima) || ultima.Sub(prima) > 366*24*time.Hour ||
		(i.FiltruZile != "lucratoare" && i.FiltruZile != "toate" && i.FiltruZile != "weekend") || (i.ModExistente != "pastreaza" && i.ModExistente != "inlocuieste") {
		return p, surse, eroare("ai.plan_incomplet")
	}
	deLa := time.Date(prima.Year(), prima.Month(), 1, 0, 0, 0, 0, time.UTC).Format(time.DateOnly)
	panaLa := time.Date(ultima.Year(), ultima.Month()+1, 0, 0, 0, 0, 0, time.UTC).Format(time.DateOnly)
	var err error
	surse, err = sursePentruPlanificare(ctx, q, deLa, panaLa)
	if err != nil {
		return p, surse, err
	}
	program, ok := surse.coduri[strings.ToUpper(i.CodSchimb)]
	if !ok || !program.Activ {
		return p, surse, eroare("ai.plan_program_invalid")
	}
	// Rezolvăm numele afișat sau ID-ul grupului, fără potriviri parțiale.
	grupID := i.Selectie.Valoare
	switch i.Selectie.Tip {
	case "toti", "persoane":
	case "sectie":
		lista, e := q.ListaSectii(ctx)
		if e != nil {
			return p, surse, e
		}
		gasite := 0
		for _, g := range lista {
			if g.ID == i.Selectie.Valoare || normalPlan(numePlan(g.Denumire)) == normalPlan(i.Selectie.Valoare) {
				grupID = g.ID
				gasite++
			}
		}
		if gasite != 1 {
			return p, surse, eroare("ai.plan_selectie_invalida")
		}
	case "departament":
		lista, e := q.ListaDepartamente(ctx)
		if e != nil {
			return p, surse, e
		}
		gasite := 0
		for _, g := range lista {
			if g.ID == i.Selectie.Valoare || normalPlan(numePlan(g.Denumire)) == normalPlan(i.Selectie.Valoare) {
				grupID = g.ID
				gasite++
			}
		}
		if gasite != 1 {
			return p, surse, eroare("ai.plan_selectie_invalida")
		}
	case "categorie":
		gasite := 0
		for _, g := range surse.categorii {
			if g.ID == i.Selectie.Valoare || normalPlan(g.Denumire.In(i18n.Romana)) == normalPlan(i.Selectie.Valoare) {
				grupID = g.ID
				gasite++
			}
		}
		if gasite != 1 {
			return p, surse, eroare("ai.plan_selectie_invalida")
		}
	default:
		return p, surse, eroare("ai.plan_selectie_invalida")
	}
	marci := map[string]bool{}
	for _, m := range i.Selectie.Marci {
		if m == "" || marci[m] {
			return p, surse, eroare("ai.plan_selectie_invalida")
		}
		marci[m] = true
	}
	gasite := 0
	p.SelectieExplicata = []ExplicatieSelectieAI{}
	for _, a := range surse.oameni {
		if a.Activ != 1 {
			p.SelectieExplicata = append(p.SelectieExplicata, ExplicatieSelectieAI{a.ID, a.Marca, strings.TrimSpace(a.Nume + " " + a.Prenume), false, "inactiv"})
			continue
		}
		inclus := false
		switch i.Selectie.Tip {
		case "toti":
			inclus = true
		case "sectie":
			inclus = textSauGol(a.SectieID) == grupID
		case "departament":
			inclus = textSauGol(a.DepartamentID) == grupID
		case "categorie":
			inclus = textSauGol(a.CategorieProgramID) == grupID
		case "persoane":
			inclus = marci[a.Marca]
			if inclus {
				gasite++
			}
		}
		if inclus {
			p.Persoane = append(p.Persoane, PersoanaPlanAI{a.ID, a.Marca, strings.TrimSpace(a.Nume + " " + a.Prenume)})
		}
		motivSelectie := "alt_grup"
		if inclus {
			motivSelectie = i.Selectie.Tip
		}
		p.SelectieExplicata = append(p.SelectieExplicata, ExplicatieSelectieAI{a.ID, a.Marca, strings.TrimSpace(a.Nume + " " + a.Prenume), inclus, motivSelectie})
	}
	if len(p.Persoane) == 0 || (i.Selectie.Tip == "persoane" && gasite != len(marci)) {
		return p, surse, eroare("ai.plan_selectie_invalida")
	}
	sort.Slice(p.Persoane, func(x, y int) bool { return p.Persoane[x].Marca < p.Persoane[y].Marca })
	sort.Slice(p.SelectieExplicata, func(x, y int) bool { return p.SelectieExplicata[x].Marca < p.SelectieExplicata[y].Marca })
	stari := map[string]string{}
	for d := time.Date(prima.Year(), prima.Month(), 1, 0, 0, 0, 0, time.UTC); !d.After(ultima); d = d.AddDate(0, 1, 0) {
		luna := d.Format("2006-01")
		p.Luni = append(p.Luni, luna)
		stare := "DESCHISA"
		err := tx.QueryRowContext(ctx, `SELECT stare FROM stare_luni WHERE an=? AND luna=?`, d.Year(), int(d.Month())).Scan(&stare)
		if err != nil && err != sql.ErrNoRows {
			return p, surse, err
		}
		stari[luna] = stare
	}
	fara := map[string]bool{}
	if i.Selectie.FaraPontaj {
		for _, c := range surse.ciorne {
			if strings.TrimSpace(c.Valoare) != "" {
				fara[c.AngajatID+"/"+c.Data[:7]] = true
			}
		}
	}
	// Intervale din schițe, pontaje efective și planificări, inclusiv turele vecine.
	intervale := append([]calcul.IntervalPlanificare(nil), surse.intervale...)
	for _, v := range surse.pozitii {
		if v.Activ != 1 {
			continue
		}
		pr := surse.schimburi[v.SchimbID]
		d, _ := time.Parse(time.DateOnly, v.Data)
		start := d.Add(time.Duration(pr.OraStart) * time.Minute)
		stop := start.Add(time.Duration(pr.DurataMinute) * time.Minute)
		intervale = append(intervale, calcul.IntervalPlanificare{AngajatID: v.AngajatID, Cheie: "p/" + v.ID, Inceput: start, Sfarsit: stop})
	}
	for _, a := range p.Persoane {
		for d := prima; !d.After(ultima); d = d.AddDate(0, 0, 1) {
			weekend := d.Weekday() == time.Saturday || d.Weekday() == time.Sunday
			if (i.FiltruZile == "lucratoare" && weekend) || (i.FiltruZile == "weekend" && !weekend) {
				continue
			}
			if len(p.Randuri) >= 10000 {
				return p, surse, eroare("ai.plan_prea_mare")
			}
			data := d.Format(time.DateOnly)
			cheie := a.ID + "/" + data
			veche := surse.ciorne[cheie]
			r := RandPlanAI{AngajatID: a.ID, Marca: a.Marca, Nume: a.Nume, Data: data, Inainte: veche.Valoare, Dupa: program.Cod, Stare: "noua"}
			if veche.Revizie > 0 {
				r.Revizie = strconv.FormatInt(veche.Revizie, 10)
			}
			ignora := map[string]bool{"c/" + cheie: true}
			planExistent := false
			sectiuneDiferita := false
			sectiune := "GARZI"
			if calcul.ProgramInSectiune(program, "DUPA_AMIAZA") {
				sectiune = "DUPA_AMIAZA"
			}
			for _, v := range surse.pozitii {
				if v.Activ == 1 && v.AngajatID == a.ID && v.Data == data {
					planExistent = true
					sectiuneDiferita = sectiuneDiferita || v.Sectiune != sectiune
					ignora["p/"+v.ID] = true
				}
			}
			motiv := ""
			switch {
			case stari[data[:7]] != "DESCHISA":
				motiv = "luna_inchisa"
			case fara[a.ID+"/"+data[:7]]:
				motiv = "fara_pontaj"
			case surse.absente[cheie]:
				motiv = "absenta"
			}
			ang := surse.oameni[a.ID]
			if id := textSauGol(ang.CategorieProgramID); id != "" && motiv == "" {
				cat, ok := surse.categorii[id]
				permis := false
				for _, id := range cat.SchimburiID {
					permis = permis || id == program.ID
				}
				if !ok || !cat.Activ || !permis {
					motiv = "program_nepermis"
				}
			}
			if motiv == "" {
				switch {
				case strings.EqualFold(strings.TrimSpace(veche.Valoare), program.Cod):
					r.Stare = "identica"
				case (strings.TrimSpace(veche.Valoare) != "" || planExistent) && i.ModExistente == "pastreaza":
					r.Stare = "ocupata"
				case planExistent && !calcul.ProgramDePlanificat(program):
					motiv = "planificare_existenta"
				case planExistent && sectiuneDiferita:
					motiv = "planificare_existenta"
				default:
					motiv = calcul.VerificaTuraPlanificata(a.ID, data, program, surse.absente, intervale, ignora)
					if veche.Valoare != "" {
						r.Stare = "inlocuire"
					}
				}
			}
			if motiv != "" {
				r.Stare = "blocata"
				r.Motiv = motiv
			}
			switch r.Stare {
			case "noua":
				p.Noi++
			case "inlocuire":
				p.Inlocuiri++
			case "identica":
				p.Identice++
			case "ocupata":
				p.Ocupate++
			default:
				p.Blocate++
			}
			if r.Stare == "noua" || r.Stare == "inlocuire" {
				start := d.Add(time.Duration(program.OraStart) * time.Minute)
				intervale = append(intervale, calcul.IntervalPlanificare{AngajatID: a.ID, Cheie: "nou/" + cheie, Inceput: start, Sfarsit: start.Add(time.Duration(program.DurataMinute) * time.Minute)})
			}
			p.Randuri = append(p.Randuri, r)
		}
	}
	p.SumarLuni, p.SumarPersoane = rezumaPlanAI(p)
	date, err := json.Marshal(struct {
		Plan  PlanAI
		Surse surseAmprentaAI
		Stari map[string]string
	}{p, surseAmprentaAI{surse.oameni, surse.categorii, surse.schimburi, surse.absente, surse.ciorne, surse.pozitii, surse.intervale}, stari})
	if err != nil {
		return p, surse, err
	}
	suma := sha256.Sum256(date)
	p.Token = hex.EncodeToString(suma[:])
	return p, surse, nil
}

// Câmpuri exportate pentru serializarea completă și stabilă a surselor tokenului.
type surseAmprentaAI struct {
	Oameni    map[string]acces.ListaAngajatiRow
	Categorii map[string]CategorieProgram
	Programe  map[string]calcul.ProgramPrecompletare
	Absente   map[string]bool
	Ciorne    map[string]acces.PontajCiorne
	Pozitii   []acces.Planificare
	Intervale []calcul.IntervalPlanificare
}

// AplicaPlanAI reconstruiește planul în tranzacția de scriere. Nu acceptă
// celulele venite din UI și nu salvează parțial un interval dintre luni.
func (s *Serviciu) AplicaPlanAI(ctx context.Context, i ai.IntentiePlan, token string) (PlanAI, error) {
	var plan PlanAI
	err := s.modificaTx(ctx, func(tx *sql.Tx, q *acces.Queries) error {
		p, surse, err := construiestePlanAI(ctx, tx, q, i)
		if err != nil {
			return err
		}
		if token == "" || token != p.Token {
			return eroare("ai.plan_expirat")
		}
		if p.Noi+p.Inlocuiri == 0 {
			return eroare("ai.fara_zile")
		}
		program := surse.coduri[strings.ToUpper(i.CodSchimb)]
		instantanee := []InstantaneuCelulaAI{}
		for _, r := range p.Randuri {
			if r.Stare != "noua" && r.Stare != "inlocuire" {
				continue
			}
			instantaneu, err := citesteInstantaneuAI(ctx, tx, r.AngajatID, r.Data)
			if err != nil {
				return err
			}
			if calcul.ProgramDePlanificat(program) {
				sectiune := "GARZI"
				if calcul.ProgramInSectiune(program, "DUPA_AMIAZA") {
					sectiune = "DUPA_AMIAZA"
				}
				var veche *acces.Planificare
				for x := range surse.pozitii {
					v := &surse.pozitii[x]
					if v.Activ == 1 && v.AngajatID == r.AngajatID && v.Data == r.Data {
						veche = v
						break
					}
				}
				id := idNou()
				if veche != nil {
					if veche.Sectiune != sectiune {
						return eroare("ai.plan_expirat")
					}
					id = veche.ID
					n, e := q.ModificaPozitiePlanificare(ctx, acces.ModificaPozitiePlanificareParams{ID: id, Revizie: veche.Revizie, SchimbID: program.ID, Activ: 1, ModificatLa: moment()})
					if e != nil {
						return e
					}
					if n != 1 {
						return eroare("ai.plan_expirat")
					}
				} else if err := q.AdaugaPozitiePlanificare(ctx, acces.AdaugaPozitiePlanificareParams{ID: id, AngajatID: r.AngajatID, Data: r.Data, SchimbID: program.ID, Sectiune: sectiune, Activ: 1, CreatLa: moment(), ModificatLa: moment()}); err != nil {
					return err
				}
				if err := jurnal(ctx, q, "planificare", id, veche, PozitiePlanificare{ID: id, AngajatID: r.AngajatID, Data: r.Data, SchimbID: program.ID, Sectiune: sectiune, Activ: true}); err != nil {
					return err
				}
			}
			revizie, err := q.SalveazaCiorna(ctx, acces.SalveazaCiornaParams{AngajatID: r.AngajatID, Data: r.Data, Valoare: r.Dupa, ModificatLa: moment()})
			if err != nil {
				return err
			}
			var inainte any
			if v, ok := surse.ciorne[r.AngajatID+"/"+r.Data]; ok {
				inainte = v
			}
			if err := jurnal(ctx, q, "pontaj_ciorne", r.AngajatID+"/"+r.Data, inainte, CelulaSchita{AngajatID: r.AngajatID, Data: r.Data, Valoare: r.Dupa, Revizie: strconv.FormatInt(revizie, 10)}); err != nil {
				return err
			}
			// Recitirea verifică rezultatul efectiv înainte de commit, cu rollback la eroare.
			v, err := q.CitesteCiorna(ctx, acces.CitesteCiornaParams{AngajatID: r.AngajatID, Data: r.Data})
			if err != nil {
				return err
			}
			if v.Valoare != r.Dupa || v.Revizie != revizie {
				return eroare("ai.plan_verificare_esuat")
			}
			p.Aplicate++
			dupa, err := citesteInstantaneuAI(ctx, tx, r.AngajatID, r.Data)
			if err != nil {
				return err
			}
			instantaneu.Dupa = &dupa
			instantanee = append(instantanee, instantaneu)
		}
		if calcul.ProgramDePlanificat(program) {
			pozitii, err := q.ListaPlanificare(ctx, acces.ListaPlanificareParams{DataDeLa: i.DeLa, DataPanaLa: i.PanaLa})
			if err != nil {
				return err
			}
			verificate := map[string]bool{}
			for _, v := range pozitii {
				if v.Activ == 1 && v.SchimbID == program.ID {
					verificate[v.AngajatID+"/"+v.Data] = true
				}
			}
			for _, r := range p.Randuri {
				if (r.Stare == "noua" || r.Stare == "inlocuire") && !verificate[r.AngajatID+"/"+r.Data] {
					return eroare("ai.plan_verificare_esuat")
				}
			}
		}
		p.OperatieID = idNou()
		continut, err := json.Marshal(p)
		if err != nil {
			return err
		}
		imagini, err := json.Marshal(instantanee)
		if err != nil {
			return err
		}
		if _, err := tx.ExecContext(ctx, `INSERT INTO operatii_ai(id,creat_la,plan,instantanee) VALUES(?,?,?,?)`, p.OperatieID, moment(), string(continut), string(imagini)); err != nil {
			return err
		}
		if err := jurnal(ctx, q, "operatii_ai", p.OperatieID, nil, p); err != nil {
			return err
		}
		plan = p
		return nil
	})
	return plan, err
}
