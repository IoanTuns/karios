package main

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	wruntime "github.com/wailsapp/wails/v2/pkg/runtime"

	"pontix/internal/calcul"
	"pontix/internal/datebaza"
	"pontix/internal/exportsalarizare"
	"pontix/internal/i18n"
	"pontix/internal/identitate"
)

// F10-T03 (D108): exportul pentru salarizare după un model configurabil (CSV sau Excel).

type RaspunsModeleExport struct {
	Modele []exportsalarizare.Model `json:"modele"`
	Model  *exportsalarizare.Model  `json:"model,omitempty"`
	Surse  []string                 `json:"surse"` // sursele posibile ale coloanelor
	Eroare *MesajEroare             `json:"eroare,omitempty"`
}

type RaspunsExportSalarizare struct {
	Cale        string       `json:"cale,omitempty"`
	Randuri     int          `json:"randuri"`
	Necalculate int          `json:"necalculate"`
	Anulat      bool         `json:"anulat,omitempty"`
	Eroare      *MesajEroare `json:"eroare,omitempty"`
}

func surseExport() []string {
	rez := append([]string{}, exportsalarizare.SurseIdentitate...)
	for _, c := range exportsalarizare.ColoaneTotal() {
		rez = append(rez, exportsalarizare.PrefixTotal+c)
	}
	return append(rez, exportsalarizare.SursaConstanta)
}

// ListaModeleExportSalarizare întoarce modelele și sursele posibile ale coloanelor.
func (a *App) ListaModeleExportSalarizare() RaspunsModeleExport {
	r := RaspunsModeleExport{Surse: surseExport(), Modele: []exportsalarizare.Model{}}
	if a.db == nil {
		return r
	}
	m, err := datebaza.Nou(a.db).ListaModeleExport(a.contextLucru())
	if err != nil {
		r.Eroare = eroareInterfata(err)
		return r
	}
	r.Modele = m
	return r
}

// SalveazaModelExportSalarizare adaugă sau modifică un model (doar administratorul).
func (a *App) SalveazaModelExportSalarizare(m exportsalarizare.Model) RaspunsModeleExport {
	if e := a.doarAdministrator(); e != nil {
		r := a.ListaModeleExportSalarizare()
		r.Eroare = e
		return r
	}
	salvat, err := datebaza.Nou(a.db).SalveazaModelExport(a.contextLucru(), m)
	r := a.ListaModeleExportSalarizare()
	if err != nil {
		r.Eroare = eroareInterfata(err)
		return r
	}
	r.Model = &salvat
	return r
}

// randuriExportSalarizare: un rând pe marcă, pentru angajații activi sau cu date în lună,
// ordonat după marcă; sectieID gol = toate secțiile.
func (a *App) randuriExportSalarizare(an, luna int, sectieID string) ([]exportsalarizare.Rand, error) {
	db := datebaza.Nou(a.db)
	ctx := a.contextLucru()
	date, err := db.Citeste(ctx)
	if err != nil {
		return nil, err
	}
	schita, err := db.CitesteSchita(ctx, an, luna)
	if err != nil {
		return nil, err
	}
	l := a.limba()
	sectii, departamente := map[string]string{}, map[string]string{}
	for _, s := range date.Sectii {
		sectii[s.ID] = s.Denumire.In(l)
	}
	for _, d := range date.Departamente {
		departamente[d.ID] = d.Denumire.In(l)
	}
	totaluri := map[string]map[string]calcul.ValoarePontaj{}
	for _, t := range schita.Calcul.Totaluri {
		totaluri[t.AngajatID] = t.Valori
	}
	rez := []exportsalarizare.Rand{}
	angajati := append([]datebaza.Angajat{}, date.Angajati...)
	sort.Slice(angajati, func(i, j int) bool { return angajati[i].Marca < angajati[j].Marca })
	perioada := fmt.Sprintf("%04d-%02d", an, luna)
	for _, x := range angajati {
		valori, cuDate := totaluri[x.ID]
		if (!x.Activ && !cuDate) || (sectieID != "" && x.SectieID != sectieID) {
			continue
		}
		rez = append(rez, exportsalarizare.Rand{
			Identitate: map[string]string{
				"marca": x.Marca, "nume": x.Nume, "prenume": x.Prenume,
				"nume_complet": identitate.Nume(x.Nume, x.Prenume, x.Garda || x.AngajatPrincipalID != "", l),
				"sectie":       sectii[x.SectieID], "departament": departamente[x.DepartamentID], "functie": x.Functie,
				"an": strconv.Itoa(an), "luna": fmt.Sprintf("%02d", luna), "perioada": perioada,
			},
			Valori: valori,
		})
	}
	return rez, nil
}

// ExportaSalarizare scrie fișierul după model și îl salvează unde alege utilizatorul.
func (a *App) ExportaSalarizare(modelID string, an, luna int, sectieID string) RaspunsExportSalarizare {
	if a.db == nil || an < 1 || an > 9999 || luna < 1 || luna > 12 {
		return RaspunsExportSalarizare{Eroare: &MesajEroare{Cod: "date.eroare"}}
	}
	m, err := datebaza.Nou(a.db).ModelExport(a.contextLucru(), modelID)
	if err != nil {
		return RaspunsExportSalarizare{Eroare: eroareInterfata(err)}
	}
	if err := exportsalarizare.Valideaza(&m); err != nil {
		return RaspunsExportSalarizare{Eroare: eroareInterfata(err)}
	}
	randuri, err := a.randuriExportSalarizare(an, luna, sectieID)
	if err != nil {
		return RaspunsExportSalarizare{Eroare: eroareInterfata(err)}
	}
	continut, ext, necalculate, err := exportsalarizare.Genereaza(m, randuri)
	if err != nil {
		return RaspunsExportSalarizare{Eroare: eroareInterfata(err)}
	}
	r := RaspunsExportSalarizare{Randuri: len(randuri), Necalculate: necalculate}
	if a.ctx == nil {
		return r
	}
	folder, _ := folderRapoarte()
	nume := fmt.Sprintf("Pontix_Salarizare_%04d_%02d%s", an, luna, ext)
	cale, err := wruntime.SaveFileDialog(a.ctx, wruntime.SaveDialogOptions{
		Title:                i18n.T(a.limba(), "export_salarizare.dialog_titlu", nil),
		DefaultDirectory:     folder,
		DefaultFilename:      nume,
		Filters:              []wruntime.FileFilter{{DisplayName: strings.ToUpper(strings.TrimPrefix(ext, ".")), Pattern: "*" + ext}},
		CanCreateDirectories: true,
	})
	if err != nil {
		return RaspunsExportSalarizare{Eroare: eroareInterfata(err)}
	}
	if cale == "" {
		r.Anulat = true
		return r
	}
	if !strings.EqualFold(filepath.Ext(cale), ext) {
		cale += ext
	}
	if err := os.WriteFile(cale, continut, 0o600); err != nil {
		return RaspunsExportSalarizare{Eroare: eroareInterfata(err)}
	}
	r.Cale = cale
	return r
}
