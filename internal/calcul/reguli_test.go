package calcul

import (
	"testing"
)

func TestEvalueazaReguli(t *testing.T) {
	reguli := []RegulaConfigurabila{
		{
			ID:              "reg-1",
			Cod:             "REPAUS_G24",
			TipRegula:       "REPAUS_MINIM",
			Severitate:      "AVERTISMENT",
			CoduriSchimb:    "G24",
			ValoareNumerica: 24,
			Activ:           true,
		},
		{
			ID:              "reg-2",
			Cod:             "PLAFON_GARZI",
			TipRegula:       "PLAFON_LUNAR",
			Severitate:      "AVERTISMENT",
			CoduriSchimb:    "G24",
			ValoareNumerica: 2,
			Activ:           true,
		},
		{
			ID:              "reg-3",
			Cod:             "SUCCESIUNE_G24_Z7",
			TipRegula:       "SUCCESIUNE_INTERZISA",
			Severitate:      "BLOCANT",
			CoduriSchimb:    "G24,Z7",
			ValoareNumerica: 0,
			Activ:           true,
		},
	}

	pozitii := []PozitiePentruRegula{
		// Popescu: Gardă G24 pe 2026-10-01 (07:00-07:00), urmată imediat pe 2026-10-02 de Z7 (07:00-14:00) -> 0h repaus!
		{AngajatID: "ang1", Nume: "Popescu", Data: "2026-10-01", CodSchimb: "G24", OraStart: 420, DurataMin: 1440, Activ: true},
		{AngajatID: "ang1", Nume: "Popescu", Data: "2026-10-02", CodSchimb: "Z7", OraStart: 420, DurataMin: 420, Activ: true},

		// Ionescu: 3 gărzi în octombrie (plafon = 2)
		{AngajatID: "ang2", Nume: "Ionescu", Data: "2026-10-03", CodSchimb: "G24", OraStart: 420, DurataMin: 1440, Activ: true},
		{AngajatID: "ang2", Nume: "Ionescu", Data: "2026-10-10", CodSchimb: "G24", OraStart: 420, DurataMin: 1440, Activ: true},
		{AngajatID: "ang2", Nume: "Ionescu", Data: "2026-10-17", CodSchimb: "G24", OraStart: 420, DurataMin: 1440, Activ: true},
	}

	incalcari := EvalueazaReguli(reguli, pozitii)

	if len(incalcari) != 3 {
		t.Fatalf("așteptat 3 încălcări, obținut: %d (%+v)", len(incalcari), incalcari)
	}

	coduriGasite := make(map[string]bool)
	for _, inc := range incalcari {
		coduriGasite[inc.RegulaCod] = true
	}

	if !coduriGasite["REPAUS_G24"] {
		t.Errorf("lipsă încălcare REPAUS_G24")
	}
	if !coduriGasite["SUCCESIUNE_G24_Z7"] {
		t.Errorf("lipsă încălcare SUCCESIUNE_G24_Z7")
	}
	if !coduriGasite["PLAFON_GARZI"] {
		t.Errorf("lipsă încălcare PLAFON_GARZI")
	}
}
