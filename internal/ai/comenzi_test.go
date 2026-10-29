package ai

import (
	"encoding/json"
	"os"
	"testing"
)

// TestSetulDeComenziReale verifică ruta aleasă pentru comenzi reale (F14-T01, F14-T27).
// Fiecare eșec raportat de utilizator se adaugă în testdata/comenzi.json, ca să nu reapară.
func TestSetulDeComenziReale(t *testing.T) {
	b, err := os.ReadFile("testdata/comenzi.json")
	if err != nil {
		t.Fatal(err)
	}
	var comenzi []struct {
		Mesaj string `json:"mesaj"`
		Ruta  Ruta   `json:"ruta"`
		Sursa string `json:"sursa"`
	}
	if err := json.Unmarshal(b, &comenzi); err != nil {
		t.Fatal(err)
	}
	info := contextConcediu()
	info.CuPontaj = func(an, luna int) map[string]bool { return map[string]bool{} }
	for _, c := range comenzi {
		if d := Ruteaza(c.Mesaj, info); d.Ruta != c.Ruta {
			t.Errorf("%q: ruta %q, așteptat %q (%s)", c.Mesaj, d.Ruta, c.Ruta, c.Sursa)
		}
	}
	t.Logf("%d comenzi verificate", len(comenzi))
}
