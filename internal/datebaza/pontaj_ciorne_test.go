package datebaza

import (
	"context"
	"fmt"
	"strings"
	"testing"
	"time"

	"pontix/internal/bd"
)

func TestSchitaPersistentaFaraInterpretare(t *testing.T) {
	m := pregateste(t)
	ctx := context.Background()
	if err := m.serviciu.SalveazaAngajat(ctx, angajatTest("test")); err != nil {
		t.Fatal(err)
	}
	id := citeste(t, m.serviciu).Angajati[0].ID
	celule := []CelulaSchita{{AngajatID: id, Data: "2026-10-01", Valoare: "7"}, {AngajatID: id, Data: "2026-10-02", Valoare: "N12"}, {AngajatID: id, Data: "2026-10-03", Valoare: "cev"}}
	if err := m.serviciu.SalveazaSchita(ctx, 2026, 10, celule); err != nil {
		t.Fatal(err)
	}
	if err := m.db.Close(); err != nil {
		t.Fatal(err)
	}
	var err error
	m.db, err = bd.Deschide(m.cale, m.cheie)
	if err != nil {
		t.Fatal(err)
	}
	m.serviciu = Nou(m.db)
	luna, err := m.serviciu.CitesteSchita(ctx, 2026, 10)
	if err != nil || len(luna.Celule) != 3 {
		t.Fatalf("schiță: %+v, %v", luna, err)
	}
	for i, c := range luna.Celule {
		if c.Valoare != celule[i].Valoare || c.Revizie == "" {
			t.Fatalf("valoare interpretată sau revizie lipsă: %+v", c)
		}
	}
	var n int
	if err := m.db.QueryRow(`SELECT count(*) FROM pontaj`).Scan(&n); err != nil || n != 0 {
		t.Fatalf("schița a devenit pontaj efectiv: %d, %v", n, err)
	}
	if err := m.serviciu.SchimbaActiv(ctx, id, false); err != nil {
		t.Fatal(err)
	}
	if luna, err := m.serviciu.CitesteSchita(ctx, 2026, 10); err != nil || len(luna.Celule) != 3 {
		t.Fatalf("schiță pierdută la dezactivare: %+v, %v", luna, err)
	}
	if alta, err := m.serviciu.CitesteSchita(ctx, 2026, 11); err != nil || len(alta.Celule) != 0 {
		t.Fatalf("luni amestecate: %+v, %v", alta, err)
	}
}

func TestSchitaConflicteSiLotAtomic(t *testing.T) {
	m := pregateste(t)
	ctx := context.Background()
	if err := m.serviciu.SalveazaAngajat(ctx, angajatTest("test")); err != nil {
		t.Fatal(err)
	}
	id := citeste(t, m.serviciu).Angajati[0].ID
	if err := m.serviciu.SalveazaSchita(ctx, 2026, 10, []CelulaSchita{{AngajatID: id, Data: "2026-10-01", Valoare: "7"}}); err != nil {
		t.Fatal(err)
	}
	luna, err := m.serviciu.CitesteSchita(ctx, 2026, 10)
	if err != nil {
		t.Fatal(err)
	}
	veche := luna.Celule[0]
	externa := veche
	externa.Valoare = "12"
	if err := Nou(m.db).SalveazaSchita(ctx, 2026, 10, []CelulaSchita{externa}); err != nil {
		t.Fatal(err)
	}
	veche.Valoare = "24"
	verificaCod(t, m.serviciu.SalveazaSchita(ctx, 2026, 10, []CelulaSchita{{AngajatID: id, Data: "2026-10-02", Valoare: "FP"}, veche}), "pontaj.conflict")
	luna, err = m.serviciu.CitesteSchita(ctx, 2026, 10)
	if err != nil || len(luna.Celule) != 1 || luna.Celule[0].Valoare != "12" {
		t.Fatalf("conflictul a suprascris sau a lăsat un lot parțial: %+v, %v", luna, err)
	}
	actuala := luna.Celule[0]
	actuala.Valoare = ""
	if err := m.serviciu.SalveazaSchita(ctx, 2026, 10, []CelulaSchita{actuala}); err != nil {
		t.Fatal(err)
	}
	luna, err = m.serviciu.CitesteSchita(ctx, 2026, 10)
	if err != nil || len(luna.Celule) != 1 || luna.Celule[0].Valoare != "" || luna.Celule[0].Revizie == actuala.Revizie {
		t.Fatalf("golirea celulei nu păstrează înregistrarea/revizia: %+v, %v", luna, err)
	}
	var n int
	if err := m.db.QueryRow(`SELECT count(*) FROM jurnal_modificari WHERE tabel = 'pontaj_ciorne'`).Scan(&n); err != nil || n != 3 {
		t.Fatalf("jurnalul lotului anulat: %d, %v", n, err)
	}
}

func TestValidareSchita(t *testing.T) {
	m := pregateste(t)
	ctx := context.Background()
	if err := m.serviciu.SalveazaAngajat(ctx, angajatTest("test")); err != nil {
		t.Fatal(err)
	}
	id := citeste(t, m.serviciu).Angajati[0].ID
	for _, c := range []CelulaSchita{{AngajatID: id, Data: "2026-11-01", Valoare: "7"}, {AngajatID: id, Data: "2026-10-32", Valoare: "7"}, {Data: "2026-10-01", Valoare: "7"}, {AngajatID: id, Data: "2026-10-01", Valoare: strings.Repeat("X", 17)}, {AngajatID: id, Data: "2026-10-01", Valoare: "a\tb"}} {
		verificaCod(t, m.serviciu.SalveazaSchita(ctx, 2026, 10, []CelulaSchita{c}), "pontaj.celula_invalida")
	}
	c := CelulaSchita{AngajatID: id, Data: "2026-10-01", Valoare: "7"}
	verificaCod(t, m.serviciu.SalveazaSchita(ctx, 2026, 10, []CelulaSchita{c, c}), "pontaj.lot_invalid")
	c.AngajatID = "inexistent"
	verificaCod(t, m.serviciu.SalveazaSchita(ctx, 2026, 10, []CelulaSchita{c}), "date.inexistent")
	verificaCod(t, m.serviciu.SalveazaSchita(ctx, 2026, 13, nil), "pontaj.perioada_invalida")
	if _, err := m.serviciu.CitesteSchita(ctx, 0, 1); err == nil {
		t.Fatal("an invalid acceptat")
	}
	if err := m.serviciu.SalveazaSchita(ctx, 2028, 2, []CelulaSchita{{AngajatID: id, Data: "2028-02-29", Valoare: "7"}}); err != nil {
		t.Fatal(err)
	}
}

func TestSchita50Pe31(t *testing.T) {
	m := pregateste(t)
	ctx := context.Background()
	tx, err := m.db.Begin()
	if err != nil {
		t.Fatal(err)
	}
	defer tx.Rollback()
	celule := make([]CelulaSchita, 0, 50*31)
	for r := 0; r < 50; r++ {
		id := fmt.Sprintf("test-%02d", r)
		if _, err := tx.Exec(`INSERT INTO angajati (id, nume, prenume, marca, creat_la, modificat_la) VALUES (?, 'Exemplu', 'Test', ?, 'x', 'x')`, id, id); err != nil {
			t.Fatal(err)
		}
		for zi := 1; zi <= 31; zi++ {
			celule = append(celule, CelulaSchita{AngajatID: id, Data: fmt.Sprintf("2026-10-%02d", zi), Valoare: "7"})
		}
	}
	if err := tx.Commit(); err != nil {
		t.Fatal(err)
	}
	inceput := time.Now()
	if err := m.serviciu.SalveazaSchita(ctx, 2026, 10, celule); err != nil {
		t.Fatal(err)
	}
	luna, err := m.serviciu.CitesteSchita(ctx, 2026, 10)
	if err != nil || len(luna.Celule) != 1550 {
		t.Fatalf("50×31: %d, %v", len(luna.Celule), err)
	}
	t.Logf("Salvare și citire 50×31 pe SQLite criptat: %s", time.Since(inceput))
}
