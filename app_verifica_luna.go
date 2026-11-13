package main

import "pontix/internal/datebaza"

type RaspunsVerificareLuna struct {
	Verificare *datebaza.VerificareLuna `json:"verificare,omitempty"`
	Eroare     *MesajEroare             `json:"eroare,omitempty"`
}

// VerificaLuna adună constatările lunii înainte de închidere (F12-T06); nu schimbă nimic.
func (a *App) VerificaLuna(an, luna int) RaspunsVerificareLuna {
	if a.db == nil || an < 1 || an > 9999 || luna < 1 || luna > 12 {
		return RaspunsVerificareLuna{Eroare: &MesajEroare{Cod: "date.eroare"}}
	}
	v, err := datebaza.Nou(a.db).VerificaLuna(a.contextLucru(), an, luna)
	if err != nil {
		return RaspunsVerificareLuna{Eroare: eroareInterfata(err)}
	}
	return RaspunsVerificareLuna{Verificare: &v}
}
