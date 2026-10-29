package main

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"

	wruntime "github.com/wailsapp/wails/v2/pkg/runtime"

	"pontix/internal/ai"
)

// EvenimentSarcinaAI e numele evenimentului Wails trimis la terminarea unei sarcini (F14-T07).
const EvenimentSarcinaAI = "ai:sarcina_terminata"

// RaspunsSarcinaAI întoarce identificatorul sarcinii pornite.
type RaspunsSarcinaAI struct {
	ID     string       `json:"id"`
	Eroare *MesajEroare `json:"eroare,omitempty"`
}

// RezultatSarcinaAI e conținutul evenimentului de terminare.
type RezultatSarcinaAI struct {
	ID      string              `json:"id"`
	Raspuns *ai.RaspunsAsistent `json:"raspuns,omitempty"`
	Eroare  *MesajEroare        `json:"eroare,omitempty"`
	Oprita  bool                `json:"oprita"`
}

type sarcinaAI struct {
	id       string
	anuleaza context.CancelFunc
}

// PornesteSarcinaAI pornește interogarea în fundal și întoarce imediat. Rezultatul vine
// prin evenimentul EvenimentSarcinaAI, deci nu se pierde dacă utilizatorul schimbă ecranul.
// Rulează o singură sarcină o dată.
func (a *App) PornesteSarcinaAI(cerere CerereInterogareAI) RaspunsSarcinaAI {
	ctx, anuleaza := context.WithCancel(a.contextLucru())
	s := &sarcinaAI{id: idSarcina(), anuleaza: anuleaza}
	a.aiMu.Lock()
	if a.aiSarcina != nil {
		a.aiMu.Unlock()
		anuleaza()
		return RaspunsSarcinaAI{Eroare: &MesajEroare{Cod: "ai.sarcina_activa"}}
	}
	a.aiSarcina = s
	a.aiMu.Unlock()

	go func() {
		rez := a.ruleazaSarcinaAI(ctx, s.id, cerere)
		a.aiMu.Lock()
		if a.aiSarcina == s {
			a.aiSarcina = nil
		}
		a.aiMu.Unlock()
		anuleaza()
		a.emiteSarcinaAI(rez)
	}()
	return RaspunsSarcinaAI{ID: s.id}
}

func (a *App) ruleazaSarcinaAI(ctx context.Context, id string, cerere CerereInterogareAI) RezultatSarcinaAI {
	raspuns := a.interogheazaAI(ctx, cerere)
	rez := RezultatSarcinaAI{ID: id, Raspuns: raspuns.Raspuns, Eroare: raspuns.Eroare}
	if errors.Is(ctx.Err(), context.Canceled) {
		rez = RezultatSarcinaAI{ID: id, Oprita: true}
	}
	return rez
}

// OpresteSarcinaAI anulează sarcina în curs; cererea către motor se întrerupe.
func (a *App) OpresteSarcinaAI() bool {
	a.aiMu.Lock()
	defer a.aiMu.Unlock()
	if a.aiSarcina == nil {
		return false
	}
	a.aiSarcina.anuleaza()
	return true
}

// SarcinaAIActiva întoarce identificatorul sarcinii în curs sau "".
func (a *App) SarcinaAIActiva() string {
	a.aiMu.Lock()
	defer a.aiMu.Unlock()
	if a.aiSarcina == nil {
		return ""
	}
	return a.aiSarcina.id
}

func (a *App) emiteSarcinaAI(rez RezultatSarcinaAI) {
	if a.ctx == nil {
		return
	}
	wruntime.EventsEmit(a.ctx, EvenimentSarcinaAI, rez)
}

func idSarcina() string {
	var b [8]byte
	if _, err := rand.Read(b[:]); err != nil {
		panic(err)
	}
	return hex.EncodeToString(b[:])
}
