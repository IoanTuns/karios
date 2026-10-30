package seif

import (
	"bytes"
	"path/filepath"
	"runtime"
	"testing"
)

func seifNou(t *testing.T) *Seif {
	t.Helper()
	s, err := CreeazaNou(filepath.Join(t.TempDir(), "seif"), lungimeTest)
	if err != nil {
		t.Fatal(err)
	}
	return s
}

func TestCopiaWindows(t *testing.T) {
	if runtime.GOOS != "windows" {
		t.Skip("DPAPI există doar pe Windows")
	}
	s := seifNou(t)
	if err := s.AdaugaWindows(); err != nil {
		t.Fatal(err)
	}
	if s.AreCopiaAplicatiei() || len(s.CopiiCuPrefix(PrefixWindows)) != 1 {
		t.Fatalf("copii: %v", s.Copii())
	}
	r, err := Incarca(s.Cale())
	if err != nil {
		t.Fatal(err)
	}
	if len(r.CheieDate) != 0 {
		t.Fatal("Incarca nu trebuie să deschidă cheia")
	}
	if err := r.DeschideWindows(); err != nil || !bytes.Equal(r.CheieDate, s.CheieDate) {
		t.Fatalf("deschidere W: %v", err)
	}
}

func TestFaraCopieWEsteBlocat(t *testing.T) {
	s := seifNou(t)
	if err := s.SeteazaSecret(PrefixParola+"ana", "parola123"); err != nil {
		t.Fatal(err)
	}
	r, _ := Incarca(s.Cale())
	if err := r.DeschideWindows(); cod(err) != "seif.blocat" {
		t.Fatalf("fără W: %v", err)
	}
}

func TestDeblocareCuParolaSiRecuperare(t *testing.T) {
	s := seifNou(t)
	for id, p := range map[string]string{PrefixParola + "ana": "parola123", PrefixParola + "dan": "altaparola9"} {
		if err := s.SeteazaSecret(id, p); err != nil {
			t.Fatal(err)
		}
	}
	cod1, err := s.SeteazaRecuperare()
	if err != nil || len(cod1) != 29 || s.RecuperareConfirmata() {
		t.Fatalf("recuperare: %q %v confirmată=%v", cod1, err, s.RecuperareConfirmata())
	}
	r, _ := Incarca(s.Cale())
	if cont, err := r.DeschideCuParola("altaparola9"); err != nil || cont != "dan" || !bytes.Equal(r.CheieDate, s.CheieDate) {
		t.Fatalf("parolă: %q %v", cont, err)
	}
	r, _ = Incarca(s.Cale())
	if _, err := r.DeschideCuParola("gresit123"); cod(err) != "seif.secret_gresit" {
		t.Fatalf("parolă greșită: %v", err)
	}
	r, _ = Incarca(s.Cale())
	// Cheia de recuperare se acceptă cu litere mici și fără cratime.
	if err := r.DeschideCuRecuperare(NormalizeazaRecuperare(cod1)[:10] + " " + NormalizeazaRecuperare(cod1)[10:]); err != nil || !bytes.Equal(r.CheieDate, s.CheieDate) {
		t.Fatalf("recuperare: %v", err)
	}
	if err := r.ConfirmaRecuperare("AAAAA-BBBBB-CCCCC-DDDDD-EEEEE"); cod(err) != "seif.recuperare_gresita" {
		t.Fatalf("confirmare greșită: %v", err)
	}
	if err := r.ConfirmaRecuperare(cod1); err != nil || !r.RecuperareConfirmata() {
		t.Fatalf("confirmare: %v", err)
	}
	if r2, _ := Incarca(s.Cale()); !r2.RecuperareConfirmata() {
		t.Fatal("confirmarea nu s-a salvat")
	}
	// O cheie nouă anulează confirmarea și pe cea veche.
	cod2, _ := r.SeteazaRecuperare()
	if r.RecuperareConfirmata() || r.DeschideCuRecuperare(cod1) == nil || r.DeschideCuRecuperare(cod2) != nil {
		t.Fatal("cheia de recuperare nu a fost înlocuită")
	}
}

func TestCopiaRedenumitaSeVerificaCuIdulOriginal(t *testing.T) {
	s := seifNou(t)
	if err := s.SeteazaSecret(PrefixParola+"ana", "parola123"); err != nil {
		t.Fatal(err)
	}
	s.RedenumesteCopie(PrefixParola+"ana", PrefixParolaVeche+"ana")
	// Cheia de date s-a schimbat (trecerea D054); copia veche dovedește doar parola.
	s.CheieDate = bytes.Repeat([]byte{7}, lungimeTest)
	if err := s.Salveaza(); err != nil {
		t.Fatal(err)
	}
	r, _ := Incarca(s.Cale())
	r.CheieDate = s.CheieDate
	if !r.VerificaSecret(PrefixParolaVeche+"ana", "parola123") || r.VerificaSecret(PrefixParolaVeche+"ana", "gresit123") {
		t.Fatal("copia P0 nu verifică parola")
	}
	if _, err := r.DeschideCuParola("parola123"); err == nil {
		t.Fatal("P0 nu trebuie să deschidă datele")
	}
}
