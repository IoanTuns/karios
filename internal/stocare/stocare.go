// Package stocare leagă seiful de baza de date la pornirea aplicației (F05-T12)
// și face backup-ul automat înainte de migrări (F05-T06).
package stocare

import (
	"context"
	"database/sql"
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"pontix/internal/bd"
	"pontix/internal/i18n"
	"pontix/internal/seif"
	"pontix/internal/versiune"
)

const (
	NumeSeif = "seif"
	NumeBaza = "pontix.db"

	// MaxBackupuri: câte copii automate se păstrează (decizie Ioan Tuns, 2026-10-25).
	MaxBackupuri = 7

	prefixBackup = "pontix-"
	sufixBackup  = ".db"
)

// Deschide deschide seiful și baza din folderul de date și aplică migrările.
// La prima pornire (fără seif și fără bază) le creează pe amândouă.
func Deschide(ctx context.Context, radacina string) (*sql.DB, error) {
	return deschide(ctx, radacina, bd.Migrari, time.Now)
}

func deschide(ctx context.Context, radacina string, migrari fs.FS, acum func() time.Time) (*sql.DB, error) {
	caleSeif := filepath.Join(radacina, NumeSeif)
	caleBaza := filepath.Join(radacina, NumeBaza)

	s, err := seif.Deschide(caleSeif)
	if errors.Is(err, fs.ErrNotExist) {
		// Fără seif, o bază existentă nu mai poate fi citită. Un seif nou ar avea
		// altă cheie și ar ascunde problema, deci se raportează ca seif deteriorat.
		if _, errBaza := os.Stat(caleBaza); errBaza == nil {
			return nil, seif.ErrDeteriorat
		}
		s, err = seif.Creeaza(caleSeif, bd.LungimeCheie)
	}
	if err != nil {
		return nil, err
	}

	db, err := bd.Deschide(caleBaza, s.CheieDate)
	if err != nil {
		return nil, err
	}
	if err := migreaza(ctx, db, radacina, s.CheieDate, migrari, acum); err != nil {
		db.Close()
		return nil, err
	}
	return db, nil
}

// migreaza face backup înainte de migrări, dacă baza are deja date (D002 punctul 3),
// apoi aplică migrările. Dacă backup-ul eșuează, migrarea nu se mai face.
func migreaza(ctx context.Context, db *sql.DB, radacina string, cheie []byte, migrari fs.FS, acum func() time.Time) error {
	versiuneBaza, inAsteptare, err := bd.StareMigrari(ctx, db, migrari)
	if err != nil {
		return err
	}
	if !inAsteptare {
		return nil
	}
	if versiuneBaza > 0 {
		dir := filepath.Join(radacina, "backup")
		if err := os.MkdirAll(dir, 0o700); err != nil {
			return err
		}
		nume := prefixBackup + acum().UTC().Format("2006-01-02-150405") + "-v" + versiune.Aplicatie + sufixBackup
		if err := bd.Backup(ctx, db, filepath.Join(dir, nume), cheie); err != nil {
			return &i18n.Eroare{Cod: "stocare.backup_esuat", Param: map[string]string{"detalii": err.Error()}}
		}
		if err := pastreazaUltimele(dir, MaxBackupuri); err != nil {
			return err
		}
	}
	_, err = bd.Migreaza(ctx, db, migrari)
	return err
}

// pastreazaUltimele șterge backup-urile automate mai vechi, astfel încât să rămână
// cel mult n. Numele încep cu data, deci ordinea alfabetică e ordinea în timp.
// Alte fișiere din folder nu se ating.
func pastreazaUltimele(dir string, n int) error {
	intrari, err := os.ReadDir(dir)
	if err != nil {
		return err
	}
	var backupuri []string
	for _, i := range intrari {
		if !i.IsDir() && strings.HasPrefix(i.Name(), prefixBackup) && strings.HasSuffix(i.Name(), sufixBackup) {
			backupuri = append(backupuri, i.Name())
		}
	}
	sort.Strings(backupuri)
	for len(backupuri) > n {
		if err := os.Remove(filepath.Join(dir, backupuri[0])); err != nil {
			return err
		}
		backupuri = backupuri[1:]
	}
	return nil
}
