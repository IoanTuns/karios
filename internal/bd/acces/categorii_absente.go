package acces

import "context"

func (q *Queries) CategoriiTipuriAbsenta(ctx context.Context) (map[string]string, error) {
	r, e := q.db.QueryContext(ctx, "SELECT cod,coalesce(categorie_id,categorie) FROM coduri_absenta")
	if e != nil {
		return nil, e
	}
	defer r.Close()
	rez := map[string]string{}
	for r.Next() {
		var cod, cat string
		if e := r.Scan(&cod, &cat); e != nil {
			return nil, e
		}
		rez[cod] = cat
	}
	return rez, r.Err()
}
func (q *Queries) AsociazaCategorieAbsenta(ctx context.Context, cod, cat string) error {
	_, e := q.db.ExecContext(ctx, "UPDATE coduri_absenta SET categorie_id=? WHERE cod=?", cat, cod)
	return e
}

func (q *Queries) DenumiriCategoriiAbsente(ctx context.Context) (map[string]string, error) {
	r, e := q.db.QueryContext(ctx, "SELECT t.cod,c.denumire FROM coduri_absenta t JOIN categorii_absente c ON c.cod=t.categorie_id")
	if e != nil {
		return nil, e
	}
	defer r.Close()
	rez := map[string]string{}
	for r.Next() {
		var cod, den string
		if e := r.Scan(&cod, &den); e != nil {
			return nil, e
		}
		rez[cod] = den
	}
	return rez, r.Err()
}

func (q *Queries) CategorieTipAbsentaActiva(ctx context.Context, cod string) (bool, error) {
	var activ bool
	e := q.db.QueryRowContext(ctx, "SELECT coalesce(c.activ,1) FROM coduri_absenta t LEFT JOIN categorii_absente c ON c.cod=t.categorie_id WHERE t.cod=?", cod).Scan(&activ)
	return activ, e
}
