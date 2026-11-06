package datebaza

import (
	"context"
	"database/sql"
)

func citesteLegaturiGarzi(ctx context.Context, tx *sql.Tx) (map[string]string, error) {
	rows, err := tx.QueryContext(ctx, "SELECT angajat_id,principal_id FROM angajati_garzi")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	rez := map[string]string{}
	for rows.Next() {
		var id, principal string
		if err := rows.Scan(&id, &principal); err != nil {
			return nil, err
		}
		rez[id] = principal
	}
	return rez, rows.Err()
}
func valideazaLegaturaGarda(ctx context.Context, tx *sql.Tx, a Angajat) error {
	if a.AngajatPrincipalID == "" {
		return nil
	}
	if a.ID != "" && a.ID == a.AngajatPrincipalID {
		return eroare("date.garda_legatura_invalida")
	}
	var n int
	err := tx.QueryRowContext(ctx, `SELECT count(*) FROM angajati WHERE id=? AND sters_la IS NULL`, a.AngajatPrincipalID).Scan(&n)
	if err != nil {
		return err
	}
	if n != 1 {
		return eroare("date.garda_legatura_invalida")
	}
	err = tx.QueryRowContext(ctx, `SELECT count(*) FROM angajati_garzi WHERE angajat_id=? OR principal_id=?`, a.AngajatPrincipalID, a.ID).Scan(&n)
	if err != nil {
		return err
	}
	if n != 0 {
		return eroare("date.garda_legatura_invalida")
	}
	return nil
}
