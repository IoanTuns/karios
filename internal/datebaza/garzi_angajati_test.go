package datebaza

import (
	"context"
	"testing"
)

func TestMarcaGardaLegaturaSiProtectii(t *testing.T) {
	m := pregateste(t)
	ctx := context.Background()
	for _, marca := range []string{"N-01", "G-01", "G-02"} {
		if err := m.serviciu.SalveazaAngajat(ctx, angajatTest(marca)); err != nil {
			t.Fatal(err)
		}
	}
	lista := citeste(t, m.serviciu).Angajati
	dupaMarca := map[string]Angajat{}
	for _, a := range lista {
		dupaMarca[a.Marca] = a
	}
	n, g, h := dupaMarca["N-01"], dupaMarca["G-01"], dupaMarca["G-02"]
	g.AngajatPrincipalID = n.ID
	if err := m.serviciu.SalveazaAngajat(ctx, g); err != nil {
		t.Fatal(err)
	}
	for _, a := range citeste(t, m.serviciu).Angajati {
		if a.ID == g.ID && a.AngajatPrincipalID != n.ID {
			t.Fatal("legatura nepersistata")
		}
	}
	g.Nume = "Nume modificat"
	if err := m.serviciu.SalveazaAngajat(ctx, g); err != nil {
		t.Fatal("editare garda existenta", err)
	}
	h.AngajatPrincipalID = g.ID
	verificaCod(t, m.serviciu.SalveazaAngajat(ctx, h), "date.garda_legatura_invalida")
	n.AngajatPrincipalID = h.ID
	verificaCod(t, m.serviciu.SalveazaAngajat(ctx, n), "date.garda_legatura_invalida")
	g.AngajatPrincipalID = g.ID
	verificaCod(t, m.serviciu.SalveazaAngajat(ctx, g), "date.garda_legatura_invalida")
	g.AngajatPrincipalID = "inexistent"
	verificaCod(t, m.serviciu.SalveazaAngajat(ctx, g), "date.garda_legatura_invalida")
	if _, err := m.db.Exec("INSERT INTO angajati_garzi VALUES(?,?)", h.ID, g.ID); err == nil {
		t.Fatal("DB permite lant")
	}
	if _, err := m.db.Exec("UPDATE angajati_garzi SET principal_id=angajat_id WHERE angajat_id=?", g.ID); err == nil {
		t.Fatal("DB permite sine")
	}
	g.AngajatPrincipalID = ""
	if err := m.serviciu.SalveazaAngajat(ctx, g); err != nil {
		t.Fatal(err)
	}
	var nr int
	if err := m.db.QueryRow("SELECT count(*) FROM jurnal_modificari WHERE tabel='angajati_garzi'").Scan(&nr); err != nil || nr != 2 {
		t.Fatal("jurnal legaturi", nr, err)
	}
}
