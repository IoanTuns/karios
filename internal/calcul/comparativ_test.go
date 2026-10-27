package calcul

import (
	"testing"
)

func TestCalculeazaPlanificatVsRealizat(t *testing.T) {
	angajati := []AngajatComparativ{
		{ID: "ang1", Nume: "Popescu", Prenume: "Ion", Marca: "DEMO-001", Functie: "Medic primar", Activ: true},
		{ID: "ang2", Nume: "Dumitru", Prenume: "Andrei", Marca: "DEMO-002", Functie: "Medic specialist", Activ: true},
	}

	pozitii := []PozitieComparativ{
		{AngajatID: "ang1", Data: "2026-10-01", CodSchimb: "Z7", Ore: 7, Activ: true},
		{AngajatID: "ang1", Data: "2026-10-02", CodSchimb: "Z7", Ore: 7, Activ: true},
		{AngajatID: "ang2", Data: "2026-10-03", CodSchimb: "Z12", Ore: 12, Activ: true},
	}

	celule := []CelulaComparativ{
		// Popescu: 1 oct a făcut Z7 (identic), 2 oct a fost înlocuit cu G24 (modificat), 5 oct a lucrat Z7 neplanificat
		{AngajatID: "ang1", Data: "2026-10-01", Valoare: "Z7", Ore: 7},
		{AngajatID: "ang1", Data: "2026-10-02", Valoare: "G24", Ore: 24},
		{AngajatID: "ang1", Data: "2026-10-05", Valoare: "Z7", Ore: 7},
		// Dumitru: 3 oct nu a venit (neefectuat)
	}

	absente := []AbsentaComparativ{
		{AngajatID: "ang2", Data: "2026-10-03", Cod: "CM"},
	}

	rap := CalculeazaPlanificatVsRealizat(2026, 10, angajati, pozitii, celule, absente)

	if rap.An != 2026 || rap.Luna != 10 {
		t.Fatalf("perioadă incorectă: %d/%d", rap.An, rap.Luna)
	}
	if len(rap.Randuri) != 2 {
		t.Fatalf("așteptat 2 rânduri, obținut %d", len(rap.Randuri))
	}

	// Verificare Popescu
	rPopescu := rap.Randuri[0]
	if rPopescu.TotalOrePlanificate != 14 {
		t.Errorf("ore planificate Popescu: primit %d, așteptat 14", rPopescu.TotalOrePlanificate)
	}
	if rPopescu.TotalOreRealizate != 38 {
		t.Errorf("ore realizate Popescu: primit %d, așteptat 38", rPopescu.TotalOreRealizate)
	}
	if rPopescu.DiferentaOre != 24 {
		t.Errorf("diferență ore Popescu: primit %d, așteptat 24", rPopescu.DiferentaOre)
	}
	if rPopescu.Zile[0].Stare != "IDENTIC" {
		t.Errorf("ziua 1 Popescu: primit %s, așteptat IDENTIC", rPopescu.Zile[0].Stare)
	}
	if rPopescu.Zile[1].Stare != "MODIFICAT" {
		t.Errorf("ziua 2 Popescu: primit %s, așteptat MODIFICAT", rPopescu.Zile[1].Stare)
	}
	if rPopescu.Zile[4].Stare != "NEPLANIFICAT" {
		t.Errorf("ziua 5 Popescu: primit %s, așteptat NEPLANIFICAT", rPopescu.Zile[4].Stare)
	}

	// Verificare Dumitru
	rDumitru := rap.Randuri[1]
	if rDumitru.Zile[2].Stare != "ABSENTA" {
		t.Errorf("ziua 3 Dumitru: primit %s, așteptat ABSENTA", rDumitru.Zile[2].Stare)
	}
	if rap.TotalDiscrepante == 0 {
		t.Errorf("așteptat discrepanțe detectate")
	}
}
