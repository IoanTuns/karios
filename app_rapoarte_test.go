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
	// Cinci modele standard si cinci modele predefinite SQL + pdfme (D090, migrarea 0033).
	if modele.Eroare != nil || len(modele.Configuratie.Modele) != 10 {
		t.Fatalf("modele: %+v", modele)
	}
	vizual := a.PregatesteRaport(rapoarte.Cerere{ModelID: "raport-pontaj-sql", An: 2026, Luna: 10, Limba: "ro"})
	if vizual.Eroare != nil || !vizual.GenerarePdfme || vizual.Sablon == "" || vizual.Intrari["tabel"] == "" || vizual.PDF != "" {
		t.Fatalf("model vizual: %+v", vizual.Eroare)
	}
	// Modelul vechi a fost convertit la primul acces (F09-T23): PDF-ul vine din interfata.
	if len(modele.Conversie) != 5 {
		t.Fatalf("conversie: %+v", modele.Conversie)
	}
	prima := a.PregatesteRaport(rapoarte.Cerere{ModelID: "raport-pontaj", An: 2026, Luna: 10, Limba: "ro"})
	if prima.Eroare != nil || prima.Document == nil || prima.Token == "" || !prima.GenerarePdfme {
		t.Fatalf("generare: %+v", prima.Eroare)
	}
	if r := a.SalveazaPDFRaport(prima.Token); r.Eroare == nil || r.Eroare.Cod != "rapoarte.pdfme_negenerat" {
		t.Fatalf("PDF salvat inainte de generare: %+v", r)
	}
	// PDF-ul generat de pdfme in interfata; aici, o pagina simpla.
	interfata, err := rapoarte.AnexaPDF(rapoarte.Document{Pagina: prima.Document.Pagina, Metadate: "pagina pdfme", Limba: "ro"})
	if err != nil {
		t.Fatal(err)
	}
	final := a.FinalizeazaRaportPdfme(prima.Token, base64.StdEncoding.EncodeToString(interfata))
	if final.Eroare != nil || final.Pagini < 2 {
		t.Fatalf("finalizare: %+v pagini %d", final.Eroare, final.Pagini)
	}
	pdf, err := base64.StdEncoding.DecodeString(final.PDF)
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
