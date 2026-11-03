package rapoarte

import (
	"bytes"
	"context"
	"github.com/xuri/excelize/v2"
	"pontix/internal/i18n"
	"reflect"
	"strings"
	"testing"
)

func TestExplicatiiOptionalePersistentaSiExporturi(t *testing.T) {
	m := pregateste(t)
	fixture(t, m)
	ctx := context.Background()
	for _, limba := range []string{"ro", "en"} {
		for _, id := range []string{"raport-pontaj", "raport-pontaj-detaliat", "raport-fisa", "raport-absente", "raport-garzi", "flexibil"} {
			t.Run(limba+"-"+id, func(t *testing.T) {
				var d Definitie
				if id == "flexibil" {
					d = modelFlexTest()
				} else {
					d = initiala(t, m, id)
					d.ID = ""
					d.Revizie = ""
				}
				c := Cerere{An: 2026, Luna: 10, Limba: limba, AngajatID: "a"}
				var etalon Document
				for _, ascunde := range []bool{false, true, false} {
					d.Continut.AscundeExplicatii = ascunde
					salvat, e := m.s.Salveaza(ctx, d)
					if e != nil {
						t.Fatal(e)
					}
					d = initiala(t, m, salvat)
					if d.Continut.AscundeExplicatii != ascunde {
						t.Fatal("optiune nepersistata")
					}
					c.ModelID = salvat
					doc, e := m.s.Pregateste(ctx, c, momentFix)
					if e != nil {
						t.Fatal(e)
					}
					if (len(doc.Legenda) == 0) != ascunde {
						t.Fatal("legenda nu respecta alegerea")
					}
					if doc.Metadate == "" {
						t.Fatal("istoric pierdut")
					}
					if etalon.Sursa == "" {
						etalon = doc
					} else if !reflect.DeepEqual(etalon.Coloane, doc.Coloane) || !reflect.DeepEqual(etalon.Randuri, doc.Randuri) {
						t.Fatal("optiunea schimba datele")
					}
					if doc.Flexibil == nil {
						p, e := ConstruiestePDF(doc)
						if e != nil {
							t.Fatal(e)
						}
						txt := textePaginaAnexa(p.GetStructure())
						titlu := i18n.T(i18n.Limba(limba), "rapoarte.anexa_calcul", nil)
						if strings.Contains(txt, titlu) == ascunde {
							t.Fatal("titlu calcul nepotrivit in PDF")
						}
					}
					if _, e := GenereazaPDF(doc); e != nil {
						t.Fatal(e)
					}
					x, e := GenereazaExcel(doc)
					if e != nil {
						t.Fatal(e)
					}
					f, e := excelize.OpenReader(bytes.NewReader(x))
					if e != nil {
						t.Fatal(e)
					}
					randuri, e := f.GetRows(f.GetSheetName(0))
					f.Close()
					if e != nil {
						t.Fatal(e)
					}
					var celule []string
					for _, r := range randuri {
						celule = append(celule, r...)
					}
					for _, legenda := range etalon.Legenda {
						if strings.Contains(strings.Join(celule, "\n"), legenda) == ascunde {
							t.Fatal("explicatii nepotrivite in Excel")
						}
					}
				}
			})
		}
	}
}
