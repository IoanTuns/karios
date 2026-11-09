package main

import "pontix/internal/datebaza"

// RaspunsIstoricContract: schimbarile normei si ale dreptului la concediu ale unui angajat (F07-T43, D101).
type RaspunsIstoricContract struct {
	Schimbari []datebaza.SchimbareContract `json:"schimbari"`
	Eroare    *MesajEroare                 `json:"eroare"`
}

func (a *App) istoricContract(angajatID string, err error) RaspunsIstoricContract {
	if err != nil {
		return RaspunsIstoricContract{Eroare: eroareInterfata(err)}
	}
	lista, err := datebaza.Nou(a.db).ListaIstoricContract(a.contextLucru(), angajatID)
	if err != nil {
		return RaspunsIstoricContract{Eroare: eroareInterfata(err)}
	}
	return RaspunsIstoricContract{Schimbari: lista}
}

// ListaIstoricContract intoarce schimbarile unui angajat, in ordinea datei.
func (a *App) ListaIstoricContract(angajatID string) RaspunsIstoricContract {
	return a.istoricContract(angajatID, nil)
}

// AdaugaSchimbareContract adauga o schimbare si intoarce lista actualizata.
func (a *App) AdaugaSchimbareContract(s datebaza.SchimbareContract) RaspunsIstoricContract {
	return a.istoricContract(s.AngajatID, datebaza.Nou(a.db).AdaugaSchimbareContract(a.contextLucru(), s))
}

// StergeSchimbareContract sterge schimbarea de la data deLa si intoarce lista actualizata.
func (a *App) StergeSchimbareContract(angajatID, deLa string) RaspunsIstoricContract {
	return a.istoricContract(angajatID, datebaza.Nou(a.db).StergeSchimbareContract(a.contextLucru(), angajatID, deLa))
}
