package stocare

import (
	"context"
	"io/fs"
	"path/filepath"
	"testing"
	"testing/fstest"
	"time"

	"pontix/internal/audit"
	"pontix/internal/bd"
	"pontix/internal/seif"
)

// cuMigrareNoua întoarce migrările reale plus o migrare de probă după ultima.
func cuMigrareNoua(t *testing.T) fs.FS {
	t.Helper()
	m := fstest.MapFS{}
	intrari, err := fs.ReadDir(bd.Migrari, ".")
	if err != nil {
		t.Fatal(err)
	}
	for _, i := range intrari {
		date, err := fs.ReadFile(bd.Migrari, i.Name())
		if err != nil {
			t.Fatal(err)
		}
		m[i.Name()] = &fstest.MapFile{Data: date}
	}
	m["0999_proba.sql"] = &fstest.MapFile{Data: []byte("-- +goose Up\nCREATE TABLE proba (x INTEGER) STRICT;\nINSERT INTO setari (cheie, valoare, modificat_la) VALUES ('din_migrare', 'da', 'x');\n")}
	return m
}

func modificaExtern(t *testing.T, dir, sql string) {
	t.Helper()
	s, err := seif.Deschide(filepath.Join(dir, NumeSeif))
	if err != nil {
		t.Fatal(err)
	}
	db, err := bd.Deschide(filepath.Join(dir, NumeBaza), s.CheieDate)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	if _, err := db.Exec(sql); err != nil {
		t.Fatal(err)
	}
}

func TestPornireaSigileazaSiDetecteaza(t *testing.T) {
	ctx := context.Background()
	dir := t.TempDir()
	db, paza, err := Deschide(ctx, dir)
	if err != nil {
		t.Fatal(err)
	}
	if r := paza.Rezultat(); r.Stare != audit.StareOK {
		t.Fatalf("prima pornire: %+v", r)
	}
	var n int
	if err := db.QueryRow(`SELECT count(*) FROM jurnal_modificari WHERE tabel = 'audit' AND rand_id = 'initializare'`).Scan(&n); err != nil || n != 1 {
		t.Fatalf("inițializarea în jurnal: %d, %v", n, err)
	}
	if _, err := db.Exec(`INSERT INTO setari (cheie, valoare, modificat_la) VALUES ('proba', 'da', 'x')`); err != nil {
		t.Fatal(err)
	}
	db.Close()

	db, paza, err = Deschide(ctx, dir)
	if err != nil {
		t.Fatal(err)
	}
	if paza.Blocata() {
		t.Fatalf("scrierea aplicației a produs incident: %+v", paza.Rezultat())
	}
	db.Close()

	modificaExtern(t, dir, `UPDATE setari SET valoare = 'nu' WHERE cheie = 'proba'`)
	db, paza, err = Deschide(ctx, dir)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	r := paza.Rezultat()
	if !paza.Blocata() || len(r.Tabele) != 1 || r.Tabele[0] != "setari" {
		t.Fatalf("modificarea externă: %+v", r)
	}
}

func TestMigrareaResigileazaDoarBazaIntacta(t *testing.T) {
	ctx := context.Background()
	acum := func() time.Time { return time.Date(2026, 10, 28, 8, 0, 0, 0, time.UTC) }

	dir := t.TempDir()
	db, _, err := deschide(ctx, dir, bd.Migrari, acum)
	if err != nil {
		t.Fatal(err)
	}
	db.Close()
	db, paza, err := deschide(ctx, dir, cuMigrareNoua(t), acum)
	if err != nil {
		t.Fatal(err)
	}
	if r := paza.Rezultat(); r.Stare != audit.StareOK {
		t.Fatalf("după migrare: %+v", r)
	}
	var n int
	if err := db.QueryRow(`SELECT count(*) FROM jurnal_modificari WHERE tabel = 'audit' AND rand_id = 'migrare'`).Scan(&n); err != nil || n != 1 {
		t.Fatalf("migrarea în jurnal: %d, %v", n, err)
	}
	db.Close()

	// O modificare externă urmată de o migrare rămâne incident.
	dir = t.TempDir()
	db, _, err = deschide(ctx, dir, bd.Migrari, acum)
	if err != nil {
		t.Fatal(err)
	}
	db.Close()
	modificaExtern(t, dir, `INSERT INTO setari (cheie, valoare, modificat_la) VALUES ('strecurat', 'da', 'x')`)
	db, paza, err = deschide(ctx, dir, cuMigrareNoua(t), acum)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	if !paza.Blocata() {
		t.Fatalf("migrarea a ascuns modificarea: %+v", paza.Rezultat())
	}
}
