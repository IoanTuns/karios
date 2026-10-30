package stocare

import (
	"context"
	"crypto/hmac"
	"database/sql"

	"pontix/internal/audit"
	"pontix/internal/bd"
	"pontix/internal/seif"
)

// Contorul din seif detectează înlocuirea bazei cu o copie mai veche, dar validă (F06-T09,
// D044 punctul 7): seiful ține ultimul id din jurnal văzut de aplicație, autentificat cu cheia
// de audit; o bază cu mai puține intrări sau un contor falsificat dau incident.

// verificaSecventa compară baza cu contorul din seif. Fără contor (prima pornire cu F06-T09)
// nu e incident: contorul se scrie acum.
func verificaSecventa(ctx context.Context, caleBaza string, s *seif.Seif) (bool, error) {
	c := s.Secventa()
	if c == nil {
		return false, nil
	}
	db, err := bd.Deschide(caleBaza, s.CheieDate)
	if err != nil {
		return false, err
	}
	defer db.Close()
	n, err := audit.UltimulRandJurnal(ctx, db)
	if err != nil {
		return false, err
	}
	falsificat := !hmac.Equal(c.Mac, audit.MacSecventa(s.CheieDate, c.Numar))
	return falsificat || n < c.Numar, nil
}

func scrieSecventa(ctx context.Context, db *sql.DB, s *seif.Seif) error {
	n, err := audit.UltimulRandJurnal(ctx, db)
	if err != nil {
		return err
	}
	if c := s.Secventa(); c != nil && c.Numar == n {
		return nil
	}
	return s.SeteazaSecventa(seif.Secventa{Numar: n, Mac: audit.MacSecventa(s.CheieDate, n)})
}

// ActualizeazaSecventa scrie în seif ultimul id din jurnalul bazei deschise; se apelează la
// închidere și după confirmarea unui incident.
func ActualizeazaSecventa(ctx context.Context, db *sql.DB, s *seif.Seif) error {
	return scrieSecventa(ctx, db, s)
}
