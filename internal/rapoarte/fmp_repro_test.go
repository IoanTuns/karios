package rapoarte

import (
	"context"
	"testing"

	"pontix/internal/datebaza"
)

// F08-T15: tipul FMP creat dupa salvarea configuratiei pontajului nu intra in CVE, are
// coloana lui si orele in OL; raportul detaliat arata FMP si CVE in locul coloanelor vechi.
func TestFMPDupaConfiguratieSalvata(t *testing.T) {
	for _, salvata := range []bool{false, true} {
		m := pregateste(t)
		fixture(t, m)
		ctx := context.Background()
		s := datebaza.Nou(m.db)
		if salvata {
			c, _, err := s.ConfiguratiePontaj(ctx)
			if err != nil {
				t.Fatal(err)
			}
			if err := s.SalveazaConfiguratiePontaj(ctx, c); err != nil {
				t.Fatal(err)
			}
		}
		for _, q := range []string{
			`INSERT INTO categorii_absente(cod,denumire,natura,efect_ore) VALUES('FORMARE','{"ro":"Formare"}','SPECIAL','NORMA')`,
			`INSERT INTO coduri_absenta(cod,denumire,activ,creat_la,modificat_la,categorie,efect_ore,categorie_id) VALUES('FMP','{"ro":"Formare profesionala"}',1,'x','x','SPECIAL','NORMA','FORMARE')`,
			`INSERT INTO absente_zile(angajat_id,data,cod,activ,creat_la,modificat_la,stare_aprobare,minute_creditate) VALUES('a','2026-10-27','FMP',1,'x','x','APROBATA',480)`,
		} {
			if _, err := m.db.Exec(q); err != nil {
				t.Fatal(err)
			}
		}
		c, _, err := s.ConfiguratiePontaj(ctx)
		if err != nil {
			t.Fatal(err)
		}
		reguli := map[string][]string{}
		surse := map[string]string{}
		for _, r := range c.Reguli {
			reguli[r.Coloana], surse[r.Coloana] = r.Coduri, r.Sursa
		}
		contine := func(l []string, cod string) bool {
			for _, x := range l {
				if x == cod {
					return true
				}
			}
			return false
		}
		if !contine(reguli["CVE"], "FMP") || surse["FMP"] != "ZILE_ABSENTA" || !contine(reguli["FMP"], "FMP") || !contine(reguli["OL"], "FMP") {
			t.Fatalf("salvata=%v: reguli %v %v", salvata, surse, reguli)
		}
		doc, err := m.s.Pregateste(ctx, Cerere{ModelID: "raport-pontaj-detaliat", An: 2026, Luna: 10, Limba: "ro"}, momentFix)
		if err != nil {
			t.Fatal(err)
		}
		valori, titluri := map[string]string{}, map[string]string{}
		for i, col := range doc.Coloane {
			valori[col.Cod], titluri[col.Cod] = doc.Randuri[0].Valori[i], col.Titlu
		}
		if valori["total_form_prof"] != "1" || titluri["total_form_prof"] != "FMP (Zile)" || valori["total_CEV"] != "0" || titluri["total_CEV"] != "CVE (Zile)" {
			t.Fatalf("salvata=%v: raport %v %v", salvata, valori, titluri)
		}
	}
}
