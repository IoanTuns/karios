package main

import (
	"context"
	"testing"
)

func TestInchidereFaraModificariNuSeOpreste(t *testing.T) {
	a := &App{}
	if a.inainteDeInchidere(context.Background()) {
		t.Fatal("fereastra fără modificări nesalvate nu trebuie oprită")
	}
	a.SeteazaModificariNesalvate(true, "ro")
	a.SeteazaModificariNesalvate(false, "ro")
	if a.inainteDeInchidere(context.Background()) {
		t.Fatal("după salvare, fereastra se închide normal")
	}
}
