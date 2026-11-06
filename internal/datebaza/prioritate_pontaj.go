package datebaza

import (
	"context"
	"pontix/internal/bd/acces"
	"pontix/internal/calcul"
	"time"
)

func retragePontajPrioritar(ctx context.Context, q *acces.Queries, p acces.PontajPrioritate) error {
	d, err := time.Parse(time.DateOnly, p.Data)
	if err != nil {
		return err
	}
	deschisa, err := q.LunaDeschisaPrioritate(ctx, d.Year(), int(d.Month()))
	if err != nil {
		return err
	}
	if !deschisa {
		return eroare("pontaj.luna_inchisa")
	}
	// Se verifica toate lunile intersectate, nu numai data de inceput a turei.
	a, e1 := time.Parse(time.RFC3339, p.Start)
	b, e2 := time.Parse(time.RFC3339, p.Sfarsit)
	if e1 == nil && e2 == nil && b.After(a) {
		startLocal := a.In(calcul.Locatie)
		for z := time.Date(startLocal.Year(), startLocal.Month(), 1, 0, 0, 0, 0, calcul.Locatie); z.Before(b); z = z.AddDate(0, 1, 0) {
			ok, e := q.LunaDeschisaPrioritate(ctx, z.Year(), int(z.Month()))
			if e != nil {
				return e
			}
			if !ok {
				return eroare("pontaj.luna_inchisa")
			}
		}
		z := b.Add(-time.Nanosecond).In(calcul.Locatie)
		ok, e := q.LunaDeschisaPrioritate(ctx, z.Year(), int(z.Month()))
		if e != nil {
			return e
		}
		if !ok {
			return eroare("pontaj.luna_inchisa")
		}
	}
	acum := moment()
	n, err := q.RetragePontajPrioritate(ctx, p, acum)
	if err != nil {
		return err
	}
	if n != 1 {
		return eroare("pontaj.conflict")
	}
	nou := map[string]any{}
	for k, v := range p.Instantaneu {
		nou[k] = v
	}
	nou["sters_la"] = acum
	nou["modificat_la"] = acum
	return jurnal(ctx, q, "pontaj", p.ID, p.Instantaneu, nou)
}
func pontajIntersecteazaZi(p acces.PontajPrioritate, zi string) bool {
	if p.Data == zi {
		return true
	}
	d, e := time.ParseInLocation(time.DateOnly, zi, calcul.Locatie)
	if e != nil {
		return false
	}
	a, e1 := time.Parse(time.RFC3339, p.Start)
	b, e2 := time.Parse(time.RFC3339, p.Sfarsit)
	return e1 == nil && e2 == nil && a.Before(d.AddDate(0, 0, 1)) && b.After(d)
}
