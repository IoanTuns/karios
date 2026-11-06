package datebaza

import (
	"context"
	"pontix/internal/bd/acces"
	"time"
)

func initializarePentru(m map[string]acces.InitializarePontaj, id string) *acces.InitializarePontaj {
	r, ok := m[id]
	if !ok {
		return nil
	}
	return &r
}
func salveazaInitializare(ctx context.Context, q *acces.Queries, id string, r acces.InitializarePontaj) error {
	d, e := time.Parse("2006-01", r.Luna)
	if e != nil || d.Year() < 2000 || d.Year() > 2100 || r.COZile < 0 || r.COZile > 366 || r.Minute < -60000000 || r.Minute > 60000000 {
		return eroare("sold.initializare_invalida")
	}
	toate, e := q.ListaInitializariPontaj(ctx)
	if e != nil {
		return e
	}
	v, exista := toate[id]
	if exista && v == r {
		return nil
	}
	if (exista && v.Revizie != r.Revizie) || (!exista && r.Revizie != 0) {
		return eroare("pontaj.conflict")
	}
	for _, l := range []string{r.Luna, v.Luna} {
		if l == "" {
			continue
		}
		d, _ := time.Parse("2006-01", l)
		ok, e := q.LunaDeschisaPrioritate(ctx, d.Year(), int(d.Month()))
		if e != nil {
			return e
		}
		if !ok {
			return eroare("pontaj.luna_inchisa")
		}
	}
	prima := r.Luna
	if v.Luna != "" && v.Luna < prima {
		prima = v.Luna
	}
	inchisa, e := q.SoldAfecteazaLunaInchisa(ctx, prima)
	if e != nil {
		return e
	}
	if inchisa {
		return eroare("pontaj.luna_inchisa")
	}
	// CO initial nu dubleaza absentele deja prezente inaintea referintei.
	if r.COZile > 0 {
		zile, e := q.ListaAbsenteInterval(ctx, acces.ListaAbsenteIntervalParams{DeLa: r.Luna[:4] + "-01-01", PanaLa: d.AddDate(0, 0, -1).Format(time.DateOnly)})
		if e != nil {
			return e
		}
		for _, z := range zile {
			if z.AngajatID == id && z.Activ == 1 {
				tip, e := q.CitesteTipAbsenta(ctx, z.Cod)
				if e != nil {
					return e
				}
				if tip.Categorie != "ODIHNA" {
					continue
				}
				return eroare("sold.co_dublat")
			}
		}
	}
	n, e := q.SalveazaInitializarePontaj(ctx, id, r, moment())
	if e != nil {
		return e
	}
	if n != 1 {
		return eroare("pontaj.conflict")
	}
	r.Revizie++
	return jurnal(ctx, q, "initializari_pontaj", id, v, r)
}
