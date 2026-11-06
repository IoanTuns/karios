package rapoarte

import (
	_ "embed"
	"encoding/json"
	"regexp"
	"strings"
	"unicode"
	"unicode/utf8"
)

// Constructorul Go al sabloanelor pdfme (F09-T23): aceeasi logica precum
// sablonDinColoane din frontend/src/raportPdfme.ts, pe proprietatile implicite
// ale pdfme 6.2.8 salvate in pdfme_baza.json. Il folosesc conversia modelelor
// vechi si adaptarea tabelului la coloanele curente.

//go:embed pdfme_baza.json
var pdfmeBazaJSON []byte

type obiect = map[string]any

func bazaPdfme(tip string) obiect {
	var b map[string]json.RawMessage
	if err := json.Unmarshal(pdfmeBazaJSON, &b); err != nil {
		panic(err)
	}
	var o obiect
	if err := json.Unmarshal(b[tip], &o); err != nil {
		panic(err)
	}
	return o
}

// ColoanaSablon este o coloana a tabelului de date: codul SQL si titlul.
type ColoanaSablon struct {
	Cod, Titlu string
	Pondere    float64
}

// OptiuniSablon corespunde optiunilor din frontend.
type OptiuniSablon struct {
	TitluriDinRaport bool
	StilClasic       bool
	Nume             map[string]string // numele elementelor (antet_stanga, titlu, pagina ...)
	TextPagina       string
}

var ziSablon = regexp.MustCompile(`^z\d{2}$`)

func titluColoanaSablon(c ColoanaSablon) string {
	t := strings.TrimSpace(c.Titlu)
	if t == "" {
		t = strings.TrimSpace(c.Cod)
	}
	r, n := utf8.DecodeRuneInString(t)
	if n == 0 {
		return t
	}
	return string(unicode.ToUpper(r)) + t[n:]
}

func dimensiuniPagina(p Pagina) (w, h, m, font float64) {
	l, i := 210.0, 297.0
	if p.Format == "A3" {
		l, i = 297, 420
	}
	w, h = l, i
	if p.Orientare == "PEISAJ" {
		w, h = i, l
	}
	m = float64(p.MargineMM)
	if m == 0 {
		m = 10
	}
	if m < 5 {
		m = 5
	}
	font = float64(p.FontPT)
	if font == 0 {
		font = 9
	}
	if font < 6 {
		font = 6
	}
	return
}

func cuProprietati(baza obiect, extra obiect) obiect {
	for k, v := range extra {
		baza[k] = v
	}
	return baza
}

func cutiePdfme(v float64) obiect {
	return obiect{"top": v, "right": v, "bottom": v, "left": v}
}

// tabelSablon construieste tabelul de date pentru coloanele date.
func tabelSablon(coloane []ColoanaSablon, x, y, latime, font float64, opt OptiuniSablon) obiect {
	total := 0.0
	ponderi := make([]float64, len(coloane))
	for i, c := range coloane {
		ponderi[i] = c.Pondere
		if ponderi[i] <= 0 {
			ponderi[i] = 1
		}
		total += ponderi[i]
	}
	if total == 0 {
		total = 1
	}
	head, latimi, coduri, gol := []any{}, []any{}, []any{}, []any{}
	for i, c := range coloane {
		head = append(head, titluColoanaSablon(c))
		latimi = append(latimi, ponderi[i]*100/total)
		coduri = append(coduri, c.Cod)
		gol = append(gol, "")
	}
	continut, _ := json.Marshal([]any{gol})
	t := cuProprietati(bazaPdfme("table"), obiect{
		"name": TabelPdfme, "position": obiect{"x": x, "y": y}, "width": latime, "height": 20.0,
		"showHead": true, "repeatHead": true, "head": head, "headWidthPercentages": latimi,
		"content": string(continut), "pontixColoane": coduri,
	})
	if opt.TitluriDinRaport {
		t["pontixTitluri"] = true
	}
	spatiu := 2.0
	if len(coloane) > 12 {
		spatiu = 0.8
	}
	for cheie, fontName := range map[string]string{"headStyles": "GoBold", "bodyStyles": "Go"} {
		s := t[cheie].(obiect)
		s["fontName"], s["fontSize"], s["padding"] = fontName, font, cutiePdfme(spatiu)
	}
	centrate := obiect{}
	for i, c := range coloane {
		if ziSablon.MatchString(c.Cod) {
			centrate[itoa(i)] = "center"
		}
	}
	if len(centrate) > 0 {
		cs, _ := t["columnStyles"].(obiect)
		if cs == nil {
			cs = obiect{}
		}
		cs["alignment"] = centrate
		t["columnStyles"] = cs
	}
	if opt.StilClasic {
		cap := t["headStyles"].(obiect)
		cap["backgroundColor"], cap["fontColor"], cap["borderColor"], cap["borderWidth"] = "#E7EFEA", "#000000", "#8C9690", cutiePdfme(0.15)
		corp := t["bodyStyles"].(obiect)
		corp["alternateBackgroundColor"], corp["borderColor"], corp["borderWidth"] = "", "#8C9690", cutiePdfme(0.15)
		t["tableStyles"] = obiect{"borderColor": "#8C9690", "borderWidth": 0.15}
	}
	return t
}

func textFixSablon(nume, continut string, x, y, latime float64, aliniere string, font float64, extra obiect) obiect {
	t := cuProprietati(bazaPdfme("text"), obiect{
		"name": nume, "content": continut, "readOnly": true, "position": obiect{"x": x, "y": y},
		"width": latime, "height": 8.0, "alignment": aliniere, "fontSize": font, "fontName": "Go",
	})
	return cuProprietati(t, extra)
}

// SablonDinColoane intoarce sablonul de pornire (antet in trei zone, titlu,
// tabelul de date, subsol si numarul paginii), ca JSON.
func SablonDinColoane(coloane []ColoanaSablon, p Pagina, opt OptiuniSablon) string {
	w, h, m, font := dimensiuniPagina(p)
	util := w - 2*m
	treime := util / 3
	n := func(k, implicit string) string {
		if v := opt.Nume[k]; v != "" {
			return v
		}
		return implicit
	}
	antet := []any{
		textFixSablon(n("antet_stanga", "antet_stanga_txt"), "{antet_stanga}", m, m, treime, "left", font, nil),
		textFixSablon(n("antet_centru", "antet_centru_txt"), "{antet_centru}", m+treime, m, treime, "center", font, nil),
		textFixSablon(n("antet_dreapta", "antet_dreapta_txt"), "{antet_dreapta}", m+2*treime, m, treime, "right", font, nil),
		textFixSablon(n("titlu", "titlu_txt"), "{titlu}", m, m+10, util, "center", font, obiect{"height": 10.0, "fontSize": font + 5, "fontName": "GoBold"}),
	}
	jos := h - m - 6
	static := []any{
		textFixSablon(n("subsol_stanga", "subsol_stanga_txt"), "{subsol_stanga}", m, jos, treime, "left", font, obiect{"fontSize": 8.0}),
		textFixSablon(n("subsol_centru", "subsol_centru_txt"), "{subsol_centru}", m+treime, jos, treime, "center", font, obiect{"fontSize": 8.0}),
		textFixSablon(n("pagina", "pagina_txt"), opt.TextPagina, m+2*treime, jos, treime, "right", font, obiect{"fontSize": 8.0}),
	}
	// Antetul si titlul se repeta pe fiecare pagina (F09-T26); continutul incepe sub ele.
	sablon := obiect{
		"basePdf": obiect{"width": w, "height": h, "padding": []any{m + 22, m, m + 8, m}, "staticSchema": append(antet, static...)},
		"schemas": []any{[]any{tabelSablon(coloane, m, m+24, util, font, opt)}},
	}
	b, _ := json.Marshal(sablon)
	return string(b)
}

// SablonPentruDocument intoarce sablonul folosit la generare. Tabelul de date
// marcat pontixTitluri (modelele predefinite si cele convertite) se adapteaza
// la coloanele vizibile curente: coloanele pastrate isi pastreaza latimea,
// cele noi primesc latimea medie, iar totalul se normalizeaza la 100%.
// Celelalte sabloane raman neschimbate.
func SablonPentruDocument(d Document) (string, error) {
	if !EstePdfme(d) {
		return "", eroare("rapoarte.pdfme_invalid")
	}
	var t obiect
	if err := json.Unmarshal([]byte(d.Flexibil.Sablon), &t); err != nil {
		return "", eroare("rapoarte.pdfme_invalid")
	}
	pagini, _ := t["schemas"].([]any)
	for _, pg := range pagini {
		elemente, _ := pg.([]any)
		for _, el := range elemente {
			tabel, _ := el.(obiect)
			if tabel == nil || tabel["type"] != "table" || tabel["name"] != TabelPdfme || tabel["readOnly"] == true || tabel["pontixTitluri"] != true {
				continue
			}
			vechi := []string{}
			if c, ok := tabel["pontixColoane"].([]any); ok {
				for _, v := range c {
					s, _ := v.(string)
					vechi = append(vechi, s)
				}
			}
			noi := make([]string, len(d.Coloane))
			for i, c := range d.Coloane {
				noi[i] = c.Cod
			}
			if strings.Join(vechi, "\x00") == strings.Join(noi, "\x00") {
				return d.Flexibil.Sablon, nil
			}
			head, _ := tabel["head"].([]any)
			latimi := []float64{}
			if l, ok := tabel["headWidthPercentages"].([]any); ok {
				for _, v := range l {
					f, _ := v.(float64)
					latimi = append(latimi, f)
				}
			}
			medie := 100.0 / float64(max(1, len(noi)))
			if len(latimi) > 0 {
				suma := 0.0
				for _, v := range latimi {
					suma += v
				}
				medie = suma / float64(len(latimi))
			}
			nouHead, nouLatimi, gol, coduri := []any{}, []float64{}, []any{}, []any{}
			for _, c := range d.Coloane {
				i := -1
				for k, v := range vechi {
					if v == c.Cod {
						i = k
					}
				}
				if i >= 0 && i < len(head) && i < len(latimi) && latimi[i] > 0 {
					nouHead, nouLatimi = append(nouHead, head[i]), append(nouLatimi, latimi[i])
				} else {
					nouHead, nouLatimi = append(nouHead, titluColoanaSablon(ColoanaSablon{Cod: c.Cod, Titlu: c.Titlu})), append(nouLatimi, medie)
				}
				gol, coduri = append(gol, ""), append(coduri, c.Cod)
			}
			total := 0.0
			for _, v := range nouLatimi {
				total += v
			}
			procente := []any{}
			for _, v := range nouLatimi {
				procente = append(procente, v*100/total)
			}
			continut, _ := json.Marshal([]any{gol})
			tabel["head"], tabel["headWidthPercentages"], tabel["content"], tabel["pontixColoane"] = nouHead, procente, string(continut), coduri
			// Alinierea pe coloane se reface dupa pozitiile noi: zilele centrate.
			if cs, ok := tabel["columnStyles"].(obiect); ok {
				centrate := obiect{}
				for i, c := range d.Coloane {
					if ziSablon.MatchString(c.Cod) {
						centrate[itoa(i)] = "center"
					}
				}
				cs["alignment"] = centrate
			}
			b, err := json.Marshal(t)
			if err != nil {
				return "", err
			}
			return string(b), nil
		}
	}
	return d.Flexibil.Sablon, nil
}
