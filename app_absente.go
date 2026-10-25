package main

import (
	"errors"
	"pontix/internal/datebaza"
	"pontix/internal/i18n"
)

type RaspunsAbsente struct {
	Absente *datebaza.AbsenteAn `json:"absente"`
	Eroare  *MesajEroare        `json:"eroare"`
}

func (a *App) CitesteAbsente(an int) RaspunsAbsente { return a.raspunsAbsente(an, nil) }

func (a *App) SalveazaTipAbsenta(initial string, tip datebaza.TipAbsenta, an int) RaspunsAbsente {
	return a.raspunsAbsente(an, datebaza.Nou(a.db).SalveazaTipAbsenta(a.contextLucru(), initial, tip))
}

func (a *App) SalveazaAbsente(zile []datebaza.ZiAbsenta, an int) RaspunsAbsente {
	return a.raspunsAbsente(an, datebaza.Nou(a.db).SalveazaAbsente(a.contextLucru(), zile))
}

func (a *App) raspunsAbsente(an int, err error) RaspunsAbsente {
	var absente datebaza.AbsenteAn
	if err == nil {
		absente, err = datebaza.Nou(a.db).CitesteAbsente(a.contextLucru(), an)
	}
	if err != nil {
		mesaj := &MesajEroare{Cod: "date.eroare"}
		var e *i18n.Eroare
		if errors.As(err, &e) {
			mesaj.Cod, mesaj.Param = e.Cod, e.Param
		}
		return RaspunsAbsente{Eroare: mesaj}
	}
	return RaspunsAbsente{Absente: &absente}
}
