package datebaza

import (
	"context"
	"testing"
)

func TestPublicareConflictPontajEfectivSiAudit(t *testing.T) {
	m, id, sc := mediuPlanificare(t)
	ctx := context.Background()
	p := planTest(t, m, 2026, 10, PozitiePlanificare{AngajatID: id, Data: "2026-10-01", SchimbID: sc["S7"], Sectiune: "DUPA_AMIAZA", Activ: true})
	if _, err := m.db.Exec(`INSERT INTO pontaj(id,angajat_id,data,start,sfarsit,creat_la,modificat_la) VALUES('ef-test',?,'2026-10-01','2026-10-01T05:00:00Z','2026-10-01T12:00:00Z','x','x')`, id); err != nil {
		t.Fatal(err)
	}
	c := cerereTest(p, "DUPA_AMIAZA")
	r := previzualizareTest(t, m, &c)
	if r.Randuri[0].PontajeEfective != 1 || r.Randuri[0].Stare != "INLOCUIRE" {
		t.Fatal(r)
	}
	if _, err := m.db.Exec("UPDATE pontaj SET sfarsit='2026-10-01T13:00:00Z' WHERE id='ef-test'"); err != nil {
		t.Fatal(err)
	}
	_, err := m.serviciu.PublicaPlanificare(ctx, c)
	verificaCod(t, err, "planificare.conflict_publicare")
	previzualizareTest(t, m, &c)
	if _, err := m.serviciu.PublicaPlanificare(ctx, c); err != nil {
		t.Fatal(err)
	}
	var n int
	if err := m.db.QueryRow("SELECT count(*) FROM pontaj WHERE id='ef-test' AND sters_la IS NOT NULL AND sfarsit='2026-10-01T13:00:00Z'").Scan(&n); err != nil || n != 1 {
		t.Fatal(n, err)
	}
	if err := m.db.QueryRow("SELECT count(*) FROM jurnal_modificari WHERE tabel='pontaj' AND rand_id='ef-test'").Scan(&n); err != nil || n != 1 {
		t.Fatal("audit", n, err)
	}
}
func TestMedicalPontajEfectivLunaAnterioaraRollback(t *testing.T) {
	m := pregateste(t)
	a := persoanaAbsente(t, m)
	ctx := context.Background()
	if _, err := m.db.Exec(`INSERT INTO pontaj(id,angajat_id,data,start,sfarsit,creat_la,modificat_la) VALUES('ef-test',?,'2026-09-30','2026-09-30T19:00:00Z','2026-10-01T05:00:00Z','x','x')`, a.ID); err != nil {
		t.Fatal(err)
	}
	if _, err := m.serviciu.SeteazaStareLuna(ctx, 2026, 9, "BLOCATA", "test", "operator"); err != nil {
		t.Fatal(err)
	}
	verificaCod(t, m.serviciu.SalveazaAbsente(ctx, []ZiAbsenta{{AngajatID: a.ID, Data: "2026-10-01", Cod: "CM", Activ: true}}), "pontaj.luna_inchisa")
	var n int
	if err := m.db.QueryRow("SELECT count(*) FROM pontaj WHERE id='ef-test' AND sters_la IS NULL").Scan(&n); err != nil || n != 1 {
		t.Fatal(n, err)
	}
	if len(absenteTest(t, m, 2026).Zile) != 0 {
		t.Fatal("medical partial")
	}
}

func TestMedicalPontajEfectivLimitaMiezulNoptii(t *testing.T) {
	m := pregateste(t)
	a := persoanaAbsente(t, m)
	if _, err := m.db.Exec(`INSERT INTO pontaj(id,angajat_id,data,start,sfarsit,creat_la,modificat_la) VALUES('limita',?,'2026-10-25','2026-10-25T18:00:00Z','2026-10-25T22:00:00Z','x','x')`, a.ID); err != nil {
		t.Fatal(err)
	}
	if err := m.serviciu.SalveazaAbsente(context.Background(), []ZiAbsenta{{AngajatID: a.ID, Data: "2026-10-26", Cod: "CM", Activ: true}}); err != nil {
		t.Fatal(err)
	}
	var n int
	if err := m.db.QueryRow("SELECT count(*) FROM pontaj WHERE id='limita' AND sters_la IS NULL").Scan(&n); err != nil || n != 1 {
		t.Fatal("tura adiacenta retrasa", n, err)
	}
}
func TestMedicalPontajEfectivRollbackAudit(t *testing.T) {
	m := pregateste(t)
	a := persoanaAbsente(t, m)
	if _, err := m.db.Exec(`INSERT INTO pontaj(id,angajat_id,data,start,sfarsit,creat_la,modificat_la) VALUES('audit-refuz',?,'2026-10-05','2026-10-05T04:00:00Z','2026-10-05T11:00:00Z','x','x')`, a.ID); err != nil {
		t.Fatal(err)
	}
	if _, err := m.db.Exec(`CREATE TRIGGER refuz_audit_ef BEFORE INSERT ON jurnal_modificari WHEN NEW.tabel='pontaj' BEGIN SELECT RAISE(ABORT,'test');END`); err != nil {
		t.Fatal(err)
	}
	if err := m.serviciu.SalveazaAbsente(context.Background(), []ZiAbsenta{{AngajatID: a.ID, Data: "2026-10-05", Cod: "CM", Activ: true}}); err == nil {
		t.Fatal("auditul trebuia sa refuze")
	}
	var n int
	if err := m.db.QueryRow("SELECT count(*) FROM pontaj WHERE id='audit-refuz' AND sters_la IS NULL").Scan(&n); err != nil || n != 1 {
		t.Fatal(n, err)
	}
	if len(absenteTest(t, m, 2026).Zile) != 0 {
		t.Fatal("medical partial")
	}
}
