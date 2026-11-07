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
	"pontix/internal/diagnoza"
	"pontix/internal/exportexcel"
	"pontix/internal/fonturi"
	"pontix/internal/i18n"
	"pontix/internal/portabil"
	"pontix/internal/rapoarte"
)

type RaspunsModeleRapoarte struct {
	// Conversie: modelele convertite la acest acces (F09-T23), pentru anuntul din interfata.
	Conversie    []rapoarte.RezultatConversie `json:"conversie,omitempty"`
	Configuratie *rapoarte.Configuratie       `json:"configuratie"`
	SalvatID     string                       `json:"salvat_id"`
	Eroare       *MesajEroare                 `json:"eroare"`
}
type RaspunsRaport struct {
	Document   *rapoarte.Document `json:"document"`
	PDF        string             `json:"pdf"`
	Token      string             `json:"token"`
	NumeFisier string             `json:"nume_fisier"`
	Pagini     int                `json:"pagini"`
	// GenerarePdfme cere interfetei sa genereze PDF-ul din Sablon si Intrari,
	// apoi sa il trimita la FinalizeazaRaportPdfme (D087).
	GenerarePdfme bool              `json:"generare_pdfme"`
	Sablon        string            `json:"sablon"`
	Intrari       map[string]string `json:"intrari"`
	Eroare        *MesajEroare      `json:"eroare"`
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

// asiguraConversiaRapoarte converteste o data pe sesiune modelele vechi in
// SQL + pdfme (F09-T23). Daca un model nu se poate converti (de exemplu inainte
// de autentificare), conversia se reia la urmatorul acces.
func (a *App) asiguraConversiaRapoarte() []rapoarte.RezultatConversie {
	a.conversieMu.Lock()
	defer a.conversieMu.Unlock()
	// In doar citire conversia nu poate scrie; o face instanta care tine datele.
	if a.rapoarteConvertite || a.db == nil || a.doarCitire {
		return nil
	}
	rez, err := rapoarte.Nou(a.db).ConvertesteModeleVechi(a.contextLucru())
	if err != nil {
		diagnoza.ObtineManager().InregistreazaEroare("rapoarte/conversie", err)
		return nil
	}
	complet := true
	for _, r := range rez {
		if r.Eroare != "" {
			complet = false
			diagnoza.ObtineManager().InregistreazaEroare("rapoarte/conversie", errors.New("model "+r.ID+": "+r.Eroare))
		}
	}
	a.rapoarteConvertite = complet
	return rez
}

func (a *App) ListaModeleRapoarte() RaspunsModeleRapoarte {
	conversie := a.asiguraConversiaRapoarte()
	r := a.listaModeleRapoarte()
	r.Conversie = conversie
	return r
}

func (a *App) listaModeleRapoarte() RaspunsModeleRapoarte {
	c, err := rapoarte.Nou(a.db).Lista(a.contextLucru())
	if err != nil {
		return RaspunsModeleRapoarte{Eroare: mesajRaport(err)}
	}
	return RaspunsModeleRapoarte{Configuratie: &c}
}
func (a *App) SalveazaModelRaport(d rapoarte.Definitie) RaspunsModeleRapoarte {
	if rapoarte.ConversieNecesara(d) || !rapoarte.SursaSQL(d.Continut.Sursa) {
		return RaspunsModeleRapoarte{Eroare: &MesajEroare{Cod: "rapoarte.doar_editor_vizual"}}
	}
	id, err := rapoarte.Nou(a.db).Salveaza(a.contextLucru(), d)
	if err != nil {
		return RaspunsModeleRapoarte{Eroare: mesajRaport(err)}
	}
	r := a.listaModeleRapoarte()
	r.SalvatID = id
	return r
}
func (a *App) PregatesteRaport(c rapoarte.Cerere) RaspunsRaport {
	a.asiguraConversiaRapoarte()
	doc, err := rapoarte.Nou(a.db).Pregateste(a.contextLucru(), c, time.Now().UTC())
	if err != nil {
		return RaspunsRaport{Eroare: mesajRaport(err)}
	}
	return a.raspunsDocumentRaport(doc, c)
}

// ColoaneModelRaport ruleaza SQL-ul unui model si intoarce doar coloanele si schema,
// fara PDF (F09-T27). Editorul le foloseste ca sa construiasca sau sa actualizeze
// tabelul sablonului; inainte rula o previzualizare in modul vechi „tabel”, al carui
// motor a fost sters in F09-T23, si primea mereu rapoarte.model_neconvertit.
func (a *App) ColoaneModelRaport(d rapoarte.Definitie, c rapoarte.Cerere) RaspunsRaport {
	doc, err := rapoarte.Nou(a.db).PrevizualizeazaFlexibil(a.contextLucru(), d, c, time.Now().UTC())
	if err != nil {
		return RaspunsRaport{Eroare: mesajRaport(err)}
	}
	return RaspunsRaport{Document: &rapoarte.Document{Coloane: doc.Coloane, Schema: doc.Schema}}
}
func (a *App) PrevizualizeazaModelRaport(d rapoarte.Definitie, c rapoarte.Cerere) RaspunsRaport {
	doc, err := rapoarte.Nou(a.db).PrevizualizeazaFlexibil(a.contextLucru(), d, c, time.Now().UTC())
	if err != nil {
		return RaspunsRaport{Eroare: mesajRaport(err)}
	}
	return a.raspunsDocumentRaport(doc, c)
}
func (a *App) raspunsDocumentRaport(doc rapoarte.Document, c rapoarte.Cerere, tokenDat ...string) RaspunsRaport {
	if !rapoarte.EstePdfme(doc) {
		// Motorul vechi a fost sters (F09-T23); modelele se convertesc la primul acces.
		return RaspunsRaport{Eroare: &MesajEroare{Cod: "rapoarte.model_neconvertit"}}
	}
	return a.raspunsPdfme(doc, c, tokenDat...)
}

// raspunsPdfme pune documentul in cache fara PDF; PDF-ul vine din interfata.
func (a *App) raspunsPdfme(doc rapoarte.Document, c rapoarte.Cerere, tokenDat ...string) RaspunsRaport {
	intrari, err := rapoarte.DatePdfme(doc)
	if err != nil {
		return RaspunsRaport{Eroare: mesajRaport(err)}
	}
	sablon, err := rapoarte.SablonPentruDocument(doc)
	if err != nil {
		return RaspunsRaport{Eroare: mesajRaport(err)}
	}
	r := a.raspunsRaportValidat(doc, c, nil, 0, tokenDat...)
	r.GenerarePdfme, r.Sablon, r.Intrari = true, sablon, intrari
	return r
}

// FinalizeazaRaportPdfme primeste PDF-ul generat de pdfme pentru un token din
// cache, il verifica, adauga anexa si il pastreaza pentru salvare si tiparire.
func (a *App) FinalizeazaRaportPdfme(token, pdfBase64 string) RaspunsRaport {
	r, ok := a.raportPentru(token)
	if !ok {
		return RaspunsRaport{Eroare: &MesajEroare{Cod: "rapoarte.expirat"}}
	}
	if r.PDF != nil || !rapoarte.EstePdfme(r.Document) || len(pdfBase64) > base64.StdEncoding.EncodedLen(64<<20) {
		return RaspunsRaport{Eroare: &MesajEroare{Cod: "rapoarte.pdfme_pdf_invalid"}}
	}
	octeti, err := base64.StdEncoding.DecodeString(pdfBase64)
	if err != nil {
		return RaspunsRaport{Eroare: &MesajEroare{Cod: "rapoarte.pdfme_pdf_invalid"}}
	}
	pdf, pagini, err := rapoarte.FinalizeazaPdfme(a.contextLucru(), r.Document, octeti)
	if err != nil {
		return RaspunsRaport{Eroare: mesajRaport(err)}
	}
	a.raportMu.Lock()
	if v, ok := a.raportCache[token]; ok && v.PDF == nil {
		v.PDF = pdf
		a.raportCache[token] = v
	}
	a.raportMu.Unlock()
	return RaspunsRaport{Document: &r.Document, PDF: base64.StdEncoding.EncodeToString(pdf), Token: token, NumeFisier: r.Nume, Pagini: pagini}
}

// fonturiSistem citeste fonturile Windows o singura data, la prima cerere (D088).
var fonturiSistem = fonturi.Nou()

type RaspunsFonturiSistem struct {
	Fonturi []fonturi.Font `json:"fonturi"`
}
type RaspunsFontSistem struct {
	Nume   string       `json:"nume"`
	Date   string       `json:"date"`
	Subset bool         `json:"subset"`
	Eroare *MesajEroare `json:"eroare"`
}

// ListaFonturiSistem da numele fonturilor Windows care pot fi incluse in PDF.
func (a *App) ListaFonturiSistem() RaspunsFonturiSistem {
	return RaspunsFonturiSistem{Fonturi: fonturiSistem.Lista()}
}

// FontSistem da datele unui font Windows, cautat dupa nume in lista aplicatiei.
func (a *App) FontSistem(nume string) RaspunsFontSistem {
	date, f, err := fonturiSistem.Date(nume)
	if err != nil {
		return RaspunsFontSistem{Nume: nume, Eroare: mesajRaport(err)}
	}
	return RaspunsFontSistem{Nume: f.Nume, Date: base64.StdEncoding.EncodeToString(date), Subset: f.Subset}
}

// FonturiRapoarte da interfetei fonturile Go pentru pdfme, in base64.
func (a *App) FonturiRapoarte() map[string]string {
	rez := map[string]string{}
	for k, v := range rapoarte.FonturiPdfme() {
		rez[k] = base64.StdEncoding.EncodeToString(v)
	}
	return rez
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
	if rapoarte.EstePdfme(doc) {
		if _, e := rapoarte.DatePdfme(doc); e != nil {
			return RaspunsRaport{Eroare: mesajRaport(e)}
		}
		if e := rapoarte.Nou(a.db).JurnalPersonalizare(ctx, tokenNou, r.Document, doc, editari, acum, r.Cerere); e != nil {
			return RaspunsRaport{Eroare: mesajRaport(e)}
		}
		return a.raspunsPdfme(doc, r.Cerere, tokenNou)
	}
	return RaspunsRaport{Eroare: &MesajEroare{Cod: "rapoarte.model_neconvertit"}}
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
	if r.PDF == nil {
		return RaspunsFisierRaport{Eroare: &MesajEroare{Cod: "rapoarte.pdfme_negenerat"}}
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
	if r.PDF == nil {
		return RaspunsFisierRaport{Eroare: &MesajEroare{Cod: "rapoarte.pdfme_negenerat"}}
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
