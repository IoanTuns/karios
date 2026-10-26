package main

import "pontix/internal/datebaza"

type RaspunsPlanificare struct {
	Planificare *datebaza.PlanificareLuna `json:"planificare"`
	Eroare      *MesajEroare              `json:"eroare"`
}
type RaspunsPublicarePlan struct {
	Publicare *datebaza.PublicarePlan `json:"publicare"`
	Eroare    *MesajEroare            `json:"eroare"`
}

func (a *App) CitestePlanificare(an, luna int) RaspunsPlanificare {
	p, err := datebaza.Nou(a.db).CitestePlanificare(a.contextLucru(), an, luna)
	if err != nil {
		return RaspunsPlanificare{Eroare: mesajPontaj(err)}
	}
	return RaspunsPlanificare{Planificare: &p}
}
func (a *App) SalveazaPlanificare(an, luna int, pozitii []datebaza.PozitiePlanificare) RaspunsPlanificare {
	if err := datebaza.Nou(a.db).SalveazaPlanificare(a.contextLucru(), an, luna, pozitii); err != nil {
		return RaspunsPlanificare{Eroare: mesajPontaj(err)}
	}
	return a.CitestePlanificare(an, luna)
}
func (a *App) PrevizualizeazaPublicarePlan(c datebaza.CererePublicarePlan) RaspunsPublicarePlan {
	p, err := datebaza.Nou(a.db).PrevizualizeazaPublicarePlan(a.contextLucru(), c)
	if err != nil {
		return RaspunsPublicarePlan{Eroare: mesajPontaj(err)}
	}
	return RaspunsPublicarePlan{Publicare: &p}
}
func (a *App) PublicaPlanificare(c datebaza.CererePublicarePlan) RaspunsPublicarePlan {
	p, err := datebaza.Nou(a.db).PublicaPlanificare(a.contextLucru(), c)
	if err != nil {
		return RaspunsPublicarePlan{Eroare: mesajPontaj(err)}
	}
	return RaspunsPublicarePlan{Publicare: &p}
}
