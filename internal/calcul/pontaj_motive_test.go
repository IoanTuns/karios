package calcul

import "testing"

// F05-T20 (I049): fiecare valoare nedefinita spune de ce, cu codul in cauza;
// motivele nu se reunesc intr-un mesaj generic si nu schimba nicio suma.
func TestMotiveValoriNedefinite(t *testing.T) {
	programe := []ProgramPontaj{
		{Cod: "Z7", OraStart: 420, DurataMinute: 420},
		{Cod: "P8", OraStart: 420, DurataMinute: 480, PauzaMinute: 30},
		{Cod: "FS", OraStart: -1, DurataMinute: 420},
		{Cod: "D0", OraStart: 420, DurataMinute: 3000},
	}
	calculeaza := func(c ConfiguratiePontaj, intrari ...IntrarePontaj) map[string]ValoarePontaj {
		t.Helper()
		r, e := CalculeazaPontaj("2026-10-01", "2026-10-31", c, programe, intrari, nil)
		if e != nil || len(r.Totaluri) != 1 {
			t.Fatalf("%v %+v", e, r)
		}
		return r.Totaluri[0].Valori
	}
	zi := func(data, valoare string) IntrarePontaj {
		return IntrarePontaj{AngajatID: "a", Data: data, Valoare: valoare}
	}
	initiala := ConfiguratiePontajInitiala()
	calendar := ConfiguratiePontajInitiala()
	calendar.Atribuire = AtribuireCalendar
	verifica := func(nume string, v ValoarePontaj, motiv, cod string) {
		t.Helper()
		if v.Definita || v.Motiv != motiv || v.MotivCod != cod {
			t.Errorf("%s: %+v, astept motiv %q cod %q", nume, v, motiv, cod)
		}
	}

	// Regula neconfigurata: motivul vine din regula, nu din date.
	nedefinita := ""
	for _, r := range initiala.Reguli {
		if r.Sursa == "NEDEFINITA" {
			nedefinita = r.Coloana
			break
		}
	}
	if nedefinita == "" {
		t.Fatal("configuratia initiala nu are nicio coloana nedefinita")
	}
	v := calculeaza(initiala, zi("2026-10-05", "Z7"))
	verifica("regula", v[nedefinita], MotivRegulaNedefinita, "")
	if !v["T"].Definita || v["T"].Motiv != "" || v["T"].Sutimi != 420*100 {
		t.Errorf("valoare definita cu motiv sau suma schimbata: %+v", v["T"])
	}

	verifica("cod necunoscut", calculeaza(initiala, zi("2026-10-05", "xyz"))["T"], MotivValoareNecunoscuta, "XYZ")
	verifica("durata", calculeaza(initiala, zi("2026-10-05", "D0"))["T"], MotivDurataInvalida, "D0")
	verifica("pauza", calculeaza(initiala, zi("2026-10-05", "P8"))["T"], MotivPauza, "P8")
	n := (*int64)(nil)
	verifica("absenta", calculeaza(initiala, IntrarePontaj{AngajatID: "a", Data: "2026-10-05", Absenta: "DEPL", StareAprobare: "APROBATA", Crediteaza: true, MinuteCreditate: n})["T"], MotivAbsentaFaraOre, "DEPL")

	// Fara ora de inceput, atribuire la inceput: totalul e definit, doar ZNL si norma nu.
	v = calculeaza(initiala, zi("2026-10-05", "FS"))
	if !v["T"].Definita {
		t.Errorf("T trebuia definit: %+v", v["T"])
	}
	verifica("ZNL fara ora", v["ZNL"], MotivFaraOraInceput, "FS")
	verifica("norma fara ora", v["ore_norma"], MotivFaraOraInceput, "FS")
	// Atribuire calendaristica: fara ora nu se poate imparti pe zile, nici totalul.
	verifica("calendar fara ora", calculeaza(calendar, zi("2026-10-05", "FS"))["T"], MotivFaraOraInceput, "FS")

	// Mai multe motive: ordinea fixa si cel mai mic cod, indiferent de ordinea zilelor.
	for _, ordine := range [][]IntrarePontaj{
		{zi("2026-10-06", "P8"), zi("2026-10-07", "zz"), zi("2026-10-05", "aa")},
		{zi("2026-10-05", "aa"), zi("2026-10-07", "zz"), zi("2026-10-06", "P8")},
	} {
		verifica("ordine", calculeaza(initiala, ordine...)["T"], MotivValoareNecunoscuta, "AA")
	}
}

// Validarea configurarii indica coloana si campul greseli, nu un mesaj generic.
func TestValidareReguliIndicaCampul(t *testing.T) {
	cazuri := []struct {
		modifica  func(r *RegulaColoana)
		cod, camp string
	}{
		{func(r *RegulaColoana) { r.FactorSutimi = 0 }, "pontaj.regula_factor", "factor"},
		{func(r *RegulaColoana) { r.FactorSutimi = 1001 }, "pontaj.regula_factor", "factor"},
		{func(r *RegulaColoana) { r.Sursa = "ORE_SCHIMBURI"; r.Coduri = nil }, "pontaj.regula_fara_coduri", "coduri"},
		{func(r *RegulaColoana) { r.Coduri = []string{"Z7"} }, "pontaj.regula_coduri", "coduri"},
		{func(r *RegulaColoana) { r.Sursa = "ZILE_ABSENTA"; r.Coduri = []string{"CO"} }, "pontaj.regula_sursa", "sursa"},
		{func(r *RegulaColoana) { r.Sursa = "NECUNOSCUTA" }, "pontaj.regula_sursa", "sursa"},
	}
	for _, c := range cazuri {
		cfg := ConfiguratiePontajInitiala()
		for i := range cfg.Reguli {
			if cfg.Reguli[i].Coloana == "T" {
				c.modifica(&cfg.Reguli[i])
			}
		}
		err := ValideazaConfiguratiePontaj(cfg)
		e, ok := err.(*Eroare)
		if !ok || e.Cod != c.cod || e.Param["coloana"] != "T" || e.Param["camp"] != c.camp {
			t.Errorf("astept %s/%s pe T, primit %v", c.cod, c.camp, err)
		}
	}
	if err := ValideazaConfiguratiePontaj(ConfiguratiePontajInitiala()); err != nil {
		t.Fatalf("configuratia initiala: %v", err)
	}
}
