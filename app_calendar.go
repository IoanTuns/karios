package main

import (
	"errors"

	"pontix/internal/datebaza"
	"pontix/internal/i18n"
)

type RaspunsCalendar struct {
	Calendar *datebaza.CalendarAn `json:"calendar"`
	Adaugate int64                `json:"adaugate"`
	Eroare   *MesajEroare         `json:"eroare"`
}

func (a *App) CitesteCalendar(an int) RaspunsCalendar { return a.raspunsCalendar(an, 0, nil) }

func (a *App) SalveazaSarbatoare(initiala string, sarbatoare datebaza.Sarbatoare, an int) RaspunsCalendar {
	return a.raspunsCalendar(an, 0, datebaza.Nou(a.db).SalveazaSarbatoare(a.contextLucru(), initiala, sarbatoare))
}

func (a *App) CompleteazaSarbatori2026() RaspunsCalendar {
	n, err := datebaza.Nou(a.db).CompleteazaSarbatori2026(a.contextLucru())
	return a.raspunsCalendar(2026, n, err)
}

func (a *App) raspunsCalendar(an int, adaugate int64, err error) RaspunsCalendar {
	var calendar datebaza.CalendarAn
	if err == nil {
		calendar, err = datebaza.Nou(a.db).CitesteCalendar(a.contextLucru(), an)
	}
	if err != nil {
		mesaj := &MesajEroare{Cod: "date.eroare"}
		var e *i18n.Eroare
		if errors.As(err, &e) {
			mesaj.Cod, mesaj.Param = e.Cod, e.Param
		}
		return RaspunsCalendar{Eroare: mesaj}
	}
	return RaspunsCalendar{Calendar: &calendar, Adaugate: adaugate}
}
