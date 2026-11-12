package rapoarte

import (
	"bytes"
	"context"
	"github.com/xuri/excelize/v2"
	"pontix/internal/calcul"
	"pontix/internal/i18n"
	"strconv"
	"strings"
	"testing"
)

func TestLegendaColoanelorVizibileSiRegulilor(t *testing.T) {
	c := calcul.ConfiguratiePontajInitiala()
	c.Reguli[4] = calcul.RegulaColoana{Coloana: "CEV", Sursa: "ZILE_ABSENTA", Coduri: []string{"EV"}, FactorSutimi: 100, NumaiAprobate: true}
	c.Conditionate = []calcul.RegulaConditionata{{Activ: true, GrupID: "med", Regula: calcul.RegulaColoana{Coloana: "ore_garda", Sursa: "ORE_SCHIMBURI_ZNL", Coduri: []string{"G24"}, FactorSutimi: 200}}}
	d := Document{Sursa: "PONTAJ_DETALIAT", Coloane: []Coloana{{Cod: "marca", Titlu: "Marca"}, {Cod: "2026-10-01", Titlu: "1"}, {Cod: "2026-10-02", Titlu: "2"}, {Cod: "total_CEV", Titlu: "CEV"}, {Cod: "total_ore_garda", Titlu: "Garda"}}}
	for _, limba := range []i18n.Limba{i18n.Romana, i18n.Engleza} {
		l := legendaColoane(d, c, limba, map[string]string{"grup/med": "Medici fictivi"})
		if len(l) != 5 {
			t.Fatal(l)
		}
		text := strings.Join(l, " ")
		for _, val := range []string{"EV", "G24", "x 2", "Medici fictivi"} {
			if !strings.Contains(text, val) {
				t.Fatal(val, text)
			}
		}
		if strings.Contains(text, "rapoarte.legenda_") {
			t.Fatal(text)
		}
	}
}
func TestDetaliatFaraOreNormaCuLegendaPDFExcel(t *testing.T) {
	m := pregateste(t)
	fixture(t, m)
	for _, limba := range []string{"ro", "en"} {
		d, err := m.s.Pregateste(context.Background(), Cerere{ModelID: "raport-pontaj-detaliat", An: 2026, Luna: 10, Limba: limba}, momentFix)
		if err != nil {
			t.Fatal(err)
		}
		for _, c := range d.Coloane {
			if c.Cod == "total_ore_norma" {
				t.Fatal("coloana ascunsa inca vizibila")
			}
		}
		if len(d.Legenda) == 0 {
			t.Fatal("legenda lipsa")
		}
		b, err := GenereazaExcel(d)
		if err != nil {
			t.Fatal(err)
		}
		f, err := excelize.OpenReader(bytes.NewReader(b))
		if err != nil {
			t.Fatal(err)
		}
		v, err := f.GetCellValue(f.GetSheetName(0), "A"+strconv.Itoa(len(d.Randuri)+6))
		f.Close()
		if err != nil || v != d.Legenda[0] {
			t.Fatal(v, err)
		}
	}
}

// F04-T08 (D102): randul despre schimbarea orei apare doar cu orele efective si nu la graficul de garzi.
func TestLegendaSchimbareOra(t *testing.T) {
	c := calcul.ConfiguratiePontajInitiala()
	d := Document{Sursa: "PONTAJ", Coloane: []Coloana{{Cod: "marca", Titlu: "Marca"}, {Cod: "2026-10-01", Titlu: "1"}}}
	cheie := i18n.T(i18n.Romana, "rapoarte.legenda_schimbare_ora_efectiva", nil)
	if strings.Contains(strings.Join(legendaColoane(d, c, i18n.Romana), " "), cheie) {
		t.Fatal("rand cu durata din program")
	}
	c.SchimbareOra = calcul.SchimbareOraEfectiva
	if !strings.Contains(strings.Join(legendaColoane(d, c, i18n.Romana), " "), cheie) {
		t.Fatal("lipseste randul cu orele efective")
	}
	d.Sursa = "GRAFIC_GARZI"
	if strings.Contains(strings.Join(legendaColoane(d, c, i18n.Romana), " "), cheie) {
		t.Fatal("rand la garzi")
	}
}
