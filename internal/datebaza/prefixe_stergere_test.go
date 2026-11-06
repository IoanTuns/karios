package datebaza

import (
	"context"
	"testing"
)

// Un prefix inactiv nefolosit se poate sterge din administrare; unul folosit este protejat.
func TestPrefixStergereAdmin(t *testing.T) {
	m := pregateste(t)
	ctx := context.Background()
	dr, ec := elementTest(t, m, "prefixe", "Dr."), elementTest(t, m, "prefixe", "Ec.")
	a := angajatTest("P-20")
	a.PrefixID = dr.ID
	if err := m.serviciu.SalveazaAngajat(ctx, a); err != nil {
		t.Fatal(err)
	}
	for _, p := range []ElementNomenclator{dr, ec} {
		p.Activ = false
		if err := m.serviciu.SalveazaElementNomenclator(ctx, "prefixe", p); err != nil {
			t.Fatal(err)
		}
	}
	dr, ec = elementTest(t, m, "prefixe", "Dr."), elementTest(t, m, "prefixe", "Ec.")
	if err := m.serviciu.StergeConfigurare(ctx, "prefixe", dr.ID, dr.Revizie); err == nil {
		t.Fatal("prefix folosit sters")
	}
	if err := m.serviciu.StergeConfigurare(ctx, "prefixe", ec.ID, ec.Revizie); err != nil {
		t.Fatal("prefix nefolosit nesters", err)
	}
}
