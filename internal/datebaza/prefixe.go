package datebaza

import (
	"context"
	"database/sql"

	"pontix/internal/bd/acces"
)

// Prefixul angajatului (D089): titulatura optionala din nomenclatorul
// prefixe, legata prin angajati_prefixe. Numele salvat nu se schimba.

func citestePrefixeAngajati(ctx context.Context, tx *sql.Tx) (map[string]string, error) {
	rows, err := tx.QueryContext(ctx, "SELECT angajat_id,prefix_id FROM angajati_prefixe")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	rez := map[string]string{}
	for rows.Next() {
		var a, p string
		if err := rows.Scan(&a, &p); err != nil {
			return nil, err
		}
		rez[a] = p
	}
	return rez, rows.Err()
}

// completeazaPrefixe pune prefixul (ID si denumire) pe fiecare angajat.
func completeazaPrefixe(angajati []Angajat, prefixe []ElementNomenclator, legaturi map[string]string) {
	for i := range angajati {
		angajati[i].PrefixID = legaturi[angajati[i].ID]
		angajati[i].Prefix = nil
		for _, p := range prefixe {
			if p.ID == angajati[i].PrefixID {
				angajati[i].Prefix = p.Denumire
			}
		}
	}
}

// salveazaPrefixAngajat valideaza si inlocuieste legatura; un prefix inactiv
// ramane la angajatii care il au deja, dar nu poate fi ales nou.
func salveazaPrefixAngajat(ctx context.Context, tx *sql.Tx, q *acces.Queries, angajatID, prefixID string) error {
	var vechi string
	if err := tx.QueryRowContext(ctx, "SELECT prefix_id FROM angajati_prefixe WHERE angajat_id=?", angajatID).Scan(&vechi); err != nil && err != sql.ErrNoRows {
		return err
	}
	if vechi == prefixID {
		return nil
	}
	if prefixID != "" {
		var activ int
		if err := tx.QueryRowContext(ctx, "SELECT activ FROM prefixe WHERE id=?", prefixID).Scan(&activ); err != nil {
			return eroare("nomenclator.invalid")
		}
		if activ != 1 {
			return eroare("nomenclator.inactiv")
		}
	}
	if _, err := tx.ExecContext(ctx, "DELETE FROM angajati_prefixe WHERE angajat_id=?", angajatID); err != nil {
		return err
	}
	if prefixID != "" {
		if _, err := tx.ExecContext(ctx, "INSERT INTO angajati_prefixe VALUES(?,?)", angajatID, prefixID); err != nil {
			return err
		}
	}
	return jurnal(ctx, q, "angajati_prefixe", angajatID, vechi, prefixID)
}
