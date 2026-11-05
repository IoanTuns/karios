package rapoarte

import (
	"context"
	"database/sql"

	"pontix/internal/i18n"
)

// prefixeRapoarte intoarce titulatura fiecarui angajat in limba raportului
// (D089). Marca de garda fara prefix propriu il preia pe al marcii principale,
// ca la functie.
func prefixeRapoarte(ctx context.Context, tx *sql.Tx, limba i18n.Limba) (map[string]string, error) {
	rows, err := tx.QueryContext(ctx, `SELECT a.id,COALESCE(p1.denumire,p2.denumire)
FROM angajati a
LEFT JOIN angajati_prefixe ap1 ON ap1.angajat_id=a.id LEFT JOIN prefixe p1 ON p1.id=ap1.prefix_id
LEFT JOIN angajati_garzi g ON g.angajat_id=a.id
LEFT JOIN angajati_prefixe ap2 ON ap2.angajat_id=g.principal_id LEFT JOIN prefixe p2 ON p2.id=ap2.prefix_id
WHERE COALESCE(p1.denumire,p2.denumire) IS NOT NULL`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	rez := map[string]string{}
	for rows.Next() {
		var id, text string
		if err := rows.Scan(&id, &text); err != nil {
			return nil, err
		}
		rez[id] = textJSON(text, limba)
	}
	return rez, rows.Err()
}
