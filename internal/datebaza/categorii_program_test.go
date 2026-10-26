package datebaza

import (
	"context"
	"pontix/internal/i18n"
	"testing"
)

func TestCategoriiProgramPersistentaConflictSiAtribuire(t *testing.T) {
	m := pregateste(t)
	ctx := context.Background()
	d := citeste(t, m.serviciu)
	if len(d.CategoriiProgram) != 6 {
		t.Fatalf("categorii initiale: %+v", d.CategoriiProgram)
	}
	programe, err := m.serviciu.ListaSchimburi(ctx)
	if err != nil {
		t.Fatal(err)
	}
	ids := map[string]string{}
	for _, p := range programe {
		ids[p.Cod] = p.ID
	}
	cat := CategorieProgram{Denumire: i18n.Text{"ro": " Combinat de test "}, SchimburiID: []string{ids["Z7"], ids["G24"]}, SchimbImplicitID: ids["Z7"], Activ: true}
	if err := m.serviciu.SalveazaCategorieProgram(ctx, cat); err != nil {
		t.Fatal(err)
	}
	for _, c := range citeste(t, m.serviciu).CategoriiProgram {
		if c.Denumire["ro"] == "Combinat de test" {
			cat = c
		}
	}
	if cat.ID == "" || cat.Revizie != "1" {
		t.Fatalf("categorie nesalvata: %+v", cat)
	}
	a := angajatTest("PROGRAM")
	a.CategorieProgramID = cat.ID
	if err := m.serviciu.SalveazaAngajat(ctx, a); err != nil {
		t.Fatal(err)
	}
	a = citeste(t, m.serviciu).Angajati[0]
	if a.CategorieProgramID != cat.ID {
		t.Fatalf("atribuire: %+v", a)
	}
	verificaCod(t, m.serviciu.SalveazaSchita(ctx, 2026, 10, []CelulaSchita{{AngajatID: a.ID, Data: "2026-10-01", Valoare: "Z7"}, {AngajatID: a.ID, Data: "2026-10-02", Valoare: "S7"}}), "program.program_nepermis")
	var partial int
	if err := m.db.QueryRow(`SELECT count(*) FROM pontaj_ciorne`).Scan(&partial); err != nil || partial != 0 {
		t.Fatalf("lot partial salvat: %d %v", partial, err)
	}
	if err := m.serviciu.SalveazaSchita(ctx, 2026, 10, []CelulaSchita{{AngajatID: a.ID, Data: "2026-10-01", Valoare: "Z7"}, {AngajatID: a.ID, Data: "2026-10-02", Valoare: "@M240"}}); err != nil {
		t.Fatal(err)
	}
	// O înregistrare istorică nu devine invalidă retroactiv prin atribuirea categoriei.
	if _, err := m.db.Exec(`INSERT INTO pontaj_ciorne(angajat_id,data,valoare,modificat_la) VALUES(?,'2026-10-03','S7','x')`, a.ID); err != nil {
		t.Fatal(err)
	}
	if err := m.serviciu.SalveazaSchita(ctx, 2026, 10, []CelulaSchita{{AngajatID: a.ID, Data: "2026-10-03", Valoare: "S7", Revizie: "1"}}); err != nil {
		t.Fatal(err)
	}
	invalida := cat
	invalida.SchimbImplicitID = ids["S7"]
	verificaCod(t, m.serviciu.SalveazaCategorieProgram(ctx, invalida), "program.implicit_invalid")
	cat.Activ = false
	if err := m.serviciu.SalveazaCategorieProgram(ctx, cat); err != nil {
		t.Fatal(err)
	}
	verificaCod(t, m.serviciu.SalveazaCategorieProgram(ctx, cat), "program.conflict")
	if err := m.serviciu.SalveazaAngajat(ctx, a); err != nil {
		t.Fatal("atribuirea veche inactiva nu trebuie pierduta:", err)
	}
	nou := angajatTest("NOU")
	nou.CategorieProgramID = cat.ID
	verificaCod(t, m.serviciu.SalveazaAngajat(ctx, nou), "program.categorie_inactiva")
	var n int
	if err := m.db.QueryRow(`SELECT count(*) FROM jurnal_modificari WHERE tabel='categorii_program' AND rand_id=?`, cat.ID).Scan(&n); err != nil || n != 2 {
		t.Fatalf("jurnal/rollback: %d %v", n, err)
	}
}

func TestPrecompletareCitesteAbsenteSiSchitaCurentaFaraScriere(t *testing.T) {
	m := pregateste(t)
	ctx := context.Background()
	a := angajatTest("PREFILL")
	a.CategorieProgramID = "program-dimineata"
	if err := m.serviciu.SalveazaAngajat(ctx, a); err != nil {
		t.Fatal(err)
	}
	a = citeste(t, m.serviciu).Angajati[0]
	if _, err := m.db.Exec(`INSERT INTO coduri_absenta(cod,denumire,categorie,creat_la,modificat_la) VALUES('EV_TEST','{"ro":"Liber special de test"}','SPECIAL','x','x')`); err != nil {
		t.Fatal(err)
	}
	for i, cod := range []string{"CO", "CM", "EV_TEST"} {
		if _, err := m.db.Exec(`INSERT INTO absente_zile(angajat_id,data,cod,activ,creat_la,modificat_la) VALUES(?,?,?,1,'x','x')`, a.ID, []string{"2026-10-05", "2026-10-06", "2026-10-09"}[i], cod); err != nil {
			t.Fatal(err)
		}
	}
	if _, err := m.db.Exec(`INSERT INTO sarbatori(data,denumire,creat_la,modificat_la) VALUES('2026-10-07','{"ro":"Liber personalizat"}','x','x')`); err != nil {
		t.Fatal(err)
	}
	mod := []CelulaSchita{{AngajatID: a.ID, Data: "2026-10-08", Valoare: "@M240"}}
	p, err := m.serviciu.PrevizualizeazaPrecompletare(ctx, 2026, 10, []string{a.ID}, mod)
	if err != nil {
		t.Fatal(err)
	}
	if len(p.Celule) != 17 || p.Omise["absenta"] != 3 || p.Omise["libera"] != 1 || p.Omise["existent"] != 1 {
		t.Fatalf("plan: %+v", p)
	}
	var n int
	if err := m.db.QueryRow(`SELECT count(*) FROM pontaj_ciorne`).Scan(&n); err != nil || n != 0 {
		t.Fatalf("previzualizarea a scris date: %d %v", n, err)
	}
}
