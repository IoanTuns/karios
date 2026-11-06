package main

import "pontix/internal/datebaza"

// StergeConfigurare cere administrator autentificat, fara exceptii pentru baze fara conturi.
func (a *App) StergeConfigurare(tip, id, revizie string) RaspunsDate {
	cont := a.contCurent()
	if cont == nil || !cont.Administrator() {
		return RaspunsDate{Eroare: faraDrept()}
	}
	actual, err := a.serviciuConturi().DupaID(a.contextLucru(), cont.ID)
	if err != nil || !actual.Activ || !actual.Administrator() {
		return RaspunsDate{Eroare: faraDrept()}
	}
	return a.raspunsDate(datebaza.Nou(a.db).StergeConfigurare(a.contextLucru(), tip, id, revizie))
}
