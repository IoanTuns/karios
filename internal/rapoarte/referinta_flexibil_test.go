package rapoarte

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
	"testing"

	"pontix/internal/i18n"
	"pontix/internal/versiune"
)

// TestReferintaRaportFlexibil (F09-T06, D009): rapoartele flexibile au si ele un fisier de
// referinta, pe versiunea lor (internal/rapoarte/flexibil). Modelul atinge tabelele principale
// din setul de date: lunar, angajati, istoric_contract, programe (cu culoare), plus o coloana
// calculata prin formula si una prin SQL.
func TestReferintaRaportFlexibil(t *testing.T) {
	m := pregateste(t)
	fixture(t, m)
	d := modelFlexTest()
	d.Continut.Flexibil.SQL = `SELECT l.angajat_id, l.angajat, l.marca, l.total_ore_cf_contract, l.total_ore_lucratoare, l.concediu_CO,
		a.norma_zilnica, a.concediu_zile_an,
		(SELECT count(*) FROM istoric_contract i WHERE i.angajat_id=l.angajat_id) AS schimbari_contract,
		(SELECT count(*) FROM programe) AS programe
		FROM lunar l JOIN angajati a ON a.angajat_id=l.angajat_id ORDER BY l.marca`
	d.Continut.Flexibil.Calculate = []ColoanaCalculata{
		{Cod: "diferenta", Denumire: i18n.Text{"ro": "Diferenta"}, Unitate: "ORE", Mod: "FORMULA", Expresie: "total_ore_cf_contract-total_ore_lucratoare", Zecimale: 2, Activ: true},
		{Cod: "zile", Unitate: "ZILE", Mod: "SQL", Expresie: "SELECT count(*) FROM absente WHERE angajat_id=:angajat_id AND activ=1", Activ: true},
	}
	doc, err := m.s.PrevizualizeazaFlexibil(context.Background(), d, Cerere{An: 2026, Luna: 10, Limba: "ro"}, momentFix)
	if err != nil {
		t.Fatal(err)
	}
	doc.Metadate = strings.ReplaceAll(doc.Metadate, versiune.Aplicatie, "<aplicatie>")
	primit, err := json.MarshalIndent(doc, "", "  ")
	if err != nil {
		t.Fatal(err)
	}
	cale := filepath.Join("testdata", "raport-flexibil-v"+strconv.Itoa(doc.VersiuneRaport)+".json")
	if *actualizeazaReferinte {
		continut, _ := json.MarshalIndent(struct {
			Document json.RawMessage `json:"document"`
		}{primit}, "", "  ")
		if err := os.WriteFile(cale, continut, 0o600); err != nil {
			t.Fatal(err)
		}
	}
	brut, err := os.ReadFile(cale)
	if err != nil {
		t.Fatal(err)
	}
	var fisier struct {
		Document json.RawMessage `json:"document"`
	}
	if err := json.Unmarshal(brut, &fisier); err != nil {
		t.Fatal(err)
	}
	var asteptat, obtinut any
	_ = json.Unmarshal(fisier.Document, &asteptat)
	_ = json.Unmarshal(primit, &obtinut)
	if !reflect.DeepEqual(asteptat, obtinut) {
		t.Fatalf("raportul flexibil difera de referinta; verifica schimbarea si versiunea raportului: %s", cale)
	}
}

// F09-T37 (D107): tabelul de date pontaj_detaliat are si coloana total_OS (ore suplimentare).
func TestPontajDetaliatSQLAreOreSuplimentare(t *testing.T) {
	m := pregateste(t)
	fixture(t, m)
	d := modelFlexTest()
	d.Continut.Flexibil.SQL = `SELECT marca, total_T, total_OS FROM pontaj_detaliat WHERE grup = 0 ORDER BY marca`
	doc, err := m.s.PrevizualizeazaFlexibil(context.Background(), d, Cerere{An: 2026, Luna: 10, Limba: "ro"}, momentFix)
	if err != nil {
		t.Fatal(err)
	}
	if len(doc.Coloane) != 3 || doc.Coloane[2].Cod != "total_OS" || len(doc.Randuri) == 0 {
		t.Fatalf("coloane %+v, randuri %d", doc.Coloane, len(doc.Randuri))
	}
}
