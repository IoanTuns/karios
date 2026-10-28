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

	"pontix/internal/audit"
	"pontix/internal/i18n"
)

// LungimeCheie este lungimea cheii de date: AES-256-XTS folosește două chei AES-256.
const LungimeCheie = 64

//go:embed migrari/*.sql
var fisiereMigrari embed.FS

// Deschide baza criptată de la cale, o creează dacă nu există și verifică cheia.
// O cheie greșită sau un fișier deteriorat dau *i18n.Eroare „bd.nu_se_deschide”.
// Conexiunile au funcțiile auditului, dar nu și trigger-ele lui: așa se verifică
// și se migrează baza (D044).
func Deschide(cale string, cheie []byte) (*sql.DB, error) {
	return deschide(cale, cheie, nil)
}

// DeschideCuAudit deschide baza pentru lucru: fiecare conexiune are trigger-ele care
// țin sigiliile la zi, iar paza refuză scrierile cât timp există un incident (D044).
func DeschideCuAudit(cale string, cheie []byte, paza *audit.Paza) (*sql.DB, error) {
	return deschide(cale, cheie, paza)
}

func deschide(cale string, cheie []byte, paza *audit.Paza) (*sql.DB, error) {
	if len(cheie) != LungimeCheie {
		return nil, &i18n.Eroare{Cod: "bd.cheie_invalida"}
	}
	hexCheie := hex.EncodeToString(cheie)

	// Cheia se dă prin PRAGMA la fiecare conexiune, nu în textul conexiunii, ca să
	// nu fie vizibilă altor părți ale programului.
	initializare := func(c *sqlite3.Conn) error {
		if err := c.Exec(`PRAGMA hexkey='` + hexCheie + `';
			PRAGMA journal_mode=WAL;
			PRAGMA synchronous=NORMAL;
			PRAGMA foreign_keys=ON;
			PRAGMA temp_store=MEMORY;
			PRAGMA busy_timeout=5000;`); err != nil {
			return err
		}
		return audit.PregatesteConexiune(c, cheie, paza)
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

// Migrari sunt migrările incluse în .exe.
var Migrari = func() fs.FS {
	surse, err := fs.Sub(fisiereMigrari, "migrari")
	if err != nil {
		panic(err)
	}
	return surse
}()

// StareMigrari întoarce versiunea bazei (0 pentru o bază nouă) și dacă există
// migrări neaplicate.
func StareMigrari(ctx context.Context, db *sql.DB, surse fs.FS) (versiune int64, inAsteptare bool, err error) {
	p, err := furnizorMigrari(ctx, db, surse)
	if err != nil {
		return 0, false, err
	}
	if versiune, err = p.GetDBVersion(ctx); err != nil {
		return 0, false, err
	}
	inAsteptare, err = p.HasPending(ctx)
	return versiune, inAsteptare, err
}

// Migreaza aplică migrările care lipsesc și întoarce câte s-au aplicat.
func Migreaza(ctx context.Context, db *sql.DB, surse fs.FS) (int, error) {
	p, err := furnizorMigrari(ctx, db, surse)
	if err != nil {
		return 0, err
	}
	rezultate, err := p.Up(ctx)
	return len(rezultate), err
}

// Backup scrie o copie consecventă a bazei la cale, criptată cu aceeași cheie
// (VACUUM INTO, D002). Merge și cu baza deschisă, în modul WAL.
func Backup(ctx context.Context, db *sql.DB, cale string, cheie []byte) error {
	// Fișierul țintă are nevoie de cheie la deschidere; aici apare scurt în textul
	// conexiunii, doar pe durata comenzii.
	tinta := "file:" + filepath.ToSlash(cale) + "?vfs=xts&hexkey=" + hex.EncodeToString(cheie)
	_, err := db.ExecContext(ctx, `VACUUM INTO ?`, tinta)
	return err
}
