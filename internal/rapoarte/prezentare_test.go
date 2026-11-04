package rapoarte

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/pdfcpu/pdfcpu/pkg/api"
	"github.com/xuri/excelize/v2"
	"io"
	"math"
	"pontix/internal/i18n"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"testing"
)

func logic(v bool) *bool { return &v }
func TestEditareCopieAtomicaCuIstoricSiExport(t *testing.T) {
	m := pregateste(t)
	fixture(t, m)
	ctx := context.Background()
	for _, limba := range []string{"ro", "en"} {
		d, e := m.s.Pregateste(ctx, Cerere{ModelID: "raport-pontaj", An: 2026, Luna: 10, Limba: limba}, momentFix)
		if e != nil {
			t.Fatal(e)
		}
		k := "celula:0:" + d.Coloane[0].Cod
		veche := d.Randuri[0].Valori[0]
		d.Stiluri = map[string]StilText{k: {Aldin: logic(false)}}
		modificari := []EditareTextRaport{{Cheie: k, Text: "=2+2", Stil: &StilText{Aldin: logic(true), Cursiv: logic(true), FontPT: 14, Aliniere: "R"}}, {Cheie: "antet_centru", Text: "Unitate demonstrativa"}, {Cheie: "eticheta:" + d.Coloane[0].Cod, Text: "Eticheta proprie"}}
		n, e := EditeazaDocument(d, modificari, "Operator demo", momentFix)
		if e != nil {
			t.Fatal(e)
		}
		if !n.Personalizat || !strings.Contains(n.Metadate, "Operator demo") || n.Randuri[0].Valori[0] != "=2+2" || d.Randuri[0].Valori[0] != veche || d.Antet[1] == n.Antet[1] || *d.Stiluri[k].Aldin {
			t.Fatal("copie/istoric incorect")
		}
		if _, e := GenereazaPDF(n); e != nil {
			t.Fatal(e)
		}
		x, e := GenereazaExcel(n)
		if e != nil {
			t.Fatal(e)
		}
		f, e := excelize.OpenReader(bytes.NewReader(x))
		if e != nil {
			t.Fatal(e)
		}
		foaie := f.GetSheetName(0)
		v, _ := f.GetCellValue(foaie, "A5")
		formula, _ := f.GetCellFormula(foaie, "A5")
		id, _ := f.GetCellStyle(foaie, "A5")
		s, _ := f.GetStyle(id)
		f.Close()
		if v != "=2+2" || formula != "" || !s.Font.Bold || !s.Font.Italic || s.Font.Size != 14 || s.Alignment.Horizontal != "right" {
			t.Fatal(v, formula, s)
		}
		if e := m.s.JurnalPersonalizare(ctx, "token-demo", d, n, modificari, momentFix); e != nil {
			t.Fatal(e)
		}
		var inainte, dupa string
		if e := m.db.QueryRow("SELECT valoare_veche,valoare_noua FROM jurnal_modificari WHERE tabel='rapoarte_personalizate' ORDER BY id DESC LIMIT 1").Scan(&inainte, &dupa); e != nil {
			t.Fatal(e)
		}
		var a, b struct {
			Modificari []struct {
				Text string
				Stil *StilText
			}
		}
		json.Unmarshal([]byte(inainte), &a)
		json.Unmarshal([]byte(dupa), &b)
		if a.Modificari[0].Text != veche || b.Modificari[0].Text != "=2+2" || *a.Modificari[0].Stil.Aldin || !*b.Modificari[0].Stil.Aldin {
			t.Fatal("jurnal incorect", inainte, dupa)
		}
		for _, bad := range [][]EditareTextRaport{{{Cheie: "metadate", Text: "fals"}}, {{Cheie: "celula:-1:x", Text: "x"}}, {{Cheie: k, Text: "ok"}, {Cheie: k, Text: "duplicat"}}, {{Cheie: k, Text: "x", Stil: &StilText{FontPT: math.NaN()}}}, {{Cheie: k, Text: strings.Repeat("x", 4001)}}} {
			if _, e := EditeazaDocument(d, bad, "demo", momentFix); e == nil {
				t.Fatal("patch invalid acceptat", bad)
			}
		}
	}
}
func TestPrezentarePersistentaEticheteStiluri(t *testing.T) {
	m := pregateste(t)
	fixture(t, m)
	ctx := context.Background()
	d := initiala(t, m, "raport-pontaj")
	d.Continut.Prezentare = &PrezentareRaport{Stiluri: map[string]StilText{"antet_centru": {FontPT: 14, Aldin: logic(true), Cursiv: logic(true), Aliniere: "C"}}, Etichete: map[string]i18n.Text{"angajat": {"ro": "Nume si prenume", "en": "Full name"}}}
	id, e := m.s.Salveaza(ctx, d)
	if e != nil {
		t.Fatal(e)
	}
	d = initiala(t, m, id)
	if d.Continut.Prezentare.Stiluri["antet_centru"].FontPT != 14 {
		t.Fatal("stil pierdut")
	}
	for _, limba := range []string{"ro", "en"} {
		doc, e := m.s.Pregateste(ctx, Cerere{ModelID: id, An: 2026, Luna: 10, Limba: limba}, momentFix)
		if e != nil {
			t.Fatal(e)
		}
		if doc.Coloane[0].Titlu != d.Continut.Prezentare.Etichete["angajat"].In(i18n.Limba(limba)) {
			t.Fatal(doc.Coloane)
		}
		if _, e := GenereazaPDF(doc); e != nil {
			t.Fatal(e)
		}
	}
	d.Continut.Prezentare.Stiluri["antet_centru"] = StilText{FontPT: 5}
	if _, e := m.s.Salveaza(ctx, d); e == nil {
		t.Fatal("marime invalida acceptata")
	}
}
func TestGarziExtindereCalendarBoldSiTotalLunarNeschimbat(t *testing.T) {
	m := pregateste(t)
	dateGarziDetaliat(t, m)
	ctx := context.Background()
	for _, zi := range []string{"2026-11-28", "2026-11-29", "2026-11-30", "2026-12-01"} {
		_, e := m.db.Exec(`INSERT INTO planificare(id,angajat_id,data,schimb_id,sectiune,creat_la,modificat_la) SELECT ?, 'g', ?, id,'GARZI','x','x' FROM tipuri_schimb WHERE cod='G24'`, zi, zi)
		if e != nil {
			t.Fatal(e)
		}
	}
	d := modelFlexTest()
	f := d.Continut.Flexibil
	f.SQL = "SELECT data AS data_calendar,ore FROM garzi WHERE data>=:de_la AND data<:pana_la ORDER BY data"
	f.ColoanaData = "data_calendar"
	f.BoldWeekend = true
	f.BoldSarbatori = true
	f.Ascunse = []string{"data_calendar"}
	for _, luna := range []int{11, 12} {
		c := Cerere{An: 2026, Luna: luna, Limba: "ro", AngajatID: "g"}
		var normal, extins Document
		for _, extinde := range []bool{false, true} {
			f.ExtindeZileLibere = extinde
			x, e := m.s.PrevizualizeazaFlexibil(ctx, d, c, momentFix)
			if e != nil {
				t.Fatal(e)
			}
			if extinde {
				extins = x
			} else {
				normal = x
			}
		}
		if len(extins.Randuri) != 4 || len(extins.Coloane) != 1 {
			t.Fatal("interval/hide", luna, extins.Randuri)
		}
		for _, r := range extins.Randuri {
			if !r.Aldin {
				t.Fatal("zi libera fara bold", luna, r)
			}
		}
		if (luna == 11 && len(normal.Randuri) != 3) || (luna == 12 && len(normal.Randuri) != 1) {
			t.Fatal("interval implicit schimbat", luna, normal.Randuri)
		}
		sql := f.SQL
		ascunse := f.Ascunse
		f.SQL = "SELECT total_ore_garda,total_ore_cf_contract FROM lunar"
		f.Ascunse = nil
		f.BoldWeekend = false
		f.BoldSarbatori = false
		f.ExtindeZileLibere = false
		a, e := m.s.PrevizualizeazaFlexibil(ctx, d, c, momentFix)
		if e != nil {
			t.Fatal(e)
		}
		f.ExtindeZileLibere = true
		b, e := m.s.PrevizualizeazaFlexibil(ctx, d, c, momentFix)
		if e != nil || !reflect.DeepEqual(a.Randuri, b.Randuri) {
			t.Fatal("totalul lunar schimbat", a.Randuri, b.Randuri, e)
		}
		f.SQL = sql
		f.Ascunse = ascunse
		f.BoldWeekend = true
		f.BoldSarbatori = true
	}
	f.ColoanaData = "ore"
	if _, e := m.s.PrevizualizeazaFlexibil(ctx, d, Cerere{An: 2026, Luna: 11, Limba: "ro"}, momentFix); e == nil {
		t.Fatal("data invalida acceptata")
	}
}
func TestBoldRandPDFExcelPaginare(t *testing.T) {
	d := Document{Titlu: "Demo", Pagina: Pagina{Format: "A4", Orientare: "PORTRET", MargineMM: 10, FontPT: 8}, Flexibil: &RaportFlexibil{ModAspect: "MARKDOWN", Markdown: "{{tabel}}"}, Coloane: []Coloana{{Cod: "x", Titlu: "Coloana"}}, Metadate: "Istoric", Limba: i18n.Romana}
	for i := 0; i < 130; i++ {
		d.Randuri = append(d.Randuri, Rand{Valori: []string{fmt.Sprintf("RAND_%03d", i)}, Aldin: i%2 == 1})
	}
	pdf, e := GenereazaPDF(d)
	if e != nil {
		t.Fatal(e)
	}
	fonturi := map[string]string{}
	re := regexp.MustCompile(`(?s)/(F\S+) [-0-9.]+ Tf|\((.*?)\)\s*Tj`)
	e = api.ExtractContent(context.Background(), bytes.NewReader(pdf), nil, func(r io.Reader, n int) error {
		b, e := io.ReadAll(r)
		if e != nil {
			return e
		}
		font := ""
		for _, v := range re.FindAllStringSubmatch(string(b), -1) {
			if v[1] != "" {
				font = v[1]
			} else {
				var s strings.Builder
				for i := 1; i < len(v[2]); i += 2 {
					s.WriteByte(v[2][i])
				}
				if strings.HasPrefix(s.String(), "RAND_") {
					fonturi[s.String()] = font
				}
			}
		}
		return nil
	}, nil)
	if e != nil {
		t.Fatal(e)
	}
	if len(fonturi) != 130 || fonturi["RAND_000"] == fonturi["RAND_001"] {
		t.Fatal("bold PDF pierdut", fonturi)
	}
	x, e := GenereazaExcel(d)
	if e != nil {
		t.Fatal(e)
	}
	f, e := excelize.OpenReader(bytes.NewReader(x))
	if e != nil {
		t.Fatal(e)
	}
	defer f.Close()
	for i := 0; i < 130; i++ {
		id, _ := f.GetCellStyle(f.GetSheetName(0), "A"+strconv.Itoa(i+5))
		s, _ := f.GetStyle(id)
		if s.Font.Bold != (i%2 == 1) {
			t.Fatal("bold Excel pierdut", i)
		}
	}
}
