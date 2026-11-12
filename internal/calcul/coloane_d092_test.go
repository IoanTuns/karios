package calcul

import "testing"

// D092: TZLN, TOCN, OL, CVE, CFP, DLG, FMP, cu valorile cerute de Ioan Tuns.
func TestColoaneD092(t *testing.T) {
	// Octombrie 2026: 22 de zile luni-vineri. Sarbatoare fictiva intr-o miercuri (2026-10-14)
	// si una intr-o sambata (2026-10-17), care nu se scade de doua ori.
	libere := map[string]bool{"2026-10-14": true, "2026-10-17": true}
	sapte, patru := int64(420), int64(240)
	credit := int64(420)
	c := ConfiguratiePontajInitiala()
	c.Norme = map[string]*int64{"a": &sapte, "b": &patru, "g": nil}
	intrari := []IntrarePontaj{
		{AngajatID: "a", Data: "2026-10-05", Valoare: "8"},
		{AngajatID: "a", Data: "2026-10-06", Absenta: "FMP", StareAprobare: "APROBATA", Crediteaza: true, MinuteCreditate: &credit},
		{AngajatID: "a", Data: "2026-10-07", Absenta: "DLG", StareAprobare: "APROBATA", Crediteaza: true, MinuteCreditate: &credit},
		{AngajatID: "a", Data: "2026-10-08", Absenta: "CO", StareAprobare: "APROBATA"},
		{AngajatID: "a", Data: "2026-10-09", Absenta: "CFP", StareAprobare: "APROBATA"},
		{AngajatID: "a", Data: "2026-10-12", Absenta: "CVE", StareAprobare: "APROBATA"},
		{AngajatID: "a", Data: "2026-10-13", Absenta: "SPEC", StareAprobare: "APROBATA"},
		{AngajatID: "a", Data: "2026-10-15", Absenta: "CM", StareAprobare: "APROBATA"},
	}
	r, err := CalculeazaPontaj("2026-10-01", "2026-10-31", c, nil, intrari, libere)
	if err != nil {
		t.Fatal(err)
	}
	val := map[string]map[string]ValoarePontaj{}
	for _, x := range r.Totaluri {
		val[x.AngajatID] = x.Valori
	}
	verifica := func(id, col string, asteptat string) {
		t.Helper()
		v, ok := val[id][col]
		if !ok || TextValoarePontaj(v) != asteptat {
			t.Errorf("%s/%s = %q (%+v), astept %q", id, col, TextValoarePontaj(v), v, asteptat)
		}
	}
	verifica("a", "TZLN", "21") // 22 - 1 sarbatoare luni-vineri
	verifica("a", "TOCN", "147")
	verifica("b", "TOCN", "84") // norma partiala 4 h
	verifica("b", "TZLN", "21") // si fara pontaj
	if v := val["g"]["TOCN"]; v.Definita || v.Motiv != MotivFaraNorma {
		t.Errorf("TOCN fara norma: %+v", v)
	}
	verifica("a", "OL", "15") // 8 h salvate + 7 h FMP; DLG nu intra
	verifica("a", "T", "22")  // totalul vechi include toate creditele
	verifica("a", "CO", "1")
	verifica("a", "CM", "1")
	verifica("a", "CFP", "1")
	verifica("a", "DLG", "1")
	verifica("a", "FMP", "1")
	verifica("a", "CVE", "2") // CVE si SPEC; fara CO, CM, CFP, DLG, FMP
	// TOCN nu scade cu concediile: aceeasi valoare cu sau fara absente.
	r2, _ := CalculeazaPontaj("2026-10-01", "2026-10-31", c, nil, nil, libere)
	for _, x := range r2.Totaluri {
		if x.AngajatID == "a" && TextValoarePontaj(x.Valori["TOCN"]) != "147" {
			t.Errorf("TOCN depinde de concedii: %+v", x.Valori["TOCN"])
		}
	}
	// Grila arata cele 7 coloane cerute (D092) si orele suplimentare dupa OL (D107).
	if g := ColoaneGrilaPontaj(); len(g) != 8 || g[0] != "OL" || g[1] != "OS" || g[7] != "FMP" {
		t.Errorf("grila: %v", g)
	}
	if len(ColoaneRapoarte()) != 12 || ColoaneRapoarte()[11] != "ore_garda" {
		t.Errorf("rapoartele trebuie sa pastreze lista veche: %v", ColoaneRapoarte())
	}
}

func TestZileLucratoareNete(t *testing.T) {
	cazuri := []struct {
		deLa, panaLa string
		libere       map[string]bool
		n            int64
	}{
		{"2026-10-01", "2026-10-31", nil, 22},
		{"2026-02-01", "2026-02-28", nil, 20},
		{"2026-12-01", "2026-12-31", map[string]bool{"2026-12-01": true, "2026-12-25": true, "2026-12-26": true}, 21},
	}
	for _, c := range cazuri {
		if n := zileLucratoareNete(c.deLa, c.panaLa, c.libere); n != c.n {
			t.Errorf("%s..%s: %d, astept %d", c.deLa, c.panaLa, n, c.n)
		}
	}
}
