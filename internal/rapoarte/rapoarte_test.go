package rapoarte

import (
	"bytes"
	"context"
	"crypto/rand"
	"database/sql"
	"encoding/json"
	"errors"
	"flag"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/pdfcpu/pdfcpu/pkg/api"

	"pontix/internal/bd"
	"pontix/internal/bd/acces"
	"pontix/internal/calcul"
	"pontix/internal/datebaza"
	"pontix/internal/i18n"
	"pontix/internal/versiune"
)

var actualizeazaReferinte = flag.Bool("actualizeaza-rapoarte", false, "Scrie referințele inițiale ale rapoartelor.")
var directorDemo = flag.String("director-demo", "", "Folder opțional pentru PDF-uri exclusiv fictive.")

type mediu struct {
	db    *sql.DB
	s     *Serviciu
	cale  string
	cheie []byte
}

func pregateste(t *testing.T) *mediu {
	t.Helper()
	m := &mediu{cale: filepath.Join(t.TempDir(), "rapoarte.db"), cheie: make([]byte, bd.LungimeCheie)}
	if _, err := rand.Read(m.cheie); err != nil {
		t.Fatal(err)
	}
	var err error
	m.db, err = bd.Deschide(m.cale, m.cheie)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { m.db.Close() })
	if _, err := bd.Migreaza(context.Background(), m.db, bd.Migrari); err != nil {
		t.Fatal(err)
	}
	m.s = Nou(m.db)
	return m
}
func codEroare(t *testing.T, err error, cod string) {
	t.Helper()
	var e *i18n.Eroare
	if !errors.As(err, &e) || e.Cod != cod {
		t.Fatalf("eroare %v, așteptat %s", err, cod)
	}
}
func initiala(t *testing.T, m *mediu, id string) Definitie {
	t.Helper()
	c, err := m.s.Lista(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	for _, d := range c.Modele {
		if d.ID == id {
			return d
		}
	}
	t.Fatalf("lipsește %s", id)
	return Definitie{}
}
func fixture(t *testing.T, m *mediu) {
	t.Helper()
	for _, comanda := range []string{
		`INSERT INTO sectii(id,denumire,creat_la,modificat_la) VALUES ('s','{"ro":"Laborator demo","en":"Demo laboratory"}','x','x')`,
		`INSERT INTO departamente(id,sectie_id,denumire,creat_la,modificat_la) VALUES ('d','s','{"ro":"Analize demo","en":"Demo analyses"}','x','x')`,
		`INSERT INTO angajati(id,nume,prenume,marca,sectie_id,departament_id,creat_la,modificat_la) VALUES ('a','Șerban','Ana','DEMO-R01','s','d','x','x')`,
		`INSERT INTO angajati(id,nume,prenume,marca,sectie_id,departament_id,activ,creat_la,modificat_la) VALUES ('b','Țincu','Răzvan','DEMO-R02','s','d',0,'x','x')`,
		`INSERT INTO angajati(id,nume,prenume,marca,activ,creat_la,modificat_la) VALUES ('c','Fictiv','Inactiv','DEMO-R03',0,'x','x')`,
		`INSERT INTO pontaj_ciorne(angajat_id,data,valoare,modificat_la) VALUES ('a','2026-10-23','IN9','x'),('b','2026-10-01','Z7','x'),('a','2028-02-29','N9','x')`,
		`INSERT INTO absente_zile(angajat_id,data,cod,activ,creat_la,modificat_la) VALUES ('a','2026-10-24','CO',1,'x','x'),('a','2026-10-25','CM',0,'x','x'),('b','2027-01-01','CM',1,'x','x')`,
		`INSERT INTO planificare(id,angajat_id,data,schimb_id,sectiune,activ,creat_la,modificat_la) VALUES ('p1','a','2026-10-23','20000000-0000-4000-8000-000000000004','GARZI',1,'x','x')`,
	} {
		if _, err := m.db.Exec(comanda); err != nil {
			t.Fatal(err)
		}
	}
}

var momentFix = time.Date(2026, 10, 25, 12, 30, 0, 0, time.UTC)

func TestDetaliatFolosesteAceleasiOreSiReguliCaGrila(t *testing.T) {
	m := pregateste(t)
	fixture(t, m)
	ctx := context.Background()
	s := datebaza.Nou(m.db)
	c, _, err := s.ConfiguratiePontaj(ctx)
	if err != nil {
		t.Fatal(err)
	}
	c.Atribuire = calcul.AtribuireCalendar
	if err = s.SalveazaConfiguratiePontaj(ctx, c); err != nil {
		t.Fatal(err)
	}
	doc, err := m.s.Pregateste(ctx, Cerere{ModelID: "raport-pontaj-detaliat", An: 2026, Luna: 10, Limba: "ro"}, momentFix)
	if err != nil {
		t.Fatal(err)
	}
	if len(doc.Coloane) != 45 {
		t.Fatal(len(doc.Coloane))
	}
	rez, _, err := datebaza.CalculeazaPontaj(ctx, acces.New(m.db), 2026, 10, nil)
	if err != nil {
		t.Fatal(err)
	}
	valori := map[string]string{}
	for i, col := range doc.Coloane {
		if strings.ContainsAny(col.Titlu, "\r\n") {
			t.Fatalf("separator de rand in titlul PDF: %q", col.Titlu)
		}
		valori[col.Cod] = doc.Randuri[0].Valori[i]
	}
	if valori["2026-10-23"] != "2" || valori["2026-10-24"] != "7" || valori["total_T"] != "9" || valori["total_ZNL"] != "7" || valori["total_ore_norma"] != "2" || valori["total_CO"] != "1" || valori["total_ore_garda"] != "—" {
		t.Fatal(valori)
	}
	for _, r := range rez.Totaluri {
		if r.AngajatID == "a" {
			for cod, v := range r.Valori {
				if valori["total_"+cod] != calcul.TextValoarePontaj(v) {
					t.Fatal(cod, v, valori)
				}
			}
		}
	}
	model := initiala(t, m, "raport-pontaj-detaliat")
	model.Continut.Coloane = []string{"angajat", "totaluri", "zile"}
	codEroare(t, Valideaza(model), "rapoarte.continut_invalid")
	var raw string
	if err = m.db.QueryRow(`SELECT valoare FROM pontaj_ciorne WHERE angajat_id='a' AND data='2026-10-23'`).Scan(&raw); err != nil || raw != "IN9" {
		t.Fatal(raw, err)
	}
}

func TestModeleInitialeReferinta(t *testing.T) {
	m := pregateste(t)
	c, err := m.s.Lista(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	for _, d := range c.Modele {
		if err := Valideaza(d); err != nil {
			t.Fatal(err)
		}
	}
	date, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		t.Fatal(err)
	}
	cale := filepath.Join("testdata", "modele-v4.json")
	if *actualizeazaReferinte {
		if err := os.MkdirAll("testdata", 0700); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(cale, date, 0600); err != nil {
			t.Fatal(err)
		}
	}
	asteptat, err := os.ReadFile(cale)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(asteptat, date) {
		t.Fatal("modelele inițiale diferă de referință")
	}
}

func TestDefinitiiIndependentePersistenteSiConflicte(t *testing.T) {
	m := pregateste(t)
	ctx := context.Background()
	d := initiala(t, m, "raport-pontaj")
	originalFisa := initiala(t, m, "raport-fisa")
	d.Antet.Stanga = i18n.Text{"ro": "Unitate demo", "en": "Demo organisation"}
	d.Subsol.Dreapta = i18n.Text{"ro": "Semnătură demo"}
	id, err := m.s.Salveaza(ctx, d)
	if err != nil || id != d.ID {
		t.Fatalf("salvare: %s %v", id, err)
	}
	_, err = m.s.Salveaza(ctx, d)
	codEroare(t, err, "rapoarte.conflict")
	actual := initiala(t, m, d.ID)
	if actual.Revizie != "2" || actual.Antet.Stanga["ro"] != "Unitate demo" {
		t.Fatal("editare nepăstrată")
	}
	document, err := m.s.Pregateste(ctx, Cerere{ModelID: d.ID, Revizie: actual.Revizie, An: 2026, Luna: 10, Limba: "en"}, momentFix)
	if err != nil || document.Antet[0] != "Demo organisation" || document.Subsol[2] != "Semnătură demo" {
		t.Fatalf("componentele personalizate nu sunt folosite: %+v %v", document, err)
	}
	neschimbata := initiala(t, m, "raport-fisa")
	v1, _ := json.Marshal(originalFisa)
	v2, _ := json.Marshal(neschimbata)
	if !bytes.Equal(v1, v2) {
		t.Fatal("alt raport a fost schimbat")
	}
	actual.ID, actual.Revizie = "", ""
	actual.Denumire = i18n.Text{"ro": "Copie independentă"}
	copieID, err := m.s.Salveaza(ctx, actual)
	if err != nil || copieID == d.ID {
		t.Fatalf("duplicare: %s %v", copieID, err)
	}
	copie := initiala(t, m, copieID)
	copie.Activ = false
	if _, err := m.s.Salveaza(ctx, copie); err != nil {
		t.Fatal(err)
	}
	if err := m.db.Close(); err != nil {
		t.Fatal(err)
	}
	m.db, err = bd.Deschide(m.cale, m.cheie)
	if err != nil {
		t.Fatal(err)
	}
	m.s = Nou(m.db)
	if initiala(t, m, copieID).Activ || initiala(t, m, d.ID).Antet.Stanga["en"] != "Demo organisation" {
		t.Fatal("redeschidere incorectă")
	}
	_, err = m.s.Pregateste(ctx, Cerere{ModelID: copieID, An: 2026, Luna: 10, Limba: "ro"}, momentFix)
	codEroare(t, err, "rapoarte.inactiv")
}

func TestValidariDefinitiiSiJurnalAtomic(t *testing.T) {
	m := pregateste(t)
	for _, schimba := range []struct {
		cod    string
		aplica func(*Definitie)
	}{
		{"rapoarte.nume_invalid", func(d *Definitie) { d.Denumire = i18n.Text{} }},
		{"rapoarte.continut_invalid", func(d *Definitie) { d.Continut.Coloane = []string{"marca"} }},
		{"rapoarte.continut_invalid", func(d *Definitie) { d.Continut.Coloane = []string{"angajat", "angajat"} }},
		{"rapoarte.pagina_invalida", func(d *Definitie) { d.Pagina.MargineMM = 100 }},
		{"rapoarte.variabila_invalida", func(d *Definitie) { d.Antet.Stanga = i18n.Text{"ro": "{{executa_sql}}"} }},
		{"rapoarte.text_invalid", func(d *Definitie) { d.Subsol.Centru = i18n.Text{"ro": strings.Repeat("x", 1001)} }},
	} {
		d := initiala(t, m, "raport-pontaj")
		schimba.aplica(&d)
		codEroare(t, Valideaza(d), schimba.cod)
	}
	if _, err := m.db.Exec(`CREATE TRIGGER refuza_jurnal_raport BEFORE INSERT ON jurnal_modificari WHEN NEW.tabel='definitii_rapoarte' BEGIN SELECT RAISE(ABORT,'proba'); END`); err != nil {
		t.Fatal(err)
	}
	d := initiala(t, m, "raport-pontaj")
	d.Antet.Stanga = i18n.Text{"ro": "Nesalvat"}
	_, err := m.s.Salveaza(context.Background(), d)
	codEroare(t, err, "rapoarte.eroare")
	if initiala(t, m, d.ID).Revizie != "1" {
		t.Fatal("salvare parțială fără jurnal")
	}
}

func TestDateRaportFiltreCalendarSiAbsente(t *testing.T) {
	m := pregateste(t)
	fixture(t, m)
	ctx := context.Background()
	c := Cerere{ModelID: "raport-pontaj", An: 2026, Luna: 10, Limba: "ro", SectieID: "s", DepartamentID: "d"}
	doc, err := m.s.Pregateste(ctx, c, momentFix)
	if err != nil {
		t.Fatal(err)
	}
	if len(doc.Randuri) != 2 || len(doc.Coloane) != 33 {
		t.Fatalf("dimensiuni: %d %d", len(doc.Randuri), len(doc.Coloane))
	}
	if doc.Randuri[0].Valori[24] != "9" || doc.Randuri[0].Valori[25] != "0" || doc.Randuri[0].Culori[25] != "ABSENTA" {
		t.Fatalf("coduri/absențe: %+v", doc.Randuri[0])
	}
	if !strings.Contains(doc.Metadate, versiune.Aplicatie) || !strings.Contains(doc.Nota, "fără calcul") {
		t.Fatal("lipsesc versiunile sau statutul de schiță")
	}
	c.An, c.Luna = 2028, 2
	c.AngajatID = "a"
	doc, err = m.s.Pregateste(ctx, c, momentFix)
	if err != nil || len(doc.Coloane) != 31 || doc.Randuri[0].Valori[30] != "9" {
		t.Fatalf("an bisect: %+v %v", doc, err)
	}
	c.ModelID, c.An, c.Luna = "raport-absente", 2026, 10
	doc, err = m.s.Pregateste(ctx, c, momentFix)
	if err != nil || len(doc.Randuri) != 1 || doc.Randuri[0].Valori[3] != "CO" {
		t.Fatalf("absențe active/an: %+v %v", doc, err)
	}
	c.ModelID, c.AngajatID = "raport-fisa", ""
	_, err = m.s.Pregateste(ctx, c, momentFix)
	codEroare(t, err, "rapoarte.alege_angajat")
	c.AngajatID, c.Revizie = "a", "veche"
	_, err = m.s.Pregateste(ctx, c, momentFix)
	codEroare(t, err, "rapoarte.conflict")
}

func TestReferinteRapoarteSiPDF(t *testing.T) {
	m := pregateste(t)
	fixture(t, m)
	for _, id := range []string{"raport-pontaj", "raport-fisa", "raport-absente", "raport-pontaj-detaliat", "raport-garzi"} {
		t.Run(id, func(t *testing.T) {
			c := Cerere{ModelID: id, An: 2026, Luna: 10, Limba: "ro"}
			if id == "raport-fisa" {
				c.AngajatID = "a"
			}
			doc, err := m.s.Pregateste(context.Background(), c, momentFix)
			if err != nil {
				t.Fatal(err)
			}
			pdf, err := GenereazaPDF(doc)
			if err != nil {
				t.Fatal(err)
			}
			if !bytes.HasPrefix(pdf, []byte("%PDF-")) {
				t.Fatal("fișierul nu este PDF")
			}
			if err := api.Validate(context.Background(), bytes.NewReader(pdf), nil, nil); err != nil {
				t.Fatalf("PDF invalid: %v", err)
			}
			if *directorDemo != "" {
				if err := os.MkdirAll(*directorDemo, 0700); err != nil {
					t.Fatal(err)
				}
				if err := os.WriteFile(filepath.Join(*directorDemo, id+".pdf"), pdf, 0600); err != nil {
					t.Fatal(err)
				}
			}
			doc.Metadate = strings.ReplaceAll(doc.Metadate, versiune.Aplicatie, "<aplicatie>")
			structura, err := ConstruiestePDF(doc)
			if err != nil {
				t.Fatal(err)
			}
			referinta, err := json.MarshalIndent(struct {
				Document  Document `json:"document"`
				Structura any      `json:"structura"`
			}{doc, structura.GetStructure()}, "", "  ")
			if err != nil {
				t.Fatal(err)
			}
			v := strconv.Itoa(doc.VersiuneRaport)
			cale := filepath.Join("testdata", id+"-v"+v+".json")
			if *actualizeazaReferinte {
				if err := os.MkdirAll("testdata", 0700); err != nil {
					t.Fatal(err)
				}
				if err := os.WriteFile(cale, referinta, 0600); err != nil {
					t.Fatal(err)
				}
			}
			asteptat, err := os.ReadFile(cale)
			if err != nil {
				t.Fatal(err)
			}
			if !bytes.Equal(asteptat, referinta) {
				t.Fatalf("raportul diferă de referință; verifică schimbarea și versiunea raportului: %s", cale)
			}
		})
	}
}

func TestPaginarePDFCuAntetSiSubsol(t *testing.T) {
	m := pregateste(t)
	fixture(t, m)
	doc, err := m.s.Pregateste(context.Background(), Cerere{ModelID: "raport-pontaj", An: 2026, Luna: 10, Limba: "en"}, momentFix)
	if err != nil {
		t.Fatal(err)
	}
	baza := doc.Randuri[0]
	doc.Randuri = nil
	for i := 0; i < 150; i++ {
		doc.Randuri = append(doc.Randuri, baza)
	}
	pdf, err := GenereazaPDF(doc)
	if err != nil {
		t.Fatal(err)
	}
	pagini, err := api.PageCount(context.Background(), bytes.NewReader(pdf), nil)
	if err != nil || pagini < 2 {
		t.Fatalf("paginare: %d %v", pagini, err)
	}
	if *directorDemo != "" {
		if err := os.WriteFile(filepath.Join(*directorDemo, "paginare-A3-demo.pdf"), pdf, 0600); err != nil {
			t.Fatal(err)
		}
	}
}
