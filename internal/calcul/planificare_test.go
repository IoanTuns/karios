package calcul

import (
	"testing"
	"time"
)

func TestSectiuniPlanificareSiTureNominale(t *testing.T) {
	s := ProgramPrecompletare{ID: "s", OraStart: 900, DurataMinute: 420, Activ: true, Tip: "TURA", Regim: "LUCRATOARE"}
	g := ProgramPrecompletare{ID: "g", OraStart: 420, DurataMinute: 1440, Activ: true, Tip: "GARDA"}
	if !ProgramDePlanificat(s) || !ProgramInSectiune(s, "DUPA_AMIAZA") || ProgramInSectiune(s, "GARZI") || !ProgramInSectiune(g, "GARZI") {
		t.Fatal("sectiuni incorecte")
	}
	if motiv := VerificaTuraPlanificata("a", "2026-10-31", g, map[string]bool{"a/2026-11-01": true}, nil, nil); motiv != "absenta" {
		t.Fatalf("trecere de luna: %s", motiv)
	}
	start, _ := time.Parse(time.RFC3339, "2026-10-01T07:00:00Z")
	intervale := []IntervalPlanificare{{AngajatID: "a", Cheie: "anterior", Inceput: start, Sfarsit: start.Add(8 * time.Hour)}}
	if motiv := VerificaTuraPlanificata("a", "2026-10-01", s, nil, intervale, nil); motiv != "" {
		t.Fatalf("intervale adiacente: %s", motiv)
	}
	intervale[0].Sfarsit = start.Add(9 * time.Hour)
	if motiv := VerificaTuraPlanificata("a", "2026-10-01", s, nil, intervale, nil); motiv != "suprapunere" {
		t.Fatalf("suprapunere: %s", motiv)
	}
	if motiv := VerificaTuraPlanificata("a", "2026-10-01", s, nil, intervale, map[string]bool{"anterior": true}); motiv != "" {
		t.Fatalf("inlocuire explicita: %s", motiv)
	}
}

func TestAnalizeazaAcoperireZileNelucratoare(t *testing.T) {
	sarbatori := map[string]string{
		"2026-10-14": "Sfânta Parascheva",
	}
	pozitii := []PozitiePlanificareSimple{
		{AngajatID: "a1", Data: "2026-10-03", CodSchimb: "G24", Activ: true},
		{AngajatID: "a2", Data: "2026-10-03", CodSchimb: "Z12", Activ: true},
		{AngajatID: "a3", Data: "2026-10-03", CodSchimb: "N12", Activ: true},
		{AngajatID: "a4", Data: "2026-10-03", CodSchimb: "IN9", Activ: true},
		{AngajatID: "a1", Data: "2026-10-04", CodSchimb: "G24", Activ: true},
	}

	rap := AnalizeazaAcoperire(2026, 10, sarbatori, pozitii, nil)
	if rap.TotalZileNelucratoare < 9 {
		t.Fatalf("zile nelucratoare insuficiente: %d", rap.TotalZileNelucratoare)
	}
	// 2026-10-03 e sâmbătă, acoperită complet
	var gasit03, gasit04 bool
	for _, z := range rap.Zile {
		if z.Data == "2026-10-03" {
			gasit03 = true
			if !z.Complet || len(z.Lipsuri) != 0 {
				t.Fatalf("ziua 03 trebuia sa fie completa: %+v", z)
			}
		}
		if z.Data == "2026-10-04" {
			gasit04 = true
			if z.Complet || len(z.Lipsuri) != 3 {
				t.Fatalf("ziua 04 trebuia sa aiba 3 lipsuri: %+v", z)
			}
		}
	}
	if !gasit03 || !gasit04 {
		t.Fatal("zilele de weekend nu au fost gasite")
	}
}

func TestAnalizeazaEchitateSarbatori(t *testing.T) {
	pozitii := []PozitiePlanificareSimple{
		{AngajatID: "a2", Data: "2026-12-31", CodSchimb: "G24", Activ: true}, // revelion
		{AngajatID: "a3", Data: "2026-12-25", CodSchimb: "G24", Activ: true}, // craciun
	}
	istoricAnTrecut := []PozitiePlanificareSimple{
		{AngajatID: "a2", Data: "2025-12-25", CodSchimb: "G24", Activ: true}, // craciun an trecut
	}
	istoricAnViitor := []PozitiePlanificareSimple{
		{AngajatID: "a3", Data: "2027-05-02", CodSchimb: "G24", Activ: true}, // pasti an viitor
	}
	sarbatori := map[string]string{
		"2027-05-02": "Paștele",
	}

	avertismente := AnalizeazaEchitate(2026, 12, pozitii, istoricAnTrecut, istoricAnViitor, sarbatori)
	if len(avertismente) != 2 {
		t.Fatalf("asteptam 2 avertismente, am primit %d: %+v", len(avertismente), avertismente)
	}

	tipuri := map[string]bool{}
	for _, a := range avertismente {
		tipuri[a.Tip] = true
	}
	if !tipuri["ROTATIE_CRACIUN_REVELION"] || !tipuri["CONFLICT_CRACIUN_PASTI"] {
		t.Fatalf("lipsesc tipuri de avertisment: %+v", tipuri)
	}
}
