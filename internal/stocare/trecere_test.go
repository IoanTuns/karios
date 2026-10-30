package stocare

import (
	"bytes"
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"pontix/internal/audit"
	"pontix/internal/bd"
	"pontix/internal/seif"
)

// folderVechi construiește un folder de date în formatul dinaintea D054: seif cu copia A,
// bază sigilată cu o setare, un cont cu parolă, o lună închisă cu amprentă și un backup.
func folderVechi(t *testing.T) (dir string, cheieVeche []byte, amprenta string) {
	t.Helper()
	ctx := context.Background()
	dir = t.TempDir()
	s, err := seif.Creeaza(filepath.Join(dir, NumeSeif), bd.LungimeCheie)
	if err != nil {
		t.Fatal(err)
	}
	if err := s.SeteazaSecret(seif.PrefixParola+"cont-ana", "parola123"); err != nil {
		t.Fatal(err)
	}
	db, err := bd.Deschide(filepath.Join(dir, NumeBaza), s.CheieDate)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := bd.Migreaza(ctx, db, bd.Migrari); err != nil {
		t.Fatal(err)
	}
	for _, q := range []string{
		`INSERT INTO setari (cheie, valoare, modificat_la) VALUES ('proba', 'da', 'x')`,
		`INSERT INTO angajati (id, nume, marca, creat_la, modificat_la) VALUES ('a', 'Pop', 'M1', 'x', 'x')`,
		`INSERT INTO pontaj_ciorne (angajat_id, data, valoare, modificat_la) VALUES ('a', '2026-09-03', 'Z7', 'x')`,
	} {
		if _, err := db.Exec(q); err != nil {
			t.Fatal(err)
		}
	}
	if amprenta, err = audit.AmprentaLuna(ctx, db, 2026, 9); err != nil {
		t.Fatal(err)
	}
	if _, err := db.Exec(`INSERT INTO stare_luni (an, luna, stare, modificat_la, amprenta) VALUES (2026, 9, 'INCHISA', 'x', ?)`, amprenta); err != nil {
		t.Fatal(err)
	}
	if err := audit.Initializeaza(ctx, db); err != nil {
		t.Fatal(err)
	}
	if _, err := BackupManual(ctx, db, dir, s.CheieDate); err != nil {
		t.Fatal(err)
	}
	db.Close()
	return dir, s.CheieDate, amprenta
}

func TestTrecereaLaD054(t *testing.T) {
	ctx := context.Background()
	dir, cheieVeche, amprentaVeche := folderVechi(t)
	d, err := deschideTot(ctx, dir, bd.Migrari, ceas(0))
	if err != nil {
		t.Fatal(err)
	}
	defer d.DB.Close()
	if d.Blocat || d.CodRecuperare == "" {
		t.Fatalf("după trecere: blocat=%v cod=%q", d.Blocat, d.CodRecuperare)
	}
	if r := d.Paza.Rezultat(); r.Stare != audit.StareOK {
		t.Fatalf("sigiliile după trecere: %+v", r)
	}
	s := d.Seif
	if s.AreCopiaAplicatiei() || bytes.Equal(s.CheieDate, cheieVeche) || len(s.CopiiCuPrefix(seif.PrefixWindows)) != 1 || !s.AreCopia(seif.IdRecuperare) {
		t.Fatalf("seiful după trecere: %v", s.Copii())
	}
	if !s.AreCopia(seif.PrefixParolaVeche+"cont-ana") || s.AreCopia(seif.PrefixParola+"cont-ana") {
		t.Fatalf("copiile contului: %v", s.Copii())
	}
	var v string
	if err := d.DB.QueryRow(`SELECT valoare FROM setari WHERE cheie = 'proba'`).Scan(&v); err != nil || v != "da" {
		t.Fatalf("datele după trecere: %q %v", v, err)
	}
	// Cheia veche (publică) nu mai deschide baza și nici backup-urile.
	if db, err := bd.Deschide(filepath.Join(dir, NumeBaza), cheieVeche); err == nil {
		db.Close()
		t.Fatal("cheia veche deschide încă baza")
	}
	lista, _ := ListaBackupuri(dir)
	for _, b := range lista {
		if db, err := bd.Deschide(filepath.Join(dir, "backup", b.Nume), cheieVeche); err == nil {
			db.Close()
			t.Fatalf("cheia veche deschide încă backup-ul %s", b.Nume)
		}
		if db, err := bd.Deschide(filepath.Join(dir, "backup", b.Nume), s.CheieDate); err != nil {
			t.Fatalf("backup-ul %s nu se deschide cu cheia nouă: %v", b.Nume, err)
		} else {
			db.Close()
		}
	}
	// Amprenta lunii închise e recalculată, iar perechea veche → nouă e în jurnal.
	var amprentaNoua, jurnal string
	if err := d.DB.QueryRow(`SELECT amprenta FROM stare_luni WHERE an = 2026 AND luna = 9`).Scan(&amprentaNoua); err != nil || amprentaNoua == amprentaVeche {
		t.Fatalf("amprenta: %q (veche %q) %v", amprentaNoua, amprentaVeche, err)
	}
	if err := d.DB.QueryRow(`SELECT valoare_noua FROM jurnal_modificari WHERE tabel = 'audit' AND rand_id = 'trecere_d054'`).Scan(&jurnal); err != nil ||
		!strings.Contains(jurnal, amprentaVeche) || !strings.Contains(jurnal, amprentaNoua) {
		t.Fatalf("jurnalul trecerii: %q %v", jurnal, err)
	}
	for _, f := range []string{NumeSeif + sufixVechi, NumeBaza + sufixVechi, NumeSeif + sufixNou, NumeBaza + sufixNou} {
		if existaFisier(filepath.Join(dir, f)) {
			t.Fatalf("a rămas %s", f)
		}
	}
	d.DB.Close()
	// Pornirea următoare: fără trecere, cu cheia de recuperare încă neconfirmată (deci una nouă).
	d2, err := deschideTot(ctx, dir, bd.Migrari, ceas(1))
	if err != nil {
		t.Fatal(err)
	}
	defer d2.DB.Close()
	if d2.Blocat || d2.Paza.Rezultat().Stare != audit.StareOK || d2.CodRecuperare == "" || d2.CodRecuperare == d.CodRecuperare {
		t.Fatalf("a doua pornire: %+v cod=%q", d2.Paza.Rezultat(), d2.CodRecuperare)
	}
	if err := d2.Seif.ConfirmaRecuperare(d2.CodRecuperare); err != nil {
		t.Fatal(err)
	}
	d2.DB.Close()
	d3, err := deschideTot(ctx, dir, bd.Migrari, ceas(2))
	if err != nil {
		t.Fatal(err)
	}
	defer d3.DB.Close()
	if d3.CodRecuperare != "" {
		t.Fatal("cheia confirmată a fost regenerată")
	}
}

func TestTrecereIntreruptaSeTermina(t *testing.T) {
	ctx := context.Background()
	dir, cheieVeche, _ := folderVechi(t)
	// Întrerupere după înlocuirea seifului, înainte de backup-uri și curățenie.
	caleSeif, caleBaza := filepath.Join(dir, NumeSeif), filepath.Join(dir, NumeBaza)
	vechi, _ := seif.Deschide(caleSeif)
	nou, err := seif.CreeazaNou(caleSeif+sufixNou, bd.LungimeCheie)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := nou.SeteazaRecuperare(); err != nil {
		t.Fatal(err)
	}
	if err := nou.AdaugaWindows(); err != nil {
		t.Fatal(err)
	}
	db0, _ := bd.Deschide(caleBaza, vechi.CheieDate)
	if err := bd.Backup(ctx, db0, caleBaza+sufixNou, nou.CheieDate); err != nil {
		t.Fatal(err)
	}
	db0.Close()
	if err := os.Rename(caleSeif, caleSeif+sufixVechi); err != nil {
		t.Fatal(err)
	}
	if err := mutaBaza(caleBaza, caleBaza+sufixVechi); err != nil {
		t.Fatal(err)
	}
	// Pornirea termină înlocuirea, recriptează backup-urile și șterge fișierele vechi.
	d, err := deschideTot(ctx, dir, bd.Migrari, ceas(0))
	if err != nil {
		t.Fatal(err)
	}
	defer d.DB.Close()
	if d.Blocat || !bytes.Equal(d.Seif.CheieDate, nou.CheieDate) || existaFisier(caleSeif+sufixVechi) {
		t.Fatalf("trecerea întreruptă: blocat=%v", d.Blocat)
	}
	lista, _ := ListaBackupuri(dir)
	for _, b := range lista {
		if db, err := bd.Deschide(filepath.Join(dir, "backup", b.Nume), cheieVeche); err == nil {
			db.Close()
			t.Fatalf("backup rămas pe cheia veche: %s", b.Nume)
		}
	}
}

func TestDeblocareCandLipsesteCopiaW(t *testing.T) {
	ctx := context.Background()
	dir := t.TempDir()
	db, _, err := deschideTest(ctx, dir)
	if err != nil {
		t.Fatal(err)
	}
	db.Close()
	s := seifTest(t, dir)
	cheie := s.CheieDate
	if err := s.SeteazaSecret(seif.PrefixParola+"cont-ana", "parola123"); err != nil {
		t.Fatal(err)
	}
	// Alt utilizator Windows: copia W nu e a lui (o ștergem pe a noastră).
	for _, id := range s.CopiiCuPrefix(seif.PrefixWindows) {
		if err := s.StergeCopie(id); err != nil {
			t.Fatal(err)
		}
	}
	d, err := Deschide(ctx, dir)
	if err != nil || !d.Blocat || d.DB != nil {
		t.Fatalf("fără W: %+v %v", d, err)
	}
	if _, err := d.Seif.DeschideCuParola("gresita12"); err == nil {
		t.Fatal("parola greșită a deblocat datele")
	}
	cont, err := d.Seif.DeschideCuParola("parola123")
	if err != nil || cont != "cont-ana" || !bytes.Equal(d.Seif.CheieDate, cheie) {
		t.Fatalf("deblocare: %q %v", cont, err)
	}
	d2, err := Deblocheaza(ctx, dir, d.Seif)
	if err != nil || d2.DB == nil {
		t.Fatalf("după deblocare: %v", err)
	}
	d2.DB.Close()
	// Copia W s-a adăugat: pornirea următoare nu mai e blocată.
	d3, err := Deschide(ctx, dir)
	if err != nil || d3.Blocat {
		t.Fatalf("după deblocare, pornirea: %+v %v", d3, err)
	}
	d3.DB.Close()
}
