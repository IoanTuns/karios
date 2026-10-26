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

	"pontix/internal/bd/acces"
	"pontix/internal/calcul"
)

type PozitiePlanificare struct {
	ID         string `json:"id"`
	AngajatID  string `json:"angajat_id"`
	Data       string `json:"data"`
	SchimbID   string `json:"schimb_id"`
	Sectiune   string `json:"sectiune"`
	Activ      bool   `json:"activ"`
	Revizie    string `json:"revizie"`
	Stare      string `json:"stare"`
	PublicatLa string `json:"publicat_la"`
}
type PlanificareLuna struct {
	An      int                  `json:"an"`
	Luna    int                  `json:"luna"`
	Pozitii []PozitiePlanificare `json:"pozitii"`
}
type CererePublicarePlan struct {
	An          int      `json:"an"`
	Luna        int      `json:"luna"`
	Sectiune    string   `json:"sectiune"`
	PozitiiID   []string `json:"pozitii_id"`
	Inlocuieste bool     `json:"inlocuieste"`
	Token       string   `json:"token"`
}
type RandPublicarePlan struct {
	PozitieID     string `json:"pozitie_id"`
	AngajatID     string `json:"angajat_id"`
	Data          string `json:"data"`
	SchimbID      string `json:"schimb_id"`
	Valoare       string `json:"valoare"`
	Inainte       string `json:"inainte"`
	ReviziePlan   string `json:"revizie_plan"`
	ReviziePontaj string `json:"revizie_pontaj"`
	OraStart      int64  `json:"ora_start"`
	DurataMinute  int64  `json:"durata_minute"`
	Stare         string `json:"stare"`
	Motiv         string `json:"motiv"`
}
type PublicarePlan struct {
	Randuri     []RandPublicarePlan `json:"randuri"`
	Token       string              `json:"token"`
	Publicabile int                 `json:"publicabile"`
	Identice    int                 `json:"identice"`
	Blocate     int                 `json:"blocate"`
}
type sursePlanificare struct {
	pozitii   []acces.Planificare
	oameni    map[string]acces.ListaAngajatiRow
	schimburi map[string]calcul.ProgramPrecompletare
	coduri    map[string]calcul.ProgramPrecompletare
	categorii map[string]CategorieProgram
	absente   map[string]bool
	ciorne    map[string]acces.PontajCiorne
	intervale []calcul.IntervalPlanificare
}

func sursePentruPlanificare(ctx context.Context, q *acces.Queries, deLa, panaLa string) (sursePlanificare, error) {
	s := sursePlanificare{oameni: map[string]acces.ListaAngajatiRow{}, schimburi: map[string]calcul.ProgramPrecompletare{}, coduri: map[string]calcul.ProgramPrecompletare{}, categorii: map[string]CategorieProgram{}, absente: map[string]bool{}, ciorne: map[string]acces.PontajCiorne{}}
	prima, _ := time.Parse(time.DateOnly, deLa)
	ultima, _ := time.Parse(time.DateOnly, panaLa)
	extDeLa, extPanaLa := prima.AddDate(0, 0, -2).Format(time.DateOnly), ultima.AddDate(0, 0, 2).Format(time.DateOnly)
	var err error
	s.pozitii, err = q.ListaPlanificare(ctx, acces.ListaPlanificareParams{DataDeLa: extDeLa, DataPanaLa: extPanaLa})
	if err != nil {
		return s, err
	}
	oameni, err := q.ListaAngajati(ctx)
	if err != nil {
		return s, err
	}
	for _, a := range oameni {
		s.oameni[a.ID] = a
	}
	categorii, err := citesteCategoriiProgram(ctx, q)
	if err != nil {
		return s, err
	}
	for _, c := range categorii {
		s.categorii[c.ID] = c
	}
	schimburi, err := q.ListaSchimburi(ctx)
	if err != nil {
		return s, err
	}
	for _, p := range schimburi {
		program := calcul.ProgramPrecompletare{ID: p.ID, Cod: p.Cod, Tip: p.Tip, Regim: p.Regim, OraStart: p.OraStart, DurataMinute: p.DurataMinute, PauzaMinute: p.PauzaMinute, Activ: p.Activ == 1}
		s.schimburi[p.ID] = program
		s.coduri[strings.ToUpper(p.Cod)] = program
	}
	absente, err := q.ListaAbsenteInterval(ctx, acces.ListaAbsenteIntervalParams{DeLa: extDeLa, PanaLa: extPanaLa})
	if err != nil {
		return s, err
	}
	for _, a := range absente {
		if a.Activ == 1 {
			s.absente[a.AngajatID+"/"+a.Data] = true
		}
	}
	ciorne, err := q.ListaCiorneLuna(ctx, acces.ListaCiorneLunaParams{DataDeLa: extDeLa, DataPanaLa: extPanaLa})
	if err != nil {
		return s, err
	}
	for _, c := range ciorne {
		k := c.AngajatID + "/" + c.Data
		s.ciorne[k] = c
		if strings.TrimSpace(c.Valoare) == "" {
			continue
		}
		d, _ := time.Parse(time.DateOnly, c.Data)
		start, stop := d, d.AddDate(0, 0, 1)
		if p, ok := s.coduri[strings.ToUpper(strings.TrimSpace(c.Valoare))]; ok && p.OraStart >= 0 && p.DurataMinute > 0 {
			start = d.Add(time.Duration(p.OraStart) * time.Minute)
			stop = start.Add(time.Duration(p.DurataMinute) * time.Minute)
		} else if _, manual := calcul.OreManuale(c.Valoare); !manual {
			stop = d.AddDate(0, 0, 3)
		}
		s.intervale = append(s.intervale, calcul.IntervalPlanificare{AngajatID: c.AngajatID, Cheie: "c/" + k, Inceput: start, Sfarsit: stop})
	}
	efective, err := q.ListaPontajePrecompletare(ctx, acces.ListaPontajePrecompletareParams{DataDeLa: extDeLa, DataPanaLa: extPanaLa})
	if err != nil {
		return s, err
	}
	for i, p := range efective {
		d, _ := time.Parse(time.DateOnly, p.Data)
		start, stop := d, d.AddDate(0, 0, 1)
		if p.Start != nil && p.Sfarsit != nil {
			a, e1 := time.Parse(time.RFC3339, *p.Start)
			b, e2 := time.Parse(time.RFC3339, *p.Sfarsit)
			if e1 == nil && e2 == nil && b.After(a) {
				start = oraCivila(a)
				stop = oraCivila(b)
			}
		}
		s.intervale = append(s.intervale, calcul.IntervalPlanificare{AngajatID: p.AngajatID, Cheie: "e/" + strconv.Itoa(i), Inceput: start, Sfarsit: stop})
		// O zi cu pontaj efectiv nu poate fi înlocuită prin publicarea schiței.
		s.intervale = append(s.intervale, calcul.IntervalPlanificare{AngajatID: p.AngajatID, Cheie: "e/zi/" + strconv.Itoa(i), Inceput: d, Sfarsit: d.AddDate(0, 0, 1)})
	}
	return s, nil
}
func oraCivila(t time.Time) time.Time {
	v := t.In(calcul.Locatie)
	return time.Date(v.Year(), v.Month(), v.Day(), v.Hour(), v.Minute(), v.Second(), v.Nanosecond(), time.UTC)
}
func (s sursePlanificare) verifica(p acces.Planificare, cuPontaj bool, ignoraCelula bool) string {
	a, ok := s.oameni[p.AngajatID]
	if !ok || a.Activ != 1 {
		return "angajat_inactiv"
	}
	program, ok := s.schimburi[p.SchimbID]
	if !ok || !program.Activ {
		return "program_inactiv"
	}
	if !calcul.ProgramInSectiune(program, p.Sectiune) {
		return "sectiune_program"
	}
	if id := textSauGol(a.CategorieProgramID); id != "" {
		c, ok := s.categorii[id]
		if !ok || !c.Activ {
			return "categorie_inactiva"
		}
		permis := false
		for _, id := range c.SchimburiID {
			permis = permis || id == program.ID
		}
		if !permis {
			return "program_nepermis"
		}
	}
	intervale := []calcul.IntervalPlanificare{}
	if cuPontaj {
		intervale = append(intervale, s.intervale...)
	}
	for _, alt := range s.pozitii {
		if alt.Activ != 1 || alt.ID == p.ID {
			continue
		}
		schimb, ok := s.schimburi[alt.SchimbID]
		d, _ := time.Parse(time.DateOnly, alt.Data)
		start, stop := d, d.AddDate(0, 0, 3)
		if ok && schimb.OraStart >= 0 && schimb.DurataMinute > 0 {
			start = d.Add(time.Duration(schimb.OraStart) * time.Minute)
			stop = start.Add(time.Duration(schimb.DurataMinute) * time.Minute)
		}
		intervale = append(intervale, calcul.IntervalPlanificare{AngajatID: alt.AngajatID, Cheie: "p/" + alt.ID, Inceput: start, Sfarsit: stop})
	}
	ignora := map[string]bool{}
	if ignoraCelula {
		ignora["c/"+p.AngajatID+"/"+p.Data] = true
	}
	return calcul.VerificaTuraPlanificata(p.AngajatID, p.Data, program, s.absente, intervale, ignora)
}
func (s *Serviciu) CitestePlanificare(ctx context.Context, an, luna int) (PlanificareLuna, error) {
	rez := PlanificareLuna{An: an, Luna: luna, Pozitii: []PozitiePlanificare{}}
	deLa, panaLa, err := limiteLuna(an, luna)
	if err != nil {
		return rez, err
	}
	tx, err := s.db.BeginTx(ctx, &sql.TxOptions{ReadOnly: true})
	if err != nil {
		return rez, traduceEroare(err)
	}
	defer tx.Rollback()
	surse, err := sursePentruPlanificare(ctx, acces.New(tx), deLa, panaLa)
	if err != nil {
		return rez, traduceEroare(err)
	}
	for _, p := range surse.pozitii {
		if p.Data < deLa || p.Data > panaLa {
			continue
		}
		stare := "NEPUBLICATA"
		c := surse.ciorne[p.AngajatID+"/"+p.Data]
		if p.PublicatLa != "" {
			stare = "MODIFICATA"
		}
		if p.PublicatPlanRevizie == p.Revizie && p.PublicatPontajRevizie == c.Revizie && strings.EqualFold(p.PublicatValoare, strings.TrimSpace(c.Valoare)) && p.PublicatValoare == surse.schimburi[p.SchimbID].Cod {
			stare = "PUBLICATA"
		}
		if p.Activ != 1 {
			stare = "ANULATA"
		}
		rez.Pozitii = append(rez.Pozitii, PozitiePlanificare{ID: p.ID, AngajatID: p.AngajatID, Data: p.Data, SchimbID: p.SchimbID, Sectiune: p.Sectiune, Activ: p.Activ == 1, Revizie: strconv.FormatInt(p.Revizie, 10), Stare: stare, PublicatLa: p.PublicatLa})
	}
	if err := tx.Commit(); err != nil {
		return rez, traduceEroare(err)
	}
	return rez, nil
}
func (s *Serviciu) SalveazaPlanificare(ctx context.Context, an, luna int, pozitii []PozitiePlanificare) error {
	deLa, panaLa, err := limiteLuna(an, luna)
	if err != nil {
		return err
	}
	if len(pozitii) == 0 || len(pozitii) > 10000 {
		return eroare("planificare.lot_invalid")
	}
	return s.modifica(ctx, func(q *acces.Queries) error {
		surse, err := sursePentruPlanificare(ctx, q, deLa, panaLa)
		if err != nil {
			return err
		}
		vechi := map[string]acces.Planificare{}
		for _, p := range surse.pozitii {
			vechi[p.ID] = p
		}
		schimbate := map[string]acces.Planificare{}
		for _, p := range pozitii {
			if !dataValida(p.Data) || p.Data < deLa || p.Data > panaLa || (p.Sectiune != "GARZI" && p.Sectiune != "DUPA_AMIAZA") {
				return eroare("planificare.lot_invalid")
			}
			if p.ID == "" {
				if p.Revizie != "" || !p.Activ {
					return eroare("planificare.conflict")
				}
				p.ID = idNou()
			} else {
				v, ok := vechi[p.ID]
				if !ok || strconv.FormatInt(v.Revizie, 10) != p.Revizie {
					return eroare("planificare.conflict")
				}
				if v.AngajatID != p.AngajatID || v.Data != p.Data || v.Sectiune != p.Sectiune {
					return eroare("planificare.lot_invalid")
				}
			}
			if _, ok := schimbate[p.ID]; ok {
				return eroare("planificare.lot_invalid")
			}
			nou := vechi[p.ID]
			nou.ID = p.ID
			nou.AngajatID = p.AngajatID
			nou.Data = p.Data
			nou.SchimbID = p.SchimbID
			nou.Sectiune = p.Sectiune
			nou.Activ = numarBool(p.Activ)
			schimbate[p.ID] = nou
		}
		finale := []acces.Planificare{}
		for _, p := range surse.pozitii {
			if _, ok := schimbate[p.ID]; !ok {
				finale = append(finale, p)
			}
		}
		for _, p := range schimbate {
			finale = append(finale, p)
		}
		surse.pozitii = finale
		ocupate := map[string]bool{}
		for _, p := range finale {
			if p.Activ == 1 {
				k := p.AngajatID + "/" + p.Data
				if ocupate[k] {
					return eroare("planificare.duplicat")
				}
				ocupate[k] = true
			}
		}
		for _, p := range schimbate {
			if p.Activ == 1 {
				if motiv := surse.verifica(p, false, true); motiv != "" {
					return eroare("planificare.motiv_" + motiv)
				}
			}
			v, exista := vechi[p.ID]
			if exista && v.SchimbID == p.SchimbID && v.Activ == p.Activ {
				continue
			}
			if exista {
				n, err := q.ModificaPozitiePlanificare(ctx, acces.ModificaPozitiePlanificareParams{ID: p.ID, Revizie: v.Revizie, SchimbID: p.SchimbID, Activ: p.Activ, ModificatLa: moment()})
				if err != nil {
					return err
				}
				if n != 1 {
					return eroare("planificare.conflict")
				}
				p.Revizie = v.Revizie + 1
			} else {
				if err := q.AdaugaPozitiePlanificare(ctx, acces.AdaugaPozitiePlanificareParams{ID: p.ID, AngajatID: p.AngajatID, Data: p.Data, SchimbID: p.SchimbID, Sectiune: p.Sectiune, Activ: p.Activ, CreatLa: moment(), ModificatLa: moment()}); err != nil {
					return err
				}
				p.Revizie = 1
			}
			var inainte any
			if exista {
				inainte = v
			}
			if err := jurnal(ctx, q, "planificare", p.ID, inainte, p); err != nil {
				return err
			}
		}
		return nil
	})
}

func pregatestePublicare(ctx context.Context, q *acces.Queries, c CererePublicarePlan) (PublicarePlan, error) {
	rez := PublicarePlan{Randuri: []RandPublicarePlan{}}
	deLa, panaLa, err := limiteLuna(c.An, c.Luna)
	if err != nil {
		return rez, err
	}
	if len(c.PozitiiID) == 0 || len(c.PozitiiID) > 10000 || (c.Sectiune != "GARZI" && c.Sectiune != "DUPA_AMIAZA") {
		return rez, eroare("planificare.lot_invalid")
	}
	ids := map[string]bool{}
	for _, id := range c.PozitiiID {
		if id == "" || ids[id] {
			return rez, eroare("planificare.lot_invalid")
		}
		ids[id] = true
	}
	surse, err := sursePentruPlanificare(ctx, q, deLa, panaLa)
	if err != nil {
		return rez, err
	}
	for _, p := range surse.pozitii {
		if !ids[p.ID] {
			continue
		}
		if p.Data < deLa || p.Data > panaLa || p.Sectiune != c.Sectiune || p.Activ != 1 {
			return rez, eroare("planificare.lot_invalid")
		}
		program := surse.schimburi[p.SchimbID]
		veche := surse.ciorne[p.AngajatID+"/"+p.Data]
		r := RandPublicarePlan{PozitieID: p.ID, AngajatID: p.AngajatID, Data: p.Data, SchimbID: p.SchimbID, Valoare: program.Cod, Inainte: veche.Valoare, ReviziePlan: strconv.FormatInt(p.Revizie, 10), OraStart: program.OraStart, DurataMinute: program.DurataMinute, Stare: "NOUA"}
		if veche.Revizie > 0 {
			r.ReviziePontaj = strconv.FormatInt(veche.Revizie, 10)
		}
		r.Motiv = surse.verifica(p, true, true)
		if r.Motiv != "" {
			r.Stare = "BLOCATA"
		} else if strings.EqualFold(strings.TrimSpace(veche.Valoare), program.Cod) {
			r.Stare = "IDENTIC"
		} else if strings.TrimSpace(veche.Valoare) != "" {
			if c.Inlocuieste {
				r.Stare = "INLOCUIRE"
			} else {
				r.Stare = "BLOCATA"
				r.Motiv = "celula_ocupata"
			}
		}
		switch r.Stare {
		case "IDENTIC":
			rez.Identice++
		case "BLOCATA":
			rez.Blocate++
		default:
			rez.Publicabile++
		}
		rez.Randuri = append(rez.Randuri, r)
	}
	if len(rez.Randuri) != len(ids) {
		return rez, eroare("planificare.conflict")
	}
	sort.Slice(rez.Randuri, func(i, j int) bool { return rez.Randuri[i].PozitieID < rez.Randuri[j].PozitieID })
	date, err := json.Marshal(struct {
		An, Luna    int
		Sectiune    string
		Inlocuieste bool
		Randuri     []RandPublicarePlan
	}{c.An, c.Luna, c.Sectiune, c.Inlocuieste, rez.Randuri})
	if err != nil {
		return rez, err
	}
	hash := sha256.Sum256(date)
	rez.Token = hex.EncodeToString(hash[:])
	return rez, nil
}
func (s *Serviciu) PrevizualizeazaPublicarePlan(ctx context.Context, c CererePublicarePlan) (PublicarePlan, error) {
	tx, err := s.db.BeginTx(ctx, &sql.TxOptions{ReadOnly: true})
	if err != nil {
		return PublicarePlan{}, traduceEroare(err)
	}
	defer tx.Rollback()
	rez, err := pregatestePublicare(ctx, acces.New(tx), c)
	if err != nil {
		return rez, traduceEroare(err)
	}
	return rez, traduceEroare(tx.Commit())
}
func (s *Serviciu) PublicaPlanificare(ctx context.Context, c CererePublicarePlan) (PublicarePlan, error) {
	var rez PublicarePlan
	err := s.modifica(ctx, func(q *acces.Queries) error {
		var err error
		rez, err = pregatestePublicare(ctx, q, c)
		if err != nil {
			return err
		}
		if c.Token == "" || c.Token != rez.Token {
			return eroare("planificare.conflict_publicare")
		}
		for _, r := range rez.Randuri {
			if r.Stare == "BLOCATA" {
				continue
			}
			revizie, _ := strconv.ParseInt(r.ReviziePontaj, 10, 64)
			if r.Stare != "IDENTIC" {
				revizie, err = q.SalveazaCiorna(ctx, acces.SalveazaCiornaParams{AngajatID: r.AngajatID, Data: r.Data, Valoare: r.Valoare, ModificatLa: moment()})
				if err != nil {
					return err
				}
				if err := jurnal(ctx, q, "pontaj_ciorne", r.AngajatID+"/"+r.Data, r.Inainte, CelulaSchita{AngajatID: r.AngajatID, Data: r.Data, Valoare: r.Valoare, Revizie: strconv.FormatInt(revizie, 10)}); err != nil {
					return err
				}
			}
			planRev, _ := strconv.ParseInt(r.ReviziePlan, 10, 64)
			n, err := q.MarcheazaPlanificarePublicata(ctx, acces.MarcheazaPlanificarePublicataParams{ID: r.PozitieID, Revizie: planRev, PublicatPlanRevizie: planRev, PublicatPontajRevizie: revizie, PublicatValoare: r.Valoare, PublicatLa: moment()})
			if err != nil {
				return err
			}
			if n != 1 {
				return eroare("planificare.conflict_publicare")
			}
			if err := jurnal(ctx, q, "planificare_publicare", r.PozitieID, nil, r); err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return PublicarePlan{}, err
	}
	return rez, nil
}
