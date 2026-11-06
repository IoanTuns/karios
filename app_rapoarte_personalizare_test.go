package main

import (
	"context"
	"crypto/rand"
	"path/filepath"
	"pontix/internal/bd"
	"pontix/internal/i18n"
	"pontix/internal/rapoarte"
	"strings"
	"testing"
	"time"
)

func TestPersonalizareTokenAuditSiOriginalImuabil(t *testing.T) {
	cheie := make([]byte, bd.LungimeCheie)
	if _, e := rand.Read(cheie); e != nil {
		t.Fatal(e)
	}
	db, e := bd.Deschide(filepath.Join(t.TempDir(), "copie.db"), cheie)
	if e != nil {
		t.Fatal(e)
	}
	defer db.Close()
	if _, e = bd.Migreaza(context.Background(), db, bd.Migrari); e != nil {
		t.Fatal(e)
	}
	a := NewApp(db)
	d := rapoarte.Document{Titlu: "Demo", Sursa: "SQL_LIBER", Limba: i18n.Romana, GeneratLa: time.Date(2026, 11, 1, 12, 0, 0, 0, time.UTC), Pagina: rapoarte.Pagina{Format: "A4", Orientare: "PORTRET", MargineMM: 10, FontPT: 8}, Flexibil: &rapoarte.RaportFlexibil{ModAspect: rapoarte.ModPdfme, Sablon: rapoarte.SablonDinColoane([]rapoarte.ColoanaSablon{{Cod: "x", Titlu: "Eticheta"}}, rapoarte.Pagina{Format: "A4", Orientare: "PORTRET", MargineMM: 10, FontPT: 8}, rapoarte.OptiuniSablon{TextPagina: "{currentPage}"})}, Coloane: []rapoarte.Coloana{{Cod: "x", Titlu: "Eticheta"}}, Randuri: []rapoarte.Rand{{Valori: []string{"Initial"}}}, Metadate: "Istoric", Antet: []string{"", "Titlu", ""}, Subsol: []string{"", "", ""}}
	original := a.raspunsDocumentRaport(d, rapoarte.Cerere{Limba: "ro"})
	if original.Eroare != nil {
		t.Fatal(original.Eroare)
	}
	modificari := []rapoarte.EditareTextRaport{{Cheie: "celula:0:x", Text: "Personalizat"}}
	if r := a.PersonalizeazaRaport("../invalid", modificari); r.Eroare == nil || r.Eroare.Cod != "rapoarte.expirat" {
		t.Fatal(r.Eroare)
	}
	r := a.PersonalizeazaRaport(original.Token, modificari)
	if r.Eroare != nil {
		t.Fatal(r.Eroare)
	}
	if r.Token == original.Token || !r.Document.Personalizat || r.Document.Randuri[0].Valori[0] != "Personalizat" {
		t.Fatal("copie noua incorecta")
	}
	vechi, ok := a.raportPentru(original.Token)
	if !ok || vechi.Document.Randuri[0].Valori[0] != "Initial" {
		t.Fatal("original modificat")
	}
	// Copia editata cere un PDF nou din interfata, cu datele modificate.
	nou, ok := a.raportPentru(r.Token)
	if !ok || nou.PDF != nil || !r.GenerarePdfme || !strings.Contains(r.Intrari["tabel"], "Personalizat") {
		t.Fatal("PDF neactualizat")
	}
	count := func() int {
		var n int
		if e := db.QueryRow("SELECT COUNT(*) FROM jurnal_modificari WHERE tabel='rapoarte_personalizate'").Scan(&n); e != nil {
			t.Fatal(e)
		}
		return n
	}
	if count() != 1 {
		t.Fatal("jurnal lipsa")
	}
	var idAudit string
	if e := db.QueryRow("SELECT rand_id FROM jurnal_modificari WHERE tabel='rapoarte_personalizate' ORDER BY id DESC LIMIT 1").Scan(&idAudit); e != nil {
		t.Fatal(e)
	}
	if idAudit != r.Token || !strings.Contains(r.Document.Metadate, r.Token) {
		t.Fatal("copia exportata nu se leaga de jurnal")
	}

	// Din F09-T23 aranjarea in pagina o face pdfme in interfata (textul lung trece
	// pe randurile urmatoare); in Go raman limitele de lungime si cheile permise.
	bad := a.PersonalizeazaRaport(original.Token, []rapoarte.EditareTextRaport{{Cheie: "metadate", Text: "fals"}})
	if bad.Eroare == nil || count() != 1 {
		t.Fatal("istoric editabil")
	}
}
