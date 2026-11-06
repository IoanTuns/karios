package rapoarte

import (
	"context"
	"encoding/json"
	"reflect"
	"strconv"
	"strings"
	"testing"

	"pontix/internal/i18n"
)

// TestColoaneModelStandardCaRaportul: conversia foloseste aceleasi coloane si
// ponderi ca raportul standard (fisierele de referinta).
func TestColoaneModelStandardCaRaportul(t *testing.T) {
	m := pregateste(t)
	fixture(t, m)
	ctx := context.Background()
	c, err := m.s.Lista(ctx)
	if err != nil {
		t.Fatal(err)
	}
	for _, d := range c.Modele {
		if _, ok := tabelPentruSursa[d.Continut.Sursa]; !ok {
			continue
		}
		doc, err := m.s.Pregateste(ctx, Cerere{ModelID: d.ID, An: 2026, Luna: 10, Limba: "ro", AngajatID: "a"}, momentFix)
		if err != nil {
			t.Fatal(d.ID, err)
		}
		cols := coloaneModelStandard(d, i18n.Romana)
		if len(cols) != len(doc.Coloane) {
			t.Fatalf("%s: %d coloane fata de %d", d.ID, len(cols), len(doc.Coloane))
		}
		for i, col := range doc.Coloane {
			if cols[i].Cod != codSQLStandard(col.Cod) || cols[i].Pondere != float64(col.Pondere) {
				t.Errorf("%s coloana %d: %+v fata de %+v", d.ID, i, cols[i], col)
			}
		}
	}
}

// TestConversieModeleStandardEchivalenta: dupa ConvertesteModeleVechi, fiecare
// model standard da aceleasi valori, titluri si culori ca inainte, iar
// definitia veche ramane in jurnal.
func TestConversieModeleStandardEchivalenta(t *testing.T) {
	m := pregateste(t)
	fixture(t, m)
	ctx := context.Background()
	vechi := map[string]Document{}
	for _, id := range []string{"raport-pontaj", "raport-pontaj-detaliat", "raport-fisa", "raport-absente", "raport-garzi"} {
		doc, err := m.s.Pregateste(ctx, Cerere{ModelID: id, An: 2026, Luna: 10, Limba: "en", AngajatID: "a"}, momentFix)
		if err != nil {
			t.Fatal(err)
		}
		vechi[id] = doc
	}
	rez, err := m.s.ConvertesteModeleVechi(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if len(rez) != 5 {
		t.Fatalf("convertite: %+v", rez)
	}
	for _, r := range rez {
		if r.Eroare != "" {
			t.Fatalf("%s: %s", r.ID, r.Eroare)
		}
	}
	for id, v := range vechi {
		nou, err := m.s.Pregateste(ctx, Cerere{ModelID: id, An: 2026, Luna: 10, Limba: "en", AngajatID: "a"}, momentFix)
		if err != nil {
			t.Fatal(id, err)
		}
		revVeche, _ := strconv.Atoi(v.RevizieModel)
		if !EstePdfme(nou) || nou.RevizieModel != strconv.Itoa(revVeche+1) {
			t.Fatalf("%s: neconvertit (revizie %s)", id, nou.RevizieModel)
		}
		if len(nou.Coloane) != len(v.Coloane) || len(nou.Randuri) != len(v.Randuri) {
			t.Fatalf("%s: dimensiuni diferite", id)
		}
		for j := range v.Coloane {
			if nou.Coloane[j].Titlu != v.Coloane[j].Titlu {
				t.Errorf("%s coloana %d: %q fata de %q", id, j, nou.Coloane[j].Titlu, v.Coloane[j].Titlu)
			}
		}
		for i := range v.Randuri {
			if !reflect.DeepEqual(nou.Randuri[i].Valori, v.Randuri[i].Valori) {
				t.Errorf("%s rand %d: valori diferite", id, i)
			}
		}
		intrari, err := DatePdfme(nou)
		if err != nil {
			t.Fatalf("%s: %v", id, err)
		}
		// Textele finale de antet/subsol sunt cele ale raportului vechi.
		final := []string{intrari["antet_stanga"], intrari["antet_centru"], intrari["antet_dreapta"]}
		if nou.Titlu != v.Titlu || !reflect.DeepEqual(final, v.Antet) || intrari["subsol_dreapta"] != v.Subsol[2] {
			t.Errorf("%s: antet %q fata de %q", id, final, v.Antet)
		}
	}
	var n int
	if err := m.db.QueryRow(`SELECT count(*) FROM jurnal_modificari WHERE tabel='definitii_rapoarte' AND rand_id='raport-pontaj' AND valoare_veche LIKE '%PONTAJ_LUNAR%' AND valoare_noua LIKE '%SQL_LIBER%'`).Scan(&n); err != nil || n != 1 {
		t.Fatalf("originalul nu este in jurnal: %d %v", n, err)
	}
	// A doua rulare nu mai gaseste nimic de convertit.
	if rez, err := m.s.ConvertesteModeleVechi(ctx); err != nil || len(rez) != 0 {
		t.Fatalf("a doua conversie: %+v %v", rez, err)
	}
}

func salveazaSiConverteste(t *testing.T, m *mediu, d Definitie) (Definitie, []string) {
	t.Helper()
	ctx := context.Background()
	id, err := m.s.Salveaza(ctx, d)
	if err != nil {
		t.Fatal(err)
	}
	rez, err := m.s.ConvertesteModeleVechi(ctx)
	if err != nil {
		t.Fatal(err)
	}
	var note []string
	for _, r := range rez {
		if r.Eroare != "" {
			t.Fatalf("%s: %s", r.ID, r.Eroare)
		}
		if r.ID == id {
			note = r.Note
		}
	}
	c, err := m.s.Lista(ctx)
	if err != nil {
		t.Fatal(err)
	}
	for _, x := range c.Modele {
		if x.ID == id {
			return x, note
		}
	}
	t.Fatal("model lipsa")
	return Definitie{}, nil
}

func elementeSablon(t *testing.T, d Definitie) []map[string]any {
	t.Helper()
	var s struct{ Schemas [][]map[string]any }
	if err := json.Unmarshal([]byte(d.Continut.Flexibil.Sablon), &s); err != nil {
		t.Fatal(err)
	}
	return s.Schemas[0]
}

func TestConversieModeleFlexibile(t *testing.T) {
	m := pregateste(t)
	fixture(t, m)
	ctx := context.Background()
	c := Cerere{An: 2026, Luna: 10, Limba: "ro"}

	// Tabel: tabelul urmeaza coloanele SQL la generare.
	tabel := modelFlexTest()
	tabel.Continut.Sursa = "SQL_LIBER"
	tabel.Continut.Flexibil.SQL = "SELECT angajat, marca FROM angajati"
	tabel.Antet.Imagini = map[string]Imagine{"stanga": imagineFictiva(t, "png")}
	conv, _ := salveazaSiConverteste(t, m, tabel)
	if conv.Continut.Flexibil.ModAspect != ModPdfme {
		t.Fatal("tabel neconvertit")
	}
	doc, err := m.s.Pregateste(ctx, Cerere{ModelID: conv.ID, An: 2026, Luna: 10, Limba: "ro"}, momentFix)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := DatePdfme(doc); err != nil {
		t.Fatalf("tabel: %v", err)
	}
	// Imaginea si antetul se repeta pe fiecare pagina (F09-T26): sunt elemente fixe.
	var sab struct {
		BasePdf struct {
			Padding      []float64        `json:"padding"`
			StaticSchema []map[string]any `json:"staticSchema"`
		} `json:"basePdf"`
		Schemas [][]map[string]any `json:"schemas"`
	}
	if err := json.Unmarshal([]byte(conv.Continut.Flexibil.Sablon), &sab); err != nil {
		t.Fatal(err)
	}
	imagineFixa, antetFix := false, false
	for _, e := range sab.BasePdf.StaticSchema {
		imagineFixa = imagineFixa || (e["type"] == "image" && strings.HasPrefix(e["content"].(string), "data:image/png;base64,"))
		antetFix = antetFix || e["content"] == "{antet_stanga}"
	}
	if !imagineFixa || !antetFix || len(sab.Schemas[0]) != 1 {
		t.Fatalf("imagine/antet nefixe: imagine=%v antet=%v pagina1=%d", imagineFixa, antetFix, len(sab.Schemas[0]))
	}
	// Zona de continut incepe sub imagine si antet.
	if tabelY := sab.Schemas[0][0]["position"].(map[string]any)["y"].(float64); sab.BasePdf.Padding[0] > tabelY+0.001 || sab.BasePdf.Padding[0] < 32 {
		t.Fatalf("zona de continut: padding %v, tabel la %v", sab.BasePdf.Padding[0], tabelY)
	}

	// Markdown: titlu, paragraf cu variabila, tabel.
	md := modelFlexTest()
	md.Continut.Sursa = "SQL_LIBER"
	md.Continut.Flexibil.SQL = "SELECT angajat, marca FROM angajati"
	md.Continut.Flexibil.ModAspect = "MARKDOWN"
	md.Continut.Flexibil.Markdown = "# {{titlu}}\n\nPerioada **{{perioada}}**\n- primul\n- al doilea\n\n{{tabel}}"
	conv, _ = salveazaSiConverteste(t, m, md)
	el := elementeSablon(t, conv)
	texte := []string{}
	areTabel := false
	for _, e := range el {
		if e["type"] == "table" {
			areTabel = true
		} else if e["type"] == "text" {
			texte = append(texte, e["content"].(string))
		}
	}
	if !areTabel || !strings.Contains(strings.Join(texte, "|"), "{titlu}") || !strings.Contains(strings.Join(texte, "|"), "Perioada **{perioada}**\n• primul\n• al doilea") {
		t.Fatalf("markdown convertit: %q tabel=%v", texte, areTabel)
	}
	doc, err = m.s.PrevizualizeazaFlexibil(ctx, conv, c, momentFix)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := DatePdfme(doc); err != nil {
		t.Fatalf("markdown: %v", err)
	}

	// Pozitionare libera: element fix si element repetat -> text + tabel, cu nota.
	adevarat := true
	liber := modelFlexTest()
	liber.Continut.Sursa = "SQL_LIBER"
	liber.Continut.Flexibil.SQL = "SELECT angajat, marca FROM angajati"
	liber.Continut.Flexibil.ModAspect = "LIBER"
	liber.Continut.Flexibil.InaltimeRand = 10
	liber.Continut.Flexibil.Elemente = []ElementRaport{
		{ID: "1", Text: "{{titlu}}", X: 5, Y: 5, Latime: 80, Inaltime: 12, Font: 14, Stil: &StilText{Aldin: &adevarat, Aliniere: "C"}},
		{ID: "2", Text: "{{angajat}} - {{marca}}", X: 5, Y: 30, Latime: 120, Inaltime: 8, Font: 9, RepetaRand: true},
	}
	conv, note := salveazaSiConverteste(t, m, liber)
	el = elementeSablon(t, conv)
	gasitTitlu, gasitTabel := false, false
	for _, e := range el {
		if e["type"] == "text" && e["content"] == "{titlu}" && e["fontName"] == "GoBold" && e["alignment"] == "center" {
			gasitTitlu = true
		}
		gasitTabel = gasitTabel || e["type"] == "table"
	}
	if !gasitTitlu || !gasitTabel || len(note) != 1 {
		t.Fatalf("liber convertit: titlu=%v tabel=%v note=%v", gasitTitlu, gasitTabel, note)
	}
	doc, err = m.s.PrevizualizeazaFlexibil(ctx, conv, c, momentFix)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := DatePdfme(doc); err != nil {
		t.Fatalf("liber: %v", err)
	}
}
