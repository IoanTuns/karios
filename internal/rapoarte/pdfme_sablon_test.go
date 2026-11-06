package rapoarte

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"testing"
)

// egalJSON compara doua valori JSON, cu toleranta pentru numerele reale
// (impartirile din TypeScript si Go pot diferi pe ultima zecimala).
func egalJSON(a, b any, cale string) error {
	switch x := a.(type) {
	case map[string]any:
		y, ok := b.(map[string]any)
		if !ok || len(x) != len(y) {
			return fmt.Errorf("%s: obiecte diferite (%d/%d chei)", cale, len(x), len(y))
		}
		for k := range x {
			if err := egalJSON(x[k], y[k], cale+"."+k); err != nil {
				return err
			}
		}
	case []any:
		y, ok := b.([]any)
		if !ok || len(x) != len(y) {
			return fmt.Errorf("%s: liste diferite", cale)
		}
		for i := range x {
			if err := egalJSON(x[i], y[i], fmt.Sprintf("%s[%d]", cale, i)); err != nil {
				return err
			}
		}
	case float64:
		y, ok := b.(float64)
		if !ok || math.Abs(x-y) > 1e-9 {
			return fmt.Errorf("%s: %v fata de %v", cale, a, b)
		}
	default:
		if a != b {
			return fmt.Errorf("%s: %v fata de %v", cale, a, b)
		}
	}
	return nil
}

// TestSablonGoReproduceSabloaneleDinMigrare: constructorul Go da aceleasi
// sabloane ca generatorul TypeScript al migrarii 0033.
func TestSablonGoReproduceSabloaneleDinMigrare(t *testing.T) {
	m := pregateste(t)
	c, err := m.s.Lista(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	nume := map[string]string{"antet_stanga": "Antet stânga", "antet_centru": "Antet centru", "antet_dreapta": "Antet dreapta", "titlu": "Titlu", "subsol_stanga": "Subsol stânga", "subsol_centru": "Subsol centru", "pagina": "Număr pagină"}
	for id, fisier := range map[string]string{"raport-pontaj-sql": "raport-pontaj-v17.json", "raport-pontaj-detaliat-sql": "raport-pontaj-detaliat-v18.json", "raport-fisa-sql": "raport-fisa-v17.json", "raport-absente-sql": "raport-absente-v14.json", "raport-garzi-sql": "raport-garzi-v14.json"} {
		var d *Definitie
		for i := range c.Modele {
			if c.Modele[i].ID == id {
				d = &c.Modele[i]
			}
		}
		if d == nil {
			t.Fatalf("model lipsa %s", id)
		}
		b, err := os.ReadFile(filepath.Join("testdata", fisier))
		if err != nil {
			t.Fatal(err)
		}
		var ref struct {
			Document struct {
				Coloane []struct {
					Cod, Titlu string
					Pondere    float64
				}
			}
		}
		if err := json.Unmarshal(b, &ref); err != nil {
			t.Fatal(err)
		}
		coloane := []ColoanaSablon{}
		for _, col := range ref.Document.Coloane {
			coloane = append(coloane, ColoanaSablon{Cod: codSQLStandard(col.Cod), Titlu: col.Titlu, Pondere: col.Pondere})
		}
		goSablon := SablonDinColoane(coloane, d.Pagina, OptiuniSablon{TitluriDinRaport: true, StilClasic: true, Nume: nume, TextPagina: "Pagina {currentPage} din {totalPages}"})
		var a, e any
		if err := json.Unmarshal([]byte(goSablon), &a); err != nil {
			t.Fatal(err)
		}
		if err := json.Unmarshal([]byte(d.Continut.Flexibil.Sablon), &e); err != nil {
			t.Fatal(err)
		}
		if err := egalJSON(a, e, id); err != nil {
			t.Errorf("sablon Go diferit: %v", err)
		}
	}
}

// Un tabel marcat pontixTitluri urmeaza SQL-ul schimbat: coloanele pastrate
// isi pastreaza titlul si latimea, cele noi se adauga, iar datele se potrivesc.
func TestSablonSeAdapteazaLaColoaneleSQL(t *testing.T) {
	m := pregateste(t)
	fixture(t, m)
	ctx := context.Background()
	c, err := m.s.Lista(ctx)
	if err != nil {
		t.Fatal(err)
	}
	var d Definitie
	for _, x := range c.Modele {
		if x.ID == "raport-pontaj-sql" {
			d = x
		}
	}
	d.ID = ""
	d.Continut.Flexibil.SQL = "SELECT rand, angajat, z01, z02, functie FROM pontaj_lunar ORDER BY rand"
	doc, err := m.s.PrevizualizeazaFlexibil(ctx, d, Cerere{An: 2026, Luna: 10, Limba: "ro"}, momentFix)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := DatePdfme(doc); err != nil {
		t.Fatalf("datele nu se potrivesc cu sablonul adaptat: %v", err)
	}
	text, err := SablonPentruDocument(doc)
	if err != nil {
		t.Fatal(err)
	}
	var vechi, nou struct {
		Schemas [][]struct {
			Name                 string
			Head                 []string
			HeadWidthPercentages []float64
			PontixColoane        []string
		}
	}
	_ = json.Unmarshal([]byte(d.Continut.Flexibil.Sablon), &vechi)
	_ = json.Unmarshal([]byte(text), &nou)
	gaseste := func(s [][]struct {
		Name                 string
		Head                 []string
		HeadWidthPercentages []float64
		PontixColoane        []string
	}) (h []string, l []float64, c []string) {
		for _, e := range s[0] {
			if e.Name == "tabel" {
				return e.Head, e.HeadWidthPercentages, e.PontixColoane
			}
		}
		return
	}
	hv, lv, _ := gaseste(vechi.Schemas)
	hn, ln, cn := gaseste(nou.Schemas)
	if fmt.Sprint(cn) != "[angajat z01 z02 functie]" || len(hn) != 4 || hn[0] != hv[0] || hn[3] != "Funcție" {
		t.Fatalf("tabel adaptat: %v %v", hn, cn)
	}
	// Proportia dintre angajat si z01 ramane cea din sablonul initial.
	if math.Abs(ln[0]/ln[1]-lv[0]/lv[2]) > 1e-9 {
		t.Fatalf("latimi: %v fata de %v", ln, lv)
	}
	total := 0.0
	for _, v := range ln {
		total += v
	}
	if math.Abs(total-100) > 1e-9 {
		t.Fatalf("total latimi %v", total)
	}
}
