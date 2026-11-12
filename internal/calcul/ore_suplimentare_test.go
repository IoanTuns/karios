package calcul

import "testing"

// F04-T04 (D107): orele suplimentare sunt orele lucrate in fiecare zi peste norma zilei; in
// weekend si in zilele libere norma e 0; garzile si persoanele fara norma zilnica au „—”.
func TestOreSuplimentare(t *testing.T) {
	p := []ProgramPontaj{{Cod: "Z10", OraStart: 7 * 60, DurataMinute: 600}, {Cod: "Z6", OraStart: 7 * 60, DurataMinute: 360}, {Cod: "P8", OraStart: 7 * 60, DurataMinute: 480, PauzaMinute: 30}}
	opt, patru := int64(480), int64(240)
	c := ConfiguratiePontajInitiala()
	c.Norme = map[string]*int64{"a": &opt, "b": &opt, "g": &opt, "fara": nil}
	c.IstoricNorme = map[string][]SchimbareContract{"b": {{DeLa: "2026-10-15", NormaMinute: &patru}}}
	c.Garzi = map[string]bool{"g": true}
	in := []IntrarePontaj{
		{AngajatID: "a", Data: "2026-10-05", Valoare: "Z10"}, // luni: 10 h, norma 8 -> 2
		{AngajatID: "a", Data: "2026-10-06", Valoare: "Z6"},  // marti: 6 h, sub norma -> 0
		{AngajatID: "a", Data: "2026-10-10", Valoare: "Z6"},  // sambata: 6 h -> 6
		{AngajatID: "a", Data: "2026-10-14", Valoare: "8"},   // zi libera (luni-vineri), ore manuale: 8 -> 8
		{AngajatID: "b", Data: "2026-10-05", Valoare: "Z10"}, // inainte de schimbare: norma 8 -> 2
		{AngajatID: "b", Data: "2026-10-16", Valoare: "Z6"},  // dupa schimbare: norma 4 -> 2
		{AngajatID: "g", Data: "2026-10-05", Valoare: "Z10"},
		{AngajatID: "fara", Data: "2026-10-05", Valoare: "Z10"},
	}
	r, err := CalculeazaPontaj("2026-10-01", "2026-10-31", c, p, in, map[string]bool{"2026-10-14": true})
	if err != nil {
		t.Fatal(err)
	}
	val := map[string]ValoarePontaj{}
	for _, x := range r.Totaluri {
		val[x.AngajatID] = x.Valori["OS"]
	}
	for id, asteptat := range map[string]string{"a": "16", "b": "4"} {
		if v := TextValoarePontaj(val[id]); v != asteptat {
			t.Errorf("%s: OS=%s, asteptat %s", id, v, asteptat)
		}
	}
	if v := val["g"]; v.Definita || v.Motiv != MotivGardaFaraOS {
		t.Errorf("garda: %+v", v)
	}
	if v := val["fara"]; v.Definita || v.Motiv != MotivFaraNorma {
		t.Errorf("fara norma: %+v", v)
	}
	// O zi necalculabila (program cu pauza nestabilita) face OS necalculat.
	r, err = CalculeazaPontaj("2026-10-01", "2026-10-31", c, p, []IntrarePontaj{{AngajatID: "a", Data: "2026-10-05", Valoare: "P8"}}, nil)
	if err != nil {
		t.Fatal(err)
	}
	if v := r.Totaluri[0].Valori["OS"]; v.Definita {
		t.Errorf("zi necalculabila: %+v", v)
	}
}
