package exportexcel

import (
	"bytes"
	"github.com/xuri/excelize/v2"
	"pontix/internal/datebaza"
	"strings"
	"testing"
)

func TestMarcajSubordonatInAmbeleFoi(t *testing.T) {
	d := datebaza.Date{Angajati: []datebaza.Angajat{{ID: "g", Nume: "Fictiv", Prenume: "Ana", Marca: "G-TEST", Activ: true, AngajatPrincipalID: "p"}, {ID: "p", Nume: "Fictiv", Prenume: "Ana", Marca: "P-TEST", Activ: true}}}
	b, e := GenereazaRegistruExcel(OptiuniExport{An: 2026, Luna: 10, Date: d, Schita: datebaza.SchitaLuna{An: 2026, Luna: 10}})
	if e != nil {
		t.Fatal(e)
	}
	f, e := excelize.OpenReader(bytes.NewReader(b))
	if e != nil {
		t.Fatal(e)
	}
	defer f.Close()
	for _, foaie := range f.GetSheetList() {
		rows, e := f.GetRows(foaie)
		if e != nil {
			t.Fatal(e)
		}
		marcate, simple := 0, 0
		for _, r := range rows {
			for _, v := range r {
				if strings.HasPrefix(v, "Fictiv Ana (") {
					marcate++
				}
				if v == "Fictiv Ana" {
					simple++
				}
			}
		}
		if marcate != 1 || simple != 1 {
			t.Fatal(foaie, marcate, simple)
		}
	}
	if d.Angajati[0].Nume != "Fictiv" {
		t.Fatal("nume modificat")
	}
}
