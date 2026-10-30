package diagnoza

import (
	"regexp"
	"strings"
)

// Jurnalul de diagnoză e necriptat, în karios/jurnale/, și poate fi trimis cuiva pentru depanare;
// de aceea nu trebuie să conțină date ale angajaților (D043, SEC-7). Valorile introduse ajung în
// mesajele de eroare de obicei între ghilimele; CNP-urile, telefoanele și emailurile se recunosc
// după formă. Tipul erorii, sursa și locul din cod rămân, ca jurnalul să fie util.
var (
	reGhilimele = regexp.MustCompile(`"[^"\n]*"|'[^'\n]*'|„[^”\n]*”|“[^”\n]*”|«[^»\n]*»`)
	reCifre     = regexp.MustCompile(`\d{6,}`)
	reEmail     = regexp.MustCompile(`[\w.+-]+@[\w-]+(\.[\w-]+)+`)
)

// LungimeMaximaIntrare limitează o intrare din jurnal (o stivă JS întreagă poate avea zeci de KB).
const LungimeMaximaIntrare = 4000

// Curata scoate din text valorile care pot fi date personale.
func Curata(text string) string {
	text = reEmail.ReplaceAllString(text, "@…")
	text = reGhilimele.ReplaceAllStringFunc(text, func(s string) string {
		r := []rune(s)
		return string(r[0]) + "…" + string(r[len(r)-1])
	})
	text = reCifre.ReplaceAllString(text, "#")
	if r := []rune(text); len(r) > LungimeMaximaIntrare {
		text = string(r[:LungimeMaximaIntrare]) + " […]"
	}
	return strings.TrimSpace(text)
}
