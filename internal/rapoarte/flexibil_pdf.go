package rapoarte

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"github.com/phpdave11/gofpdf"
	"golang.org/x/image/font/gofont/gobold"
	"golang.org/x/image/font/gofont/goregular"
	"pontix/internal/i18n"
	"strings"
)

func GenereazaPDFFlexibil(d Document) ([]byte, error) {
	f := gofpdf.New(map[bool]string{true: "L", false: "P"}[d.Pagina.Orientare == "PEISAJ"], "mm", d.Pagina.Format, "")
	f.AddUTF8FontFromBytes("Go", "", goregular.TTF)
	f.AddUTF8FontFromBytes("Go", "B", gobold.TTF)
	f.SetFont("Go", "", float64(d.Pagina.FontPT))
	m := float64(d.Pagina.MargineMM)
	f.SetMargins(m, m, m)
	f.SetAutoPageBreak(false, m)
	w, h := f.GetPageSize()
	anexa := false
	vars := map[string]string{"titlu": d.Titlu, "perioada": d.GeneratLa.Format("2006-01"), "generat_la": d.GeneratLa.Format("02.01.2006 15:04"), "tabel": ""}
	for k, v := range d.VariabileRaport {
		vars[k] = v
	}
	if len(d.Randuri) > 0 {
		for i, c := range d.Coloane {
			vars[c.Cod] = d.Randuri[0].Valori[i]
		}
	}
	inloc := func(s string) string {
		return variabila.ReplaceAllStringFunc(s, func(v string) string {
			return vars[strings.TrimSpace(strings.TrimSuffix(strings.TrimPrefix(v, "{{"), "}}"))]
		})
	}
	zone := (w - 2*m) / 3
	footerH := 6.
	maxText, maxImg := 0., 0.
	f.SetFont("Go", "", 8)
	for _, v := range d.Subsol {
		n := float64(len(f.SplitText(inloc(v), zone))) * 3.5
		if n > maxText {
			maxText = n
		}
	}
	for _, im := range d.SubsolImagini {
		if float64(im.InaltimeMM) > maxImg {
			maxImg = float64(im.InaltimeMM)
		}
	}
	footerH += maxText + maxImg
	f.SetFont("Go", "", float64(d.Pagina.FontPT))
	// Verificam coordonatele si inaltimea textului, fara decupare tacita.
	for _, e := range d.Flexibil.Elemente {
		if d.Flexibil.ModAspect != "LIBER" {
			break
		}
		if e.X+e.Latime > w-2*m || e.Y+e.Inaltime > h-2*m-footerH {
			return nil, eroare("rapoarte.aspect_prea_mare")
		}
		if strings.TrimSpace(e.Text) != "{{tabel}}" {
			f.SetFont("Go", "", e.Font)
			if float64(len(f.SplitText(inloc(e.Text), e.Latime)))*e.Font*.45 > e.Inaltime {
				return nil, eroare("rapoarte.aspect_prea_mare")
			}
		}
	}
	f.SetFont("Go", "", float64(d.Pagina.FontPT))
	imagine := func(ims map[string]Imagine, y float64) float64 {
		max := 0.
		for i, z := range []string{"stanga", "centru", "dreapta"} {
			im, ok := ims[z]
			if !ok {
				continue
			}
			data, e := base64.StdEncoding.DecodeString(im.Date)
			if e != nil {
				continue
			}
			tip := "JPG"
			if bytes.HasPrefix(data, []byte{137, 80, 78, 71}) {
				tip = "PNG"
			}
			n := fmt.Sprintf("%s-%d-%s", z, f.PageNo(), tip)
			o := gofpdf.ImageOptions{ImageType: tip}
			info := f.RegisterImageOptionsReader(n, o, bytes.NewReader(data))
			if info == nil {
				continue
			}
			ih := float64(im.InaltimeMM)
			iw := ih * info.Width() / info.Height()
			zone := (w - 2*m) / 3
			if iw > zone {
				ih *= zone / iw
				iw = zone
			}
			f.ImageOptions(n, m+float64(i)*zone+(zone-iw)*float64(i)/2, y, iw, ih, false, o, 0, "")
			if ih > max {
				max = ih
			}
		}
		return max
	}
	f.SetFooterFunc(func() {
		f.SetFont("Go", "", 8)
		y := h - m - footerH
		for i, v := range d.Subsol {
			f.SetXY(m+float64(i)*zone, y)
			f.MultiCell(zone, 3.5, inloc(v), "", []string{"L", "C", "R"}[i], false)
		}
		imagine(d.SubsolImagini, y+maxText)
		f.SetFont("Go", "", 7)
		f.SetXY(m, h-m-4)
		f.CellFormat(w-2*m, 4, fmt.Sprintf("%d", f.PageNo()), "", 0, "R", false, 0, "")
	})

	startTabel := m + 16
	nouaPagina := func() {
		f.AddPage()
		f.SetFont("Go", "B", float64(d.Pagina.FontPT)+2)
		f.SetXY(m, m)
		if anexa || d.Flexibil.ModAspect == "TABEL" || d.Flexibil.ModAspect == "" {
			f.MultiCell(w-2*m, 6, d.Titlu, "", "L", false)
		}
		f.SetFont("Go", "", float64(d.Pagina.FontPT))
		if d.Flexibil.ModAspect == "LIBER" && !anexa {
			imagine(d.AntetImagini, m)
			for _, e := range d.Flexibil.Elemente {
				if e.RepetaRand || strings.TrimSpace(e.Text) == "{{tabel}}" {
					continue
				}
				f.SetFont("Go", "", e.Font)
				f.SetXY(m+e.X, m+e.Y)
				f.MultiCell(e.Latime, e.Font*.45, inloc(e.Text), "", "L", false)
			}
			f.SetFont("Go", "", float64(d.Pagina.FontPT))
		} else {
			y := f.GetY() + imagine(d.AntetImagini, f.GetY())
			maxY := y
			for i, v := range d.Antet {
				if strings.TrimSpace(v) != "" {
					f.SetXY(m+float64(i)*zone, y)
					f.MultiCell(zone, 5, inloc(v), "", []string{"L", "C", "R"}[i], false)
					if f.GetY() > maxY {
						maxY = f.GetY()
					}
				}
			}
			y = maxY
			if d.Introducere != "" && !anexa {
				f.SetXY(m, y)
				f.MultiCell(w-2*m, 4, inloc(d.Introducere), "", "L", false)
				y = f.GetY()
			}
			startTabel = y + 3
		}
		if startTabel > h-m-footerH-10 {
			f.SetError(eroare("rapoarte.aspect_prea_mare"))
		}
		f.SetY(startTabel)
	}
	nouaPagina()
	tableX, tableY, tableW := m, startTabel, w-2*m
	tableActif := true
	tableLimite := h - m - footerH - 2
	if d.Flexibil.ModAspect == "LIBER" {
		tableActif = false
		for _, e := range d.Flexibil.Elemente {
			if strings.TrimSpace(e.Text) == "{{tabel}}" {
				tableActif = true
				tableX = m + e.X
				tableY = m + e.Y
				tableW = e.Latime
				tableLimite = m + e.Y + e.Inaltime
				break
			}
		}
	}
	var tabel func()
	tabel = func() {
		if !tableActif {
			return
		}
		if len(d.Coloane) == 0 {
			return
		}
		f.SetY(tableY)
		cw := tableW / float64(len(d.Coloane))
		if cw <= 2.5 {
			f.SetError(eroare("rapoarte.aspect_prea_mare"))
			return
		}
		line := float64(d.Pagina.FontPT) * .5
		labels := []string{}
		for _, c := range d.Coloane {
			labels = append(labels, c.Titlu)
		}
		var deseneaza func([]string, bool) error
		deseneaza = func(vals []string, bold bool) error {
			if bold {
				f.SetFont("Go", "B", float64(d.Pagina.FontPT))
			} else {
				f.SetFont("Go", "", float64(d.Pagina.FontPT))
			}
			rh := line + 2
			for _, v := range vals {
				n := float64(len(f.SplitText(v, cw-2)))*line + 2
				if n > rh {
					rh = n
				}
			}
			if rh > tableLimite-tableY-2*line-4 {
				return eroare("rapoarte.aspect_prea_mare")
			}
			if f.GetY()+rh > tableLimite {
				if bold {
					return eroare("rapoarte.aspect_prea_mare")
				}
				nouaPagina()
				if d.Flexibil.ModAspect != "LIBER" {
					tableY = startTabel
				}
				f.SetY(tableY)
				if e := deseneaza(labels, true); e != nil {
					return e
				}
				f.SetFont("Go", "", float64(d.Pagina.FontPT))
				if f.GetY()+rh > tableLimite {
					return eroare("rapoarte.aspect_prea_mare")
				}
			}
			y := f.GetY()
			for i, v := range vals {
				x := tableX + float64(i)*cw
				f.Rect(x, y, cw, rh, "D")
				f.SetXY(x+1, y+1)
				f.MultiCell(cw-2, line, v, "", "L", false)
			}
			f.SetY(y + rh)
			return nil
		}
		if e := deseneaza(labels, true); e != nil {
			f.SetError(e)
			return
		}
		for _, r := range d.Randuri {
			if e := deseneaza(r.Valori, false); e != nil {
				f.SetError(e)
				return
			}
		}
		if len(d.Randuri) == 0 {
			f.MultiCell(tableW, line, d.FaraDate, "", "L", false)
		}
	}
	if d.Flexibil.ModAspect == "LIBER" {
		repetate := []ElementRaport{}
		minY, maxY := h, 0.
		for _, e := range d.Flexibil.Elemente {
			if e.RepetaRand {
				repetate = append(repetate, e)
				if e.Y < minY {
					minY = e.Y
				}
				if e.Y+e.Inaltime > maxY {
					maxY = e.Y + e.Inaltime
				}
			}
		}
		if len(repetate) > 0 {
			if d.Flexibil.InaltimeRand < maxY-minY {
				return nil, eroare("rapoarte.aspect_prea_mare")
			}
			initial := map[string]string{}
			for k, v := range vars {
				initial[k] = v
			}
			offset := 0.
			for ri, r := range d.Randuri {
				if m+maxY+offset > h-m-footerH-2 {
					for k, v := range initial {
						vars[k] = v
					}
					nouaPagina()
					offset = 0
				}
				for i, c := range d.Coloane {
					vars[c.Cod] = r.Valori[i]
				}
				if ri < len(d.valoriRanduri) {
					for k, v := range d.valoriRanduri[ri] {
						vars[k] = v
					}
				}
				for _, e := range repetate {
					f.SetFont("Go", "", e.Font)
					v := inloc(e.Text)
					if float64(len(f.SplitText(v, e.Latime)))*e.Font*.45 > e.Inaltime {
						return nil, eroare("rapoarte.aspect_prea_mare")
					}
					f.SetXY(m+e.X, m+e.Y+offset)
					f.MultiCell(e.Latime, e.Font*.45, v, "", "L", false)
				}
				offset += d.Flexibil.InaltimeRand
			}
			for k, v := range initial {
				vars[k] = v
			}
		}
	}
	if d.Flexibil.ModAspect == "MARKDOWN" {
		lignes := strings.Split(d.Flexibil.Markdown, "\n")
		for idx := 0; idx < len(lignes); idx++ {
			linie := lignes[idx]
			if idx+1 < len(lignes) && strings.Contains(linie, "|") && strings.Contains(lignes[idx+1], "-") && strings.Trim(lignes[idx+1], "|:- \t") == "" {
				headers := strings.Split(strings.Trim(linie, " |"), "|")
				if len(headers) > 80 {
					f.SetError(eroare("rapoarte.aspect_prea_mare"))
					break
				}
				colsOld, rowsOld := d.Coloane, d.Randuri
				d.Coloane = []Coloana{}
				d.Randuri = []Rand{}
				for j, v := range headers {
					d.Coloane = append(d.Coloane, Coloana{Cod: fmt.Sprintf("c%d", j), Titlu: strings.TrimSpace(inloc(v))})
				}
				idx++
				for idx+1 < len(lignes) && strings.Contains(lignes[idx+1], "|") {
					idx++
					vs := strings.Split(strings.Trim(lignes[idx], " |"), "|")
					row := make([]string, len(headers))
					for j, v := range vs {
						if j < len(row) {
							row[j] = strings.TrimSpace(inloc(v))
						}
					}
					d.Randuri = append(d.Randuri, Rand{Valori: row})
				}
				tableY = f.GetY()
				tabel()
				d.Coloane, d.Randuri = colsOld, rowsOld
				f.Ln(2)
				continue
			}
			if strings.TrimSpace(linie) == "{{tabel}}" {
				tableY = f.GetY()
				tabel()
				continue
			}
			size := float64(d.Pagina.FontPT)
			style := ""
			v := inloc(linie)
			if strings.HasPrefix(v, "#") {
				n := len(v) - len(strings.TrimLeft(v, "#"))
				size += float64(5 - n)
				style = "B"
				v = strings.TrimSpace(strings.TrimLeft(v, "#"))
			}
			if strings.HasPrefix(v, "- ") {
				v = "\u2022 " + strings.TrimPrefix(v, "- ")
			}
			if strings.Contains(v, "**") {
				style = "B"
				v = strings.ReplaceAll(v, "**", "")
			}
			f.SetFont("Go", style, size)
			for _, l := range f.SplitText(v, w-2*m) {
				if f.GetY()+size*.5 > h-m-footerH-2 {
					nouaPagina()
					f.SetFont("Go", style, size)
				}
				f.SetX(m)
				f.CellFormat(w-2*m, size*.5, l, "", 1, "L", false, 0, "")
			}
			f.Ln(2)
		}
		f.SetFont("Go", "", float64(d.Pagina.FontPT))
	} else {
		tabel()
	}
	// Explicatiile si istoricul raman la finalul documentului.
	anexa = true
	nouaPagina()
	f.SetFont("Go", "B", 10)
	f.SetXY(m, startTabel)
	if len(d.Legenda) > 0 {
		f.MultiCell(w-2*m, 5, i18n.T(d.Limba, "rapoarte.anexa_calcul", nil), "", "L", false)
	}
	f.SetFont("Go", "", 8)
	linii := append([]string{}, d.Legenda...)
	linii = append(linii, i18n.T(d.Limba, "rapoarte.anexa_istoric", nil), d.Metadate)

	for _, v := range linii {
		for _, linie := range f.SplitText(v, w-2*m) {
			if f.GetY()+4 > h-m-footerH-2 {
				nouaPagina()
				f.SetFont("Go", "", 8)
				f.SetY(startTabel)
			}
			f.SetX(m)
			f.CellFormat(w-2*m, 4, linie, "", 1, "L", false, 0, "")
		}
		f.Ln(1)
	}
	var out bytes.Buffer
	if e := f.Output(&out); e != nil {
		return nil, eroare("rapoarte.aspect_prea_mare")
	}
	return out.Bytes(), nil
}
