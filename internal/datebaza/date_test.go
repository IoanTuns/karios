package datebaza

import (
	"context"
	"crypto/rand"
	"database/sql"
	"errors"
	"path/filepath"
	"testing"

	"pontix/internal/bd"
	"pontix/internal/i18n"
)

type mediuTest struct {
	db       *sql.DB
	serviciu *Serviciu
	cale     string
	cheie    []byte
}

func pregateste(t *testing.T) *mediuTest {
	t.Helper()
	m := &mediuTest{cale: filepath.Join(t.TempDir(), "date.db"), cheie: make([]byte, bd.LungimeCheie)}
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
	m.serviciu = Nou(m.db)
	return m
}

func citeste(t *testing.T, s *Serviciu) Date {
	t.Helper()
	d, err := s.Citeste(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	return d
}

func verificaCod(t *testing.T, err error, cod string) {
	t.Helper()
	var e *i18n.Eroare
	if !errors.As(err, &e) || e.Cod != cod {
		t.Fatalf("eroare: %v; așteptat codul %s", err, cod)
	}
}

func angajatTest(marca string) Angajat {
	return Angajat{Nume: " Exemplu ", Prenume: " Ana ", Marca: marca, TipContract: "INTREAGA", Activ: true}
}

func TestFluxAngajatCuPersistentaSiPontaj(t *testing.T) {
	m := pregateste(t)
	ctx := context.Background()
	if err := m.serviciu.SalveazaAngajat(ctx, angajatTest("0007")); err != nil {
		t.Fatal(err)
	}
	d := citeste(t, m.serviciu)
	if len(d.Angajati) != 1 || d.Angajati[0].Nume != "Exemplu" || d.Angajati[0].Prenume != "Ana" || d.Angajati[0].Marca != "0007" || d.Angajati[0].NormaMinute != 420 || d.Angajati[0].ConcediuZileAn != nil {
		t.Fatalf("angajat creat: %+v", d)
	}
	a := d.Angajati[0]
	minute, zile := int64(240), int64(0)
	a.TipContract, a.NormaPartialaMinute, a.ConcediuZileAn = "PARTIALA", &minute, &zile
	if err := m.serviciu.SalveazaAngajat(ctx, a); err != nil {
		t.Fatal(err)
	}
	if _, err := m.db.Exec(`INSERT INTO pontaj (id, angajat_id, data, start, sfarsit, creat_la, modificat_la)
		VALUES ('p1', ?, '2026-10-25', '2026-10-25T05:00:00Z', '2026-10-25T09:00:00Z', 'x', 'x')`, a.ID); err != nil {
		t.Fatal(err)
	}
	if err := m.serviciu.SchimbaActiv(ctx, a.ID, false); err != nil {
		t.Fatal(err)
	}
	if err := m.db.Close(); err != nil {
		t.Fatal(err)
	}
	var err error
	m.db, err = bd.Deschide(m.cale, m.cheie)
	if err != nil {
		t.Fatal(err)
	}
	m.serviciu = Nou(m.db)
	d = citeste(t, m.serviciu)
	if len(d.Angajati) != 1 || d.Angajati[0].ID != a.ID || d.Angajati[0].Activ || d.Angajati[0].NormaMinute != 240 || d.Angajati[0].ConcediuZileAn == nil || *d.Angajati[0].ConcediuZileAn != 0 {
		t.Fatalf("angajat după repornire: %+v", d.Angajati)
	}
	var n int
	if err := m.db.QueryRow(`SELECT count(*) FROM pontaj WHERE angajat_id = ?`, a.ID).Scan(&n); err != nil || n != 1 {
		t.Fatalf("pontaj pierdut: %d, %v", n, err)
	}
	if err := m.db.QueryRow(`SELECT count(*) FROM jurnal_modificari WHERE tabel = 'angajati' AND rand_id = ? AND cine = 'local'`, a.ID).Scan(&n); err != nil || n != 3 {
		t.Fatalf("jurnal: %d, %v", n, err)
	}
	if err := m.serviciu.SchimbaActiv(ctx, a.ID, true); err != nil {
		t.Fatal(err)
	}
	if !citeste(t, m.serviciu).Angajati[0].Activ {
		t.Fatal("reactivarea nu s-a salvat")
	}
}

func TestValidareSiRollbackAngajat(t *testing.T) {
	m := pregateste(t)
	ctx := context.Background()
	if err := m.serviciu.SalveazaAngajat(ctx, angajatTest("01")); err != nil {
		t.Fatal(err)
	}
	verificaCod(t, m.serviciu.SalveazaAngajat(ctx, angajatTest("01")), "date.marca_duplicata")
	cazuri := []struct {
		nume, cod string
		modifica  func(*Angajat)
	}{
		{"nume gol", "date.campuri_obligatorii", func(a *Angajat) { a.Nume = " " }},
		{"contract necunoscut", "date.contract_invalid", func(a *Angajat) { a.TipContract = "ALTUL" }},
		{"parțial fără normă", "date.norma_invalida", func(a *Angajat) { a.TipContract = "PARTIALA" }},
		{"concediu negativ", "date.concediu_invalid", func(a *Angajat) { n := int64(-1); a.ConcediuZileAn = &n }},
		{"secție inexistentă", "date.sectie_invalida", func(a *Angajat) { a.SectieID = "inexistenta" }},
		{"id inexistent", "date.inexistent", func(a *Angajat) { a.ID = "inexistent" }},
	}
	for _, c := range cazuri {
		t.Run(c.nume, func(t *testing.T) {
			a := angajatTest("02")
			c.modifica(&a)
			verificaCod(t, m.serviciu.SalveazaAngajat(ctx, a), c.cod)
		})
	}
	if len(citeste(t, m.serviciu).Angajati) != 1 {
		t.Fatal("validarea a lăsat angajați parțial salvați")
	}
	var n int
	if err := m.db.QueryRow(`SELECT count(*) FROM jurnal_modificari`).Scan(&n); err != nil || n != 1 {
		t.Fatalf("jurnal după operații refuzate: %d, %v", n, err)
	}
	// Un jurnal indisponibil trebuie să anuleze și modificarea angajatului.
	if _, err := m.db.Exec(`CREATE TRIGGER jurnal_blocat BEFORE INSERT ON jurnal_modificari BEGIN SELECT RAISE(ABORT, 'proba'); END`); err != nil {
		t.Fatal(err)
	}
	verificaCod(t, m.serviciu.SalveazaAngajat(ctx, angajatTest("03")), "date.eroare")
	if len(citeste(t, m.serviciu).Angajati) != 1 {
		t.Fatal("salvare fără jurnal acceptată")
	}
}

func TestDemoSiNormaGlobala(t *testing.T) {
	m := pregateste(t)
	ctx := context.Background()
	if !citeste(t, m.serviciu).DemoDisponibil {
		t.Fatal("demo indisponibil în baza goală")
	}
	if err := m.serviciu.SalveazaNorma(ctx, 450); err != nil {
		t.Fatal(err)
	}
	if err := m.serviciu.IncarcaDemo(ctx); err != nil {
		t.Fatal(err)
	}
	d := citeste(t, m.serviciu)
	if len(d.Angajati) != 12 || len(d.Sectii) != 3 || d.DemoDisponibil || d.NormaIntreagaMinute != 450 {
		t.Fatalf("date demo: %+v", d)
	}
	activi, partiali, necompletat, zero := 0, 0, 0, 0
	for _, a := range d.Angajati {
		if a.Activ {
			activi++
		}
		if a.TipContract == "PARTIALA" {
			partiali++
			if a.NormaMinute != *a.NormaPartialaMinute {
				t.Fatal("normă parțială modificată de setarea globală")
			}
		} else if a.NormaMinute != 450 {
			t.Fatal("norma globală nu a fost aplicată")
		}
		if a.ConcediuZileAn == nil {
			necompletat++
		} else if *a.ConcediuZileAn == 0 {
			zero++
		}
	}
	if activi != 10 || partiali != 3 || necompletat != 1 || zero != 1 {
		t.Fatalf("distribuție demo: %d activi, %d parțiali, %d necompletate, %d zero", activi, partiali, necompletat, zero)
	}
	verificaCod(t, m.serviciu.IncarcaDemo(ctx), "date.demo_nepermis")
	verificaCod(t, m.serviciu.SalveazaNorma(ctx, 0), "date.norma_invalida")
	if len(citeste(t, m.serviciu).Angajati) != 12 {
		t.Fatal("demo duplicat")
	}
	if err := m.serviciu.SalveazaNorma(ctx, 480); err != nil {
		t.Fatal(err)
	}
	for _, a := range citeste(t, m.serviciu).Angajati {
		if a.TipContract == "INTREAGA" && a.NormaMinute != 480 {
			t.Fatal("normă globală nesalvată")
		}
	}
}

func TestDemoRefuzatPesteDateSiTranzactional(t *testing.T) {
	ctx := context.Background()
	t.Run("angajat existent", func(t *testing.T) {
		m := pregateste(t)
		if err := m.serviciu.SalveazaAngajat(ctx, angajatTest("existenta")); err != nil {
			t.Fatal(err)
		}
		verificaCod(t, m.serviciu.IncarcaDemo(ctx), "date.demo_nepermis")
		if len(citeste(t, m.serviciu).Angajati) != 1 {
			t.Fatal("date existente modificate")
		}
		if _, err := m.db.Exec(`UPDATE angajati SET sters_la = 'x'`); err != nil {
			t.Fatal(err)
		}
		verificaCod(t, m.serviciu.IncarcaDemo(ctx), "date.demo_nepermis")
	})
	t.Run("eșec după primele inserări", func(t *testing.T) {
		m := pregateste(t)
		if _, err := m.db.Exec(`CREATE TRIGGER demo_blocat BEFORE INSERT ON angajati WHEN NEW.marca = 'DEMO-004' BEGIN SELECT RAISE(ABORT, 'proba'); END`); err != nil {
			t.Fatal(err)
		}
		verificaCod(t, m.serviciu.IncarcaDemo(ctx), "date.eroare")
		d := citeste(t, m.serviciu)
		if len(d.Angajati) != 0 || len(d.Sectii) != 0 || !d.DemoDisponibil {
			t.Fatal("încărcare demo parțială")
		}
		var n int
		if err := m.db.QueryRow(`SELECT count(*) FROM jurnal_modificari`).Scan(&n); err != nil || n != 0 {
			t.Fatalf("jurnal parțial: %d, %v", n, err)
		}
	})
}

func TestSectiiTraductibile(t *testing.T) {
	m := pregateste(t)
	ctx := context.Background()
	verificaCod(t, m.serviciu.AdaugaSectie(ctx, i18n.Text{i18n.Romana: " "}), "date.sectie_obligatorie")
	if err := m.serviciu.AdaugaSectie(ctx, i18n.Text{i18n.Romana: " Secție nouă ", i18n.Engleza: " New department "}); err != nil {
		t.Fatal(err)
	}
	d := citeste(t, m.serviciu)
	if len(d.Sectii) != 1 || d.Sectii[0].Denumire.In(i18n.Romana) != "Secție nouă" || d.Sectii[0].Denumire.In(i18n.Engleza) != "New department" {
		t.Fatalf("secție: %+v", d.Sectii)
	}
	a := angajatTest("01")
	a.SectieID = d.Sectii[0].ID
	if err := m.serviciu.SalveazaAngajat(ctx, a); err != nil {
		t.Fatal(err)
	}
	if citeste(t, m.serviciu).Angajati[0].SectieID != a.SectieID {
		t.Fatal("secție neatribuită")
	}
}
