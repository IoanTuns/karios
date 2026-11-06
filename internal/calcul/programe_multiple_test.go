package calcul

import "testing"

func TestProgrameMultipleOreSiCalendar(t *testing.T) {
	p := []ProgramPontaj{{Cod: "G5", OraStart: 840, DurataMinute: 300}, {Cod: "G12", OraStart: 1140, DurataMinute: 720}}
	for _, mod := range []string{AtribuireInceput, AtribuireCalendar} {
		c := ConfiguratiePontajInitiala()
		c.Atribuire = mod
		r, e := CalculeazaPontaj("2026-10-01", "2026-10-31", c, p, []IntrarePontaj{{AngajatID: "a", Data: "2026-10-30", Valoare: "G5+G12"}}, nil)
		if e != nil {
			t.Fatal(e)
		}
		if v := TextValoarePontaj(r.Totaluri[0].Valori["T"]); v != "17" {
			t.Fatal(mod, v, r)
		}
		if v := TextValoarePontaj(r.Totaluri[0].Valori["ZNL"]); v != "7" {
			t.Fatal(mod, v, r)
		}
		if mod == AtribuireCalendar && (len(r.Celule) != 2 || r.Celule[0].Minute != 600 || r.Celule[1].Minute != 420) {
			t.Fatal(r.Celule)
		}
	}
	c := ConfiguratiePontajInitiala()
	n := int64(420)
	r, e := CalculeazaPontaj("2026-10-01", "2026-10-31", c, p, []IntrarePontaj{{AngajatID: "a", Data: "2026-10-30", Valoare: "G5+G12", Absenta: "DEL", StareAprobare: "APROBATA", Crediteaza: true, MinuteCreditate: &n}}, nil)
	if e != nil || TextValoarePontaj(r.Totaluri[0].Valori["T"]) != "7" {
		t.Fatal(r, e)
	}
}
