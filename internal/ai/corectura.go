package ai

import (
	"strings"
	"unicode"
)

// Corectura e un cuvânt din mesaj înlocuit cu cel mai apropiat cuvânt cunoscut (F14-T15).
type Corectura struct {
	Scris    string `json:"scris"`
	Inteles  string `json:"inteles"`
	Distanta int    `json:"distanta"`
}

// vocabularComenzi sunt cuvintele pe care se sprijină recunoașterea comenzilor, fără diacritice.
// Zilele săptămânii sunt incluse ca să nu fie „corectate” în luni („marti” → „martie”).
var vocabularComenzi = []string{
	"ianuarie", "februarie", "martie", "aprilie", "mai", "iunie", "iulie", "august", "septembrie", "octombrie", "noiembrie", "decembrie",
	"luni", "marti", "miercuri", "joi", "vineri", "sambata", "duminica",
	"toti", "toate", "tot", "toata", "fiecare",
	"angajat", "angajati", "angajatii", "angajatilor", "angajatul", "salariat", "salariati", "salariatii", "salariatilor",
	"personal", "personalul", "colegi", "colegii", "persoane", "persoanele", "oameni", "oamenii", "lumea",
	"sectia", "sectie", "sectiei", "categoria", "categorie", "categoriei",
	"weekend", "weekendul", "weekenduri", "weekendurile", "lucratoare", "zilele", "inclusiv", "fiecare",
	"fara", "pontaj", "pontajul", "pontajului", "luna", "lunii", "listeaza", "schimburi", "schimburile",
	"concediu", "concediul", "concedii", "odihna", "medical", "boala", "absenta", "absente", "pana", "merge", "pleaca",
}

// pragGreseli: câte greșeli se acceptă după lungimea cuvântului (D-L): cuvintele scurte
// trebuie scrise exact, altfel „mai” ar putea deveni „mar”.
func pragGreseli(lungime int) int {
	switch {
	case lungime <= 4:
		return 0
	case lungime <= 7:
		return 1
	default:
		return 2
	}
}

// distantaDL e distanța Damerau-Levenshtein restrânsă (alinierea optimă): o literă lipsă,
// una în plus, una schimbată sau două vecine inversate costă câte 1.
func distantaDL(a, b string) int {
	x, y := []rune(a), []rune(b)
	d := make([][]int, len(x)+1)
	for i := range d {
		d[i] = make([]int, len(y)+1)
		d[i][0] = i
	}
	for j := range d[0] {
		d[0][j] = j
	}
	for i := 1; i <= len(x); i++ {
		for j := 1; j <= len(y); j++ {
			cost := 1
			if x[i-1] == y[j-1] {
				cost = 0
			}
			d[i][j] = min(d[i-1][j]+1, d[i][j-1]+1, d[i-1][j-1]+cost)
			if i > 1 && j > 1 && x[i-1] == y[j-2] && x[i-2] == y[j-1] {
				d[i][j] = min(d[i][j], d[i-2][j-2]+1)
			}
		}
	}
	return d[len(x)][len(y)]
}

// celMaiApropiat întoarce cuvântul din vocabular cel mai apropiat de c, în pragul de greșeli.
// La egalitate între două cuvinte diferite nu alege nimic, mai puțin când unul e o formă
// a celuilalt („personal” / „personalul”).
func celMaiApropiat(c string, vocabular map[string]bool) (string, int, bool) {
	prag := pragGreseli(len([]rune(c)))
	if prag == 0 {
		return "", 0, false
	}
	best, bestD, egal := "", prag+1, false
	for v := range vocabular {
		if abs(len([]rune(v))-len([]rune(c))) > prag {
			continue
		}
		d := distantaDL(c, v)
		switch {
		case d < bestD:
			best, bestD, egal = v, d, false
		case d == bestD && v != best && !strings.HasPrefix(v, best) && !strings.HasPrefix(best, v):
			egal = true
		case d == bestD && len(v) < len(best):
			best = v // dintre forme, cea mai scurtă, ca rezultatul să nu depindă de ordinea hărții
		}
	}
	if best == "" || bestD > prag || egal {
		return "", 0, false
	}
	return best, bestD, true
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

// vocabularPentru adaugă la cuvintele comenzilor cuvintele din numele secțiilor și categoriilor.
func vocabularPentru(info ContextPontix) map[string]bool {
	v := make(map[string]bool, len(vocabularComenzi)+16)
	for _, c := range vocabularComenzi {
		v[c] = true
	}
	for _, a := range info.Angajati {
		for _, c := range strings.Fields(normalizeaza(a.Sectie + " " + a.Categorie)) {
			v[c] = true
		}
	}
	return v
}

// corecteazaMesaj normalizează mesajul și înlocuiește cuvintele necunoscute cu cele mai apropiate
// cuvinte cunoscute. Cuvintele cu cifre (coduri de schimb, ore, date) nu se ating.
func corecteazaMesaj(mesaj string, info ContextPontix) (string, []Corectura) {
	vocabular := vocabularPentru(info)
	cuvinte := strings.Fields(normalizeaza(mesaj))
	var corecturi []Corectura
	for i, c := range cuvinte {
		// Punctuația de la capete rămâne pe loc: „octombie,” → „octombrie,”.
		inceput := strings.IndexFunc(c, func(r rune) bool { return unicode.IsLetter(r) })
		sfarsit := strings.LastIndexFunc(c, func(r rune) bool { return unicode.IsLetter(r) })
		if inceput < 0 {
			continue
		}
		baza := c[inceput : sfarsit+1]
		if vocabular[baza] || strings.IndexFunc(baza, func(r rune) bool { return !unicode.IsLetter(r) }) >= 0 {
			continue
		}
		if v, d, ok := celMaiApropiat(baza, vocabular); ok {
			cuvinte[i] = c[:inceput] + v + c[sfarsit+1:]
			corecturi = append(corecturi, Corectura{Scris: baza, Inteles: v, Distanta: d})
		}
	}
	return strings.Join(cuvinte, " "), corecturi
}
