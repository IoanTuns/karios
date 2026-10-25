package main

import (
	"context"
	"errors"

	"pontix/internal/datebaza"
	"pontix/internal/i18n"
)

// MesajEroare păstrează codul și parametrii pentru traducerea în React.
type MesajEroare struct {
	Cod   string            `json:"cod"`
	Param map[string]string `json:"param"`
}

type RaspunsDate struct {
	Date   *datebaza.Date `json:"date"`
	Eroare *MesajEroare   `json:"eroare"`
}

func (a *App) DateInitiale() RaspunsDate { return a.raspunsDate(nil) }

func (a *App) SalveazaAngajat(angajat datebaza.Angajat) RaspunsDate {
	return a.raspunsDate(datebaza.Nou(a.db).SalveazaAngajat(a.contextLucru(), angajat))
}

func (a *App) SchimbaActivAngajat(id string, activ bool) RaspunsDate {
	return a.raspunsDate(datebaza.Nou(a.db).SchimbaActiv(a.contextLucru(), id, activ))
}

func (a *App) SalveazaNormaIntreaga(minute int64) RaspunsDate {
	return a.raspunsDate(datebaza.Nou(a.db).SalveazaNorma(a.contextLucru(), minute))
}

func (a *App) AdaugaSectie(denumire i18n.Text) RaspunsDate {
	return a.raspunsDate(datebaza.Nou(a.db).AdaugaSectie(a.contextLucru(), denumire))
}

func (a *App) AdaugaDepartament(sectieID string, denumire i18n.Text) RaspunsDate {
	return a.raspunsDate(datebaza.Nou(a.db).AdaugaDepartament(a.contextLucru(), sectieID, denumire))
}

func (a *App) IncarcaDemo() RaspunsDate {
	return a.raspunsDate(datebaza.Nou(a.db).IncarcaDemo(a.contextLucru()))
}

func (a *App) contextLucru() context.Context {
	if a.ctx != nil {
		return a.ctx
	}
	return context.Background()
}

func (a *App) raspunsDate(err error) RaspunsDate {
	var date datebaza.Date
	if err == nil {
		date, err = datebaza.Nou(a.db).Citeste(a.contextLucru())
	}
	if err != nil {
		mesaj := &MesajEroare{Cod: "date.eroare"}
		var e *i18n.Eroare
		if errors.As(err, &e) {
			mesaj.Cod, mesaj.Param = e.Cod, e.Param
		}
		return RaspunsDate{Eroare: mesaj}
	}
	return RaspunsDate{Date: &date}
}
