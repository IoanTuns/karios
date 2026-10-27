package main

import (
	"context"
	"crypto/rand"
	"path/filepath"
	"testing"

	"pontix/internal/bd"
)

func TestAppStareLuna(t *testing.T) {
	cheie := make([]byte, bd.LungimeCheie)
	if _, err := rand.Read(cheie); err != nil {
		t.Fatal(err)
	}
	db, err := bd.Deschide(filepath.Join(t.TempDir(), "test_stare.db"), cheie)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	if _, err := bd.Migreaza(context.Background(), db, bd.Migrari); err != nil {
		t.Fatal(err)
	}

	app := &App{db: db}

	// Citire stare lună
	rInit := app.CitesteStareLuna(2026, 10)
	if rInit.Eroare != nil {
		t.Fatalf("eroare citire stare: %+v", rInit.Eroare)
	}
	if rInit.Stare.Stare != "DESCHISA" {
		t.Errorf("așteptat DESCHISA, obținut: %s", rInit.Stare.Stare)
	}

	// Închidere lună
	rInchide := app.SeteazaStareLuna(2026, 10, "INCHISA", "Închis pentru salarii")
	if rInchide.Eroare != nil {
		t.Fatalf("eroare închidere: %+v", rInchide.Eroare)
	}
	if rInchide.Stare.Stare != "INCHISA" {
		t.Errorf("așteptat INCHISA, obținut: %s", rInchide.Stare.Stare)
	}

	// Redeschidere lună
	rRedeschide := app.SeteazaStareLuna(2026, 10, "DESCHISA", "Redeschis")
	if rRedeschide.Eroare != nil {
		t.Fatalf("eroare redeschidere: %+v", rRedeschide.Eroare)
	}
	if rRedeschide.Stare.Stare != "DESCHISA" {
		t.Errorf("așteptat DESCHISA, obținut: %s", rRedeschide.Stare.Stare)
	}
}
