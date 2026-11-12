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

// TipuriSarbatori: data → tipul zilei (LEGALA sau LIBERA, F04-T11, D107).
func (q *Queries) TipuriSarbatori(ctx context.Context) (map[string]string, error) {
	rows, err := q.db.QueryContext(ctx, "SELECT data,tip FROM sarbatori")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	rez := map[string]string{}
	for rows.Next() {
		var d, t string
		if err := rows.Scan(&d, &t); err != nil {
			return nil, err
		}
		rez[d] = t
	}
	return rez, rows.Err()
}

// SeteazaTipSarbatoare schimbă tipul unei zile din calendar.
func (q *Queries) SeteazaTipSarbatoare(ctx context.Context, data, tip string) error {
	_, err := q.db.ExecContext(ctx, "UPDATE sarbatori SET tip=? WHERE data=?", tip, data)
	return err
}
