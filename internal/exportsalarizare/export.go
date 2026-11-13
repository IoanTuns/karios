// Package exportsalarizare scrie orele și zilele din Pontaj într-un fișier CSV sau Excel, după
// un model configurabil (F10-T03, D108). Pontix nu calculează bani (D100): se exportă doar ore
// și zile.
package exportsalarizare

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/xuri/excelize/v2"

	"pontix/internal/calcul"
	"pontix/internal/i18n"
)

// VersiuneExportSalarizare: versiunea formatului (D009); crește când se schimbă conținutul.
const VersiuneExportSalarizare = "1.0.0"

const (
	FormatCSV  = "CSV"
	FormatXLSX = "XLSX"

	OreZecimal   = "ZECIMAL"
	OreOreMinute = "ORE_MINUTE"

	// Prefixul surselor care iau o coloană de total din Pontaj („total:OL”).
	PrefixTotal = "total:"
	// SursaConstanta: coloana are textul fix din Valoare.
	SursaConstanta = "CONSTANTA"
)

// SurseIdentitate: coloanele care nu sunt totaluri.
var SurseIdentitate = []string{"marca", "nume", "prenume", "nume_complet", "sectie", "departament", "functie", "an", "luna", "perioada"}

type Coloana struct {
	Sursa   string `json:"sursa"`
	Titlu   string `json:"titlu"`
	Valoare string `json:"valoare,omitempty"` // doar pentru CONSTANTA
}

// Model descrie un fișier de export.
type Model struct {
	ID               string    `json:"id"`
	Denumire         string    `json:"denumire"`
	Format           string    `json:"format"`
	Separator        string    `json:"separator"` // ";", "," sau "TAB" (CSV)
	BOM              bool      `json:"bom"`
	SeparatorZecimal string    `json:"separator_zecimal"` // "," sau "."
	Zecimale         int       `json:"zecimale"`
	FormaOre         string    `json:"forma_ore"`
	Antet            bool      `json:"antet"`
	Coloane          []Coloana `json:"coloane"`
	Activ            bool      `json:"activ"`
	Revizie          string    `json:"revizie"`
}

// Rand: un angajat, cu identitatea (marca, nume…) și totalurile lunii.
type Rand struct {
	Identitate map[string]string
	Valori     map[string]calcul.ValoarePontaj
}

func eroare(cod string, param ...string) error {
	e := &i18n.Eroare{Cod: cod}
	if len(param) == 2 {
		e.Param = map[string]string{param[0]: param[1]}
	}
	return e
}

func sursaCunoscuta(s string) bool {
	if s == SursaConstanta {
		return true
	}
	if cod, ok := strings.CutPrefix(s, PrefixTotal); ok {
		for _, c := range ColoaneTotal() {
			if c == cod {
				return true
			}
		}
		return false
	}
	for _, c := range SurseIdentitate {
		if c == s {
			return true
		}
	}
	return false
}

// ColoaneTotal: coloanele de total din Pontaj care se pot exporta, plus soldurile lunii.
func ColoaneTotal() []string {
	return append(calcul.ColoanePontaj(), "ore_lucratoare", "L_INCEPUT")
}

// Valideaza verifică modelul și completează valorile implicite.
func Valideaza(m *Model) error {
	m.Denumire = strings.TrimSpace(m.Denumire)
	if m.Denumire == "" || len([]rune(m.Denumire)) > 100 {
		return eroare("export_salarizare.denumire_invalida")
	}
	if m.Format != FormatCSV && m.Format != FormatXLSX {
		return eroare("export_salarizare.format_invalid")
	}
	if m.Separator == "" {
		m.Separator = ";"
	}
	if m.Separator != ";" && m.Separator != "," && m.Separator != "TAB" {
		return eroare("export_salarizare.separator_invalid")
	}
	if m.SeparatorZecimal == "" {
		m.SeparatorZecimal = ","
	}
	if m.SeparatorZecimal != "," && m.SeparatorZecimal != "." || (m.Format == FormatCSV && m.Separator == "," && m.SeparatorZecimal == ",") {
		return eroare("export_salarizare.zecimal_invalid")
	}
	if m.Zecimale < 0 || m.Zecimale > 4 {
		return eroare("export_salarizare.zecimale_invalide")
	}
	if m.FormaOre == "" {
		m.FormaOre = OreZecimal
	}
	if m.FormaOre != OreZecimal && m.FormaOre != OreOreMinute {
		return eroare("export_salarizare.forma_ore_invalida")
	}
	if len(m.Coloane) == 0 || len(m.Coloane) > 100 {
		return eroare("export_salarizare.fara_coloane")
	}
	for i := range m.Coloane {
		c := &m.Coloane[i]
		c.Titlu = strings.TrimSpace(c.Titlu)
		if !sursaCunoscuta(c.Sursa) {
			return eroare("export_salarizare.sursa_invalida", "sursa", c.Sursa)
		}
		if c.Sursa != SursaConstanta {
			c.Valoare = ""
		}
		if len([]rune(c.Titlu)) > 100 || len([]rune(c.Valoare)) > 100 {
			return eroare("export_salarizare.titlu_invalid")
		}
	}
	return nil
}

// numar scrie o valoare de total după model; ok=false pentru valorile necalculate.
func numar(m Model, v calcul.ValoarePontaj) (string, bool) {
	if !v.Definita || v.NuSeAplica {
		return "", false
	}
	if v.Unitate == "ORE" && m.FormaOre == OreOreMinute {
		minute := v.Sutimi / 100
		semn := ""
		if minute < 0 {
			semn, minute = "-", -minute
		}
		return fmt.Sprintf("%s%d:%02d", semn, minute/60, minute%60), true
	}
	divizor := 100.0
	if v.Unitate == "ORE" {
		divizor = 6000
	}
	// Rotunjire obisnuita (0,5 in sus), nu la par: 2,25 cu o zecimala da 2,3.
	scara := math.Pow(10, float64(m.Zecimale))
	text := strconv.FormatFloat(math.Round(float64(v.Sutimi)/divizor*scara)/scara, 'f', m.Zecimale, 64)
	if m.SeparatorZecimal == "," {
		text = strings.Replace(text, ".", ",", 1)
	}
	return text, true
}

// Tabel construiește rândurile fișierului (cu antet, dacă modelul îl cere) și numără valorile
// necalculate, lăsate goale.
func Tabel(m Model, randuri []Rand) ([][]string, int) {
	rez := [][]string{}
	if m.Antet {
		antet := make([]string, len(m.Coloane))
		for i, c := range m.Coloane {
			antet[i] = c.Titlu
		}
		rez = append(rez, antet)
	}
	necalculate := 0
	for _, r := range randuri {
		linie := make([]string, len(m.Coloane))
		for i, c := range m.Coloane {
			switch {
			case c.Sursa == SursaConstanta:
				linie[i] = c.Valoare
			case strings.HasPrefix(c.Sursa, PrefixTotal):
				v, exista := r.Valori[strings.TrimPrefix(c.Sursa, PrefixTotal)]
				text, ok := numar(m, v)
				// Valorile care nu se aplica randului (F08-T18) sunt goale fara sa fie „necalculate”.
				if !ok && exista && !v.NuSeAplica {
					necalculate++
				}
				linie[i] = text
			default:
				linie[i] = r.Identitate[c.Sursa]
			}
		}
		rez = append(rez, linie)
	}
	return rez, necalculate
}

// Genereaza scrie fișierul și întoarce conținutul, extensia și numărul valorilor necalculate.
func Genereaza(m Model, randuri []Rand) ([]byte, string, int, error) {
	tabel, necalculate := Tabel(m, randuri)
	if m.Format == FormatXLSX {
		f := excelize.NewFile()
		defer f.Close()
		foaie := f.GetSheetName(0)
		for r, linie := range tabel {
			for c, text := range linie {
				celula, err := excelize.CoordinatesToCellName(c+1, r+1)
				if err != nil {
					return nil, "", 0, err
				}
				// In Excel totalurile zecimale sunt numere, ca programul de salarizare sa nu le
				// citeasca drept text; ore:minute si celelalte coloane raman text.
				numeric := (r > 0 || !m.Antet) && text != "" && strings.HasPrefix(m.Coloane[c].Sursa, PrefixTotal) && !(m.FormaOre == OreOreMinute && strings.Contains(text, ":"))
				if numeric {
					if v, err := strconv.ParseFloat(strings.Replace(text, ",", ".", 1), 64); err == nil {
						if err := f.SetCellFloat(foaie, celula, v, m.Zecimale, 64); err != nil {
							return nil, "", 0, err
						}
						continue
					}
				}
				if err := f.SetCellStr(foaie, celula, text); err != nil {
					return nil, "", 0, err
				}
			}
		}
		var b bytes.Buffer
		if err := f.Write(&b); err != nil {
			return nil, "", 0, err
		}
		return b.Bytes(), ".xlsx", necalculate, nil
	}
	var b bytes.Buffer
	if m.BOM {
		b.WriteString("\xef\xbb\xbf")
	}
	w := csv.NewWriter(&b)
	w.Comma = ';'
	switch m.Separator {
	case ",":
		w.Comma = ','
	case "TAB":
		w.Comma = '\t'
	}
	w.UseCRLF = true
	if err := w.WriteAll(tabel); err != nil {
		return nil, "", 0, err
	}
	return b.Bytes(), ".csv", necalculate, nil
}
