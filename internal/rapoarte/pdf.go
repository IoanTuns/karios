package rapoarte

import (
	"github.com/johnfercher/maroto/v2"
	"github.com/johnfercher/maroto/v2/pkg/components/row"
	"github.com/johnfercher/maroto/v2/pkg/components/text"
	"github.com/johnfercher/maroto/v2/pkg/config"
	"github.com/johnfercher/maroto/v2/pkg/consts/align"
	"github.com/johnfercher/maroto/v2/pkg/consts/border"
	"github.com/johnfercher/maroto/v2/pkg/consts/fontstyle"
	"github.com/johnfercher/maroto/v2/pkg/consts/orientation"
	"github.com/johnfercher/maroto/v2/pkg/consts/pagesize"
	"github.com/johnfercher/maroto/v2/pkg/core"
	"github.com/johnfercher/maroto/v2/pkg/fontrepository"
	"github.com/johnfercher/maroto/v2/pkg/props"
	"golang.org/x/image/font/gofont/gobold"
	"golang.org/x/image/font/gofont/goregular"
)

// ConstruiestePDF este folosit și de testele de referință ale structurii raportului.
func ConstruiestePDF(d Document) (core.Maroto, error) {
	fonturi, err := fontrepository.New().AddUTF8FontFromBytes("Go", fontstyle.Normal, goregular.TTF).AddUTF8FontFromBytes("Go", fontstyle.Bold, gobold.TTF).Load()
	if err != nil {
		return nil, traduce(err)
	}
	dimensiune := pagesize.A4
	if d.Pagina.Format == "A3" {
		dimensiune = pagesize.A3
	}
	orientare := orientation.Vertical
	if d.Pagina.Orientare == "PEISAJ" {
		orientare = orientation.Horizontal
	}
	cfg := config.NewBuilder().WithPageSize(dimensiune).WithOrientation(orientare).
		WithLeftMargin(float64(d.Pagina.MargineMM)).WithRightMargin(float64(d.Pagina.MargineMM)).
		WithTopMargin(float64(d.Pagina.MargineMM)).WithBottomMargin(float64(d.Pagina.MargineMM)).
		WithDefaultFont(&props.Font{Family: "Go", Size: float64(d.Pagina.FontPT)}).WithCustomFonts(fonturi).
		WithMaxGridSize(120).WithCreationDate(d.GeneratLa).WithTitle(d.Titlu, true).
		WithCreator("Pontix", true).WithCompression(true).
		WithPageNumber(props.PageNumber{Pattern: d.Paginare, Place: props.Bottom, Family: "Go", Style: fontstyle.Normal, Size: 7}).Build()
	m := maroto.New(cfg)
	bloc := func(valori []string, marime float64, stil fontstyle.Type) core.Row {
		r := row.New()
		for i, v := range valori {
			aliniere := align.Left
			if i == 1 {
				aliniere = align.Center
			}
			if i == 2 {
				aliniere = align.Right
			}
			r.Add(text.NewCol(40, v, props.Text{Family: "Go", Style: stil, Size: marime, Align: aliniere, Top: 1, Bottom: 2, Left: 1, Right: 1}))
		}
		return r
	}
	culori := map[string]*props.Color{"WEEKEND": {Red: 231, Green: 234, Blue: 231}, "SARBATOARE": {Red: 255, Green: 240, Blue: 199}, "ABSENTA": {Red: 223, Green: 239, Blue: 228}}
	suma := 0
	for _, c := range d.Coloane {
		suma += c.Pondere
	}
	latimi := make([]int, len(d.Coloane))
	ramas := 120
	for i, c := range d.Coloane {
		latimi[i] = c.Pondere * 120 / suma
		if latimi[i] < 1 {
			latimi[i] = 1
		}
		ramas -= latimi[i]
	}
	for i := 0; ramas > 0; i++ {
		latimi[i%len(latimi)]++
		ramas--
	}
	rand := func(valori []string, fundal []string, antet bool) core.Row {
		r := row.New()
		for i, v := range valori {
			stil := fontstyle.Normal
			culoare := (*props.Color)(nil)
			if antet {
				stil = fontstyle.Bold
				culoare = &props.Color{Red: 231, Green: 239, Blue: 234}
			}
			if i < len(fundal) && culori[fundal[i]] != nil {
				culoare = culori[fundal[i]]
			}
			aliniere := align.Left
			if len(d.Coloane[i].Cod) == 10 && d.Coloane[i].Cod[4] == '-' {
				aliniere = align.Center
			}
			col := text.NewCol(latimi[i], v, props.Text{Family: "Go", Style: stil, Size: float64(d.Pagina.FontPT), Align: aliniere, Top: 1, Bottom: 1.5, Left: 0.7, Right: 0.7})
			col.WithStyle(&props.Cell{BorderType: border.Full, BorderColor: &props.Color{Red: 140, Green: 150, Blue: 144}, BorderThickness: 0.15, BackgroundColor: culoare})
			r.Add(col)
		}
		return r
	}
	titluri, fundal := []string{}, []string{}
	for _, c := range d.Coloane {
		titluri = append(titluri, c.Titlu)
		fundal = append(fundal, c.Culoare)
	}
	antet := []core.Row{bloc(d.Antet, float64(d.Pagina.FontPT)+2, fontstyle.Bold)}
	antet = append(antet, text.NewAutoRow(d.Nota, props.Text{Family: "Go", Size: 6, Top: 1, Bottom: 2}))
	if d.Introducere != "" {
		antet = append(antet, text.NewAutoRow(d.Introducere, props.Text{Family: "Go", Size: float64(d.Pagina.FontPT), Top: 2, Bottom: 3}))
	}
	antet = append(antet, rand(titluri, fundal, true))
	if err := m.RegisterHeader(antet...); err != nil {
		return nil, eroare("rapoarte.aspect_prea_mare")
	}
	if err := m.RegisterFooter(bloc(d.Subsol, float64(d.Pagina.FontPT), fontstyle.Normal), text.NewAutoRow(d.Metadate, props.Text{Family: "Go", Size: 6, Top: 2, Bottom: 2})); err != nil {
		return nil, eroare("rapoarte.aspect_prea_mare")
	}
	for _, r := range d.Randuri {
		m.AddRows(rand(r.Valori, r.Culori, false))
	}
	if len(d.Randuri) == 0 {
		m.AddRows(text.NewAutoRow(d.FaraDate, props.Text{Family: "Go", Size: float64(d.Pagina.FontPT), Top: 4, Bottom: 4}))
	}
	return m, nil
}

func GenereazaPDF(d Document) ([]byte, error) {
	if len(d.Coloane) == 0 {
		return nil, eroare("rapoarte.continut_invalid")
	}
	m, err := ConstruiestePDF(d)
	if err != nil {
		return nil, err
	}
	pdf, err := m.Generate()
	if err != nil {
		return nil, eroare("rapoarte.aspect_prea_mare")
	}
	return pdf.GetBytes(), nil
}
