package rapoarte

import (
	"context"
	"strings"
	"testing"
)

func TestMarcajGardaContractIndependent(t *testing.T) {
	m := pregateste(t)
	fixture(t, m)
	if _, e := m.db.Exec(`INSERT INTO angajati_contracte(angajat_id,contract_id) VALUES('a','GARDA') ON CONFLICT(angajat_id) DO UPDATE SET contract_id='GARDA'`); e != nil {
		t.Fatal(e)
	}
	for _, l := range []string{"ro", "en"} {
		d, e := m.s.Pregateste(context.Background(), Cerere{ModelID: "raport-pontaj-detaliat", An: 2026, Luna: 10, Limba: l, AngajatID: "a"}, momentFix)
		if e != nil {
			t.Fatal(e)
		}
		gasit := false
		for i, c := range d.Coloane {
			if c.Cod == "angajat" {
				for _, r := range d.Randuri {
					if !r.Grup && strings.Contains(r.Valori[i], "(") {
						gasit = true
					}
				}
			}
		}
		if !gasit {
			t.Fatal("contract garda fara marcaj", l)
		}
	}
}
