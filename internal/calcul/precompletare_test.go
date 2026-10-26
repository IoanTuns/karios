package calcul

import "testing"

func TestPrecompletareCategoriiRespectaCalendarAbsenteNormeSiGarzi(t *testing.T) {
	programe := []ProgramPrecompletare{{ID: "z", Cod: "Z7", OraStart: 420, DurataMinute: 420, Activ: true}, {ID: "s", Cod: "S7", OraStart: 900, DurataMinute: 420, Activ: true}, {ID: "g", Cod: "G24", Tip: "GARDA", OraStart: 420, DurataMinute: 1440, Activ: true}}
	categorii := []CategoriePrecompletare{{ID: "dim", ImplicitID: "z", Activ: true}, {ID: "dupa", ImplicitID: "s", Activ: true}, {ID: "garzi", Activ: true}}
	angajati := []AngajatPrecompletare{{ID: "a", CategorieID: "dim", NormaMinute: 420, Activ: true}, {ID: "b", CategorieID: "dupa", NormaMinute: 420, Activ: true}, {ID: "g", CategorieID: "garzi", NormaMinute: 420, Activ: true}, {ID: "partial", CategorieID: "dim", NormaMinute: 240, Activ: true}, {ID: "fara", NormaMinute: 420, Activ: true}}
	absente := map[string]bool{"a/2026-10-05": true, "a/2026-10-07": true, "a/2026-10-08": true}
	intrari := []IntrarePontaj{{AngajatID: "a", Data: "2026-10-09", Valoare: "@M420"}, {AngajatID: "a", Data: "2026-10-12", Valoare: "G24"}}
	plan, err := PlanificaPrecompletare("2026-10-01", "2026-10-31", angajati, categorii, programe, intrari, map[string]bool{"2026-10-06": true}, absente, nil)
	if err != nil {
		t.Fatal(err)
	}
	if len(plan.Celule) != 36 || plan.Omise["absenta"] != 3 || plan.Omise["norma_diferita"] != 21 || plan.Omise["fara_implicit"] != 21 || plan.Omise["fara_categorie"] != 21 || plan.Omise["dupa_garda"] != 1 {
		t.Fatalf("plan: %+v", plan)
	}
	for _, c := range plan.Celule {
		if c.AngajatID == "a" && (c.Valoare != "Z7" || absente[c.AngajatID+"/"+c.Data] || c.Data == "2026-10-09" || c.Data == "2026-10-12" || c.Data == "2026-10-13") {
			t.Fatalf("dimineata incorecta: %+v", c)
		}
		if c.AngajatID == "b" && c.Valoare != "S7" {
			t.Fatalf("dupa-amiaza incorecta: %+v", c)
		}
		if c.AngajatID != "a" && c.AngajatID != "b" {
			t.Fatalf("program inventat: %+v", c)
		}
	}
}

func TestPrecompletareNuTraverseazaAbsentaSauZiLibera(t *testing.T) {
	a := []AngajatPrecompletare{{ID: "a", CategorieID: "n", NormaMinute: 540, Activ: true}}
	c := []CategoriePrecompletare{{ID: "n", ImplicitID: "n", Activ: true}}
	p := []ProgramPrecompletare{{ID: "n", Cod: "N9", OraStart: 1320, DurataMinute: 540, Activ: true}}
	plan, err := PlanificaPrecompletare("2026-10-01", "2026-10-02", a, c, p, nil, nil, map[string]bool{"a/2026-10-02": true}, nil)
	if err != nil || len(plan.Celule) != 0 || plan.Omise["conflict_absenta"] != 1 || plan.Omise["absenta"] != 1 {
		t.Fatalf("absenta: %+v %v", plan, err)
	}
	plan, err = PlanificaPrecompletare("2026-10-02", "2026-10-02", a, c, p, nil, nil, nil, nil)
	if err != nil || len(plan.Celule) != 0 || plan.Omise["trecere_zi_libera"] != 1 {
		t.Fatalf("weekend: %+v %v", plan, err)
	}
}

func TestPrecompletareSuprapunereLaTrecereDeLunaSiOraExacta(t *testing.T) {
	for _, sfarsit := range []int64{420, 600} {
		p := []ProgramPrecompletare{{ID: "z", Cod: "Z7", OraStart: 420, DurataMinute: 420, Activ: true}, {ID: "n", Cod: "N", OraStart: 1320, DurataMinute: 120 + sfarsit, Activ: true}}
		plan, err := PlanificaPrecompletare("2026-10-01", "2026-10-01", []AngajatPrecompletare{{ID: "a", CategorieID: "d", NormaMinute: 420, Activ: true}}, []CategoriePrecompletare{{ID: "d", ImplicitID: "z", Activ: true}}, p, []IntrarePontaj{{AngajatID: "a", Data: "2026-09-30", Valoare: "N"}}, nil, nil, nil)
		if err != nil {
			t.Fatal(err)
		}
		asteptat := 1
		if sfarsit == 600 {
			asteptat = 0
		}
		if len(plan.Celule) != asteptat {
			t.Fatalf("sfarsit %d: %+v", sfarsit, plan)
		}
	}
}

func TestPrecompletareProgrameIncerteInactiveSiPontajEfectiv(t *testing.T) {
	baza := ProgramPrecompletare{ID: "z", Cod: "Z7", OraStart: 420, DurataMinute: 420, Activ: true}
	for _, tc := range []struct {
		nume    string
		p       ProgramPrecompletare
		intrari []IntrarePontaj
		ocupate map[string]bool
		motiv   string
	}{
		{nume: "inactiv", p: ProgramPrecompletare{ID: "z"}, motiv: "schimb_inactiv"},
		{nume: "fara ora", p: ProgramPrecompletare{ID: "z", Activ: true, OraStart: -1}, motiv: "program_incomplet"},
		{nume: "necunoscut anterior", p: baza, intrari: []IntrarePontaj{{AngajatID: "a", Data: "2026-09-30", Valoare: "VECHI"}}, motiv: "suprapunere"},
		{nume: "efectiv", p: baza, ocupate: map[string]bool{"a/2026-10-01": true}, motiv: "existent"},
	} {
		t.Run(tc.nume, func(t *testing.T) {
			plan, err := PlanificaPrecompletare("2026-10-01", "2026-10-01", []AngajatPrecompletare{{ID: "a", CategorieID: "d", NormaMinute: 420, Activ: true}}, []CategoriePrecompletare{{ID: "d", ImplicitID: "z", Activ: true}}, []ProgramPrecompletare{tc.p}, tc.intrari, nil, nil, tc.ocupate)
			if err != nil || len(plan.Celule) != 0 || plan.Omise[tc.motiv] != 1 {
				t.Fatalf("plan: %+v %v", plan, err)
			}
		})
	}
}
