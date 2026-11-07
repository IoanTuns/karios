package i18n

import (
	"encoding/json"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

// R4 (analiza UX 0.92.59, F16-T37): textele romanesti au diacritice. Lista contine doar
// cuvinte care, scrise fara diacritice, sunt sigur gresite in romana (nu forme articulate
// valide ca „luna”, „norma”, „data”, nici „ca” sau „sa”). Variabilele {{…}}, codul intre
// backticks si adresele web se ignora.
var faraDiacriticeInterzise = []string{
	"si", "in", "fara", "daca", "dupa", "pana", "inca", "cand", "lucratoare", "inainte",
	"foloseste", "sarbatori", "sarbatorile", "garzi", "angajati", "angajatii", "angajatilor",
	"selectia", "selectie", "modificari", "modificarile", "absente", "absentele", "absentei",
	"inregistrare", "inregistrarea", "inregistrari", "inceput", "inseamna", "verificati",
	"configuratia", "interfata", "definitie", "sterge", "fisierul", "rand", "randul", "functie",
	"functia", "pozitie", "aplicatiei", "acelasi", "aceeasi", "ramane", "campuri", "obligatia",
	"intre", "numar", "numarul", "marci", "actiune", "actiuni",
}

func TestTexteleRomanestiAuDiacritice(t *testing.T) {
	b, err := os.ReadFile(filepath.Join("locales", "ro.json"))
	if err != nil {
		t.Fatal(err)
	}
	var texte map[string]string
	if err := json.Unmarshal(b, &texte); err != nil {
		t.Fatal(err)
	}
	protejat := regexp.MustCompile("\\{\\{[^}]*\\}\\}|`[^`]*`|\\{[a-z_]+\\}|https?://\\S+")
	interzise := map[string]bool{}
	for _, c := range faraDiacriticeInterzise {
		interzise[c] = true
	}
	cuvant := regexp.MustCompile(`[A-Za-zĂÂÎȘȚăâîșț]+`)
	verbEaza := regexp.MustCompile(`^[a-z]{2,}eaza$`)
	gresite := 0
	for cheie, text := range texte {
		for _, w := range cuvant.FindAllString(protejat.ReplaceAllString(text, " "), -1) {
			lw := strings.ToLower(w)
			if interzise[lw] || verbEaza.MatchString(lw) {
				gresite++
				if gresite <= 20 {
					t.Errorf("%s: „%s” fara diacritice in %q", cheie, w, text)
				}
			}
		}
	}
	if gresite > 20 {
		t.Errorf("inca %d cuvinte fara diacritice", gresite-20)
	}
}
