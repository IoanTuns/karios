package acces

import "context"

// ListaLegaturiIerarhice intoarce doar asocierile explicite ale marcilor subordonate.
func (q *Queries) ListaLegaturiIerarhice(ctx context.Context) (map[string]string, error) {
	rows, err := q.db.QueryContext(ctx, "SELECT angajat_id,principal_id FROM angajati_garzi")
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
