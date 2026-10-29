package main

import (
	"bytes"
	"encoding/csv"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	wruntime "github.com/wailsapp/wails/v2/pkg/runtime"
	"pontix/internal/ai"
	"pontix/internal/datebaza"
	"pontix/internal/i18n"
)

type RaspunsInstrumenteAI struct {
	Anulare    *datebaza.AnulareAI          `json:"anulare,omitempty"`
	Sabloane   []datebaza.SablonAI          `json:"sabloane"`
	Verificare *datebaza.VerificarePontajAI `json:"verificare,omitempty"`
	Eroare     *MesajEroare                 `json:"eroare,omitempty"`
}

func (a *App) PrevizualizeazaAnulareAI() RaspunsInstrumenteAI {
	r, err := datebaza.Nou(a.db).PrevizualizeazaAnulareAI(a.contextLucru())
	if err != nil {
		return RaspunsInstrumenteAI{Eroare: eroareInterfata(err)}
	}
	return RaspunsInstrumenteAI{Anulare: &r}
}
func (a *App) AnuleazaUltimaOperatieAI(token string) RaspunsInstrumenteAI {
	r, err := datebaza.Nou(a.db).AnuleazaUltimaOperatieAI(datebaza.CuAsistentAI(a.contextLucru()), token)
	if err != nil {
		return RaspunsInstrumenteAI{Eroare: eroareInterfata(err)}
	}
	return RaspunsInstrumenteAI{Anulare: &r}
}
func (a *App) ListaSabloaneAI() RaspunsInstrumenteAI {
	r, err := datebaza.Nou(a.db).ListaSabloaneAI(a.contextLucru())
	if err != nil {
		return RaspunsInstrumenteAI{Eroare: eroareInterfata(err)}
	}
	return RaspunsInstrumenteAI{Sabloane: r}
}
func (a *App) SalveazaSablonAI(nume string, intentie ai.IntentiePlan) RaspunsInstrumenteAI {
	err := datebaza.Nou(a.db).SalveazaSablonAI(a.contextLucru(), nume, intentie)
	if err != nil {
		return RaspunsInstrumenteAI{Eroare: eroareInterfata(err)}
	}
	return a.ListaSabloaneAI()
}
func (a *App) StergeSablonAI(id string) RaspunsInstrumenteAI {
	err := datebaza.Nou(a.db).StergeSablonAI(a.contextLucru(), id)
	if err != nil {
		return RaspunsInstrumenteAI{Eroare: eroareInterfata(err)}
	}
	return a.ListaSabloaneAI()
}
func (a *App) VerificaPontajAI(intentie ai.IntentiePlan) RaspunsInstrumenteAI {
	r, err := datebaza.Nou(a.db).VerificaPontajAI(a.contextLucru(), intentie)
	if err != nil {
		return RaspunsInstrumenteAI{Eroare: eroareInterfata(err)}
	}
	return RaspunsInstrumenteAI{Verificare: &r}
}

func (a *App) ExportaPlanAICSV(intentie ai.IntentiePlan, token, operatieID, limba string) RaspunsExportJurnal {
	s := datebaza.Nou(a.db)
	var p datebaza.PlanAI
	var err error
	if operatieID != "" {
		p, err = s.PlanSalvatAI(a.contextLucru(), operatieID)
	} else {
		p, err = s.PregatestePlanAI(a.contextLucru(), intentie)
		if err == nil && (token == "" || token != p.Token) {
			return RaspunsExportJurnal{Eroare: &MesajEroare{Cod: "ai.plan_expirat"}}
		}
	}
	if err != nil {
		return RaspunsExportJurnal{Eroare: eroareInterfata(err)}
	}
	continut, err := csvPlanAI(p, i18n.Limba(limba))
	if err != nil {
		return RaspunsExportJurnal{Eroare: eroareInterfata(err)}
	}
	if a.ctx == nil {
		return RaspunsExportJurnal{Anulat: true}
	}
	folder, _ := folderRapoarte()
	cale, err := wruntime.SaveFileDialog(a.ctx, wruntime.SaveDialogOptions{Title: i18n.T(i18n.Limba(limba), "ai.export_titlu", nil), DefaultDirectory: folder, DefaultFilename: "plan-ai-" + p.Intentie.DeLa + ".csv", Filters: []wruntime.FileFilter{{DisplayName: "CSV (*.csv)", Pattern: "*.csv"}}, CanCreateDirectories: true})
	if err != nil {
		return RaspunsExportJurnal{Eroare: eroareInterfata(err)}
	}
	if cale == "" {
		return RaspunsExportJurnal{Anulat: true}
	}
	if !strings.EqualFold(filepath.Ext(cale), ".csv") {
		cale += ".csv"
	}
	if err = os.WriteFile(cale, continut, 0600); err != nil {
		return RaspunsExportJurnal{Eroare: &MesajEroare{Cod: "jurnal.export_esuat", Param: map[string]string{"detalii": err.Error()}}}
	}
	return RaspunsExportJurnal{Cale: cale, Numar: len(p.Randuri)}
}

const versiuneExportPlanAI = "1.0.0"

func csvPlanAI(p datebaza.PlanAI, limba i18n.Limba) ([]byte, error) {
	var b bytes.Buffer
	b.WriteString(bomUTF8)
	w := csv.NewWriter(&b)
	w.Comma = ';'
	tr := func(cheie string) string { return i18n.T(limba, "ai."+cheie, nil) }
	scrie := func(campuri ...string) error {
		for n, c := range campuri {
			t := strings.TrimLeft(c, " \t\r\n")
			if strings.HasPrefix(t, "=") || strings.HasPrefix(t, "+") || strings.HasPrefix(t, "-") || strings.HasPrefix(t, "@") || strings.HasPrefix(c, "\t") || strings.HasPrefix(c, "\r") {
				campuri[n] = "'" + c
			}
		}
		return w.Write(campuri)
	}
	stare := "export_propunere"
	if p.OperatieID != "" {
		stare = "export_aplicat"
	}
	if p.Anulat {
		stare = "export_anulat"
	}
	if err := scrie(tr("export_format"), versiuneExportPlanAI); err != nil {
		return nil, err
	}
	if err := scrie(tr("export_stare"), tr(stare), p.OperatieID); err != nil {
		return nil, err
	}
	if err := scrie(tr("plan_de_la"), p.Intentie.DeLa, tr("plan_pana_la"), p.Intentie.PanaLa, tr("plan_program"), p.Intentie.CodSchimb); err != nil {
		return nil, err
	}
	if err := scrie(tr("plan_persoane"), tr("plan_grup_"+p.Intentie.Selectie.Tip), p.Intentie.Selectie.Valoare, tr("plan_zile_"+p.Intentie.FiltruZile), tr("plan_context_"+p.Intentie.ModExistente)); err != nil {
		return nil, err
	}
	if p.Intentie.Selectie.FaraPontaj {
		if err := scrie(tr("plan_context_fara_pontaj")); err != nil {
			return nil, err
		}
	}
	for _, s := range append(append([]datebaza.SumarPlanAI{}, p.SumarLuni...), p.SumarPersoane...) {
		if err := scrie(tr("sumar"), s.Cheie, s.Nume, tr("selectati"), strconv.Itoa(s.Selectati), tr("afectati"), strconv.Itoa(s.Afectati), tr("plan_stare_noua"), strconv.Itoa(s.Noi), tr("plan_stare_inlocuire"), strconv.Itoa(s.Inlocuiri), tr("plan_stare_ocupata"), strconv.Itoa(s.Pastrate), tr("plan_stare_blocata"), strconv.Itoa(s.Blocate), tr("plan_stare_identica"), strconv.Itoa(s.Identice)); err != nil {
			return nil, err
		}
	}
	for _, s := range p.SelectieExplicata {
		stare := "selectie_exclus"
		if s.Inclus {
			stare = "selectie_inclus"
		}
		if err := scrie(tr("selectie_explicata"), s.Marca, s.Nume, tr(stare), tr("selectie_motiv_"+s.Motiv)); err != nil {
			return nil, err
		}
	}
	if err := scrie(tr("plan_persoane"), tr("marca"), tr("plan_data"), tr("plan_inainte"), tr("plan_dupa"), tr("plan_valoare_finala"), tr("plan_stare"), tr("motiv")); err != nil {
		return nil, err
	}
	for _, r := range p.Randuri {
		final := ""
		if p.OperatieID != "" {
			final = r.Inainte
			if r.Stare == "noua" || r.Stare == "inlocuire" {
				final = r.Dupa
			}
		}
		motiv := ""
		if r.Motiv != "" {
			motiv = tr("plan_motiv_" + r.Motiv)
		}
		if err := scrie(r.Nume, r.Marca, r.Data, r.Inainte, r.Dupa, final, tr("plan_stare_"+r.Stare), motiv); err != nil {
			return nil, err
		}
	}
	w.Flush()
	return b.Bytes(), w.Error()
}
