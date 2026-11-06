package acces

import "context"

// ListaToateSarbatori este interogarea explicita fara filtru de an.
// ListaSarbatori(ctx,nil) compara an=NULL si nu intoarce niciun rand.
func (q *Queries) ListaToateSarbatori(ctx context.Context) ([]ListaSarbatoriRow, error) {
	rows, err := q.db.QueryContext(ctx, "SELECT data,denumire FROM sarbatori ORDER BY data")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	rez := []ListaSarbatoriRow{}
	for rows.Next() {
		var z ListaSarbatoriRow
		if err := rows.Scan(&z.Data, &z.Denumire); err != nil {
			return nil, err
		}
		rez = append(rez, z)
	}
	return rez, rows.Err()
}
