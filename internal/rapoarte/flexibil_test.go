package rapoarte

import (
	"context"
	"errors"
	"fmt"
	"pontix/internal/datebaza"
	"pontix/internal/i18n"
	"strings"
	"testing"
	"time"
)

func modelFlexTest() Definitie {
	return Definitie{Denumire: i18n.Text{"ro": "Raport flexibil demo", "en": "Custom demo report"}, Continut: Continut{Sursa: "CALCULATE_LUNAR", Coloane: []string{"rezultat"}, Text: i18n.Text{}, Flexibil: &RaportFlexibil{ModAspect: "TABEL", Weekend: true, Sarbatori: true}}, Pagina: Pagina{Format: "A4", Orientare: "PEISAJ", MargineMM: 10, FontPT: 8}, Activ: true}
}
func TestFlexibilSQLFormulePersistentaPDF(t *testing.T) {
	m := pregateste(t)
	fixture(t, m)
	ctx := context.Background()
	d := modelFlexTest()
	d.Continut.Flexibil.SQL = "SELECT angajat_id, angajat, marca, total_ore_cf_contract, total_ore_lucratoare, concediu_CO FROM lunar WHERE angajat_id=:angajat_id"
	d.Continut.Flexibil.Calculate = []ColoanaCalculata{{Cod: "diferenta", Denumire: i18n.Text{"ro": "Diferenta"}, Unitate: "ORE", Mod: "FORMULA", Expresie: "total_ore_cf_contract-total_ore_lucratoare", Zecimale: 2, Activ: true}, {Cod: "zile", Unitate: "ZILE", Mod: "SQL", Expresie: "SELECT count(*) FROM absente WHERE angajat_id=:angajat_id AND activ=1 AND stare_aprobare='APROBATA'", Activ: true}}
	c := Cerere{An: 2026, Luna: 10, Limba: "ro", AngajatID: "a"}
	before := 0
	m.db.QueryRow("SELECT count(*) FROM jurnal_modificari").Scan(&before)
	doc, e := m.s.PrevizualizeazaFlexibil(ctx, d, c, momentFix)
	if e != nil {
		t.Fatal(e)
	}
	if len(doc.Randuri) != 1 || len(doc.Schema) < 7 || len(doc.Coloane) != 8 {
		t.Fatalf("date %+v", doc)
	}
	after := 0
	m.db.QueryRow("SELECT count(*) FROM jurnal_modificari").Scan(&after)
	if before != after {
		t.Fatal("previzualizarea a scris")
	}
	for _, mod := range []string{"TABEL", "MARKDOWN", "LIBER"} {
		d.Continut.Flexibil.ModAspect = mod
		d.Continut.Flexibil.Markdown = "# {{titlu}}\n\nPerioada {{perioada}}\n\n{{tabel}}"
		if mod == "LIBER" {
			d.Continut.Flexibil.Elemente = []ElementRaport{{ID: "titlu", Text: "{{titlu}}", X: 0, Y: 0, Latime: 100, Inaltime: 12, Font: 12}, {ID: "tabel", Text: "{{tabel}}", X: 0, Y: 20, Latime: 270, Inaltime: 140, Font: 8}}
		}
		doc, e = m.s.PrevizualizeazaFlexibil(ctx, d, c, momentFix)
		if e != nil {
			t.Fatal(e)
		}
	}
	id, e := m.s.Salveaza(ctx, d)
	if e != nil {
		t.Fatal(e)
	}
	c.ModelID = id
	doc, e = m.s.Pregateste(ctx, c, momentFix)
	if e != nil || len(doc.Randuri) != 1 {
		t.Fatal(doc, e)
	}
	x, e := GenereazaExcel(doc)
	if e != nil || len(x) == 0 {
		t.Fatal(e)
	}
}
func TestSQLRaportIzolatRefuzaScrieriSiSecrete(t *testing.T) {
	b, e := nouaBazaRaport()
	if e != nil {
		t.Fatal(e)
	}
	defer b.db.Close()
	ctx := context.Background()
	if e := b.adauga(ctx, "angajati", []string{"angajat_id", "ore"}, [][]any{{"demo", 7}}); e != nil {
		t.Fatal(e)
	}
	for _, sql := range []string{"DELETE FROM angajati", "SELECT * FROM conturi", "SELECT * FROM sqlite_master", "ATTACH DATABASE 'x' AS a", "PRAGMA user_version", "SELECT load_extension('x') AS x", "SELECT readfile('x') AS x", "SELECT 1 AS x; DELETE FROM angajati", "SELECT :necunoscut AS x", "SELECT randomblob(99999999) AS x"} {
		if _, e := b.query(ctx, sql, nil); e == nil {
			t.Fatal("acceptat", sql)
		}
	}
	r, e := b.query(ctx, "WITH a AS (SELECT ore FROM angajati) SELECT SUM(ore) AS total FROM a", nil)
	if e != nil || fmt.Sprint(r.Randuri) != "[[7]]" {
		t.Fatal(r, e)
	}
}
func TestFlexibilConcediiCategorieSiTipDeduplicate(t *testing.T) {
	m := pregateste(t)
	fixture(t, m)
	ctx := context.Background()
	s := datebaza.Nou(m.db)
	if e := s.SalveazaAbsente(ctx, []datebaza.ZiAbsenta{{AngajatID: "a", Data: "2026-10-05", Cod: "CO", Activ: true, StareAprobare: "APROBATA"}}); e != nil {
		t.Fatal(e)
	}
	d := modelFlexTest()
	d.Continut.Flexibil.ConcediiTipuri = []string{"CO"}
	d.Continut.Flexibil.ConcediiCategorii = []string{"ODIHNA"}
	d.Continut.Flexibil.SQL = "SELECT total_zile_cf_contract, total_ore_cf_contract,total_zile_concediu,total_ore_lucratoare,concediu_CO FROM lunar WHERE angajat_id='a'"
	doc, e := m.s.PrevizualizeazaFlexibil(ctx, d, Cerere{An: 2026, Luna: 10, Limba: "ro"}, momentFix)
	if e != nil {
		t.Fatal(e)
	}
	if v := strings.Join(doc.Randuri[0].Valori, ","); v != "22,154,1,147,1" {
		t.Fatal(v)
	}
	// Un tip creat ulterior intra automat in schema, inclusiv daca nu are zile.
	if e := s.SalveazaTipAbsenta(ctx, "", datebaza.TipAbsenta{Cod: "EV-TEST", Denumire: i18n.Text{"ro": "Eveniment"}, Categorie: "SPECIAL", EfectOre: "FARA_ORE", Activ: true}); e != nil {
		t.Fatal(e)
	}
	doc, e = m.s.PrevizualizeazaFlexibil(ctx, d, Cerere{An: 2026, Luna: 10, Limba: "ro"}, momentFix)
	if e != nil {
		t.Fatal(e)
	}
	gasit := false
	for _, tab := range doc.Schema {
		for _, col := range tab.Campuri {
			gasit = gasit || col == "concediu_EV_2dTEST"
		}
	}
	if !gasit {
		t.Fatal("coloana automata lipsa")
	}
}

func TestFlexibilGarziEfectiveAuPrioritateSiDST(t *testing.T) {
	m := pregateste(t)
	dateGarziDetaliat(t, m)
	ctx := context.Background()
	for _, s := range []string{
		`INSERT INTO planificare(id,angajat_id,data,schimb_id,sectiune,creat_la,modificat_la) SELECT 'dst','g','2026-10-24',id,'GARZI','x','x' FROM tipuri_schimb WHERE cod='G24'`,
		`INSERT INTO pontaj(id,angajat_id,data,start,sfarsit,creat_la,modificat_la) VALUES('actual','g','2026-10-24','2026-10-24T04:00:00Z','2026-10-25T05:00:00Z','x','x')`,
	} {
		if _, e := m.db.Exec(s); e != nil {
			t.Fatal(e)
		}
	}
	d := modelFlexTest()
	d.Continut.Sursa = "LISTA_GARZI"
	d.Continut.Flexibil.SQL = "SELECT ore,ore_eligibile FROM garzi WHERE data='2026-10-24'"
	doc, e := m.s.PrevizualizeazaFlexibil(ctx, d, Cerere{An: 2026, Luna: 10, Limba: "ro", AngajatID: "g"}, momentFix)
	if e != nil {
		t.Fatal(e)
	}
	if len(doc.Randuri) != 1 || strings.Join(doc.Randuri[0].Valori, ",") != "25,25" {
		t.Fatalf("efectiv sau dublare: %+v", doc.Randuri)
	}
	d.Continut.Flexibil.Weekend = false
	d.Continut.Flexibil.Sarbatori = false
	doc, e = m.s.PrevizualizeazaFlexibil(ctx, d, Cerere{An: 2026, Luna: 10, Limba: "ro", AngajatID: "g"}, momentFix)
	if e != nil || doc.Randuri[0].Valori[1] != "0" {
		t.Fatal(doc, e)
	}
}
func TestFlexibilScalareMultipleSiCicluriRefuzate(t *testing.T) {
	m := pregateste(t)
	fixture(t, m)
	d := modelFlexTest()
	d.Continut.Flexibil.SQL = "SELECT 'a' AS angajat_id, 1 AS ore"
	for _, c := range []ColoanaCalculata{{Cod: "x", Unitate: "ORE", Mod: "SQL", Expresie: "SELECT 1 UNION ALL SELECT 2", Activ: true}, {Cod: "x", Unitate: "ORE", Mod: "FORMULA", Expresie: "x+1", Activ: true}} {
		d.Continut.Flexibil.Calculate = []ColoanaCalculata{c}
		if _, e := m.s.PrevizualizeazaFlexibil(context.Background(), d, Cerere{An: 2026, Luna: 10, Limba: "ro"}, momentFix); e == nil {
			t.Fatal("calcul acceptat", c)
		}
	}
}

func TestSQLRaportLimiteResurse(t *testing.T) {
	b, e := nouaBazaRaport()
	if e != nil {
		t.Fatal(e)
	}
	defer b.db.Close()
	if e := b.db.Ping(); e != nil {
		t.Fatal(e)
	}
	ctx, stop := context.WithTimeout(context.Background(), 30*time.Millisecond)
	defer stop()
	if _, e := b.query(ctx, "WITH RECURSIVE n(x) AS (SELECT 1 UNION ALL SELECT x+1 FROM n) SELECT sum(x) AS total FROM n", nil); e == nil {
		t.Fatal("interogarea fara sfarsit nu a fost oprita")
	}
	if _, e := b.query(context.Background(), "SELECT printf('%999999s', 'x') AS text", nil); e == nil {
		t.Fatal("limita de text nu a fost aplicata")
	}
}

func TestFlexibilCampuriRepetateSiAscunse(t *testing.T) {
	m := pregateste(t)
	fixture(t, m)
	d := modelFlexTest()
	d.Continut.Flexibil.SQL = "WITH RECURSIVE n(x) AS (SELECT 1 UNION ALL SELECT x+1 FROM n WHERE x<90) SELECT x AS identificator, 'Angajat demonstrativ' AS angajat FROM n"
	d.Continut.Flexibil.Ascunse = []string{"identificator"}
	d.Continut.Flexibil.ModAspect = "LIBER"
	d.Continut.Flexibil.InaltimeRand = 12
	d.Continut.Flexibil.Elemente = []ElementRaport{{ID: "nume", Text: "{{identificator}} - {{angajat}}", X: 5, Y: 20, Latime: 150, Inaltime: 10, Font: 10, RepetaRand: true}}
	doc, e := m.s.PrevizualizeazaFlexibil(context.Background(), d, Cerere{An: 2026, Luna: 10, Limba: "ro"}, momentFix)
	if e != nil {
		t.Fatal(e)
	}
	if len(doc.Coloane) != 1 || len(doc.valoriRanduri) != 90 {
		t.Fatal("ascundere sau repetare", doc.Coloane)
	}
}

// I126: eroarea numeste coloana SQL care trebuie corectata, nu doar regula generala.
func TestFlexibilEroriColoaneNumite(t *testing.T) {
	m := pregateste(t)
	fixture(t, m)
	cazuri := []struct {
		sql, cod, coloana string
		ascunse           []string
		calculate         []ColoanaCalculata
	}{
		{sql: "SELECT 1 AS zi", ascunse: []string{"rand"}, cod: "rapoarte.sql_coloana_ascunsa", coloana: "rand"},
		{sql: "SELECT 1 AS zi, 2 AS ore", calculate: []ColoanaCalculata{{Cod: "ORE", Unitate: "ORE", Mod: "FORMULA", Expresie: "zi+1", Activ: true}}, cod: "rapoarte.sql_coloana_calculata", coloana: "ORE"},
		{sql: "SELECT 'x' AS \"NUME ȘI PRENUME\"", cod: "rapoarte.sql_alias_invalid", coloana: "NUME ȘI PRENUME"},
		{sql: "SELECT 1 AS zi, 2 AS ZI", cod: "rapoarte.sql_alias_dublat", coloana: "ZI"},
		{sql: "SELECT 1 AS zi", ascunse: []string{"zi"}, cod: "rapoarte.sql_toate_ascunse"},
	}
	for _, c := range cazuri {
		d := modelFlexTest()
		d.Continut.Flexibil.SQL = c.sql
		d.Continut.Flexibil.Ascunse = c.ascunse
		d.Continut.Flexibil.Calculate = c.calculate
		_, e := m.s.PrevizualizeazaFlexibil(context.Background(), d, Cerere{An: 2026, Luna: 10, Limba: "ro"}, momentFix)
		var er *i18n.Eroare
		if !errors.As(e, &er) || er.Cod != c.cod || er.Param["coloana"] != c.coloana {
			t.Fatalf("%s: eroare %#v, astept %s «%s»", c.sql, e, c.cod, c.coloana)
		}
		if c.coloana != "" && !strings.Contains(er.Text("ro"), "«"+c.coloana+"»") {
			t.Fatalf("textul nu numeste coloana: %s", er.Text("ro"))
		}
	}
}
