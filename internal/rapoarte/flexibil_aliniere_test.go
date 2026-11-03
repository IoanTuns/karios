package rapoarte

import (
	"bytes"
	"context"
	"fmt"
	"github.com/pdfcpu/pdfcpu/pkg/api"
	"github.com/phpdave11/gofpdf"
	"golang.org/x/image/font/gofont/goregular"
	"io"
	"math"
	"os"
	"path/filepath"
	"pontix/internal/i18n"
	"regexp"
	"strconv"
	"strings"
	"testing"
)

// Verificam coordonatele din PDF-ul rezultat, nu doar optiunile rendererului.
func TestFlexibilAliniereZoneAntetSubsol(t *testing.T) {
	const k = 72.0 / 25.4
	reText := regexp.MustCompile(`(?s)BT ([-0-9.]+) ([-0-9.]+) Td \((.*?)\)\s*Tj ET`)
	reImg := regexp.MustCompile(`q ([-0-9.]+) 0 0 ([-0-9.]+) ([-0-9.]+) ([-0-9.]+) cm /I\S+ Do Q`)
	num := func(s string) float64 {
		v, e := strconv.ParseFloat(s, 64)
		if e != nil {
			t.Fatal(e)
		}
		return v
	}
	masura := gofpdf.New("P", "mm", "A4", "")
	masura.AddUTF8FontFromBytes("Go", "", goregular.TTF)
	masura.SetFont("Go", "", 8)
	codifica := func(s string) string {
		var b strings.Builder
		for _, r := range s {
			b.WriteByte(0)
			b.WriteByte(byte(r))
		}
		return b.String()
	}
	for _, limba := range []i18n.Limba{i18n.Romana, i18n.Engleza} {
		for _, pagina := range []Pagina{{Format: "A4", Orientare: "PORTRET", MargineMM: 10, FontPT: 8}, {Format: "A4", Orientare: "PEISAJ", MargineMM: 10, FontPT: 8}, {Format: "A3", Orientare: "PEISAJ", MargineMM: 15, FontPT: 8}} {
			for _, mod := range []string{"TABEL", "MARKDOWN"} {
				t.Run(string(limba)+"-"+pagina.Format+"-"+pagina.Orientare+"-"+mod, func(t *testing.T) {
					im := imagineFictiva(t, "png")
					d := Document{Titlu: "Raport demonstrativ", Limba: limba, Pagina: pagina, Flexibil: &RaportFlexibil{ModAspect: mod, Markdown: "{{tabel}}"}, Antet: []string{"A_STANGA", "A_CENTRU\n" + strings.Repeat("Titlu de garda cu mai multe cuvinte. ", 4) + "\nA_ULTIMA", "A_DREAPTA"}, Subsol: []string{"S_STANGA", "S_CENTRU", "S_DREAPTA"}, AntetImagini: map[string]Imagine{"stanga": im, "centru": im, "dreapta": im}, SubsolImagini: map[string]Imagine{"stanga": im, "centru": im, "dreapta": im}, Coloane: []Coloana{{Cod: "data", Titlu: "DATA"}, {Cod: "nume", Titlu: "NUME"}}, Legenda: []string{"Explicatie demonstrativa."}, Metadate: "Istoric demonstrativ."}
					for i := 0; i < 160; i++ {
						d.Randuri = append(d.Randuri, Rand{Valori: []string{fmt.Sprint(i + 1), "Angajat demo"}})
					}
					pdf, e := GenereazaPDF(d)
					if e != nil {
						t.Fatal(e)
					}
					w := 210.0
					if pagina.Format == "A3" {
						w = 297
					}
					if pagina.Orientare == "PEISAJ" {
						if pagina.Format == "A3" {
							w = 420
						} else {
							w = 297
						}
					}
					pagini := 0
					e = api.ExtractContent(context.Background(), bytes.NewReader(pdf), nil, func(r io.Reader, n int) error {
						pagini++
						raw, e := io.ReadAll(r)
						if e != nil {
							return e
						}
						stream := string(raw)
						poz := map[string][2]float64{}
						for _, m := range reText.FindAllStringSubmatch(stream, -1) {
							for _, tag := range []string{"A_STANGA", "A_CENTRU", "A_DREAPTA", "S_STANGA", "S_CENTRU", "S_DREAPTA", "A_ULTIMA", "DATA"} {
								if m[3] == codifica(tag) {
									poz[tag] = [2]float64{num(m[1]), num(m[2])}
								}
							}
						}
						for _, prefix := range []string{"A_", "S_"} {
							for i, sufix := range []string{"STANGA", "CENTRU", "DREAPTA"} {
								tag := prefix + sufix
								p, ok := poz[tag]
								if !ok {
									t.Fatalf("pagina%d: lipseste %s", n, tag)
								}
								textW := masura.GetStringWidth(tag) * k
								tinta := (float64(pagina.MargineMM) + masura.GetCellMargin()) * k
								if i == 1 {
									tinta = w*k/2 - textW/2
								}
								if i == 2 {
									tinta = (w-float64(pagina.MargineMM)-masura.GetCellMargin())*k - textW
								}
								if math.Abs(p[0]-tinta) > .08 {
									t.Fatalf("pagina%d %s: x %.2f, asteptat %.2f", n, tag, p[0], tinta)
								}
								if math.Abs(p[1]-poz[prefix+"STANGA"][1]) > .08 {
									t.Fatalf("pagina%d: zonele nu sunt simultane", n)
								}
							}
						}
						if data, ok := poz["DATA"]; ok && data[1] >= poz["A_ULTIMA"][1] {
							t.Fatalf("pagina%d: tabel peste antet", n)
						}
						imgs := reImg.FindAllStringSubmatch(stream, -1)
						if len(imgs) != 6 {
							t.Fatalf("pagina%d: imagini %d, asteptat6", n, len(imgs))
						}
						for i, m := range imgs {
							iw, x := num(m[1]), num(m[3])
							tinta := float64(pagina.MargineMM) * k
							if i%3 == 1 {
								tinta = w*k/2 - iw/2
							}
							if i%3 == 2 {
								tinta = (w-float64(pagina.MargineMM))*k - iw
							}
							if math.Abs(x-tinta) > .08 {
								t.Fatalf("pagina%d: imagine%d x %.2f, asteptat %.2f", n, i, x, tinta)
							}
						}
						return nil
					}, nil)
					if e != nil {
						t.Fatal(e)
					}
					if pagini < 3 {
						t.Fatal("paginare insuficienta", pagini)
					}
					if *directorDemo != "" {
						if e := os.MkdirAll(*directorDemo, 0700); e != nil {
							t.Fatal(e)
						}
						n := fmt.Sprintf("aliniere-%s-%s-%s-%s.pdf", limba, pagina.Format, pagina.Orientare, mod)
						if e := os.WriteFile(filepath.Join(*directorDemo, n), pdf, 0600); e != nil {
							t.Fatal(e)
						}
					}
				})
			}
		}
	}
}
