package bd

import (
	"context"
	"io/fs"
	"path/filepath"
	"strings"
	"testing"
	"testing/fstest"
)

func TestMigrareClasificariPastreazaProgrameSiReferinte(t *testing.T) {
	db, err := Deschide(filepath.Join(t.TempDir(), "clasificari.db"), cheieNoua(t))
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	vechi := fstest.MapFS{}
	intrari, _ := fs.ReadDir(Migrari, ".")
	for _, e := range intrari {
		if strings.HasSuffix(e.Name(), ".sql") && e.Name() < "0027" {
			b, _ := fs.ReadFile(Migrari, e.Name())
			vechi[e.Name()] = &fstest.MapFile{Data: b}
		}
	}
	ctx := context.Background()
	if _, err := Migreaza(ctx, db, vechi); err != nil {
		t.Fatal(err)
	}
	if _, err := db.Exec(`UPDATE tipuri_schimb SET regim='GENERAL',ora_start=1200,durata_minute=660,activ=0 WHERE cod='N12'`); err != nil {
		t.Fatal(err)
	}
	if n, err := Migreaza(ctx, db, Migrari); err != nil || n != 16 {
		t.Fatal(n, err)
	}
	var tip, regim, tipID, regimID string
	var ora, durata, activ int
	if err := db.QueryRow(`SELECT s.tip,s.regim,c.tip_id,c.regim_id,s.ora_start,s.durata_minute,s.activ FROM tipuri_schimb s JOIN clasificari_schimburi c ON c.schimb_id=s.id WHERE s.cod='N12'`).Scan(&tip, &regim, &tipID, &regimID, &ora, &durata, &activ); err != nil {
		t.Fatal(err)
	}
	if tip != "NOAPTE" || tipID != "NOAPTE" || regim != "LUCRATOARE" || regimID != "LUCRATOARE" || ora != 1200 || durata != 660 || activ != 0 {
		t.Fatal("migrare a reinterpretat programul", tip, regim, ora, durata, activ)
	}
	if _, err := db.Exec("UPDATE tipuri_schimb SET regim='GENERAL' WHERE cod='N12'"); err != nil {
		t.Fatal(err)
	}
	if err := db.QueryRow("SELECT regim FROM tipuri_schimb WHERE cod='N12'").Scan(&regim); err != nil || regim != "LUCRATOARE" {
		t.Fatal("alias vechi neconvertit", regim, err)
	}
	var n int
	if err := db.QueryRow("SELECT count(*) FROM pragma_foreign_key_check").Scan(&n); err != nil || n != 0 {
		t.Fatal("referinte rupte", n, err)
	}
}
