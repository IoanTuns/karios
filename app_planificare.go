package main

import (
	"pontix/internal/calcul"
	"pontix/internal/datebaza"
	"pontix/internal/diagnoza"
)

type RaspunsPlanificare struct {
	Planificare *datebaza.PlanificareLuna `json:"planificare"`
	Eroare      *MesajEroare              `json:"eroare"`
}
type RaspunsPublicarePlan struct {
	Publicare *datebaza.PublicarePlan `json:"publicare"`
	Eroare    *MesajEroare            `json:"eroare"`
}
type RaspunsAcoperirePlan struct {
	Acoperire *calcul.RaportAcoperire `json:"acoperire"`
	Eroare    *MesajEroare            `json:"eroare"`
}
type RaspunsEchitatePlan struct {
	Avertismente []calcul.AvertismentEchitate `json:"avertismente"`
	Eroare       *MesajEroare                 `json:"eroare"`
}

type RaspunsPlanificatVsRealizat struct {
	Raport *calcul.RaportPlanificatVsRealizat `json:"raport"`
	Eroare *MesajEroare                       `json:"eroare"`
}

func (a *App) CitestePlanificare(an, luna int) RaspunsPlanificare {
	stop := diagnoza.Cronometreaza("Planificare.Citeste")
	var err error
	defer stop(&err)

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
func (a *App) CitesteAcoperirePlanificare(an, luna int, posturi []string) RaspunsAcoperirePlan {
	rap, err := datebaza.Nou(a.db).CitesteAcoperirePlanificare(a.contextLucru(), an, luna, posturi)
	if err != nil {
		return RaspunsAcoperirePlan{Eroare: mesajPontaj(err)}
	}
	return RaspunsAcoperirePlan{Acoperire: &rap}
}
func (a *App) CitesteEchitatePlanificare(an, luna int) RaspunsEchitatePlan {
	av, err := datebaza.Nou(a.db).CitesteEchitatePlanificare(a.contextLucru(), an, luna)
	if err != nil {
		return RaspunsEchitatePlan{Eroare: mesajPontaj(err)}
	}
	return RaspunsEchitatePlan{Avertismente: av}
}

func (a *App) CitestePlanificatVsRealizat(an, luna int) RaspunsPlanificatVsRealizat {
	rap, err := datebaza.Nou(a.db).CitestePlanificatVsRealizat(a.contextLucru(), an, luna)
	if err != nil {
		return RaspunsPlanificatVsRealizat{Eroare: mesajPontaj(err)}
	}
	return RaspunsPlanificatVsRealizat{Raport: &rap}
}
