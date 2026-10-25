package stocare

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"testing"

	"pontix/internal/seif"
)

func TestPrimaPornireSiRedeschidere(t *testing.T) {
	ctx := context.Background()
	dir := t.TempDir()

	db, err := Deschide(ctx, dir)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := db.Exec(`INSERT INTO setari (cheie, valoare, modificat_la) VALUES ('proba', 'da', 'x')`); err != nil {
		t.Fatal(err)
	}
	db.Close()
	for _, f := range []string{NumeSeif, NumeBaza} {
		if _, err := os.Stat(filepath.Join(dir, f)); err != nil {
			t.Errorf("%s nu a fost creat: %v", f, err)
		}
	}

	// A doua pornire găsește aceeași cheie și aceleași date.
	db, err = Deschide(ctx, dir)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	var v string
	if err := db.QueryRow(`SELECT valoare FROM setari WHERE cheie = 'proba'`).Scan(&v); err != nil || v != "da" {
		t.Errorf("datele după redeschidere: %q, %v", v, err)
	}
}

// Dacă seiful lipsește, dar baza există, nu se creează un seif nou.
func TestBazaFaraSeif(t *testing.T) {
	ctx := context.Background()
	dir := t.TempDir()
	db, err := Deschide(ctx, dir)
	if err != nil {
		t.Fatal(err)
	}
	db.Close()
	if err := os.Remove(filepath.Join(dir, NumeSeif)); err != nil {
		t.Fatal(err)
	}

	if _, err := Deschide(ctx, dir); !errors.Is(err, seif.ErrDeteriorat) {
		t.Fatalf("aștepta seif.ErrDeteriorat, a primit %v", err)
	}
	if _, err := os.Stat(filepath.Join(dir, NumeSeif)); !os.IsNotExist(err) {
		t.Error("s-a creat un seif nou peste o bază existentă")
	}
}
