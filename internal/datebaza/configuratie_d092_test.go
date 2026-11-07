package datebaza

import (
	"testing"

	"pontix/internal/calcul"
)

// D092: o configuratie salvata inainte de coloanele noi le primeste la citire, doar cu
// codurile existente, iar codurile redenumite se traduc daca vechiul cod a disparut.
func TestAdapteazaConfiguratiaD092(t *testing.T) {
	c := calcul.ConfiguratiePontajInitiala()
	vechi := calcul.ConfiguratiePontaj{Atribuire: c.Atribuire, Revizie: "3"}
	for _, r := range c.Reguli[:12] {
		if r.Coloana == "CEV" {
			r.Sursa, r.Coduri = "ZILE_ABSENTA", []string{"CEV", "CO"}
		}
		vechi.Reguli = append(vechi.Reguli, r)
	}
	optiuni := []OptiuneRegula{{Cod: "CO", Tip: "ABSENTA"}, {Cod: "CM", Tip: "ABSENTA"}, {Cod: "CVE", Tip: "ABSENTA"}, {Cod: "FMP", Tip: "ABSENTA"}, {Cod: "Z7", Tip: "SCHIMB"}}
	adapteazaConfiguratia(&vechi, optiuni)
	if err := calcul.ValideazaConfiguratiePontaj(vechi); err != nil {
		t.Fatalf("configuratia adaptata e invalida: %v", err)
	}
	reguli := map[string]calcul.RegulaColoana{}
	for _, r := range vechi.Reguli {
		reguli[r.Coloana] = r
	}
	if len(vechi.Reguli) != len(calcul.ColoanePontaj()) {
		t.Fatalf("coloane: %d", len(vechi.Reguli))
	}
	if r := reguli["CEV"]; r.Coduri[0] != "CVE" || r.Coduri[1] != "CO" {
		t.Errorf("codul vechi nu a fost tradus: %v", r.Coduri)
	}
	if r := reguli["CFP"]; r.Sursa != "NEDEFINITA" || len(r.Coduri) != 0 {
		t.Errorf("CFP fara cod existent trebuia neconfigurat: %+v", r)
	}
	if r := reguli["FMP"]; r.Sursa != "ZILE_ABSENTA" || len(r.Coduri) != 1 {
		t.Errorf("FMP: %+v", r)
	}
	if r := reguli["CVE"]; r.Sursa != "ZILE_ABSENTA_ALTELE" || len(r.Coduri) != 3 { // CO, CM, FMP exista; CFP, DLG nu
		t.Errorf("CVE exclude doar codurile existente: %+v", r)
	}
	if r := reguli["OL"]; r.Sursa != "ORE_LUCRATE_CODURI" || len(r.Coduri) != 1 || r.Coduri[0] != "FMP" {
		t.Errorf("OL: %+v", r)
	}
	if reguli["TZLN"].Sursa != "ZILE_LUCRATOARE_NETE" || reguli["TOCN"].Sursa != "ORE_NORMA_LUNA" {
		t.Error("TZLN/TOCN lipsa")
	}
	// Daca vechiul cod inca exista, nu se traduce.
	cfg := calcul.ConfiguratiePontaj{Reguli: []calcul.RegulaColoana{{Coloana: "CEV", Sursa: "ZILE_ABSENTA", Coduri: []string{"CEV"}, FactorSutimi: 100}}}
	adapteazaConfiguratia(&cfg, append(optiuni, OptiuneRegula{Cod: "CEV", Tip: "ABSENTA"}))
	if cfg.Reguli[0].Coduri[0] != "CEV" {
		t.Errorf("codul existent a fost tradus: %v", cfg.Reguli[0].Coduri)
	}
}
