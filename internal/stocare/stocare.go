// Package stocare leagă seiful de baza de date la pornirea aplicației (F05-T12).
package stocare

import (
	"context"
	"database/sql"
	"errors"
	"io/fs"
	"os"
	"path/filepath"

	"pontix/internal/bd"
	"pontix/internal/seif"
)

const (
	NumeSeif = "seif"
	NumeBaza = "pontix.db"
)

// Deschide deschide seiful și baza din folderul de date și aplică migrările.
// La prima pornire (fără seif și fără bază) le creează pe amândouă.
func Deschide(ctx context.Context, radacina string) (*sql.DB, error) {
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
	if _, err := bd.Migreaza(ctx, db); err != nil {
		db.Close()
		return nil, err
	}
	return db, nil
}
