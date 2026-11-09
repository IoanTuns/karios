package bd

import (
	"context"
	"io/fs"
	"path/filepath"
	"strings"
	"testing"
	"testing/fstest"
)

func TestMigrareNomenclatoarePastreazaDatele(t *testing.T) {
	ctx := context.Background()
	db, err := Deschide(filepath.Join(t.TempDir(), "veche.db"), make([]byte, LungimeCheie))
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	vechi := fstest.MapFS{}
	entries, err := fs.ReadDir(Migrari, ".")
	if err != nil {
		t.Fatal(err)
	}
	for _, e := range entries {
		if strings.HasSuffix(e.Name(), ".sql") && e.Name() < "0024" {
			b, err := fs.ReadFile(Migrari, e.Name())
			if err != nil {
				t.Fatal(err)
			}
			vechi[e.Name()] = &fstest.MapFile{Data: b}
		}
	}
	if _, err := Migreaza(ctx, db, vechi); err != nil {
		t.Fatal(err)
	}
	if _, err := db.Exec("INSERT INTO angajati(id,nume,prenume,marca,functie,tip_contract,activ,creat_la,modificat_la) VALUES('fictiv','Exemplu','Ana','FICTIV','Medic existent','INTREAGA',1,'vechi','vechi')"); err != nil {
		t.Fatal(err)
	}
	var initial int
	_ = db.QueryRow("SELECT count(*) FROM tipuri_schimb").Scan(&initial)
	if _, err := Migreaza(ctx, db, Migrari); err != nil {
		t.Fatal(err)
	}
	var total int
	_ = db.QueryRow("SELECT count(*) FROM tipuri_schimb").Scan(&total)
	if total != initial {
		t.Fatal("programe pierdute")
	}
	var nume string
	if err := db.QueryRow("SELECT json_extract(f.denumire,'$.ro') FROM functii f JOIN angajati_functii a ON a.functie_id=f.id WHERE a.angajat_id='fictiv'").Scan(&nume); err != nil || nume != "Medic existent" {
		t.Fatal("functie veche nepreluata", err)
	}
	var categorii int
	_ = db.QueryRow("SELECT count(*) FROM categorii_schimburi WHERE id IN('LAB','CUSTOM')").Scan(&categorii)
	if categorii != 2 {
		t.Fatal("categorii initiale pierdute")
	}
	rows, err := db.Query("PRAGMA foreign_key_check")
	if err != nil {
		t.Fatal(err)
	}
	defer rows.Close()
	if rows.Next() {
		t.Fatal("referinte rupte dupa migrare")
	}
}

func TestMigrareMarciGardaNuAsociazaOmonimi(t *testing.T) {
	ctx := context.Background()
	db, err := Deschide(filepath.Join(t.TempDir(), "marci.db"), make([]byte, LungimeCheie))
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	vechi := fstest.MapFS{}
	entries, err := fs.ReadDir(Migrari, ".")
	if err != nil {
		t.Fatal(err)
	}
	for _, e := range entries {
		if strings.HasSuffix(e.Name(), ".sql") && e.Name() < "0025" {
			b, err := fs.ReadFile(Migrari, e.Name())
			if err != nil {
				t.Fatal(err)
			}
			vechi[e.Name()] = &fstest.MapFile{Data: b}
		}
	}
	if _, err := Migreaza(ctx, db, vechi); err != nil {
		t.Fatal(err)
	}
	if _, err := db.Exec(`INSERT INTO angajati(id,nume,prenume,marca,creat_la,modificat_la) VALUES('n','Exemplu','Ana','N','x','x'),('g','Exemplu','Ana','G','x','x')`); err != nil {
		t.Fatal(err)
	}
	if n, err := Migreaza(ctx, db, Migrari); err != nil || n != 13 {
		t.Fatal(n, err)
	}
	var n int
	if err := db.QueryRow("SELECT count(*) FROM angajati_garzi").Scan(&n); err != nil || n != 0 {
		t.Fatal("asociere automata", n, err)
	}
	if err := db.QueryRow("SELECT count(*) FROM angajati WHERE marca IN('N','G')").Scan(&n); err != nil || n != 2 {
		t.Fatal("marci pierdute", n, err)
	}
}
