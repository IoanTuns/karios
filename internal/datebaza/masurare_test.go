package datebaza

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"pontix/internal/audit"
	"pontix/internal/bd"
)

// TestMasurareBazaMare măsoară operațiile sensibile la volum pe o bază generată cu 300 de
// angajați și un an întreg de pontaj (F06-T10). Rulează doar la cerere:
//
//	PONTIX_MASURARE=1 go test -run MasurareBazaMare -v ./internal/datebaza/
func TestMasurareBazaMare(t *testing.T) {
	if os.Getenv("PONTIX_MASURARE") == "" {
		t.Skip("măsurare la cerere: PONTIX_MASURARE=1")
	}
	const angajati, an = 300, 2026
	m := pregateste(t)
	ctx := context.Background()

	start := time.Now()
	tx, err := m.db.BeginTx(ctx, nil)
	if err != nil {
		t.Fatal(err)
	}
	coduri := []string{"Z7", "Z12", "N12", "G24", "S7"}
	celule := 0
	for i := 0; i < angajati; i++ {
		id := fmt.Sprintf("ang-%03d", i)
		if _, err := tx.Exec(`INSERT INTO angajati (id, nume, prenume, marca, creat_la, modificat_la) VALUES (?, ?, 'Test', ?, 'x', 'x')`, id, fmt.Sprintf("Angajat%03d", i), fmt.Sprintf("M-%03d", i)); err != nil {
			t.Fatal(err)
		}
		for zi := time.Date(an, 1, 1, 0, 0, 0, 0, time.UTC); zi.Year() == an; zi = zi.AddDate(0, 0, 1) {
			if zi.Weekday() == time.Saturday || zi.Weekday() == time.Sunday {
				continue
			}
			if _, err := tx.Exec(`INSERT INTO pontaj_ciorne (angajat_id, data, valoare, modificat_la) VALUES (?, ?, ?, 'x')`, id, zi.Format(time.DateOnly), coduri[(i+zi.YearDay())%len(coduri)]); err != nil {
				t.Fatal(err)
			}
			celule++
		}
	}
	if err := tx.Commit(); err != nil {
		t.Fatal(err)
	}
	t.Logf("generare: %d angajați, %d celule de pontaj în %s", angajati, celule, time.Since(start).Round(time.Millisecond))

	masoara := func(nume string, f func() error) time.Duration {
		t.Helper()
		start := time.Now()
		if err := f(); err != nil {
			t.Fatalf("%s: %v", nume, err)
		}
		d := time.Since(start)
		t.Logf("%-45s %8s", nume, d.Round(time.Millisecond))
		return d
	}
	masoara("prima sigilare (audit.Initializeaza)", func() error { return audit.Initializeaza(ctx, m.db) })
	verificare := masoara("verificarea la pornire (audit.Verifica)", func() error {
		r, err := audit.Verifica(ctx, m.db)
		if err == nil && r.Stare != audit.StareOK {
			err = fmt.Errorf("stare %s", r.Stare)
		}
		return err
	})
	masoara("amprenta unei luni (audit.AmprentaLuna)", func() error { _, err := audit.AmprentaLuna(ctx, m.db, an, 10); return err })
	masoara("citirea pontajului unei luni (CitesteSchita)", func() error { _, err := m.serviciu.CitesteSchita(ctx, an, 10); return err })

	// Salvarea unei luni întregi, fără și cu trigger-ele auditului.
	luna := func(valoare string) []CelulaSchita {
		s, err := m.serviciu.CitesteSchita(ctx, an, 11)
		if err != nil {
			t.Fatal(err)
		}
		rez := make([]CelulaSchita, len(s.Celule))
		for i, c := range s.Celule {
			rez[i] = CelulaSchita{AngajatID: c.AngajatID, Data: c.Data, Valoare: valoare, Revizie: c.Revizie}
		}
		return rez
	}
	masoara("salvarea unei luni, fără audit", func() error { return m.serviciu.SalveazaSchita(ctx, an, 11, luna("Z7")) })
	// Starea sigilată trebuie refăcută după scrierea fără trigger-e.
	if err := audit.NoteazaEveniment(ctx, m.db, "masurare", nil, nil); err != nil {
		t.Fatal(err)
	}
	m.db.Close()
	paza := audit.NouaPaza(audit.Rezultat{Stare: audit.StareOK})
	db, err := bd.DeschideCuAudit(m.cale, m.cheie, paza)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	s := Nou(db)
	m.serviciu = s
	masoara("salvarea unei luni, cu audit (trigger-e)", func() error { return s.SalveazaSchita(ctx, an, 11, luna("N12")) })
	masoara("verificarea după salvare", func() error {
		r, err := audit.Verifica(ctx, db)
		if err == nil && r.Stare != audit.StareOK {
			err = fmt.Errorf("stare %s, tabele %v", r.Stare, r.Tabele)
		}
		return err
	})
	if verificare > 10*time.Second {
		t.Errorf("verificarea la pornire durează %s pentru %d celule", verificare, celule)
	}
}
