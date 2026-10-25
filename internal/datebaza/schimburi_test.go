package datebaza

import (
	"context"
	"strings"
	"testing"

	"pontix/internal/i18n"
)

func schimbTest() Schimb {
	// Valori exclusiv de test, nu intervale implicite ale aplicației.
	return Schimb{Cod: " test_n12 ", Denumire: i18n.Text{i18n.Romana: " Noapte de test ", i18n.Engleza: "Test night"}, OraStart: 19 * 60, DurataMinute: 720, PauzaMinute: 30, Tip: "NOAPTE", Activ: true}
}

func TestNomenclatorSchimburiSiLegaturaPontaj(t *testing.T) {
	m := pregateste(t)
	ctx := context.Background()
	if err := m.serviciu.SalveazaSchimb(ctx, schimbTest()); err != nil {
		t.Fatal(err)
	}
	lista, err := m.serviciu.ListaSchimburi(ctx)
	if err != nil || len(lista) != 12 {
		t.Fatalf("lista: %+v, %v", lista, err)
	}
	var s Schimb
	for _, schimb := range lista {
		if schimb.Cod == "TEST_N12" {
			s = schimb
		}
	}
	if s.Cod != "TEST_N12" || s.OraStart != 1140 || s.DurataMinute != 720 || s.PauzaMinute != 30 || s.Categorie != "CUSTOM" || s.Denumire.In(i18n.Romana) != "Noapte de test" {
		t.Fatalf("schimb: %+v", s)
	}
	verificaCod(t, m.serviciu.SalveazaSchimb(ctx, schimbTest()), "schimburi.cod_duplicat")
	s.OraStart = 20 * 60
	if err := m.serviciu.SalveazaSchimb(ctx, s); err != nil {
		t.Fatal(err)
	}
	if err := m.serviciu.SalveazaAngajat(ctx, angajatTest("test")); err != nil {
		t.Fatal(err)
	}
	a := citeste(t, m.serviciu).Angajati[0]
	if _, err := m.db.Exec(`INSERT INTO pontaj (id, angajat_id, data, start, sfarsit, tip_schimb_id, creat_la, modificat_la)
		VALUES ('test', ?, '2026-10-25', '2026-10-25T18:00:00Z', '2026-10-26T06:00:00Z', ?, 'x', 'x')`, a.ID, s.ID); err != nil {
		t.Fatal(err)
	}
	if err := m.serviciu.SchimbaActivSchimb(ctx, s.ID, false); err != nil {
		t.Fatal(err)
	}
	lista, err = m.serviciu.ListaSchimburi(ctx)
	var actual Schimb
	for _, schimb := range lista {
		if schimb.ID == s.ID {
			actual = schimb
		}
	}
	if err != nil || actual.Activ || actual.OraStart != 1200 {
		t.Fatalf("dezactivare/editare: %+v, %v", lista, err)
	}
	var id string
	if err := m.db.QueryRow(`SELECT tip_schimb_id FROM pontaj WHERE id = 'test'`).Scan(&id); err != nil || id != s.ID {
		t.Fatalf("legătura pontaj: %s, %v", id, err)
	}
	verificaCod(t, m.serviciu.SalveazaSchimb(ctx, schimbTest()), "schimburi.cod_duplicat")
	if err := m.serviciu.SchimbaActivSchimb(ctx, s.ID, true); err != nil {
		t.Fatal(err)
	}
}

func TestValidareSchimburi(t *testing.T) {
	m := pregateste(t)
	cazuri := []struct {
		cod      string
		modifica func(*Schimb)
	}{
		{"schimburi.cod_invalid", func(s *Schimb) { s.Cod = " " }},
		{"schimburi.cod_invalid", func(s *Schimb) { s.Cod = "N 12" }},
		{"schimburi.cod_invalid", func(s *Schimb) { s.Cod = strings.Repeat("X", 17) }},
		{"sarbatori.denumire_obligatorie", func(s *Schimb) { s.Denumire = nil }},
		{"schimburi.interval_invalid", func(s *Schimb) { s.OraStart = 1440 }},
		{"schimburi.interval_invalid", func(s *Schimb) { s.OraStart = -1 }},
		{"schimburi.interval_invalid", func(s *Schimb) { s.DurataMinute = 0 }},
		{"schimburi.interval_invalid", func(s *Schimb) { s.DurataMinute = 2881 }},
		{"schimburi.interval_invalid", func(s *Schimb) { s.PauzaMinute = 721 }},
		{"schimburi.tip_invalid", func(s *Schimb) { s.Tip = "ALTUL" }},
		{"schimburi.categorie_invalida", func(s *Schimb) { s.Categorie = "ALTCEVA" }},
		{"schimburi.inexistent", func(s *Schimb) { s.ID = "inexistent" }},
	}
	for i, c := range cazuri {
		t.Run(c.cod+string(rune('A'+i)), func(t *testing.T) {
			s := schimbTest()
			c.modifica(&s)
			verificaCod(t, m.serviciu.SalveazaSchimb(context.Background(), s), c.cod)
		})
	}
	lista, err := m.serviciu.ListaSchimburi(context.Background())
	if err != nil || len(lista) != 11 {
		t.Fatalf("salvare parțială: %+v, %v", lista, err)
	}
}
