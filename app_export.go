package main

import (
	"encoding/base64"
	"os"
	"path/filepath"
	"strings"

	wruntime "github.com/wailsapp/wails/v2/pkg/runtime"

	"pontix/internal/datebaza"
	"pontix/internal/exportexcel"
	"pontix/internal/i18n"
)

type RaspunsExportExcel struct {
	NumeFisier string       `json:"nume_fisier"`
	Continut   string       `json:"continut"` // base64
	Cale       string       `json:"cale,omitempty"`
	Anulat     bool         `json:"anulat,omitempty"`
	Eroare     *MesajEroare `json:"eroare,omitempty"`
}

// ExportaPontajExcel generează registrul complet Excel (pontaj + centralizator salarii) și oferă opțiunea de salvare pe disc.
func (a *App) ExportaPontajExcel(an, luna int) RaspunsExportExcel {
	db := datebaza.Nou(a.db)
	date, err := db.Citeste(a.contextLucru())
	if err != nil {
		return RaspunsExportExcel{Eroare: eroareInterfata(err)}
	}
	schita, err := db.CitesteSchita(a.contextLucru(), an, luna)
	if err != nil {
		return RaspunsExportExcel{Eroare: eroareInterfata(err)}
	}
	hartaSarb := make(map[string]string)
	if cal, err := db.CitesteCalendar(a.contextLucru(), an); err == nil {
		for _, s := range cal.Sarbatori {
			hartaSarb[s.Data] = s.Denumire.In(i18n.Romana)
		}
	}

	bytesXlsx, err := exportexcel.GenereazaRegistruExcel(exportexcel.OptiuniExport{
		An:        an,
		Luna:      luna,
		Date:      date,
		Schita:    schita,
		Sarbatori: hartaSarb,
	})
	if err != nil {
		return RaspunsExportExcel{Eroare: eroareInterfata(err)}
	}

	numeFisier := exportexcel.NumeFisierPontaj(an, luna)
	continutB64 := base64.StdEncoding.EncodeToString(bytesXlsx)

	folder, _ := folderRapoarte()
	if a.ctx != nil {
		cale, err := wruntime.SaveFileDialog(a.ctx, wruntime.SaveDialogOptions{
			Title:            i18n.T(a.limba(), "pontaj.export_dialog_titlu", nil),
			DefaultDirectory: folder,
			DefaultFilename:  numeFisier,
			Filters: []wruntime.FileFilter{
				{DisplayName: i18n.T(a.limba(), "pontaj.export_filtru", nil), Pattern: "*.xlsx"},
			},
			CanCreateDirectories: true,
		})
		if err == nil && cale != "" {
			if !strings.EqualFold(filepath.Ext(cale), ".xlsx") {
				cale += ".xlsx"
			}
			if err := os.WriteFile(cale, bytesXlsx, 0600); err == nil {
				return RaspunsExportExcel{
					NumeFisier: numeFisier,
					Continut:   continutB64,
					Cale:       cale,
				}
			}
		} else if cale == "" && err == nil {
			return RaspunsExportExcel{
				NumeFisier: numeFisier,
				Continut:   continutB64,
				Anulat:     true,
			}
		}
	}

	return RaspunsExportExcel{
		NumeFisier: numeFisier,
		Continut:   continutB64,
	}
}
