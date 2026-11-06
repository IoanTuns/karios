package acces

import (
	"context"
	"database/sql"
)

// PontajPrioritate pastreaza valorile complete pentru token si audit.
type PontajPrioritate struct {
	ID, AngajatID, Data, ModificatLa, Start, Sfarsit string
	Instantaneu                                      map[string]any
}

func (q *Queries) ListaPontajePrioritate(ctx context.Context, deLa, panaLa string) ([]PontajPrioritate, error) {
	rows, err := q.db.QueryContext(ctx, `SELECT * FROM pontaj WHERE sters_la IS NULL AND (data BETWEEN ?1 AND ?2 OR (julianday(start)<julianday(?2||'T23:59:59Z') AND julianday(sfarsit)>julianday(?1||'T00:00:00Z'))) ORDER BY id`, deLa, panaLa)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	cols, err := rows.Columns()
	if err != nil {
		return nil, err
	}
	rez := []PontajPrioritate{}
	for rows.Next() {
		v := make([]any, len(cols))
		d := make([]any, len(cols))
		for i := range v {
			d[i] = &v[i]
		}
		if err := rows.Scan(d...); err != nil {
			return nil, err
		}
		s := map[string]any{}
		for i, c := range cols {
			s[c] = v[i]
		}
		text := func(c string) string { v, _ := s[c].(string); return v }
		rez = append(rez, PontajPrioritate{ID: text("id"), AngajatID: text("angajat_id"), Data: text("data"), ModificatLa: text("modificat_la"), Start: text("start"), Sfarsit: text("sfarsit"), Instantaneu: s})
	}
	return rez, rows.Err()
}
func (q *Queries) RetragePontajPrioritate(ctx context.Context, p PontajPrioritate, acum string) (int64, error) {
	r, err := q.db.ExecContext(ctx, `UPDATE pontaj SET sters_la=?,modificat_la=? WHERE id=? AND modificat_la=? AND sters_la IS NULL`, acum, acum, p.ID, p.ModificatLa)
	if err != nil {
		return 0, err
	}
	return r.RowsAffected()
}
func (q *Queries) LunaDeschisaPrioritate(ctx context.Context, an, luna int) (bool, error) {
	s := "DESCHISA"
	err := q.db.QueryRowContext(ctx, "SELECT stare FROM stare_luni WHERE an=? AND luna=?", an, luna).Scan(&s)
	if err == sql.ErrNoRows {
		err = nil
	}
	return s == "DESCHISA", err
}
