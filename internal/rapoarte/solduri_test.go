package rapoarte

import (
	"context"
	"testing"
)

func TestDetaliatSoldInitialDupaMarca(t *testing.T) {
	m := pregateste(t)
	fixture(t, m)
	if _, e := m.db.Exec(`INSERT INTO initializari_pontaj(angajat_id,luna,minute,co_zile,modificat_la) VALUES('a','2026-10',90,0,'test')`); e != nil {
		t.Fatal(e)
	}
	d, e := m.s.Pregateste(context.Background(), Cerere{ModelID: "raport-pontaj-detaliat", An: 2026, Luna: 10, Limba: "ro"}, momentFix)
	if e != nil {
		t.Fatal(e)
	}
	gasit := false
	for i, c := range d.Coloane {
		if c.Cod == "marca" {
			if i+1 >= len(d.Coloane) || d.Coloane[i+1].Cod != "total_L_INCEPUT" {
				t.Fatal(d.Coloane)
			}
			if d.Randuri[0].Valori[i+1] != "1.5" {
				t.Fatal(d.Randuri[0])
			}
			gasit = true
		}
		if c.Cod == "total_ore_lucrate" {
			t.Fatal("rubrica veche afisata")
		}
	}
	if !gasit {
		t.Fatal("marca lipsa")
	}
}
