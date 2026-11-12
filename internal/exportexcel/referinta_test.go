package exportexcel

import (
	"bytes"
	"encoding/json"
	"flag"
	"os"
	"path/filepath"
	"reflect"
	"testing"
	"time"

	"github.com/xuri/excelize/v2"

	"pontix/internal/calcul"
	"pontix/internal/datebaza"
	"pontix/internal/i18n"
)

var actualizeazaReferinta = flag.Bool("actualizeaza-export", false, "Scrie referinta exportului Excel pentru versiunea curenta.")

// continutRegistru: ce vede utilizatorul in fisier — foile, in ordine, cu valorile celulelor,
// si celulele unite. Stilurile si latimile coloanelor nu intra.
type continutRegistru struct {
	Versiune string          `json:"versiune"`
	Foi      []continutFoaie `json:"foi"`
}
type continutFoaie struct {
	Nume    string     `json:"nume"`
	Randuri [][]string `json:"randuri"`
	Unite   []string   `json:"unite"`
}

func citesteContinut(t *testing.T, b []byte) continutRegistru {
	t.Helper()
	f, err := excelize.OpenReader(bytes.NewReader(b))
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	rez := continutRegistru{Versiune: VersiuneExport}
	for _, nume := range f.GetSheetList() {
		randuri, err := f.GetRows(nume)
		if err != nil {
			t.Fatal(err)
		}
		unite, err := f.GetMergeCells(nume)
		if err != nil {
			t.Fatal(err)
		}
		foaie := continutFoaie{Nume: nume, Randuri: randuri, Unite: []string{}}
		for _, u := range unite {
			foaie.Unite = append(foaie.Unite, u.GetStartAxis()+":"+u.GetEndAxis())
		}
		rez.Foi = append(rez.Foi, foaie)
	}
	return rez
}

// TestReferintaExportExcel (F10-T05, D009): orice schimbare a continutului exportului cere
// o versiune noua (VersiuneExport) si un fisier de referinta nou; cele vechi raman.
func TestReferintaExportExcel(t *testing.T) {
	date := datebaza.Date{Angajati: []datebaza.Angajat{
		{ID: "ang-1", Nume: "Fictiv", Prenume: "Unu", Marca: "T001", Functie: "Functie test", TipContract: "NORMA_INTREAGA", Activ: true},
		{ID: "ang-2", Nume: "Fictiv", Prenume: "Doi", Marca: "T002", Functie: "Functie test", TipContract: "PARTIALA", Activ: true},
	}}
	schita := datebaza.SchitaLuna{An: 2026, Luna: 10,
		Celule:  []datebaza.CelulaSchita{{AngajatID: "ang-1", Data: "2026-10-01", Valoare: "Z7"}, {AngajatID: "ang-1", Data: "2026-10-24", Valoare: "N12"}},
		Absente: []datebaza.ZiAbsenta{{AngajatID: "ang-2", Data: "2026-10-05", Cod: "CO", Activ: true}},
		Calcul: calcul.RezultatPontaj{Totaluri: []calcul.TotalPontaj{
			{AngajatID: "ang-1", Valori: map[string]calcul.ValoarePontaj{"ore_lucrate": {Sutimi: 1900, Definita: true, Unitate: "ORE"}, "ore_norma": {Sutimi: 15400, Definita: true, Unitate: "ORE"}, "ZNL": {Sutimi: 1200, Definita: true, Unitate: "ORE"}}},
			{AngajatID: "ang-2", Valori: map[string]calcul.ValoarePontaj{"CO": {Sutimi: 100, Definita: true, Unitate: "ZILE"}, "ore_norma": {Sutimi: 8400, Definita: true, Unitate: "ORE"}}},
		}},
	}
	// F10-T06 (I130): referinta romana are aceleasi texte ca 1.2.0; engleza are referinta ei.
	for _, caz := range []struct {
		limba  i18n.Limba
		prefix string
	}{{i18n.Romana, "export-excel-v"}, {i18n.Engleza, "export-excel-en-v"}} {
		t.Run(string(caz.limba), func(t *testing.T) { verificaReferinta(t, date, schita, caz.limba, caz.prefix) })
	}
}

func verificaReferinta(t *testing.T, date datebaza.Date, schita datebaza.SchitaLuna, limba i18n.Limba, prefix string) {
	b, err := GenereazaRegistruExcel(OptiuniExport{An: 2026, Luna: 10, Date: date, Schita: schita,
		Sarbatori: map[string]string{"2026-10-14": "Sarbatoare fictiva"},
		GeneratLa: time.Date(2026, 10, 25, 14, 30, 0, 0, time.UTC), Limba: limba})
	if err != nil {
		t.Fatal(err)
	}
	primit := citesteContinut(t, b)
	cale := filepath.Join("testdata", prefix+VersiuneExport+".json")
	if *actualizeazaReferinta {
		continut, _ := json.MarshalIndent(primit, "", "  ")
		if err := os.MkdirAll("testdata", 0o700); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(cale, append(continut, '\n'), 0o600); err != nil {
			t.Fatal(err)
		}
	}
	brut, err := os.ReadFile(cale)
	if err != nil {
		t.Fatalf("lipseste referinta pentru versiunea %s: %v", VersiuneExport, err)
	}
	var asteptat continutRegistru
	if err := json.Unmarshal(brut, &asteptat); err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(asteptat, primit) {
		t.Fatalf("exportul Excel difera de referinta %s; o schimbare de continut cere o versiune noua a exportului", cale)
	}
}
