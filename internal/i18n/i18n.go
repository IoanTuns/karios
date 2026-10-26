// Package i18n ține textele aplicației în mai multe limbi. Cataloagele sunt
// fișiere JSON plate (cheie → text), unul pe limbă, numite după codul limbii
// (ro.json, en.json). Cele din locales/ sunt incluse în .exe; un folder de
// lângă .exe poate corecta texte sau adăuga limbi noi, fără recompilare.
// Aceleași fișiere vor fi folosite și de interfața React.
package i18n

import (
	"embed"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

type Limba string

const (
	Romana  Limba = "ro"
	Engleza Limba = "en"
)

// Implicita este limba folosită când o cheie lipsește din limba cerută.
const Implicita = Romana

//go:embed locales/*.json
var fisiere embed.FS

// cataloage se modifică doar la pornire (AdaugaFolder), înainte de orice T.
var cataloage = incarcaIncluse()

var reLimba = regexp.MustCompile(`^[a-z]{2}(-[A-Z]{2})?$`)

func incarcaIncluse() map[Limba]map[string]string {
	rez := map[Limba]map[string]string{}
	intrari, err := fisiere.ReadDir("locales")
	if err != nil {
		panic(err)
	}
	for _, intrare := range intrari {
		data, err := fisiere.ReadFile("locales/" + intrare.Name())
		if err != nil {
			panic(err)
		}
		if err := adauga(rez, intrare.Name(), data); err != nil {
			panic(err)
		}
	}
	return rez
}

// AdaugaFolder citește cataloagele *.json din dir peste cele incluse: o cheie
// existentă își schimbă textul, o limbă nouă se adaugă. Un folder lipsă nu e eroare.
func AdaugaFolder(dir string) error {
	nume, err := filepath.Glob(filepath.Join(dir, "*.json"))
	if err != nil {
		return err
	}
	for _, n := range nume {
		data, err := os.ReadFile(n)
		if err != nil {
			return err
		}
		if err := adauga(cataloage, filepath.Base(n), data); err != nil {
			return err
		}
	}
	return nil
}

func adauga(tinta map[Limba]map[string]string, numeFisier string, data []byte) error {
	l := Limba(strings.TrimSuffix(numeFisier, ".json"))
	if !reLimba.MatchString(string(l)) {
		return fmt.Errorf("%s: numele fișierului trebuie să fie un cod de limbă (de ex. ro.json, en.json)", numeFisier)
	}
	var c map[string]string
	if err := json.Unmarshal(data, &c); err != nil {
		return fmt.Errorf("%s: %w", numeFisier, err)
	}
	if tinta[l] == nil {
		tinta[l] = map[string]string{}
	}
	for k, v := range c {
		tinta[l][k] = v
	}
	return nil
}

// Limbi întoarce limbile disponibile, cu limba implicită prima.
func Limbi() []Limba {
	rez := []Limba{Implicita}
	for l := range cataloage {
		if l != Implicita {
			rez = append(rez, l)
		}
	}
	sort.Slice(rez[1:], func(i, j int) bool { return rez[1+i] < rez[1+j] })
	return rez
}

// T întoarce textul cheii în limba cerută, cu {{param}} înlocuiți. Dacă textul
// lipsește, folosește limba implicită, iar dacă lipsește și acolo, cheia.
func T(limba Limba, cheie string, param map[string]string) string {
	text, ok := cataloage[limba][cheie]
	if !ok || text == "" {
		text, ok = cataloage[Implicita][cheie]
	}
	if !ok {
		return cheie
	}
	for k, v := range param {
		text = strings.ReplaceAll(text, "{{"+k+"}}", v)
	}
	return text
}

// Catalog întoarce toate textele unei limbi, cu golurile completate din limba
// implicită. Parametrii {{nume}} rămân neînlocuiți; îi completează interfața.
func Catalog(limba Limba) map[string]string {
	rez := make(map[string]string, len(cataloage[Implicita]))
	for cheie := range cataloage[Implicita] {
		rez[cheie] = T(limba, cheie, nil)
	}
	return rez
}

// Lipsa întoarce cheile din limba implicită care nu au traducere în limba
// cerută, ca ecranul de setări să arate ce mai e de tradus.
func Lipsa(limba Limba) []string {
	var rez []string
	for cheie := range cataloage[Implicita] {
		if cataloage[limba][cheie] == "" {
			rez = append(rez, cheie)
		}
	}
	sort.Strings(rez)
	return rez
}

// Text este un text introdus de utilizator, cu câte o variantă pe limbă
// (denumirea unei categorii, a unui cod de absență, a unei sărbători).
// Se salvează în SQLite ca JSON: {"ro":"Noapte","en":"Night"}.
type Text map[Limba]string

// In întoarce varianta din limba cerută; dacă lipsește, pe cea din limba
// implicită; dacă lipsește și aceea, prima variantă completată.
func (t Text) In(l Limba) string {
	if s := t[l]; s != "" {
		return s
	}
	if s := t[Implicita]; s != "" {
		return s
	}
	limbi := make([]string, 0, len(t))
	for k, v := range t {
		if v != "" {
			limbi = append(limbi, string(k))
		}
	}
	if len(limbi) == 0 {
		return ""
	}
	sort.Strings(limbi)
	return t[Limba(limbi[0])]
}
