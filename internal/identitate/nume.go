// Package identitate formateaza numele de afisare fara a modifica identitatea salvata.
package identitate

import (
	"pontix/internal/i18n"
	"strings"
)

func Nume(nume, prenume string, subordonat bool, limba i18n.Limba) string {
	complet := strings.TrimSpace(nume + " " + prenume)
	if !subordonat {
		return complet
	}
	return i18n.T(limba, "ui.nume_subordonat", map[string]string{"nume": complet, "tip": i18n.T(limba, "rapoarte.subtip_garda", nil)})
}
