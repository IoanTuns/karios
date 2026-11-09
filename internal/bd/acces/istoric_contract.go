package acces

import "context"

// SchimbareContract este un rand din istoric_contract (F07-T43, D101).
type SchimbareContract struct {
	AngajatID      string `json:"angajat_id"`
	DeLa           string `json:"de_la"`
	NormaMinute    *int64 `json:"norma_minute"`
	ConcediuZileAn *int64 `json:"concediu_zile_an"`
	CreatLa        string `json:"creat_la"`
}

// ListaIstoricContract intoarce schimbarile pe angajat, in ordinea datei.
func (q *Queries) ListaIstoricContract(ctx context.Context) (map[string][]SchimbareContract, error) {
	rows, err := q.db.QueryContext(ctx, "SELECT angajat_id,de_la,norma_minute,concediu_zile_an,creat_la FROM istoric_contract ORDER BY angajat_id,de_la")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	rez := map[string][]SchimbareContract{}
	for rows.Next() {
		var s SchimbareContract
		if err := rows.Scan(&s.AngajatID, &s.DeLa, &s.NormaMinute, &s.ConcediuZileAn, &s.CreatLa); err != nil {
			return nil, err
		}
		rez[s.AngajatID] = append(rez[s.AngajatID], s)
	}
	return rez, rows.Err()
}

func (q *Queries) AdaugaSchimbareContract(ctx context.Context, s SchimbareContract) (int64, error) {
	r, err := q.db.ExecContext(ctx, "INSERT INTO istoric_contract(angajat_id,de_la,norma_minute,concediu_zile_an,creat_la) VALUES(?,?,?,?,?) ON CONFLICT(angajat_id,de_la) DO NOTHING", s.AngajatID, s.DeLa, s.NormaMinute, s.ConcediuZileAn, s.CreatLa)
	if err != nil {
		return 0, err
	}
	return r.RowsAffected()
}

func (q *Queries) StergeSchimbareContract(ctx context.Context, angajatID, deLa string) (int64, error) {
	r, err := q.db.ExecContext(ctx, "DELETE FROM istoric_contract WHERE angajat_id=? AND de_la=?", angajatID, deLa)
	if err != nil {
		return 0, err
	}
	return r.RowsAffected()
}
