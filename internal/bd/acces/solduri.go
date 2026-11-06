package acces

import (
	"context"
	"database/sql"
)

type InitializarePontaj struct {
	Luna    string `json:"luna"`
	Minute  int64  `json:"minute"`
	COZile  int64  `json:"co_zile"`
	Revizie int64  `json:"revizie"`
}

func (q *Queries) ListaInitializariPontaj(ctx context.Context) (map[string]InitializarePontaj, error) {
	rows, err := q.db.QueryContext(ctx, "SELECT angajat_id,luna,minute,co_zile,revizie FROM initializari_pontaj")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	rez := map[string]InitializarePontaj{}
	for rows.Next() {
		var id string
		var r InitializarePontaj
		if err := rows.Scan(&id, &r.Luna, &r.Minute, &r.COZile, &r.Revizie); err != nil {
			return nil, err
		}
		rez[id] = r
	}
	return rez, rows.Err()
}
func (q *Queries) SalveazaInitializarePontaj(ctx context.Context, id string, r InitializarePontaj, acum string) (int64, error) {
	var n sql.Result
	var err error
	if r.Revizie == 0 {
		n, err = q.db.ExecContext(ctx, `INSERT INTO initializari_pontaj(angajat_id,luna,minute,co_zile,modificat_la) VALUES(?,?,?,?,?) ON CONFLICT(angajat_id) DO NOTHING`, id, r.Luna, r.Minute, r.COZile, acum)
	} else {
		n, err = q.db.ExecContext(ctx, `UPDATE initializari_pontaj SET luna=?,minute=?,co_zile=?,modificat_la=?,revizie=revizie+1 WHERE angajat_id=? AND revizie=?`, r.Luna, r.Minute, r.COZile, acum, id, r.Revizie)
	}

	if err != nil {
		return 0, err
	}
	return n.RowsAffected()
}
func (q *Queries) PrimaLunaPontaj(ctx context.Context) (string, error) {
	var s string
	err := q.db.QueryRowContext(ctx, `SELECT coalesce(substr(min(data),1,7),'') FROM (SELECT data FROM pontaj_ciorne WHERE trim(valoare)<>'' UNION ALL SELECT data FROM absente_zile WHERE activ=1 UNION ALL SELECT data FROM planificare WHERE activ=1)`).Scan(&s)
	return s, err
}
func (q *Queries) NormeSoldPontaj(ctx context.Context) (map[string]*int64, error) {
	rows, err := q.db.QueryContext(ctx, `SELECT a.id,CASE WHEN ag.principal_id IS NOT NULL OR coalesce(tc.garda,0)=1 OR tc.normare IN ('FARA','LUNARA') THEN NULL WHEN tc.normare='ZILNICA' THEN tc.norma_minute WHEN tc.normare='INDIVIDUALA' OR (tc.normare IS NULL AND a.tip_contract='PARTIALA') THEN a.norma_partiala_minute ELSE CAST((SELECT valoare FROM setari WHERE cheie='norma_intreaga_minute') AS INTEGER) END FROM angajati a LEFT JOIN angajati_contracte ac ON ac.angajat_id=a.id LEFT JOIN tipuri_contract tc ON tc.id=ac.contract_id LEFT JOIN angajati_garzi ag ON ag.angajat_id=a.id WHERE a.sters_la IS NULL`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	r := map[string]*int64{}
	for rows.Next() {
		var id string
		var n *int64
		if err := rows.Scan(&id, &n); err != nil {
			return nil, err
		}
		r[id] = n
	}
	return r, rows.Err()
}

func (q *Queries) SoldAfecteazaLunaInchisa(ctx context.Context, luna string) (bool, error) {
	var n int
	err := q.db.QueryRowContext(ctx, "SELECT count(*) FROM stare_luni WHERE stare<>'DESCHISA' AND printf('%04d-%02d',an,luna)>=?", luna).Scan(&n)
	return n > 0, err
}
