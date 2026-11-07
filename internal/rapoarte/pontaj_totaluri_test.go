package rapoarte

import (
	"context"
	"testing"
)

// F08-T13 (D092): TZLN si TOCN sunt in SQL, in tabelul separat pontaj_totaluri, cu aceleasi
// valori ca total_zile_cf_contract si total_ore_cf_contract din lunar; tabelul lunar nu se schimba.
func TestPontajTotaluriInSQL(t *testing.T) {
	m := pregateste(t)
	fixture(t, m)
	ctx := context.Background()
	c := Cerere{An: 2026, Luna: 10, Limba: "ro"}
	d := modelFlexTest()
	d.Continut.Sursa = "SQL_LIBER"
	d.Continut.Flexibil.SQL = "SELECT p.angajat_id, p.tzln, p.tocn, l.total_zile_cf_contract, l.total_ore_cf_contract, p.ol, p.co, p.cve FROM pontaj_totaluri p JOIN lunar l ON l.angajat_id = p.angajat_id ORDER BY p.angajat_id"
	doc, err := m.s.PrevizualizeazaFlexibil(ctx, d, c, momentFix)
	if err != nil {
		t.Fatal(err)
	}
	if len(doc.Randuri) == 0 {
		t.Fatal("pontaj_totaluri fara randuri")
	}
	cuNorma := 0
	for _, r := range doc.Randuri {
		v := r.Valori
		if v[1] != v[3] || v[2] != v[4] {
			t.Errorf("%s: tzln/tocn %q/%q difera de lunar %q/%q", v[0], v[1], v[2], v[3], v[4])
		}
		if v[2] != "" && v[2] != "—" {
			cuNorma++
		}
	}
	if cuNorma == 0 {
		t.Errorf("niciun angajat cu TOCN: %v", doc.Randuri)
	}
	// Tabelul lunar pastreaza exact coloanele de dinainte (SELECT * din rapoartele existente).
	d.Continut.Flexibil.SQL = "SELECT * FROM lunar"
	lunar, err := m.s.PrevizualizeazaFlexibil(ctx, d, c, momentFix)
	if err != nil {
		t.Fatal(err)
	}
	for _, col := range lunar.Coloane {
		switch col.Cod {
		case "tzln", "tocn", "ol":
			t.Errorf("coloana %s a aparut in lunar", col.Cod)
		}
	}
	gasit := false
	for _, s := range doc.Schema {
		gasit = gasit || s.Cod == "pontaj_totaluri"
	}
	if !gasit {
		t.Error("pontaj_totaluri lipseste din schema editorului")
	}
}
