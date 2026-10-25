package main

import (
	"errors"

	"pontix/internal/datebaza"
	"pontix/internal/i18n"
)

type RaspunsSchita struct {
	Schita *datebaza.SchitaLuna `json:"schita"`
	Eroare *MesajEroare         `json:"eroare"`
}

func (a *App) CitesteSchita(an, luna int) RaspunsSchita { return a.raspunsSchita(an, luna, nil) }

func (a *App) SalveazaSchita(an, luna int, celule []datebaza.CelulaSchita) RaspunsSchita {
	return a.raspunsSchita(an, luna, datebaza.Nou(a.db).SalveazaSchita(a.contextLucru(), an, luna, celule))
}

func (a *App) raspunsSchita(an, luna int, err error) RaspunsSchita {
	var schita datebaza.SchitaLuna
	if err == nil {
		schita, err = datebaza.Nou(a.db).CitesteSchita(a.contextLucru(), an, luna)
	}
	if err != nil {
		mesaj := &MesajEroare{Cod: "date.eroare"}
		var e *i18n.Eroare
		if errors.As(err, &e) {
			mesaj.Cod, mesaj.Param = e.Cod, e.Param
		}
		return RaspunsSchita{Eroare: mesaj}
	}
	return RaspunsSchita{Schita: &schita}
}
