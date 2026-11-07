package bd

import (
	"context"
	"crypto/rand"
	"path/filepath"
	"testing"
)

// F08-T12 (D092): codurile vechi se redenumesc cu zilele lor; un cod nou deja existent
// nu se amesteca cu cel vechi, iar codurile neatinse raman cum sunt.
func TestMigrare0035RedenumesteCodurileAbsentei(t *testing.T) {
	ctx := context.Background()
	cheie := make([]byte, LungimeCheie)
	if _, err := rand.Read(cheie); err != nil {
		t.Fatal(err)
	}
	db, err := Deschide(filepath.Join(t.TempDir(), "d092.db"), cheie)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	p, err := furnizorMigrari(ctx, db, Migrari)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := p.UpTo(ctx, 34); err != nil {
		t.Fatal(err)
	}
	exec := func(q string, args ...any) {
		t.Helper()
		if _, err := db.ExecContext(ctx, q, args...); err != nil {
			t.Fatalf("%s: %v", q, err)
		}
	}
	exec(`INSERT INTO angajati (id, nume, prenume, marca, creat_la, modificat_la) VALUES ('a', 'Pop', 'Ana', 'M1', 'x', 'x')`)
	// CEV, CFS, DEPL se redenumesc; FP nu, pentru ca FMP exista deja.
	for _, cod := range []string{"CEV", "CFS", "DEPL", "FP", "FMP"} {
		exec(`INSERT INTO coduri_absenta (cod, denumire, categorie, creat_la, modificat_la) VALUES (?, '{"ro":"Tip `+cod+`"}', 'SPECIAL', 'x', 'x')`, cod)
	}
	exec(`UPDATE coduri_absenta SET efect_ore = 'NORMA' WHERE cod = 'DEPL'`)
	for i, cod := range []string{"CEV", "CFS", "DEPL", "FP", "CO"} {
		exec(`INSERT INTO absente_zile (angajat_id, data, cod, creat_la, modificat_la) VALUES ('a', ?, ?, 'x', 'x')`, "2026-10-0"+string(rune('1'+i)), cod)
	}
	if _, err := p.UpTo(ctx, 35); err != nil {
		t.Fatal(err)
	}
	coduri := map[string]string{}
	rows, err := db.QueryContext(ctx, `SELECT cod, denumire, efect_ore FROM coduri_absenta`)
	if err != nil {
		t.Fatal(err)
	}
	for rows.Next() {
		var cod, den, efect string
		if err := rows.Scan(&cod, &den, &efect); err != nil {
			t.Fatal(err)
		}
		coduri[cod] = den + "|" + efect
	}
	rows.Close()
	for _, vechi := range []string{"CEV", "CFS", "DEPL"} {
		if _, ok := coduri[vechi]; ok {
			t.Errorf("%s a ramas", vechi)
		}
	}
	if coduri["CVE"] != `{"ro":"Tip CEV"}|ANTERIOR` || coduri["CFP"] == "" || coduri["DLG"] != `{"ro":"Tip DEPL"}|NORMA` {
		t.Errorf("tipurile redenumite nu si-au pastrat datele: %v", coduri)
	}
	if coduri["FP"] == "" || coduri["FMP"] != `{"ro":"Tip FMP"}|ANTERIOR` {
		t.Errorf("FP si FMP existente nu trebuiau amestecate: %v", coduri)
	}
	zile := map[string]string{}
	rows, err = db.QueryContext(ctx, `SELECT data, cod FROM absente_zile`)
	if err != nil {
		t.Fatal(err)
	}
	for rows.Next() {
		var d, c string
		if err := rows.Scan(&d, &c); err != nil {
			t.Fatal(err)
		}
		zile[d] = c
	}
	rows.Close()
	asteptat := map[string]string{"2026-10-01": "CVE", "2026-10-02": "CFP", "2026-10-03": "DLG", "2026-10-04": "FP", "2026-10-05": "CO"}
	for d, c := range asteptat {
		if zile[d] != c {
			t.Errorf("%s: %s, astept %s", d, zile[d], c)
		}
	}
	var fk int
	if err := db.QueryRowContext(ctx, `SELECT count(*) FROM pragma_foreign_key_check`).Scan(&fk); err != nil || fk != 0 {
		t.Fatalf("legaturi rupte: %d %v", fk, err)
	}
	// Coborarea reface codurile vechi.
	if _, err := p.DownTo(ctx, 34); err != nil {
		t.Fatal(err)
	}
	var n int
	if err := db.QueryRowContext(ctx, `SELECT count(*) FROM absente_zile WHERE cod IN ('CEV','CFS','DEPL')`).Scan(&n); err != nil || n != 3 {
		t.Fatalf("coborare: %d %v", n, err)
	}
}
