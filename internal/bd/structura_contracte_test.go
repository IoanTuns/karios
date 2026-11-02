package bd

import (
	"context"
	"io/fs"
	"path/filepath"
	"strings"
	"testing"
	"testing/fstest"
)

func TestMigrareStructuraContractePastreazaDatele(t *testing.T) {
	db, err := Deschide(filepath.Join(t.TempDir(), "structura.db"), make([]byte, LungimeCheie))
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	ctx := context.Background()
	vechi := fstest.MapFS{}
	entries, err := fs.ReadDir(Migrari, ".")
	if err != nil {
		t.Fatal(err)
	}
	for _, e := range entries {
		if strings.HasSuffix(e.Name(), ".sql") && e.Name() < "0026" {
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
	if _, err := db.Exec(`INSERT INTO angajati(id,nume,prenume,marca,tip_contract,norma_partiala_minute,creat_la,modificat_la) VALUES('n','Exemplu','Ana','N','PARTIALA',225,'x','x'),('g','Exemplu','Ana','G','INTREAGA',NULL,'x','x');INSERT INTO angajati_garzi VALUES('g','n')`); err != nil {
		t.Fatal(err)
	}
	if n, err := Migreaza(ctx, db, Migrari); err != nil || n != 4 {
		t.Fatal(n, err)
	}
	var contract string
	if err := db.QueryRow("SELECT contract_id FROM angajati_contracte WHERE angajat_id='g'").Scan(&contract); err != nil || contract != "GARDA" {
		t.Fatal(contract, err)
	}
	if err := db.QueryRow("SELECT contract_id FROM angajati_contracte WHERE angajat_id='n'").Scan(&contract); err != nil || contract != "PARTIALA" {
		t.Fatal(contract, err)
	}
	var minute int
	if err := db.QueryRow("SELECT norma_partiala_minute FROM angajati WHERE id='n'").Scan(&minute); err != nil || minute != 225 {
		t.Fatal(minute, err)
	}
	rows, err := db.Query("PRAGMA foreign_key_check")
	if err != nil {
		t.Fatal(err)
	}
	defer rows.Close()
	if rows.Next() {
		t.Fatal("referinte rupte")
	}
}
