package conturi

import (
	"context"
	"database/sql"
	"errors"
	"path/filepath"
	"testing"
	"time"

	"pontix/internal/audit"
	"pontix/internal/bd"
	"pontix/internal/i18n"
	"pontix/internal/seif"
)

var ctx = context.Background()

type mediu struct {
	s     *Serviciu
	db    *sql.DB
	cale  string
	cheie []byte
	seif  string
	ceas  time.Time
}

func pregateste(t *testing.T) *mediu {
	t.Helper()
	dir := t.TempDir()
	m := &mediu{cale: filepath.Join(dir, "pontix.db"), seif: filepath.Join(dir, "seif"), ceas: time.Date(2026, 10, 29, 9, 0, 0, 0, time.UTC)}
	sf, err := seif.Creeaza(m.seif, bd.LungimeCheie)
	if err != nil {
		t.Fatal(err)
	}
	m.cheie = sf.CheieDate
	m.db, err = bd.Deschide(m.cale, m.cheie)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { m.db.Close() })
	if _, err := bd.Migreaza(ctx, m.db, bd.Migrari); err != nil {
		t.Fatal(err)
	}
	m.s = Nou(m.db, m.seif, nil)
	m.s.acum = func() time.Time { return m.ceas }
	return m
}

func cerere(nume, rol string) CerereCont {
	return CerereCont{Nume: nume, NumeAfisat: "Ana Pop", Rol: rol, Parola: "parola123",
		Intrebari: []string{"intrebare_1", "intrebare_4", "intrebare_7"}, Raspunsuri: []string{"Școala nr. 5", "Bobiță", "Brașov"}}
}

func cod(t *testing.T, err error, asteptat string) {
	t.Helper()
	var e *i18n.Eroare
	if !errors.As(err, &e) || e.Cod != asteptat {
		t.Fatalf("eroare %v; așteptat %s", err, asteptat)
	}
}

var aici = Identitate{SID: "S-1-5-21-1", Calculator: "PC-ATI"}

func TestCreareSiValidari(t *testing.T) {
	m := pregateste(t)
	c, err := m.s.Creeaza(ctx, cerere("ana", RolOperator))
	if err != nil || c.Rol != RolAdministrator {
		t.Fatalf("primul cont trebuie să fie administrator: %+v, %v", c, err)
	}
	_, err = m.s.Creeaza(ctx, cerere("ANA", RolOperator))
	cod(t, err, "cont.nume_existent")
	slab := cerere("dan", RolOperator)
	for _, p := range []string{"scurt1", "faracifre", "12345678"} {
		slab.Parola = p
		_, err = m.s.Creeaza(ctx, slab)
		cod(t, err, "cont.parola_slaba")
	}
	repetat := cerere("dan", RolOperator)
	repetat.Intrebari = []string{"intrebare_1", "intrebare_1", "intrebare_2"}
	_, err = m.s.Creeaza(ctx, repetat)
	cod(t, err, "cont.intrebari_invalide")
	op, err := m.s.Creeaza(ctx, cerere("dan", RolOperator))
	if err != nil || op.Rol != RolOperator {
		t.Fatalf("al doilea cont: %+v, %v", op, err)
	}
	if n, _ := m.s.Numar(ctx); n != 2 {
		t.Fatalf("număr conturi: %d", n)
	}
	sf, err := seif.Deschide(m.seif)
	if err != nil || len(sf.Copii()) != 5 {
		t.Fatalf("copiile din seif (A + P și R pentru două conturi): %v, %v", sf.Copii(), err)
	}
}

func TestAutentificareBlocareSiJurnal(t *testing.T) {
	m := pregateste(t)
	if _, err := m.s.Creeaza(ctx, cerere("ana", RolOperator)); err != nil {
		t.Fatal(err)
	}
	if c, err := m.s.Autentifica(ctx, " Ana ", "parola123", aici); err != nil || c.Nume != "ana" {
		t.Fatalf("autentificare: %+v, %v", c, err)
	}
	_, err := m.s.Autentifica(ctx, "nimeni", "parola123", aici)
	cod(t, err, "cont.date_gresite")
	for i := 0; i < pragBlocare; i++ {
		_, err = m.s.Autentifica(ctx, "ana", "gresit123", aici)
		cod(t, err, "cont.date_gresite")
	}
	_, err = m.s.Autentifica(ctx, "ana", "parola123", aici)
	cod(t, err, "cont.blocat")
	m.ceas = m.ceas.Add(61 * time.Second)
	if _, err := m.s.Autentifica(ctx, "ana", "parola123", aici); err != nil {
		t.Fatalf("după expirarea blocării: %v", err)
	}
	var esuate, reusite int
	m.db.QueryRow(`SELECT count(*) FROM jurnal_modificari WHERE tabel = 'conturi' AND camp = 'autentificare_esuata'`).Scan(&esuate)
	m.db.QueryRow(`SELECT count(*) FROM jurnal_modificari WHERE tabel = 'conturi' AND camp = 'autentificare' AND cine = 'ana'`).Scan(&reusite)
	if esuate != pragBlocare || reusite != 2 {
		t.Fatalf("jurnal: %d eșuate, %d reușite", esuate, reusite)
	}
}

func TestBlocareProgresiva(t *testing.T) {
	m := pregateste(t)
	c, _ := m.s.Creeaza(ctx, cerere("ana", RolOperator))
	for i, durata := range []time.Duration{time.Minute, 5 * time.Minute, 30 * time.Minute, 120 * time.Minute, 120 * time.Minute} {
		if i == 0 {
			for j := 0; j < pragBlocare-1; j++ {
				m.s.Autentifica(ctx, "ana", "gresit123", aici)
			}
		}
		m.s.Autentifica(ctx, "ana", "gresit123", aici)
		c, _ = m.s.DupaID(ctx, c.ID)
		pana, _ := time.Parse(time.RFC3339Nano, c.BlocatPanaLa)
		if !pana.Equal(m.ceas.Add(durata)) {
			t.Fatalf("pasul %d: blocat până la %s, așteptat +%s", i, c.BlocatPanaLa, durata)
		}
		m.ceas = pana.Add(time.Second)
	}
}

func TestAutoLoginSiIesire(t *testing.T) {
	m := pregateste(t)
	c, _ := m.s.Creeaza(ctx, cerere("ana", RolOperator))
	if _, ok, _ := m.s.AutoLogin(ctx, aici); ok {
		t.Fatal("auto-login fără autentificare anterioară")
	}
	if _, err := m.s.Autentifica(ctx, "ana", "parola123", aici); err != nil {
		t.Fatal(err)
	}
	if x, ok, err := m.s.AutoLogin(ctx, aici); !ok || err != nil || x.ID != c.ID {
		t.Fatalf("auto-login: %+v %v %v", x, ok, err)
	}
	if _, ok, _ := m.s.AutoLogin(ctx, Identitate{SID: "S-alt", Calculator: "PC-ATI"}); ok {
		t.Fatal("auto-login pentru alt utilizator Windows")
	}
	if _, ok, _ := m.s.AutoLogin(ctx, Identitate{SID: aici.SID, Calculator: "ALT-PC"}); ok {
		t.Fatal("auto-login pe alt calculator")
	}
	if m.s.UltimulNume(ctx) != "ana" {
		t.Fatal("ultimul nume")
	}
	if err := m.s.Iesire(ctx, c); err != nil {
		t.Fatal(err)
	}
	if _, ok, _ := m.s.AutoLogin(ctx, aici); ok {
		t.Fatal("auto-login după ieșirea din cont")
	}
}

func TestRecuperareCuRaspunsuri(t *testing.T) {
	m := pregateste(t)
	m.s.Creeaza(ctx, cerere("ana", RolOperator))
	q, err := m.s.IntrebariPentru(ctx, "ana")
	if err != nil || len(q) != 3 || q[1] != "intrebare_4" {
		t.Fatalf("întrebări: %v %v", q, err)
	}
	_, err = m.s.Recupereaza(ctx, "ana", []string{"x", "y", "z"}, "nouaParola1", aici)
	cod(t, err, "cont.raspunsuri_gresite")
	// Răspunsurile se compară fără diacritice, majuscule și spații în plus.
	if _, err := m.s.Recupereaza(ctx, "ana", []string{"  scoala NR. 5 ", "bobita", "BRASOV"}, "nouaParola1", aici); err != nil {
		t.Fatalf("recuperare: %v", err)
	}
	if _, err := m.s.Autentifica(ctx, "ana", "parola123", aici); err == nil {
		t.Fatal("parola veche merge încă")
	}
	if _, err := m.s.Autentifica(ctx, "ana", "nouaParola1", aici); err != nil {
		t.Fatalf("parola nouă: %v", err)
	}
}

func TestAdministrare(t *testing.T) {
	m := pregateste(t)
	admin, _ := m.s.Creeaza(ctx, cerere("ana", RolAdministrator))
	op, _ := m.s.Creeaza(ctx, cerere("dan", RolOperator))
	admin.Rol = RolOperator
	_, err := m.s.Modifica(ctx, admin)
	cod(t, err, "cont.ultimul_administrator")
	admin.Rol, admin.Activ = RolAdministrator, false
	_, err = m.s.Modifica(ctx, admin)
	cod(t, err, "cont.ultimul_administrator")
	op.Rol = RolAdministrator
	if _, err := m.s.Modifica(ctx, op); err != nil {
		t.Fatal(err)
	}
	admin.Activ = false
	if x, err := m.s.Modifica(ctx, admin); err != nil || x.Activ {
		t.Fatalf("dezactivare cu alt administrator: %+v %v", x, err)
	}
	_, err = m.s.Autentifica(ctx, "ana", "parola123", aici)
	cod(t, err, "cont.date_gresite")
	if err := m.s.ReseteazaParola(ctx, op.ID, "resetata99"); err != nil {
		t.Fatal(err)
	}
	if _, err := m.s.Autentifica(ctx, "dan", "resetata99", aici); err != nil {
		t.Fatalf("după resetare: %v", err)
	}
	cod(t, m.s.SchimbaParola(ctx, op, "gresita1", "alta12345"), "cont.parola_veche_gresita")
	if err := m.s.SchimbaParola(ctx, op, "resetata99", "alta12345"); err != nil {
		t.Fatal(err)
	}
	if err := m.s.SchimbaIntrebari(ctx, op, "alta12345", []string{"intrebare_2", "intrebare_3", "intrebare_5"}, []string{"a", "b", "c"}); err != nil {
		t.Fatal(err)
	}
	if _, err := m.s.Recupereaza(ctx, "dan", []string{"A", "B", "C"}, "ultima1234", aici); err != nil {
		t.Fatalf("recuperare cu întrebările noi: %v", err)
	}
}

func TestPazaRefuzaScrierileInainteDeAutentificare(t *testing.T) {
	m := pregateste(t)
	if err := audit.Initializeaza(ctx, m.db); err != nil {
		t.Fatal(err)
	}
	m.db.Close()
	paza := audit.NouaPaza(audit.Rezultat{Stare: audit.StareOK})
	paza.CereAutentificare(true)
	db, err := bd.DeschideCuAudit(m.cale, m.cheie, paza)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	s := Nou(db, m.seif, paza)
	s.acum = m.s.acum
	if _, err := db.Exec(`INSERT INTO setari (cheie, valoare, modificat_la) VALUES ('x', 'y', 'z')`); !audit.EsteEditareBlocata(err) {
		t.Fatalf("scriere înainte de autentificare: %v", err)
	}
	// Crearea primului cont și autentificarea merg (operații de sistem).
	if _, err := s.Creeaza(ctx, cerere("ana", RolAdministrator)); err != nil {
		t.Fatal(err)
	}
	if _, err := s.Autentifica(ctx, "ana", "parola123", aici); err != nil {
		t.Fatal(err)
	}
	paza.SeteazaAutentificat(true)
	if _, err := db.Exec(`INSERT INTO setari (cheie, valoare, modificat_la) VALUES ('x', 'y', 'z')`); err != nil {
		t.Fatalf("scriere după autentificare: %v", err)
	}
	if rez, err := audit.Verifica(ctx, db); err != nil || rez.Stare != audit.StareOK {
		t.Fatalf("sigiliile după operațiile conturilor: %+v %v", rez, err)
	}
}
