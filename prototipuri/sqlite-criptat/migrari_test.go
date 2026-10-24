package sqlitecriptat

// Prototipul F06-T01: migrările cu goose pe baza criptată (D009 punctul 3).

import (
	"context"
	"path/filepath"
	"testing"
	"testing/fstest"

	"github.com/pressly/goose/v3"
)

var migrari = fstest.MapFS{
	"0001_angajati.sql": {Data: []byte(`-- +goose Up
CREATE TABLE angajati (id INTEGER PRIMARY KEY, nume TEXT NOT NULL);
`)},
	"0002_sectii.sql": {Data: []byte(`-- +goose Up
ALTER TABLE angajati ADD COLUMN sectie TEXT;
`)},
}

func TestMigrari(t *testing.T) {
	ctx := context.Background()
	v := variante[1] // xts, varianta propusă
	db := deschide(t, filepath.Join(t.TempDir(), "pontix.db"), v.vfs, v.cheie)

	p, err := goose.NewProvider(goose.DialectSQLite3, db, migrari)
	if err != nil {
		t.Fatal(err)
	}

	// 1. Prima pornire aplică ambele migrări, în ordine.
	rezultate, err := p.Up(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if len(rezultate) != 2 {
		t.Fatalf("aplicate %d migrări, așteptat 2", len(rezultate))
	}

	// 2. Istoricul cerut de D009: versiunea aplicației lângă fiecare migrare.
	// Goose ține doar versiunea migrării și data; versiunea aplicației o adăugăm noi.
	if _, err := db.Exec(`CREATE TABLE IF NOT EXISTS istoric_migrari (
		migrare INTEGER PRIMARY KEY, versiune_aplicatie TEXT NOT NULL,
		aplicata_la TEXT NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ')))`); err != nil {
		t.Fatal(err)
	}
	for _, r := range rezultate {
		if _, err := db.Exec(`INSERT INTO istoric_migrari (migrare, versiune_aplicatie) VALUES (?, ?)`,
			r.Source.Version, "0.3.1"); err != nil {
			t.Fatal(err)
		}
	}
	var cuData int
	if err := db.QueryRow(`SELECT count(*) FROM goose_db_version WHERE version_id > 0 AND tstamp IS NOT NULL`).Scan(&cuData); err != nil || cuData != 2 {
		t.Errorf("goose_db_version: %d migrări cu dată, %v", cuData, err)
	}

	// 3. A doua pornire nu mai aplică nimic.
	if rezultate, err = p.Up(ctx); err != nil || len(rezultate) != 0 {
		t.Errorf("a doua pornire: %d migrări, %v", len(rezultate), err)
	}

	// 4. Protecția la downgrade: o versiune veche a aplicației știe doar de 0001,
	// dar baza e deja la 0002, deci trebuie să refuze deschiderea.
	veche, err := goose.NewProvider(goose.DialectSQLite3, db, fstest.MapFS{"0001_angajati.sql": migrari["0001_angajati.sql"]})
	if err != nil {
		t.Fatal(err)
	}
	versiuneBaza, err := veche.GetDBVersion(ctx)
	if err != nil {
		t.Fatal(err)
	}
	surse := veche.ListSources()
	cunoscuta := surse[len(surse)-1].Version
	if versiuneBaza <= cunoscuta {
		t.Errorf("baza la %d, aplicația veche cunoaște până la %d: downgrade-ul n-ar fi detectat", versiuneBaza, cunoscuta)
	}
}
