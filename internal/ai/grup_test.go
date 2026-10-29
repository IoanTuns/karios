package ai

import "testing"

func contextGrup() ContextPontix {
	return ContextPontix{
		An: 2026, Luna: 10,
		Angajati: []InfoAngajat{
			{Marca: "M1", Nume: "Pop", Prenume: "Ana", Sectie: "ATI", Categorie: "Dimineață"},
			{Marca: "M2", Nume: "Ion", Prenume: "Dan", Sectie: "ATI", Categorie: "Gărzi"},
			{Marca: "M3", Nume: "Rus", Prenume: "Ema", Sectie: "Chirurgie", Categorie: "Dimineață"},
		},
		Schimburi: []InfoSchimb{
			{Cod: "Z7", Denumire: "Zi 07–14", Ore: 7, Interval: "07:00–14:00"},
			{Cod: "S7", Denumire: "Seară 15–22", Ore: 7, Interval: "15:00–22:00"},
			{Cod: "G24", Denumire: "Gardă", Ore: 24, Interval: "08:00–08:00"},
		},
	}
}

func TestComandaGrup(t *testing.T) {
	cazuri := []struct {
		mesaj    string
		marci    int
		cod      string
		filtru   string
		inceput  string
		recunosc bool
	}{
		{"Toti angajatii pe luna octombrie sunt pe zi de la 7 la 14", 3, "Z7", "lucratoare", "2026-10-01", true},
		{"Toți angajații în noiembrie pe S7", 3, "S7", "lucratoare", "2026-11-01", true},
		{"Secția ATI de la 15 la 22 toate zilele", 2, "S7", "toate", "2026-10-01", true},
		{"Pune categoria dimineață pe Z7 inclusiv weekend", 2, "Z7", "toate", "2026-10-01", true},
		{"tot personalul G24 in weekend", 3, "G24", "weekend", "2026-10-01", true},
		{"Toti angajatii pe 15 octombrie pe Z7", 0, "", "", "", false}, // zi anume: rămâne la model
		{"Toti angajatii de la 9 la 17", 0, "", "", "", false},         // niciun schimb cu acest interval
		{"Pune-l pe Pop pe Z7 toata luna", 0, "", "", "", false},       // un singur angajat
		{"Listeaza toti angajatii", 0, "", "", "", false},              // fără schimb
	}
	for _, c := range cazuri {
		r, ok := AnalizeazaComandaGrup(c.mesaj, contextGrup())
		if ok != c.recunosc {
			t.Errorf("%q: recunoscută=%v", c.mesaj, ok)
			continue
		}
		if !ok {
			continue
		}
		if len(r.Propuneri) != 1 {
			t.Fatalf("%q: propuneri %+v", c.mesaj, r)
		}
		p := r.Propuneri[0]
		if p.Tip != TipPlanificaGrup || len(p.Parametri["marci"].([]any)) != c.marci || p.Parametri["cod_schimb"] != c.cod ||
			p.Parametri["filtru_zile"] != c.filtru || p.Parametri["data_inceput"] != c.inceput {
			t.Errorf("%q: %+v", c.mesaj, p.Parametri)
		}
	}
}

func TestComandaGrupSfarsitLuna(t *testing.T) {
	r, ok := AnalizeazaComandaGrup("toti angajatii in februarie pe Z7", contextGrup())
	if !ok || r.Propuneri[0].Parametri["data_sfarsit"] != "2026-02-28" {
		t.Fatalf("%v %+v", ok, r)
	}
}
