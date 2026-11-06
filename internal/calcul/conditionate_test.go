package calcul

import (
	"errors"
	"testing"
)

func TestRegulaConditionataTotal200DoarOreEligibile(t *testing.T) {
	c := ConfiguratiePontajInitiala()
	c.Conditionate = []RegulaConditionata{{Activ: true, GrupID: "medici", ContractID: "GARDA", Garda: "DA", Regula: RegulaColoana{Coloana: "ore_garda", Sursa: "ORE_SCHIMBURI_ZNL", Coduri: []string{"N9"}, FactorSutimi: 200}}}
	c.Profile = map[string]ProfilPontaj{"m": {GrupID: "medici", ContractID: "GARDA", Garda: true}, "a": {GrupID: "asistenti", ContractID: "GARDA", Garda: true}, "n": {GrupID: "medici", ContractID: "INTREAGA"}}
	intrari := []IntrarePontaj{}
	for _, id := range []string{"m", "a", "n"} {
		intrari = append(intrari, IntrarePontaj{AngajatID: id, Data: "2026-10-02", Valoare: "N9"})
	}
	calc := func(libere map[string]bool) RezultatPontaj {
		r, e := CalculeazaPontaj("2026-10-01", "2026-10-31", c, []ProgramPontaj{{Cod: "N9", OraStart: 1320, DurataMinute: 540}}, intrari, libere)
		if e != nil {
			t.Fatal(e)
		}
		return r
	}
	for _, r := range calc(nil).Totaluri {
		if TextValoarePontaj(r.Valori["T"]) != "9" || TextValoarePontaj(r.Valori["ZNL"]) != "7" {
			t.Fatal(r)
		}
		v := r.Valori["ore_garda"]
		if r.AngajatID == "m" {
			if !v.Definita || TextValoarePontaj(v) != "14" {
				t.Fatal("7 ore eligibile x 2 total, fara 300%", r)
			}
		} else if v.Definita {
			t.Fatal("regula aplicata altui profil", r)
		}
	}
	for _, r := range calc(map[string]bool{"2026-10-02": true, "2026-10-03": true}).Totaluri {
		if r.AngajatID == "m" && TextValoarePontaj(r.Valori["ore_garda"]) != "18" {
			t.Fatal("suprapunere weekend si sarbatoare", r)
		}
	}
	c.Conditionate = append(c.Conditionate, c.Conditionate[0])
	_, e := CalculeazaPontaj("2026-10-01", "2026-10-31", c, nil, intrari, nil)
	var eroare *Eroare
	if !errors.As(e, &eroare) || eroare.Cod != "pontaj.reguli_suprapuse" {
		t.Fatal("reguli ambigue", e)
	}
}
