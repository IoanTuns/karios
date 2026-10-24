package i18n

import (
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"sort"
	"strings"
	"testing"
)

var reParam = regexp.MustCompile(`\{\{(\w+)\}\}`)

func parametri(text string) []string {
	var rez []string
	for _, m := range reParam.FindAllStringSubmatch(text, -1) {
		rez = append(rez, m[1])
	}
	sort.Strings(rez)
	return rez
}

// cuFolder încarcă fișierele date într-un folder temporar și reface cataloagele
// incluse la sfârșitul testului.
func cuFolder(t *testing.T, fisiere map[string]string) error {
	t.Helper()
	t.Cleanup(func() { cataloage = incarcaIncluse() })
	dir := t.TempDir()
	for nume, continut := range fisiere {
		if err := os.WriteFile(filepath.Join(dir, nume), []byte(continut), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	return AdaugaFolder(dir)
}

// Limbile incluse în .exe trebuie să fie complete, cu aceiași parametri.
func TestCataloageIncluseComplete(t *testing.T) {
	incluse := []Limba{Romana, Engleza}
	for cheie, text := range cataloage[Implicita] {
		for _, l := range incluse {
			tradus := cataloage[l][cheie]
			if tradus == "" {
				t.Errorf("%s: lipsește traducerea pentru %q", l, cheie)
				continue
			}
			if !reflect.DeepEqual(parametri(text), parametri(tradus)) {
				t.Errorf("%s: cheia %q are parametrii %v, în %s are %v",
					l, cheie, parametri(tradus), Implicita, parametri(text))
			}
		}
	}
	for _, l := range incluse {
		for cheie := range cataloage[l] {
			if _, ok := cataloage[Implicita][cheie]; !ok {
				t.Errorf("%s: cheia %q nu există în %s", l, cheie, Implicita)
			}
		}
	}
}

func TestT(t *testing.T) {
	param := map[string]string{"tip": "SAMBATA", "de_la": "06:00", "pana_la": "07:00"}
	if got, want := T(Engleza, "reguli.gol", param), "SAMBATA: there is a gap between 06:00 and 07:00."; got != want {
		t.Errorf("en: %q, așteptat %q", got, want)
	}
	if got, want := T(Romana, "reguli.gol", param), "SAMBATA: există un gol între 06:00 și 07:00."; got != want {
		t.Errorf("ro: %q, așteptat %q", got, want)
	}
	if got := T("de", "tip_zi.SAMBATA", nil); got != "Sâmbătă" {
		t.Errorf("o limbă necunoscută ar trebui să folosească limba implicită, a dat %q", got)
	}
	if got := T(Engleza, "cheie.inexistenta", nil); got != "cheie.inexistenta" {
		t.Errorf("o cheie lipsă ar trebui să se întoarcă neschimbată, a dat %q", got)
	}
}

func TestAdaugaFolder(t *testing.T) {
	err := cuFolder(t, map[string]string{
		"en.json": `{"tip_zi.SARBATOARE": "Bank holiday"}`,
		"hu.json": `{"tip_zi.SAMBATA": "Szombat", "tip_zi.DUMINICA": ""}`,
	})
	if err != nil {
		t.Fatal(err)
	}

	if got := T(Engleza, "tip_zi.SARBATOARE", nil); got != "Bank holiday" {
		t.Errorf("textul corectat din folder: %q", got)
	}
	if got := T(Engleza, "tip_zi.SAMBATA", nil); got != "Saturday" {
		t.Errorf("celelalte texte incluse trebuie să rămână: %q", got)
	}
	if got := T("hu", "tip_zi.SAMBATA", nil); got != "Szombat" {
		t.Errorf("limba nouă: %q", got)
	}
	if got := T("hu", "tip_zi.DUMINICA", nil); got != "Duminică" {
		t.Errorf("un text gol trebuie să folosească limba implicită: %q", got)
	}
	if got, want := Limbi(), []Limba{"ro", "en", "hu"}; !reflect.DeepEqual(got, want) {
		t.Errorf("Limbi() = %v, așteptat %v", got, want)
	}

	lipsa := Lipsa("hu")
	if len(lipsa) != len(cataloage[Implicita])-1 {
		t.Errorf("hu are tradusă o singură cheie, Lipsa a dat %d chei", len(lipsa))
	}
	if Lipsa(Engleza) != nil {
		t.Errorf("engleza este completă, Lipsa a dat %v", Lipsa(Engleza))
	}
}

func TestAdaugaFolderErori(t *testing.T) {
	if err := cuFolder(t, map[string]string{"romana.json": `{}`}); err == nil {
		t.Error("aștepta eroare pentru un nume de fișier care nu e cod de limbă")
	}
	if err := cuFolder(t, map[string]string{"fr.json": `{"a": `}); err == nil {
		t.Error("aștepta eroare pentru JSON invalid")
	}
	if err := AdaugaFolder(filepath.Join(t.TempDir(), "nu-exista")); err != nil {
		t.Errorf("un folder lipsă nu e eroare: %v", err)
	}
}

func TestCatalog(t *testing.T) {
	err := cuFolder(t, map[string]string{"hu.json": `{"tip_zi.SAMBATA": "Szombat"}`})
	if err != nil {
		t.Fatal(err)
	}
	c := Catalog("hu")
	if len(c) != len(cataloage[Implicita]) {
		t.Errorf("catalogul are %d chei, așteptat %d", len(c), len(cataloage[Implicita]))
	}
	if c["tip_zi.SAMBATA"] != "Szombat" || c["tip_zi.DUMINICA"] != "Duminică" {
		t.Errorf("traducerea și rezerva din limba implicită: %q, %q", c["tip_zi.SAMBATA"], c["tip_zi.DUMINICA"])
	}
	if got := Catalog(Engleza)["reguli.gol"]; !strings.Contains(got, "{{de_la}}") {
		t.Errorf("parametrii trebuie să rămână pentru interfață: %q", got)
	}
}

func TestText(t *testing.T) {
	noapte := Text{Romana: "Noapte", Engleza: "Night"}
	if got := noapte.In(Engleza); got != "Night" {
		t.Errorf("en: %q", got)
	}
	if got := noapte.In("hu"); got != "Noapte" {
		t.Errorf("o limbă lipsă trebuie să folosească limba implicită: %q", got)
	}
	doarEn := Text{Romana: "", Engleza: "Night"}
	if got := doarEn.In(Romana); got != "Night" {
		t.Errorf("fără română, trebuie să folosească prima variantă completată: %q", got)
	}
	if got := (Text{}).In(Romana); got != "" {
		t.Errorf("un text gol trebuie să rămână gol: %q", got)
	}
}
