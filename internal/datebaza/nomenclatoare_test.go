package datebaza

import (
	"context"
	"pontix/internal/i18n"
	"testing"
)

func elementTest(t *testing.T, m *mediuTest, tip, nume string) ElementNomenclator {
	t.Helper()
	lista := citeste(t, m.serviciu).Functii
	if tip == "categorii" {
		lista = citeste(t, m.serviciu).CategoriiSchimburi
	}
	for _, e := range lista {
		if numeElement(e) == nume {
			return e
		}
	}
	t.Fatalf("element lipsa: %s", nume)
	return ElementNomenclator{}
}
func TestCategoriiSchimburiConfigurabile(t *testing.T) {
	m := pregateste(t)
	ctx := context.Background()
	c := ElementNomenclator{Denumire: i18n.Text{"ro": "Ambulatoriu", "en": "Outpatient"}, Activ: true}
	if err := m.serviciu.SalveazaElementNomenclator(ctx, "categorii", c); err != nil {
		t.Fatal(err)
	}
	c = elementTest(t, m, "categorii", "Ambulatoriu")
	s := schimbTest()
	s.Categorie = c.ID
	if err := m.serviciu.SalveazaSchimb(ctx, s); err != nil {
		t.Fatal(err)
	}
	lista, err := m.serviciu.ListaSchimburi(ctx)
	if err != nil {
		t.Fatal(err)
	}
	var salvat Schimb
	for _, e := range lista {
		if e.Cod == "TEST_N12" {
			salvat = e
		}
	}
	if salvat.ID == "" || salvat.Categorie != c.ID {
		t.Fatal("program neasociat categoriei")
	}
	vechi := c
	c.Denumire["ro"] = "Ambulatoriu nou"
	if err := m.serviciu.SalveazaElementNomenclator(ctx, "categorii", c); err != nil {
		t.Fatal(err)
	}
	verificaCod(t, m.serviciu.SalveazaElementNomenclator(ctx, "categorii", vechi), "nomenclator.conflict")
	c = elementTest(t, m, "categorii", "Ambulatoriu nou")
	c.Activ = false
	if err := m.serviciu.SalveazaElementNomenclator(ctx, "categorii", c); err != nil {
		t.Fatal(err)
	}
	if err := m.serviciu.SalveazaSchimb(ctx, salvat); err != nil {
		t.Fatal("asocierea existenta nu poate fi pastrata", err)
	}
	s.Cod = "ALTNOU"
	verificaCod(t, m.serviciu.SalveazaSchimb(ctx, s), "nomenclator.inactiv")
	c.ID = ""
	c.Activ = true
	verificaCod(t, m.serviciu.SalveazaElementNomenclator(ctx, "categorii", c), "nomenclator.duplicat")
}
func TestFunctieAngajatPersistentaRedenumireSiDezactivare(t *testing.T) {
	m := pregateste(t)
	ctx := context.Background()
	f := ElementNomenclator{Denumire: i18n.Text{"ro": "Medic test", "en": "Test physician"}, Activ: true}
	if err := m.serviciu.SalveazaElementNomenclator(ctx, "functii", f); err != nil {
		t.Fatal(err)
	}
	f = elementTest(t, m, "functii", "Medic test")
	a := angajatTest("FUNC-01")
	a.FunctieID = f.ID
	if err := m.serviciu.SalveazaAngajat(ctx, a); err != nil {
		t.Fatal(err)
	}
	a = citeste(t, m.serviciu).Angajati[0]
	if a.FunctieID != f.ID || a.Functie != "Medic test" {
		t.Fatal("functie neasociata", a)
	}
	f.Denumire["ro"] = "Medic redenumit"
	if err := m.serviciu.SalveazaElementNomenclator(ctx, "functii", f); err != nil {
		t.Fatal(err)
	}
	a = citeste(t, m.serviciu).Angajati[0]
	if a.FunctieID != f.ID || a.Functie != "Medic redenumit" {
		t.Fatal("redenemire pierduta")
	}
	var nume string
	if err := m.db.QueryRow("SELECT functie FROM angajati WHERE id=?", a.ID).Scan(&nume); err != nil || nume != "Medic redenumit" {
		t.Fatal("textul folosit de rapoarte nu este sincronizat", err)
	}
	f = elementTest(t, m, "functii", "Medic redenumit")
	f.Activ = false
	if err := m.serviciu.SalveazaElementNomenclator(ctx, "functii", f); err != nil {
		t.Fatal(err)
	}
	if err := m.serviciu.SalveazaAngajat(ctx, a); err != nil {
		t.Fatal("asociere existenta refuzata", err)
	}
	b := angajatTest("FUNC-02")
	b.FunctieID = f.ID
	verificaCod(t, m.serviciu.SalveazaAngajat(ctx, b), "nomenclator.inactiv")
	a.FunctieID = ""
	a.Functie = ""
	if err := m.serviciu.SalveazaAngajat(ctx, a); err != nil {
		t.Fatal(err)
	}
	if citeste(t, m.serviciu).Angajati[0].FunctieID != "" {
		t.Fatal("asociere eliminata nepersistata")
	}
	var n int
	_ = m.db.QueryRow("SELECT count(*) FROM jurnal_modificari WHERE tabel='functii'").Scan(&n)
	if n != 3 {
		t.Fatalf("jurnal incomplet %d", n)
	}
}
