// Prototipul F05-T09: verifică dacă ncruces/go-sqlite3 cu VFS de criptare
// îndeplinește cerințele din D008, înainte să intre în aplicație.
package sqlitecriptat

import (
	"bytes"
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	_ "github.com/ncruces/go-sqlite3/driver"
	_ "github.com/ncruces/go-sqlite3/vfs/adiantum"
	_ "github.com/ncruces/go-sqlite3/vfs/xts"
)

// marcaj apare în date; nu trebuie să se vadă în niciun fișier de pe disc.
const marcaj = "Popescu Ion Marcaj"

var variante = []struct {
	vfs    string
	cheie  string // hex
	gresit string // hex, altă cheie de aceeași lungime
}{
	{"adiantum", strings.Repeat("ab", 32), strings.Repeat("cd", 32)},
	{"xts", strings.Repeat("ab", 64), strings.Repeat("cd", 64)}, // AES-256-XTS: 64 de octeți
}

func dsn(cale, vfs, cheie string) string {
	return fmt.Sprintf("file:%s?vfs=%s&hexkey=%s", filepath.ToSlash(cale), vfs, cheie)
}

func deschide(t *testing.T, cale, vfs, cheie string) *sql.DB {
	t.Helper()
	db, err := sql.Open("sqlite3", dsn(cale, vfs, cheie))
	if err != nil {
		t.Fatal(err)
	}
	db.SetMaxOpenConns(1)
	t.Cleanup(func() { db.Close() })
	return db
}

func numara(db *sql.DB) (int, error) {
	var n int
	err := db.QueryRow(`SELECT count(*) FROM segmente`).Scan(&n)
	return n, err
}

func contineMarcaj(t *testing.T, cai ...string) []string {
	t.Helper()
	var gasit []string
	for _, c := range cai {
		data, err := os.ReadFile(c)
		if os.IsNotExist(err) {
			continue
		}
		if err != nil {
			t.Fatal(err)
		}
		if bytes.Contains(data, []byte(marcaj)) || bytes.Contains(data, []byte("CREATE TABLE")) {
			gasit = append(gasit, filepath.Base(c))
		}
	}
	return gasit
}

// TestControl arată că verificarea de text lizibil chiar funcționează:
// într-o bază necriptată, marcajul trebuie găsit.
func TestControl(t *testing.T) {
	cale := filepath.Join(t.TempDir(), "necriptat.db")
	db, err := sql.Open("sqlite3", "file:"+filepath.ToSlash(cale))
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	if _, err := db.Exec(`CREATE TABLE segmente (angajat TEXT)`); err != nil {
		t.Fatal(err)
	}
	if _, err := db.Exec(`INSERT INTO segmente VALUES (?)`, marcaj); err != nil {
		t.Fatal(err)
	}
	db.Close()
	if gasit := contineMarcaj(t, cale); len(gasit) == 0 {
		t.Fatal("verificarea nu găsește textul nici într-o bază necriptată")
	}
}

func TestPrototip(t *testing.T) {
	for _, v := range variante {
		t.Run(v.vfs, func(t *testing.T) {
			dir := t.TempDir()
			cale := filepath.Join(dir, "pontix.db")
			db := deschide(t, cale, v.vfs, v.cheie)

			// 1. Modul WAL și setările din D002/D008.
			var mod string
			if err := db.QueryRow(`PRAGMA journal_mode=WAL`).Scan(&mod); err != nil || mod != "wal" {
				t.Fatalf("journal_mode: %q, %v", mod, err)
			}
			for _, p := range []string{`PRAGMA foreign_keys=ON`, `PRAGMA temp_store=memory`} {
				if _, err := db.Exec(p); err != nil {
					t.Fatalf("%s: %v", p, err)
				}
			}

			// 2. O lună de pontaj: 50 de angajați × 31 de zile × 3 segmente.
			if _, err := db.Exec(`CREATE TABLE segmente (
				id INTEGER PRIMARY KEY, angajat TEXT NOT NULL, data TEXT NOT NULL,
				categorie TEXT NOT NULL, minute INTEGER NOT NULL)`); err != nil {
				t.Fatal(err)
			}
			start := time.Now()
			tx, err := db.Begin()
			if err != nil {
				t.Fatal(err)
			}
			ins, err := tx.Prepare(`INSERT INTO segmente (angajat, data, categorie, minute) VALUES (?, ?, ?, ?)`)
			if err != nil {
				t.Fatal(err)
			}
			for a := 0; a < 50; a++ {
				for z := 1; z <= 31; z++ {
					for _, c := range []string{"ZI", "NOAPTE", "WEEKEND"} {
						if _, err := ins.Exec(fmt.Sprintf("%s %02d", marcaj, a), fmt.Sprintf("2026-10-%02d", z), c, 240); err != nil {
							t.Fatal(err)
						}
					}
				}
			}
			ins.Close()
			if err := tx.Commit(); err != nil {
				t.Fatal(err)
			}
			scriere := time.Since(start)

			start = time.Now()
			rows, err := db.Query(`SELECT angajat, categorie, sum(minute) FROM segmente GROUP BY angajat, categorie`)
			if err != nil {
				t.Fatal(err)
			}
			totaluri := 0
			for rows.Next() {
				totaluri++
			}
			rows.Close()
			citire := time.Since(start)
			t.Logf("4650 de rânduri scrise în %v; totalurile (%d grupuri) calculate în %v", scriere, totaluri, citire)

			// 3. Pe disc nu se vede nimic lizibil, nici în jurnalul WAL.
			if gasit := contineMarcaj(t, cale, cale+"-wal"); len(gasit) > 0 {
				t.Errorf("text lizibil pe disc în %v", gasit)
			}

			// 4. Backup cu VACUUM INTO, criptat cu aceeași cheie.
			backup := filepath.Join(dir, "backup.db")
			if _, err := db.Exec(`VACUUM INTO ?`, dsn(backup, v.vfs, v.cheie)); err != nil {
				t.Fatalf("VACUUM INTO: %v", err)
			}
			if gasit := contineMarcaj(t, backup); len(gasit) > 0 {
				t.Errorf("backup-ul conține text lizibil")
			}
			if n, err := numara(deschide(t, backup, v.vfs, v.cheie)); err != nil || n != 4650 {
				t.Errorf("backup redeschis: %d rânduri, %v", n, err)
			}
			if _, err := numara(deschide(t, backup, v.vfs, v.gresit)); err == nil {
				t.Error("backup-ul s-a deschis cu o cheie greșită")
			}

			// 5. Cheie greșită sau fără criptare: baza nu se poate citi.
			db.Close()
			if _, err := numara(deschide(t, cale, v.vfs, v.gresit)); err == nil {
				t.Error("baza s-a deschis cu o cheie greșită")
			}
			fara, err := sql.Open("sqlite3", "file:"+filepath.ToSlash(cale))
			if err != nil {
				t.Fatal(err)
			}
			if _, err := numara(fara); err == nil {
				t.Error("baza s-a citit fără criptare")
			}
			fara.Close()

			// 6. Redeschisă cu cheia corectă, datele sunt toate acolo.
			if n, err := numara(deschide(t, cale, v.vfs, v.cheie)); err != nil || n != 4650 {
				t.Errorf("redeschidere: %d rânduri, %v", n, err)
			}
		})
	}
}
