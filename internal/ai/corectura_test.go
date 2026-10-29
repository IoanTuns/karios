package ai

import (
	"strings"
	"testing"
)

func TestDistantaDL(t *testing.T) {
	cazuri := []struct {
		a, b string
		d    int
	}{
		{"octombie", "octombrie", 1},
		{"ocotmbrie", "octombrie", 1}, // două litere inversate
		{"angajati", "angajatii", 1},
		{"noiebrie", "noiembrie", 1},
		{"sectia", "sectia", 0},
		{"mai", "mar", 1},
	}
	for _, c := range cazuri {
		if d := distantaDL(c.a, c.b); d != c.d {
			t.Errorf("%s/%s: %d, așteptat %d", c.a, c.b, d, c.d)
		}
	}
}

func TestCorecturaMesajului(t *testing.T) {
	info := contextGrup()
	cazuri := []struct {
		mesaj, contine string
		corecturi      int
	}{
		{"Toti angajati fara pontaj in octombie lucreaza de la 7 la 14", "in octombrie lucreaza", 1},
		{"tot personalu in noiebrie pe Z7", "personal in noiembrie pe z7", 2},
		{"sectai ATI", "sectia ati", 1},
		{"pune-l pe Pop marti pe Z7", "marti", 0}, // ziua nu devine „martie”
		{"mai mult de doua garzi", "mai mult", 0}, // cuvintele scurte nu se corectează
		{"G42 in octombrie", "g42", 0},            // codurile cu cifre nu se ating
	}
	for _, c := range cazuri {
		m, corecturi := corecteazaMesaj(c.mesaj, info)
		if !strings.Contains(m, c.contine) || len(corecturi) != c.corecturi {
			t.Errorf("%q → %q, corecturi %+v", c.mesaj, m, corecturi)
		}
	}
}

func TestComandaGrupCuGreseli(t *testing.T) {
	cazuri := []struct {
		mesaj  string
		luna   string
		marci  int
		corect bool
	}{
		{"Toti angajati in octombie lucreaza de la 7 la 14", "2026-10-01", 3, true},
		{"toti angajatii in noiebrie pe Z7", "2026-11-01", 3, true},
		{"sectai ATI in ocotmbrie pe S7", "2026-10-01", 2, true},
		{"toti angajatii pe Z7 mai mult", "2026-10-01", 3, false}, // „mai” ca adverb: luna rămâne cea curentă
		{"toti angajatii in mai pe Z7", "2026-05-01", 3, false},
	}
	for _, c := range cazuri {
		r, ok := AnalizeazaComandaGrup(c.mesaj, contextGrup())
		if !ok || len(r.Propuneri) != 1 {
			t.Errorf("%q: nerecunoscută", c.mesaj)
			continue
		}
		p := r.Propuneri[0].Parametri
		corecturi := p["corecturi"].([]any)
		if p["data_inceput"] != c.luna || len(p["marci"].([]any)) != c.marci || (len(corecturi) > 0) != c.corect {
			t.Errorf("%q: %v %v corecturi %v", c.mesaj, p["data_inceput"], len(p["marci"].([]any)), corecturi)
		}
		if !strings.Contains(p["interpretare"].(string), " · ") {
			t.Errorf("%q: interpretare %q", c.mesaj, p["interpretare"])
		}
	}
}
