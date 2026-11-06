package rapoarte

import (
	"context"
	"reflect"
	"strings"
	"testing"

	"pontix/internal/i18n"
)

// TestRapoarteStandardDinSQLEchivalente verifica D090: pentru fiecare model
// standard, raportul SQL care selecteaza aceleasi coloane din tabelul calculat
// da aceleasi titluri, valori, culori si aceeasi legenda ca motorul vechi.
func TestRapoarteStandardDinSQLEchivalente(t *testing.T) {
	m := pregateste(t)
	fixture(t, m)
	ctx := context.Background()
	for _, caz := range []struct{ model, tabel, angajat string }{
		{"raport-pontaj", "pontaj_lunar", ""},
		{"raport-pontaj-detaliat", "pontaj_detaliat", ""},
		{"raport-fisa", "fisa_individuala", "a"},
		{"raport-absente", "absente_anuale", ""},
		{"raport-garzi", "grafic_garzi", ""},
	} {
		t.Run(caz.model, func(t *testing.T) {
			c := Cerere{ModelID: caz.model, An: 2026, Luna: 10, Limba: "ro", AngajatID: caz.angajat}
			vechi, err := m.s.Pregateste(ctx, c, momentFix)
			if err != nil {
				t.Fatal(err)
			}
			coduri := []string{"rand"}
			for _, col := range vechi.Coloane {
				coduri = append(coduri, codSQLStandard(col.Cod))
			}
			d := modelFlexTest()
			d.Continut.Sursa = "SQL_LIBER"
			d.Continut.AscundeExplicatii = false
			d.Continut.Flexibil.SQL = "SELECT " + strings.Join(coduri, ", ") + " FROM " + caz.tabel + " ORDER BY rand"
			d.Continut.Flexibil.Ascunse = []string{"rand"}
			c.ModelID = ""
			nou, err := m.s.PrevizualizeazaFlexibil(ctx, d, c, momentFix)
			if err != nil {
				t.Fatal(err)
			}
			if len(nou.Coloane) != len(vechi.Coloane) || len(nou.Randuri) != len(vechi.Randuri) {
				t.Fatalf("dimensiuni: %d/%d coloane, %d/%d randuri", len(nou.Coloane), len(vechi.Coloane), len(nou.Randuri), len(vechi.Randuri))
			}
			for j := range vechi.Coloane {
				if nou.Coloane[j].Cod != codSQLStandard(vechi.Coloane[j].Cod) || nou.Coloane[j].Titlu != vechi.Coloane[j].Titlu || nou.Coloane[j].Culoare != vechi.Coloane[j].Culoare {
					t.Errorf("coloana %d: %+v fata de %+v", j, nou.Coloane[j], vechi.Coloane[j])
				}
			}
			for i := range vechi.Randuri {
				v, n := vechi.Randuri[i], nou.Randuri[i]
				if !reflect.DeepEqual(n.Valori, v.Valori) {
					t.Errorf("rand %d valori: %q fata de %q", i, n.Valori, v.Valori)
				}
				if strings.Join(n.Culori, "") != strings.Join(v.Culori, "") || (len(strings.Join(v.Culori, "")) > 0 && !reflect.DeepEqual(n.Culori, v.Culori)) {
					t.Errorf("rand %d culori: %q fata de %q", i, n.Culori, v.Culori)
				}
				if n.Grup != v.Grup || n.Aldin != v.Aldin {
					t.Errorf("rand %d grup/aldin diferite", i)
				}
			}
			// Legenda veche apare intreaga, inaintea explicatiei SQL.
			if len(nou.Legenda) < len(vechi.Legenda) || !reflect.DeepEqual(nou.Legenda[:len(vechi.Legenda)], vechi.Legenda) {
				t.Errorf("legenda: %q fata de %q", nou.Legenda, vechi.Legenda)
			}
		})
	}
}

func TestTabelStandardDoarCandEsteFolosit(t *testing.T) {
	m := pregateste(t)
	fixture(t, m)
	d := modelFlexTest()
	d.Continut.Sursa = "SQL_LIBER"
	d.Continut.Flexibil.SQL = "SELECT angajat FROM angajati"
	doc, err := m.s.PrevizualizeazaFlexibil(context.Background(), d, Cerere{An: 2026, Luna: 10, Limba: "ro"}, momentFix)
	if err != nil {
		t.Fatal(err)
	}
	for _, s := range doc.Schema {
		for _, ts := range tabeleStandard {
			if s.Cod == ts.Tabel {
				t.Fatalf("tabelul %s nu trebuia calculat", ts.Tabel)
			}
		}
	}
	// Fisa individuala cere un angajat, ca motorul vechi.
	d.Continut.Flexibil.SQL = "SELECT * FROM fisa_individuala"
	if _, err := m.s.PrevizualizeazaFlexibil(context.Background(), d, Cerere{An: 2026, Luna: 10, Limba: "ro"}, momentFix); err == nil {
		t.Fatal("fisa fara angajat acceptata")
	}
}

func TestCodSQLStandard(t *testing.T) {
	for in, out := range map[string]string{"2026-10-01": "z01", "2026-02-28": "z28", "total_L_INCEPUT": "total_L_INCEPUT", "angajat": "angajat"} {
		if codSQLStandard(in) != out {
			t.Errorf("%s -> %s", in, codSQLStandard(in))
		}
	}
	if !folosesteTabel("select * FROM Pontaj_Lunar p", "pontaj_lunar") || folosesteTabel("SELECT * FROM pontaj_lunar_x", "pontaj_lunar") {
		t.Fatal("detectarea tabelului")
	}
	_ = i18n.Romana
}

// Culorile urmeaza codul coloanei, nu pozitia: SQL-ul poate reordona si filtra.
func TestTabelStandardReordonatPastreazaCulorile(t *testing.T) {
	m := pregateste(t)
	fixture(t, m)
	ctx := context.Background()
	c := Cerere{ModelID: "raport-pontaj", An: 2026, Luna: 10, Limba: "ro"}
	vechi, err := m.s.Pregateste(ctx, c, momentFix)
	if err != nil {
		t.Fatal(err)
	}
	idx := map[string]int{}
	for j, col := range vechi.Coloane {
		idx[codSQLStandard(col.Cod)] = j
	}
	d := modelFlexTest()
	d.Continut.Sursa = "SQL_LIBER"
	d.Continut.Flexibil.SQL = "SELECT z03, angajat, z26, rand FROM pontaj_lunar WHERE rand = 0"
	d.Continut.Flexibil.Ascunse = []string{"rand"}
	c.ModelID = ""
	nou, err := m.s.PrevizualizeazaFlexibil(ctx, d, c, momentFix)
	if err != nil {
		t.Fatal(err)
	}
	if len(nou.Randuri) != 1 || len(nou.Coloane) != 3 {
		t.Fatalf("dimensiuni %d/%d", len(nou.Randuri), len(nou.Coloane))
	}
	v := vechi.Randuri[0]
	for j, cod := range []string{"z03", "angajat", "z26"} {
		if nou.Randuri[0].Valori[j] != v.Valori[idx[cod]] {
			t.Errorf("%s: valoare %q fata de %q", cod, nou.Randuri[0].Valori[j], v.Valori[idx[cod]])
		}
		if len(v.Culori) > idx[cod] && len(nou.Randuri[0].Culori) == 3 && nou.Randuri[0].Culori[j] != v.Culori[idx[cod]] {
			t.Errorf("%s: culoare %q fata de %q", cod, nou.Randuri[0].Culori[j], v.Culori[idx[cod]])
		}
	}
	if strings.Join(v.Culori, "") != "" && len(nou.Randuri[0].Culori) != 3 {
		t.Fatalf("culori pierdute: %q", nou.Randuri[0].Culori)
	}
}

func TestCompleteazaZileLunaScurta(t *testing.T) {
	cols := []string{"rand", "angajat", "z01", "z28", "total"}
	rows := [][]any{{0, "A", "8", "8", "16"}}
	c, r := completeazaZile(cols, rows)
	if strings.Join(c, ",") != "rand,angajat,z01,z28,z29,z30,z31,total" {
		t.Fatalf("coloane: %v", c)
	}
	if len(r[0]) != 8 || r[0][4] != nil || r[0][7] != "16" {
		t.Fatalf("rand: %v", r[0])
	}
	if c2, _ := completeazaZile([]string{"rand", "data"}, nil); len(c2) != 2 {
		t.Fatal("tabel fara zile modificat")
	}
}
