package datebaza

import (
	"context"
	"pontix/internal/bd/acces"
)

// zileMedicalFamilie extinde protectia programelor, fara a duplica absentele.
func zileMedicalFamilie(ctx context.Context, q *acces.Queries, deLa, panaLa string) (map[string]bool, error) {
	legaturi, err := q.ListaLegaturiIerarhice(ctx)
	if err != nil {
		return nil, err
	}
	tipuri, err := q.ListaTipuriAbsenta(ctx)
	if err != nil {
		return nil, err
	}
	medical := map[string]bool{}
	for _, tip := range tipuri {
		medical[tip.Cod] = tip.Categorie == "MEDICAL"
	}
	zile, err := q.ListaAbsenteInterval(ctx, acces.ListaAbsenteIntervalParams{DeLa: deLa, PanaLa: panaLa})
	if err != nil {
		return nil, err
	}
	principal := func(id string) string {
		if p := legaturi[id]; p != "" {
			return p
		}
		return id
	}
	rez := map[string]bool{}
	for _, zi := range zile {
		if zi.Activ != 1 || zi.StareAprobare == "RESPINSA" || !medical[zi.Cod] {
			continue
		}
		p := principal(zi.AngajatID)
		rez[p+"/"+zi.Data] = true
		rez[zi.AngajatID+"/"+zi.Data] = true
		for copil := range legaturi {
			if principal(copil) == p {
				rez[copil+"/"+zi.Data] = true
			}
		}
	}
	return rez, nil
}
