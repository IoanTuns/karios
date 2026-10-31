package datebaza

import (
	"context"
	"testing"
)

func TestAbsenteNumaiZileLucratoare(t *testing.T) {
	for _, cod := range []string{"CO", "CM"} {
		for _, data := range []string{"2026-10-03", "2026-10-04", "2026-12-01"} {
			t.Run(cod+data, func(t *testing.T) {
				m := pregateste(t)
				a := persoanaAbsente(t, m)
				verificaCod(t, m.serviciu.SalveazaAbsente(context.Background(), []ZiAbsenta{{AngajatID: a.ID, Data: data, Cod: cod, Activ: true}}), "absente.zi_nelucratoare")
				if len(absenteTest(t, m, 2026).Zile) != 0 {
					t.Fatal("absenta nelucratoare salvata")
				}
			})
		}
	}
}

func TestMedicalAnuleazaIntreagaTuraPeToateMarcile(t *testing.T) {
	m := pregateste(t)
	ctx := context.Background()
	a := persoanaAbsente(t, m)
	for _, marca := range []string{"G-TEST-1", "G-TEST-2", "ALTA"} {
		if err := m.serviciu.SalveazaAngajat(ctx, angajatTest(marca)); err != nil {
			t.Fatal(err)
		}
	}
	ids := []string{a.ID}
	alt := ""
	for _, p := range citeste(t, m.serviciu).Angajati {
		if p.Marca == "ALTA" {
			alt = p.ID
		}
		if p.Marca == "G-TEST-1" || p.Marca == "G-TEST-2" {
			if _, err := m.db.Exec("INSERT INTO angajati_garzi VALUES(?,?)", p.ID, a.ID); err != nil {
				t.Fatal(err)
			}
			ids = append(ids, p.ID)
		}
	}
	for _, id := range append(append([]string{}, ids...), alt) {
		if _, err := m.db.Exec(`INSERT INTO pontaj_ciorne(angajat_id,data,valoare,modificat_la) VALUES(?,'2026-10-04','N12','x'),(?,'2026-10-05','Z7','x')`, id, id); err != nil {
			t.Fatal(err)
		}
		if _, err := m.db.Exec(`INSERT INTO planificare(id,angajat_id,data,schimb_id,sectiune,creat_la,modificat_la) SELECT ?,?,'2026-10-04',id,'GARZI','x','x' FROM tipuri_schimb WHERE cod='N12'`, id, id); err != nil {
			t.Fatal(err)
		}
	}
	if err := m.serviciu.SalveazaAbsente(ctx, []ZiAbsenta{{AngajatID: a.ID, Data: "2026-10-05", Cod: "CO", Activ: true}}); err != nil {
		t.Fatal(err)
	}
	// Medical nou: nu cere revizia concediului pe care il inlocuieste prioritar.
	if err := m.serviciu.SalveazaAbsente(ctx, []ZiAbsenta{{AngajatID: a.ID, Data: "2026-10-05", Cod: "CM", Activ: true}}); err != nil {
		t.Fatal(err)
	}
	for _, id := range ids {
		var n int
		if err := m.db.QueryRow("SELECT count(*) FROM pontaj_ciorne WHERE angajat_id=? AND valoare<>''", id).Scan(&n); err != nil || n != 0 {
			t.Fatalf("schita %s: %d %v", id, n, err)
		}
		if err := m.db.QueryRow("SELECT count(*) FROM planificare WHERE angajat_id=? AND activ=1", id).Scan(&n); err != nil || n != 0 {
			t.Fatalf("plan %s: %d %v", id, n, err)
		}
		var rev string
		if err := m.db.QueryRow("SELECT revizie FROM pontaj_ciorne WHERE angajat_id=? AND data='2026-10-04'", id).Scan(&rev); err != nil {
			t.Fatal(err)
		}
		verificaCod(t, m.serviciu.SalveazaSchita(ctx, 2026, 10, []CelulaSchita{{AngajatID: id, Data: "2026-10-04", Valoare: "N12", Revizie: rev}}), "absente.celula_protejata")
	}
	var n int
	if err := m.db.QueryRow("SELECT count(*) FROM pontaj_ciorne WHERE angajat_id=? AND valoare<>''", alt).Scan(&n); err != nil || n != 2 {
		t.Fatalf("alta persoana modificata: %d %v", n, err)
	}
	zile := absenteTest(t, m, 2026).Zile
	if len(zile) != 1 || zile[0].Cod != "CM" || zile[0].AngajatID != a.ID {
		t.Fatal("medical duplicat sau CO ramas", zile)
	}
	z := zile[0]
	z.Cod = "CO"
	z.Revizie = ""
	verificaCod(t, m.serviciu.SalveazaAbsente(ctx, []ZiAbsenta{z}), "absente.medical_prioritar")
}

func TestMedicalTuraLunaAnterioaraBlocataRollback(t *testing.T) {
	m := pregateste(t)
	ctx := context.Background()
	a := persoanaAbsente(t, m)
	if _, err := m.db.Exec(`INSERT INTO pontaj_ciorne(angajat_id,data,valoare,modificat_la) VALUES(?,'2026-09-30','N12','x')`, a.ID); err != nil {
		t.Fatal(err)
	}
	if _, err := m.serviciu.SeteazaStareLuna(ctx, 2026, 9, "BLOCATA", "test", "operator"); err != nil {
		t.Fatal(err)
	}
	verificaCod(t, m.serviciu.SalveazaAbsente(ctx, []ZiAbsenta{{AngajatID: a.ID, Data: "2026-10-01", Cod: "CM", Activ: true}}), "pontaj.luna_inchisa")
	var v string
	if err := m.db.QueryRow("SELECT valoare FROM pontaj_ciorne WHERE angajat_id=?", a.ID).Scan(&v); err != nil || v != "N12" {
		t.Fatal("rollback incomplet", v, err)
	}
	if len(absenteTest(t, m, 2026).Zile) != 0 {
		t.Fatal("medical salvat partial")
	}
}

func TestMedicalNuAnuleazaTuraTerminataLaMiezulNoptii(t *testing.T) {
	m := pregateste(t)
	ctx := context.Background()
	a := persoanaAbsente(t, m)
	if _, err := m.db.Exec("UPDATE tipuri_schimb SET ora_start=1200,durata_minute=240 WHERE cod='N12'"); err != nil {
		t.Fatal(err)
	}
	if _, err := m.db.Exec(`INSERT INTO pontaj_ciorne(angajat_id,data,valoare,modificat_la) VALUES(?,'2026-10-04','N12','x')`, a.ID); err != nil {
		t.Fatal(err)
	}
	if err := m.serviciu.SalveazaAbsente(ctx, []ZiAbsenta{{AngajatID: a.ID, Data: "2026-10-05", Cod: "CM", Activ: true}}); err != nil {
		t.Fatal(err)
	}
	var v string
	if err := m.db.QueryRow("SELECT valoare FROM pontaj_ciorne WHERE angajat_id=?", a.ID).Scan(&v); err != nil || v != "N12" {
		t.Fatal("tura fara suprapunere anulata", v, err)
	}
}

func TestMedicalEditabilSiAnulabilExplicit(t *testing.T) {
	m := pregateste(t)
	ctx := context.Background()
	a := persoanaAbsente(t, m)
	if err := m.serviciu.SalveazaAbsente(ctx, []ZiAbsenta{{AngajatID: a.ID, Data: "2026-10-05", Cod: "CM", Activ: true}}); err != nil {
		t.Fatal(err)
	}
	z := absenteTest(t, m, 2026).Zile[0]
	veche := z
	z.Cod = "CO"
	if err := m.serviciu.SalveazaAbsente(ctx, []ZiAbsenta{z}); err != nil {
		t.Fatal("corectia explicita refuzata", err)
	}
	verificaCod(t, m.serviciu.SalveazaAbsente(ctx, []ZiAbsenta{veche}), "absente.conflict")
	z = absenteTest(t, m, 2026).Zile[0]
	if z.Cod != "CO" {
		t.Fatal(z)
	}
	z.Cod = "CM"
	if err := m.serviciu.SalveazaAbsente(ctx, []ZiAbsenta{z}); err != nil {
		t.Fatal(err)
	}
	z = absenteTest(t, m, 2026).Zile[0]
	z.Activ = false
	if err := m.serviciu.SalveazaAbsente(ctx, []ZiAbsenta{z}); err != nil {
		t.Fatal("anularea medicalului refuzata", err)
	}
	if absenteTest(t, m, 2026).Zile[0].Activ {
		t.Fatal("medical ramas activ")
	}
}
