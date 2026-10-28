package datebaza

import (
	"context"
	"strings"
	"testing"
)

func TestSabloaneImport(t *testing.T) {
	// Verifică generarea șablonului de angajați în format CSV și XLSX
	csvAng, numeCsvAng, err := GenereazaSablonAngajati("csv")
	if err != nil {
		t.Fatalf("eroare sablon csv angajati: %v", err)
	}
	if !strings.HasSuffix(numeCsvAng, ".csv") || len(csvAng) == 0 {
		t.Fatalf("sablon csv angajati invalid: nume=%s, lungime=%d", numeCsvAng, len(csvAng))
	}
	if !strings.Contains(string(csvAng), "Popescu") || !strings.Contains(string(csvAng), "Marca") {
		t.Fatalf("continut csv angajati incomplet: %s", string(csvAng))
	}

	xlsxAng, numeXlsxAng, err := GenereazaSablonAngajati("xlsx")
	if err != nil {
		t.Fatalf("eroare sablon xlsx angajati: %v", err)
	}
	if !strings.HasSuffix(numeXlsxAng, ".xlsx") || len(xlsxAng) == 0 {
		t.Fatalf("sablon xlsx angajati invalid: nume=%s, lungime=%d", numeXlsxAng, len(xlsxAng))
	}

	// Verifică generarea șablonului de sărbători
	csvSarb, numeCsvSarb, err := GenereazaSablonSarbatori("csv")
	if err != nil {
		t.Fatalf("eroare sablon csv sarbatori: %v", err)
	}
	if !strings.HasSuffix(numeCsvSarb, ".csv") || len(csvSarb) == 0 {
		t.Fatalf("sablon csv sarbatori invalid")
	}

	xlsxSarb, numeXlsxSarb, err := GenereazaSablonSarbatori("xlsx")
	if err != nil {
		t.Fatalf("eroare sablon xlsx sarbatori: %v", err)
	}
	if !strings.HasSuffix(numeXlsxSarb, ".xlsx") || len(xlsxSarb) == 0 {
		t.Fatalf("sablon xlsx sarbatori invalid")
	}
}

func TestAnalizaSiExecutieImportAngajati(t *testing.T) {
	m := pregateste(t)
	ctx := context.Background()

	// 1. Analizează fișierul șablon CSV generat
	csvContinut, _, err := GenereazaSablonAngajati("csv")
	if err != nil {
		t.Fatal(err)
	}

	rez, err := m.serviciu.AnalizeazaImportAngajati(ctx, csvContinut, "test_angajati.csv")
	if err != nil {
		t.Fatalf("eroare analiza import angajati: %v", err)
	}

	if rez.TotalRanduri != 3 || rez.TotalValide != 3 || rez.TotalInvalide != 0 || rez.TotalExistente != 0 {
		t.Fatalf("rezultat analiza neasteptat: %+v", rez)
	}

	// 2. Execută importul
	adaugati, err := m.serviciu.ExecutaImportAngajati(ctx, rez.Randuri, false)
	if err != nil {
		t.Fatalf("eroare executie import angajati: %v", err)
	}
	if adaugati != 3 {
		t.Fatalf("asteptat 3 adaugati, obtinut %d", adaugati)
	}

	// 3. Verifică datele salvate în baza de date
	date, err := m.serviciu.Citeste(ctx)
	if err != nil {
		t.Fatal(err)
	}

	gasitIon := false
	for _, a := range date.Angajati {
		if a.Marca == "M001" {
			gasitIon = true
			if a.Nume != "Popescu" || a.Prenume != "Ion" || a.TipContract != "INTREAGA" || a.NormaMinute != 420 {
				t.Fatalf("date incorecte pentru Ion: %+v", a)
			}
			if a.SectieID == "" {
				t.Fatalf("sectie lipsa: %+v", a)
			}
			if a.DepartamentID == "" {
				t.Fatalf("departament lipsa: %+v", a)
			}
		}
	}
	if !gasitIon {
		t.Fatal("angajatul M001 nu a fost gasit in baza de date")
	}

	// 4. Analizează din nou același fișier -> acum rândurile sunt marcate ca `ExistaDeja`
	rez2, err := m.serviciu.AnalizeazaImportAngajati(ctx, csvContinut, "test_angajati.csv")
	if err != nil {
		t.Fatal(err)
	}
	if rez2.TotalExistente != 3 {
		t.Fatalf("asteptat 3 existente, obtinut %d", rez2.TotalExistente)
	}

	// 5. Dacă executăm fără `actualizeazaExistente`, nu ar trebui să modifice nimic
	adaugati2, err := m.serviciu.ExecutaImportAngajati(ctx, rez2.Randuri, false)
	if err != nil {
		t.Fatal(err)
	}
	if adaugati2 != 0 {
		t.Fatalf("asteptat 0 modificati fara actualizare, obtinut %d", adaugati2)
	}

	// 6. Testează și fișierul XLSX
	xlsxContinut, _, err := GenereazaSablonAngajati("xlsx")
	if err != nil {
		t.Fatal(err)
	}
	rezXlsx, err := m.serviciu.AnalizeazaImportAngajati(ctx, xlsxContinut, "angajati.xlsx")
	if err != nil {
		t.Fatalf("eroare analiza xlsx: %v", err)
	}
	if rezXlsx.Format != "Excel" || rezXlsx.TotalRanduri != 3 {
		t.Fatalf("rezultat xlsx neasteptat: %+v", rezXlsx)
	}
}

func TestAnalizaSiExecutieImportSarbatori(t *testing.T) {
	m := pregateste(t)
	ctx := context.Background()

	// Fișier CSV sărbători cu format diferit de dată (ex. 15.08.2026 și 2026-11-30)
	csvData := []byte("Data;Denumire;Denumire_EN\n15.08.2026;Adormirea Maicii Domnului;Assumption of Mary\n2026-11-30;Sfântul Andrei;Saint Andrew\n")

	rez, err := m.serviciu.AnalizeazaImportSarbatori(ctx, csvData, "sarbatori.csv")
	if err != nil {
		t.Fatalf("eroare analiza sarbatori: %v", err)
	}
	if rez.TotalRanduri != 2 || rez.TotalValide != 2 {
		t.Fatalf("rezultat analiza sarbatori invalid: %+v", rez)
	}
	if rez.Randuri[0].Data != "2026-08-15" || rez.Randuri[1].Data != "2026-11-30" {
		t.Fatalf("conversie format data incorect: %+v", rez.Randuri)
	}

	salvate, err := m.serviciu.ExecutaImportSarbatori(ctx, rez.Randuri, true)
	if err != nil {
		t.Fatalf("eroare executie sarbatori: %v", err)
	}
	if salvate != 2 {
		t.Fatalf("asteptat 2 salvate, obtinut %d", salvate)
	}

	// Verifică în calendarul 2026
	cal, err := m.serviciu.CitesteCalendar(ctx, 2026)
	if err != nil {
		t.Fatal(err)
	}
	gasit := false
	for _, s := range cal.Sarbatori {
		if s.Data == "2026-08-15" {
			gasit = true
			break
		}
	}
	if !gasit {
		t.Fatal("sarbatoarea importata 2026-08-15 nu apare in calendar")
	}
}
