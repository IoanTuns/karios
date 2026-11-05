package bd

import (
	"context"
	"io/fs"
	"path/filepath"
	"strings"
	"testing"
	"testing/fstest"
)

func TestMigrareProgrameMultiplePastreazaCiornaSiRevizia(t *testing.T) {
	db, e := Deschide(filepath.Join(t.TempDir(), "multi.db"), cheieNoua(t))
	if e != nil {
		t.Fatal(e)
	}
	defer db.Close()
	vechi := fstest.MapFS{}
	intrari, _ := fs.ReadDir(Migrari, ".")
	for _, v := range intrari {
		if strings.HasSuffix(v.Name(), ".sql") && v.Name() < "0031" {
			b, _ := fs.ReadFile(Migrari, v.Name())
			vechi[v.Name()] = &fstest.MapFile{Data: b}
		}
	}
	ctx := context.Background()
	if _, e = Migreaza(ctx, db, vechi); e != nil {
		t.Fatal(e)
	}
	if _, e = db.Exec(`INSERT INTO angajati(id,nume,marca,creat_la,modificat_la) VALUES('a','Exemplu','DEMO','x','x');INSERT INTO pontaj_ciorne(angajat_id,data,valoare,modificat_la,revizie) VALUES('a','2026-10-30','S7','vechi',9)`); e != nil {
		t.Fatal(e)
	}
	if n, e := Migreaza(ctx, db, Migrari); e != nil || n != 2 {
		t.Fatal(n, e)
	}
	var val, data string
	var rev int
	if e = db.QueryRow(`SELECT valoare,modificat_la,revizie FROM pontaj_ciorne`).Scan(&val, &data, &rev); e != nil || val != "S7" || data != "vechi" || rev != 9 {
		t.Fatal(val, data, rev, e)
	}
	if _, e = db.Exec(`UPDATE pontaj_ciorne SET valoare='IN9+S7'`); e != nil {
		t.Fatal(e)
	}
	var fk int
	if e = db.QueryRow(`PRAGMA foreign_keys`).Scan(&fk); e != nil || fk != 1 {
		t.Fatal(fk, e)
	}
	rows, e := db.Query(`PRAGMA foreign_key_check`)
	if e != nil {
		t.Fatal(e)
	}
	defer rows.Close()
	if rows.Next() {
		t.Fatal("legatura invalida")
	}
}
