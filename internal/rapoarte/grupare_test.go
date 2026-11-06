package rapoarte

import (
	"context"
	"testing"
)

func TestRaportGrupuriSefiAlfabeticSiFiltre(t *testing.T) {
	m := pregateste(t)
	dateGarziDetaliat(t, m)
	for _, sql := range []string{
		`INSERT INTO grupuri_profesionale VALUES('med','{"ro":"Medici","en":"Doctors"}',1,'x',2),('as','{"ro":"Asistenti","en":"Nurses"}',1,'x',1)`,
		`INSERT INTO functii(id,denumire,activ,modificat_la,grup_id,conducere) VALUES('m','{"ro":"Medic"}',1,'x','med',0),('sef','{"ro":"Medic sef"}',1,'x','med',1),('as','{"ro":"Asistent"}',1,'x','as',0)`,
		`INSERT INTO angajati(id,nume,prenume,marca,creat_la,modificat_la) VALUES('sef','Zef','Sef','SEF','x','x'),('alfa','Ana','Fictiva','ALFA','x','x'),('as','Personal','Fictiv','AS','x','x')`,
		`INSERT INTO angajati_functii VALUES('a','m'),('sef','sef'),('alfa','m'),('as','as')`,
	} {
		if _, err := m.db.Exec(sql); err != nil {
			t.Fatal(err)
		}
	}
	c := Cerere{ModelID: "raport-pontaj-detaliat", An: 2026, Luna: 10, Limba: "ro"}
	citeste := func() Document {
		d, err := m.s.Pregateste(context.Background(), c, momentFix)
		if err != nil {
			t.Fatal(err)
		}
		return d
	}
	d := citeste()

	marci := []string{}
	grupuri := []string{}
	for i, r := range d.Randuri {
		if r.Grup {
			grupuri = append(grupuri, r.Valori[0])
			continue
		}
		marci = append(marci, valoriRand(d, i)["marca"])
	}
	if len(grupuri) != 3 || grupuri[0] != "Asistenti" || grupuri[1] != "Medici" {
		t.Fatal(grupuri)
	}
	if len(marci) < 5 || marci[0] != "AS" || marci[1] != "SEF" || marci[2] != "ALFA" || marci[3] != "DEMO-R01" || marci[4] != "GARDA-R01" {
		t.Fatal(marci)
	}
	c.GrupID = "med"
	c.Garda = "DA"
	d = citeste()
	if len(d.Randuri) != 2 || valoriRand(d, 1)["marca"] != "GARDA-R01" {
		t.Fatal("filtre grup/garda", d.Randuri)
	}
	c.Garda = ""
	c.Pozitie = "CONDUCERE"
	d = citeste()
	if len(d.Randuri) != 2 || valoriRand(d, 1)["marca"] != "SEF" {
		t.Fatal("filtru pozitie", d.Randuri)
	}
	c.Pozitie = ""
	c.FunctieID = "m"
	c.ContractID = "INTREAGA"
	d = citeste()
	if len(d.Randuri) != 4 {
		t.Fatal("functie/contract", d.Randuri)
	}
}

func TestGraficGarziFiltratNuIncludeDupaAmiaza(t *testing.T) {
	m := pregateste(t)
	dateGarziDetaliat(t, m)
	d, err := m.s.Pregateste(context.Background(), Cerere{ModelID: "raport-garzi", An: 2026, Luna: 10, Limba: "ro", Garda: "DA"}, momentFix)
	if err != nil {
		t.Fatal(err)
	}
	if len(d.Randuri) != 2 {
		t.Fatal(d.Randuri)
	}
	r, total := valoriRand(d, 0), valoriRand(d, 1)
	if r["marca"] != "GARDA-R01" || r["total_garzi"] != "1" || total["total_garzi"] != "1" || total["total_ore"] != "24" || total["2026-10-11"] != "0" {
		t.Fatal(r, total)
	}
}
