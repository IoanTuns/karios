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
