package datebaza

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"reflect"
	"testing"
)

func instantaneuSimulare(t *testing.T, m *mediuTest) string {
	t.Helper()
	r := map[string][]string{}
	for _, tab := range []string{"pontaj", "pontaj_ciorne", "planificare", "absente_zile", "jurnal_modificari", "absente_documente"} {
		rows, e := m.db.Query("SELECT * FROM " + tab + " ORDER BY 1")
		if e != nil {
			t.Fatal(e)
		}
		cols, _ := rows.Columns()
		for rows.Next() {
			v := make([]any, len(cols))
			p := make([]any, len(cols))
			for i := range v {
				p[i] = &v[i]
			}
			if e := rows.Scan(p...); e != nil {
				t.Fatal(e)
			}
			b, _ := json.Marshal(v)
			r[tab] = append(r[tab], string(b))
		}
		if e := rows.Err(); e != nil {
			t.Fatal(e)
		}
		rows.Close()
	}
	b, _ := json.Marshal(r)
	return string(b)
}
func TestPrevizualizareMedicalFaraPersistenta(t *testing.T) {
	m := pregateste(t)
	ctx := context.Background()
	a := persoanaAbsente(t, m)
	g := angajatTest("SIM-GARDA")
	if e := m.serviciu.SalveazaAngajat(ctx, g); e != nil {
		t.Fatal(e)
	}
	var gid string
	for _, p := range citeste(t, m.serviciu).Angajati {
		if p.Marca == "SIM-GARDA" {
			gid = p.ID
		}
	}
	if _, e := m.db.Exec("INSERT INTO angajati_garzi VALUES(?,?)", gid, a.ID); e != nil {
		t.Fatal(e)
	}
	for _, id := range []string{a.ID, gid} {
		if _, e := m.db.Exec(`INSERT INTO pontaj(id,angajat_id,data,start,sfarsit,creat_la,modificat_la) VALUES(?,?,'2026-10-04','2026-10-04T19:00:00Z','2026-10-05T05:00:00Z','x','x')`, "sim-"+id, id); e != nil {
			t.Fatal(e)
		}
		if _, e := m.db.Exec(`INSERT INTO pontaj_ciorne(angajat_id,data,valoare,modificat_la) VALUES(?,'2026-10-04','N12','x')`, id); e != nil {
			t.Fatal(e)
		}
		if _, e := m.db.Exec(`INSERT INTO planificare(id,angajat_id,data,schimb_id,sectiune,creat_la,modificat_la) SELECT ?,?,'2026-10-04',id,'GARZI','x','x' FROM tipuri_schimb WHERE cod='N12'`, id, id); e != nil {
			t.Fatal(e)
		}
	}
	if e := m.serviciu.SalveazaAbsente(ctx, []ZiAbsenta{{AngajatID: a.ID, Data: "2026-10-05", Cod: "CO", Activ: true}}); e != nil {
		t.Fatal(e)
	}
	inainte := instantaneuSimulare(t, m)
	zile := []ZiAbsenta{{AngajatID: a.ID, Data: "2026-10-05", Cod: "CM", Activ: true, DocumentNume: "test.pdf", DocumentMime: "application/pdf", DocumentContinutB64: base64.StdEncoding.EncodeToString([]byte("%PDF-1.4\n1 0 obj<</Type/Catalog>>endobj\n%%EOF"))}}
	copie := append([]ZiAbsenta(nil), zile...)
	ef, e := m.serviciu.PrevizualizeazaAbsente(ctx, zile, 2026, 10, nil)
	if e != nil {
		t.Fatal(e)
	}
	if !reflect.DeepEqual(zile, copie) {
		t.Fatal("cererea a fost modificata")
	}
	if instantaneuSimulare(t, m) != inainte {
		t.Fatal("simularea a persistat date")
	}
	surse := map[string]int{}
	for _, v := range ef {
		surse[v.Sursa]++
	}
	if surse["pontaj"] != 1 || surse["planificare"] != 2 || surse["pontaj_ciorne"] != 2 || surse["absente_zile"] != 1 {
		t.Fatalf("efecte incomplete: %+v", surse)
	}
	if e := m.serviciu.SalveazaAbsente(ctx, zile); e != nil {
		t.Fatal(e)
	}
	if instantaneuSimulare(t, m) == inainte {
		t.Fatal("salvarea nu a persistat")
	}
}
func TestPrevizualizareAbsentaBlocataSiConversie(t *testing.T) {
	m := pregateste(t)
	ctx := context.Background()
	a := persoanaAbsente(t, m)
	if e := m.serviciu.SalveazaSchita(ctx, 2026, 10, []CelulaSchita{{AngajatID: a.ID, Data: "2026-10-01", Valoare: "Z7"}}); e != nil {
		t.Fatal(e)
	}
	s, e := m.serviciu.CitesteSchita(ctx, 2026, 10)
	if e != nil {
		t.Fatal(e)
	}
	inainte := instantaneuSimulare(t, m)
	z := []ZiAbsenta{{AngajatID: a.ID, Data: "2026-10-01", Cod: "CO", Activ: true}}
	ef, e := m.serviciu.PrevizualizeazaAbsente(ctx, z, 2026, 10, s.Celule)
	if e != nil || len(ef) < 2 {
		t.Fatalf("conversie: %v %+v", e, ef)
	}
	if instantaneuSimulare(t, m) != inainte {
		t.Fatal("conversie persistata")
	}
	z[0].Data = "2026-10-03"
	ef, e = m.serviciu.PrevizualizeazaAbsente(ctx, z, 2026, 10, nil)
	verificaCod(t, e, "absente.zi_nelucratoare")
	if ef != nil || instantaneuSimulare(t, m) != inainte {
		t.Fatal("eroare cu efecte")
	}
}
func TestCampuriEfectFaraDocumente(t *testing.T) {
	r := campuriEfect(map[string]any{"angajat_id": "a", "cod": "CO", "document_continut_b64": "secret", "nr_inregistrare": "privat", "document_nume": "privat.pdf"})
	if len(r) != 2 || r["cod"] != "CO" {
		t.Fatal(r)
	}
}

func TestPrevizualizareRefuzaLunaInchisaFaraPersistenta(t *testing.T) {
	m := pregateste(t)
	ctx := context.Background()
	a := persoanaAbsente(t, m)
	if _, err := m.serviciu.SeteazaStareLuna(ctx, 2026, 10, "INCHISA", "Test inchidere", "test"); err != nil {
		t.Fatal(err)
	}
	inainte := instantaneuSimulare(t, m)
	ef, err := m.serviciu.PrevizualizeazaAbsente(ctx, []ZiAbsenta{{AngajatID: a.ID, Data: "2026-10-01", Cod: "CO", Activ: true}}, 2026, 10, nil)
	verificaCod(t, err, "pontaj.luna_inchisa")
	if ef != nil || instantaneuSimulare(t, m) != inainte {
		t.Fatal("simularea lunii inchise a modificat date")
	}
}
