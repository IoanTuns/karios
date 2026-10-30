package stocare

import (
	"context"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"testing/fstest"
	"time"

	"pontix/internal/bd"
	"pontix/internal/seif"
)

func TestPrimaPornireSiRedeschidere(t *testing.T) {
	ctx := context.Background()
	dir := t.TempDir()

	db, _, err := deschideTest(ctx, dir)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := db.Exec(`INSERT INTO setari (cheie, valoare, modificat_la) VALUES ('proba', 'da', 'x')`); err != nil {
		t.Fatal(err)
	}
	db.Close()
	for _, f := range []string{NumeSeif, NumeBaza} {
		if _, err := os.Stat(filepath.Join(dir, f)); err != nil {
			t.Errorf("%s nu a fost creat: %v", f, err)
		}
	}

	// A doua pornire găsește aceeași cheie și aceleași date.
	db, _, err = deschideTest(ctx, dir)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	var v string
	if err := db.QueryRow(`SELECT valoare FROM setari WHERE cheie = 'proba'`).Scan(&v); err != nil || v != "da" {
		t.Errorf("datele după redeschidere: %q, %v", v, err)
	}
}

// Dacă seiful lipsește, dar baza există, nu se creează un seif nou.
func TestBazaFaraSeif(t *testing.T) {
	ctx := context.Background()
	dir := t.TempDir()
	db, _, err := deschideTest(ctx, dir)
	if err != nil {
		t.Fatal(err)
	}
	db.Close()
	if err := os.Remove(filepath.Join(dir, NumeSeif)); err != nil {
		t.Fatal(err)
	}

	if _, _, err := deschideTest(ctx, dir); !errors.Is(err, seif.ErrDeteriorat) {
		t.Fatalf("aștepta seif.ErrDeteriorat, a primit %v", err)
	}
	if _, err := os.Stat(filepath.Join(dir, NumeSeif)); !os.IsNotExist(err) {
		t.Error("s-a creat un seif nou peste o bază existentă")
	}
}

// migrariPana întoarce migrările reale plus n-1 migrări de probă (0002 … 000n).
func migrariPana(t *testing.T, n int) fs.FS {
	t.Helper()
	initiala, err := fs.ReadFile(bd.Migrari, "0001_schema_initiala.sql")
	if err != nil {
		t.Fatal(err)
	}
	m := fstest.MapFS{"0001_schema_initiala.sql": {Data: initiala}}
	for i := 2; i <= n; i++ {
		m[fmt.Sprintf("%04d_proba.sql", i)] = &fstest.MapFile{Data: []byte(fmt.Sprintf(
			"-- +goose Up\nCREATE TABLE proba_%d (x INTEGER) STRICT;\n", i))}
	}
	return m
}

func ceas(minute int) func() time.Time {
	return func() time.Time { return time.Date(2026, 10, 25, 8, minute, 0, 0, time.UTC) }
}

func backupuri(t *testing.T, dir string) []string {
	t.Helper()
	intrari, err := os.ReadDir(filepath.Join(dir, "backup"))
	if os.IsNotExist(err) {
		return nil
	}
	if err != nil {
		t.Fatal(err)
	}
	var rez []string
	for _, i := range intrari {
		rez = append(rez, i.Name())
	}
	return rez
}

func TestBackupInainteDeMigrare(t *testing.T) {
	ctx := context.Background()
	dir := t.TempDir()

	// Baza nouă: nimic de salvat, deci niciun backup.
	db, _, err := deschide(ctx, dir, migrariPana(t, 1), ceas(0))
	if err != nil {
		t.Fatal(err)
	}
	if _, err := db.Exec(`INSERT INTO setari (cheie, valoare, modificat_la) VALUES ('proba', 'inainte', 'x')`); err != nil {
		t.Fatal(err)
	}
	db.Close()
	if b := backupuri(t, dir); len(b) != 0 {
		t.Fatalf("backup la o bază nouă: %v", b)
	}

	// O versiune nouă aduce migrarea 0002: întâi backup, apoi migrare.
	db, _, err = deschide(ctx, dir, migrariPana(t, 2), ceas(1))
	if err != nil {
		t.Fatal(err)
	}
	db.Close()
	b := backupuri(t, dir)
	if len(b) != 1 {
		t.Fatalf("backup-uri după migrare: %v", b)
	}
	if !strings.HasPrefix(b[0], "pontix-2026-10-25-080100-v") {
		t.Errorf("numele backup-ului: %s", b[0])
	}

	// Backup-ul are datele de dinainte și schema de dinainte de migrare.
	s := seifTest(t, dir)
	copie, err := bd.Deschide(filepath.Join(dir, "backup", b[0]), s.CheieDate)
	if err != nil {
		t.Fatal(err)
	}
	var v string
	if err := copie.QueryRow(`SELECT valoare FROM setari WHERE cheie = 'proba'`).Scan(&v); err != nil || v != "inainte" {
		t.Errorf("datele din backup: %q, %v", v, err)
	}
	var n int
	copie.QueryRow(`SELECT count(*) FROM sqlite_schema WHERE name = 'proba_2'`).Scan(&n)
	if n != 0 {
		t.Error("backup-ul conține deja migrarea 0002")
	}
	copie.Close()

	// Fără migrări noi, nu se face alt backup.
	db, _, err = deschide(ctx, dir, migrariPana(t, 2), ceas(2))
	if err != nil {
		t.Fatal(err)
	}
	db.Close()
	if b := backupuri(t, dir); len(b) != 1 {
		t.Errorf("backup fără migrări noi: %v", b)
	}
}

// Verifică traseul de pornire cu migrările reale din F07-T06, nu doar cu
// migrări sintetice: copia păstrează schema veche, baza curentă primește câmpurile noi.
func TestBackupMigrarilorContractuale(t *testing.T) {
	ctx := context.Background()
	dir := t.TempDir()
	db, _, err := deschide(ctx, dir, migrariPana(t, 1), ceas(0))
	if err != nil {
		t.Fatal(err)
	}
	if _, err := db.Exec(`INSERT INTO angajati (id, nume, marca, norma_minute, creat_la, modificat_la)
		VALUES ('test', 'Exemplu', '001', 480, 'x', 'x')`); err != nil {
		db.Close()
		t.Fatal(err)
	}
	db.Close()

	db, _, err = deschide(ctx, dir, bd.Migrari, ceas(1))
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	var tip string
	var minute int
	if err := db.QueryRow(`SELECT tip_contract, norma_efectiva_minute FROM angajati_norme WHERE id = 'test'`).Scan(&tip, &minute); err != nil || tip != "INTREAGA" || minute != 420 {
		t.Fatalf("contract după migrare: %s, %d minute, %v", tip, minute, err)
	}
	copii := backupuri(t, dir)
	if len(copii) != 1 {
		t.Fatalf("backup-uri pentru cele două migrări noi: %v", copii)
	}
	s := seifTest(t, dir)
	copie, err := bd.Deschide(filepath.Join(dir, "backup", copii[0]), s.CheieDate)
	if err != nil {
		t.Fatal(err)
	}
	defer copie.Close()
	if err := copie.QueryRow(`SELECT norma_minute FROM angajati WHERE id = 'test'`).Scan(&minute); err != nil || minute != 480 {
		t.Fatalf("norma originală din backup: %d, %v", minute, err)
	}
	if _, err := copie.Exec(`SELECT prenume FROM angajati`); err == nil {
		t.Error("backup-ul conține deja schema nouă")
	}
	// Închiderea copiei elimină fișierele WAL/SHM create la verificarea ei.
	if err := copie.Close(); err != nil {
		t.Fatal(err)
	}
	// Redeschiderea aceleiași versiuni nu repetă migrările sau backup-ul.
	db.Close()
	db, _, err = deschide(ctx, dir, bd.Migrari, ceas(2))
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	if copii := backupuri(t, dir); len(copii) != 1 {
		t.Errorf("backup repetat la redeschidere: %v", copii)
	}
}

// După 10 actualizări rămân doar cele mai noi 7 backup-uri; alte fișiere nu se ating.
func TestCelMult7Backupuri(t *testing.T) {
	ctx := context.Background()
	dir := t.TempDir()
	db, _, err := deschide(ctx, dir, migrariPana(t, 1), ceas(0))
	if err != nil {
		t.Fatal(err)
	}
	db.Close()
	if err := os.MkdirAll(filepath.Join(dir, "backup"), 0o700); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "backup", "notite.txt"), []byte("ale utilizatorului"), 0o600); err != nil {
		t.Fatal(err)
	}

	for i := 2; i <= 11; i++ {
		db, _, err := deschide(ctx, dir, migrariPana(t, i), ceas(i))
		if err != nil {
			t.Fatal(err)
		}
		db.Close()
	}

	b := backupuri(t, dir)
	var automate []string
	for _, f := range b {
		if strings.HasPrefix(f, "pontix-") {
			automate = append(automate, f)
		}
	}
	if len(automate) != MaxBackupuri {
		t.Fatalf("backup-uri păstrate: %d, așteptat %d (%v)", len(automate), MaxBackupuri, automate)
	}
	if !strings.HasPrefix(automate[0], "pontix-2026-10-25-080500-") || !strings.HasPrefix(automate[6], "pontix-2026-10-25-081100-") {
		t.Errorf("nu s-au păstrat cele mai noi: %v", automate)
	}
	if len(b) != MaxBackupuri+1 {
		t.Errorf("fișierul utilizatorului din backup/ a fost atins: %v", b)
	}
}

func TestBackupManualPastreazaCopiileAutomateSiCriptarea(t *testing.T) {
	ctx := context.Background()
	dir := t.TempDir()
	db, _, err := deschide(ctx, dir, migrariPana(t, 1), ceas(0))
	if err != nil {
		t.Fatal(err)
	}
	if _, err := db.Exec(`INSERT INTO setari (cheie, valoare, modificat_la) VALUES ('manual', 'ok', 'x')`); err != nil {
		db.Close()
		t.Fatal(err)
	}
	cale1, err := backupManual(ctx, db, dir, seifTest(t, dir).CheieDate, ceas(5))
	if err != nil {
		db.Close()
		t.Fatal(err)
	}
	cale2, err := backupManual(ctx, db, dir, seifTest(t, dir).CheieDate, ceas(5))
	if err != nil {
		db.Close()
		t.Fatal(err)
	}
	db.Close()
	if cale1 == cale2 {
		t.Fatalf("două copii în aceeași secundă au primit aceeași cale: %s", cale1)
	}
	if !strings.HasPrefix(filepath.Base(cale1), "pontix-manual-2026-10-25-080500-v") {
		t.Fatalf("nume neașteptat pentru primul backup manual: %s", filepath.Base(cale1))
	}
	if !strings.HasPrefix(filepath.Base(cale2), "pontix-manual-2026-10-25-080500-2-v") {
		t.Fatalf("nume neașteptat pentru al doilea backup manual: %s", filepath.Base(cale2))
	}

	s := seifTest(t, dir)
	copie, err := bd.Deschide(cale1, s.CheieDate)
	if err != nil {
		t.Fatal(err)
	}
	var val string
	if err := copie.QueryRow(`SELECT valoare FROM setari WHERE cheie = 'manual'`).Scan(&val); err != nil || val != "ok" {
		copie.Close()
		t.Fatalf("valoare din backup manual: %q, %v", val, err)
	}
	copie.Close()

	for i := 2; i <= 10; i++ {
		db, _, err := deschide(ctx, dir, migrariPana(t, i), ceas(10+i))
		if err != nil {
			t.Fatal(err)
		}
		db.Close()
	}
	b := backupuri(t, dir)
	var automate, manuale int
	for _, f := range b {
		if strings.HasPrefix(f, "pontix-manual-") {
			manuale++
		} else if strings.HasPrefix(f, "pontix-") {
			automate++
		}
	}
	if manuale != 2 || automate != MaxBackupuri {
		t.Fatalf("după migrări: manuale=%d (așteptat 2), automate=%d (așteptat %d), toate=%v", manuale, automate, MaxBackupuri, b)
	}
}
