package main

import (
	"context"
	"testing"

	"pontix/internal/seif"
	"pontix/internal/stocare"
)

func porneste(t *testing.T, dir string) *App {
	t.Helper()
	d, err := stocare.Deschide(context.Background(), dir)
	if err != nil {
		t.Fatal(err)
	}
	a := NewApp(nil)
	a.radacina = dir
	a.aplicaDeschiderea(d)
	t.Cleanup(func() {
		if a.db != nil {
			a.db.Close()
		}
	})
	return a
}

func TestCheiaDeRecuperareSiDeblocarea(t *testing.T) {
	dir := t.TempDir()
	a := porneste(t, dir)
	if a.StareSeif().Blocat {
		t.Fatal("prima pornire e blocată")
	}
	// Înainte de autentificare, codul nu se arată nimănui.
	if r := a.StareRecuperare(); !r.Necesara || r.Cod != "" {
		t.Fatalf("recuperare fără cont: %+v", r)
	}
	if r := a.CreeazaPrimulAdministrator(cerereCont("ana")); r.Eroare != nil {
		t.Fatal(r.Eroare)
	}
	r := a.StareRecuperare()
	if !r.Necesara || r.Cod == "" {
		t.Fatalf("recuperare pentru administrator: %+v", r)
	}
	if c := a.ConfirmaCheiaRecuperare("AAAAA-BBBBB-CCCCC-DDDDD-EEEEE"); c.Eroare == nil || c.Eroare.Cod != "seif.recuperare_gresita" {
		t.Fatalf("cheie greșită: %+v", c)
	}
	cod := r.Cod
	if c := a.ConfirmaCheiaRecuperare(cod); c.Eroare != nil || c.Necesara {
		t.Fatalf("confirmare: %+v", c)
	}
	if r := a.CreeazaCont(cerereCont("dan")); r.Eroare != nil {
		t.Fatal(r.Eroare)
	}
	// Alt utilizator Windows: seiful nu are copia lui W.
	for _, id := range a.seif.CopiiCuPrefix(seif.PrefixWindows) {
		if err := a.seif.StergeCopie(id); err != nil {
			t.Fatal(err)
		}
	}
	a.db.Close()
	a.db = nil

	b := porneste(t, dir)
	if !b.StareSeif().Blocat {
		t.Fatal("fără copia W datele nu sunt blocate")
	}
	if r := b.DeblocheazaCuParola("dan", "gresit123"); r.Eroare == nil || r.Eroare.Cod != "seif.secret_gresit" || !b.StareSeif().Blocat {
		t.Fatalf("parolă greșită: %+v", r.Eroare)
	}
	if r := b.DeblocheazaCuParola("dan", "parola123"); r.Eroare != nil || !r.Autentificat || r.Cont.Nume != "dan" {
		t.Fatalf("deblocare cu parola: %+v", r)
	}
	if b.StareRecuperare().Necesara {
		t.Fatal("cheia confirmată se cere din nou")
	}
	for _, id := range b.seif.CopiiCuPrefix(seif.PrefixWindows) {
		if err := b.seif.StergeCopie(id); err != nil {
			t.Fatal(err)
		}
	}
	b.db.Close()
	b.db = nil

	c := porneste(t, dir)
	if r := c.DeblocheazaCuRecuperare("gresit"); r.Eroare == nil || !c.StareSeif().Blocat {
		t.Fatalf("cheie de recuperare greșită: %+v", r.Eroare)
	}
	if r := c.DeblocheazaCuRecuperare(cod); r.Eroare != nil || c.StareSeif().Blocat || r.Autentificat {
		t.Fatalf("deblocare cu cheia de recuperare: %+v", r)
	}
	if r := c.Autentifica("ana", "parola123"); !r.Autentificat {
		t.Fatalf("autentificare după recuperare: %+v", r)
	}
	// Operatorul nu vede și nu schimbă cheia de recuperare.
	c.IesireDinCont()
	c.Autentifica("dan", "parola123")
	if r := c.GenereazaCheieRecuperareNoua(); r.Eroare == nil {
		t.Fatal("operatorul a generat o cheie nouă")
	}
}
