package bd

import (
	"bytes"
	"context"
	"crypto/rand"
	"errors"
	"os"
	"path/filepath"
	"testing"

	"pontix/internal/i18n"
)

func cheieNoua(t *testing.T) []byte {
	t.Helper()
	c := make([]byte, LungimeCheie)
	if _, err := rand.Read(c); err != nil {
		t.Fatal(err)
	}
	return c
}

func deschideMigrata(t *testing.T, cale string, cheie []byte) {
	t.Helper()
	db, err := Deschide(cale, cheie)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	if _, err := Migreaza(context.Background(), db, Migrari); err != nil {
		t.Fatal(err)
	}
}

func TestDeschideSiMigreaza(t *testing.T) {
	cale := filepath.Join(t.TempDir(), "pontix.db")
	cheie := cheieNoua(t)
	db, err := Deschide(cale, cheie)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	n, err := Migreaza(context.Background(), db, Migrari)
	if err != nil {
		t.Fatal(err)
	}
	if n != 9 {
		t.Errorf("migrări aplicate: %d, așteptat 9", n)
	}

	// Setările fiecărei conexiuni.
	var mod string
	var fk int
	if err := db.QueryRow(`PRAGMA journal_mode`).Scan(&mod); err != nil || mod != "wal" {
		t.Errorf("journal_mode = %q, %v", mod, err)
	}
	if err := db.QueryRow(`PRAGMA foreign_keys`).Scan(&fk); err != nil || fk != 1 {
		t.Errorf("foreign_keys = %d, %v", fk, err)
	}

	for _, tabel := range []string{"setari", "sectii", "angajati", "tipuri_schimb", "coduri_absenta",
		"sarbatori", "categorii_ore", "seturi_reguli", "reguli_interval", "pontaj", "jurnal_modificari"} {
		var exista int
		if err := db.QueryRow(`SELECT count(*) FROM sqlite_schema WHERE type = 'table' AND name = ?`, tabel).Scan(&exista); err != nil || exista != 1 {
			t.Errorf("tabelul %s lipsește (%v)", tabel, err)
		}
	}

	// A doua rulare nu mai aplică nimic.
	if n, err := Migreaza(context.Background(), db, Migrari); err != nil || n != 0 {
		t.Errorf("a doua migrare: %d, %v", n, err)
	}
}

func TestConstrangeri(t *testing.T) {
	cale := filepath.Join(t.TempDir(), "pontix.db")
	cheie := cheieNoua(t)
	deschideMigrata(t, cale, cheie)
	db, err := Deschide(cale, cheie)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	acum := "2026-10-25T08:00:00Z"
	if _, err := db.Exec(`INSERT INTO angajati (id, nume, marca, creat_la, modificat_la) VALUES ('a1', 'Test', '101', ?, ?)`, acum, acum); err != nil {
		t.Fatal(err)
	}

	gresite := []struct {
		nume string
		sql  string
	}{
		{"angajat inexistent (cheie străină)", `INSERT INTO pontaj (id, angajat_id, data, start, sfarsit, creat_la, modificat_la)
			VALUES ('p1', 'lipsa', '2026-10-25', '2026-10-25T05:00:00Z', '2026-10-25T17:00:00Z', 'x', 'x')`},
		{"tură care se termină înainte să înceapă", `INSERT INTO pontaj (id, angajat_id, data, start, sfarsit, creat_la, modificat_la)
			VALUES ('p2', 'a1', '2026-10-25', '2026-10-25T17:00:00Z', '2026-10-25T05:00:00Z', 'x', 'x')`},
		{"rând fără ore și fără absență", `INSERT INTO pontaj (id, angajat_id, data, creat_la, modificat_la)
			VALUES ('p3', 'a1', '2026-10-25', 'x', 'x')`},
		{"marcă duplicată", `INSERT INTO angajati (id, nume, marca, creat_la, modificat_la) VALUES ('a2', 'Alt', '101', 'x', 'x')`},
		{"denumire care nu e JSON", `INSERT INTO sectii (id, denumire, creat_la, modificat_la) VALUES ('s1', 'Secția 1', 'x', 'x')`},
		{"dată invalidă", `INSERT INTO sarbatori (data, denumire, creat_la, modificat_la) VALUES ('2026-13-01', '{"ro":"x"}', 'x', 'x')`},
		{"interval după 24:00", `INSERT INTO seturi_reguli (id, denumire, valabil_de_la, creat_la, modificat_la) VALUES ('r1', 'x', '2026-01-01', 'x', 'x');
			INSERT INTO categorii_ore (cod, denumire, creat_la, modificat_la) VALUES ('ZI', '{"ro":"Zi"}', 'x', 'x');
			INSERT INTO reguli_interval (id, set_id, tip_zi, ora_start, ora_sfarsit, categorie_cod) VALUES ('i1', 'r1', 'LUCRATOARE', 360, 1500, 'ZI')`},
	}
	for _, g := range gresite {
		if _, err := db.Exec(g.sql); err == nil {
			t.Errorf("%s: aștepta refuz", g.nume)
		}
	}

	// O tură corectă și o absență se acceptă.
	if _, err := db.Exec(`INSERT INTO coduri_absenta (cod, denumire, creat_la, modificat_la) VALUES ('TEST_CM', '{"ro":"Absență test","en":"Test absence"}', ?, ?)`, acum, acum); err != nil {
		t.Fatal(err)
	}
	for _, s := range []string{
		`INSERT INTO pontaj (id, angajat_id, data, start, sfarsit, creat_la, modificat_la) VALUES ('p4', 'a1', '2026-10-24', '2026-10-24T16:00:00Z', '2026-10-25T05:00:00Z', 'x', 'x')`,
		`INSERT INTO pontaj (id, angajat_id, data, cod_absenta, creat_la, modificat_la) VALUES ('p5', 'a1', '2026-10-26', 'TEST_CM', 'x', 'x')`,
	} {
		if _, err := db.Exec(s); err != nil {
			t.Errorf("aștepta acceptare: %v", err)
		}
	}
}

func TestCheieGresita(t *testing.T) {
	cale := filepath.Join(t.TempDir(), "pontix.db")
	deschideMigrata(t, cale, cheieNoua(t))

	_, err := Deschide(cale, cheieNoua(t))
	var e *i18n.Eroare
	if !errors.As(err, &e) || e.Cod != "bd.nu_se_deschide" {
		t.Fatalf("aștepta bd.nu_se_deschide, a primit %v", err)
	}
	for _, l := range i18n.Limbi() {
		if e.Text(l) == e.Cod {
			t.Errorf("%s: mesajul nu are traducere", l)
		}
	}

	if _, err := Deschide(cale, []byte("scurta")); !errors.As(err, &e) || e.Cod != "bd.cheie_invalida" {
		t.Errorf("cheie de lungime greșită: %v", err)
	}
}

func TestNimicLizibilPeDisc(t *testing.T) {
	cale := filepath.Join(t.TempDir(), "pontix.db")
	cheie := cheieNoua(t)
	db, err := Deschide(cale, cheie)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := Migreaza(context.Background(), db, Migrari); err != nil {
		t.Fatal(err)
	}
	if _, err := db.Exec(`INSERT INTO angajati (id, nume, marca, creat_la, modificat_la) VALUES ('a1', 'Popescu Marcaj', '7', 'x', 'x')`); err != nil {
		t.Fatal(err)
	}
	for _, f := range []string{cale, cale + "-wal"} {
		data, err := os.ReadFile(f)
		if err != nil {
			t.Fatal(err)
		}
		for _, s := range []string{"Popescu Marcaj", "CREATE TABLE", "angajati"} {
			if bytes.Contains(data, []byte(s)) {
				t.Errorf("%s conține textul lizibil %q", filepath.Base(f), s)
			}
		}
	}
	db.Close()
}

func TestBackup(t *testing.T) {
	dir := t.TempDir()
	cale := filepath.Join(dir, "pontix.db")
	cheie := cheieNoua(t)
	db, err := Deschide(cale, cheie)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	if _, err := Migreaza(context.Background(), db, Migrari); err != nil {
		t.Fatal(err)
	}
	if _, err := db.Exec(`INSERT INTO angajati (id, nume, marca, creat_la, modificat_la) VALUES ('a1', 'Popescu Marcaj', '7', 'x', 'x')`); err != nil {
		t.Fatal(err)
	}

	copie := filepath.Join(dir, "backup.db")
	if err := Backup(context.Background(), db, copie, cheie); err != nil {
		t.Fatal(err)
	}
	data, err := os.ReadFile(copie)
	if err != nil {
		t.Fatal(err)
	}
	if bytes.Contains(data, []byte("Popescu Marcaj")) {
		t.Error("backup-ul conține text lizibil")
	}

	b, err := Deschide(copie, cheie)
	if err != nil {
		t.Fatal(err)
	}
	defer b.Close()
	var nume string
	if err := b.QueryRow(`SELECT nume FROM angajati WHERE id = 'a1'`).Scan(&nume); err != nil || nume != "Popescu Marcaj" {
		t.Errorf("backup redeschis: %q, %v", nume, err)
	}
	if _, err := Deschide(copie, cheieNoua(t)); err == nil {
		t.Error("backup-ul s-a deschis cu altă cheie")
	}
}
