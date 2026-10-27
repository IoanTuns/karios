package exportexcel

import (
	"bytes"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/xuri/excelize/v2"

	"pontix/internal/calcul"
	"pontix/internal/datebaza"
)

// VersiuneExport reprezintă versiunea formatului de export Excel (D009).
const VersiuneExport = "1.0.0"

// OptiuniExport conține parametrii pentru generarea registrului Excel.
type OptiuniExport struct {
	An        int
	Luna      int
	Date      datebaza.Date
	Schita    datebaza.SchitaLuna
	Sarbatori map[string]string // data AAAA-LL-ZZ -> denumire
}

// GenereazaRegistruExcel creează un fișier XLSX complet cu foaia de pontaj și centralizatorul.
func GenereazaRegistruExcel(opt OptiuniExport) ([]byte, error) {
	f := excelize.NewFile()
	defer f.Close()

	numeLuni := []string{"Ianuarie", "Februarie", "Martie", "Aprilie", "Mai", "Iunie", "Iulie", "August", "Septembrie", "Octombrie", "Noiembrie", "Decembrie"}
	numeLuna := ""
	if opt.Luna >= 1 && opt.Luna <= 12 {
		numeLuna = numeLuni[opt.Luna-1]
	}

	numarZile := time.Date(opt.An, time.Month(opt.Luna)+1, 0, 12, 0, 0, 0, time.UTC).Day()

	// 1. Foaie Pontaj Lunar
	numeFoaiePontaj := "Pontaj Lunar"
	indexPontaj, err := f.NewSheet(numeFoaiePontaj)
	if err != nil {
		return nil, err
	}
	f.DeleteSheet("Sheet1")
	f.SetActiveSheet(indexPontaj)

	stilTitlu, _ := f.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Bold: true, Size: 14, Color: "163F25"},
		Alignment: &excelize.Alignment{Vertical: "center"},
	})
	stilSubtitlu, _ := f.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Italic: true, Size: 10, Color: "496356"},
		Alignment: &excelize.Alignment{Vertical: "center"},
	})
	stilAntet, _ := f.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Bold: true, Size: 10, Color: "FFFFFF"},
		Fill:      excelize.Fill{Type: "pattern", Color: []string{"2C4540"}, Pattern: 1},
		Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center", WrapText: true},
		Border: []excelize.Border{
			{Type: "left", Color: "D7DED9", Style: 1},
			{Type: "right", Color: "D7DED9", Style: 1},
			{Type: "top", Color: "D7DED9", Style: 1},
			{Type: "bottom", Color: "D7DED9", Style: 1},
		},
	})
	stilWeekendAntet, _ := f.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Bold: true, Size: 10, Color: "FFFFFF"},
		Fill:      excelize.Fill{Type: "pattern", Color: []string{"3E5C55"}, Pattern: 1},
		Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center"},
		Border: []excelize.Border{
			{Type: "left", Color: "D7DED9", Style: 1},
			{Type: "right", Color: "D7DED9", Style: 1},
			{Type: "top", Color: "D7DED9", Style: 1},
			{Type: "bottom", Color: "D7DED9", Style: 1},
		},
	})
	stilCelulaNorm, _ := f.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Size: 9},
		Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center"},
		Border: []excelize.Border{
			{Type: "left", Color: "E0E0E0", Style: 1},
			{Type: "right", Color: "E0E0E0", Style: 1},
			{Type: "top", Color: "E0E0E0", Style: 1},
			{Type: "bottom", Color: "E0E0E0", Style: 1},
		},
	})
	stilCelulaWeekend, _ := f.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Size: 9},
		Fill:      excelize.Fill{Type: "pattern", Color: []string{"F4F6F5"}, Pattern: 1},
		Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center"},
		Border: []excelize.Border{
			{Type: "left", Color: "E0E0E0", Style: 1},
			{Type: "right", Color: "E0E0E0", Style: 1},
			{Type: "top", Color: "E0E0E0", Style: 1},
			{Type: "bottom", Color: "E0E0E0", Style: 1},
		},
	})
	stilTextStanga, _ := f.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Size: 9},
		Alignment: &excelize.Alignment{Horizontal: "left", Vertical: "center"},
		Border: []excelize.Border{
			{Type: "left", Color: "E0E0E0", Style: 1},
			{Type: "right", Color: "E0E0E0", Style: 1},
			{Type: "top", Color: "E0E0E0", Style: 1},
			{Type: "bottom", Color: "E0E0E0", Style: 1},
		},
	})
	stilTotalBold, _ := f.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Bold: true, Size: 9, Color: "163F25"},
		Fill:      excelize.Fill{Type: "pattern", Color: []string{"E8F5E9"}, Pattern: 1},
		Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center"},
		Border: []excelize.Border{
			{Type: "left", Color: "C8E6C9", Style: 1},
			{Type: "right", Color: "C8E6C9", Style: 1},
			{Type: "top", Color: "C8E6C9", Style: 1},
			{Type: "bottom", Color: "C8E6C9", Style: 1},
		},
	})

	// Antet raport
	_ = f.SetCellValue(numeFoaiePontaj, "A1", "PONTIX — FOAIE COLECTIVĂ DE PREZENȚĂ ȘI PONTAJ")
	_ = f.SetCellStyle(numeFoaiePontaj, "A1", "A1", stilTitlu)
	_ = f.SetCellValue(numeFoaiePontaj, "A2", fmt.Sprintf("Luna: %s %d | Generat la: %s", numeLuna, opt.An, time.Now().Format("02.01.2006 15:04")))
	_ = f.SetCellStyle(numeFoaiePontaj, "A2", "A2", stilSubtitlu)

	// Coloane fixe
	anteteFixe := []string{"Nr. crt.", "Nume și prenume", "Marca", "Funcție"}
	randAntet := 4

	colIdx := 1
	for _, col := range anteteFixe {
		celula, _ := excelize.CoordinatesToCellName(colIdx, randAntet)
		_ = f.SetCellValue(numeFoaiePontaj, celula, col)
		_ = f.SetCellStyle(numeFoaiePontaj, celula, celula, stilAntet)
		colIdx++
	}

	// Coloane zile (1..numarZile)
	for zi := 1; zi <= numarZile; zi++ {
		tZi := time.Date(opt.An, time.Month(opt.Luna), zi, 12, 0, 0, 0, time.UTC)
		esteWeekend := tZi.Weekday() == time.Saturday || tZi.Weekday() == time.Sunday
		celula, _ := excelize.CoordinatesToCellName(colIdx, randAntet)
		_ = f.SetCellValue(numeFoaiePontaj, celula, zi)
		if esteWeekend {
			_ = f.SetCellStyle(numeFoaiePontaj, celula, celula, stilWeekendAntet)
		} else {
			_ = f.SetCellStyle(numeFoaiePontaj, celula, celula, stilAntet)
		}
		_ = f.SetColWidth(numeFoaiePontaj, celula[:len(celula)-1], celula[:len(celula)-1], 4.2)
		colIdx++
	}

	// Coloane totaluri
	anteteTotaluri := []string{"Total Ore", "Ore Noapte / Gardă", "Ore ZNL", "Normă", "CO (zile)", "CM (zile)"}
	for _, tot := range anteteTotaluri {
		celula, _ := excelize.CoordinatesToCellName(colIdx, randAntet)
		_ = f.SetCellValue(numeFoaiePontaj, celula, tot)
		_ = f.SetCellStyle(numeFoaiePontaj, celula, celula, stilAntet)
		_ = f.SetColWidth(numeFoaiePontaj, celula[:len(celula)-1], celula[:len(celula)-1], 11)
		colIdx++
	}

	// Indexare ciorne și totaluri
	hartaCelule := make(map[string]string)
	for _, c := range opt.Schita.Celule {
		hartaCelule[c.AngajatID+"/"+c.Data] = c.Valoare
	}
	hartaAbsente := make(map[string]string)
	for _, a := range opt.Schita.Absente {
		hartaAbsente[a.AngajatID+"/"+a.Data] = a.Cod
	}
	hartaTotaluri := make(map[string]calcul.TotalPontaj)
	for _, t := range opt.Schita.Calcul.Totaluri {
		hartaTotaluri[t.AngajatID] = t
	}

	// Sortăm angajații alfabetic
	angajati := append([]datebaza.Angajat{}, opt.Date.Angajati...)
	sort.Slice(angajati, func(i, j int) bool {
		if angajati[i].Nume != angajati[j].Nume {
			return angajati[i].Nume < angajati[j].Nume
		}
		return angajati[i].Prenume < angajati[j].Prenume
	})

	randCurent := randAntet + 1
	nrCrt := 1

	for _, a := range angajati {
		if !a.Activ {
			continue
		}
		cIdx := 1

		// Nr. crt.
		cName, _ := excelize.CoordinatesToCellName(cIdx, randCurent)
		_ = f.SetCellValue(numeFoaiePontaj, cName, nrCrt)
		_ = f.SetCellStyle(numeFoaiePontaj, cName, cName, stilCelulaNorm)
		cIdx++

		// Nume
		cName, _ = excelize.CoordinatesToCellName(cIdx, randCurent)
		_ = f.SetCellValue(numeFoaiePontaj, cName, strings.TrimSpace(a.Nume+" "+a.Prenume))
		_ = f.SetCellStyle(numeFoaiePontaj, cName, cName, stilTextStanga)
		cIdx++

		// Marca
		cName, _ = excelize.CoordinatesToCellName(cIdx, randCurent)
		_ = f.SetCellValue(numeFoaiePontaj, cName, a.Marca)
		_ = f.SetCellStyle(numeFoaiePontaj, cName, cName, stilCelulaNorm)
		cIdx++

		// Functie
		cName, _ = excelize.CoordinatesToCellName(cIdx, randCurent)
		_ = f.SetCellValue(numeFoaiePontaj, cName, a.Functie)
		_ = f.SetCellStyle(numeFoaiePontaj, cName, cName, stilTextStanga)
		cIdx++

		// Zile 1..numarZile
		for zi := 1; zi <= numarZile; zi++ {
			data := fmt.Sprintf("%04d-%02d-%02d", opt.An, opt.Luna, zi)
			k := a.ID + "/" + data
			valoare := hartaCelule[k]
			if absCod, ok := hartaAbsente[k]; ok {
				valoare = absCod
			}

			tZi := time.Date(opt.An, time.Month(opt.Luna), zi, 12, 0, 0, 0, time.UTC)
			esteWeekend := tZi.Weekday() == time.Saturday || tZi.Weekday() == time.Sunday

			cName, _ = excelize.CoordinatesToCellName(cIdx, randCurent)
			_ = f.SetCellValue(numeFoaiePontaj, cName, valoare)
			if esteWeekend {
				_ = f.SetCellStyle(numeFoaiePontaj, cName, cName, stilCelulaWeekend)
			} else {
				_ = f.SetCellStyle(numeFoaiePontaj, cName, cName, stilCelulaNorm)
			}
			cIdx++
		}

		// Totaluri calculate din pontaj
		tot := hartaTotaluri[a.ID]
		extrageVal := func(cheie string) string {
			if v, ok := tot.Valori[cheie]; ok && v.Definita {
				return calcul.TextValoarePontaj(v)
			}
			return "0"
		}

		// Total ore
		cName, _ = excelize.CoordinatesToCellName(cIdx, randCurent)
		_ = f.SetCellValue(numeFoaiePontaj, cName, extrageVal("ore_lucrate"))
		_ = f.SetCellStyle(numeFoaiePontaj, cName, cName, stilTotalBold)
		cIdx++

		// Ore noapte / garda
		cName, _ = excelize.CoordinatesToCellName(cIdx, randCurent)
		_ = f.SetCellValue(numeFoaiePontaj, cName, extrageVal("ore_garda"))
		_ = f.SetCellStyle(numeFoaiePontaj, cName, cName, stilCelulaNorm)
		cIdx++

		// Ore ZNL
		cName, _ = excelize.CoordinatesToCellName(cIdx, randCurent)
		_ = f.SetCellValue(numeFoaiePontaj, cName, extrageVal("ZNL"))
		_ = f.SetCellStyle(numeFoaiePontaj, cName, cName, stilCelulaNorm)
		cIdx++

		// Norma
		cName, _ = excelize.CoordinatesToCellName(cIdx, randCurent)
		_ = f.SetCellValue(numeFoaiePontaj, cName, extrageVal("ore_norma"))
		_ = f.SetCellStyle(numeFoaiePontaj, cName, cName, stilCelulaNorm)
		cIdx++

		// CO
		cName, _ = excelize.CoordinatesToCellName(cIdx, randCurent)
		_ = f.SetCellValue(numeFoaiePontaj, cName, extrageVal("CO"))
		_ = f.SetCellStyle(numeFoaiePontaj, cName, cName, stilCelulaNorm)
		cIdx++

		// CM
		cName, _ = excelize.CoordinatesToCellName(cIdx, randCurent)
		_ = f.SetCellValue(numeFoaiePontaj, cName, extrageVal("CM"))
		_ = f.SetCellStyle(numeFoaiePontaj, cName, cName, stilCelulaNorm)
		cIdx++

		nrCrt++
		randCurent++
	}

	// Setare lățimi coloane fixe
	_ = f.SetColWidth(numeFoaiePontaj, "A", "A", 6)
	_ = f.SetColWidth(numeFoaiePontaj, "B", "B", 24)
	_ = f.SetColWidth(numeFoaiePontaj, "C", "C", 11)
	_ = f.SetColWidth(numeFoaiePontaj, "D", "D", 20)

	// 2. Foaie Centralizator Salarizare
	numeFoaieCentralizator := "Centralizator Salarii"
	_, err = f.NewSheet(numeFoaieCentralizator)
	if err == nil {
		_ = f.SetCellValue(numeFoaieCentralizator, "A1", "CENTRALIZATOR LUNAR PENTRU SALARIZARE")
		_ = f.SetCellStyle(numeFoaieCentralizator, "A1", "A1", stilTitlu)
		_ = f.SetCellValue(numeFoaieCentralizator, "A2", fmt.Sprintf("Luna: %s %d | Export securizat Pontix", numeLuna, opt.An))
		_ = f.SetCellStyle(numeFoaieCentralizator, "A2", "A2", stilSubtitlu)

		anteteCentr := []string{
			"Nr. crt.", "Nume și prenume", "Marca", "Funcție", "Contract",
			"Ore Normă", "Ore Efectiv Lucrate", "Ore Gărzi", "Ore Nelucrătoare (ZNL)",
			"Zile Concediu Odihnă", "Zile Concediu Medical",
		}
		cIdx := 1
		for _, ac := range anteteCentr {
			cell, _ := excelize.CoordinatesToCellName(cIdx, 4)
			_ = f.SetCellValue(numeFoaieCentralizator, cell, ac)
			_ = f.SetCellStyle(numeFoaieCentralizator, cell, cell, stilAntet)
			_ = f.SetColWidth(numeFoaieCentralizator, cell[:len(cell)-1], cell[:len(cell)-1], 15)
			cIdx++
		}
		_ = f.SetColWidth(numeFoaieCentralizator, "A", "A", 6)
		_ = f.SetColWidth(numeFoaieCentralizator, "B", "B", 24)
		_ = f.SetColWidth(numeFoaieCentralizator, "C", "C", 11)
		_ = f.SetColWidth(numeFoaieCentralizator, "D", "D", 20)

		rC := 5
		nrC := 1
		for _, a := range angajati {
			if !a.Activ {
				continue
			}
			tot := hartaTotaluri[a.ID]
			extrageFloat := func(cheie string) float64 {
				if v, ok := tot.Valori[cheie]; ok && v.Definita {
					if v.Unitate == "ORE" {
						return float64(v.Sutimi) / 6000.0
					}
					return float64(v.Sutimi) / 100.0
				}
				return 0
			}

			_ = f.SetCellValue(numeFoaieCentralizator, fmt.Sprintf("A%d", rC), nrC)
			_ = f.SetCellValue(numeFoaieCentralizator, fmt.Sprintf("B%d", rC), strings.TrimSpace(a.Nume+" "+a.Prenume))
			_ = f.SetCellValue(numeFoaieCentralizator, fmt.Sprintf("C%d", rC), a.Marca)
			_ = f.SetCellValue(numeFoaieCentralizator, fmt.Sprintf("D%d", rC), a.Functie)
			_ = f.SetCellValue(numeFoaieCentralizator, fmt.Sprintf("E%d", rC), a.TipContract)

			_ = f.SetCellValue(numeFoaieCentralizator, fmt.Sprintf("F%d", rC), extrageFloat("ore_norma"))
			_ = f.SetCellValue(numeFoaieCentralizator, fmt.Sprintf("G%d", rC), extrageFloat("ore_lucrate"))
			_ = f.SetCellValue(numeFoaieCentralizator, fmt.Sprintf("H%d", rC), extrageFloat("ore_garda"))
			_ = f.SetCellValue(numeFoaieCentralizator, fmt.Sprintf("I%d", rC), extrageFloat("ZNL"))
			_ = f.SetCellValue(numeFoaieCentralizator, fmt.Sprintf("J%d", rC), extrageFloat("CO"))
			_ = f.SetCellValue(numeFoaieCentralizator, fmt.Sprintf("K%d", rC), extrageFloat("CM"))

			for col := 1; col <= 11; col++ {
				cell, _ := excelize.CoordinatesToCellName(col, rC)
				if col == 2 || col == 4 {
					_ = f.SetCellStyle(numeFoaieCentralizator, cell, cell, stilTextStanga)
				} else if col == 7 {
					_ = f.SetCellStyle(numeFoaieCentralizator, cell, cell, stilTotalBold)
				} else {
					_ = f.SetCellStyle(numeFoaieCentralizator, cell, cell, stilCelulaNorm)
				}
			}

			nrC++
			rC++
		}

		// Rând TOTAL GENERAL cu formule Excel =SUM(...)
		if rC > 5 {
			_ = f.SetCellValue(numeFoaieCentralizator, fmt.Sprintf("B%d", rC), "TOTAL GENERAL")
			_ = f.SetCellFormula(numeFoaieCentralizator, fmt.Sprintf("F%d", rC), fmt.Sprintf("SUM(F5:F%d)", rC-1))
			_ = f.SetCellFormula(numeFoaieCentralizator, fmt.Sprintf("G%d", rC), fmt.Sprintf("SUM(G5:G%d)", rC-1))
			_ = f.SetCellFormula(numeFoaieCentralizator, fmt.Sprintf("H%d", rC), fmt.Sprintf("SUM(H5:H%d)", rC-1))
			_ = f.SetCellFormula(numeFoaieCentralizator, fmt.Sprintf("I%d", rC), fmt.Sprintf("SUM(I5:I%d)", rC-1))
			_ = f.SetCellFormula(numeFoaieCentralizator, fmt.Sprintf("J%d", rC), fmt.Sprintf("SUM(J5:J%d)", rC-1))
			_ = f.SetCellFormula(numeFoaieCentralizator, fmt.Sprintf("K%d", rC), fmt.Sprintf("SUM(K5:K%d)", rC-1))

			for col := 1; col <= 11; col++ {
				cell, _ := excelize.CoordinatesToCellName(col, rC)
				_ = f.SetCellStyle(numeFoaieCentralizator, cell, cell, stilTotalBold)
			}
		}
	}

	buf := new(bytes.Buffer)
	if err := f.Write(buf); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// NumeFisierPontaj generează numele standard de fișier pentru exportul lunar.
func NumeFisierPontaj(an, luna int) string {
	return fmt.Sprintf("Pontix_Pontaj_%04d_%02d.xlsx", an, luna)
}
