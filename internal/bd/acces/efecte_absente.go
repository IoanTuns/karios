package acces

import "context"

func (q *Queries) EfecteAbsente(ctx context.Context) (map[string]string, error) {
	r, e := q.db.QueryContext(ctx, "SELECT cod, efect_ore FROM coduri_absenta")
	if e != nil {
		return nil, e
	}
	defer r.Close()
	rez := map[string]string{}
	for r.Next() {
		var cod, efect string
		if e := r.Scan(&cod, &efect); e != nil {
			return nil, e
		}
		rez[cod] = efect
	}
	return rez, r.Err()
}
func (q *Queries) SalveazaEfectAbsenta(ctx context.Context, cod, efect string) error {
	_, e := q.db.ExecContext(ctx, "UPDATE coduri_absenta SET efect_ore=? WHERE cod=?", efect, cod)
	return e
}
func (q *Queries) MinuteAbsente(ctx context.Context, deLa, panaLa string) (map[string]*int64, error) {
	r, e := q.db.QueryContext(ctx, "SELECT angajat_id,data,minute_creditate FROM absente_zile WHERE data BETWEEN ? AND ?", deLa, panaLa)
	if e != nil {
		return nil, e
	}
	defer r.Close()
	rez := map[string]*int64{}
	for r.Next() {
		var id, data string
		var n *int64
		if e := r.Scan(&id, &data, &n); e != nil {
			return nil, e
		}
		rez[id+"/"+data] = n
	}
	return rez, r.Err()
}
