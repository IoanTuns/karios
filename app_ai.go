package main

import (
	"context"
	"fmt"
	"os"
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
	if filepath.Base(radacina) == portabil.NumeFolderDate {
		return filepath.Join(radacina, "ai")
	}
	return filepath.Join(radacina, portabil.NumeFolderDate, "ai")
}

const fisierFaraAutoPornire = "fara-autostart"

// EsteAutoPornireAI verifică dacă pornirea automată a AI este activă (implicit true dacă nu există fișierul de dezactivare).
func (a *App) EsteAutoPornireAI() bool {
	folder := a.folderAI()
	caleFisier := filepath.Join(folder, fisierFaraAutoPornire)
	_, err := os.Stat(caleFisier)
	return os.IsNotExist(err)
}

// SeteazaAutoPornireAI activează sau dezactivează pornirea automată a motorului AI la pornirea aplicației.
func (a *App) SeteazaAutoPornireAI(activa bool) bool {
	folder := a.folderAI()
	caleFisier := filepath.Join(folder, fisierFaraAutoPornire)
	if activa {
		_ = os.Remove(caleFisier)
	} else {
		_ = os.WriteFile(caleFisier, []byte("fara-autostart"), 0644)
	}
	return a.EsteAutoPornireAI()
}

// autoPornesteAIDacaEPosibil pornește motorul în fundal la lansarea aplicației dacă fișierele sunt prezente și opțiunea e activă.
func (a *App) autoPornesteAIDacaEPosibil() {
	if !a.EsteAutoPornireAI() {
		return
	}
	folder := a.folderAI()
	stare := ai.DetecteazaMediu(folder, ai.PortImplicit)
	if !stare.Instalat || stare.Activ || len(stare.ModeleDisponibile) == 0 {
		return
	}
	mgr := a.asiguratManagerAI()
	_ = mgr.Porneste(context.Background(), "")
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
			interval := ""
			if s.DurataMinute > 0 {
				startH := s.OraStart / 60
				startM := s.OraStart % 60
				endMin := (s.OraStart + s.DurataMinute) % 1440
				endH := endMin / 60
				endM := endMin % 60
				interval = fmt.Sprintf("%02d:%02d–%02d:%02d", startH, startM, endH, endM)
			}
			contextPontix.Schimburi = append(contextPontix.Schimburi, ai.InfoSchimb{
				Cod:      s.Cod,
				Denumire: s.Denumire.In(i18n.Romana),
				Ore:      int(s.DurataMinute / 60),
				Interval: interval,
			})
		}
	}

	client := ai.NouClient(mgr.Port())
	raspuns, err := client.Interogheaza(a.contextLucru(), cerere.Istoric, contextPontix)
	if err != nil {
		return RaspunsInterogareAI{
			Eroare: &MesajEroare{
				Cod:   "ai.eroare_generare",
				Param: map[string]string{"detalii": err.Error()},
			},
		}
	}

	return RaspunsInterogareAI{Raspuns: &raspuns}
}
