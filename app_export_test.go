package main

import (
	"context"
	"crypto/rand"
	"path/filepath"
	"testing"

	"pontix/internal/bd"
	"pontix/internal/datebaza"
)

func TestAppExportaPontajExcel(t *testing.T) {
	cheie := make([]byte, bd.LungimeCheie)
	if _, err := rand.Read(cheie); err != nil {
		t.Fatal(err)
	}
	db, err := bd.Deschide(filepath.Join(t.TempDir(), "test_export.db"), cheie)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	if _, err := bd.Migreaza(context.Background(), db, bd.Migrari); err != nil {
		t.Fatal(err)
	}

	app := &App{db: db}
	rDemo := app.IncarcaDemo()
	if rDemo.Eroare != nil {
		t.Fatalf("eroare la incarcare demo: %+v", rDemo.Eroare)
	}

	srv := datebaza.Nou(db)
	data, _ := srv.Citeste(context.Background())
	if len(data.Angajati) == 0 {
		t.Fatal("nu există angajați în demo")
	}
	angID := data.Angajati[0].ID

	err = srv.SalveazaSchita(context.Background(), 2026, 10, []datebaza.CelulaSchita{
		{AngajatID: angID, Data: "2026-10-01", Valoare: "Z7"},
	})
	if err != nil {
		t.Fatalf("eroare salvare schiță: %v", err)
	}

	resp := app.ExportaPontajExcel(2026, 10)
	if resp.Eroare != nil {
		t.Fatalf("eroare la export: %+v", resp.Eroare)
	}
	if resp.Continut == "" {
		t.Fatal("conținutul base64 este gol")
	}
	if resp.NumeFisier != "Pontix_Pontaj_2026_10.xlsx" {
		t.Errorf("nume fișier neașteptat: %s", resp.NumeFisier)
	}
}
