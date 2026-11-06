package calcul

import "testing"

func TestNormaDatorataReduceNumaiConcediileConfirmate(t *testing.T) {
	norma := int64(420)
	absente := []AbsentaNorma{{"2026-10-01", "CO", true}, {"2026-10-01", "CM", true}, {"2026-10-02", "CEV", true}, {"2026-10-05", "CFS", true}, {"2026-10-06", "CM", true}, {"2026-10-03", "CO", true}, {"2026-10-07", "CO", false}, {"2026-10-08", "FORM_PROF", true}, {"2026-10-09", "DELEG", true}, {"2026-10-12", "COMPENSATORIU", true}, {"2026-10-26", "CO", true}, {"2026-09-30", "CO", true}}
	r, err := CalculeazaNormaDatorata("2026-10-01", "2026-10-31", &norma, map[string]bool{"2026-10-26": true}, absente)
	if err != nil {
		t.Fatal(err)
	}
	if !r.Definita || r.ZileLucratoare != 21 || r.ZileConcediu != 4 || r.ZileDatorate != 17 || r.Minute != 7140 {
		t.Fatal(r)
	}
	partiala := int64(240)
	r, err = CalculeazaNormaDatorata("2026-10-01", "2026-10-31", &partiala, nil, nil)
	if err != nil || r.Minute != 5280 || r.ZileDatorate != 22 {
		t.Fatal(r, err)
	}
}
func TestNormaDatorataFaraNormaNuEsteZeroDefinit(t *testing.T) {
	r, err := CalculeazaNormaDatorata("2026-10-01", "2026-10-31", nil, nil, nil)
	if err != nil || r.Definita || r.ZileDatorate != 22 {
		t.Fatal(r, err)
	}
	norma := int64(420)
	r, err = CalculeazaNormaDatorata("2026-10-03", "2026-10-04", &norma, nil, nil)
	if err != nil || !r.Definita || r.Minute != 0 {
		t.Fatal(r, err)
	}
	for _, c := range [][2]string{{"gresit", "2026-10-31"}, {"2026-10-31", "2026-10-01"}} {
		if _, err := CalculeazaNormaDatorata(c[0], c[1], &norma, nil, nil); err == nil {
			t.Fatal("perioada invalida acceptata", c)
		}
	}
}
