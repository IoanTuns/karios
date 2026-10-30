package main

import (
	"context"
	"fmt"

	"pontix/internal/audit"
	"pontix/internal/conturi"
	"pontix/internal/diagnoza"
	"pontix/internal/stocare"
)

// RaspunsStareSeif spune interfeței dacă datele trebuie deblocate înainte de orice (D054).
type RaspunsStareSeif struct {
	Blocat bool `json:"blocat"`
}

// RaspunsRecuperare: cheia de recuperare încă neconfirmată (D054 punctul 3). Codul se trimite
// doar administratorului autentificat.
type RaspunsRecuperare struct {
	Necesara bool         `json:"necesara"`
	Cod      string       `json:"cod,omitempty"`
	Eroare   *MesajEroare `json:"eroare,omitempty"`
}

// aplicaDeschiderea preia baza, paza și seiful deschise la pornire sau la deblocare.
func (a *App) aplicaDeschiderea(d *stocare.Deschidere) {
	a.seif = d.Seif
	if d.Blocat {
		return
	}
	a.db, a.paza, a.codRecuperare = d.DB, d.Paza, d.CodRecuperare
	// Până la autentificare (sau crearea primului cont), scrierile sunt refuzate (D045).
	a.paza.CereAutentificare(true)
	if r := a.paza.Rezultat(); r.Stare == audit.StareIncident {
		// Doar numele tabelelor și id-ul rândului din jurnal, fără date (D043).
		diagnoza.ObtineManager().InregistreazaEroare("pornire/audit", fmt.Errorf("incident: tabele %v, rând rupt %d", r.Tabele, r.RandRupt))
	}
}

// StareSeif: datele sunt blocate când seiful nu are copia W a utilizatorului Windows curent.
func (a *App) StareSeif() RaspunsStareSeif {
	return RaspunsStareSeif{Blocat: a.db == nil}
}

func (a *App) deblocheaza(ctx context.Context) error {
	d, err := stocare.Deblocheaza(ctx, a.radacina, a.seif)
	if err != nil {
		return err
	}
	a.aplicaDeschiderea(d)
	return nil
}

// DeblocheazaCuParola deschide datele cu parola unui cont și autentifică acel cont.
func (a *App) DeblocheazaCuParola(nume, parola string) RaspunsStareCont {
	ctx := a.contextLucru()
	if a.db != nil {
		return a.Autentifica(nume, parola)
	}
	if _, err := a.seif.DeschideCuParola(parola); err != nil {
		return RaspunsStareCont{Eroare: eroareInterfata(err)}
	}
	if err := a.deblocheaza(ctx); err != nil {
		return RaspunsStareCont{Eroare: eroareInterfata(err)}
	}
	return a.Autentifica(nume, parola)
}

// DeblocheazaCuRecuperare deschide datele cu cheia de recuperare; urmează autentificarea obișnuită.
func (a *App) DeblocheazaCuRecuperare(cod string) RaspunsStareCont {
	if a.db == nil {
		if err := a.seif.DeschideCuRecuperare(cod); err != nil {
			return RaspunsStareCont{Eroare: eroareInterfata(err)}
		}
		if err := a.deblocheaza(a.contextLucru()); err != nil {
			return RaspunsStareCont{Eroare: eroareInterfata(err)}
		}
	}
	return a.stareCont(a.contextLucru())
}

// StareRecuperare spune dacă cheia de recuperare așteaptă confirmarea administratorului.
func (a *App) StareRecuperare() RaspunsRecuperare {
	if a.seif == nil || a.db == nil || a.seif.RecuperareConfirmata() {
		return RaspunsRecuperare{}
	}
	r := RaspunsRecuperare{Necesara: true}
	if c := a.contCurent(); c != nil && c.Rol == conturi.RolAdministrator {
		r.Cod = a.codRecuperare
	}
	return r
}

// ConfirmaCheiaRecuperare: administratorul retastează cheia de recuperare (D054 punctul 3).
func (a *App) ConfirmaCheiaRecuperare(cod string) RaspunsRecuperare {
	if c := a.contCurent(); c == nil || c.Rol != conturi.RolAdministrator {
		return RaspunsRecuperare{Necesara: true, Eroare: faraDrept()}
	}
	if err := a.seif.ConfirmaRecuperare(cod); err != nil {
		r := a.StareRecuperare()
		r.Eroare = eroareInterfata(err)
		return r
	}
	a.codRecuperare = ""
	return RaspunsRecuperare{}
}

// GenereazaCheieRecuperareNoua: administratorul poate înlocui oricând cheia de recuperare
// (de ex. dacă hârtia s-a pierdut); cea veche nu mai funcționează, iar cea nouă trebuie confirmată.
func (a *App) GenereazaCheieRecuperareNoua() RaspunsRecuperare {
	if c := a.contCurent(); c == nil || c.Rol != conturi.RolAdministrator {
		return RaspunsRecuperare{Eroare: faraDrept()}
	}
	cod, err := a.seif.SeteazaRecuperare()
	if err != nil {
		return RaspunsRecuperare{Eroare: eroareInterfata(err)}
	}
	a.codRecuperare = cod
	return a.StareRecuperare()
}
