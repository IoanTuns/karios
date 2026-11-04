package rapoarte

import (
	"github.com/johnfercher/maroto/v2/pkg/consts/align"
	"github.com/johnfercher/maroto/v2/pkg/consts/fontstyle"
	"github.com/johnfercher/maroto/v2/pkg/props"
	"strings"
)

func stilProps(d Document, p props.Text, chei ...string) props.Text {
	if len(d.Stiluri) == 0 {
		return p
	}
	a := string(p.Align)
	if a == "" {
		a = "L"
	}
	s := StilPentru(d, p.Size, a, strings.Contains(string(p.Style), "B"), chei...)
	if strings.Contains(string(p.Style), "I") {
		s.Cursiv = true
	}
	p.Size = s.FontPT
	p.Style = fontstyle.Type(s.FontPDF())
	p.Align = align.Type(s.Aliniere)
	return p
}
