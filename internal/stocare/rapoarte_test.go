package stocare

import (
	"context"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"testing/fstest"
	"time"

	"pontix/internal/bd"
)

func TestProbaRapoarteFaceBackupInainteaCompletariiCererilor(t *testing.T) {
	ctx := context.Background()
	surse := fstest.MapFS{}
	intrari, err := fs.ReadDir(bd.Migrari, ".")
	if err != nil {
		t.Fatal(err)
	}
	for _, intrare := range intrari {
		if intrare.Name() >= "0014_" || strings.HasPrefix(intrare.Name(), "0012_") {
			continue
		}
		continut, err := fs.ReadFile(bd.Migrari, intrare.Name())
		if err != nil {
			t.Fatal(err)
		}
		surse[intrare.Name()] = &fstest.MapFile{Data: continut}
	}
	dir := t.TempDir()
	acum := func() time.Time { return time.Date(2026, 10, 25, 13, 0, 0, 0, time.UTC) }
	db, _, err := deschide(ctx, dir, surse, acum)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := db.Exec(`UPDATE definitii_rapoarte SET revizie=6,antet=json_set(antet,'$.stanga.ro','Model fictiv pastrat') WHERE id='raport-pontaj'`); err != nil {
		db.Close()
		t.Fatal(err)
	}
	if err := db.Close(); err != nil {
		t.Fatal(err)
	}
	db, _, err = deschide(ctx, dir, bd.Migrari, acum)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	var curent int
	if err := db.QueryRow(`SELECT count(*) FROM pragma_table_info('absente_zile') WHERE name='stare_aprobare'`).Scan(&curent); err != nil || curent != 1 {
		t.Fatalf("actualizare cereri: %d %v", curent, err)
	}
	copii, err := os.ReadDir(filepath.Join(dir, "backup"))
	if err != nil || len(copii) != 1 {
		t.Fatalf("backup: %d %v", len(copii), err)
	}
	s := seifTest(t, dir)
	copie, err := bd.Deschide(filepath.Join(dir, "backup", copii[0].Name()), s.CheieDate)
	if err != nil {
		t.Fatal(err)
	}
	defer copie.Close()
	var coloane, revizie int
	var antet string
	if err := copie.QueryRow(`SELECT count(*) FROM pragma_table_info('absente_zile') WHERE name='stare_aprobare'`).Scan(&coloane); err != nil || coloane != 0 {
		t.Fatalf("backupul nu precede completarea 0012: %d %v", coloane, err)
	}
	if err := copie.QueryRow(`SELECT revizie,json_extract(antet,'$.stanga.ro') FROM definitii_rapoarte WHERE id='raport-pontaj'`).Scan(&revizie, &antet); err != nil || revizie != 6 || antet != "Model fictiv pastrat" {
		t.Fatalf("modelul din backup nu este pastrat: %d %q %v", revizie, antet, err)
	}
	if err := db.QueryRow(`SELECT revizie,json_extract(antet,'$.stanga.ro') FROM definitii_rapoarte WHERE id='raport-pontaj'`).Scan(&revizie, &antet); err != nil || revizie != 6 || antet != "Model fictiv pastrat" {
		t.Fatalf("modelul actualizat nu este pastrat: %d %q %v", revizie, antet, err)
	}
}
