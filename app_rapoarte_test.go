package main

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/base64"
	"path/filepath"
	"testing"

	"pontix/internal/bd"
	"pontix/internal/rapoarte"
)

func TestRaportCacheTokenSiPDF(t *testing.T) {
	cheie := make([]byte, bd.LungimeCheie)
	if _, err := rand.Read(cheie); err != nil {
		t.Fatal(err)
	}
	db, err := bd.Deschide(filepath.Join(t.TempDir(), "rapoarte.db"), cheie)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	if _, err := bd.Migreaza(context.Background(), db, bd.Migrari); err != nil {
		t.Fatal(err)
	}
	a := NewApp(db)
	modele := a.ListaModeleRapoarte()
	if modele.Eroare != nil || len(modele.Configuratie.Modele) != 4 {
		t.Fatalf("modele: %+v", modele)
	}
	prima := a.PregatesteRaport(rapoarte.Cerere{ModelID: "raport-pontaj", An: 2026, Luna: 10, Limba: "ro"})
	if prima.Eroare != nil || prima.Document == nil || prima.Token == "" {
		t.Fatalf("generare: %+v", prima.Eroare)
	}
	if prima.Pagini != 1 {
		t.Fatalf("numar pagini in interfata: %d", prima.Pagini)
	}
	pdf, err := base64.StdEncoding.DecodeString(prima.PDF)
	if err != nil || !bytes.HasPrefix(pdf, []byte("%PDF-")) {
		t.Fatalf("PDF invalid: %v", err)
	}
	pastrat, ok := a.raportPentru(prima.Token)
	if !ok || !bytes.Equal(pdf, pastrat.PDF) {
		t.Fatal("fișierul salvat ar diferi de previzualizare")
	}
	for i := 0; i < 4; i++ {
		r := a.PregatesteRaport(rapoarte.Cerere{ModelID: "raport-pontaj", An: 2026, Luna: 10, Limba: "en"})
		if r.Eroare != nil {
			t.Fatal(r.Eroare)
		}
	}
	if _, ok := a.raportPentru(prima.Token); ok {
		t.Fatal("cache fără limită")
	}
	// Tokenul invalid este refuzat înainte de orice dialog, acces pe disc sau imprimare.
	for _, r := range []RaspunsFisierRaport{a.SalveazaPDFRaport(prima.Token), a.DeschidePDFPentruTiparire("../../extern.pdf")} {
		if r.Eroare == nil || r.Eroare.Cod != "rapoarte.expirat" {
			t.Fatalf("token necontrolat: %+v", r)
		}
	}
}
