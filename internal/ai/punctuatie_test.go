package ai

import "testing"

// I118: numele urmat sau înconjurat de punctuație ajunge în contextul modelului
// și în selecția planului; cratima din mărci rămâne semnificativă.
func TestNumeCuPunctuatie(t *testing.T) {
	info := contextConcediu()
	for _, m := range []string{"Ce program are Pop Ana?", "Ce program are Pop Ana.", "Programul lui (Pop Ana)", "Pop Ana, te rog", "„Pop Ana”: program", "Ce program are M1?"} {
		c := ContextPentruCerere(m, info)
		if len(c.Angajati) != 1 || c.Angajati[0].Marca != "M1" {
			t.Errorf("%q: context %+v", m, c.Angajati)
		}
	}
	c := ContextPentruCerere("Ce program are DEMO-001?", info)
	if len(c.Angajati) != 1 || c.Angajati[0].Marca != "DEMO-001" {
		t.Errorf("marca cu cratima: %+v", c.Angajati)
	}
	for _, m := range []string{"Pop Ana: Z7.", "Z7 pentru Pop Ana."} {
		i, _ := InterpreteazaPlan(m, info)
		if i.Selectie.Tip != "persoane" || len(i.Selectie.Marci) != 1 || i.Selectie.Marci[0] != "M1" {
			t.Errorf("%q: selectie %+v", m, i.Selectie)
		}
	}
	// Punctuația singură nu se potrivește cu nimic.
	if contineExpresie("pop ana", "?") || contineExpresie("pop ana", "") {
		t.Error("expresie goala potrivita")
	}
}
