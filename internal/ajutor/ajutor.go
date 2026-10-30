// Package ajutor pastreaza catalogul comun pentru panou, manual si asistent.
package ajutor

import (
	"embed"
	"encoding/json"
	"golang.org/x/text/unicode/norm"
	"pontix/internal/i18n"
	"pontix/internal/versiune"
	"sort"
	"strings"
	"unicode"
)

//go:embed catalog.json
var fisiere embed.FS

type Definitie struct {
	ID      string   `json:"id"`
	Ecran   string   `json:"ecran"`
	Capitol int      `json:"capitol"`
	Termeni []string `json:"termeni"`
}
type Articol struct {
	ID       string   `json:"id"`
	Ecran    string   `json:"ecran"`
	Capitol  int      `json:"capitol"`
	Versiune string   `json:"versiune"`
	Titlu    string   `json:"titlu"`
	Pasi     []string `json:"pasi"`
	Rezultat string   `json:"rezultat"`
	Probleme string   `json:"probleme"`
	Termeni  []string `json:"termeni"`
}

var catalog = incarca()

func incarca() []Definitie {
	b, err := fisiere.ReadFile("catalog.json")
	if err != nil {
		panic(err)
	}
	var d []Definitie
	if err = json.Unmarshal(b, &d); err != nil {
		panic(err)
	}
	return d
}
func Normalizat(text string) string {
	return strings.Map(func(r rune) rune {
		if unicode.Is(unicode.Mn, r) {
			return -1
		}
		return unicode.ToLower(r)
	}, norm.NFD.String(text))
}
func Lista(limba i18n.Limba) []Articol {
	rez := make([]Articol, 0, len(catalog))
	for _, d := range catalog {
		prefix := "ajutor.articol." + d.ID + "."
		a := Articol{ID: d.ID, Ecran: d.Ecran, Capitol: d.Capitol, Versiune: versiune.Aplicatie, Termeni: append([]string{}, d.Termeni...), Titlu: i18n.T(limba, prefix+"titlu", nil), Rezultat: i18n.T(limba, prefix+"rezultat", nil), Probleme: i18n.T(limba, prefix+"probleme", nil)}
		for _, pas := range []string{"pas1", "pas2", "pas3"} {
			a.Pasi = append(a.Pasi, i18n.T(limba, prefix+pas, nil))
		}
		rez = append(rez, a)
	}
	return rez
}
func Gaseste(id string, limba i18n.Limba) (Articol, bool) {
	for _, a := range Lista(limba) {
		if a.ID == id {
			return a, true
		}
	}
	return Articol{}, false
}

// Cauta nu interpreteaza comenzi si nu acceseaza date operationale.
func Cauta(text string, limba i18n.Limba) []Articol {
	cuvinte := strings.FieldsFunc(Normalizat(text), func(r rune) bool { return !unicode.IsLetter(r) && !unicode.IsDigit(r) })
	ignorate := map[string]bool{"cum": true, "fac": true, "face": true, "pot": true, "pentru": true, "unde": true, "care": true, "este": true, "how": true, "can": true, "the": true, "what": true, "where": true, "help": true, "ajutor": true, "ghid": true, "please": true, "manual": true, "utilizare": true}
	type potrivire struct {
		a    Articol
		scor int
	}
	gasite := []potrivire{}
	for _, a := range Lista(limba) {
		scor := 0
		termeni := " " + Normalizat(strings.Join(a.Termeni, " ")+" "+a.Titlu) + " "
		for _, c := range cuvinte {
			if len(c) < 3 || ignorate[c] {
				continue
			}
			if strings.Contains(termeni, c) {
				scor++
			}
		}
		if scor > 0 {
			gasite = append(gasite, potrivire{a, scor})
		}
	}
	sort.SliceStable(gasite, func(i, j int) bool { return gasite[i].scor > gasite[j].scor })
	rez := []Articol{}
	for i, p := range gasite {
		if i == 3 {
			break
		}
		rez = append(rez, p.a)
	}
	return rez
}
