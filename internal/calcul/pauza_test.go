package calcul

import "testing"

// F04-T05 (D107): pauza se scade dupa modul ales pe program.
func TestPauzaPeModuri(t *testing.T) {
	ora := func(h, m int64) *int64 { v := h*60 + m; return &v }
	cazuri := []struct {
		nume      string
		p         ProgramPontaj
		atribuire string
		celule    []int64
		necalc    bool
	}{
		{"nestabilita", ProgramPontaj{Cod: "X", OraStart: 7 * 60, DurataMinute: 480, PauzaMinute: 30}, AtribuireInceput, nil, true},
		{"nu se scade", ProgramPontaj{Cod: "X", OraStart: 7 * 60, DurataMinute: 480, PauzaMinute: 30, ModPauza: PauzaNuSeScade}, AtribuireInceput, []int64{480}, false},
		{"din total", ProgramPontaj{Cod: "X", OraStart: 7 * 60, DurataMinute: 480, PauzaMinute: 30, ModPauza: PauzaDinTotal}, AtribuireInceput, []int64{450}, false},
		// Noapte 19-07 (vineri spre sambata), pauza 01:00-01:30: vineri 5 h, sambata 6,5 h.
		{"la ora pe zile", ProgramPontaj{Cod: "X", OraStart: 19 * 60, DurataMinute: 720, PauzaMinute: 30, ModPauza: PauzaLaOra, PauzaOra: ora(1, 0)}, AtribuireCalendar, []int64{300, 390}, false},
		// Mijlocul turei 19-07 e la 01:00; pauza de 60 min 00:30-01:30 cade sambata: vineri 5 h, sambata 6 h.
		{"la mijloc pe zile", ProgramPontaj{Cod: "X", OraStart: 19 * 60, DurataMinute: 720, PauzaMinute: 60, ModPauza: PauzaLaMijloc}, AtribuireCalendar, []int64{300, 360}, false},
		// Pauza peste miezul noptii (23:30-00:30): cate 30 de minute din fiecare zi.
		{"la ora peste miezul noptii", ProgramPontaj{Cod: "X", OraStart: 19 * 60, DurataMinute: 720, PauzaMinute: 60, ModPauza: PauzaLaOra, PauzaOra: ora(23, 30)}, AtribuireCalendar, []int64{270, 390}, false},
		// Din total pe zile: pauza se ia din ziua inceperii.
		{"din total pe zile", ProgramPontaj{Cod: "X", OraStart: 19 * 60, DurataMinute: 720, PauzaMinute: 60, ModPauza: PauzaDinTotal}, AtribuireCalendar, []int64{240, 420}, false},
		// Fara ora de inceput: doar totalul, cu pauza scazuta.
		{"fara ora de inceput", ProgramPontaj{Cod: "X", OraStart: -1, DurataMinute: 480, PauzaMinute: 30, ModPauza: PauzaLaMijloc}, AtribuireInceput, []int64{450}, false},
	}
	for _, caz := range cazuri {
		t.Run(caz.nume, func(t *testing.T) {
			c := ConfiguratiePontajInitiala()
			c.Atribuire = caz.atribuire
			r, err := CalculeazaPontaj("2026-10-01", "2026-10-31", c, []ProgramPontaj{caz.p}, []IntrarePontaj{{AngajatID: "a", Data: "2026-10-09", Valoare: "X"}}, nil)
			if err != nil {
				t.Fatal(err)
			}
			if caz.necalc {
				if v := r.Totaluri[0].Valori["T"]; v.Definita || v.Motiv != MotivPauza {
					t.Fatalf("asteptat necalculat cu motivul pauzei: %+v", v)
				}
				return
			}
			if len(r.Celule) != len(caz.celule) {
				t.Fatalf("celule %+v", r.Celule)
			}
			for i, m := range caz.celule {
				if r.Celule[i].Minute != m {
					t.Errorf("%s: %d, asteptat %d", r.Celule[i].Data, r.Celule[i].Minute, m)
				}
			}
		})
	}
}

// F04-T12 (I133): un program cu modul pauzei ales se precompleteaza si se planifica; norma se
// compara cu orele fara pauza; fara mod ales, programul ramane „incomplet”.
func TestPauzaInPrecompletareSiPlanificare(t *testing.T) {
	a := []AngajatPrecompletare{{ID: "a", CategorieID: "c", NormaMinute: 480, Activ: true}}
	c := []CategoriePrecompletare{{ID: "c", ImplicitID: "p", Activ: true}}
	cuPauza := ProgramPrecompletare{ID: "p", Cod: "P8", OraStart: 7 * 60, DurataMinute: 510, PauzaMinute: 30, ModPauza: PauzaDinTotal, Activ: true}
	plan, err := PlanificaPrecompletare("2026-10-05", "2026-10-05", a, c, []ProgramPrecompletare{cuPauza}, nil, nil, nil, nil)
	if err != nil || len(plan.Celule) != 1 || plan.Celule[0].Valoare != "P8" {
		t.Fatalf("precompletare cu pauza din total: %+v %v", plan, err)
	}
	nestabilit := cuPauza
	nestabilit.ModPauza = PauzaNestabilita
	plan, _ = PlanificaPrecompletare("2026-10-05", "2026-10-05", a, c, []ProgramPrecompletare{nestabilit}, nil, nil, nil, nil)
	if len(plan.Celule) != 0 || plan.Omise["program_incomplet"] != 1 {
		t.Fatalf("pauza nestabilita: %+v", plan)
	}
	if m := VerificaTuraPlanificata("a", "2026-10-05", cuPauza, nil, nil, nil); m != "" {
		t.Fatalf("planificare cu pauza: %q", m)
	}
	if m := VerificaTuraPlanificata("a", "2026-10-05", nestabilit, nil, nil, nil); m != "program_incomplet" {
		t.Fatalf("planificare fara mod: %q", m)
	}
	if cuPauza.MinuteLucrate() != 480 {
		t.Fatal(cuPauza.MinuteLucrate())
	}
}
