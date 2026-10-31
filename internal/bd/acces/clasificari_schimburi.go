package acces

import "context"

type ClasificareSchimb struct{ TipID, RegimID string }

func (q *Queries) ListaClasificariSchimburi(ctx context.Context) (map[string]ClasificareSchimb, error) {
	rows, err := q.db.QueryContext(ctx, "SELECT schimb_id,tip_id,regim_id FROM clasificari_schimburi")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	rez := map[string]ClasificareSchimb{}
	for rows.Next() {
		var id string
		var c ClasificareSchimb
		if err := rows.Scan(&id, &c.TipID, &c.RegimID); err != nil {
			return nil, err
		}
		rez[id] = c
	}
	return rez, rows.Err()
}
