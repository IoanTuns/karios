package datebaza

import (
	"context"
	"pontix/internal/i18n"
	"testing"
)

func configurareStergereTest(t *testing.T, m *mediuTest) (string, string) {
	t.Helper()
	ctx := context.Background()
	if err := m.serviciu.SalveazaElementNomenclator(ctx, "grupuri", ElementNomenclator{Denumire: i18n.Text{"ro": "Configuratie anonima test"}, Activ: false}); err != nil {
		t.Fatal(err)
	}
	for _, e := range citeste(t, m.serviciu).GrupuriProfesionale {
		if e.Denumire.In("ro") == "Configuratie anonima test" {
			return e.ID, e.Revizie
		}
	}
	t.Fatal("configuratie lipsa")
	return "", ""
}

func TestStergereConfigurareInactiveNeutilizataCuJurnal(t *testing.T) {
	m := pregateste(t)
	id, rev := configurareStergereTest(t, m)
	if err := m.serviciu.StergeConfigurare(context.Background(), "grupuri", id, rev); err != nil {
		t.Fatal(err)
	}
	var n int
	if err := m.db.QueryRow("SELECT count(*) FROM grupuri_profesionale WHERE id=?", id).Scan(&n); err != nil || n != 0 {
		t.Fatal(n, err)
	}
	if err := m.db.QueryRow("SELECT count(*) FROM jurnal_modificari WHERE tabel='grupuri_profesionale' AND rand_id=?", id).Scan(&n); err != nil || n < 2 {
		t.Fatal("jurnal lipsa", n, err)
	}
}

func TestStergereConfigurareProtectii(t *testing.T) {
	for _, caz := range []string{"activa", "revizie", "fk", "json", "protejat", "tip_invalid"} {
		t.Run(caz, func(t *testing.T) {
			m := pregateste(t)
			id, rev := configurareStergereTest(t, m)
			tip := "grupuri"
			cod := ""
			switch caz {
			case "activa":
				if _, err := m.db.Exec("UPDATE grupuri_profesionale SET activ=1 WHERE id=?", id); err != nil {
					t.Fatal(err)
				}
				cod = "stergere.activa"
			case "revizie":
				rev = "veche"
				cod = "stergere.conflict"
			case "fk":
				if _, err := m.db.Exec(`INSERT INTO functii(id,denumire,activ,modificat_la,grup_id) VALUES('f-test','{"ro":"Functie test"}',0,'x',?)`, id); err != nil {
					t.Fatal(err)
				}
				cod = "stergere.folosit"
			case "json":
				if _, err := m.db.Exec(`INSERT INTO setari(cheie,valoare,modificat_la) VALUES('test-ref',json_object('grup_id',?),'x')`, id); err != nil {
					t.Fatal(err)
				}
				cod = "stergere.folosit"
			case "protejat":
				tip = "tipuri_program"
				id = "NORMAL"
				rev = "initial"
				cod = "stergere.protejat"
			case "tip_invalid":
				tip = "angajati"
				cod = "stergere.selectie_invalida"
			}
			verificaCod(t, m.serviciu.StergeConfigurare(context.Background(), tip, id, rev), cod)
			var n int
			tabela := "grupuri_profesionale"
			if caz == "protejat" {
				tabela = "tipuri_program"
			}
			if err := m.db.QueryRow("SELECT count(*) FROM "+tabela+" WHERE id=?", id).Scan(&n); err != nil || n != 1 {
				t.Fatal("date pierdute la refuz", n, err)
			}
		})
	}
}

func TestStergereProgramSiProtectiePontaj(t *testing.T) {
	for _, folosit := range []bool{false, true} {
		t.Run(map[bool]string{false: "nefolosit", true: "pontaj"}[folosit], func(t *testing.T) {
			m := pregateste(t)
			ctx := context.Background()
			s := schimbTest()
			s.Cod = "STERGE-TEST"
			s.Activ = false
			if err := m.serviciu.SalveazaSchimb(ctx, s); err != nil {
				t.Fatal(err)
			}
			lista, err := m.serviciu.ListaSchimburi(ctx)
			if err != nil {
				t.Fatal(err)
			}
			for _, v := range lista {
				if v.Cod == s.Cod {
					s = v
				}
			}
			if folosit {
				a := persoanaAbsente(t, m)
				if _, err := m.db.Exec("INSERT INTO pontaj_ciorne(angajat_id,data,valoare,modificat_la) VALUES(?,'2026-10-05',?,'x')", a.ID, s.Cod); err != nil {
					t.Fatal(err)
				}
			}
			err = m.serviciu.StergeConfigurare(ctx, "schimburi", s.ID, s.Revizie)
			if folosit {
				verificaCod(t, err, "stergere.folosit")
				return
			}
			if err != nil {
				t.Fatal(err)
			}
			var n int
			if err := m.db.QueryRow("SELECT count(*) FROM clasificari_schimburi WHERE schimb_id=?", s.ID).Scan(&n); err != nil || n != 0 {
				t.Fatal("metadate orfane", n, err)
			}
		})
	}
}

func TestStergereCategoriePastreazaProgrameSiRollback(t *testing.T) {
	for _, rollback := range []bool{false, true} {
		t.Run(map[bool]string{false: "succes", true: "rollback"}[rollback], func(t *testing.T) {
			m := pregateste(t)
			ctx := context.Background()
			lista, err := m.serviciu.ListaSchimburi(ctx)
			if err != nil {
				t.Fatal(err)
			}
			schimb := lista[0]
			c := CategorieProgram{Denumire: i18n.Text{"ro": "Categorie eliminabila test"}, SchimburiID: []string{schimb.ID}, Activ: false}
			if err := m.serviciu.SalveazaCategorieProgram(ctx, c); err != nil {
				t.Fatal(err)
			}
			for _, v := range citeste(t, m.serviciu).CategoriiProgram {
				if v.Denumire.In("ro") == "Categorie eliminabila test" {
					c = v
				}
			}
			if rollback {
				if _, err := m.db.Exec("CREATE TRIGGER refuz_test BEFORE DELETE ON categorii_program BEGIN SELECT RAISE(ABORT,'refuz');END"); err != nil {
					t.Fatal(err)
				}
			}
			err = m.serviciu.StergeConfigurare(ctx, "categorii_program", c.ID, c.Revizie)
			if rollback {
				if err == nil {
					t.Fatal("eroare nesemnalata")
				}
			} else if err != nil {
				t.Fatal(err)
			}
			var n int
			asteptat := 0
			if rollback {
				asteptat = 1
			}
			if err := m.db.QueryRow("SELECT count(*) FROM categorii_program_schimburi WHERE categorie_id=?", c.ID).Scan(&n); err != nil || n != asteptat {
				t.Fatal("lista proprie incorecta", n, err)
			}
			if err := m.db.QueryRow("SELECT count(*) FROM tipuri_schimb WHERE id=?", schimb.ID).Scan(&n); err != nil || n != 1 {
				t.Fatal("program afectat", n, err)
			}
		})
	}
}

func TestStergereConfigurareReferintaAI(t *testing.T) {
	m := pregateste(t)
	id, rev := configurareStergereTest(t, m)
	if _, err := m.db.Exec(`INSERT INTO sabloane_ai(id,nume,intentie,modificat_la) VALUES('test-ai','Test anonim',json_object('grup_id',?),'x')`, id); err != nil {
		t.Fatal(err)
	}
	verificaCod(t, m.serviciu.StergeConfigurare(context.Background(), "grupuri", id, rev), "stergere.folosit")
}

func TestStergereConfigurareTipAbsentaSiModel(t *testing.T) {
	for _, tip := range []string{"tipuri_absenta", "modele_raport"} {
		t.Run(tip, func(t *testing.T) {
			m := pregateste(t)
			var tabela, cheie string
			if tip == "tipuri_absenta" {
				tabela, cheie = "coduri_absenta", "cod"
				if _, err := m.db.Exec(`INSERT INTO coduri_absenta(cod,denumire,activ,creat_la,modificat_la) VALUES('TEST-STERGERE','{"ro":"Absenta anonima"}',0,'x','x')`); err != nil {
					t.Fatal(err)
				}
			} else {
				tabela, cheie = "definitii_rapoarte", "id"
				if _, err := m.db.Exec(`INSERT INTO definitii_rapoarte(id,denumire,antet,continut,subsol,pagina,activ,creat_la,modificat_la) SELECT 'TEST-STERGERE',denumire,antet,continut,subsol,pagina,0,'x','x' FROM definitii_rapoarte WHERE id='raport-pontaj'`); err != nil {
					t.Fatal(err)
				}
			}
			rev := "x"
			if tip == "modele_raport" {
				rev = "1"
			}
			if err := m.serviciu.StergeConfigurare(context.Background(), tip, "TEST-STERGERE", rev); err != nil {
				t.Fatal(err)
			}
			var n int
			if err := m.db.QueryRow("SELECT count(*) FROM " + tabela + " WHERE " + cheie + "='TEST-STERGERE'").Scan(&n); err != nil || n != 0 {
				t.Fatal(n, err)
			}
		})
	}
}
