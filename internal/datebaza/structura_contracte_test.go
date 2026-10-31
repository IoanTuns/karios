package datebaza

import (
	"context"
	"pontix/internal/bd/acces"
	"pontix/internal/calcul"
	"pontix/internal/i18n"
	"testing"
)

func TestStructuraContractGardaSiReguliPersistente(t *testing.T) {
	m := pregateste(t)
	ctx := context.Background()
	salveaza := func(tip string, e ElementNomenclator) ElementNomenclator {
		t.Helper()
		if err := m.serviciu.SalveazaElementNomenclator(ctx, tip, e); err != nil {
			t.Fatal(err)
		}
		d := citeste(t, m.serviciu)
		lista := d.GrupuriProfesionale
		if tip == "functii" {
			lista = d.Functii
		}
		if tip == "contracte" {
			lista = d.TipuriContract
		}
		for _, v := range lista {
			if numeElement(v) == numeElement(e) {
				return v
			}
		}
		t.Fatal("element lipsa")
		return e
	}
	grup := salveaza("grupuri", ElementNomenclator{Denumire: i18n.Text{"ro": "Medici fictivi", "en": "Test doctors"}, Activ: true, Ordine: 2})
	functie := salveaza("functii", ElementNomenclator{Denumire: i18n.Text{"ro": "Medic sef fictiv"}, Activ: true, GrupID: grup.ID, Conducere: true})
	a := angajatTest("CONTRACT-G")
	a.ContractID = "GARDA"
	a.FunctieID = functie.ID
	if err := m.serviciu.SalveazaAngajat(ctx, a); err != nil {
		t.Fatal(err)
	}
	d := citeste(t, m.serviciu)
	a = d.Angajati[0]
	if a.NormaDefinita || a.NormaMinute != 0 || !a.Garda || !a.Conducere || a.GrupID != grup.ID {
		t.Fatal(a)
	}
	cfg, _, err := m.serviciu.ConfiguratiePontaj(ctx)
	if err != nil {
		t.Fatal(err)
	}
	cfg.Conditionate = []calcul.RegulaConditionata{{Activ: true, GrupID: grup.ID, ContractID: "GARDA", Garda: "DA", Regula: calcul.RegulaColoana{Coloana: "ore_garda", Sursa: "ORE_SCHIMBURI_ZNL", Coduri: []string{"G24"}, FactorSutimi: 200}}}
	if err := m.serviciu.SalveazaConfiguratiePontaj(ctx, cfg); err != nil {
		t.Fatal(err)
	}
	nou, _, err := m.serviciu.ConfiguratiePontaj(ctx)
	if err != nil || len(nou.Conditionate) != 1 || nou.Conditionate[0].GrupID != grup.ID {
		t.Fatal("reguli nepersistate", nou, err)
	}
	if err := m.serviciu.SalveazaSchita(ctx, 2026, 10, []CelulaSchita{{AngajatID: a.ID, Data: "2026-10-03", Valoare: "G24"}}); err != nil {
		t.Fatal(err)
	}
	rez, _, err := CalculeazaPontaj(ctx, acces.New(m.db), 2026, 10, nil)
	if err != nil {
		t.Fatal(err)
	}
	if len(rez.Totaluri) != 1 || calcul.TextValoarePontaj(rez.Totaluri[0].Valori["ore_garda"]) != "48" {
		t.Fatal("regula nu foloseste profilul real", rez)
	}
	var norma int64 = 300
	c := salveaza("contracte", ElementNomenclator{Denumire: i18n.Text{"ro": "Contract 5h"}, Activ: true, Normare: "ZILNICA", NormaMinute: &norma})
	a.ContractID = c.ID
	if err := m.serviciu.SalveazaAngajat(ctx, a); err != nil {
		t.Fatal(err)
	}
	a = citeste(t, m.serviciu).Angajati[0]
	if a.NormaMinute != 300 || !a.NormaDefinita {
		t.Fatal(a)
	}
	c.Activ = false
	c = salveaza("contracte", c)
	b := angajatTest("ALT-CONTRACT")
	b.ContractID = c.ID
	verificaCod(t, m.serviciu.SalveazaAngajat(ctx, b), "nomenclator.inactiv")
	if err := m.serviciu.SalveazaAngajat(ctx, a); err != nil {
		t.Fatal("contract dezactivat existent", err)
	}
	fSeparata := salveaza("functii", ElementNomenclator{Denumire: i18n.Text{"ro": "Alta functie fara grup"}, Activ: true})
	copil := angajatTest("GARDA-SEPARATA")
	copil.AngajatPrincipalID = a.ID
	copil.ContractID = "GARDA"
	copil.FunctieID = fSeparata.ID
	if err := m.serviciu.SalveazaAngajat(ctx, copil); err != nil {
		t.Fatal(err)
	}
	profile, err := acces.New(m.db).ListaProfilePontaj(ctx)
	if err != nil {
		t.Fatal(err)
	}
	for _, p := range profile {
		if p.FunctieID == fSeparata.ID && p.GrupID != "" {
			t.Fatal("functia proprie fara grup a mostenit incorect grupul principal", p)
		}
	}

}

func TestRevizieNomenclatorMonotona(t *testing.T) {
	veche := "2026-10-31T10:00:00Z"
	for _, acum := range []string{veche, "2026-10-31T09:59:59Z"} {
		if r := revizieNomenclatorDupa(veche, acum); r != "2026-10-31T10:00:00.000000001Z" {
			t.Fatal(r)
		}
	}
	if r := revizieNomenclatorDupa(veche, "2026-10-31T10:00:01Z"); r != "2026-10-31T10:00:01Z" {
		t.Fatal(r)
	}
}

func TestSalvareGardaPastreazaFunctiaMostenitaInactiva(t *testing.T) {
	m := pregateste(t)
	ctx := context.Background()
	f := ElementNomenclator{ID: "f", Denumire: i18n.Text{"ro": "Functie inactiva"}, Activ: false}
	if _, err := m.db.Exec(`INSERT INTO functii(id,denumire,activ,modificat_la) VALUES('f','{"ro":"Functie inactiva"}',0,'x');INSERT INTO angajati(id,nume,prenume,marca,creat_la,modificat_la) VALUES('p','Exemplu','Ana','P','x','x'),('g','Exemplu','Ana','G','x','x');INSERT INTO angajati_functii VALUES('p','f');INSERT INTO angajati_garzi VALUES('g','p');INSERT INTO angajati_contracte VALUES('p','INTREAGA'),('g','GARDA')`); err != nil {
		t.Fatal(err)
	}
	var g Angajat
	for _, a := range citeste(t, m.serviciu).Angajati {
		if a.ID == "g" {
			g = a
		}
	}
	if g.FunctieID != f.ID {
		t.Fatal("functie nemostenita", g)
	}
	g.ConcediuZileAn = nil
	if err := m.serviciu.SalveazaAngajat(ctx, g); err != nil {
		t.Fatal("editarea neutra a garzii refuzata", err)
	}
	var n int
	if err := m.db.QueryRow("SELECT count(*) FROM angajati_functii WHERE angajat_id='g'").Scan(&n); err != nil || n != 0 {
		t.Fatal("functie mostenita fixata accidental", n, err)
	}
}
