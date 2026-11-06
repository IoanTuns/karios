package rapoarte

import (
	"context"
	"encoding/json"
	"reflect"
	"strings"
	"testing"
)

// TestModelePredefiniteSQLEchivalente: fiecare model predefinit SQL + pdfme (D090)
// da aceleasi titluri, valori si culori ca modelul standard vechi, iar sablonul
// lui are cate o coloana de tabel pentru fiecare coloana vizibila.
func TestModelePredefiniteSQLEchivalente(t *testing.T) {
	m := pregateste(t)
	fixture(t, m)
	ctx := context.Background()
	for _, caz := range []struct{ vechi, nou, angajat string }{
		{"raport-pontaj", "raport-pontaj-sql", ""},
		{"raport-pontaj-detaliat", "raport-pontaj-detaliat-sql", ""},
		{"raport-fisa", "raport-fisa-sql", "a"},
		{"raport-absente", "raport-absente-sql", ""},
		{"raport-garzi", "raport-garzi-sql", ""},
	} {
		t.Run(caz.nou, func(t *testing.T) {
			for _, limba := range []string{"ro", "en"} {
				c := Cerere{ModelID: caz.vechi, An: 2026, Luna: 10, Limba: limba, AngajatID: caz.angajat}
				vechi, err := m.s.Pregateste(ctx, c, momentFix)
				if err != nil {
					t.Fatal(err)
				}
				c.ModelID = caz.nou
				nou, err := m.s.Pregateste(ctx, c, momentFix)
				if err != nil {
					t.Fatal(err)
				}
				if !EstePdfme(nou) {
					t.Fatal("modelul nou nu foloseste editorul vizual")
				}
				if len(nou.Coloane) != len(vechi.Coloane) || len(nou.Randuri) != len(vechi.Randuri) {
					t.Fatalf("%s: dimensiuni %d/%d coloane, %d/%d randuri", limba, len(nou.Coloane), len(vechi.Coloane), len(nou.Randuri), len(vechi.Randuri))
				}
				for j := range vechi.Coloane {
					if nou.Coloane[j].Titlu != vechi.Coloane[j].Titlu {
						t.Errorf("%s coloana %d: %q fata de %q", limba, j, nou.Coloane[j].Titlu, vechi.Coloane[j].Titlu)
					}
				}
				for i := range vechi.Randuri {
					if !reflect.DeepEqual(nou.Randuri[i].Valori, vechi.Randuri[i].Valori) || strings.Join(nou.Randuri[i].Culori, "|") != strings.Join(vechi.Randuri[i].Culori, "|") && strings.Join(vechi.Randuri[i].Culori, "") != "" {
						t.Errorf("%s rand %d diferit", limba, i)
					}
				}
				intrari, err := DatePdfme(nou)
				if err != nil {
					t.Fatalf("%s: datele pdfme: %v", limba, err)
				}
				var st stiluriTabelPdfme
				if err := json.Unmarshal([]byte(intrari[IntrareStiluriPdfme]), &st); err != nil {
					t.Fatal(err)
				}
				if len(st.Titluri) != len(vechi.Coloane) || st.Titluri[0] != vechi.Coloane[0].Titlu {
					t.Fatalf("%s: titluri %q", limba, st.Titluri)
				}
				// Octombrie are 31 de zile: nicio coloana de zi ascunsa.
				if len(st.Ascunse) != 0 {
					t.Fatalf("%s: zile ascunse in octombrie: %v", limba, st.Ascunse)
				}
			}
		})
	}
}

// In februarie coloanele z29-z31 se ascund in PDF; valorile raman goale.
func TestModelPredefinitFebruarieAscundeZilele(t *testing.T) {
	m := pregateste(t)
	fixture(t, m)
	doc, err := m.s.Pregateste(context.Background(), Cerere{ModelID: "raport-pontaj-sql", An: 2026, Luna: 2, Limba: "ro"}, momentFix)
	if err != nil {
		t.Fatal(err)
	}
	intrari, err := DatePdfme(doc)
	if err != nil {
		t.Fatal(err)
	}
	var st stiluriTabelPdfme
	if err := json.Unmarshal([]byte(intrari[IntrareStiluriPdfme]), &st); err != nil {
		t.Fatal(err)
	}
	asteptat := []string{"z29", "z30", "z31"}
	if len(st.Ascunse) != 3 {
		t.Fatalf("ascunse: %v", st.Ascunse)
	}
	for k, j := range st.Ascunse {
		if doc.Coloane[j].Cod != asteptat[k] {
			t.Fatalf("coloana ascunsa %d este %s", j, doc.Coloane[j].Cod)
		}
	}
}
