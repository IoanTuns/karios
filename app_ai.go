package main

import (
	"fmt"
	"path/filepath"

	"pontix/internal/ai"
	"pontix/internal/datebaza"
	"pontix/internal/i18n"
	"pontix/internal/portabil"
)

type RaspunsStareAI struct {
	Stare  ai.StareMediu `json:"stare"`
	Eroare *MesajEroare  `json:"eroare"`
}

type CerereInterogareAI struct {
	Istoric []ai.MesajConversatie `json:"istoric"`
	An      int                   `json:"an"`
	Luna    int                   `json:"luna"`
}

type RaspunsInterogareAI struct {
	Raspuns *ai.RaspunsAsistent `json:"raspuns"`
	Eroare  *MesajEroare        `json:"eroare"`
}

func (a *App) folderAI() string {
	radacina := a.radacina
	if radacina == "" {
		foldere, err := portabil.Pregateste()
		if err == nil {
			return foldere.AI
		}
	}
	return filepath.Join(radacina, portabil.NumeFolderDate, "ai")
}

func (a *App) asiguratManagerAI() *ai.Manager {
	if a.aiManager == nil {
		a.aiManager = ai.NouManager(a.folderAI(), ai.PortImplicit)
	}
	return a.aiManager
}

// StareAI raportează disponibilitatea motorului și starea procesului local.
func (a *App) StareAI() RaspunsStareAI {
	folder := a.folderAI()
	stare := ai.DetecteazaMediu(folder, ai.PortImplicit)
	return RaspunsStareAI{Stare: stare}
}

// PornesteAI lansează motorul local în fundal dacă fișierele sunt prezente.
func (a *App) PornesteAI(numeModel string) RaspunsStareAI {
	mgr := a.asiguratManagerAI()
	err := mgr.Porneste(a.contextLucru(), numeModel)
	stare := ai.DetecteazaMediu(a.folderAI(), mgr.Port())
	if err != nil {
		return RaspunsStareAI{
			Stare:  stare,
			Eroare: eroareInterfata(fmt.Errorf("eroare_pornire_ai: %w", err)),
		}
	}
	return RaspunsStareAI{Stare: stare}
}

// OpresteAI oprește procesul copil al serverului AI local.
func (a *App) OpresteAI() RaspunsStareAI {
	mgr := a.asiguratManagerAI()
	_ = mgr.Opreste()
	stare := ai.DetecteazaMediu(a.folderAI(), mgr.Port())
	return RaspunsStareAI{Stare: stare}
}

// InterogheazaAI trimite mesajul utilizatorului către asistentul local.
func (a *App) InterogheazaAI(cerere CerereInterogareAI) RaspunsInterogareAI {
	mgr := a.asiguratManagerAI()
	if !mgr.EsteActiv() {
		// Încearcă pornirea automată dacă este instalat
		if err := mgr.Porneste(a.contextLucru(), ""); err != nil {
			return RaspunsInterogareAI{
				Eroare: eroareInterfata(fmt.Errorf("ai_inactiv: %w", err)),
			}
		}
	}

	// Colectează contextul de lucru (angajați și schimburi)
	date, err := datebaza.Nou(a.db).Citeste(a.contextLucru())
	if err != nil {
		return RaspunsInterogareAI{Eroare: eroareInterfata(err)}
	}

	schimburi, _ := datebaza.Nou(a.db).ListaSchimburi(a.contextLucru())

	contextPontix := ai.ContextPontix{
		Limba: "ro",
		An:    cerere.An,
		Luna:  cerere.Luna,
	}

	for _, ang := range date.Angajati {
		if ang.Activ {
			contextPontix.Angajati = append(contextPontix.Angajati, ai.InfoAngajat{
				ID:      ang.ID,
				Nume:    ang.Nume,
				Prenume: ang.Prenume,
				Marca:   ang.Marca,
				Functie: ang.Functie,
			})
		}
	}

	for _, s := range schimburi {
		if s.Activ {
			contextPontix.Schimburi = append(contextPontix.Schimburi, ai.InfoSchimb{
				Cod:      s.Cod,
				Denumire: s.Denumire.In(i18n.Romana),
				Ore:      int(s.DurataMinute / 60),
			})
		}
	}

	client := ai.NouClient(mgr.Port())
	raspuns, err := client.Interogheaza(a.contextLucru(), cerere.Istoric, contextPontix)
	if err != nil {
		return RaspunsInterogareAI{Eroare: eroareInterfata(err)}
	}

	return RaspunsInterogareAI{Raspuns: &raspuns}
}

