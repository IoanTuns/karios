package datebaza

import (
	"context"
	"testing"
)

func TestStareLunaSiBlocareEditare(t *testing.T) {
	m := pregateste(t)
	ctx := context.Background()

	// 1. Verificare stare inițială implicită: DESCHISA
	st, err := m.serviciu.CitesteStareLuna(ctx, 2026, 10)
	if err != nil {
		t.Fatalf("eroare citire stare lună: %v", err)
	}
	if st.Stare != "DESCHISA" {
		t.Errorf("așteptat stare inițială DESCHISA, obținut: %s", st.Stare)
	}

	inchisa, err := m.serviciu.EsteLunaInchisa(ctx, 2026, 10)
	if err != nil || inchisa {
		t.Errorf("luna nu ar trebui să fie închisă: inchisa=%v, err=%v", inchisa, err)
	}

	// 2. Setare stare: INCHISA (ex: Transmis la salarii)
	stMod, err := m.serviciu.SeteazaStareLuna(ctx, 2026, 10, "INCHISA", "Transmis la salarizare", "Admin")
	if err != nil {
		t.Fatalf("eroare setare stare lună: %v", err)
	}
	if stMod.Stare != "INCHISA" || stMod.Motiv != "Transmis la salarizare" {
		t.Errorf("stare neactualizată corect: %+v", stMod)
	}

	inchisa, _ = m.serviciu.EsteLunaInchisa(ctx, 2026, 10)
	if !inchisa {
		t.Errorf("așteptat luna închisă=true")
	}

	// 3. Verificare respingere salvare schiță pontaj pe lună închisă
	errSalvare := m.serviciu.SalveazaSchita(ctx, 2026, 10, []CelulaSchita{
		{AngajatID: "ang1", Data: "2026-10-01", Valoare: "Z7"},
	})
	verificaCod(t, errSalvare, "pontaj.luna_inchisa")

	// 4. Verificare respingere salvare planificare pe lună închisă
	errPlan := m.serviciu.SalveazaPlanificare(ctx, 2026, 10, []PozitiePlanificare{
		{AngajatID: "ang1", Data: "2026-10-01", SchimbID: "s1", Sectiune: "GARZI", Activ: true},
	})
	verificaCod(t, errPlan, "pontaj.luna_inchisa")

	// 5. Redeschidere lună: DESCHISA
	stRedeschis, err := m.serviciu.SeteazaStareLuna(ctx, 2026, 10, "DESCHISA", "Corecție necesară", "Admin")
	if err != nil {
		t.Fatalf("eroare redeschidere: %v", err)
	}
	if stRedeschis.Stare != "DESCHISA" {
		t.Errorf("așteptat stare DESCHISA, obținut: %s", stRedeschis.Stare)
	}

	inchisa, _ = m.serviciu.EsteLunaInchisa(ctx, 2026, 10)
	if inchisa {
		t.Errorf("luna ar trebui să fie deschisă după redeschidere")
	}
}

func TestAmprentaLaInchidereSiJurnal(t *testing.T) {
	m := pregateste(t)
	ctx := context.Background()
	if err := m.serviciu.SalveazaAngajat(ctx, angajatTest("M1")); err != nil {
		t.Fatal(err)
	}
	id := citeste(t, m.serviciu).Angajati[0].ID
	if err := m.serviciu.SalveazaSchita(ctx, 2026, 10, []CelulaSchita{{AngajatID: id, Data: "2026-10-01", Valoare: "D8"}}); err != nil {
		t.Fatal(err)
	}
	st, err := m.serviciu.SeteazaStareLuna(ctx, 2026, 10, "INCHISA", "Transmis", "Admin")
	if err != nil || len(st.Amprenta) != 19 || !st.AmprentaCorespunde {
		t.Fatalf("închidere: %+v, %v", st, err)
	}
	citita, err := m.serviciu.CitesteStareLuna(ctx, 2026, 10)
	if err != nil || citita.Amprenta != st.Amprenta || !citita.AmprentaCorespunde {
		t.Fatalf("citire: %+v, %v", citita, err)
	}
	// O modificare făcută pe lângă aplicație schimbă datele, nu și amprenta păstrată.
	if _, err := m.db.Exec(`UPDATE pontaj_ciorne SET valoare = 'N12' WHERE data = '2026-10-01'`); err != nil {
		t.Fatal(err)
	}
	if citita, _ = m.serviciu.CitesteStareLuna(ctx, 2026, 10); citita.AmprentaCorespunde {
		t.Fatal("amprenta corespunde după modificarea datelor")
	}
	blocata, err := m.serviciu.SeteazaStareLuna(ctx, 2026, 10, "BLOCATA", "", "Admin")
	if err != nil || blocata.Amprenta != st.Amprenta {
		t.Fatalf("blocarea unei luni închise păstrează amprenta: %+v, %v", blocata, err)
	}
	deschisa, err := m.serviciu.SeteazaStareLuna(ctx, 2026, 10, "DESCHISA", "Corecție", "Admin")
	if err != nil || deschisa.Amprenta != "" {
		t.Fatalf("redeschidere: %+v, %v", deschisa, err)
	}
	var n int
	if err := m.db.QueryRow(`SELECT count(*) FROM jurnal_modificari WHERE tabel = 'stare_luni' AND rand_id = '2026-10'`).Scan(&n); err != nil || n != 3 {
		t.Fatalf("intrări de jurnal pentru stare_luni: %d, %v", n, err)
	}
}
