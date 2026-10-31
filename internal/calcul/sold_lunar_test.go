package calcul

import "testing"

func TestSoldLunarSemnSiMinute(t *testing.T) {
	for _, c := range []struct{ inceput, lucrate, norma, final int64 }{{120, 9000, 9240, -120}, {-60, 9300, 9240, 0}, {15, 9240, 9240, 15}} {
		r := CalculeazaSoldLunar(c.inceput, c.lucrate, NormaDatorata{Minute: c.norma, Definita: true}, true)
		if !r.Definit || r.Final != c.final {
			t.Fatal(r, c)
		}
	}
	if CalculeazaSoldLunar(0, 100, NormaDatorata{}, true).Definit {
		t.Fatal("garda cu norma presupusa")
	}
}
