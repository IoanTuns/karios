package main

import (
	"encoding/base64"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"time"

	wruntime "github.com/wailsapp/wails/v2/pkg/runtime"

	"pontix/internal/datebaza"
	"pontix/internal/exportexcel"
	"pontix/internal/i18n"
	"pontix/internal/portabil"
	"pontix/internal/rapoarte"
)

type RaspunsModeleRapoarte struct {
	Configuratie *rapoarte.Configuratie `json:"configuratie"`
	SalvatID     string                 `json:"salvat_id"`
	Eroare       *MesajEroare           `json:"eroare"`
}
type RaspunsRaport struct {
	Document   *rapoarte.Document `json:"document"`
	PDF        string             `json:"pdf"`
	Token      string             `json:"token"`
	NumeFisier string             `json:"nume_fisier"`
	Pagini     int                `json:"pagini"`
	Eroare     *MesajEroare       `json:"eroare"`
}
type RaspunsFisierRaport struct {
	Anulat bool         `json:"anulat"`
	Cale   string       `json:"cale"`
	Eroare *MesajEroare `json:"eroare"`
}
type raportPregatit struct {
	Cerere   rapoarte.Cerere
	Document rapoarte.Document
	PDF      []byte
	Nume     string
	Limba    i18n.Limba
}

func mesajRaport(err error) *MesajEroare {
	if err == nil {
		return nil
	}
	mesaj := &MesajEroare{Cod: "rapoarte.eroare"}
	var e *i18n.Eroare
	if errors.As(err, &e) {
		mesaj.Cod, mesaj.Param = e.Cod, e.Param
	}
	return mesaj
}
func (a *App) ListaModeleRapoarte() RaspunsModeleRapoarte {
	c, err := rapoarte.Nou(a.db).Lista(a.contextLucru())
	if err != nil {
		return RaspunsModeleRapoarte{Eroare: mesajRaport(err)}
	}
	return RaspunsModeleRapoarte{Configuratie: &c}
}
func (a *App) SalveazaModelRaport(d rapoarte.Definitie) RaspunsModeleRapoarte {
	id, err := rapoarte.Nou(a.db).Salveaza(a.contextLucru(), d)
	if err != nil {
		return RaspunsModeleRapoarte{Eroare: mesajRaport(err)}
	}
	r := a.ListaModeleRapoarte()
	r.SalvatID = id
	return r
}
func (a *App) PregatesteRaport(c rapoarte.Cerere) RaspunsRaport {
	doc, err := rapoarte.Nou(a.db).Pregateste(a.contextLucru(), c, time.Now().UTC())
	if err != nil {
		return RaspunsRaport{Eroare: mesajRaport(err)}
	}
	return a.raspunsDocumentRaport(doc, c)
}
func (a *App) PrevizualizeazaModelRaport(d rapoarte.Definitie, c rapoarte.Cerere) RaspunsRaport {
	doc, err := rapoarte.Nou(a.db).PrevizualizeazaFlexibil(a.contextLucru(), d, c, time.Now().UTC())
	if err != nil {
		return RaspunsRaport{Eroare: mesajRaport(err)}
	}
	return a.raspunsDocumentRaport(doc, c)
}
func (a *App) raspunsDocumentRaport(doc rapoarte.Document, c rapoarte.Cerere) RaspunsRaport {
	pdf, err := rapoarte.GenereazaPDF(doc)
	if err != nil {
		return RaspunsRaport{Eroare: mesajRaport(err)}
	}
	pagini, err := rapoarte.NumaraPaginiPDF(a.contextLucru(), pdf)
	if err != nil {
		return RaspunsRaport{Eroare: mesajRaport(err)}
	}
	return a.raspunsRaportValidat(doc, c, pdf, pagini)
}
func (a *App) raspunsRaportValidat(doc rapoarte.Document, c rapoarte.Cerere, pdf []byte, pagini int, tokenDat ...string) RaspunsRaport {
	token := rapoarte.IDNou()
	if len(tokenDat) > 0 {
		token = tokenDat[0]
	}
	nume := doc.Sursa + "-" + doc.GeneratLa.Format("20060102-150405") + "-" + token[:8] + ".pdf"
	a.raportMu.Lock()
	if a.raportCache == nil {
		a.raportCache = map[string]raportPregatit{}
	}
	a.raportCache[token] = raportPregatit{Document: doc, PDF: pdf, Nume: nume, Limba: i18n.Limba(c.Limba), Cerere: c}
	a.raportOrdine = append(a.raportOrdine, token)
	if len(a.raportOrdine) > 4 {
		delete(a.raportCache, a.raportOrdine[0])
		a.raportOrdine = a.raportOrdine[1:]
	}
	a.raportMu.Unlock()
	return RaspunsRaport{Document: &doc, PDF: base64.StdEncoding.EncodeToString(pdf), Token: token, NumeFisier: nume, Pagini: pagini}
}
func (a *App) PersonalizeazaRaport(token string, editari []rapoarte.EditareTextRaport) RaspunsRaport {
	r, ok := a.raportPentru(token)
	if !ok {
		return RaspunsRaport{Eroare: &MesajEroare{Cod: "rapoarte.expirat"}}
	}
	ctx := a.contextLucru()
	acum := time.Now().UTC()
	doc, e := rapoarte.EditeazaDocument(r.Document, editari, datebaza.Actor(ctx), acum)
	if e != nil {
		return RaspunsRaport{Eroare: mesajRaport(e)}
	}
	tokenNou := rapoarte.IDNou()
	doc.Metadate += " | " + i18n.T(r.Limba, "rapoarte.copie_id", map[string]string{"id": tokenNou})
	pdf, e := rapoarte.GenereazaPDF(doc)
	if e != nil {
		return RaspunsRaport{Eroare: mesajRaport(e)}
	}
	pagini, e := rapoarte.NumaraPaginiPDF(ctx, pdf)
	if e != nil {
		return RaspunsRaport{Eroare: mesajRaport(e)}
	}
	if e := rapoarte.Nou(a.db).JurnalPersonalizare(ctx, tokenNou, r.Document, doc, editari, acum, r.Cerere); e != nil {
		return RaspunsRaport{Eroare: mesajRaport(e)}
	}
	return a.raspunsRaportValidat(doc, r.Cerere, pdf, pagini, tokenNou)
}
func (a *App) raportPentru(token string) (raportPregatit, bool) {
	a.raportMu.Lock()
	defer a.raportMu.Unlock()
	r, ok := a.raportCache[token]
	return r, ok
}
func folderRapoarte() (string, error) {
	exe, err := os.Executable()
	if err != nil {
		return "", err
	}
	folder := filepath.Join(filepath.Dir(exe), portabil.NumeFolderDate, "rapoarte")
	if err := os.MkdirAll(folder, 0700); err != nil {
		return "", err
	}
	return folder, nil
}
func (a *App) SalveazaPDFRaport(token string) RaspunsFisierRaport {
	r, ok := a.raportPentru(token)
	if !ok {
		return RaspunsFisierRaport{Eroare: &MesajEroare{Cod: "rapoarte.expirat"}}
	}
	folder, err := folderRapoarte()
	if err != nil {
		return RaspunsFisierRaport{Eroare: mesajRaport(err)}
	}
	cale, err := wruntime.SaveFileDialog(a.ctx, wruntime.SaveDialogOptions{Title: i18n.T(r.Limba, "rapoarte.salveaza_pdf", nil), DefaultDirectory: folder, DefaultFilename: r.Nume, Filters: []wruntime.FileFilter{{DisplayName: i18n.T(r.Limba, "rapoarte.filtru_pdf", nil), Pattern: "*.pdf"}}, CanCreateDirectories: true})
	if err != nil {
		return RaspunsFisierRaport{Eroare: mesajRaport(err)}
	}
	if cale == "" {
		return RaspunsFisierRaport{Anulat: true}
	}
	if !strings.EqualFold(filepath.Ext(cale), ".pdf") {
		cale += ".pdf"
	}
	if err := os.WriteFile(cale, r.PDF, 0600); err != nil {
		return RaspunsFisierRaport{Eroare: mesajRaport(err)}
	}
	return RaspunsFisierRaport{Cale: cale}
}
func (a *App) DeschidePDFPentruTiparire(token string) RaspunsFisierRaport {
	r, ok := a.raportPentru(token)
	if !ok {
		return RaspunsFisierRaport{Eroare: &MesajEroare{Cod: "rapoarte.expirat"}}
	}
	folder, err := folderRapoarte()
	if err != nil {
		return RaspunsFisierRaport{Eroare: mesajRaport(err)}
	}
	// Tokenul este generat de aplicație și verificat în cache; nu acceptăm căi din interfață.
	cale := filepath.Join(folder, token+"-"+rapoarte.IDNou()[:8]+".pdf")
	if err := os.WriteFile(cale, r.PDF, 0600); err != nil {
		return RaspunsFisierRaport{Eroare: mesajRaport(err)}
	}
	if err := deschidePDFAplicatie(cale); err != nil {
		return RaspunsFisierRaport{Eroare: &MesajEroare{Cod: "rapoarte.deschidere_esuata"}}
	}
	return RaspunsFisierRaport{Cale: cale}
}

func (a *App) SalveazaExcelRaport(token string) RaspunsFisierRaport {
	r, ok := a.raportPentru(token)
	if !ok {
		return RaspunsFisierRaport{Eroare: &MesajEroare{Cod: "rapoarte.expirat"}}
	}
	continut, err := rapoarte.GenereazaExcel(r.Document)
	if err != nil {
		return RaspunsFisierRaport{Eroare: &MesajEroare{Cod: "excel.salvare_esuata"}}
	}
	folder, err := folderRapoarte()
	if err != nil {
		return RaspunsFisierRaport{Eroare: &MesajEroare{Cod: "excel.salvare_esuata"}}
	}
	cale, err := wruntime.SaveFileDialog(a.ctx, wruntime.SaveDialogOptions{Title: i18n.T(r.Limba, "pontaj.export_excel", nil), DefaultDirectory: folder, DefaultFilename: strings.TrimSuffix(r.Nume, ".pdf") + ".xlsx", Filters: []wruntime.FileFilter{{DisplayName: i18n.T(r.Limba, "pontaj.export_excel", nil), Pattern: "*.xlsx"}}, CanCreateDirectories: true})
	if err != nil {
		return RaspunsFisierRaport{Eroare: &MesajEroare{Cod: "excel.salvare_esuata"}}
	}
	if cale == "" {
		return RaspunsFisierRaport{Anulat: true}
	}
	if !strings.EqualFold(filepath.Ext(cale), ".xlsx") {
		cale += ".xlsx"
	}
	if err := exportexcel.SalveazaVerificat(cale, continut); err != nil {
		return RaspunsFisierRaport{Eroare: &MesajEroare{Cod: "excel.salvare_esuata"}}
	}
	return RaspunsFisierRaport{Cale: cale}
}
