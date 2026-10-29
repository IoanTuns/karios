package main

import (
	"context"
	"path/filepath"
	"testing"

	"pontix/internal/bd"
	"pontix/internal/conturi"
	"pontix/internal/datebaza"
	"pontix/internal/seif"
	"pontix/internal/stocare"
)

func appCuConturi(t *testing.T) *App {
	t.Helper()
	dir := t.TempDir()
	sf, err := seif.Creeaza(filepath.Join(dir, stocare.NumeSeif), bd.LungimeCheie)
	if err != nil {
		t.Fatal(err)
	}
	db, err := bd.Deschide(filepath.Join(dir, stocare.NumeBaza), sf.CheieDate)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { db.Close() })
	if _, err := bd.Migreaza(context.Background(), db, bd.Migrari); err != nil {
		t.Fatal(err)
	}
	a := NewApp(db)
	a.radacina = dir
	return a
}

func cerereCont(nume string) conturi.CerereCont {
	return conturi.CerereCont{Nume: nume, Parola: "parola123", Intrebari: []string{"intrebare_1", "intrebare_2", "intrebare_3"}, Raspunsuri: []string{"a", "b", "c"}}
}

func TestFluxConturiInApp(t *testing.T) {
	a := appCuConturi(t)
	if r := a.StareCont(); !r.NecesitaCont || r.Autentificat || len(r.Intrebari) != 10 {
		t.Fatalf("start: %+v", r)
	}
	if r := a.CreeazaPrimulAdministrator(cerereCont("ana")); r.Eroare != nil || !r.Autentificat || r.Cont.Rol != conturi.RolAdministrator {
		t.Fatalf("primul administrator: %+v", r)
	}
	if r := a.CreeazaPrimulAdministrator(cerereCont("altul")); r.Eroare == nil || r.Eroare.Cod != "cont.exista_deja" {
		t.Fatalf("al doilea „prim” administrator: %+v", r.Eroare)
	}
	if got := datebaza.Actor(a.contextLucru()); got != "ana" {
		t.Fatalf("autorul din context: %q", got)
	}
	if r := a.CreeazaCont(cerereCont("dan")); r.Eroare != nil || len(r.Conturi) != 2 {
		t.Fatalf("cont operator: %+v", r)
	}
	if r := a.IesireDinCont(); r.Autentificat {
		t.Fatal("ieșire")
	}
	if r := a.Autentifica("dan", "parola123"); !r.Autentificat || r.Cont.Rol != conturi.RolOperator {
		t.Fatalf("operator: %+v", r)
	}
	if r := a.ListaConturi(); r.Eroare == nil || r.Eroare.Cod != "cont.doar_administrator" {
		t.Fatalf("operatorul nu gestionează conturi: %+v", r)
	}
	if r := a.Autentifica("dan", "gresit123"); r.Eroare == nil || r.Eroare.Cod != "cont.date_gresite" {
		t.Fatalf("parolă greșită: %+v", r.Eroare)
	}
}
