package main

import (
	"bytes"
	"encoding/csv"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	wruntime "github.com/wailsapp/wails/v2/pkg/runtime"

	"pontix/internal/datebaza"
	"pontix/internal/i18n"
)

// RaspunsJurnal întoarce o pagină din jurnalul de modificări (F12-T02).
type RaspunsJurnal struct {
	Pagina datebaza.PaginaJurnal `json:"pagina"`
	Eroare *MesajEroare          `json:"eroare,omitempty"`
}

// RaspunsExportJurnal spune unde s-a salvat exportul.
type RaspunsExportJurnal struct {
	Cale   string       `json:"cale,omitempty"`
	Numar  int          `json:"numar"`
	Anulat bool         `json:"anulat,omitempty"`
	Eroare *MesajEroare `json:"eroare,omitempty"`
}

func (a *App) CitesteJurnal(filtru datebaza.FiltruJurnal) RaspunsJurnal {
	p, err := datebaza.Nou(a.db).CitesteJurnal(a.contextLucru(), filtru)
	if err != nil {
		return RaspunsJurnal{Eroare: eroareInterfata(err)}
	}
	return RaspunsJurnal{Pagina: p}
}

// ExportaJurnalCSV salvează intrările filtrate într-un fișier CSV ales de utilizator.
// Separatorul e „;” și fișierul începe cu BOM, ca Excel să-l deschidă corect.
func (a *App) ExportaJurnalCSV(filtru datebaza.FiltruJurnal, limba string) RaspunsExportJurnal {
	intrari, err := datebaza.Nou(a.db).TotJurnalul(a.contextLucru(), filtru)
	if err != nil {
		return RaspunsExportJurnal{Eroare: eroareInterfata(err)}
	}
	continut, err := csvJurnal(intrari, i18n.Limba(limba))
	if err != nil {
		return RaspunsExportJurnal{Eroare: eroareInterfata(err)}
	}
	if a.ctx == nil {
		return RaspunsExportJurnal{Anulat: true}
	}
	folder, _ := folderRapoarte()
	cale, err := wruntime.SaveFileDialog(a.ctx, wruntime.SaveDialogOptions{
		Title:                i18n.T(i18n.Limba(limba), "jurnal.export_titlu", nil),
		DefaultDirectory:     folder,
		DefaultFilename:      "jurnal-modificari-" + time.Now().Format("2006-01-02") + ".csv",
		Filters:              []wruntime.FileFilter{{DisplayName: "CSV (*.csv)", Pattern: "*.csv"}},
		CanCreateDirectories: true,
	})
	if err != nil {
		return RaspunsExportJurnal{Eroare: eroareInterfata(err)}
	}
	if cale == "" {
		return RaspunsExportJurnal{Anulat: true}
	}
	if !strings.EqualFold(filepath.Ext(cale), ".csv") {
		cale += ".csv"
	}
	if err := os.WriteFile(cale, continut, 0o600); err != nil {
		return RaspunsExportJurnal{Eroare: &MesajEroare{Cod: "jurnal.export_esuat", Param: map[string]string{"detalii": err.Error()}}}
	}
	return RaspunsExportJurnal{Cale: cale, Numar: len(intrari)}
}

func csvJurnal(intrari []datebaza.IntrareJurnal, limba i18n.Limba) ([]byte, error) {
	var b bytes.Buffer
	b.WriteString(bomUTF8)
	w := csv.NewWriter(&b)
	w.Comma = ';'
	antet := []string{"id", "cand", "cine", "tabel", "rand", "camp", "vechi", "nou"}
	for i, c := range antet {
		antet[i] = i18n.T(limba, "jurnal.coloana_"+c, nil)
	}
	if err := w.Write(antet); err != nil {
		return nil, err
	}
	for _, x := range intrari {
		if err := w.Write([]string{strconv.FormatInt(x.ID, 10), x.Cand, x.Cine, x.Tabel, x.RandID, x.Camp, x.Vechi, x.Nou}); err != nil {
			return nil, err
		}
	}
	w.Flush()
	return b.Bytes(), w.Error()
}

// bomUTF8 marchează fișierul ca UTF-8 pentru Excel.
const bomUTF8 = "\xef\xbb\xbf"
