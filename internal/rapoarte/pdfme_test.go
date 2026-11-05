package rapoarte

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/phpdave11/gofpdf"

	"pontix/internal/i18n"
)

const sablonPdfmeTest = `{"basePdf":{"width":210,"height":297,"padding":[10,10,10,10]},"schemas":[[
{"name":"titlu","type":"text","readOnly":true,"content":"{titlu}","position":{"x":10,"y":10},"width":100,"height":10},
{"name":"tabel","type":"table","head":["Angajat","Ore"],"position":{"x":10,"y":30},"width":190,"height":20}]]}`

func documentPdfmeTest(sablon string) Document {
	return Document{
		Flexibil:        &RaportFlexibil{ModAspect: ModPdfme, Sablon: sablon},
		Titlu:           "Gărzi noiembrie",
		VariabileRaport: map[string]string{"an": "2026", "luna": "11", "perioada": "2026-11"},
		Antet:           []string{"Spital {{an}}", "", "Secția {{luna}}"},
		Subsol:          []string{"", "", ""},
		Coloane:         []Coloana{{Cod: "angajat", Titlu: "Angajat"}, {Cod: "ore", Titlu: "Ore"}},
		Randuri:         []Rand{{Valori: []string{"Popescu Ș.", "12"}}, {Valori: []string{"Ionescu Ț.", "24"}}},
		Pagina:          Pagina{Format: "A4", Orientare: "PORTRET", MargineMM: 10, FontPT: 9},
		Limba:           i18n.Limba("ro"),
		Legenda:         []string{"SQL: SELECT angajat, ore FROM garzi"},
		Metadate:        "Gărzi | rev | 0.92.33 | test",
		GeneratLa:       time.Date(2026, 11, 5, 10, 30, 0, 0, time.UTC),
	}
}

func codEroarePdfme(err error) string {
	var e *i18n.Eroare
	if errors.As(err, &e) {
		return e.Cod
	}
	return ""
}

func TestSablonPdfmeValidare(t *testing.T) {
	if _, err := citesteSablonPdfme(sablonPdfmeTest); err != nil {
		t.Fatalf("sablon valid refuzat: %v", err)
	}
	gresite := map[string]string{
		"gol":          ``,
		"json invalid": `{`,
		"pdf extern":   `{"basePdf":"data:application/pdf;base64,AAAA","schemas":[[]]}`,
		"fara margini": `{"basePdf":{"width":210,"height":297},"schemas":[[]]}`,
		"fara pagini":  `{"basePdf":{"width":210,"height":297,"padding":[0,0,0,0]},"schemas":[]}`,
		"pagina mica":  `{"basePdf":{"width":10,"height":297,"padding":[0,0,0,0]},"schemas":[[]]}`,
		"tip svg":      `{"basePdf":{"width":210,"height":297,"padding":[0,0,0,0]},"schemas":[[{"name":"x","type":"svg"}]]}`,
		"fara nume":    `{"basePdf":{"width":210,"height":297,"padding":[0,0,0,0]},"schemas":[[{"name":"","type":"text"}]]}`,
		"doua tabele":  `{"basePdf":{"width":210,"height":297,"padding":[0,0,0,0]},"schemas":[[{"name":"tabel","type":"table"}],[{"name":"tabel","type":"table"}]]}`,
		"prea mare":    `{"basePdf":{"width":210,"height":297,"padding":[0,0,0,0]},"schemas":[[]],"x":"` + strings.Repeat("a", maxSablonPdfme) + `"}`,
	}
	for nume, s := range gresite {
		if _, err := citesteSablonPdfme(s); codEroarePdfme(err) != "rapoarte.pdfme_invalid" {
			t.Errorf("%s: astept pdfme_invalid, primit %v", nume, err)
		}
	}
}

func TestValideazaFlexibilPdfme(t *testing.T) {
	f := &RaportFlexibil{ModAspect: ModPdfme, Sablon: sablonPdfmeTest}
	if err := valideazaFlexibil(f); err != nil {
		t.Fatalf("model pdfme valid refuzat: %v", err)
	}
	f.Sablon = `{}`
	if codEroarePdfme(valideazaFlexibil(f)) != "rapoarte.pdfme_invalid" {
		t.Fatal("sablon gol acceptat")
	}
	// Un sablon ramas pe alt mod de aspect nu se pastreaza tacit.
	if codEroarePdfme(valideazaFlexibil(&RaportFlexibil{ModAspect: "TABEL", Sablon: sablonPdfmeTest})) != "rapoarte.pdfme_invalid" {
		t.Fatal("sablon acceptat pe modul TABEL")
	}
}

func TestDatePdfme(t *testing.T) {
	d := documentPdfmeTest(sablonPdfmeTest)
	intrari, err := DatePdfme(d)
	if err != nil {
		t.Fatal(err)
	}
	if intrari["titlu"] != "Gărzi noiembrie" || intrari["perioada"] != "2026-11" || intrari["generat_la"] != "05.11.2026 10:30" {
		t.Fatalf("variabile gresite: %v", intrari)
	}
	if intrari["antet_stanga"] != "Spital 2026" || intrari["antet_dreapta"] != "Secția 11" {
		t.Fatalf("antet neinlocuit: %q %q", intrari["antet_stanga"], intrari["antet_dreapta"])
	}
	var randuri [][]string
	if err := json.Unmarshal([]byte(intrari[TabelPdfme]), &randuri); err != nil || len(randuri) != 2 || randuri[1][0] != "Ionescu Ț." {
		t.Fatalf("tabel gresit: %v %v", randuri, err)
	}
}

func TestDatePdfmeColoaneDiferite(t *testing.T) {
	d := documentPdfmeTest(strings.Replace(sablonPdfmeTest, `["Angajat","Ore"]`, `["Angajat"]`, 1))
	_, err := DatePdfme(d)
	var e *i18n.Eroare
	if !errors.As(err, &e) || e.Cod != "rapoarte.pdfme_coloane" || e.Param["sablon"] != "1" || e.Param["date"] != "2" {
		t.Fatalf("astept pdfme_coloane 1/2, primit %v", err)
	}
}

func TestGenereazaPDFRefuzaPdfme(t *testing.T) {
	if _, err := GenereazaPDF(documentPdfmeTest(sablonPdfmeTest)); codEroarePdfme(err) != "rapoarte.pdfme_in_interfata" {
		t.Fatalf("motorul vechi nu trebuie sa deseneze sablonul pdfme: %v", err)
	}
}

// pdfSimplu tine locul PDF-ului generat de pdfme in interfata.
func pdfSimplu(t *testing.T, pagini int) []byte {
	t.Helper()
	f := gofpdf.New("P", "mm", "A4", "")
	for range pagini {
		f.AddPage()
	}
	var b bytes.Buffer
	if err := f.Output(&b); err != nil {
		t.Fatal(err)
	}
	return b.Bytes()
}

func TestFinalizeazaPdfmeAdaugaAnexa(t *testing.T) {
	ctx := context.Background()
	d := documentPdfmeTest(sablonPdfmeTest)
	pdf, pagini, err := FinalizeazaPdfme(ctx, d, pdfSimplu(t, 2))
	if err != nil {
		t.Fatal(err)
	}
	if pagini != 3 {
		t.Fatalf("astept 2 pagini pdfme + 1 anexa, primit %d", pagini)
	}
	if n, err := NumaraPaginiPDF(ctx, pdf); err != nil || n != 3 {
		t.Fatalf("PDF final: %d pagini, %v", n, err)
	}
	// Fara explicatii si istoric nu se adauga anexa.
	d.Legenda, d.Metadate = nil, ""
	if _, pagini, err = FinalizeazaPdfme(ctx, d, pdfSimplu(t, 1)); err != nil || pagini != 1 {
		t.Fatalf("fara anexa: %d pagini, %v", pagini, err)
	}
}

func TestFinalizeazaPdfmeRefuzaDateInvalide(t *testing.T) {
	ctx := context.Background()
	d := documentPdfmeTest(sablonPdfmeTest)
	for nume, pdf := range map[string][]byte{"gol": nil, "text": []byte("nu este pdf"), "prea mare": make([]byte, maxPDFPdfme+1)} {
		if _, _, err := FinalizeazaPdfme(ctx, d, pdf); codEroarePdfme(err) != "rapoarte.pdfme_pdf_invalid" {
			t.Errorf("%s: astept pdfme_pdf_invalid, primit %v", nume, err)
		}
	}
	d.Flexibil.ModAspect = "TABEL"
	if _, _, err := FinalizeazaPdfme(ctx, d, pdfSimplu(t, 1)); codEroarePdfme(err) != "rapoarte.pdfme_pdf_invalid" {
		t.Fatalf("document fara pdfme acceptat: %v", err)
	}
}
