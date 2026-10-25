package datebaza

import (
	"context"
	"io/fs"
	"path/filepath"
	"strings"
	"testing"
	"testing/fstest"

	"pontix/internal/bd"
	"pontix/internal/i18n"
)

func TestSectieCuMaiMulteDepartamente(t *testing.T) {
	m := pregateste(t)
	ctx := context.Background()
	for _, nume := range []string{"Nord", "Sud"} {
		if err := m.serviciu.AdaugaSectie(ctx, i18n.Text{i18n.Romana: nume}); err != nil {
			t.Fatal(err)
		}
	}
	d := citeste(t, m.serviciu)
	nord, sud := d.Sectii[0].ID, d.Sectii[1].ID
	for _, nume := range []string{" Laborator ", "Biochimie"} {
		if err := m.serviciu.AdaugaDepartament(ctx, nord, i18n.Text{i18n.Romana: nume, i18n.Engleza: "Test department"}); err != nil {
			t.Fatal(err)
		}
	}
	if err := m.serviciu.AdaugaDepartament(ctx, sud, i18n.Text{i18n.Romana: "Laborator"}); err != nil {
		t.Fatal(err)
	}
	d = citeste(t, m.serviciu)
	if len(d.Departamente) != 3 {
		t.Fatalf("departamente: %+v", d.Departamente)
	}
	var deptNord, deptSud string
	for _, dep := range d.Departamente {
		if dep.SectieID == nord {
			deptNord = dep.ID
		} else if dep.SectieID == sud {
			deptSud = dep.ID
		}
	}
	a := angajatTest("01")
	a.SectieID, a.DepartamentID = nord, deptNord
	if err := m.serviciu.SalveazaAngajat(ctx, a); err != nil {
		t.Fatal(err)
	}
	a = citeste(t, m.serviciu).Angajati[0]
	if a.DepartamentID != deptNord {
		t.Fatal("departament neatribuit")
	}
	if err := m.serviciu.SalveazaSchita(ctx, 2026, 10, []CelulaSchita{{AngajatID: a.ID, Data: "2026-10-01", Valoare: "7"}}); err != nil {
		t.Fatal(err)
	}
	a.SectieID = sud
	verificaCod(t, m.serviciu.SalveazaAngajat(ctx, a), "date.departament_invalid")
	a.DepartamentID = ""
	if err := m.serviciu.SalveazaAngajat(ctx, a); err != nil {
		t.Fatal(err)
	}
	a.DepartamentID = deptSud
	if err := m.serviciu.SalveazaAngajat(ctx, a); err != nil {
		t.Fatal(err)
	}
	d = citeste(t, m.serviciu)
	if d.Angajati[0].SectieID != sud || d.Angajati[0].DepartamentID != deptSud {
		t.Fatal("schimbarea organizării nu s-a salvat")
	}
	luna, err := m.serviciu.CitesteSchita(ctx, 2026, 10)
	if err != nil || len(luna.Celule) != 1 || luna.Celule[0].Valoare != "7" {
		t.Fatalf("schiță pierdută: %+v, %v", luna, err)
	}
	verificaCod(t, m.serviciu.AdaugaDepartament(ctx, "inexistenta", i18n.Text{i18n.Romana: "Test"}), "date.sectie_invalida")
	verificaCod(t, m.serviciu.AdaugaDepartament(ctx, nord, i18n.Text{i18n.Romana: " "}), "date.departament_obligatoriu")
	a.DepartamentID, a.SectieID = deptNord, ""
	verificaCod(t, m.serviciu.SalveazaAngajat(ctx, a), "date.departament_invalid")

	// Integritatea este apărată și în afara serviciului, direct de bază.
	if _, err := m.db.Exec(`INSERT INTO angajati (id, nume, prenume, marca, sectie_id, departament_id, creat_la, modificat_la) VALUES ('direct', 'Exemplu', 'Test', '02', ?, ?, 'x', 'x')`, nord, deptNord); err != nil {
		t.Fatal(err)
	}
	for _, proba := range []struct {
		sql  string
		args []any
	}{
		{`UPDATE angajati SET sectie_id = ? WHERE id = 'direct'`, []any{sud}},
		{`UPDATE departamente SET sectie_id = ? WHERE id = ?`, []any{sud, deptNord}},
		{`INSERT INTO angajati (id, nume, marca, sectie_id, departament_id, creat_la, modificat_la) VALUES ('gresit', 'Test', '03', ?, ?, 'x', 'x')`, []any{sud, deptNord}},
	} {
		if _, err := m.db.Exec(proba.sql, proba.args...); err == nil {
			t.Fatal("relație secție–departament invalidă acceptată de bază")
		}
	}
}

func TestMigrareDepartamentePastreazaOrganizareaSiSchimburile(t *testing.T) {
	ctx := context.Background()
	surse := fstest.MapFS{}
	intrari, err := fs.ReadDir(bd.Migrari, ".")
	if err != nil {
		t.Fatal(err)
	}
	for _, intrare := range intrari {
		if strings.HasPrefix(intrare.Name(), "0006_") {
			continue
		}
		continut, err := fs.ReadFile(bd.Migrari, intrare.Name())
		if err != nil {
			t.Fatal(err)
		}
		surse[intrare.Name()] = &fstest.MapFile{Data: continut}
	}
	cheie := make([]byte, bd.LungimeCheie)
	for i := range cheie {
		cheie[i] = byte(i + 1)
	}
	db, err := bd.Deschide(filepath.Join(t.TempDir(), "veche.db"), cheie)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	if _, err := bd.Migreaza(ctx, db, surse); err != nil {
		t.Fatal(err)
	}
	for _, comanda := range []string{
		`INSERT INTO sectii (id, denumire, creat_la, modificat_la) VALUES ('s', '{"ro":"Secție existentă"}', 'x', 'x')`,
		`INSERT INTO angajati (id, nume, marca, sectie_id, creat_la, modificat_la) VALUES ('a', 'Exemplu', '01', 's', 'x', 'x')`,
		`INSERT INTO tipuri_schimb (id, cod, denumire, ora_start, durata_minute, pauza_minute, tip, creat_la, modificat_la) VALUES ('t', 'CONFIG', '{"ro":"Configurat"}', 330, 300, 12, 'NORMAL', 'initial', 'initial')`,
	} {
		if _, err := db.Exec(comanda); err != nil {
			t.Fatal(err)
		}
	}
	if n, err := bd.Migreaza(ctx, db, bd.Migrari); err != nil || n != 1 {
		t.Fatalf("migrare: %d, %v", n, err)
	}
	d := citeste(t, Nou(db))
	if len(d.Angajati) != 1 || d.Angajati[0].SectieID != "s" || d.Angajati[0].DepartamentID != "" || len(d.Departamente) != 0 {
		t.Fatalf("date vechi schimbate: %+v", d)
	}
	schimburi, err := Nou(db).ListaSchimburi(ctx)
	if err != nil || len(schimburi) != 1 || schimburi[0].Cod != "CONFIG" || schimburi[0].OraStart != 330 || schimburi[0].DurataMinute != 300 || schimburi[0].PauzaMinute != 12 {
		t.Fatalf("configurație schimbată: %+v, %v", schimburi, err)
	}
}
