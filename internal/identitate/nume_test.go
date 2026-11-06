package identitate

import (
	"pontix/internal/i18n"
	"testing"
)

func TestNumeSubordonatSiOmonim(t *testing.T) {
	for _, limba := range []i18n.Limba{i18n.Romana, i18n.Engleza} {
		simplu := Nume("Fictiv", "Ana", false, limba)
		marcat := Nume("Fictiv", "Ana", true, limba)
		if simplu != "Fictiv Ana" || marcat == simplu || marcat != "Fictiv Ana ("+i18n.T(limba, "rapoarte.subtip_garda", nil)+")" {
			t.Fatal(simplu, marcat)
		}
	}
}
