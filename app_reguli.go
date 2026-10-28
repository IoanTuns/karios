package main

import (
	"pontix/internal/calcul"
	"pontix/internal/datebaza"
)

type RaspunsListaReguli struct {
	Reguli []calcul.RegulaConfigurabila `json:"reguli"`
	Eroare *MesajEroare                 `json:"eroare,omitempty"`
}

type RaspunsSalvareRegula struct {
	Regula *calcul.RegulaConfigurabila `json:"regula,omitempty"`
	Eroare *MesajEroare                `json:"eroare,omitempty"`
}

type RaspunsVerificaReguli struct {
	Incalcari []calcul.IncalcareRegula `json:"incalcari"`
	Eroare    *MesajEroare             `json:"eroare,omitempty"`
}

func (a *App) ListaReguliConfigurabile() RaspunsListaReguli {
	reg, err := datebaza.Nou(a.db).ListaReguliConfigurabile(a.contextLucru())
	if err != nil {
		return RaspunsListaReguli{Eroare: eroareInterfata(err)}
	}
	return RaspunsListaReguli{Reguli: reg}
}

func (a *App) SalveazaRegulaConfigurabila(r calcul.RegulaConfigurabila) RaspunsSalvareRegula {
	srv := datebaza.Nou(a.db)
	if err := srv.SalveazaRegulaConfigurabila(a.contextLucru(), r); err != nil {
		return RaspunsSalvareRegula{Eroare: eroareInterfata(err)}
	}
	return RaspunsSalvareRegula{Regula: &r}
}

func (a *App) VerificaReguliPlanificare(an, luna int) RaspunsVerificaReguli {
	inc, err := datebaza.Nou(a.db).VerificaReguliPlanificare(a.contextLucru(), an, luna)
	if err != nil {
		return RaspunsVerificaReguli{Eroare: eroareInterfata(err)}
	}
	return RaspunsVerificaReguli{Incalcari: inc}
}
