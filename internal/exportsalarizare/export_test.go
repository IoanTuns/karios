package exportsalarizare

import (
	"bytes"
	"encoding/json"
	"flag"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/xuri/excelize/v2"

	"pontix/internal/calcul"
)

var actualizeaza = flag.Bool("actualizeaza-export-salarizare", false, "Scrie referintele exportului pentru salarizare.")

func ore(h float64) calcul.ValoarePontaj {
	return calcul.ValoarePontaj{Sutimi: int64(h * 6000), Unitate: "ORE", Definita: true}
}
func zile(n int64) calcul.ValoarePontaj {
	return calcul.ValoarePontaj{Sutimi: n * 100, Unitate: "ZILE", Definita: true}
}

func randuriTest() []Rand {
	return []Rand{
		{Identitate: map[string]string{"marca": "T001", "nume": "Fictiv", "prenume": "Unu", "nume_complet": "Fictiv Unu", "perioada": "2026-10"},
			Valori: map[string]calcul.ValoarePontaj{"OL": ore(168.5), "OS": ore(2.25), "CO": zile(2), "CM": zile(0)}},
		{Identitate: map[string]string{"marca": "T002", "nume": "Fictiv", "prenume": "Doi", "nume_complet": "Fictiv Doi", "perioada": "2026-10"},
			Valori: map[string]calcul.ValoarePontaj{"OL": ore(84), "OS": {Unitate: "ORE", Motiv: calcul.MotivGardaFaraOS}, "CO": zile(0), "CM": zile(3)}},
	}
}

func modelTest(format string) Model {
	return Model{Denumire: "Test", Format: format, Antet: true, BOM: true, Zecimale: 2, Coloane: []Coloana{
		{Sursa: "marca", Titlu: "Marca"}, {Sursa: "nume_complet", Titlu: "Nume"}, {Sursa: "perioada", Titlu: "Luna"},
		{Sursa: "total:OL", Titlu: "Ore lucrate"}, {Sursa: "total:OS", Titlu: "Ore suplimentare"},
		{Sursa: "total:CO", Titlu: "CO"}, {Sursa: "total:CM", Titlu: "CM"}, {Sursa: SursaConstanta, Titlu: "Cod", Valoare: "PX"},
	}}
}

func TestValidareModel(t *testing.T) {
	m := modelTest(FormatCSV)
	if err := Valideaza(&m); err != nil || m.Separator != ";" || m.SeparatorZecimal != "," || m.FormaOre != OreZecimal {
		t.Fatal(m, err)
	}
	for nume, schimba := range map[string]func(*Model){
		"format":   func(m *Model) { m.Format = "PDF" },
		"sursa":    func(m *Model) { m.Coloane[0].Sursa = "cnp" },
		"total":    func(m *Model) { m.Coloane[3].Sursa = "total:SALARIU" },
		"zecimal":  func(m *Model) { m.Separator, m.SeparatorZecimal = ",", "," },
		"coloane":  func(m *Model) { m.Coloane = nil },
		"denumire": func(m *Model) { m.Denumire = " " },
	} {
		m := modelTest(FormatCSV)
		schimba(&m)
		if Valideaza(&m) == nil {
			t.Errorf("%s: model invalid acceptat", nume)
		}
	}
}

func TestFormaOre(t *testing.T) {
	m := modelTest(FormatCSV)
	m.FormaOre, m.Antet = OreOreMinute, false
	if err := Valideaza(&m); err != nil {
		t.Fatal(err)
	}
	tabel, necalculate := Tabel(m, randuriTest())
	if tabel[0][3] != "168:30" || tabel[0][4] != "2:15" || tabel[0][5] != "2,00" || tabel[1][4] != "" || necalculate != 1 {
		t.Fatal(tabel, necalculate)
	}
	m.FormaOre, m.SeparatorZecimal, m.Zecimale = OreZecimal, ".", 1
	if tabel, _ = Tabel(m, randuriTest()); tabel[0][3] != "168.5" || tabel[0][4] != "2.3" {
		t.Fatal(tabel)
	}
}

// TestReferintaExportSalarizare (D009): CSV-ul exact (octeti) si continutul fisierului Excel.
func TestReferintaExportSalarizare(t *testing.T) {
	m := modelTest(FormatCSV)
	if err := Valideaza(&m); err != nil {
		t.Fatal(err)
	}
	csvDate, ext, necalculate, err := Genereaza(m, randuriTest())
	if err != nil || ext != ".csv" || necalculate != 1 || !bytes.HasPrefix(csvDate, []byte("\xef\xbb\xbf")) {
		t.Fatal(ext, necalculate, err)
	}
	mx := modelTest(FormatXLSX)
	if err := Valideaza(&mx); err != nil {
		t.Fatal(err)
	}
	xlsx, ext, _, err := Genereaza(mx, randuriTest())
	if err != nil || ext != ".xlsx" {
		t.Fatal(ext, err)
	}
	f, err := excelize.OpenReader(bytes.NewReader(xlsx))
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	randuri, _ := f.GetRows(f.GetSheetName(0))
	tip, _ := f.GetCellType(f.GetSheetName(0), "D2")
	primit := map[string]any{"versiune": VersiuneExportSalarizare, "csv": strings.Split(string(csvDate), "\r\n"), "xlsx": randuri, "xlsx_D2_numeric": tip != excelize.CellTypeSharedString && tip != excelize.CellTypeInlineString}
	cale := filepath.Join("testdata", "export-salarizare-v"+VersiuneExportSalarizare+".json")
	if *actualizeaza {
		continut, _ := json.MarshalIndent(primit, "", "  ")
		_ = os.MkdirAll("testdata", 0o700)
		if err := os.WriteFile(cale, append(continut, '\n'), 0o600); err != nil {
			t.Fatal(err)
		}
	}
	brut, err := os.ReadFile(cale)
	if err != nil {
		t.Fatal(err)
	}
	var asteptat, obtinut any
	_ = json.Unmarshal(brut, &asteptat)
	text, _ := json.Marshal(primit)
	_ = json.Unmarshal(text, &obtinut)
	if !reflect.DeepEqual(asteptat, obtinut) {
		t.Fatalf("exportul pentru salarizare difera de referinta %s", cale)
	}
}

// F08-T18: o valoare care nu se aplica randului e goala si nu se numara ca necalculata.
func TestNuSeAplicaNuENecalculata(t *testing.T) {
	m := modelTest(FormatCSV)
	if err := Valideaza(&m); err != nil {
		t.Fatal(err)
	}
	r := randuriTest()
	r[1].Valori["OS"] = calcul.ValoarePontaj{Unitate: "ORE", Motiv: calcul.MotivGardaFaraOS, NuSeAplica: true}
	tabel, necalculate := Tabel(m, r)
	if tabel[2][4] != "" || necalculate != 0 {
		t.Fatal(tabel, necalculate)
	}
}
