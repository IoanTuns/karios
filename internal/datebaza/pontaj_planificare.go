package datebaza

import (
	"context"
	"strings"

	"pontix/internal/bd/acces"
	"pontix/internal/calcul"
)

type schimbarePontajPlan struct {
	Celula  CelulaSchita
	Inainte string
}

// Sincronizarea folosește aceeași tranzacție cu schița. Validarea se face
// după întregul lot, astfel încât programele înlocuite să nu creeze conflicte false.
func sincronizeazaPontajPlan(ctx context.Context, q *acces.Queries, deLa, panaLa string, modificari []schimbarePontajPlan) error {
	surse, err := sursePentruPlanificare(ctx, q, deLa, panaLa)
	if err != nil {
		return err
	}
	active := map[string]acces.Planificare{}
	for _, p := range surse.pozitii {
		if p.Activ == 1 {
			active[p.AngajatID+"/"+p.Data] = p
		}
	}
	verifica := map[string]bool{}
	for _, m := range modificari {
		c := m.Celula
		k := c.AngajatID + "/" + c.Data
		program, cunoscut := surse.coduri[strings.ToUpper(c.Valoare)]
		special := cunoscut && calcul.ProgramDePlanificat(program)
		v, exista := active[k]
		codVechi := surse.schimburi[v.SchimbID].Cod
		corespunde := exista && strings.EqualFold(strings.TrimSpace(m.Inainte), codVechi)
		if !special {
			if !corespunde {
				continue
			}
			if n, err := q.ModificaPozitiePlanificare(ctx, acces.ModificaPozitiePlanificareParams{ID: v.ID, Revizie: v.Revizie, SchimbID: v.SchimbID, Activ: 0, ModificatLa: moment()}); err != nil {
				return err
			} else if n != 1 {
				return eroare("planificare.conflict")
			}
			nou := v
			nou.Activ = 0
			nou.Revizie++
			nou.ModificatLa = moment()
			if err := jurnal(ctx, q, "planificare", v.ID, v, nou); err != nil {
				return err
			}
			continue
		}
		if exista && v.SchimbID != program.ID && !corespunde {
			return eroare("pontaj.plan_independent")
		}
		p := v
		if !exista || v.SchimbID != program.ID {
			if exista {
				if n, err := q.ModificaPozitiePlanificare(ctx, acces.ModificaPozitiePlanificareParams{ID: v.ID, Revizie: v.Revizie, SchimbID: v.SchimbID, Activ: 0, ModificatLa: moment()}); err != nil {
					return err
				} else if n != 1 {
					return eroare("planificare.conflict")
				}
				nou := v
				nou.Activ = 0
				nou.Revizie++
				nou.ModificatLa = moment()
				if err := jurnal(ctx, q, "planificare", v.ID, v, nou); err != nil {
					return err
				}
			}
			sectiune := "GARZI"
			if calcul.ProgramInSectiune(program, "DUPA_AMIAZA") {
				sectiune = "DUPA_AMIAZA"
			}
			p = acces.Planificare{ID: idNou(), AngajatID: c.AngajatID, Data: c.Data, SchimbID: program.ID, Sectiune: sectiune, Activ: 1, Revizie: 1, CreatLa: moment(), ModificatLa: moment()}
			if err := q.AdaugaPozitiePlanificare(ctx, acces.AdaugaPozitiePlanificareParams{ID: p.ID, AngajatID: p.AngajatID, Data: p.Data, SchimbID: p.SchimbID, Sectiune: p.Sectiune, Activ: 1, CreatLa: p.CreatLa, ModificatLa: p.ModificatLa}); err != nil {
				return err
			}
			if err := jurnal(ctx, q, "planificare", p.ID, nil, p); err != nil {
				return err
			}
		}
		verifica[p.ID] = true
		ciorna := surse.ciorne[k]
		if p.PublicatPlanRevizie != p.Revizie || p.PublicatPontajRevizie != ciorna.Revizie || p.PublicatValoare != program.Cod {
			nou := p
			nou.PublicatPlanRevizie = p.Revizie
			nou.PublicatPontajRevizie = ciorna.Revizie
			nou.PublicatValoare = program.Cod
			nou.PublicatLa = moment()
			n, err := q.MarcheazaPlanificarePublicata(ctx, acces.MarcheazaPlanificarePublicataParams{ID: p.ID, Revizie: p.Revizie, PublicatPlanRevizie: nou.PublicatPlanRevizie, PublicatPontajRevizie: nou.PublicatPontajRevizie, PublicatValoare: nou.PublicatValoare, PublicatLa: nou.PublicatLa})
			if err != nil {
				return err
			}
			if n != 1 {
				return eroare("planificare.conflict")
			}
			if err := jurnal(ctx, q, "planificare", p.ID, p, nou); err != nil {
				return err
			}
		}
	}
	finale, err := sursePentruPlanificare(ctx, q, deLa, panaLa)
	if err != nil {
		return err
	}
	for _, p := range finale.pozitii {
		if verifica[p.ID] {
			if motiv := finale.verifica(p, true, true); motiv != "" {
				return eroare("planificare.motiv_" + motiv)
			}
		}
	}
	return nil
}
