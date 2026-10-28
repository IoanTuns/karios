package main

import (
	"pontix/internal/audit"
)

// RaspunsAudit întoarce rezultatul verificării sigiliilor (D044).
type RaspunsAudit struct {
	Rezultat audit.Rezultat `json:"rezultat"`
	Eroare   *MesajEroare   `json:"eroare,omitempty"`
}

// StareAudit spune dacă la pornire s-au găsit modificări făcute în afara aplicației.
func (a *App) StareAudit() RaspunsAudit {
	if a.paza == nil {
		return RaspunsAudit{Rezultat: audit.Rezultat{Stare: audit.StareInactiv, Tabele: []string{}}}
	}
	return RaspunsAudit{Rezultat: a.paza.Rezultat()}
}

// ConfirmaIncidentAudit acceptă starea curentă a datelor după un incident, cu motiv
// obligatoriu scris în jurnal, și deblochează editarea.
func (a *App) ConfirmaIncidentAudit(motiv string) RaspunsAudit {
	if a.paza == nil {
		return a.StareAudit()
	}
	rez, err := audit.Confirma(a.contextLucru(), a.db, a.paza, motiv)
	if err != nil {
		return RaspunsAudit{Rezultat: rez, Eroare: eroareInterfata(err)}
	}
	return RaspunsAudit{Rezultat: rez}
}
