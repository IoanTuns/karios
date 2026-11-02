package calcul

import "testing"

func TestAbsenteCreditateAprobate(t *testing.T) {
	n := int64(420)
	for _, stare := range []string{"APROBATA", "IN_ASTEPTARE", "RESPINSA"} {
		r, e := CalculeazaPontaj("2026-10-01", "2026-10-31", ConfiguratiePontajInitiala(), nil, []IntrarePontaj{{AngajatID: "a", Data: "2026-10-05", Absenta: "DEPL", StareAprobare: stare, Crediteaza: true, MinuteCreditate: &n}}, nil)
		if e != nil {
			t.Fatal(e)
		}
		asteptat := int64(0)
		if stare == "APROBATA" {
			asteptat = n
		}
		if len(r.Celule) != 1 || r.Celule[0].Minute != asteptat {
			t.Fatalf("%s: %+v", stare, r)
		}
	}
	r, e := CalculeazaPontaj("2026-10-01", "2026-10-31", ConfiguratiePontajInitiala(), nil, []IntrarePontaj{{AngajatID: "a", Data: "2026-10-05", Absenta: "DEPL", StareAprobare: "APROBATA", Crediteaza: true}}, nil)
	if e != nil || r.Celule[0].Calculabil || r.Totaluri[0].Valori["ore_lucrate"].Definita {
		t.Fatalf("date lipsa: %+v %v", r, e)
	}
}
