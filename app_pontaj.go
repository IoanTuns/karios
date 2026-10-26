package main

import (
	"errors"

	"pontix/internal/calcul"
	"pontix/internal/datebaza"
	"pontix/internal/i18n"
)

type RaspunsSchita struct {
	Schita *datebaza.SchitaLuna `json:"schita"`
	Eroare *MesajEroare         `json:"eroare"`
}

type RaspunsPrecompletare struct {
	Plan   *calcul.PlanPrecompletare `json:"plan"`
	Eroare *MesajEroare              `json:"eroare"`
}

func (a *App) PrevizualizeazaPrecompletare(an, luna int, angajati []string, modificari []datebaza.CelulaSchita) RaspunsPrecompletare {
	p, err := datebaza.Nou(a.db).PrevizualizeazaPrecompletare(a.contextLucru(), an, luna, angajati, modificari)
	if err != nil {
		return RaspunsPrecompletare{Eroare: mesajPontaj(err)}
	}
	return RaspunsPrecompletare{Plan: &p}
}

func (a *App) CitesteSchita(an, luna int) RaspunsSchita { return a.raspunsSchita(an, luna, nil) }

func (a *App) PrevizualizeazaSchita(an, luna int, celule []datebaza.CelulaSchita) RaspunsSchita {
	s, err := datebaza.Nou(a.db).CitesteSchita(a.contextLucru(), an, luna, celule)
	if err != nil {
		return RaspunsSchita{Eroare: mesajPontaj(err)}
	}
	return RaspunsSchita{Schita: &s}
}
func mesajPontaj(err error) *MesajEroare {
	var e *i18n.Eroare
	if errors.As(err, &e) {
		return &MesajEroare{Cod: e.Cod, Param: e.Param}
	}
	return &MesajEroare{Cod: "date.eroare"}
}

type RaspunsConfiguratiePontaj struct {
	Configuratie *calcul.ConfiguratiePontaj `json:"configuratie"`
	Optiuni      []datebaza.OptiuneRegula   `json:"optiuni"`
	Eroare       *MesajEroare               `json:"eroare"`
}

func (a *App) CitesteConfiguratiePontaj() RaspunsConfiguratiePontaj {
	c, o, err := datebaza.Nou(a.db).ConfiguratiePontaj(a.contextLucru())
	if err != nil {
		return RaspunsConfiguratiePontaj{Eroare: mesajPontaj(err)}
	}
	return RaspunsConfiguratiePontaj{Configuratie: &c, Optiuni: o}
}
func (a *App) SalveazaConfiguratiePontaj(c calcul.ConfiguratiePontaj) RaspunsConfiguratiePontaj {
	if err := datebaza.Nou(a.db).SalveazaConfiguratiePontaj(a.contextLucru(), c); err != nil {
		return RaspunsConfiguratiePontaj{Eroare: mesajPontaj(err)}
	}
	return a.CitesteConfiguratiePontaj()
}

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
