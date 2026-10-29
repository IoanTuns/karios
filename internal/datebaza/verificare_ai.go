package datebaza

import (
	"context"
	"database/sql"
	"sort"
	"strings"
	"time"

	"pontix/internal/ai"
	"pontix/internal/bd/acces"
	"pontix/internal/calcul"
)

type SemnalPontajAI struct {
	AngajatID  string `json:"angajat_id"`
	Nume       string `json:"nume"`
	Marca      string `json:"marca"`
	Data       string `json:"data"`
	Valoare    string `json:"valoare"`
	Motiv      string `json:"motiv"`
	RegulaCod  string `json:"regula_cod"`
	RegulaNume string `json:"regula_nume"`
	Severitate string `json:"severitate"`
}
type VerificarePontajAI struct {
	Celule  int              `json:"celule"`
	Semnale []SemnalPontajAI `json:"semnale"`
}

// Verificarea este consultativă: nu deduce norme sau tarife și nu modifică pontajul.
func (s *Serviciu) VerificaPontajAI(ctx context.Context, i ai.IntentiePlan) (VerificarePontajAI, error) {
	r := VerificarePontajAI{Semnale: []SemnalPontajAI{}}
	tx, err := s.db.BeginTx(ctx, &sql.TxOptions{ReadOnly: true})
	if err != nil {
		return r, traduceEroare(err)
	}
	defer tx.Rollback()
	q := acces.New(tx)
	programe, err := q.ListaSchimburi(ctx)
	if err != nil {
		return r, traduceEroare(err)
	}
	i.CodSchimb = ""
	for _, p := range programe {
		if p.Activ == 1 {
			i.CodSchimb = p.Cod
			break
		}
	}
	i.ModExistente = "pastreaza"
	p, surse, err := construiestePlanAI(ctx, tx, q, i)
	if err != nil {
		return r, traduceEroare(err)
	}
	intervale := []calcul.IntervalPlanificare{}
	for _, interval := range surse.intervale {
		// Protecția unei zile cu pontaj efectiv este pentru scriere, nu un
		// interval de muncă. Verificarea consultativă folosește orele efective.
		if !strings.HasPrefix(interval.Cheie, "e/zi/") {
			intervale = append(intervale, interval)
		}
	}
	for _, v := range surse.pozitii {
		if v.Activ == 1 {
			pr := surse.schimburi[v.SchimbID]
			d, _ := time.Parse(time.DateOnly, v.Data)
			start := d.Add(time.Duration(pr.OraStart) * time.Minute)
			intervale = append(intervale, calcul.IntervalPlanificare{AngajatID: v.AngajatID, Cheie: "p/" + v.ID, Inceput: start, Sfarsit: start.Add(time.Duration(pr.DurataMinute) * time.Minute)})
		}
	}
	for _, v := range p.Randuri {
		if i.Selectie.FaraPontaj && v.Motiv == "fara_pontaj" {
			continue
		}
		r.Celule++
		cheie := v.AngajatID + "/" + v.Data
		valoare := strings.TrimSpace(surse.ciorne[cheie].Valoare)
		semnal := func(motiv string) {
			r.Semnale = append(r.Semnale, SemnalPontajAI{AngajatID: v.AngajatID, Nume: v.Nume, Marca: v.Marca, Data: v.Data, Valoare: valoare, Motiv: motiv})
		}
		if valoare == "" {
			if !surse.absente[cheie] {
				semnal("gol")
			}
			continue
		}
		program, ok := surse.coduri[strings.ToUpper(valoare)]
		if !ok {
			if _, manual := calcul.OreManuale(valoare); manual {
				semnal("ore_manuale")
				if surse.absente[cheie] {
					semnal("absenta")
				}
			} else {
				semnal("cod_invalid")
			}
			continue
		}
		if !program.Activ {
			semnal("program_inactiv")
		}
		ang := surse.oameni[v.AngajatID]
		if id := textSauGol(ang.CategorieProgramID); id != "" {
			cat, ok := surse.categorii[id]
			permis := false
			for _, id := range cat.SchimburiID {
				permis = permis || id == program.ID
			}
			if !ok || !cat.Activ || !permis {
				semnal("program_nepermis")
			}
		}
		ignora := map[string]bool{"c/" + cheie: true}
		for _, poz := range surse.pozitii {
			if poz.Activ == 1 && poz.AngajatID == v.AngajatID && poz.Data == v.Data {
				if poz.SchimbID == program.ID {
					ignora["p/"+poz.ID] = true
				} else {
					semnal("plan_diferit")
				}
			}
		}
		if motiv := calcul.VerificaTuraPlanificata(v.AngajatID, v.Data, program, surse.absente, intervale, ignora); motiv != "" {
			semnal(motiv)
		}
	}
	// Regulile configurate se evaluează pe schița întregii luni, pentru a nu
	// transforma un plafon lunar într-un plafon al intervalului vizibil.
	reguli := []calcul.RegulaConfigurabila{}
	rows, err := tx.QueryContext(ctx, `SELECT id,cod,denumire,descriere,tip_regula,severitate,coduri_schimb,valoare_numerica,activ FROM reguli_configurabile WHERE activ=1 ORDER BY id`)
	if err != nil {
		return r, traduceEroare(err)
	}
	for rows.Next() {
		var v calcul.RegulaConfigurabila
		if err := rows.Scan(&v.ID, &v.Cod, &v.Denumire, &v.Descriere, &v.TipRegula, &v.Severitate, &v.CoduriSchimb, &v.ValoareNumerica, &v.Activ); err != nil {
			rows.Close()
			return r, traduceEroare(err)
		}
		reguli = append(reguli, v)
	}
	err = rows.Err()
	rows.Close()
	if err != nil {
		return r, traduceEroare(err)
	}
	selectati := map[string]PersoanaPlanAI{}
	for _, a := range p.Persoane {
		selectati[a.ID] = a
	}
	reguliDupaID := map[string]calcul.RegulaConfigurabila{}
	for _, v := range reguli {
		reguliDupaID[v.ID] = v
	}
	for _, luna := range p.Luni {
		pozitii := []calcul.PozitiePentruRegula{}
		for _, c := range surse.ciorne {
			a, ok := selectati[c.AngajatID]
			if !ok || !strings.HasPrefix(c.Data, luna) {
				continue
			}
			pr, ok := surse.coduri[strings.ToUpper(strings.TrimSpace(c.Valoare))]
			if !ok {
				continue
			}
			pozitii = append(pozitii, calcul.PozitiePentruRegula{AngajatID: a.ID, Nume: a.Nume, Data: c.Data, CodSchimb: pr.Cod, OraStart: int(pr.OraStart), DurataMin: int(pr.DurataMinute), Activ: true})
		}
		for _, v := range calcul.EvalueazaReguli(reguli, pozitii) {
			a := selectati[v.AngajatID]
			r.Semnale = append(r.Semnale, SemnalPontajAI{AngajatID: a.ID, Nume: a.Nume, Marca: a.Marca, Data: v.Data, Motiv: "regula", RegulaCod: v.RegulaCod, RegulaNume: reguliDupaID[v.RegulaID].Denumire, Severitate: v.Severitate})
		}
	}
	sort.Slice(r.Semnale, func(a, b int) bool {
		x, y := r.Semnale[a], r.Semnale[b]
		return x.Marca+x.Data+x.Motiv+x.RegulaCod < y.Marca+y.Data+y.Motiv+y.RegulaCod
	})
	return r, traduceEroare(tx.Commit())
}
