package calcul

import "testing"

func TestPontajOreCalendarSiAtribuire(t *testing.T) {
	for _, mod := range []string{AtribuireInceput, AtribuireCalendar} {
		t.Run(mod, func(t *testing.T) {
			c := ConfiguratiePontajInitiala()
			c.Atribuire = mod
			r, err := CalculeazaPontaj("2026-10-01", "2026-10-31", c, []ProgramPontaj{{Cod: "N9", OraStart: 1320, DurataMinute: 540}}, []IntrarePontaj{{AngajatID: "a", Data: "2026-10-23", Valoare: "N9"}}, nil)
			if err != nil {
				t.Fatal(err)
			}
			v := r.Totaluri[0].Valori
			for cod, asteptat := range map[string]string{"T": "9", "ore_norma": "2", "ZNL": "7", "ore_garda": "—"} {
				if obtinut := TextValoarePontaj(v[cod]); obtinut != asteptat {
					t.Errorf("%s=%s; asteptat %s", cod, obtinut, asteptat)
				}
			}
			if mod == AtribuireInceput && (len(r.Celule) != 1 || r.Celule[0].Minute != 540) {
				t.Fatal(r.Celule)
			}
			if mod == AtribuireCalendar && (len(r.Celule) != 2 || r.Celule[0].Minute != 120 || r.Celule[1].Minute != 420) {
				t.Fatal(r.Celule)
			}
		})
	}
}

func TestPontajTrecereLunaSiAn(t *testing.T) {
	p := []ProgramPontaj{{Cod: "N9", OraStart: 1320, DurataMinute: 540}}
	in := []IntrarePontaj{{AngajatID: "a", Data: "2026-12-31", Valoare: "N9"}}
	c := ConfiguratiePontajInitiala()
	r, e := CalculeazaPontaj("2026-12-01", "2026-12-31", c, p, in, map[string]bool{"2027-01-01": true})
	if e != nil || TextValoarePontaj(r.Totaluri[0].Valori["ZNL"]) != "7" {
		t.Fatal(r, e)
	}
	c.Atribuire = AtribuireCalendar
	for _, caz := range []struct{ inceput, sfarsit, total string }{{"2026-12-01", "2026-12-31", "2"}, {"2027-01-01", "2027-01-31", "7"}} {
		r, e = CalculeazaPontaj(caz.inceput, caz.sfarsit, c, p, in, map[string]bool{"2027-01-01": true})
		if e != nil || TextValoarePontaj(r.Totaluri[0].Valori["T"]) != caz.total {
			t.Fatal(r, e)
		}
	}
}

func TestPontajReguliExpliciteSiAbsente(t *testing.T) {
	c := ConfiguratiePontajInitiala()
	for i := range c.Reguli {
		if c.Reguli[i].Coloana == "ore_garda" {
			c.Reguli[i] = RegulaColoana{Coloana: "ore_garda", Sursa: "ORE_SCHIMBURI", Coduri: []string{"G24"}, FactorSutimi: 200}
		}
		if c.Reguli[i].Coloana == "CEV" {
			c.Reguli[i] = RegulaColoana{Coloana: "CEV", Sursa: "ZILE_ABSENTA", Coduri: []string{"SPECIAL"}, FactorSutimi: 100, NumaiAprobate: true}
		}
	}
	r, e := CalculeazaPontaj("2026-10-01", "2026-10-31", c, []ProgramPontaj{{Cod: "G24", OraStart: 420, DurataMinute: 1440}}, []IntrarePontaj{{AngajatID: "a", Data: "2026-10-05", Valoare: "G24"}, {AngajatID: "a", Data: "2026-10-06", Valoare: "G24", Absenta: "CO"}, {AngajatID: "a", Data: "2026-10-07", Absenta: "SPECIAL", StareAprobare: "APROBATA"}, {AngajatID: "a", Data: "2026-10-08", Absenta: "SPECIAL", StareAprobare: "IN_ASTEPTARE"}}, map[string]bool{"2026-10-05": true})
	if e != nil {
		t.Fatal(e)
	}
	for cod, asteptat := range map[string]string{"T": "24", "ZNL": "17", "ore_norma": "7", "CO": "1", "CEV": "1", "ore_garda": "48"} {
		if got := TextValoarePontaj(r.Totaluri[0].Valori[cod]); got != asteptat {
			t.Errorf("%s=%s; asteptat %s", cod, got, asteptat)
		}
	}
}

func TestPontajNecunoscutSiMinute(t *testing.T) {
	for _, caz := range []struct {
		raw          string
		p            ProgramPontaj
		mod          string
		total, norma string
	}{{"@M450", ProgramPontaj{}, AtribuireInceput, "7.5", "7.5"}, {"7", ProgramPontaj{Cod: "7", OraStart: 420, DurataMinute: 360}, AtribuireInceput, "6", "6"}, {"VECHI", ProgramPontaj{}, AtribuireInceput, "—", "—"}, {"P", ProgramPontaj{Cod: "P", OraStart: 420, DurataMinute: 480, PauzaMinute: 30}, AtribuireInceput, "—", "—"}, {"P", ProgramPontaj{Cod: "P", OraStart: -1, DurataMinute: 480}, AtribuireInceput, "8", "—"}, {"P", ProgramPontaj{Cod: "P", OraStart: -1, DurataMinute: 480}, AtribuireCalendar, "—", "—"}} {
		c := ConfiguratiePontajInitiala()
		c.Atribuire = caz.mod
		r, e := CalculeazaPontaj("2026-10-01", "2026-10-31", c, []ProgramPontaj{caz.p}, []IntrarePontaj{{AngajatID: "a", Data: "2026-10-05", Valoare: caz.raw}}, nil)
		if e != nil {
			t.Fatal(e)
		}
		if TextValoarePontaj(r.Totaluri[0].Valori["T"]) != caz.total || TextValoarePontaj(r.Totaluri[0].Valori["ore_norma"]) != caz.norma {
			t.Fatal(caz, r)
		}
	}
	for _, raw := range []string{"1.111", "NaN", "-1", "49", "@M-1", "@M2881"} {
		if _, ok := OreManuale(raw); ok {
			t.Errorf("acceptat %q", raw)
		}
	}
	c := ConfiguratiePontajInitiala()
	c.Reguli[0].FactorSutimi = 0
	if ValideazaConfiguratiePontaj(c) == nil {
		t.Fatal("factor invalid acceptat")
	}
}
