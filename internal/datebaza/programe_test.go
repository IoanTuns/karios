package datebaza

import (
	"context"
	"io/fs"
	"path/filepath"
	"testing"
	"testing/fstest"

	"pontix/internal/bd"
)

func TestProgrameConfirmateSiIntrareSeparata(t *testing.T) {
	m := pregateste(t)
	lista, err := m.serviciu.ListaSchimburi(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	dupaCod := map[string]Schimb{}
	for _, s := range lista {
		dupaCod[s.Cod] = s
	}
	cazuri := []struct {
		cod                    string
		start, durata, sfarsit int64
		regim                  string
	}{
		{"Z7", 420, 420, 840, "LUCRATOARE"},
		{"S7", 900, 420, 1320, "LUCRATOARE"},
		{"N9", 1320, 540, 420, "LUCRATOARE"},
		{"IN9", 1320, 540, 420, "INTRARE_ZI_LIBERA"},
		{"Z12", 420, 720, 1140, "NELUCRATOARE"},
		{"N12", 1140, 720, 420, "NELUCRATOARE"},
		{"G24", 420, 1440, 420, "NELUCRATOARE"},
	}
	for _, c := range cazuri {
		s, ok := dupaCod[c.cod]
		if !ok || s.OraStart != c.start || s.DurataMinute != c.durata || (s.OraStart+s.DurataMinute)%1440 != c.sfarsit || s.Regim != c.regim || s.Categorie != "LAB" || !s.Activ {
			t.Fatalf("program %s: %+v", c.cod, s)
		}
	}
	if dupaCod["N9"].ID == dupaCod["IN9"].ID {
		t.Fatal("intrarea specială nu este distinctă")
	}
	if dupaCod["S7"].OraStart-(dupaCod["Z7"].OraStart+dupaCod["Z7"].DurataMinute) != 60 {
		t.Fatal("golul 14–15 a fost eliminat")
	}
	for _, codVechi := range []string{"7", "12", "S1", "S3"} {
		if dupaCod[codVechi].Categorie != "CUSTOM" {
			t.Fatalf("codul %s nu este încadrat în CUSTOM: %+v", codVechi, dupaCod[codVechi])
		}
	}
	if dupaCod["12"].OraStart != 840 || dupaCod["S3"].DurataMinute != 480 || dupaCod["7"].Regim != "GENERAL" {
		t.Fatal("codurile vechi au fost reinterpretate")
	}
	s := dupaCod["IN9"]
	s.Denumire["ro"] = "Intrare personalizată"
	s.Regim = "NELUCRATOARE"
	s.Categorie = "CUSTOM"
	if err := m.serviciu.SalveazaSchimb(context.Background(), s); err != nil {
		t.Fatal(err)
	}
	if err := m.serviciu.SchimbaActivSchimb(context.Background(), s.ID, false); err != nil {
		t.Fatal(err)
	}
	if n, err := m.serviciu.CompleteazaSchimburiInitiale(context.Background(), "LAB"); err != nil || n != 0 {
		t.Fatalf("completare: %d %v", n, err)
	}
	_, err = m.serviciu.CompleteazaSchimburiInitiale(context.Background(), "CUSTOM")
	verificaCod(t, err, "schimburi.categorie_invalida")
	lista, err = m.serviciu.ListaSchimburi(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	for _, actual := range lista {
		if actual.ID == s.ID && (actual.Regim != "NELUCRATOARE" || actual.Categorie != "CUSTOM" || actual.Activ || actual.Denumire["ro"] != "Intrare personalizată") {
			t.Fatal("editare/dezactivare pierdută")
		}
	}
	s.Regim = "INVENTAT"
	verificaCod(t, m.serviciu.SalveazaSchimb(context.Background(), s), "schimburi.regim_invalid")
	s.Regim = "NELUCRATOARE"
	s.Categorie = "INVENTATA"
	verificaCod(t, m.serviciu.SalveazaSchimb(context.Background(), s), "schimburi.categorie_invalida")
}

func TestMigrareProgramePastreazaPersonalizariSiPontaj(t *testing.T) {
	m := pregateste(t)
	ctx := context.Background()
	db, err := bd.Deschide(filepath.Join(t.TempDir(), "programe-vechi.db"), m.cheie)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	surse := fstest.MapFS{}
	intrari, err := fs.ReadDir(bd.Migrari, ".")
	if err != nil {
		t.Fatal(err)
	}
	for _, r := range intrari {
		if r.Name() >= "0009_" {
			continue
		}
		continut, err := fs.ReadFile(bd.Migrari, r.Name())
		if err != nil {
			t.Fatal(err)
		}
		surse[r.Name()] = &fstest.MapFile{Data: continut}
	}
	if _, err := bd.Migreaza(ctx, db, surse); err != nil {
		t.Fatal(err)
	}
	for _, comanda := range []string{
		`UPDATE tipuri_schimb SET ora_start=1200, durata_minute=660, activ=0, modificat_la='editat' WHERE cod='N12'`,
		`UPDATE tipuri_schimb SET sters_la='sters' WHERE cod='G24'`,
		`INSERT INTO tipuri_schimb(id,cod,denumire,ora_start,durata_minute,tip,activ,creat_la,modificat_la) VALUES ('intrare-config','in9','{"ro":"Intrare existentă"}',1260,600,'NOAPTE',0,'initial','editat')`,
		`INSERT INTO coduri_absenta(cod,denumire,categorie,activ,creat_la,modificat_la) VALUES ('Z12','{"ro":"Cod ocupat"}','SPECIAL',1,'initial','initial')`,
		`INSERT INTO angajati(id,nume,marca,creat_la,modificat_la) VALUES ('a','Exemplu','P001','initial','initial')`,
		`INSERT INTO pontaj_ciorne(angajat_id,data,valoare,modificat_la) VALUES ('a','2026-10-23','IN9','revizie-veche')`,
		`INSERT INTO pontaj(id,angajat_id,data,start,sfarsit,tip_schimb_id,creat_la,modificat_la) VALUES ('p','a','2026-10-22','2026-10-22T16:00:00Z','2026-10-23T04:00:00Z','10000000-0000-4000-8000-000000000003','initial','initial')`,
	} {
		if _, err := db.Exec(comanda); err != nil {
			t.Fatal(err)
		}
	}
	migrare, err := fs.ReadFile(bd.Migrari, "0009_programe_lucru.sql")
	if err != nil {
		t.Fatal(err)
	}
	surse["0009_programe_lucru.sql"] = &fstest.MapFile{Data: migrare}
	if n, err := bd.Migreaza(ctx, db, surse); err != nil || n != 1 {
		t.Fatalf("migrare: %d %v", n, err)
	}
	if _, err := bd.Migreaza(ctx, db, bd.Migrari); err != nil {
		t.Fatal(err)
	}
	lista, err := Nou(db).ListaSchimburi(ctx)
	if err != nil {
		t.Fatal(err)
	}
	for _, s := range lista {
		switch s.Cod {
		case "N12":
			if s.OraStart != 1200 || s.DurataMinute != 660 || s.Activ || s.Regim != "GENERAL" || s.Categorie != "LAB" {
				t.Fatalf("N12 personalizat schimbat: %+v", s)
			}
		case "in9":
			if s.ID != "intrare-config" || s.OraStart != 1260 || s.Regim != "GENERAL" || s.Categorie != "LAB" || s.Activ {
				t.Fatalf("intrare existentă schimbată: %+v", s)
			}
		case "7", "12", "S1", "S3":
			if s.Categorie != "CUSTOM" {
				t.Fatalf("cod vechi neîncadrat în CUSTOM: %+v", s)
			}
		case "IN9", "Z12", "G24":
			t.Fatalf("cod recreat peste configurație/absență/ștergere: %+v", s)
		}
	}
	if n, err := Nou(db).CompleteazaSchimburiInitiale(ctx); err != nil || n != 0 {
		t.Fatalf("completare peste coduri existente: %d %v", n, err)
	}
	var valoare, revizie, referinta string
	if err := db.QueryRow(`SELECT valoare,modificat_la FROM pontaj_ciorne WHERE angajat_id='a'`).Scan(&valoare, &revizie); err != nil || valoare != "IN9" || revizie != "revizie-veche" {
		t.Fatalf("schiță schimbată: %s %s %v", valoare, revizie, err)
	}
	if err := db.QueryRow(`SELECT tip_schimb_id FROM pontaj WHERE id='p'`).Scan(&referinta); err != nil || referinta != "10000000-0000-4000-8000-000000000003" {
		t.Fatalf("referință schimbată: %s %v", referinta, err)
	}
}
