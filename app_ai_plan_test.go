package main

import (
	"context"
	"testing"

	"pontix/internal/ai"
	"pontix/internal/datebaza"
)

func intentieDemoAI() ai.IntentiePlan {
	return ai.IntentiePlan{Selectie: ai.SelectiePlan{Tip: "toti", Marci: []string{}}, DeLa: "2026-10-01", PanaLa: "2026-11-30", CodSchimb: "Z7", FiltruZile: "lucratoare", ModExistente: "pastreaza"}
}

func TestPlanAIDouaLuniTotiAngajatii(t *testing.T) {
	a := appCuDemo(t)
	r := a.InterogheazaAI(CerereInterogareAI{An: 2026, Luna: 10, Istoric: []ai.MesajConversatie{{Rol: "user", Continut: "Toti angajati in octombrie si noiembrie lucreaza de la 7la 14"}}})
	if r.Eroare != nil || r.Raspuns == nil || r.Raspuns.Intentie == nil || len(r.Raspuns.Propuneri) > 0 {
		t.Fatalf("interogare: %+v", r)
	}
	i := *r.Raspuns.Intentie
	if i.Selectie.Tip != "toti" || i.DeLa != "2026-10-01" || i.PanaLa != "2026-11-30" || len(i.Clarificari) < 2 {
		t.Fatalf("intentie: %+v", i)
	}
	if p := a.PrevizualizeazaPlanAI(i); p.Eroare == nil {
		t.Fatal("acceptă lipsa clarificărilor")
	}
	i.FiltruZile = "lucratoare"
	i.ModExistente = "pastreaza"
	i.AnDeConfirmat = false // Utilizatorul confirmă anul afișat din context.
	i.CodSchimb = "Z7"      // Utilizatorul alege explicit între programele cu același interval.
	p := a.PrevizualizeazaPlanAI(i)
	if p.Eroare != nil {
		t.Fatalf("previzualizare: %+v", p.Eroare)
	}
	d := a.DateInitiale()
	activi := 0
	for _, ang := range d.Date.Angajati {
		if ang.Activ {
			activi++
		}
	}
	if len(p.Plan.Persoane) != activi || len(p.Plan.Luni) != 2 || p.Plan.Noi != activi*43 {
		t.Fatalf("plan: persoane=%d, luni=%v, noi=%d, blocate=%d", len(p.Plan.Persoane), p.Plan.Luni, p.Plan.Noi, p.Plan.Blocate)
	}
	var inainte int
	if err := a.db.QueryRow(`SELECT count(*) FROM pontaj_ciorne`).Scan(&inainte); err != nil {
		t.Fatal(err)
	}
	if inainte != 0 {
		t.Fatal("previzualizarea a scris date")
	}
	rez := a.ConfirmaPlanAI(p.Plan.Intentie, p.Plan.Token)
	if rez.Eroare != nil || rez.Plan.Aplicate != activi*43 {
		t.Fatalf("aplicare: %+v", rez.Eroare)
	}
	for _, luna := range []int{10, 11} {
		schita, err := datebaza.Nou(a.db).CitesteSchita(context.Background(), 2026, luna)
		if err != nil {
			t.Fatal(err)
		}
		pe := map[string]int{}
		for _, c := range schita.Celule {
			if c.Valoare == "Z7" {
				pe[c.AngajatID]++
			}
		}
		if len(pe) != activi {
			t.Fatalf("luna %d: %d persoane din %d", luna, len(pe), activi)
		}
	}
	if repetat := a.ConfirmaPlanAI(p.Plan.Intentie, p.Plan.Token); repetat.Eroare == nil || repetat.Eroare.Cod != "ai.plan_expirat" {
		t.Fatalf("reaplicare token: %+v", repetat)
	}
}

func TestPlanAIRefuzaDateSchimbate(t *testing.T) {
	for _, schimbare := range []string{"celula", "angajat", "absenta", "luna", "intentie"} {
		t.Run(schimbare, func(t *testing.T) {
			a := appCuDemo(t)
			i := intentieDemoAI()
			i.PanaLa = "2026-10-02"
			p := a.PrevizualizeazaPlanAI(i)
			if p.Eroare != nil {
				t.Fatal(p.Eroare)
			}
			ang := p.Plan.Persoane[0]
			s := datebaza.Nou(a.db)
			ctx := context.Background()
			switch schimbare {
			case "celula":
				if err := s.SalveazaSchita(ctx, 2026, 10, []datebaza.CelulaSchita{{AngajatID: ang.ID, Data: "2026-10-01", Valoare: "7"}}); err != nil {
					t.Fatal(err)
				}
			case "angajat":
				if err := s.SchimbaActiv(ctx, ang.ID, false); err != nil {
					t.Fatal(err)
				}
			case "absenta":
				if err := s.SalveazaAbsente(ctx, []datebaza.ZiAbsenta{{AngajatID: ang.ID, Data: "2026-10-01", Cod: "CO", Activ: true}}); err != nil {
					t.Fatal(err)
				}
			case "luna":
				if _, err := s.SeteazaStareLuna(ctx, 2026, 10, "BLOCATA", "test", "operator"); err != nil {
					t.Fatal(err)
				}
			case "intentie":
				i.PanaLa = "2026-10-03"
			}
			r := a.ConfirmaPlanAI(i, p.Plan.Token)
			if r.Eroare == nil || r.Eroare.Cod != "ai.plan_expirat" {
				t.Fatalf("nu refuză schimbarea: %+v", r)
			}
		})
	}
}

func TestPlanAIExcluderiSiSelectii(t *testing.T) {
	a := appCuDemo(t)
	i := intentieDemoAI()
	i.DeLa = "2026-10-01"
	i.PanaLa = "2026-10-02"
	p := a.PrevizualizeazaPlanAI(i)
	if p.Eroare != nil {
		t.Fatal(p.Eroare)
	}
	ang := p.Plan.Persoane[0]
	s := datebaza.Nou(a.db)
	ctx := context.Background()
	if err := s.SalveazaAbsente(ctx, []datebaza.ZiAbsenta{{AngajatID: ang.ID, Data: "2026-10-01", Cod: "CO", Activ: true}}); err != nil {
		t.Fatal(err)
	}
	if err := s.SalveazaSchita(ctx, 2026, 10, []datebaza.CelulaSchita{{AngajatID: ang.ID, Data: "2026-10-02", Valoare: "6"}}); err != nil {
		t.Fatal(err)
	}
	i.Selectie = ai.SelectiePlan{Tip: "persoane", Marci: []string{ang.Marca}}
	p = a.PrevizualizeazaPlanAI(i)
	if p.Eroare != nil {
		t.Fatal(p.Eroare)
	}
	if p.Plan.Noi != 0 || p.Plan.Ocupate != 1 || p.Plan.Blocate != 1 {
		t.Fatalf("excluderi: %+v", p.Plan)
	}
	i.ModExistente = "inlocuieste"
	p = a.PrevizualizeazaPlanAI(i)
	if p.Eroare != nil || p.Plan.Inlocuiri != 1 || p.Plan.Blocate != 1 {
		t.Fatalf("înlocuiri: %+v", p)
	}
	r := a.ConfirmaPlanAI(i, p.Plan.Token)
	if r.Eroare != nil || r.Plan.Aplicate != 1 {
		t.Fatalf("aplicare: %+v", r)
	}
	date, err := s.Citeste(ctx)
	if err != nil {
		t.Fatal(err)
	}
	for _, tip := range []string{"sectie", "departament", "categorie"} {
		valoare := ""
		for _, a := range date.Angajati {
			if a.Activ {
				switch tip {
				case "sectie":
					valoare = a.SectieID
				case "departament":
					valoare = a.DepartamentID
				case "categorie":
					valoare = a.CategorieProgramID
				}
				if valoare != "" {
					break
				}
			}
		}
		if valoare == "" {
			continue
		}
		i.Selectie = ai.SelectiePlan{Tip: tip, Valoare: valoare}
		p = a.PrevizualizeazaPlanAI(i)
		if p.Eroare != nil || len(p.Plan.Persoane) == 0 {
			t.Fatalf("%s: %+v", tip, p)
		}
		for _, pers := range p.Plan.Persoane {
			gasit := false
			for _, a := range date.Angajati {
				if a.ID == pers.ID {
					switch tip {
					case "sectie":
						gasit = a.SectieID == valoare
					case "departament":
						gasit = a.DepartamentID == valoare
					case "categorie":
						gasit = a.CategorieProgramID == valoare
					}
				}
			}
			if !gasit {
				t.Fatal("persoană în afara grupului")
			}
		}
	}
	i.Selectie = ai.SelectiePlan{Tip: "persoane", Marci: []string{"inexistent"}}
	if p := a.PrevizualizeazaPlanAI(i); p.Eroare == nil {
		t.Fatal("acceptă marca inexistentă")
	}
}

func TestPlanAIRollbackIntreLuni(t *testing.T) {
	a := appCuDemo(t)
	i := intentieDemoAI()
	i.DeLa = "2026-10-30"
	i.PanaLa = "2026-11-02"
	p := a.PrevizualizeazaPlanAI(i)
	if p.Eroare != nil {
		t.Fatal(p.Eroare)
	}
	if _, err := a.db.Exec(`CREATE TRIGGER proba_rollback_ai BEFORE INSERT ON pontaj_ciorne WHEN NEW.data='2026-11-02' BEGIN SELECT RAISE(ABORT, 'proba'); END`); err != nil {
		t.Fatal(err)
	}
	if r := a.ConfirmaPlanAI(i, p.Plan.Token); r.Eroare == nil {
		t.Fatal("lipsește eroarea simulată")
	}
	var numar int
	if err := a.db.QueryRow(`SELECT count(*) FROM pontaj_ciorne`).Scan(&numar); err != nil {
		t.Fatal(err)
	}
	if numar != 0 {
		t.Fatalf("salvare parțială: %d", numar)
	}
}
