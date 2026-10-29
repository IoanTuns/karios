package ai

import (
	"strings"
	"testing"
)

func TestListariDinCod(t *testing.T) {
	info := contextGrup()
	cazuri := []struct {
		mesaj    string
		contine  []string
		lipseste []string
	}{
		{"Listeaza toti angajati", []string{"Angajați activi: 3", "Pop Ana (M1)", "Rus Ema (M3)", "**ATI** (2)"}, nil},
		{"Listează angajații din secția ATI", []string{"secția ATI: 2", "Ion Dan (M2)"}, []string{"Rus Ema"}},
		{"Câți angajați avem?", []string{"Angajați activi: 3"}, nil},
		{"Ce schimburi există?", []string{"Schimburi active: 3", "Z7 · Zi 07–14 · 07:00–14:00 · 7 ore"}, nil},
	}
	for _, c := range cazuri {
		r, ok := AnalizeazaInterogareSinteza(c.mesaj, info)
		if !ok {
			t.Errorf("%q: nerecunoscută", c.mesaj)
			continue
		}
		for _, s := range c.contine {
			if !strings.Contains(r, s) {
				t.Errorf("%q: lipsește %q în\n%s", c.mesaj, s, r)
			}
		}
		for _, s := range c.lipseste {
			if strings.Contains(r, s) {
				t.Errorf("%q: nu trebuia să conțină %q", c.mesaj, s)
			}
		}
	}
	// O comandă de planificare nu e listare.
	if _, ok := AnalizeazaInterogareSinteza("Pune-l pe Pop pe Z7 pe 15 octombrie", info); ok {
		t.Error("comanda de planificare a fost tratată ca listare")
	}
}
