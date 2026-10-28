package exportexcel

import (
	"bytes"
	"testing"

	"github.com/xuri/excelize/v2"

	"pontix/internal/calcul"
	"pontix/internal/datebaza"
)

func TestGenereazaRegistruExcel(t *testing.T) {
	date := datebaza.Date{
		Angajati: []datebaza.Angajat{
			{
				ID:          "ang-1",
				Nume:        "Popescu",
				Prenume:     "Ion",
				Marca:       "M001",
				Functie:     "Medic primar",
				TipContract: "NORMA_INTREAGA",
				Activ:       true,
			},
			{
				ID:          "ang-2",
				Nume:        "Ionescu",
				Prenume:     "Maria",
				Marca:       "M002",
				Functie:     "Asistent medical",
				TipContract: "PARTIALA",
				Activ:       true,
			},
		},
	}

	schita := datebaza.SchitaLuna{
		An:   2026,
		Luna: 10,
		Celule: []datebaza.CelulaSchita{
			{AngajatID: "ang-1", Data: "2026-10-01", Valoare: "Z7"},
			{AngajatID: "ang-1", Data: "2026-10-02", Valoare: "Z7"},
		},
		Absente: []datebaza.ZiAbsenta{
			{AngajatID: "ang-2", Data: "2026-10-05", Cod: "CO", Activ: true},
		},
		Calcul: calcul.RezultatPontaj{
			Totaluri: []calcul.TotalPontaj{
				{
					AngajatID: "ang-1",
					Valori: map[string]calcul.ValoarePontaj{
						"ore_lucrate": {Sutimi: 84000, Definita: true, Unitate: "ORE"}, // 14h
						"ore_norma":   {Sutimi: 168000, Definita: true, Unitate: "ORE"},
					},
				},
				{
					AngajatID: "ang-2",
					Valori: map[string]calcul.ValoarePontaj{
						"CO":        {Sutimi: 100, Definita: true, Unitate: "ZILE"},
						"ore_norma": {Sutimi: 84000, Definita: true, Unitate: "ORE"},
					},
				},
			},
		},
	}

	b, err := GenereazaRegistruExcel(OptiuniExport{
		An:        2026,
		Luna:      10,
		Date:      date,
		Schita:    schita,
		Sarbatori: map[string]string{"2026-10-14": "Sf. Parascheva"},
	})
	if err != nil {
		t.Fatalf("eroare generare registru excel: %v", err)
	}
	if len(b) == 0 {
		t.Fatal("registrul excel generat este gol")
	}

	// Citire înapoi cu excelize
	f, err := excelize.OpenReader(bytes.NewReader(b))
	if err != nil {
		t.Fatalf("eroare deschidere XLSX generat: %v", err)
	}
	defer f.Close()

	foi := f.GetSheetList()
	if len(foi) != 2 {
		t.Fatalf("așteptat 2 foi de calcul, obținut: %v", foi)
	}
	if foi[0] != "Pontaj Lunar" || foi[1] != "Centralizator Salarii" {
		t.Errorf("foile au denumiri incorecte: %v", foi)
	}

	// Verificăm celula din foaia Pontaj Lunar
	titlu, err := f.GetCellValue("Pontaj Lunar", "A1")
	if err != nil || titlu == "" {
		t.Errorf("titlul foii de pontaj lipsește sau e invalid: %v", err)
	}

	numeAng1, _ := f.GetCellValue("Pontaj Lunar", "B5")
	if numeAng1 != "Ionescu Maria" { // sortat alfabetic
		t.Errorf("angajatul de pe primul rând așteptat 'Ionescu Maria', obținut: %q", numeAng1)
	}

	// Verificăm formula SUM din Centralizator Salarii
	formulaF, _ := f.GetCellFormula("Centralizator Salarii", "F7")
	if formulaF != "SUM(F5:F6)" {
		t.Errorf("formula SUM din total general invalidă: %q", formulaF)
	}

	numeFisier := NumeFisierPontaj(2026, 10)
	if numeFisier != "Pontix_Pontaj_2026_10.xlsx" {
		t.Errorf("nume fișier generat greșit: %s", numeFisier)
	}
}
