package rapoarte

import "testing"

func TestCelulaZiFaraOre(t *testing.T) {
	cazuri := []struct {
		v, culoare, cod  string
		absenta, libera  bool
		asteptat, culAst string
	}{
		{"0", "ABSENTA", "CO", true, false, "CO", "ABSENTA"},
		{"", "ABSENTA", "CM", true, false, "CM", "ABSENTA"},
		{"0", "ABSENTA", "CO", true, true, "", CuloareZiLibera},
		{"0", "", "", false, true, "", CuloareZiLibera},
		{"8", "", "D8", false, true, "8", ""},
		{"4", "ABSENTA", "FMP", true, false, "4", "ABSENTA"},
		{"—", "", "", false, true, "—", ""},
		{"0", "", "", false, false, "0", ""},
	}
	for _, c := range cazuri {
		v, cul := celulaZiFaraOre(c.v, c.culoare, c.cod, c.absenta, c.libera)
		if v != c.asteptat || cul != c.culAst {
			t.Errorf("%+v: %q %q", c, v, cul)
		}
	}
}

// F07-T42 (D097): fundalul PDF accepta culorile cu nume si culoarea unui program.
func TestFundalCelulaPdfme(t *testing.T) {
	for c, ast := range map[string]string{"WEEKEND": "#E7EAE7", CuloareZiLibera: "#D2D8D3", "#F8D7DC": "#F8D7DC", "": "", "ROSU": "", "#12345": ""} {
		if f := fundalCelulaPdfme(c); f != ast {
			t.Errorf("%q: %q, asteptat %q", c, f, ast)
		}
	}
}
