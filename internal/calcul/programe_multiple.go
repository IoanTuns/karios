package calcul

import (
	"sort"
	"strings"
)

// CoduriPrograme expandeaza proiectia zilnica; + nu este permis in codurile individuale.
func CoduriPrograme(valoare string) []string {
	coduri := strings.Split(strings.ToUpper(strings.TrimSpace(valoare)), "+")
	for i := range coduri {
		coduri[i] = strings.TrimSpace(coduri[i])
	}
	return coduri
}

// ValoarePrograme produce aceeasi reprezentare indiferent de ordinea pozitiilor.
func ValoarePrograme(coduri []string) string {
	c := append([]string(nil), coduri...)
	sort.Strings(c)
	return strings.Join(c, "+")
}

func AceleasiPrograme(a, b string) bool {
	return ValoarePrograme(CoduriPrograme(a)) == ValoarePrograme(CoduriPrograme(b))
}
