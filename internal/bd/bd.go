// Package bd deschide baza de date criptată (D008) și îi aplică migrările (D010).
package bd

import (
	"context"
	"database/sql"
	"embed"
	"encoding/hex"
	"io/fs"
	"path/filepath"

	"github.com/ncruces/go-sqlite3"
	"github.com/ncruces/go-sqlite3/driver"
	_ "github.com/ncruces/go-sqlite3/vfs/xts" // criptarea AES-256-XTS (D008)
	"github.com/pressly/goose/v3"

	"pontix/internal/i18n"
)

// LungimeCheie este lungimea cheii de date: AES-256-XTS folosește două chei AES-256.
const LungimeCheie = 64

//go:embed migrari/*.sql
var fisiereMigrari embed.FS

// Deschide baza criptată de la cale, o creează dacă nu există și verifică cheia.
// O cheie greșită sau un fișier deteriorat dau *i18n.Eroare „bd.nu_se_deschide”.
func Deschide(cale string, cheie []byte) (*sql.DB, error) {
	if len(cheie) != LungimeCheie {
		return nil, &i18n.Eroare{Cod: "bd.cheie_invalida"}
	}
	hexCheie := hex.EncodeToString(cheie)

	// Cheia se dă prin PRAGMA la fiecare conexiune, nu în textul conexiunii, ca să
	// nu fie vizibilă altor părți ale programului.
	initializare := func(c *sqlite3.Conn) error {
		return c.Exec(`PRAGMA hexkey='` + hexCheie + `';
			PRAGMA journal_mode=WAL;
			PRAGMA synchronous=NORMAL;
			PRAGMA foreign_keys=ON;
			PRAGMA temp_store=MEMORY;
			PRAGMA busy_timeout=5000;`)
	}

	db, err := driver.Open("file:"+filepath.ToSlash(cale)+"?vfs=xts", initializare)
	if err != nil {
		return nil, nuSeDeschide(cale)
	}
	if err := db.Ping(); err != nil {
		db.Close()
		return nil, nuSeDeschide(cale)
	}
	// Cu o cheie greșită, abia prima citire eșuează.
	var n int
	if err := db.QueryRow(`SELECT count(*) FROM sqlite_schema`).Scan(&n); err != nil {
		db.Close()
		return nil, nuSeDeschide(cale)
	}
	return db, nil
}

func nuSeDeschide(cale string) error {
	return &i18n.Eroare{Cod: "bd.nu_se_deschide", Param: map[string]string{"fisier": cale}}
}

// Migreaza aplică migrările care lipsesc și întoarce câte s-au aplicat.
func Migreaza(ctx context.Context, db *sql.DB) (int, error) {
	surse, err := fs.Sub(fisiereMigrari, "migrari")
	if err != nil {
		return 0, err
	}
	p, err := goose.NewProvider(goose.DialectSQLite3, db, surse)
	if err != nil {
		return 0, err
	}
	rezultate, err := p.Up(ctx)
	return len(rezultate), err
}
