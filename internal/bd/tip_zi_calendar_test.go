package bd

import (
	"context"
	"io/fs"
	"path/filepath"
	"strings"
	"testing"
	"testing/fstest"
)

// F04-T11 (D107): migrarea 0040 clasifica zilele existente: pentru 2025-2027 o zi care nu e in
// calendarul legal devine zi libera adaugata; ceilalti ani raman sarbatoare legala.
func TestMigrareTipZiCalendar(t *testing.T) {
	db, e := Deschide(filepath.Join(t.TempDir(), "calendar.db"), cheieNoua(t))
	if e != nil {
		t.Fatal(e)
	}
	defer db.Close()
	vechi := fstest.MapFS{}
	intrari, _ := fs.ReadDir(Migrari, ".")
	for _, v := range intrari {
		if strings.HasSuffix(v.Name(), ".sql") && v.Name() < "0040" {
			b, _ := fs.ReadFile(Migrari, v.Name())
			vechi[v.Name()] = &fstest.MapFile{Data: b}
		}
	}
	ctx := context.Background()
	if _, e = Migreaza(ctx, db, vechi); e != nil {
		t.Fatal(e)
	}
	if _, e = db.Exec(`INSERT OR IGNORE INTO sarbatori(data,denumire,creat_la,modificat_la) VALUES
		('2026-01-01','{"ro":"Anul Nou"}','x','x'),
		('2026-05-04','{"ro":"Punte"}','x','x'),
		('2028-03-01','{"ro":"Ceva"}','x','x')`); e != nil {
		t.Fatal(e)
	}
	if _, e = Migreaza(ctx, db, Migrari); e != nil {
		t.Fatal(e)
	}
	for data, asteptat := range map[string]string{"2026-01-01": "LEGALA", "2026-05-04": "LIBERA", "2028-03-01": "LEGALA"} {
		var tip string
		if e := db.QueryRow(`SELECT tip FROM sarbatori WHERE data=?`, data).Scan(&tip); e != nil || tip != asteptat {
			t.Errorf("%s: %q (%v), asteptat %s", data, tip, e, asteptat)
		}
	}
}
