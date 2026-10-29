package ai

import (
	"strings"
	"testing"
)

func TestIntentiePlanRegresii(t *testing.T) {
	info := contextGrup()
	for n := range info.Angajati {
		info.Angajati[n].Departament = "Bloc operator"
	}
	// 40 de comenzi anonimizate RO/EN: selecții, luni, clarificări și interogări.
	cazuri := []struct {
		mesaj, grup, cod, inceput, sfarsit, zile, existente string
		recunoscut                                          bool
	}{
		{"Toti angajati in octombrie si noiembrie lucreaza de la 7la 14", "toti", "Z7", "2026-10-01", "2026-11-30", "", "", true},
		{"Toți angajații în octombrie și noiembrie lucrează de la 7 la 14", "toti", "Z7", "2026-10-01", "2026-11-30", "", "", true},
		{"toti angajatii Z7", "toti", "Z7", "2026-10-01", "2026-10-31", "", "", true},
		{"tot personalul S7", "toti", "S7", "2026-10-01", "2026-10-31", "", "", true},
		{"fiecare angajat Z7 in februarie", "toti", "Z7", "2026-02-01", "2026-02-28", "", "", true},
		{"toti angajatii Z7 in februarie 2028", "toti", "Z7", "2028-02-01", "2028-02-29", "", "", true},
		{"toti angajatii Z7 in decembrie si ianuarie", "toti", "Z7", "2026-12-01", "2027-01-31", "", "", true},
		{"sectia ATI Z7", "sectie", "Z7", "2026-10-01", "2026-10-31", "", "", true},
		{"sectia Chirurgie S7", "sectie", "S7", "2026-10-01", "2026-10-31", "", "", true},
		{"toti din sectia ATI pe Z7", "sectie", "Z7", "2026-10-01", "2026-10-31", "", "", true},
		{"categoria dimineață Z7", "categorie", "Z7", "2026-10-01", "2026-10-31", "", "", true},
		{"departamentul Bloc operator Z7", "departament", "Z7", "2026-10-01", "2026-10-31", "", "", true},
		{"Pop Ana Z7", "persoane", "Z7", "2026-10-01", "2026-10-31", "", "", true},
		{"Ana Pop S7", "persoane", "S7", "2026-10-01", "2026-10-31", "", "", true},
		{"M1 si M2 Z7", "persoane", "Z7", "2026-10-01", "2026-10-31", "", "", true},
		{"toti angajatii Z7 luni-vineri pastreaza", "toti", "Z7", "2026-10-01", "2026-10-31", "lucratoare", "pastreaza", true},
		{"toti angajatii Z7 toate zilele inlocuieste", "toti", "Z7", "2026-10-01", "2026-10-31", "toate", "inlocuieste", true},
		{"toti angajatii Z7 weekend doar celulele goale", "toti", "Z7", "2026-10-01", "2026-10-31", "weekend", "pastreaza", true},
		{"toti angajatii Z7 inclusiv weekend", "toti", "Z7", "2026-10-01", "2026-10-31", "toate", "", true},
		{"toti angajatii Z7 zile lucratoare", "toti", "Z7", "2026-10-01", "2026-10-31", "lucratoare", "", true},
		{"toti angajatii Z7 2026-10-05", "toti", "Z7", "2026-10-05", "2026-10-05", "", "", true},
		{"toti angajatii Z7 2026-10-30 pana la 2026-11-03", "toti", "Z7", "2026-10-30", "2026-11-03", "", "", true},
		{"toti angajatii pe 15 octombrie Z7", "toti", "Z7", "2026-10-15", "2026-10-15", "", "", true},
		{"toti angajatii pe 31 februarie Z7", "toti", "Z7", "", "", "", "", true},
		{"toti angajatii maine Z7", "toti", "Z7", "", "", "", "", true},
		{"all employees Z7 in october weekdays keep existing", "toti", "Z7", "2026-10-01", "2026-10-31", "lucratoare", "pastreaza", true},
		{"all staff S7 in november every day replace existing", "toti", "S7", "2026-11-01", "2026-11-30", "toate", "inlocuieste", true},
		{"all employees Z7 in october and november", "toti", "Z7", "2026-10-01", "2026-11-30", "", "", true},
		{"schedule M1 Z7 empty cells", "persoane", "Z7", "2026-10-01", "2026-10-31", "", "pastreaza", true},
		{"toti angajatii de la 9 la 17", "toti", "", "2026-10-01", "2026-10-31", "", "", true},
		{"planifica toti angajatii", "toti", "", "2026-10-01", "2026-10-31", "", "", true},
		{"sectia necunoscuta Z7", "sectie", "Z7", "2026-10-01", "2026-10-31", "", "", true},
		{"planifica Maria pe Z7", "", "Z7", "2026-10-01", "2026-10-31", "", "", true},
		{"toti angajatii Z7 in mai", "toti", "Z7", "2026-05-01", "2026-05-31", "", "", true},
		{"toti angajatii Z7 in martie si aprilie", "toti", "Z7", "2026-03-01", "2026-04-30", "", "", true},
		{"listeaza toti angajatii", "", "", "", "", "", "", false},
		{"cine lucreaza pe Z7", "", "", "", "", "", "", false},
		{"cati angajati sunt pe Z7", "", "", "", "", "", "", false},
		{"ajutor", "", "", "", "", "", "", false},
		{"who works Z7", "", "", "", "", "", "", false},
	}
	for _, c := range cazuri {
		t.Run(c.mesaj, func(t *testing.T) {
			i, ok := InterpreteazaPlan(c.mesaj, info)
			if ok != c.recunoscut {
				t.Fatalf("recunoscut=%v, intentie=%+v", ok, i)
			}
			if !ok {
				return
			}
			if i.Selectie.Tip != c.grup || i.CodSchimb != c.cod || i.DeLa != c.inceput || i.PanaLa != c.sfarsit || i.FiltruZile != c.zile || i.ModExistente != c.existente {
				t.Fatalf("intentie=%+v", i)
			}
			if (c.zile == "" || c.existente == "") && len(i.Clarificari) == 0 {
				t.Fatal("lipsesc clarificările")
			}
		})
	}
}

func TestContextSpecificNuTrimiteEchipa(t *testing.T) {
	info := contextGrup()
	for n := range info.Angajati {
		info.Angajati[n].ID = info.Angajati[n].Marca
	}
	info.Absente = []InfoAbsenta{{AngajatID: "M1", Data: "2026-10-05", Cod: "CO"}, {AngajatID: "M2", Data: "2026-10-06", Cod: "CO"}}
	info.TurePlanificate = []InfoPozitiePlanificare{{AngajatID: "M1", Data: "2026-10-01", CodSchimb: "Z7"}, {AngajatID: "M2", Data: "2026-10-02", CodSchimb: "S7"}}
	for _, c := range []struct {
		mesaj, domeniu string
		absente, ture  int
	}{
		{"ajutor", "ajutor", 0, 0}, {"absente M1", "absente", 1, 0}, {"pontaj M1", "pontaj", 1, 1}, {"planificare M1", "planificare", 1, 1},
	} {
		r := ContextPentruCerere(c.mesaj, info)
		if r.Domeniu != c.domeniu || len(r.Absente) != c.absente || len(r.TurePlanificate) != c.ture || !r.ContextIncomplet {
			t.Fatalf("%s: %+v", c.mesaj, r)
		}
		for _, a := range r.Angajati {
			if a.Marca != "M1" {
				t.Fatalf("persoană nerelevantă: %+v", a)
			}
		}
		if !strings.Contains(construiestePromptSistem(r), "CONTEXT PARȚIAL") {
			t.Fatal("nu semnalează contextul parțial")
		}
	}
}

func TestIntentieNuAscundeAmbiguitati(t *testing.T) {
	for _, mesaj := range []string{"toti angajatii exceptand M1 Z7", "sectia ATI categoria Garzi Z7"} {
		i, ok := InterpreteazaPlan(mesaj, contextGrup())
		if !ok || i.Selectie.Tip != "" {
			t.Fatalf("selecție dedusă pentru %q: %+v", mesaj, i)
		}
	}
	i, _ := InterpreteazaPlan("toti angajatii Z7 octombrie si decembrie", contextGrup())
	if i.DeLa != "" || i.PanaLa != "" {
		t.Fatalf("a introdus luna intermediară: %+v", i)
	}
	i, _ = InterpreteazaPlan("toti angajatii Z7 si S7", contextGrup())
	if i.CodSchimb != "" {
		t.Fatalf("a ales primul program: %+v", i)
	}
	if _, ok := InterpreteazaPlan("sterge Z7 pentru M1", contextGrup()); ok {
		t.Fatal("ștergerea a devenit planificare")
	}
}
