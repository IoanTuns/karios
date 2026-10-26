package datebaza

import (
	"context"
	"io/fs"
	"path/filepath"
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
		if intrare.Name() >= "0006_" {
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
	migrare, err := fs.ReadFile(bd.Migrari, "0006_departamente.sql")
	if err != nil {
		t.Fatal(err)
	}
	surse["0006_departamente.sql"] = &fstest.MapFile{Data: migrare}
	if n, err := bd.Migreaza(ctx, db, surse); err != nil || n != 1 {
		t.Fatalf("migrare: %d, %v", n, err)
	}
	d := citeste(t, Nou(db))
	if len(d.Angajati) != 1 || d.Angajati[0].SectieID != "s" || d.Angajati[0].DepartamentID != "" || len(d.Departamente) != 0 {
		t.Fatalf("date vechi schimbate: %+v", d)
	}
	// Verificăm schema istorică 0006 fără serviciul care cere schema curentă.
	var cod string
	var start, durata, pauza, numar int
	if err := db.QueryRow(`SELECT cod, ora_start, durata_minute, pauza_minute FROM tipuri_schimb WHERE id = 't'`).Scan(&cod, &start, &durata, &pauza); err != nil || cod != "CONFIG" || start != 330 || durata != 300 || pauza != 12 {
		t.Fatalf("configurație schimbată: %s %d %d %d %v", cod, start, durata, pauza, err)
	}
	if err := db.QueryRow(`SELECT count(*) FROM tipuri_schimb`).Scan(&numar); err != nil || numar != 1 {
		t.Fatalf("număr schimburi: %d %v", numar, err)
	}
}
