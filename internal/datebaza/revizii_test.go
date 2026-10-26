package datebaza

import (
	"context"
	"fmt"
	"io/fs"
	"path/filepath"
	"testing"
	"testing/fstest"

	"pontix/internal/bd"
)

func TestReviziiDistincteCuAcelasiTimestamp(t *testing.T) {
	m := pregateste(t)
	ctx := context.Background()
	a := persoanaAbsente(t, m)
	const fix = "2026-10-25T10:00:00Z"
	for _, tabel := range []string{"pontaj_ciorne", "absente_zile"} {
		for _, eveniment := range []string{"INSERT", "UPDATE"} {
			comanda := fmt.Sprintf(`CREATE TRIGGER fixeaza_%s_%s AFTER %s ON %s WHEN NEW.modificat_la <> '%s'
                BEGIN UPDATE %s SET modificat_la = '%s' WHERE angajat_id = NEW.angajat_id AND data = NEW.data; END`, tabel, eveniment, eveniment, tabel, fix, tabel, fix)
			if _, err := m.db.Exec(comanda); err != nil {
				t.Fatal(err)
			}
		}
	}
	if err := m.serviciu.SalveazaSchita(ctx, 2026, 10, []CelulaSchita{{AngajatID: a.ID, Data: "2026-10-01", Valoare: "Z7"}}); err != nil {
		t.Fatal(err)
	}
	schita, err := m.serviciu.CitesteSchita(ctx, 2026, 10)
	if err != nil {
		t.Fatal(err)
	}
	veche := schita.Celule[0]
	noua := veche
	noua.Valoare = "IN9"
	if err := Nou(m.db).SalveazaSchita(ctx, 2026, 10, []CelulaSchita{noua}); err != nil {
		t.Fatal(err)
	}
	schita, err = m.serviciu.CitesteSchita(ctx, 2026, 10)
	if err != nil || schita.Celule[0].Revizie == veche.Revizie {
		t.Fatalf("revizie identică la ceas fix: %+v %v", schita, err)
	}
	verificaCod(t, m.serviciu.SalveazaSchita(ctx, 2026, 10, []CelulaSchita{veche}), "pontaj.conflict")

	if err := m.serviciu.SalveazaAbsente(ctx, []ZiAbsenta{{AngajatID: a.ID, Data: "2026-10-02", Cod: "CO", Activ: true}}); err != nil {
		t.Fatal(err)
	}
	ziVeche := absenteTest(t, m, 2026).Zile[0]
	anulata := ziVeche
	anulata.Activ = false
	if err := Nou(m.db).SalveazaAbsente(ctx, []ZiAbsenta{anulata}); err != nil {
		t.Fatal(err)
	}
	ziNoua := absenteTest(t, m, 2026).Zile[0]
	if ziNoua.Revizie == ziVeche.Revizie {
		t.Fatal("revizia anulării este identică")
	}
	verificaCod(t, m.serviciu.SalveazaAbsente(ctx, []ZiAbsenta{ziVeche}), "absente.conflict")
	ziNoua.Activ = true
	if err := m.serviciu.SalveazaAbsente(ctx, []ZiAbsenta{ziNoua}); err != nil {
		t.Fatal(err)
	}
	verificaCod(t, m.serviciu.SalveazaAbsente(ctx, []ZiAbsenta{ziNoua}), "absente.conflict")
	for _, tabel := range []string{"pontaj_ciorne", "absente_zile"} {
		var timestamp string
		if err := m.db.QueryRow(`SELECT modificat_la FROM ` + tabel + ` LIMIT 1`).Scan(&timestamp); err != nil || timestamp != fix {
			t.Fatalf("ceasul nu a fost fix: %s %v", timestamp, err)
		}
	}
}

func TestMigrareReviziiPastreazaContinutulSiMomentele(t *testing.T) {
	m := pregateste(t)
	ctx := context.Background()
	db, err := bd.Deschide(filepath.Join(t.TempDir(), "revizii-vechi.db"), m.cheie)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	surse := fstest.MapFS{}
	intrari, err := fs.ReadDir(bd.Migrari, ".")
	if err != nil {
		t.Fatal(err)
	}
	for _, r := range intrari {
		if r.Name() >= "0010_" {
			continue
		}
		continut, err := fs.ReadFile(bd.Migrari, r.Name())
		if err != nil {
			t.Fatal(err)
		}
		surse[r.Name()] = &fstest.MapFile{Data: continut}
	}
	if _, err := bd.Migreaza(ctx, db, surse); err != nil {
		t.Fatal(err)
	}
	for _, comanda := range []string{
		`INSERT INTO angajati(id,nume,marca,creat_la,modificat_la) VALUES ('a','Exemplu','R01','initial','initial')`,
		`INSERT INTO pontaj_ciorne(angajat_id,data,valoare,modificat_la) VALUES ('a','2026-10-01','IN9','moment-vechi')`,
		`INSERT INTO absente_zile(angajat_id,data,cod,activ,creat_la,modificat_la) VALUES ('a','2026-10-02','CO',0,'creare-veche','moment-vechi')`,
	} {
		if _, err := db.Exec(comanda); err != nil {
			t.Fatal(err)
		}
	}
	migrare, err := fs.ReadFile(bd.Migrari, "0010_revizii_concurenta.sql")
	if err != nil {
		t.Fatal(err)
	}
	surse["0010_revizii_concurenta.sql"] = &fstest.MapFile{Data: migrare}
	if n, err := bd.Migreaza(ctx, db, surse); err != nil || n != 1 {
		t.Fatalf("migrare: %d %v", n, err)
	}
	if _, err := bd.Migreaza(ctx, db, bd.Migrari); err != nil {
		t.Fatal(err)
	}
	s := Nou(db)
	schita, err := s.CitesteSchita(ctx, 2026, 10)
	if err != nil || schita.Celule[0].Valoare != "IN9" || schita.Celule[0].Revizie != "1" || len(schita.Absente) != 1 || schita.Absente[0].Activ {
		t.Fatalf("date schimbate: %+v %v", schita, err)
	}
	for _, tabel := range []string{"pontaj_ciorne", "absente_zile"} {
		var moment string
		if err := db.QueryRow(`SELECT modificat_la FROM ` + tabel).Scan(&moment); err != nil || moment != "moment-vechi" {
			t.Fatalf("moment schimbat: %s %v", moment, err)
		}
	}
	verificaCod(t, s.SalveazaSchita(ctx, 2026, 10, []CelulaSchita{{AngajatID: "a", Data: "2026-10-01", Valoare: "Z7", Revizie: "moment-vechi"}}), "pontaj.conflict")
}
