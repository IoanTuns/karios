package rapoarte

import (
	"bytes"
	"github.com/xuri/excelize/v2"
	"testing"
)

func TestExcelRespectaDocumentulGrupat(t *testing.T) {
	d := Document{Titlu: "Fictiv", Coloane: []Coloana{{Cod: "angajat", Titlu: "Angajat"}, {Cod: "marca", Titlu: "Marca"}}, Randuri: []Rand{{Grup: true, Valori: []string{"Medici", ""}}, {Valori: []string{"=1+1", "GARDA-01"}}}}
	b, err := GenereazaExcel(d)
	if err != nil {
		t.Fatal(err)
	}
	f, err := excelize.OpenReader(bytes.NewReader(b))
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	s := f.GetSheetName(0)
	for cell, val := range map[string]string{"A5": "Medici", "A6": "=1+1", "B6": "GARDA-01"} {
		v, err := f.GetCellValue(s, cell)
		if err != nil || v != val {
			t.Fatal(cell, v, err)
		}
	}
	if formula, err := f.GetCellFormula(s, "A6"); err != nil || formula != "" {
		t.Fatal("formula injectata", formula, err)
	}
}
