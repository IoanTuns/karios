package ai

import (
	"strings"
	"testing"
)

func TestIntrebariDespreOZi(t *testing.T) {
	info := contextConcediu()
	info.Absente = []InfoAbsenta{{Nume: "Pop Ana", Marca: "M1", Data: "2026-10-15", Cod: "CO"}, {Nume: "Rus Ema", Marca: "M3", Data: "2026-10-16", Cod: "CM"}}
	info.TurePlanificate = []InfoPozitiePlanificare{{Nume: "Ion Dan", Marca: "M2", Data: "2026-10-15", CodSchimb: "G24"}}
	cazuri := []struct {
		mesaj    string
		contine  []string
		lipseste []string
	}{
		{"Cine lipseste pe 15 octombrie?", []string{"15 octombrie 2026", "Absenți (1)", "Pop Ana (M1) · CO"}, []string{"Programați", "Rus Ema"}},
		{"Cine lucreaza pe 15?", []string{"Programați (1)", "Ion Dan (M2) · G24"}, []string{"Absenți"}},
		{"Ce avem pe 15 octombrie?", []string{"Absenți (1)", "Programați (1)"}, nil},
		{"Cine lipseste pe 3 noiembrie?", []string{"doar din luna de lucru (octombrie 2026)"}, nil},
	}
	for _, c := range cazuri {
		r, ok := sintetizeazaZi(c.mesaj, info)
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
				t.Errorf("%q: nu trebuia %q", c.mesaj, s)
			}
		}
	}
	if _, ok := sintetizeazaZi("Pune-l pe Pop Ana pe Z7 pe 15 octombrie", info); ok {
		t.Error("comanda a fost tratată ca întrebare")
	}
}
