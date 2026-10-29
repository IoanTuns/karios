package main

import (
	"context"

	wruntime "github.com/wailsapp/wails/v2/pkg/runtime"

	"pontix/internal/i18n"
)

// SeteazaModificariNesalvate primește din interfață dacă vreun ecran are modificări
// nesalvate (F05-T24). `beforeunload` nu oprește închiderea ferestrei Wails.
func (a *App) SeteazaModificariNesalvate(da bool, limba string) {
	a.inchidereMu.Lock()
	a.modificariNesalvate = da
	a.inchidereMu.Unlock()
	a.SeteazaLimbaInterfata(limba)
}

// SeteazaLimbaInterfata reține limba interfeței pentru dialogurile native (F03-T10).
func (a *App) SeteazaLimbaInterfata(limba string) {
	l := i18n.Limba(limba)
	if l != i18n.Romana && l != i18n.Engleza {
		return
	}
	a.inchidereMu.Lock()
	a.limbaInterfata = limba
	a.inchidereMu.Unlock()
}

// limba întoarce limba interfeței sau limba implicită.
func (a *App) limba() i18n.Limba {
	a.inchidereMu.Lock()
	defer a.inchidereMu.Unlock()
	if a.limbaInterfata == "" {
		return i18n.Implicita
	}
	return i18n.Limba(a.limbaInterfata)
}

// inainteDeInchidere întreabă, dacă există modificări nesalvate, dacă fereastra se închide
// totuși. Întoarce true ca să oprească închiderea.
func (a *App) inainteDeInchidere(ctx context.Context) bool {
	a.inchidereMu.Lock()
	nesalvate := a.modificariNesalvate
	a.inchidereMu.Unlock()
	if !nesalvate {
		return false
	}
	limba := a.limba()
	da, nu := i18n.T(limba, "ui.inchide_oricum", nil), i18n.T(limba, "ui.ramai", nil)
	raspuns, err := wruntime.MessageDialog(ctx, wruntime.MessageDialogOptions{
		Type:          wruntime.QuestionDialog,
		Title:         i18n.T(limba, "ui.modificari_nesalvate_titlu", nil),
		Message:       i18n.T(limba, "ui.modificari_nesalvate_inchidere", nil),
		Buttons:       []string{da, nu},
		DefaultButton: nu,
		CancelButton:  nu,
	})
	if err != nil {
		return true // dacă dialogul nu se poate afișa, nu se pierd datele
	}
	return raspuns != da
}
