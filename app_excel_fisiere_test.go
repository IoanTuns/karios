package main

import (
	"bytes"
	"fmt"
	"github.com/xuri/excelize/v2"
	"os"
	"path/filepath"
	"pontix/internal/datebaza"
	"pontix/internal/exportexcel"
	"pontix/internal/rapoarte"
	"testing"
)

func TestRegistreExcelFisiereFinale(t *testing.T) {
	director := t.TempDir()
	if v := os.Getenv("PONTIX_VERIFICARE_EXCEL"); v != "" {
		director = v
		if e := os.MkdirAll(v, 0700); e != nil {
			t.Fatal(e)
		}
	}
	persoane := []datebaza.Angajat{}
	for i := 0; i < 300; i++ {
		persoane = append(persoane, datebaza.Angajat{ID: fmt.Sprintf("anon-%d", i), Nume: fmt.Sprintf("Exemplu %03d", i), Prenume: "Test", Marca: fmt.Sprintf("TEST-%03d", i), Activ: true})
	}
	pontaj, e := exportexcel.GenereazaRegistruExcel(exportexcel.OptiuniExport{An: 2026, Luna: 10, Date: datebaza.Date{Angajati: persoane}})
	if e != nil {
		t.Fatal(e)
	}
	doc := rapoarte.Document{Titlu: "Raport anonim", Coloane: []rapoarte.Coloana{{Cod: "marca", Titlu: "Marca"}}}
	for i := 1; i <= 44; i++ {
		doc.Coloane = append(doc.Coloane, rapoarte.Coloana{Cod: fmt.Sprint(i), Titlu: fmt.Sprint(i)})
	}
	for i := 0; i < 300; i++ {
		r := rapoarte.Rand{Valori: []string{fmt.Sprintf("TEST-%03d", i)}}
		for j := 1; j <= 44; j++ {
			r.Valori = append(r.Valori, "7")
		}
		doc.Randuri = append(doc.Randuri, r)
	}
	raport, e := rapoarte.GenereazaExcel(doc)
	if e != nil {
		t.Fatal(e)
	}
	for nume, b := range map[string][]byte{"pontaj-anonim.xlsx": pontaj, "raport-anonim.xlsx": raport} {
		p := filepath.Join(director, nume)
		if e := exportexcel.SalveazaVerificat(p, b); e != nil {
			t.Fatal(e)
		}
		citit, e := os.ReadFile(p)
		if e != nil || !bytes.Equal(b, citit) {
			t.Fatal(e)
		}
		f, e := excelize.OpenFile(p)
		if e != nil {
			t.Fatal(e)
		}
		if len(f.GetSheetList()) == 0 {
			t.Fatal("fara foi")
		}
		f.Close()
	}
}
