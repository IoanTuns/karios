package i18n

import "testing"

// F03-T11 (D106): corecturile au prioritate, o limbă nouă apare în Limbi, iar ștergerea lor
// readuce textele incluse.
func TestCorecturiPesteCataloage(t *testing.T) {
	defer SeteazaCorecturi(nil)
	SeteazaCorecturi(map[Limba]map[string]string{
		Engleza: {"ui.salveaza": "Store"},
		"hu":    {"ui.salveaza": "Mentés"},
	})
	if got := T(Engleza, "ui.salveaza", nil); got != "Store" {
		t.Fatalf("corectura en: %q", got)
	}
	if got := Catalog(Engleza)["ui.salveaza"]; got != "Store" {
		t.Fatalf("catalog en: %q", got)
	}
	if got := T("hu", "ui.salveaza", nil); got != "Mentés" {
		t.Fatalf("limba noua: %q", got)
	}
	if got := T("hu", "ui.renunta", nil); got != TextInclus(Romana, "ui.renunta") {
		t.Fatalf("rezerva pentru limba noua: %q", got)
	}
	gasit := false
	for _, l := range Limbi() {
		gasit = gasit || l == "hu"
	}
	if !gasit || Limbi()[0] != Implicita {
		t.Fatal("Limbi", Limbi())
	}
	if n := len(Lipsa("hu")); n != len(cataloage[Implicita])-1 {
		t.Fatalf("Lipsa hu: %d", n)
	}
	SeteazaCorecturi(nil)
	if got := T(Engleza, "ui.salveaza", nil); got != TextInclus(Engleza, "ui.salveaza") {
		t.Fatalf("dupa stergere: %q", got)
	}
}
