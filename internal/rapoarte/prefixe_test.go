package rapoarte

import (
	"context"
	"testing"
)

// adaugaPrefix leaga angajatul de prefixul de pornire cu denumirea data (D089).
func adaugaPrefix(t *testing.T, m *mediu, angajat, prefix string) {
	t.Helper()
	if _, err := m.db.Exec(`INSERT INTO angajati_prefixe SELECT ?,id FROM prefixe WHERE json_extract(denumire,'$.ro')=?`, angajat, prefix); err != nil {
		t.Fatal(err)
	}
}

func TestPrefixInRaportulStandard(t *testing.T) {
	m := pregateste(t)
	fixture(t, m)
	adaugaPrefix(t, m, "a", "Dr.")
	doc, err := m.s.Pregateste(context.Background(), Cerere{ModelID: "raport-pontaj", An: 2026, Luna: 10, Limba: "ro", SectieID: "s", DepartamentID: "d"}, momentFix)
	if err != nil {
		t.Fatal(err)
	}
	idx := -1
	for i, c := range doc.Coloane {
		if c.Cod == "angajat" {
			idx = i
		}
	}
	if idx < 0 {
		t.Fatal("coloana angajat lipsa")
	}
	gasit := map[string]bool{}
	for _, r := range doc.Randuri {
		gasit[r.Valori[idx]] = true
	}
	// Doar angajatul cu prefix il primeste; ceilalti raman neschimbati.
	if !gasit["Dr. Șerban Ana"] || !gasit["Țincu Răzvan"] {
		t.Fatalf("nume in raport: %v", gasit)
	}
}

func TestPrefixInSQLRaportFlexibil(t *testing.T) {
	m := pregateste(t)
	fixture(t, m)
	adaugaPrefix(t, m, "a", "Biol.")
	d := modelFlexTest()
	d.Continut.Sursa = "SQL_LIBER"
	d.Continut.Flexibil.SQL = "SELECT angajat, prefix, nume FROM angajati WHERE angajat_id='a'"
	doc, err := m.s.PrevizualizeazaFlexibil(context.Background(), d, Cerere{An: 2026, Luna: 10, Limba: "ro"}, momentFix)
	if err != nil {
		t.Fatal(err)
	}
	if len(doc.Randuri) != 1 {
		t.Fatalf("randuri: %+v", doc.Randuri)
	}
	v := doc.Randuri[0].Valori
	if v[0] != "Biol. Șerban Ana" || v[1] != "Biol." || v[2] != "Șerban" {
		t.Fatalf("valori SQL: %v", v)
	}
}
