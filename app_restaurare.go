package main

import (
	wruntime "github.com/wailsapp/wails/v2/pkg/runtime"

	"pontix/internal/audit"
	"pontix/internal/stocare"
)

// RaspunsBackupuri întoarce copiile din backup/ și restaurarea programată, dacă există (F05-T25).
type RaspunsBackupuri struct {
	Backupuri  []stocare.Backup `json:"backupuri"`
	Programata string           `json:"programata"`
	Eroare     *MesajEroare     `json:"eroare,omitempty"`
}

// RaspunsVerificareBackup întoarce rezultatul verificării sigiliilor unei copii.
type RaspunsVerificareBackup struct {
	Rezultat audit.Rezultat `json:"rezultat"`
	Eroare   *MesajEroare   `json:"eroare,omitempty"`
}

func (a *App) ListaBackupuri() RaspunsBackupuri {
	l, err := stocare.ListaBackupuri(a.radacina)
	if err != nil {
		return RaspunsBackupuri{Eroare: eroareInterfata(err)}
	}
	return RaspunsBackupuri{Backupuri: l, Programata: stocare.RestaurareProgramata(a.radacina)}
}

func (a *App) VerificaBackup(nume string) RaspunsVerificareBackup {
	r, err := stocare.VerificaBackup(a.contextLucru(), a.radacina, nume, a.seif.CheieDate)
	if err != nil {
		return RaspunsVerificareBackup{Eroare: eroareInterfata(err)}
	}
	return RaspunsVerificareBackup{Rezultat: r}
}

// ProgrameazaRestaurare verifică din nou copia și o programează pentru următoarea pornire.
// O copie cu sigiliile rupte se refuză. Doar administratorul poate restaura.
func (a *App) ProgrameazaRestaurare(nume string) RaspunsBackupuri {
	if !a.esteAdministrator() {
		return RaspunsBackupuri{Eroare: faraDrept()}
	}
	r, err := stocare.VerificaBackup(a.contextLucru(), a.radacina, nume, a.seif.CheieDate)
	if err != nil {
		return RaspunsBackupuri{Eroare: eroareInterfata(err)}
	}
	if r.Stare == audit.StareIncident {
		return RaspunsBackupuri{Eroare: &MesajEroare{Cod: "stocare.backup_modificat"}}
	}
	if err := stocare.ProgrameazaRestaurare(a.radacina, nume); err != nil {
		return RaspunsBackupuri{Eroare: eroareInterfata(err)}
	}
	return a.ListaBackupuri()
}

func (a *App) AnuleazaRestaurarea() RaspunsBackupuri {
	if err := stocare.AnuleazaRestaurare(a.radacina); err != nil {
		return RaspunsBackupuri{Eroare: eroareInterfata(err)}
	}
	return a.ListaBackupuri()
}

// InchideAplicatia închide fereastra; restaurarea programată se aplică la repornire.
func (a *App) InchideAplicatia() {
	if a.ctx != nil {
		wruntime.Quit(a.ctx)
	}
}
