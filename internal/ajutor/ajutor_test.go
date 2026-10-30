package ajutor

import (
	"os"
	"pontix/internal/i18n"
	"pontix/internal/versiune"
	"regexp"
	"strings"
	"testing"
)

func TestCatalogAcoperaNavigareaSiTraducerile(t *testing.T) {
	b, err := os.ReadFile("../../frontend/src/App.tsx")
	if err != nil {
		t.Fatal(err)
	}
	m := regexp.MustCompile("type Ecran = ([^;]+);").FindSubmatch(b)
	if len(m) != 2 {
		t.Fatal("schema navigarii lipseste")
	}
	ecrane := map[string]bool{}
	ids := map[string]bool{}
	for _, l := range []i18n.Limba{i18n.Romana, i18n.Engleza} {
		for _, a := range Lista(l) {
			ecrane[a.Ecran] = true
			if a.Versiune != versiune.Aplicatie || a.Capitol < 1 || a.Capitol > 20 || len(a.Pasi) != 3 {
				t.Fatal(a)
			}
			texte := append([]string{a.Titlu, a.Rezultat, a.Probleme}, a.Pasi...)
			for _, text := range texte {
				if text == "" || strings.HasPrefix(text, "ajutor.") {
					t.Fatal(l, a.ID, text)
				}
			}
		}
	}
	for _, d := range catalog {
		if ids[d.ID] || !regexp.MustCompile("^[a-z]+[.][a-z]+$").MatchString(d.ID) {
			t.Fatal(d.ID)
		}
		ids[d.ID] = true
	}
	for _, ecran := range strings.Split(string(m[1]), "|") {
		ecran = strings.Trim(strings.TrimSpace(ecran), "'")
		if !ecrane[ecran] {
			t.Fatal("ecran fara ajutor", ecran)
		}
	}
	if !ecrane["acces"] {
		t.Fatal("lipseste ajutorul pentru acces")
	}
}
func TestCautareAccenteSiLipsaAcoperire(t *testing.T) {
	for _, c := range []struct{ mesaj, id string }{{"Cum inregistrez un concediu?", "absente.inregistreaza"}, {"Cum înregistrez un concediu?", "absente.inregistreaza"}, {"How do I record leave?", "absente.inregistreaza"}, {"backup", "backup.restaureaza"}, {"S7 sincronizare", "pontaj.sincronizare"}} {
		r := Cauta(c.mesaj, i18n.Romana)
		if len(r) == 0 || r[0].ID != c.id {
			t.Fatal(c, r)
		}
	}
	if len(Cauta("telematicxyz", i18n.Romana)) != 0 {
		t.Fatal("rezultate inventate")
	}
	if _, ok := Gaseste("../../necunoscut", i18n.Romana); ok {
		t.Fatal("id invalid")
	}
}
