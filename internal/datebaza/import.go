package datebaza

import (
	"bytes"
	"context"
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"

	"github.com/xuri/excelize/v2"

	"pontix/internal/bd/acces"
	"pontix/internal/i18n"
)

type RandImportAngajat struct {
	Nume           string `json:"nume"`
	Prenume        string `json:"prenume"`
	Marca          string `json:"marca"`
	Functie        string `json:"functie"`
	Sectie         string `json:"sectie"`
	Departament    string `json:"departament"`
	TipContract    string `json:"tip_contract"`
	NormaMinute    int64  `json:"norma_minute"`
	ConcediuZileAn *int64 `json:"concediu_zile_an"`
	Valid          bool   `json:"valid"`
	MesajEroare    string `json:"mesaj_eroare"`
	ExistaDeja     bool   `json:"exista_deja"`
}

type RezultatAnalizaImportAngajati struct {
	NumeFisier     string              `json:"nume_fisier"`
	Format         string              `json:"format"`
	Randuri        []RandImportAngajat `json:"randuri"`
	TotalRanduri   int                 `json:"total_randuri"`
	TotalValide    int                 `json:"total_valide"`
	TotalInvalide  int                 `json:"total_invalide"`
	TotalExistente int                 `json:"total_existente"`
}

type CerereExecutaImportAngajati struct {
	Randuri               []RandImportAngajat `json:"randuri"`
	ActualizeazaExistente bool                `json:"actualizeaza_existente"`
}

type RandImportSarbatoare struct {
	Data        string `json:"data"`
	DenumireRO  string `json:"denumire_ro"`
	DenumireEN  string `json:"denumire_en"`
	Valid       bool   `json:"valid"`
	MesajEroare string `json:"mesaj_eroare"`
	ExistaDeja  bool   `json:"exista_deja"`
}

type RezultatAnalizaImportSarbatori struct {
	NumeFisier     string                 `json:"nume_fisier"`
	Format         string                 `json:"format"`
	Randuri        []RandImportSarbatoare `json:"randuri"`
	TotalRanduri   int                    `json:"total_randuri"`
	TotalValide    int                    `json:"total_valide"`
	TotalInvalide  int                    `json:"total_invalide"`
	TotalExistente int                    `json:"total_existente"`
}

type CerereExecutaImportSarbatori struct {
	Randuri              []RandImportSarbatoare `json:"randuri"`
	InlocuiesteExistente bool                   `json:"inlocuieste_existente"`
}

func curataAntet(s string) string {
	s = strings.ToLower(strings.TrimSpace(s))
	repl := strings.NewReplacer(" ", "", "_", "", "-", "", "ă", "a", "â", "a", "î", "i", "ș", "s", "ț", "t", "ş", "s", "ţ", "t")
	return repl.Replace(s)
}

func citesteLiniiFisier(continut []byte, numeFisier string) ([][]string, error) {
	ext := strings.ToLower(numeFisier)
	if strings.HasSuffix(ext, ".xlsx") {
		f, err := excelize.OpenReader(bytes.NewReader(continut))
		if err != nil {
			return nil, fmt.Errorf("eroare deschidere Excel: %w", err)
		}
		defer f.Close()
		foi := f.GetSheetList()
		if len(foi) == 0 {
			return nil, errors.New("fisierul Excel nu contine nicio foaie")
		}
		return f.GetRows(foi[0])
	}

	// CSV
	curat := bytes.TrimPrefix(continut, []byte("\xef\xbb\xbf")) // Elimina UTF-8 BOM daca exista
	var sep rune = ','
	if linie, _, _ := bytes.Cut(curat, []byte("\n")); bytes.Count(linie, []byte(";")) > bytes.Count(linie, []byte(",")) {
		sep = ';'
	}
	r := csv.NewReader(bytes.NewReader(curat))
	r.Comma = sep
	r.FieldsPerRecord = -1
	r.LazyQuotes = true
	var randuri [][]string
	for {
		linie, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("eroare citire CSV: %w", err)
		}
		randuri = append(randuri, linie)
	}
	return randuri, nil
}

// AnalizeazaImportAngajati parsează fișierul și verifică validitatea și duplicatele.
func (s *Serviciu) AnalizeazaImportAngajati(ctx context.Context, continut []byte, numeFisier string) (RezultatAnalizaImportAngajati, error) {
	linii, err := citesteLiniiFisier(continut, numeFisier)
	if err != nil {
		return RezultatAnalizaImportAngajati{}, err
	}
	if len(linii) == 0 {
		return RezultatAnalizaImportAngajati{}, errors.New("fisierul este gol")
	}

	antet := linii[0]
	colNume, colPrenume, colMarca := -1, -1, -1
	colFunctie, colSectie, colDepartament := -1, -1, -1
	colContract, colNorma, colConcediu := -1, -1, -1

	for i, h := range antet {
		c := curataAntet(h)
		switch {
		case c == "nume" || c == "numefamilie" || c == "lastname":
			colNume = i
		case c == "prenume" || c == "firstname":
			colPrenume = i
		case c == "marca" || c == "nrmarca" || c == "cod" || c == "id" || c == "employeeid":
			colMarca = i
		case c == "functie" || c == "post" || c == "ocupatie" || c == "role" || c == "job":
			colFunctie = i
		case c == "sectie" || c == "sectia" || c == "unit":
			colSectie = i
		case c == "departament" || c == "subdepartament" || c == "department" || c == "birou":
			colDepartament = i
		case c == "contract" || c == "tipcontract":
			colContract = i
		case c == "norma" || c == "normaore" || c == "normaminute" || c == "orezi" || c == "hours":
			colNorma = i
		case c == "concediu" || c == "concediuan" || c == "concediuzile" || c == "leave" || c == "co":
			colConcediu = i
		}
	}

	if colMarca == -1 && colNume == -1 {
		return RezultatAnalizaImportAngajati{}, errors.New("antet invalid: nu au fost identificate coloanele 'Nume' sau 'Marca'")
	}

	// Citim angajații existenți
	existenti, err := acces.New(s.db).ListaAngajati(ctx)
	if err != nil {
		return RezultatAnalizaImportAngajati{}, err
	}
	hartaExistenti := map[string]bool{}
	for _, a := range existenti {
		hartaExistenti[strings.ToUpper(strings.TrimSpace(a.Marca))] = true
	}

	marciFisier := map[string]bool{}
	var randuri []RandImportAngajat

	getVal := func(row []string, idx int) string {
		if idx >= 0 && idx < len(row) {
			return strings.TrimSpace(row[idx])
		}
		return ""
	}

	for rIdx := 1; rIdx < len(linii); rIdx++ {
		linia := linii[rIdx]
		// Ignoră rândurile complet goale
		completGol := true
		for _, cell := range linia {
			if strings.TrimSpace(cell) != "" {
				completGol = false
				break
			}
		}
		if completGol {
			continue
		}

		nume := getVal(linia, colNume)
		prenume := getVal(linia, colPrenume)
		marca := getVal(linia, colMarca)
		functie := getVal(linia, colFunctie)
		sectie := getVal(linia, colSectie)
		departament := getVal(linia, colDepartament)
		tipContractStr := strings.ToUpper(getVal(linia, colContract))
		normaStr := getVal(linia, colNorma)
		concediuStr := getVal(linia, colConcediu)

		tipContract := "INTREAGA"
		if strings.Contains(tipContractStr, "PART") {
			tipContract = "PARTIALA"
		}

		normaMinute := int64(420) // implicit 7h = 420 min
		if normaStr != "" {
			normaStrCurata := strings.TrimSuffix(strings.TrimSuffix(normaStr, "h"), "H")
			normaStrCurata = strings.ReplaceAll(normaStrCurata, ",", ".")
			if oreVal, err := strconv.ParseFloat(normaStrCurata, 64); err == nil && oreVal > 0 && oreVal <= 24 {
				normaMinute = int64(oreVal * 60)
			}
		} else if tipContract == "PARTIALA" {
			normaMinute = 240 // implicit 4 ore pentru parțială dacă nu e precizat
		}

		var concediuPtr *int64
		if concediuStr != "" {
			if cz, err := strconv.ParseInt(concediuStr, 10, 64); err == nil && cz >= 0 && cz <= 365 {
				concediuPtr = &cz
			}
		}

		r := RandImportAngajat{
			Nume:           nume,
			Prenume:        prenume,
			Marca:          marca,
			Functie:        functie,
			Sectie:         sectie,
			Departament:    departament,
			TipContract:    tipContract,
			NormaMinute:    normaMinute,
			ConcediuZileAn: concediuPtr,
			Valid:          true,
		}

		marcaUpper := strings.ToUpper(marca)
		if marca == "" {
			r.Valid = false
			r.MesajEroare = "Marca lipsește"
		} else if nume == "" {
			r.Valid = false
			r.MesajEroare = "Numele lipsește"
		} else if marciFisier[marcaUpper] {
			r.Valid = false
			r.MesajEroare = "Marcă duplicat în fișier"
		} else {
			marciFisier[marcaUpper] = true
			if hartaExistenti[marcaUpper] {
				r.ExistaDeja = true
			}
		}

		randuri = append(randuri, r)
	}

	rez := RezultatAnalizaImportAngajati{
		NumeFisier: numeFisier,
		Format:     "CSV",
		Randuri:    randuri,
	}
	if strings.HasSuffix(strings.ToLower(numeFisier), ".xlsx") {
		rez.Format = "Excel"
	}
	rez.TotalRanduri = len(randuri)
	for _, r := range randuri {
		if r.Valid {
			rez.TotalValide++
			if r.ExistaDeja {
				rez.TotalExistente++
			}
		} else {
			rez.TotalInvalide++
		}
	}

	return rez, nil
}

// ExecutaImportAngajati execută importul efectiv într-o tranzacție atomică.
func (s *Serviciu) ExecutaImportAngajati(ctx context.Context, randuri []RandImportAngajat, actualizeazaExistente bool) (int, error) {
	if len(randuri) == 0 {
		return 0, nil
	}

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return 0, traduceEroare(err)
	}
	defer tx.Rollback()

	q := acces.New(tx)
	acum := moment()

	// Încarcă secțiile și departamentele existente
	sectiiExistente, err := q.ListaSectii(ctx)
	if err != nil {
		return 0, traduceEroare(err)
	}
	hartaSectii := map[string]string{} // denumire normalizata -> ID
	for _, sc := range sectiiExistente {
		var den i18n.Text
		_ = json.Unmarshal([]byte(sc.Denumire), &den)
		hartaSectii[curataAntet(den.In(i18n.Romana))] = sc.ID
		hartaSectii[curataAntet(den.In(i18n.Engleza))] = sc.ID
	}

	deptExistente, err := q.ListaDepartamente(ctx)
	if err != nil {
		return 0, traduceEroare(err)
	}
	hartaDept := map[string]string{} // sectieID + "/" + denumire normalizata -> ID
	for _, dp := range deptExistente {
		var den i18n.Text
		_ = json.Unmarshal([]byte(dp.Denumire), &den)
		hartaDept[dp.SectieID+"/"+curataAntet(den.In(i18n.Romana))] = dp.ID
		hartaDept[dp.SectieID+"/"+curataAntet(den.In(i18n.Engleza))] = dp.ID
	}

	// Încarcă angajații existenți
	angajatiExistenti, err := q.ListaAngajati(ctx)
	if err != nil {
		return 0, traduceEroare(err)
	}
	hartaAngajati := map[string]acces.ListaAngajatiRow{} // marca normalizata -> Angajat
	for _, a := range angajatiExistenti {
		hartaAngajati[strings.ToUpper(strings.TrimSpace(a.Marca))] = a
	}

	adaugatiSauModificati := 0

	for _, r := range randuri {
		if !r.Valid {
			continue
		}

		// Rezolvă secția
		var sectieID string
		if r.Sectie != "" {
			cheieSectie := curataAntet(r.Sectie)
			if id, ok := hartaSectii[cheieSectie]; ok {
				sectieID = id
			} else {
				// Creează secție nouă
				sectieID = idNou()
				denJSON, _ := json.Marshal(i18n.Text{"ro": r.Sectie, "en": r.Sectie})
				if err := q.AdaugaSectie(ctx, acces.AdaugaSectieParams{
					ID:          sectieID,
					Denumire:    string(denJSON),
					CreatLa:     acum,
					ModificatLa: acum,
				}); err != nil {
					return 0, traduceEroare(err)
				}
				hartaSectii[cheieSectie] = sectieID
			}
		}

		// Rezolvă departamentul
		var deptID string
		if r.Departament != "" && sectieID != "" {
			cheieDept := sectieID + "/" + curataAntet(r.Departament)
			if id, ok := hartaDept[cheieDept]; ok {
				deptID = id
			} else {
				// Creează departament nou
				deptID = idNou()
				denJSON, _ := json.Marshal(i18n.Text{"ro": r.Departament, "en": r.Departament})
				if err := q.AdaugaDepartament(ctx, acces.AdaugaDepartamentParams{
					ID:          deptID,
					SectieID:    sectieID,
					Denumire:    string(denJSON),
					CreatLa:     acum,
					ModificatLa: acum,
				}); err != nil {
					return 0, traduceEroare(err)
				}
				hartaDept[cheieDept] = deptID
			}
		}

		marcaUpper := strings.ToUpper(strings.TrimSpace(r.Marca))
		if exist, ok := hartaAngajati[marcaUpper]; ok {
			if !actualizeazaExistente {
				continue
			}
			// Actualizează angajat existent
			functie := r.Functie
			if functie == "" && exist.Functie != nil {
				functie = *exist.Functie
			}
			finalSectieID := sectieID
			if finalSectieID == "" && exist.SectieID != nil {
				finalSectieID = *exist.SectieID
			}
			finalDeptID := deptID
			if finalDeptID == "" && exist.DepartamentID != nil {
				finalDeptID = *exist.DepartamentID
			}
			var normaPartiala *int64
			if r.TipContract == "PARTIALA" {
				normaPartiala = &r.NormaMinute
			}
			concediu := r.ConcediuZileAn
			if concediu == nil {
				concediu = exist.ConcediuZileAn
			}

			if _, err := q.ModificaAngajat(ctx, acces.ModificaAngajatParams{
				ID:                  exist.ID,
				Nume:                r.Nume,
				Prenume:             r.Prenume,
				Marca:               r.Marca,
				Functie:             textOptional(functie),
				SectieID:            textOptional(finalSectieID),
				DepartamentID:       textOptional(finalDeptID),
				CategorieProgramID:  exist.CategorieProgramID,
				TipContract:         r.TipContract,
				NormaPartialaMinute: normaPartiala,
				ConcediuZileAn:      concediu,
				Activ:               exist.Activ,
				ModificatLa:         acum,
			}); err != nil {
				return 0, traduceEroare(err)
			}
			adaugatiSauModificati++
		} else {
			// Adaugă angajat nou
			id := idNou()
			var normaPartiala *int64
			if r.TipContract == "PARTIALA" {
				normaPartiala = &r.NormaMinute
			}
			if err := q.AdaugaAngajat(ctx, acces.AdaugaAngajatParams{
				ID:                  id,
				Nume:                r.Nume,
				Prenume:             r.Prenume,
				Marca:               r.Marca,
				Functie:             textOptional(r.Functie),
				SectieID:            textOptional(sectieID),
				DepartamentID:       textOptional(deptID),
				TipContract:         r.TipContract,
				NormaPartialaMinute: normaPartiala,
				ConcediuZileAn:      r.ConcediuZileAn,
				Activ:               1,
				CreatLa:             acum,
				ModificatLa:         acum,
			}); err != nil {
				return 0, traduceEroare(err)
			}
			hartaAngajati[marcaUpper] = acces.ListaAngajatiRow{ID: id, Marca: r.Marca}
			adaugatiSauModificati++
		}
	}

	if err := tx.Commit(); err != nil {
		return 0, traduceEroare(err)
	}

	return adaugatiSauModificati, nil
}

// AnalizeazaImportSarbatori parsează fișierul de sărbători legale.
func (s *Serviciu) AnalizeazaImportSarbatori(ctx context.Context, continut []byte, numeFisier string) (RezultatAnalizaImportSarbatori, error) {
	linii, err := citesteLiniiFisier(continut, numeFisier)
	if err != nil {
		return RezultatAnalizaImportSarbatori{}, err
	}
	if len(linii) == 0 {
		return RezultatAnalizaImportSarbatori{}, errors.New("fisierul este gol")
	}

	antet := linii[0]
	colData, colRO, colEN := -1, -1, -1
	for i, h := range antet {
		c := curataAntet(h)
		switch {
		case c == "data" || c == "date" || c == "ziua":
			colData = i
		case c == "denumire" || c == "denumirero" || c == "nume" || c == "sarbatoare":
			colRO = i
		case c == "denumireen" || c == "en" || c == "english":
			colEN = i
		}
	}
	if colData == -1 && colRO == -1 {
		return RezultatAnalizaImportSarbatori{}, errors.New("antet invalid: nu au fost identificate coloanele 'Data' sau 'Denumire'")
	}

	// Citim toate sărbătorile existente
	existente, err := acces.New(s.db).ListaSarbatori(ctx, nil)
	if err != nil {
		return RezultatAnalizaImportSarbatori{}, err
	}
	hartaExistente := map[string]bool{}
	for _, z := range existente {
		hartaExistente[z.Data] = true
	}

	dateFisier := map[string]bool{}
	var randuri []RandImportSarbatoare

	getVal := func(row []string, idx int) string {
		if idx >= 0 && idx < len(row) {
			return strings.TrimSpace(row[idx])
		}
		return ""
	}

	for rIdx := 1; rIdx < len(linii); rIdx++ {
		linia := linii[rIdx]
		completGol := true
		for _, cell := range linia {
			if strings.TrimSpace(cell) != "" {
				completGol = false
				break
			}
		}
		if completGol {
			continue
		}

		rawDate := getVal(linia, colData)
		denRO := getVal(linia, colRO)
		denEN := getVal(linia, colEN)
		if denEN == "" {
			denEN = denRO
		}

		stdDate := formateazaData(rawDate)
		r := RandImportSarbatoare{
			Data:       stdDate,
			DenumireRO: denRO,
			DenumireEN: denEN,
			Valid:      true,
		}

		if stdDate == "" {
			r.Valid = false
			r.MesajEroare = "Format dată invalid (așteptat AAAA-LL-ZZ sau ZZ.LL.AAAA)"
		} else if denRO == "" {
			r.Valid = false
			r.MesajEroare = "Denumirea lipsește"
		} else if dateFisier[stdDate] {
			r.Valid = false
			r.MesajEroare = "Dată duplicată în fișier"
		} else {
			dateFisier[stdDate] = true
			if hartaExistente[stdDate] {
				r.ExistaDeja = true
			}
		}

		randuri = append(randuri, r)
	}

	rez := RezultatAnalizaImportSarbatori{
		NumeFisier: numeFisier,
		Format:     "CSV",
		Randuri:    randuri,
	}
	if strings.HasSuffix(strings.ToLower(numeFisier), ".xlsx") {
		rez.Format = "Excel"
	}
	rez.TotalRanduri = len(randuri)
	for _, r := range randuri {
		if r.Valid {
			rez.TotalValide++
			if r.ExistaDeja {
				rez.TotalExistente++
			}
		} else {
			rez.TotalInvalide++
		}
	}

	return rez, nil
}

func formateazaData(d string) string {
	d = strings.TrimSpace(d)
	for _, format := range []string{time.DateOnly, "02.01.2006", "02/01/2006", "02-01-2006", "2006/01/02"} {
		if t, err := time.Parse(format, d); err == nil && t.Year() >= 1900 && t.Year() <= 2100 {
			return t.Format(time.DateOnly)
		}
	}
	return ""
}

// ExecutaImportSarbatori inserează sărbătorile în baza de date.
func (s *Serviciu) ExecutaImportSarbatori(ctx context.Context, randuri []RandImportSarbatoare, inlocuiesteExistente bool) (int, error) {
	if len(randuri) == 0 {
		return 0, nil
	}

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return 0, traduceEroare(err)
	}
	defer tx.Rollback()

	q := acces.New(tx)
	acum := moment()

	existente, err := q.ListaSarbatori(ctx, nil)
	if err != nil {
		return 0, traduceEroare(err)
	}
	hartaExistente := map[string]bool{}
	for _, z := range existente {
		hartaExistente[z.Data] = true
	}

	salvate := 0
	for _, r := range randuri {
		if !r.Valid {
			continue
		}
		if hartaExistente[r.Data] && !inlocuiesteExistente {
			continue
		}
		denJSON, _ := json.Marshal(i18n.Text{"ro": r.DenumireRO, "en": r.DenumireEN})
		if hartaExistente[r.Data] {
			if _, err := q.ModificaSarbatoare(ctx, acces.ModificaSarbatoareParams{
				DataNoua:     r.Data,
				DataInitiala: r.Data,
				Denumire:     string(denJSON),
				ModificatLa:  acum,
			}); err != nil {
				return 0, traduceEroare(err)
			}
		} else {
			if _, err := q.AdaugaSarbatoare(ctx, acces.AdaugaSarbatoareParams{
				Data:        r.Data,
				Denumire:    string(denJSON),
				CreatLa:     acum,
				ModificatLa: acum,
			}); err != nil {
				return 0, traduceEroare(err)
			}
			hartaExistente[r.Data] = true
		}
		salvate++
	}

	if err := tx.Commit(); err != nil {
		return 0, traduceEroare(err)
	}

	return salvate, nil
}

// GenereazaSablonAngajati creează fișierul demonstrativ .xlsx sau .csv.
func GenereazaSablonAngajati(format string) ([]byte, string, error) {
	if strings.ToLower(format) == "xlsx" {
		f := excelize.NewFile()
		defer f.Close()
		sheet := "Angajati"
		f.SetSheetName("Sheet1", sheet)
		headers := []string{"Nume", "Prenume", "Marca", "Functie", "Sectie", "Departament", "Tip_Contract", "Norma_Ore_Zi", "Zile_Concediu_An"}
		for colIdx, h := range headers {
			cell, _ := excelize.CoordinatesToCellName(colIdx+1, 1)
			_ = f.SetCellValue(sheet, cell, h)
		}
		mostre := [][]any{
			{"Popescu", "Ion", "M001", "Medic specialist", "Laborator", "Biochimie", "INTREAGA", 7, 22},
			{"Ionescu", "Maria", "M002", "Asistent medical", "Laborator", "Hematologie", "INTREAGA", 7, 21},
			{"Radu", "Elena", "M003", "Biolog", "Radiologie", "Imagistică", "PARTIALA", 4, 21},
		}
		for rIdx, r := range mostre {
			for cIdx, v := range r {
				cell, _ := excelize.CoordinatesToCellName(cIdx+1, rIdx+2)
				_ = f.SetCellValue(sheet, cell, v)
			}
		}
		buf, err := f.WriteToBuffer()
		if err != nil {
			return nil, "", err
		}
		return buf.Bytes(), "sablon_angajati.xlsx", nil
	}

	// CSV cu UTF-8 BOM
	var buf bytes.Buffer
	buf.Write([]byte("\xef\xbb\xbf")) // BOM
	w := csv.NewWriter(&buf)
	w.Comma = ';'
	_ = w.Write([]string{"Nume", "Prenume", "Marca", "Functie", "Sectie", "Departament", "Tip_Contract", "Norma_Ore_Zi", "Zile_Concediu_An"})
	_ = w.Write([]string{"Popescu", "Ion", "M001", "Medic specialist", "Laborator", "Biochimie", "INTREAGA", "7", "22"})
	_ = w.Write([]string{"Ionescu", "Maria", "M002", "Asistent medical", "Laborator", "Hematologie", "INTREAGA", "7", "21"})
	_ = w.Write([]string{"Radu", "Elena", "M003", "Biolog", "Radiologie", "Imagistică", "PARTIALA", "4", "21"})
	w.Flush()
	return buf.Bytes(), "sablon_angajati.csv", nil
}

// GenereazaSablonSarbatori creează fișierul demonstrativ pentru sărbători legale.
func GenereazaSablonSarbatori(format string) ([]byte, string, error) {
	if strings.ToLower(format) == "xlsx" {
		f := excelize.NewFile()
		defer f.Close()
		sheet := "Sarbatori"
		f.SetSheetName("Sheet1", sheet)
		headers := []string{"Data", "Denumire_RO", "Denumire_EN"}
		for colIdx, h := range headers {
			cell, _ := excelize.CoordinatesToCellName(colIdx+1, 1)
			_ = f.SetCellValue(sheet, cell, h)
		}
		mostre := [][]any{
			{"2026-01-01", "Anul Nou", "New Year's Day"},
			{"2026-01-02", "A doua zi de Anul Nou", "Day after New Year's Day"},
			{"2026-05-01", "Ziua Muncii", "Labour Day"},
			{"2026-12-25", "Crăciunul", "Christmas Day"},
			{"2026-12-26", "A doua zi de Crăciun", "Second day of Christmas"},
		}
		for rIdx, r := range mostre {
			for cIdx, v := range r {
				cell, _ := excelize.CoordinatesToCellName(cIdx+1, rIdx+2)
				_ = f.SetCellValue(sheet, cell, v)
			}
		}
		buf, err := f.WriteToBuffer()
		if err != nil {
			return nil, "", err
		}
		return buf.Bytes(), "sablon_sarbatori.xlsx", nil
	}

	var buf bytes.Buffer
	buf.Write([]byte("\xef\xbb\xbf"))
	w := csv.NewWriter(&buf)
	w.Comma = ';'
	_ = w.Write([]string{"Data", "Denumire_RO", "Denumire_EN"})
	_ = w.Write([]string{"2026-01-01", "Anul Nou", "New Year's Day"})
	_ = w.Write([]string{"2026-01-02", "A doua zi de Anul Nou", "Day after New Year's Day"})
	_ = w.Write([]string{"2026-05-01", "Ziua Muncii", "Labour Day"})
	_ = w.Write([]string{"2026-12-25", "Crăciunul", "Christmas Day"})
	w.Flush()
	return buf.Bytes(), "sablon_sarbatori.csv", nil
}
