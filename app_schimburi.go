package main

import (
	"errors"

	"pontix/internal/datebaza"
	"pontix/internal/i18n"
)

type RaspunsSchimburi struct {
	Adaugate  int64             `json:"adaugate"`
	Schimburi []datebaza.Schimb `json:"schimburi"`
	Eroare    *MesajEroare      `json:"eroare"`
}

func (a *App) ListaSchimburi() RaspunsSchimburi { return a.raspunsSchimburi(nil) }

func (a *App) CompleteazaSchimburiInitiale() RaspunsSchimburi {
	n, err := datebaza.Nou(a.db).CompleteazaSchimburiInitiale(a.contextLucru())
	r := a.raspunsSchimburi(err)
	r.Adaugate = n
	return r
}

func (a *App) SalveazaSchimb(schimb datebaza.Schimb) RaspunsSchimburi {
	return a.raspunsSchimburi(datebaza.Nou(a.db).SalveazaSchimb(a.contextLucru(), schimb))
}

func (a *App) SchimbaActivSchimb(id string, activ bool) RaspunsSchimburi {
	return a.raspunsSchimburi(datebaza.Nou(a.db).SchimbaActivSchimb(a.contextLucru(), id, activ))
}

func (a *App) raspunsSchimburi(err error) RaspunsSchimburi {
	var schimburi []datebaza.Schimb
	if err == nil {
		schimburi, err = datebaza.Nou(a.db).ListaSchimburi(a.contextLucru())
	}
	if err != nil {
		mesaj := &MesajEroare{Cod: "date.eroare"}
		var e *i18n.Eroare
		if errors.As(err, &e) {
			mesaj.Cod, mesaj.Param = e.Cod, e.Param
		}
		return RaspunsSchimburi{Eroare: mesaj}
	}
	return RaspunsSchimburi{Schimburi: schimburi}
}
