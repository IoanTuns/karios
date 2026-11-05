package rapoarte

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/pdfcpu/pdfcpu/pkg/api"
	"os"
	"path/filepath"
	"pontix/internal/calcul"
	"pontix/internal/datebaza"
	"pontix/internal/versiune"
	"strings"
	"testing"
)

func dateGarziDetaliat(t *testing.T, m *mediu) {
	t.Helper()
	fixture(t, m)
	for _, sql := range []string{
		`INSERT INTO angajati(id,nume,prenume,marca,sectie_id,departament_id,creat_la,modificat_la) VALUES('g','Alt nume','Ana','GARDA-R01','s','d','x','x'),('omonim','Serban','Ana','OMONIM','s','d','x','x')`,
		`INSERT INTO angajati_garzi VALUES('g','a')`,
		`INSERT INTO pontaj_ciorne(angajat_id,data,valoare,modificat_la) VALUES('a','2026-10-05','Z7','x')`,
		`INSERT INTO planificare(id,angajat_id,data,schimb_id,sectiune,activ,creat_la,modificat_la) SELECT 'g1','g','2026-10-05',id,'GARZI',1,'x','x' FROM tipuri_schimb WHERE cod='G24'`,
		`INSERT INTO planificare(id,angajat_id,data,schimb_id,sectiune,activ,creat_la,modificat_la) SELECT 'g0','g','2026-09-30',id,'GARZI',1,'x','x' FROM tipuri_schimb WHERE cod='G24'`,
		`INSERT INTO planificare(id,angajat_id,data,schimb_id,sectiune,activ,creat_la,modificat_la) SELECT 'g2','g','2026-10-10',id,'GARZI',0,'x','x' FROM tipuri_schimb WHERE cod='G24'`,
		`INSERT INTO planificare(id,angajat_id,data,schimb_id,sectiune,activ,creat_la,modificat_la) SELECT 'g3','g','2026-10-11',id,'DUPA_AMIAZA',1,'x','x' FROM tipuri_schimb WHERE cod='S7'`,
	} {
		if _, err := m.db.Exec(sql); err != nil {
			t.Fatal(err)
		}
	}
}
func valoriRand(doc Document, i int) map[string]string {
	r := map[string]string{}
	for j, c := range doc.Coloane {
		r[c.Cod] = doc.Randuri[i].Valori[j]
	}
	return r
}
func TestDetaliatGarziMarciDistincteFaraDublare(t *testing.T) {
	m := pregateste(t)
	dateGarziDetaliat(t, m)
	ctx := context.Background()
	cerere := Cerere{ModelID: "raport-pontaj-detaliat", An: 2026, Luna: 10, Limba: "ro", AngajatID: "a"}
	citeste := func() Document {
		d, e := m.s.Pregateste(ctx, cerere, momentFix)
		if e != nil {
			t.Fatal(e)
		}
		return d
	}
	doc := citeste()
	if len(doc.Randuri) != 2 {
		t.Fatal("familie cu doua marci", len(doc.Randuri))
	}
	normal, garda := valoriRand(doc, 0), valoriRand(doc, 1)
	if normal["marca"] != "DEMO-R01" || garda["marca"] != "GARDA-R01" || normal["2026-10-05"] != "7" || garda["2026-10-05"] != "24" || garda["total_T"] != "24" {
		t.Fatal(normal, garda)
	}
	if garda["2026-10-10"] != "" || garda["2026-10-11"] != "" || garda["total_ore_garda"] != "—" {
		t.Fatal("surse nepermise sau procent inventat", garda)
	}
	if _, err := m.db.Exec("UPDATE angajati SET activ=0 WHERE id='g'"); err != nil {
		t.Fatal(err)
	}
	if d := citeste(); len(d.Randuri) != 2 {
		t.Fatal("garda inactiva cu date omisa")
	}
	cerere.SectieID = "s"
	if _, err := m.db.Exec("UPDATE angajati SET sectie_id=NULL,departament_id=NULL WHERE id='g'"); err != nil {
		t.Fatal(err)
	}
	if d := citeste(); len(d.Randuri) != 1 {
		t.Fatal("filtru sectie ignorat")
	}
	cerere.SectieID = ""
	var nr int
	if err := m.db.QueryRow("SELECT count(*) FROM pontaj_ciorne WHERE angajat_id='g'").Scan(&nr); err != nil || nr != 0 {
		t.Fatal("raportul a scris pontaj", nr, err)
	}
	if _, err := m.db.Exec(`INSERT INTO pontaj_ciorne(angajat_id,data,valoare,modificat_la) VALUES('g','2026-10-05','G24','x')`); err != nil {
		t.Fatal(err)
	}
	if v := valoriRand(citeste(), 1); v["total_T"] != "24" {
		t.Fatal("dublare dupa publicare", v)
	}
	if _, err := m.db.Exec(`UPDATE pontaj_ciorne SET valoare='@M0' WHERE angajat_id='g'`); err != nil {
		t.Fatal(err)
	}
	if v := valoriRand(citeste(), 1); v["total_T"] != "0" {
		t.Fatal("pontaj manual ignorat", v)
	}
	if _, err := m.db.Exec(`DELETE FROM pontaj_ciorne WHERE angajat_id='g'; INSERT INTO absente_zile(angajat_id,data,cod,activ,creat_la,modificat_la) VALUES('g','2026-10-05','CO',1,'x','x')`); err != nil {
		t.Fatal(err)
	}
	if v := valoriRand(citeste(), 1); v["total_T"] != "0" || v["total_CO"] != "1" {
		t.Fatal("absenta ignorata", v)
	}
	cerere.AngajatID = "g"
	cerere.Limba = "en"
	if d := citeste(); len(d.Randuri) != 2 {
		t.Fatal("filtru dupa marca garda")
	}
	cerere.AngajatID = "omonim"
	if d := citeste(); len(d.Randuri) != 1 {
		t.Fatal("omonim asociat automat")
	}
}
func TestDetaliatGardaDinLunaPrecedenta(t *testing.T) {
	m := pregateste(t)
	dateGarziDetaliat(t, m)
	ctx := context.Background()
	s := datebaza.Nou(m.db)
	c, _, err := s.ConfiguratiePontaj(ctx)
	if err != nil {
		t.Fatal(err)
	}
	c.Atribuire = calcul.AtribuireCalendar
	if err := s.SalveazaConfiguratiePontaj(ctx, c); err != nil {
		t.Fatal(err)
	}
	d, err := m.s.Pregateste(ctx, Cerere{ModelID: "raport-pontaj-detaliat", An: 2026, Luna: 10, Limba: "ro", AngajatID: "a"}, momentFix)
	if err != nil {
		t.Fatal(err)
	}
	v := valoriRand(d, 1)
	if v["2026-10-01"] != "7" || v["2026-10-05"] != "17" || v["2026-10-06"] != "7" || v["total_T"] != "31" {
		t.Fatal("tura la granita lunii", v)
	}
}

func TestReferintaDetaliatCuMarciGarda(t *testing.T) {
	m := pregateste(t)
	dateGarziDetaliat(t, m)
	d, err := m.s.Pregateste(context.Background(), Cerere{ModelID: "raport-pontaj-detaliat", An: 2026, Luna: 10, Limba: "ro", AngajatID: "a"}, momentFix)
	if err != nil {
		t.Fatal(err)
	}
	pdf, err := GenereazaPDF(d)
	if err != nil {
		t.Fatal(err)
	}
	if err := api.Validate(context.Background(), bytes.NewReader(pdf), nil, nil); err != nil {
		t.Fatal(err)
	}
	if *directorDemo != "" {
		if err := os.MkdirAll(*directorDemo, 0700); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(*directorDemo, "pontaj-detaliat-marci-garda.pdf"), pdf, 0600); err != nil {
			t.Fatal(err)
		}
	}
	d.Metadate = strings.ReplaceAll(d.Metadate, versiune.Aplicatie, "<aplicatie>")
	date, err := json.MarshalIndent(d, "", "  ")
	if err != nil {
		t.Fatal(err)
	}
	cale := filepath.Join("testdata", "pontaj-detaliat-marci-garda-v18.json")
	if *actualizeazaReferinte {
		if err := os.WriteFile(cale, date, 0600); err != nil {
			t.Fatal(err)
		}
	}
	asteptat, err := os.ReadFile(cale)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(asteptat, date) {
		t.Fatal("referinta marcilor de garda difera")
	}
}
