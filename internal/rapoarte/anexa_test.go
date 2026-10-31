package rapoarte

import (
	"context"
	"fmt"
	"github.com/johnfercher/go-tree/node"
	"github.com/johnfercher/maroto/v2/pkg/core"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func textePaginaAnexa(n *node.Node[core.Structure]) string {
	var texte []string
	var viziteaza func(*node.Node[core.Structure])
	viziteaza = func(n *node.Node[core.Structure]) {
		d := n.GetData()
		if d.Type == "text" {
			texte = append(texte, fmt.Sprint(d.Value))
		}
		for _, f := range n.GetNexts() {
			viziteaza(f)
		}
	}
	viziteaza(n)
	return strings.Join(texte, "\n")
}
func TestAnexaDoarDupaTabelCuIstoricPeUltimaPagina(t *testing.T) {
	m := pregateste(t)
	fixture(t, m)
	for _, limba := range []string{"ro", "en"} {
		for _, randuri := range []int{0, 1, 150} {
			for _, explicatii := range []int{1, 200} {
				t.Run(fmt.Sprintf("%s-%d-%d", limba, randuri, explicatii), func(t *testing.T) {
					d, err := m.s.Pregateste(context.Background(), Cerere{ModelID: "raport-pontaj", An: 2026, Luna: 10, Limba: limba}, momentFix)
					if err != nil {
						t.Fatal(err)
					}
					d.Randuri = nil
					for i := 0; i < randuri; i++ {
						v := make([]string, len(d.Coloane))
						v[0] = "RAND_TABEL_UNIC"
						d.Randuri = append(d.Randuri, Rand{Valori: v})
					}
					d.Metadate = "ISTORIC_UNIC"
					d.Legenda = nil
					for i := 0; i < explicatii; i++ {
						d.Legenda = append(d.Legenda, fmt.Sprintf("CALCUL_UNIC_%03d: explicatie pentru verificarea paginarii.", i))
					}
					structura, err := ConstruiestePDF(d)
					if err != nil {
						t.Fatal(err)
					}
					pagini := structura.GetStructure().GetNexts()
					if len(pagini) < 2 {
						t.Fatal("anexa nu are pagina separata")
					}
					vazutaAnexa := false
					numarIstoric := 0
					for i, p := range pagini {
						txt := textePaginaAnexa(p)
						numarIstoric += strings.Count(txt, "ISTORIC_UNIC")
						if strings.Contains(txt, "CALCUL_UNIC") {
							vazutaAnexa = true
							if i == 0 {
								t.Fatal("calcul pe pagina tabelului")
							}
						}
						if vazutaAnexa && strings.Contains(txt, "RAND_TABEL_UNIC") {
							t.Fatal("rand de tabel in anexa")
						}
						if i < len(pagini)-1 && strings.Contains(txt, "ISTORIC_UNIC") {
							t.Fatal("istoric inaintea ultimei pagini")
						}
					}
					if !vazutaAnexa || numarIstoric != 1 {
						t.Fatal("calcul/istoric lipsa sau repetat")
					}
					pdf, err := GenereazaPDF(d)
					if err != nil {
						t.Fatal(err)
					}
					n, err := NumaraPaginiPDF(context.Background(), pdf)
					if err != nil || n != len(pagini) {
						t.Fatalf("PDF %d structura %d: %v", n, len(pagini), err)
					}
					if *directorDemo != "" && randuri == 1 && explicatii == 1 {
						if err := os.WriteFile(filepath.Join(*directorDemo, "anexa-"+limba+".pdf"), pdf, 0600); err != nil {
							t.Fatal(err)
						}
					}
				})
			}
		}
	}
}

func TestAnexaRaportDetaliatReal(t *testing.T) {
	m := pregateste(t)
	fixture(t, m)
	for _, limba := range []string{"ro", "en"} {
		d, err := m.s.Pregateste(context.Background(), Cerere{ModelID: "raport-pontaj-detaliat", An: 2026, Luna: 10, Limba: limba}, momentFix)
		if err != nil {
			t.Fatal(err)
		}
		str, err := ConstruiestePDF(d)
		if err != nil {
			t.Fatal(err)
		}
		pagini := str.GetStructure().GetNexts()
		ultim := textePaginaAnexa(pagini[len(pagini)-1])
		if !strings.Contains(ultim, d.Metadate) {
			t.Fatal("istoric lipsa de pe ultima pagina")
		}
		for _, l := range d.Legenda {
			if !strings.Contains(ultim, l) {
				t.Fatalf("explicatie lipsa de pe ultima pagina: %s", l)
			}
		}
		pdf, err := GenereazaPDF(d)
		if err != nil {
			t.Fatal(err)
		}
		if *directorDemo != "" {
			if err := os.WriteFile(filepath.Join(*directorDemo, "detaliat-anexa-"+limba+".pdf"), pdf, 0600); err != nil {
				t.Fatal(err)
			}
		}
		t.Logf("%s: %d pagini, %d explicatii", limba, len(pagini), len(d.Legenda))
	}
}
