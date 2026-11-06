package main

import (
	"context"
	"pontix/internal/datebaza"
	"pontix/internal/i18n"
	"testing"
)

func TestStergereConfigurareNumaiAdministratorAutentificat(t *testing.T) {
	a := appCuConturi(t)
	if r := a.StergeConfigurare("grupuri", "x", "x"); r.Eroare == nil || r.Eroare.Cod != "cont.doar_administrator" {
		t.Fatal("fara sesiune", r.Eroare)
	}
	if r := a.CreeazaPrimulAdministrator(cerereCont("admin-test")); r.Eroare != nil {
		t.Fatal(r.Eroare)
	}
	if r := a.CreeazaCont(cerereCont("operator-test")); r.Eroare != nil {
		t.Fatal(r.Eroare)
	}
	if err := datebaza.Nou(a.db).SalveazaElementNomenclator(context.Background(), "grupuri", datebaza.ElementNomenclator{Denumire: i18n.Text{"ro": "Stergere test"}, Activ: false}); err != nil {
		t.Fatal(err)
	}
	d, err := datebaza.Nou(a.db).Citeste(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	var id, rev string
	for _, e := range d.GrupuriProfesionale {
		if e.Denumire.In("ro") == "Stergere test" {
			id, rev = e.ID, e.Revizie
		}
	}
	a.IesireDinCont()
	if r := a.Autentifica("operator-test", "parola123"); r.Eroare != nil {
		t.Fatal(r.Eroare)
	}
	if r := a.StergeConfigurare("grupuri", id, rev); r.Eroare == nil || r.Eroare.Cod != "cont.doar_administrator" {
		t.Fatal("operator autorizat", r.Eroare)
	}
	a.IesireDinCont()
	if r := a.Autentifica("admin-test", "parola123"); r.Eroare != nil {
		t.Fatal(r.Eroare)
	}
	for _, schimbare := range []string{"rol='OPERATOR'", "activ=0"} {
		if _, err := a.db.Exec("UPDATE conturi SET " + schimbare + " WHERE nume='admin-test'"); err != nil {
			t.Fatal(err)
		}
		if r := a.StergeConfigurare("grupuri", id, rev); r.Eroare == nil || r.Eroare.Cod != "cont.doar_administrator" {
			t.Fatal("sesiune veche autorizata", r.Eroare)
		}
		if _, err := a.db.Exec("UPDATE conturi SET rol='ADMINISTRATOR',activ=1 WHERE nume='admin-test'"); err != nil {
			t.Fatal(err)
		}
	}
	if r := a.StergeConfigurare("grupuri", id, rev); r.Eroare != nil {
		t.Fatal("administrator refuzat", r.Eroare)
	}
}
