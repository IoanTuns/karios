package rapoarte

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"math"
	"regexp"
	"strconv"
	"strings"

	"github.com/pdfcpu/pdfcpu/pkg/api"
	"github.com/phpdave11/gofpdf"
	"golang.org/x/image/font/gofont/gobold"
	"golang.org/x/image/font/gofont/gobolditalic"
	"golang.org/x/image/font/gofont/goitalic"
	"golang.org/x/image/font/gofont/goregular"

	"pontix/internal/i18n"
)

// Rapoartele cu aspect PDFME (D087): sablonul este editat cu pdfme in interfata,
// care genereaza si PDF-ul din datele pregatite aici. Go pastreaza SQL-ul,
// calculele si anexa cu explicatiile si istoricul tehnic (D009).
const (
	ModPdfme = "PDFME"
	// TabelPdfme este numele tabelului din sablon care primeste randurile SQL.
	TabelPdfme        = "tabel"
	maxSablonPdfme    = 4 << 20
	maxPDFPdfme       = 64 << 20
	maxPaginiSablon   = 50
	maxElementePagina = 300
)

// Tipurile de elemente pdfme acceptate. SVG-ul si campurile de formular
// (semnatura, bife, selectii) nu au rost intr-un raport tiparit.
var tipuriPdfme = map[string]bool{"text": true, "multiVariableText": true, "table": true, "image": true, "line": true, "rectangle": true, "ellipse": true}

type elementPdfme struct {
	Name     string   `json:"name"`
	Type     string   `json:"type"`
	Head     []string `json:"head"`
	ReadOnly bool     `json:"readOnly"`
}
type paginaPdfme struct {
	Width   float64   `json:"width"`
	Height  float64   `json:"height"`
	Padding []float64 `json:"padding"`
}
type sablonPdfme struct {
	BasePdf json.RawMessage  `json:"basePdf"`
	Schemas [][]elementPdfme `json:"schemas"`
}

// EstePdfme spune daca documentul se genereaza cu pdfme in interfata.
func EstePdfme(d Document) bool {
	return d.Flexibil != nil && d.Flexibil.ModAspect == ModPdfme
}

func citesteSablonPdfme(text string) (sablonPdfme, error) {
	var s sablonPdfme
	if len(text) == 0 || len(text) > maxSablonPdfme || json.Unmarshal([]byte(text), &s) != nil {
		return s, eroare("rapoarte.pdfme_invalid")
	}
	var p paginaPdfme
	if json.Unmarshal(s.BasePdf, &p) != nil || len(p.Padding) != 4 {
		return s, eroare("rapoarte.pdfme_invalid")
	}
	for _, v := range append([]float64{p.Width, p.Height}, p.Padding...) {
		if math.IsNaN(v) || math.IsInf(v, 0) || v < 0 || v > 2000 {
			return s, eroare("rapoarte.pdfme_invalid")
		}
	}
	if p.Width < 50 || p.Height < 50 || len(s.Schemas) == 0 || len(s.Schemas) > maxPaginiSablon {
		return s, eroare("rapoarte.pdfme_invalid")
	}
	tabele := 0
	for _, pagina := range s.Schemas {
		if len(pagina) > maxElementePagina {
			return s, eroare("rapoarte.pdfme_invalid")
		}
		for _, e := range pagina {
			if !tipuriPdfme[e.Type] || e.Name == "" || len(e.Name) > 64 {
				return s, eroare("rapoarte.pdfme_invalid")
			}
			if e.Type == "table" && e.Name == TabelPdfme && !e.ReadOnly {
				tabele++
			}
		}
	}
	if tabele > 1 {
		return s, eroare("rapoarte.pdfme_invalid")
	}
	return s, nil
}

// DatePdfme construieste intrarile pdfme: variabilele raportului ca text si
// randurile vizibile ca JSON pentru tabelul TabelPdfme.
func DatePdfme(d Document) (map[string]string, error) {
	if !EstePdfme(d) {
		return nil, eroare("rapoarte.pdfme_invalid")
	}
	text, err := SablonPentruDocument(d)
	if err != nil {
		return nil, err
	}
	s, err := citesteSablonPdfme(text)
	if err != nil {
		return nil, err
	}
	for _, pagina := range s.Schemas {
		for _, e := range pagina {
			if e.Type == "table" && e.Name == TabelPdfme && !e.ReadOnly && len(e.Head) != len(d.Coloane) {
				return nil, &i18n.Eroare{Cod: "rapoarte.pdfme_coloane", Param: map[string]string{"sablon": itoa(len(e.Head)), "date": itoa(len(d.Coloane))}}
			}
		}
	}
	vars := map[string]string{"titlu": d.Titlu, "generat_la": d.GeneratLa.Format("02.01.2006 15:04"), "data_curenta": d.GeneratLa.Format("02.01.2006"), "introducere": d.Introducere, "nota": d.Nota}
	for k, v := range d.VariabileRaport {
		vars[k] = v
	}
	inloc := func(s string) string {
		return variabila.ReplaceAllStringFunc(s, func(v string) string {
			return vars[strings.TrimSpace(strings.TrimSuffix(strings.TrimPrefix(v, "{{"), "}}"))]
		})
	}
	zone := []string{"stanga", "centru", "dreapta"}
	for i, z := range zone {
		if i < len(d.Antet) {
			vars["antet_"+z] = inloc(d.Antet[i])
		}
		if i < len(d.Subsol) {
			vars["subsol_"+z] = inloc(d.Subsol[i])
		}
	}
	vars["introducere"] = inloc(d.Introducere)
	randuri := make([][]string, 0, len(d.Randuri))
	for _, r := range d.Randuri {
		randuri = append(randuri, r.Valori)
	}
	tabel, err := json.Marshal(randuri)
	if err != nil {
		return nil, err
	}
	vars[TabelPdfme] = string(tabel)
	if st := stiluriPdfme(d); st != "" {
		vars[IntrareStiluriPdfme] = st
	}
	return vars, nil
}

// IntrareStiluriPdfme poarta stilurile pe rand si pe celula ale tabelului
// (F09-T22): randurile aldine (weekend, sarbatori), celulele editate in copia
// raportului (D085) si stilurile generale ale corpului si capului de tabel.
// Indicii sunt in tot corpul tabelului; celulele au cheia "rand:coloana".
const IntrareStiluriPdfme = "pontix_stiluri"

// stilPdfme este stilul trimis desenarii tabelului; pe langa StilText are fundalul.
type stilPdfme struct {
	Aldin    *bool   `json:"aldin,omitempty"`
	Cursiv   *bool   `json:"cursiv,omitempty"`
	FontPT   float64 `json:"font_pt,omitempty"`
	Aliniere string  `json:"aliniere,omitempty"`
	Fundal   string  `json:"fundal,omitempty"`
}

func dinStilText(s StilText) stilPdfme {
	return stilPdfme{Aldin: s.Aldin, Cursiv: s.Cursiv, FontPT: s.FontPT, Aliniere: s.Aliniere}
}

type stiluriTabelPdfme struct {
	Cap        *stilPdfme           `json:"cap,omitempty"`
	Corp       *stilPdfme           `json:"corp,omitempty"`
	CapColoane map[string]stilPdfme `json:"cap_coloane,omitempty"`
	Randuri    map[string]stilPdfme `json:"randuri,omitempty"`
	Celule     map[string]stilPdfme `json:"celule,omitempty"`
	// Ascunse: coloanele zilelor care nu exista in luna raportului (z29-z31).
	Ascunse []int `json:"ascunse,omitempty"`
	// Titluri: titlurile coloanelor in limba raportului; desenarea le foloseste
	// doar pentru tabelele marcate pontixTitluri (modelele predefinite, D090).
	Titluri []string `json:"titluri,omitempty"`
}

// Culorile motorului vechi (pdf.go), ca PDF-ul pdfme sa arate la fel.
var culoriPdfme = map[string]string{"WEEKEND": "#E7EAE7", "SARBATOARE": "#FFF0C7", "ABSENTA": "#DFEFE4", CuloareZiLibera: "#D2D8D3"}

const fundalGrupPdfme = "#E7EFEA"

var coloanaZi = regexp.MustCompile(`^z(\d{2})$`)

func stiluriPdfme(d Document) string {
	st := stiluriTabelPdfme{CapColoane: map[string]stilPdfme{}, Randuri: map[string]stilPdfme{}, Celule: map[string]stilPdfme{}}
	aldin := true
	for i, r := range d.Randuri {
		switch {
		case r.Grup:
			st.Randuri[strconv.Itoa(i)] = stilPdfme{Aldin: &aldin, Fundal: fundalGrupPdfme}
		case r.Aldin:
			st.Randuri[strconv.Itoa(i)] = stilPdfme{Aldin: &aldin}
		}
		for j, c := range r.Culori {
			if f := culoriPdfme[c]; f != "" && j < len(d.Coloane) {
				st.Celule[strconv.Itoa(i)+":"+strconv.Itoa(j)] = stilPdfme{Fundal: f}
			}
		}
	}
	indice := map[string]int{}
	areTitluri := false
	for j, c := range d.Coloane {
		indice[c.Cod] = j
		if f := culoriPdfme[c.Culoare]; f != "" {
			st.CapColoane[strconv.Itoa(j)] = stilPdfme{Fundal: f}
		}
		if m := coloanaZi.FindStringSubmatch(c.Cod); m != nil && d.zileLuna > 0 {
			if n, _ := strconv.Atoi(m[1]); n > d.zileLuna {
				st.Ascunse = append(st.Ascunse, j)
			}
		}
		areTitluri = areTitluri || c.Titlu != ""
	}
	if areTitluri {
		for _, c := range d.Coloane {
			st.Titluri = append(st.Titluri, c.Titlu)
		}
	}
	for k, v := range d.Stiluri {
		sp := dinStilText(v)
		switch {
		case k == "tabel":
			st.Corp = &sp
		case k == "antet_tabel":
			st.Cap = &sp
		case strings.HasPrefix(k, "celula:"):
			parti := strings.SplitN(k, ":", 3)
			if len(parti) != 3 {
				continue
			}
			i, err := strconv.Atoi(parti[1])
			j, ok := indice[parti[2]]
			if err != nil || !ok || i < 0 || i >= len(d.Randuri) {
				continue
			}
			cheie := parti[1] + ":" + strconv.Itoa(j)
			// Editarea pastreaza fundalul calculat al celulei (weekend, absenta).
			sp.Fundal = st.Celule[cheie].Fundal
			st.Celule[cheie] = sp
		}
	}
	if st.Cap == nil && st.Corp == nil && len(st.CapColoane) == 0 && len(st.Randuri) == 0 && len(st.Celule) == 0 && len(st.Ascunse) == 0 && len(st.Titluri) == 0 {
		return ""
	}
	b, err := json.Marshal(st)
	if err != nil {
		return ""
	}
	return string(b)
}

func itoa(n int) string {
	b, _ := json.Marshal(n)
	return string(b)
}

// FonturiPdfme sunt fonturile Go folosite si de rapoartele vechi, ca PDF-ul
// pdfme sa aiba aceleasi diacritice si acelasi aspect.
func FonturiPdfme() map[string][]byte {
	return map[string][]byte{"Go": goregular.TTF, "GoBold": gobold.TTF, "GoItalic": goitalic.TTF, "GoBoldItalic": gobolditalic.TTF}
}

// AnexaPDF deseneaza anexa finala (explicatii si istoric tehnic D009), care
// nu face parte din sablonul editabil.
func AnexaPDF(d Document) ([]byte, error) {
	if len(d.Legenda) == 0 && d.Metadate == "" {
		return nil, nil
	}
	format := d.Pagina.Format
	if format == "" {
		format = "A4"
	}
	f := gofpdf.New(map[bool]string{true: "L", false: "P"}[d.Pagina.Orientare == "PEISAJ"], "mm", format, "")
	f.AddUTF8FontFromBytes("Go", "", goregular.TTF)
	f.AddUTF8FontFromBytes("Go", "B", gobold.TTF)
	m := float64(d.Pagina.MargineMM)
	if m <= 0 {
		m = 10
	}
	f.SetMargins(m, m, m)
	f.SetAutoPageBreak(true, m)
	f.AddPage()
	w, _ := f.GetPageSize()
	if len(d.Legenda) > 0 {
		f.SetFont("Go", "B", 10)
		f.MultiCell(w-2*m, 5, i18n.T(d.Limba, "rapoarte.anexa_calcul", nil), "", "L", false)
		f.SetFont("Go", "", 8)
		for _, v := range d.Legenda {
			f.MultiCell(w-2*m, 4, v, "", "L", false)
			f.Ln(1)
		}
	}
	if d.Metadate != "" {
		f.SetFont("Go", "B", 10)
		f.MultiCell(w-2*m, 5, i18n.T(d.Limba, "rapoarte.anexa_istoric", nil), "", "L", false)
		f.SetFont("Go", "", 8)
		f.MultiCell(w-2*m, 4, d.Metadate, "", "L", false)
	}
	var out bytes.Buffer
	if err := f.Output(&out); err != nil {
		return nil, eroare("rapoarte.aspect_prea_mare")
	}
	return out.Bytes(), nil
}

// FinalizeazaPdfme verifica PDF-ul generat de pdfme si ii adauga anexa.
func FinalizeazaPdfme(ctx context.Context, d Document, pdf []byte) ([]byte, int, error) {
	if !EstePdfme(d) || len(pdf) == 0 || len(pdf) > maxPDFPdfme {
		return nil, 0, eroare("rapoarte.pdfme_pdf_invalid")
	}
	if api.Validate(ctx, bytes.NewReader(pdf), nil, nil) != nil {
		return nil, 0, eroare("rapoarte.pdfme_pdf_invalid")
	}
	anexa, err := AnexaPDF(d)
	if err != nil {
		return nil, 0, err
	}
	final := pdf
	if anexa != nil {
		var out bytes.Buffer
		if err := api.MergeRaw(ctx, []io.ReadSeeker{bytes.NewReader(pdf), bytes.NewReader(anexa)}, &out, false, nil); err != nil {
			return nil, 0, traduce(err)
		}
		final = out.Bytes()
	}
	pagini, err := NumaraPaginiPDF(ctx, final)
	if err != nil {
		return nil, 0, err
	}
	return final, pagini, nil
}
