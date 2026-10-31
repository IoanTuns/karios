package acces

import "context"

type ClasificareSchimb struct {
	TipID, RegimID string
	Revizie        string `json:"-"`
}

func (q *Queries) ListaClasificariSchimburi(ctx context.Context) (map[string]ClasificareSchimb, error) {
	rows, err := q.db.QueryContext(ctx, "SELECT s.id,COALESCE(c.tip_id,s.tip),COALESCE(c.regim_id,s.regim),s.modificat_la FROM tipuri_schimb s LEFT JOIN clasificari_schimburi c ON c.schimb_id=s.id")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	rez := map[string]ClasificareSchimb{}
	for rows.Next() {
		var id string
		var c ClasificareSchimb
		if err := rows.Scan(&id, &c.TipID, &c.RegimID, &c.Revizie); err != nil {
			return nil, err
		}
		rez[id] = c
	}
	return rez, rows.Err()
}
