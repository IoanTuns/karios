package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"pontix/internal/conturi"
	"pontix/internal/datebaza"
	"pontix/internal/diagnoza"
	"pontix/internal/i18n"
)

// F03-T11 (D004, D106): traducerile din folderul karios/traduceri se importă în baza criptată.

const (
	numeFolderTraduceri    = "traduceri"
	maxFisiereTraduceri    = 20
	maxOctetiFisierTraduce = 5 << 20
)

type RaspunsTraduceri struct {
	Folder    string                      `json:"folder"`
	Corecturi map[string]int              `json:"corecturi"` // limbă → câte texte corectate
	Rezumate  []datebaza.RezumatTraduceri `json:"rezumate,omitempty"`
	Cale      string                      `json:"cale,omitempty"`
	Eroare    *MesajEroare                `json:"eroare,omitempty"`
}

func (a *App) folderTraduceri() string { return filepath.Join(a.radacina, numeFolderTraduceri) }

// incarcaTraduceri pune corecturile din bază peste cataloagele incluse (la deschiderea bazei
// și după fiecare import sau ștergere).
func (a *App) incarcaTraduceri() {
	if a.db == nil {
		i18n.SeteazaCorecturi(nil)
		return
	}
	c, err := datebaza.Nou(a.db).ListaTraduceri(a.contextLucru())
	if err != nil {
		diagnoza.ObtineManager().InregistreazaEroare("traduceri/incarcare", err)
		return
	}
	i18n.SeteazaCorecturi(c)
}

func (a *App) stareTraduceri() RaspunsTraduceri {
	r := RaspunsTraduceri{Folder: a.folderTraduceri(), Corecturi: map[string]int{}}
	if a.db == nil {
		return r
	}
	c, err := datebaza.Nou(a.db).ListaTraduceri(a.contextLucru())
	if err != nil {
		r.Eroare = eroareInterfata(err)
		return r
	}
	for l, texte := range c {
		r.Corecturi[string(l)] = len(texte)
	}
	return r
}

// StareTraduceri: folderul de import și câte texte corectate are fiecare limbă.
func (a *App) StareTraduceri() RaspunsTraduceri { return a.stareTraduceri() }

// citesteFolderTraduceri citește fișierele <limba>.json din folder (cheie → text).
func (a *App) citesteFolderTraduceri() ([]datebaza.FisierTraduceri, error) {
	nume, err := filepath.Glob(filepath.Join(a.folderTraduceri(), "*.json"))
	if err != nil {
		return nil, err
	}
	if len(nume) == 0 {
		return nil, &i18n.Eroare{Cod: "traduceri.folder_gol", Param: map[string]string{"folder": a.folderTraduceri()}}
	}
	if len(nume) > maxFisiereTraduceri {
		return nil, &i18n.Eroare{Cod: "traduceri.prea_multe_fisiere"}
	}
	sort.Strings(nume)
	rez := []datebaza.FisierTraduceri{}
	for _, n := range nume {
		fisier := filepath.Base(n)
		info, err := os.Stat(n)
		if err != nil {
			return nil, err
		}
		if info.Size() > maxOctetiFisierTraduce {
			return nil, &i18n.Eroare{Cod: "traduceri.fisier_invalid", Param: map[string]string{"fisier": fisier}}
		}
		data, err := os.ReadFile(n)
		if err != nil {
			return nil, err
		}
		var texte map[string]string
		limba := strings.TrimSuffix(fisier, filepath.Ext(fisier))
		if err := json.Unmarshal(data, &texte); err != nil || !i18n.CodLimbaValid(limba) {
			return nil, &i18n.Eroare{Cod: "traduceri.fisier_invalid", Param: map[string]string{"fisier": fisier}}
		}
		rez = append(rez, datebaza.FisierTraduceri{Limba: limba, Texte: texte})
	}
	return rez, nil
}

func (a *App) doarAdministrator() *MesajEroare {
	if c := a.contCurent(); c == nil || c.Rol != conturi.RolAdministrator {
		return faraDrept()
	}
	return nil
}

// VerificaFolderTraduceri arată ce ar importa fișierele din folder, fără să scrie nimic.
func (a *App) VerificaFolderTraduceri() RaspunsTraduceri {
	r := a.stareTraduceri()
	if a.db == nil {
		r.Eroare = &MesajEroare{Cod: "date.eroare"}
		return r
	}
	fisiere, err := a.citesteFolderTraduceri()
	if err != nil {
		r.Eroare = eroareInterfata(err)
		return r
	}
	r.Rezumate, err = datebaza.Nou(a.db).PrevizualizeazaTraduceri(a.contextLucru(), fisiere)
	if err != nil {
		r.Eroare = eroareInterfata(err)
	}
	return r
}

// ImportaFolderTraduceri importă fișierele din folder (doar administratorul).
func (a *App) ImportaFolderTraduceri() RaspunsTraduceri {
	if e := a.doarAdministrator(); e != nil {
		r := a.stareTraduceri()
		r.Eroare = e
		return r
	}
	fisiere, err := a.citesteFolderTraduceri()
	if err != nil {
		r := a.stareTraduceri()
		r.Eroare = eroareInterfata(err)
		return r
	}
	rezumate, err := datebaza.Nou(a.db).ImportaTraduceri(a.contextLucru(), fisiere)
	a.incarcaTraduceri()
	r := a.stareTraduceri()
	r.Rezumate = rezumate
	if err != nil {
		r.Eroare = eroareInterfata(err)
	}
	return r
}

// StergeTraduceriLimba șterge corecturile unei limbi (doar administratorul).
func (a *App) StergeTraduceriLimba(limba string) RaspunsTraduceri {
	if e := a.doarAdministrator(); e != nil {
		r := a.stareTraduceri()
		r.Eroare = e
		return r
	}
	err := datebaza.Nou(a.db).StergeTraduceri(a.contextLucru(), limba)
	a.incarcaTraduceri()
	r := a.stareTraduceri()
	if err != nil {
		r.Eroare = eroareInterfata(err)
	}
	return r
}

// ScrieFisierTraduceri scrie karios/traduceri/<limba>.json cu textele de acum ale limbii
// (inclusiv corecturile), ca punct de plecare pentru traducător. Un fișier existent se
// păstrează cu sufixul .vechi.
func (a *App) ScrieFisierTraduceri(limba string) RaspunsTraduceri {
	if e := a.doarAdministrator(); e != nil {
		r := a.stareTraduceri()
		r.Eroare = e
		return r
	}
	r := a.stareTraduceri()
	if !i18n.CodLimbaValid(limba) {
		r.Eroare = &MesajEroare{Cod: "traduceri.limba_invalida"}
		return r
	}
	if err := os.MkdirAll(a.folderTraduceri(), 0o700); err != nil {
		r.Eroare = eroareInterfata(err)
		return r
	}
	cale := filepath.Join(a.folderTraduceri(), limba+".json")
	if _, err := os.Stat(cale); err == nil {
		_ = os.Rename(cale, cale+".vechi")
	}
	data, err := json.MarshalIndent(i18n.Catalog(i18n.Limba(limba)), "", "  ")
	if err == nil {
		err = os.WriteFile(cale, append(data, '\n'), 0o600)
	}
	if err != nil {
		r.Eroare = eroareInterfata(fmt.Errorf("traduceri: %w", err))
		return r
	}
	r.Cale = cale
	return r
}
