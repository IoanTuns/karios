package rapoarte

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"regexp"
	"strings"
	"unicode/utf8"

	"pontix/internal/calcul"
	"pontix/internal/i18n"
)

// Conversia modelelor vechi in SQL + sablon pdfme (D087 punctul 7, F09-T23).
// Modelul isi pastreaza identificatorul; Salveaza creeaza o revizie noua si
// pune in jurnal definitia veche intreaga, deci originalul ramane in istoric.

var tabelPentruSursa = func() map[string]string {
	m := map[string]string{}
	for _, t := range tabeleStandard {
		m[t.Sursa] = t.Tabel
	}
	return m
}()

// ConversieNecesara: modelele care nu folosesc inca editorul vizual.
func ConversieNecesara(d Definitie) bool {
	return d.Continut.Flexibil == nil || d.Continut.Flexibil.ModAspect != ModPdfme
}

// coloaneModelStandard expandeaza coloanele unui model standard exact ca
// pregatesteStandard: zilele z01-z31, totalurile din calcul.ColoanePontaj,
// L Inceput dupa marca in pontajul detaliat.
func coloaneModelStandard(d Definitie, limba i18n.Limba) []ColoanaSablon {
	t := func(k string) string { return i18n.T(limba, "rapoarte."+k, nil) }
	rez := []ColoanaSablon{}
	for _, cod := range d.Continut.Coloane {
		switch cod {
		case "totaluri":
			for _, col := range calcul.ColoanePontaj() {
				if col == "ore_lucrate" && d.Continut.Sursa == "PONTAJ_DETALIAT" {
					col = "ore_lucratoare"
				}
				if col == "ore_norma" {
					continue
				}
				rez = append(rez, ColoanaSablon{Cod: "total_" + col, Titlu: i18n.T(limba, "pontaj.total_"+col, nil), Pondere: 6})
			}
		case "zile":
			for zi := 1; zi <= 31; zi++ {
				rez = append(rez, ColoanaSablon{Cod: fmt.Sprintf("z%02d", zi), Titlu: itoa(zi), Pondere: 3})
			}
		default:
			rez = append(rez, ColoanaSablon{Cod: cod, Titlu: t("coloana_" + cod), Pondere: float64(ponderiColoane[cod])})
			if cod == "marca" && d.Continut.Sursa == "PONTAJ_DETALIAT" {
				rez = append(rez, ColoanaSablon{Cod: "total_L_INCEPUT", Titlu: i18n.T(limba, "pontaj.total_L_INCEPUT", nil), Pondere: 6})
			}
		}
	}
	return rez
}

func optiuniConversie() OptiuniSablon {
	nume := map[string]string{}
	for _, k := range []string{"antet_stanga", "antet_centru", "antet_dreapta", "titlu", "subsol_stanga", "subsol_centru", "pagina"} {
		nume[k] = i18n.T(i18n.Romana, "rapoarte.pdfme_el_"+k, nil)
	}
	return OptiuniSablon{TitluriDinRaport: true, StilClasic: true, Nume: nume, TextPagina: i18n.T(i18n.Romana, "rapoarte.pdfme_pagina", nil)}
}

func flexibilGol() *RaportFlexibil {
	return &RaportFlexibil{InaltimeRand: 18, Ascunse: []string{}, Elemente: []ElementRaport{}, Calculate: []ColoanaCalculata{}, ConcediiTipuri: []string{}, ConcediiCategorii: []string{}, Weekend: true, Sarbatori: true}
}

// ConvertesteDefinitie intoarce modelul in forma SQL + pdfme si notele despre
// ce s-a aproximat. Modelele deja in editorul vizual raman neschimbate.
func ConvertesteDefinitie(d Definitie) (Definitie, []string, error) {
	if !ConversieNecesara(d) {
		return d, nil, nil
	}
	note := []string{}
	opt := optiuniConversie()
	if tabel, ok := tabelPentruSursa[d.Continut.Sursa]; ok {
		cols := coloaneModelStandard(d, i18n.Romana)
		coduri := []string{"rand"}
		for _, c := range cols {
			coduri = append(coduri, c.Cod)
		}
		f := flexibilGol()
		f.Ascunse = []string{"rand"}
		f.SQL = "SELECT " + strings.Join(coduri, ", ") + " FROM " + tabel + " ORDER BY rand"
		f.ModAspect = ModPdfme
		f.Sablon = SablonDinColoane(cols, d.Pagina, opt)
		d.Continut.Sursa, d.Continut.Coloane, d.Continut.Flexibil = "SQL_LIBER", []string{"rezultat"}, f
	} else if d.Continut.Flexibil != nil {
		f := *d.Continut.Flexibil
		switch f.ModAspect {
		case "", "TABEL":
			// Tabelul urmeaza singur coloanele SQL la generare (pontixTitluri).
			f.Sablon = SablonDinColoane(nil, d.Pagina, opt)
		case "MARKDOWN":
			f.Sablon, note = sablonDinMarkdown(f.Markdown, d.Pagina, opt)
			f.Markdown = ""
		case "LIBER":
			f.Sablon, note = sablonDinElemente(f.Elemente, f.InaltimeRand, d.Pagina, opt)
			f.Elemente = []ElementRaport{}
		default:
			return d, nil, eroare("rapoarte.flex_invalid")
		}
		f.ModAspect = ModPdfme
		d.Continut.Flexibil = &f
	} else {
		return d, nil, eroare("rapoarte.flex_invalid")
	}
	if sablon, n := adaugaImaginiSablon(d.Continut.Flexibil.Sablon, d); n != "" {
		note = append(note, n)
	} else {
		d.Continut.Flexibil.Sablon = sablon
	}
	return d, note, nil
}

var variabilaVeche = regexp.MustCompile(`\{\{\s*([^{}]+?)\s*\}\}`)

// variabilePdfme transforma {{nume}} in {nume}, sintaxa pdfme.
func variabilePdfme(s string) string {
	return variabilaVeche.ReplaceAllString(s, "{$1}")
}

// inaltimeText estimeaza inaltimea (mm) unui text la latimea si marimea date.
func inaltimeText(text string, latime, font float64) float64 {
	litera := font * 0.3528 * 0.52
	pe := math.Max(1, math.Floor(latime/litera))
	linii := 0.0
	for _, l := range strings.Split(text, "\n") {
		linii += math.Max(1, math.Ceil(float64(utf8.RuneCountInString(l))/pe))
	}
	return linii*font*0.3528*1.25 + 1.5
}

func fontVarianteGo() obiect {
	return obiect{"bold": "GoBold", "italic": "GoItalic", "boldItalic": "GoBoldItalic"}
}

// bazaFaraCorp: sablonul de pornire fara titlu si tabel (antet, subsol, pagina).
func bazaFaraCorp(p Pagina, opt OptiuniSablon) (obiect, []any) {
	var t obiect
	_ = json.Unmarshal([]byte(SablonDinColoane(nil, p, opt)), &t)
	pagina := t["schemas"].([]any)[0].([]any)
	return t, pagina[:3]
}

func tabelAuto(p Pagina, opt OptiuniSablon, x, y, latime float64) obiect {
	_, _, _, font := dimensiuniPagina(p)
	return tabelSablon(nil, x, y, latime, font, opt)
}

// sablonDinMarkdown: titlurile, paragrafele si listele devin texte, unul sub
// altul; {{tabel}} devine tabelul de date. Aldinul **...** ramane prin
// formatul Markdown in rand al pdfme.
func sablonDinMarkdown(md string, p Pagina, opt OptiuniSablon) (string, []string) {
	t, elemente := bazaFaraCorp(p, opt)
	w, _, m, font := dimensiuniPagina(p)
	util := w - 2*m
	y := m + 10
	nr := 0
	bloc := []string{}
	note := []string{}
	adaugaText := func(text string, marime float64, aldin bool) {
		if strings.TrimSpace(text) == "" {
			return
		}
		nr++
		h := inaltimeText(text, util, marime)
		el := textFixSablon(fmt.Sprintf("Text %d", nr), variabilePdfme(text), m, y, util, "left", marime, obiect{"height": h, "textFormat": "inline-markdown", "fontVariants": fontVarianteGo()})
		if aldin {
			el["fontName"] = "GoBold"
		}
		elemente = append(elemente, el)
		y += h + 2
	}
	golesteBloc := func() {
		adaugaText(strings.Join(bloc, "\n"), font, false)
		bloc = bloc[:0]
	}
	areTabel := false
	for _, linie := range strings.Split(strings.ReplaceAll(md, "\r\n", "\n"), "\n") {
		l := strings.TrimSpace(linie)
		switch {
		case l == "":
			golesteBloc()
		case l == "{{tabel}}" || l == "{{ tabel }}":
			golesteBloc()
			if !areTabel {
				elemente = append(elemente, tabelAuto(p, opt, m, y, util))
				y += 24
				areTabel = true
			}
		case strings.HasPrefix(l, "#"):
			golesteBloc()
			nivel := len(l) - len(strings.TrimLeft(l, "#"))
			adaugaText(strings.TrimSpace(l[nivel:]), font+math.Max(1, 6-float64(nivel)), true)
		case strings.HasPrefix(l, "- "):
			bloc = append(bloc, "• "+strings.TrimSpace(l[2:]))
		default:
			bloc = append(bloc, l)
		}
	}
	golesteBloc()
	if strings.Contains(md, "|") {
		note = append(note, "Tabelele Markdown scrise cu | au fost pastrate ca text.")
	}
	t["schemas"] = []any{elemente}
	b, _ := json.Marshal(t)
	return string(b), note
}

// sablonDinElemente: elementele pozitionate devin texte la aceleasi
// coordonate (fata de marginea utila, sub antet); {{tabel}} devine tabelul
// de date, iar elementele repetate pe fiecare rand devin un tabel, pentru ca
// pdfme nu are elemente repetate.
func sablonDinElemente(el []ElementRaport, pas float64, p Pagina, opt OptiuniSablon) (string, []string) {
	t, elemente := bazaFaraCorp(p, opt)
	w, _, m, _ := dimensiuniPagina(p)
	util := w - 2*m
	sus := m + 10
	note := []string{}
	repetate := []ElementRaport{}
	areTabel := false
	for i, e := range el {
		text := strings.TrimSpace(e.Text)
		x, y := m+e.X, sus+e.Y
		if e.RepetaRand {
			repetate = append(repetate, e)
			continue
		}
		if text == "{{tabel}}" {
			elemente = append(elemente, tabelAuto(p, opt, x, y, math.Min(e.Latime, util)))
			areTabel = true
			continue
		}
		aliniere := "left"
		font := "Go"
		if e.Stil != nil {
			aliniere = map[string]string{"C": "center", "R": "right"}[e.Stil.Aliniere]
			if aliniere == "" {
				aliniere = "left"
			}
			aldin := e.Stil.Aldin != nil && *e.Stil.Aldin
			cursiv := e.Stil.Cursiv != nil && *e.Stil.Cursiv
			font = []string{"Go", "GoBold", "GoItalic", "GoBoldItalic"}[map[bool]int{true: 1}[aldin]+map[bool]int{true: 2}[cursiv]]
		}
		marime := e.Font
		if marime <= 0 {
			marime = 10
		}
		elemente = append(elemente, textFixSablon(fmt.Sprintf("Element %d", i+1), variabilePdfme(e.Text), x, y, e.Latime, aliniere, marime, obiect{"height": e.Inaltime, "fontName": font}))
	}
	if len(repetate) > 0 {
		y := math.Inf(1)
		for _, e := range repetate {
			y = math.Min(y, sus+e.Y)
		}
		if !areTabel {
			elemente = append(elemente, tabelAuto(p, opt, m, y, util))
		}
		note = append(note, fmt.Sprintf("%d elemente repetate pe fiecare rand au fost inlocuite cu tabelul de date.", len(repetate)))
	}
	t["schemas"] = []any{elemente}
	b, _ := json.Marshal(t)
	return string(b), note
}

// adaugaImaginiSablon pune imaginile din antet (prima pagina, deasupra
// textelor) si din subsol (pe fiecare pagina, deasupra textelor de subsol).
// Daca sablonul ar depasi limita, imaginile raman doar in jurnal.
func adaugaImaginiSablon(sablon string, d Definitie) (string, string) {
	if len(d.Antet.Imagini)+len(d.Subsol.Imagini) == 0 {
		return sablon, ""
	}
	var t obiect
	if err := json.Unmarshal([]byte(sablon), &t); err != nil {
		return sablon, ""
	}
	w, h, m, _ := dimensiuniPagina(d.Pagina)
	treime := (w - 2*m) / 3
	imagine := func(im Imagine, zona string, y float64) (obiect, float64) {
		_, cfg, format, err := dateImagine(im)
		if err != nil {
			return nil, 0
		}
		ih := float64(im.InaltimeMM)
		iw := math.Min(treime, ih*float64(cfg.Width)/float64(cfg.Height))
		x := m
		switch zona {
		case "centru":
			x = m + treime + (treime-iw)/2
		case "dreapta":
			x = m + 3*treime - iw
		}
		mime := map[string]string{"png": "image/png", "jpeg": "image/jpeg"}[format]
		return cuProprietati(bazaPdfme("image"), obiect{"name": "Imagine " + zona, "content": "data:" + mime + ";base64," + im.Date, "position": obiect{"x": x, "y": y}, "width": iw, "height": ih, "readOnly": true}), ih
	}
	pagini := t["schemas"].([]any)
	prima := pagini[0].([]any)
	maxSus := 0.0
	imaginiSus := []any{}
	for _, zona := range []string{"stanga", "centru", "dreapta"} {
		if im, ok := d.Antet.Imagini[zona]; ok {
			if el, ih := imagine(im, zona, m); el != nil {
				imaginiSus = append(imaginiSus, el)
				maxSus = math.Max(maxSus, ih)
			}
		}
	}
	if maxSus > 0 {
		// Restul primei pagini coboara sub imagini.
		for _, e := range prima {
			if o, ok := e.(obiect); ok {
				if poz, ok := o["position"].(obiect); ok {
					poz["y"] = poz["y"].(float64) + maxSus + 2
				}
			}
		}
		pagini[0] = append(imaginiSus, prima...)
	}
	baza := t["basePdf"].(obiect)
	static := baza["staticSchema"].([]any)
	maxJos := 0.0
	for _, zona := range []string{"stanga", "centru", "dreapta"} {
		if im, ok := d.Subsol.Imagini[zona]; ok {
			ih := float64(im.InaltimeMM)
			if el, _ := imagine(im, zona, h-m-6-ih-1); el != nil {
				static = append(static, el)
				maxJos = math.Max(maxJos, ih)
			}
		}
	}
	baza["staticSchema"] = static
	if maxJos > 0 {
		if pad, ok := baza["padding"].([]any); ok && len(pad) == 4 {
			pad[2] = pad[2].(float64) + maxJos + 1
		}
	}
	b, _ := json.Marshal(t)
	if len(b) > maxSablonPdfme {
		return sablon, "Imaginile din antet/subsol depasesc limita sablonului; au ramas doar in istoricul modelului."
	}
	return string(b), ""
}

// RezultatConversie descrie conversia unui model, pentru jurnal si interfata.
type RezultatConversie struct {
	ID       string    `json:"id"`
	Denumire i18n.Text `json:"denumire"`
	Note     []string  `json:"note,omitempty"`
	Eroare   string    `json:"eroare,omitempty"`
}

// ConvertesteModeleVechi converteste toate modelele care nu folosesc editorul
// vizual. Fiecare model se salveaza separat (revizie noua si jurnal); un model
// care nu se poate converti ramane neschimbat si este raportat.
func (s *Serviciu) ConvertesteModeleVechi(ctx context.Context) ([]RezultatConversie, error) {
	c, err := s.Lista(ctx)
	if err != nil {
		return nil, err
	}
	rez := []RezultatConversie{}
	for _, d := range c.Modele {
		if !ConversieNecesara(d) {
			continue
		}
		r := RezultatConversie{ID: d.ID, Denumire: d.Denumire}
		nou, note, err := ConvertesteDefinitie(d)
		if err == nil {
			_, err = s.Salveaza(ctx, nou)
		}
		if err != nil {
			r.Eroare = err.Error()
		}
		r.Note = note
		rez = append(rez, r)
	}
	return rez, nil
}

// SursaSQL: sursele modelelor noi; sursele standard se folosesc din SQL (D090).
func SursaSQL(s string) bool {
	return esteFlexibil(s)
}
