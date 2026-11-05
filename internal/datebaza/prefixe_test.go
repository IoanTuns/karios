package datebaza

import (
	"context"
	"testing"

	"pontix/internal/i18n"
)

func dupaMarca(t *testing.T, s *Serviciu) map[string]Angajat {
	t.Helper()
	rez := map[string]Angajat{}
	for _, a := range citeste(t, s).Angajati {
		rez[a.Marca] = a
	}
	return rez
}

func TestPrefixePredefinite(t *testing.T) {
	m := pregateste(t)
	nume := map[string]bool{}
	for _, p := range citeste(t, m.serviciu).Prefixe {
		if !p.Activ || p.Denumire["ro"] != p.Denumire["en"] {
			t.Fatalf("prefix de pornire neasteptat: %+v", p)
		}
		nume[p.Denumire["ro"]] = true
	}
	for _, p := range []string{"Dr.", "Biol.", "Chim.", "Farm.", "As.", "Ing.", "Ec.", "Psih."} {
		if !nume[p] {
			t.Errorf("lipseste prefixul %s", p)
		}
	}
}

func TestPrefixAngajatOptionalSiJurnal(t *testing.T) {
	m := pregateste(t)
	ctx := context.Background()
	dr := elementTest(t, m, "prefixe", "Dr.")
	// Fara prefix: salvarea merge (prefixul este optional).
	if err := m.serviciu.SalveazaAngajat(ctx, angajatTest("P-01")); err != nil {
		t.Fatal(err)
	}
	a := dupaMarca(t, m.serviciu)["P-01"]
	if a.PrefixID != "" || a.Prefix != nil {
		t.Fatalf("prefix neasteptat: %+v", a)
	}
	a.PrefixID = dr.ID
	if err := m.serviciu.SalveazaAngajat(ctx, a); err != nil {
		t.Fatal(err)
	}
	a = dupaMarca(t, m.serviciu)["P-01"]
	if a.PrefixID != dr.ID || a.Prefix["ro"] != "Dr." {
		t.Fatalf("prefix nepersistat: %+v", a)
	}
	var nume, prenume string
	if err := m.db.QueryRow("SELECT nume,prenume FROM angajati WHERE id=?", a.ID).Scan(&nume, &prenume); err != nil || nume != a.Nume || prenume != a.Prenume {
		t.Fatal("numele salvat nu trebuie sa se schimbe", err)
	}
	// Prefix inexistent refuzat.
	b := a
	b.PrefixID = "inexistent"
	verificaCod(t, m.serviciu.SalveazaAngajat(ctx, b), "nomenclator.invalid")
	// Eliminarea prefixului.
	a.PrefixID = ""
	if err := m.serviciu.SalveazaAngajat(ctx, a); err != nil {
		t.Fatal(err)
	}
	if dupaMarca(t, m.serviciu)["P-01"].PrefixID != "" {
		t.Fatal("prefix neeliminat")
	}
	var n int
	if err := m.db.QueryRow("SELECT count(*) FROM jurnal_modificari WHERE tabel='angajati_prefixe'").Scan(&n); err != nil || n != 2 {
		t.Fatalf("jurnal prefix: %d %v", n, err)
	}
}

func TestPrefixInactivRamaneLaAngajatiiExistenti(t *testing.T) {
	m := pregateste(t)
	ctx := context.Background()
	biol := elementTest(t, m, "prefixe", "Biol.")
	a := angajatTest("P-02")
	a.PrefixID = biol.ID
	if err := m.serviciu.SalveazaAngajat(ctx, a); err != nil {
		t.Fatal(err)
	}
	biol.Activ = false
	if err := m.serviciu.SalveazaElementNomenclator(ctx, "prefixe", biol); err != nil {
		t.Fatal(err)
	}
	a = dupaMarca(t, m.serviciu)["P-02"]
	a.Nume = "Nume modificat"
	if err := m.serviciu.SalveazaAngajat(ctx, a); err != nil {
		t.Fatal("prefixul inactiv existent trebuie pastrat", err)
	}
	b := angajatTest("P-03")
	b.PrefixID = biol.ID
	verificaCod(t, m.serviciu.SalveazaAngajat(ctx, b), "nomenclator.inactiv")
}

func TestPrefixNomenclatorEditabilSiFaraDuplicat(t *testing.T) {
	m := pregateste(t)
	ctx := context.Background()
	nou := ElementNomenclator{Denumire: i18n.Text{"ro": "Conf. dr.", "en": "Assoc. Prof."}, Activ: true}
	if err := m.serviciu.SalveazaElementNomenclator(ctx, "prefixe", nou); err != nil {
		t.Fatal(err)
	}
	if elementTest(t, m, "prefixe", "Conf. dr.").Denumire["en"] != "Assoc. Prof." {
		t.Fatal("forma engleza pierduta")
	}
	verificaCod(t, m.serviciu.SalveazaElementNomenclator(ctx, "prefixe", ElementNomenclator{Denumire: i18n.Text{"ro": "dr."}, Activ: true}), "nomenclator.duplicat")
}

func TestPrefixMarcaGardaAfisatDinPrincipal(t *testing.T) {
	m := pregateste(t)
	ctx := context.Background()
	dr := elementTest(t, m, "prefixe", "Dr.")
	n := angajatTest("N-10")
	n.PrefixID = dr.ID
	if err := m.serviciu.SalveazaAngajat(ctx, n); err != nil {
		t.Fatal(err)
	}
	if err := m.serviciu.SalveazaAngajat(ctx, angajatTest("G-10")); err != nil {
		t.Fatal(err)
	}
	l := dupaMarca(t, m.serviciu)
	g := l["G-10"]
	g.AngajatPrincipalID = l["N-10"].ID
	if err := m.serviciu.SalveazaAngajat(ctx, g); err != nil {
		t.Fatal(err)
	}
	g = dupaMarca(t, m.serviciu)["G-10"]
	if g.PrefixID != "" || g.Prefix["ro"] != "Dr." {
		t.Fatalf("garda trebuie sa afiseze prefixul principal fara sa il fixeze: %+v", g)
	}
	// O salvare neutra a marcii de garda nu fixeaza prefixul mostenit.
	if err := m.serviciu.SalveazaAngajat(ctx, g); err != nil {
		t.Fatal(err)
	}
	var c int
	if err := m.db.QueryRow("SELECT count(*) FROM angajati_prefixe WHERE angajat_id=?", g.ID).Scan(&c); err != nil || c != 0 {
		t.Fatalf("prefix fixat pe garda: %d %v", c, err)
	}
}
