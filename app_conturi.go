package main

import (
	"context"
	"path/filepath"

	"pontix/internal/conturi"
	"pontix/internal/datebaza"
	"pontix/internal/stocare"
)

// RaspunsStareCont spune interfeței ce ecran să arate: crearea primului administrator,
// autentificarea sau aplicația (F15, D045).
type RaspunsStareCont struct {
	NecesitaCont bool          `json:"necesita_cont"`
	Autentificat bool          `json:"autentificat"`
	Cont         *conturi.Cont `json:"cont,omitempty"`
	UltimulNume  string        `json:"ultimul_nume"`
	Intrebari    []string      `json:"intrebari"`
	Eroare       *MesajEroare  `json:"eroare,omitempty"`
}

type RaspunsIntrebari struct {
	Intrebari []string     `json:"intrebari"`
	Eroare    *MesajEroare `json:"eroare,omitempty"`
}

type RaspunsConturi struct {
	Conturi []conturi.Cont `json:"conturi"`
	Eroare  *MesajEroare   `json:"eroare,omitempty"`
}

type RaspunsOperatieCont struct {
	Succes bool         `json:"succes"`
	Eroare *MesajEroare `json:"eroare,omitempty"`
}

func (a *App) serviciuConturi() *conturi.Serviciu {
	return conturi.Nou(a.db, filepath.Join(a.radacina, stocare.NumeSeif), a.paza)
}

// contCurent întoarce contul autentificat sau nil.
func (a *App) contCurent() *conturi.Cont {
	a.contMu.Lock()
	defer a.contMu.Unlock()
	return a.cont
}

func (a *App) seteazaCont(c *conturi.Cont) {
	a.contMu.Lock()
	a.cont = c
	a.contMu.Unlock()
	a.paza.SeteazaAutentificat(c != nil)
}

// actorCurent e autorul din jurnal: numele contului, sau „local” fără conturi.
func (a *App) actorCurent() string {
	if c := a.contCurent(); c != nil {
		return c.Nume
	}
	return datebaza.ActorLocal
}

// numeAfisatCurent e numele care apare în stare_luni.modificat_de.
func (a *App) numeAfisatCurent() string {
	if c := a.contCurent(); c != nil {
		if c.NumeAfisat != "" {
			return c.NumeAfisat
		}
		return c.Nume
	}
	return "Utilizator"
}

// esteAdministrator: fără conturi (teste, baze vechi înainte de primul cont) nu se restricționează.
func (a *App) esteAdministrator() bool {
	if c := a.contCurent(); c != nil {
		return c.Administrator()
	}
	return a.paza == nil
}

func (a *App) stareCont(ctx context.Context) RaspunsStareCont {
	s := a.serviciuConturi()
	rez := RaspunsStareCont{Intrebari: conturi.Intrebari}
	n, err := s.Numar(ctx)
	if err != nil {
		rez.Eroare = eroareInterfata(err)
		return rez
	}
	rez.NecesitaCont = n == 0
	if c := a.contCurent(); c != nil {
		rez.Autentificat, rez.Cont = true, c
	}
	rez.UltimulNume = s.UltimulNume(ctx)
	return rez
}

// StareCont întoarce starea autentificării; la primul apel încearcă auto-login-ul.
func (a *App) StareCont() RaspunsStareCont {
	ctx := a.contextLucru()
	a.contMu.Lock()
	incearca := !a.autoLoginIncercat && a.cont == nil
	a.autoLoginIncercat = true
	a.contMu.Unlock()
	if incearca {
		if c, ok, err := a.serviciuConturi().AutoLogin(ctx, conturi.IdentitateCurenta()); err == nil && ok {
			a.seteazaCont(&c)
		}
	}
	return a.stareCont(ctx)
}

// CreeazaPrimulAdministrator creează contul de administrator când nu există niciun cont și
// îl autentifică (D045 punctul 4).
func (a *App) CreeazaPrimulAdministrator(cerere conturi.CerereCont) RaspunsStareCont {
	ctx := a.contextLucru()
	s := a.serviciuConturi()
	if n, err := s.Numar(ctx); err != nil || n > 0 {
		r := a.stareCont(ctx)
		r.Eroare = &MesajEroare{Cod: "cont.exista_deja"}
		return r
	}
	cerere.Rol = conturi.RolAdministrator
	if _, err := s.Creeaza(datebaza.CuActor(ctx, cerere.Nume), cerere); err != nil {
		r := a.stareCont(ctx)
		r.Eroare = eroareInterfata(err)
		return r
	}
	return a.Autentifica(cerere.Nume, cerere.Parola)
}

func (a *App) Autentifica(nume, parola string) RaspunsStareCont {
	ctx := a.contextLucru()
	c, err := a.serviciuConturi().Autentifica(ctx, nume, parola, conturi.IdentitateCurenta())
	if err != nil {
		r := a.stareCont(ctx)
		r.Eroare = eroareInterfata(err)
		return r
	}
	a.seteazaCont(&c)
	return a.stareCont(ctx)
}

// IesireDinCont oprește auto-login-ul și cere din nou parola.
func (a *App) IesireDinCont() RaspunsStareCont {
	ctx := a.contextLucru()
	if c := a.contCurent(); c != nil {
		if err := a.serviciuConturi().Iesire(ctx, *c); err != nil {
			r := a.stareCont(ctx)
			r.Eroare = eroareInterfata(err)
			return r
		}
	}
	a.seteazaCont(nil)
	return a.stareCont(ctx)
}

func (a *App) IntrebariRecuperare(nume string) RaspunsIntrebari {
	q, err := a.serviciuConturi().IntrebariPentru(a.contextLucru(), nume)
	if err != nil {
		return RaspunsIntrebari{Eroare: eroareInterfata(err)}
	}
	return RaspunsIntrebari{Intrebari: q}
}

func (a *App) RecupereazaCont(nume string, raspunsuri []string, parolaNoua string) RaspunsStareCont {
	ctx := a.contextLucru()
	c, err := a.serviciuConturi().Recupereaza(ctx, nume, raspunsuri, parolaNoua, conturi.IdentitateCurenta())
	if err != nil {
		r := a.stareCont(ctx)
		r.Eroare = eroareInterfata(err)
		return r
	}
	a.seteazaCont(&c)
	return a.stareCont(ctx)
}

func faraDrept() *MesajEroare { return &MesajEroare{Cod: "cont.doar_administrator"} }

func (a *App) ListaConturi() RaspunsConturi {
	if !a.esteAdministrator() {
		return RaspunsConturi{Eroare: faraDrept()}
	}
	l, err := a.serviciuConturi().Lista(a.contextLucru())
	if err != nil {
		return RaspunsConturi{Eroare: eroareInterfata(err)}
	}
	return RaspunsConturi{Conturi: l}
}

func (a *App) CreeazaCont(cerere conturi.CerereCont) RaspunsConturi {
	if !a.esteAdministrator() {
		return RaspunsConturi{Eroare: faraDrept()}
	}
	if _, err := a.serviciuConturi().Creeaza(a.contextLucru(), cerere); err != nil {
		return RaspunsConturi{Eroare: eroareInterfata(err)}
	}
	return a.ListaConturi()
}

func (a *App) ModificaCont(c conturi.Cont) RaspunsConturi {
	if !a.esteAdministrator() {
		return RaspunsConturi{Eroare: faraDrept()}
	}
	nou, err := a.serviciuConturi().Modifica(a.contextLucru(), c)
	if err != nil {
		return RaspunsConturi{Eroare: eroareInterfata(err)}
	}
	if cur := a.contCurent(); cur != nil && cur.ID == nou.ID {
		if !nou.Activ {
			a.seteazaCont(nil)
		} else {
			a.seteazaCont(&nou)
		}
	}
	return a.ListaConturi()
}

func (a *App) ReseteazaParolaCont(id, parolaNoua string) RaspunsConturi {
	if !a.esteAdministrator() {
		return RaspunsConturi{Eroare: faraDrept()}
	}
	if err := a.serviciuConturi().ReseteazaParola(a.contextLucru(), id, parolaNoua); err != nil {
		return RaspunsConturi{Eroare: eroareInterfata(err)}
	}
	return a.ListaConturi()
}

func (a *App) SchimbaParolaMea(veche, noua string) RaspunsOperatieCont {
	c := a.contCurent()
	if c == nil {
		return RaspunsOperatieCont{Eroare: &MesajEroare{Cod: "cont.neautentificat"}}
	}
	if err := a.serviciuConturi().SchimbaParola(a.contextLucru(), *c, veche, noua); err != nil {
		return RaspunsOperatieCont{Eroare: eroareInterfata(err)}
	}
	return RaspunsOperatieCont{Succes: true}
}

func (a *App) SchimbaIntrebarileMele(parola string, intrebari, raspunsuri []string) RaspunsOperatieCont {
	c := a.contCurent()
	if c == nil {
		return RaspunsOperatieCont{Eroare: &MesajEroare{Cod: "cont.neautentificat"}}
	}
	if err := a.serviciuConturi().SchimbaIntrebari(a.contextLucru(), *c, parola, intrebari, raspunsuri); err != nil {
		return RaspunsOperatieCont{Eroare: eroareInterfata(err)}
	}
	if nou, err := a.serviciuConturi().DupaID(a.contextLucru(), c.ID); err == nil {
		a.seteazaCont(&nou)
	}
	return RaspunsOperatieCont{Succes: true}
}
