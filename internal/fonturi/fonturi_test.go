package fonturi

import (
	"bytes"
	"encoding/binary"
	"errors"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"golang.org/x/image/font/gofont/gobold"
	"golang.org/x/image/font/gofont/goregular"
	"golang.org/x/image/font/sfnt"

	"pontix/internal/i18n"
)

// colectie construieste un .ttc din fonturi independente, cu offseturi absolute.
func colectie(t *testing.T, fonturi ...[]byte) []byte {
	t.Helper()
	out := make([]byte, 12+4*len(fonturi))
	copy(out, "ttcf")
	binary.BigEndian.PutUint32(out[4:], 0x00010000)
	binary.BigEndian.PutUint32(out[8:], uint32(len(fonturi)))
	for k, f := range fonturi {
		n := int(binary.BigEndian.Uint16(f[4:]))
		start := len(out)
		binary.BigEndian.PutUint32(out[12+4*k:], uint32(start))
		out = append(out, f[:12+16*n]...)
		for i := 0; i < n; i++ {
			rec := f[12+16*i:]
			off, lung := binary.BigEndian.Uint32(rec[8:]), binary.BigEndian.Uint32(rec[12:])
			binary.BigEndian.PutUint32(out[start+12+16*i+8:], uint32(len(out)))
			out = append(out, f[off:off+lung]...)
			for len(out)%4 != 0 {
				out = append(out, 0)
			}
		}
	}
	return out
}

// cuFsType copiaza fontul si schimba fsType din tabelul OS/2.
func cuFsType(t *testing.T, f []byte, fs uint16) []byte {
	t.Helper()
	c := append([]byte(nil), f...)
	_, tab, _, err := tabele(bytes.NewReader(c), 0)
	if err != nil {
		t.Fatal(err)
	}
	binary.BigEndian.PutUint16(c[tab["OS/2"].offset+8:], fs)
	return c
}

func scrie(t *testing.T, dir, nume string, b []byte) {
	t.Helper()
	if err := os.WriteFile(filepath.Join(dir, nume), b, 0600); err != nil {
		t.Fatal(err)
	}
}

func numeFonturi(l []Font) map[string]Font {
	m := map[string]Font{}
	for _, f := range l {
		m[f.Nume] = f
	}
	return m
}

func TestListaSiLicente(t *testing.T) {
	sistem, utilizator := t.TempDir(), t.TempDir()
	scrie(t, sistem, "GOREG.TTF", goregular.TTF)
	scrie(t, sistem, "restrictionat.ttf", cuFsType(t, gobold.TTF, fsRestrictionat))
	scrie(t, sistem, "vechi.fon", []byte("MZ bitmap"))
	scrie(t, sistem, "stricat.ttf", []byte("nu este font"))
	// Acelasi font si in folderul utilizatorului: apare o singura data.
	scrie(t, utilizator, "copie.ttf", goregular.TTF)
	l := numeFonturi(Nou(sistem, utilizator).Lista())
	if len(l) != 1 {
		t.Fatalf("astept doar Go Regular, primit %v", l)
	}
	f, ok := l["Go Regular"]
	if !ok || f.Familie != "Go" || f.Stil != "Regular" || !f.Subset {
		t.Fatalf("antet gresit: %+v", f)
	}
}

func TestPermis(t *testing.T) {
	cazuri := map[uint16]bool{0: true, 0x0004: true, 0x0008: true, fsRestrictionat: false, fsRestrictionat | 0x0004: true, fsDoarBitmap: false, fsFaraSubset: true}
	for fs, astept := range cazuri {
		if permis(fs) != astept {
			t.Errorf("fsType %#04x: astept %v", fs, astept)
		}
	}
}

func TestFaraSubset(t *testing.T) {
	d := t.TempDir()
	scrie(t, d, "intreg.ttf", cuFsType(t, goregular.TTF, fsFaraSubset))
	l := Nou(d).Lista()
	if len(l) != 1 || l[0].Subset {
		t.Fatalf("fontul fara subset trebuie marcat: %+v", l)
	}
}

func TestDateDinColectie(t *testing.T) {
	d := t.TempDir()
	scrie(t, d, "familie.ttc", colectie(t, goregular.TTF, gobold.TTF))
	s := Nou(d)
	if l := numeFonturi(s.Lista()); len(l) != 2 {
		t.Fatalf("astept doua fonturi din colectie, primit %v", l)
	}
	b, f, err := s.Date("Go Bold")
	if err != nil {
		t.Fatal(err)
	}
	if f.Stil != "Bold" {
		t.Fatalf("font gresit: %+v", f)
	}
	// Fontul extras trebuie sa fie un fisier independent, citibil ca sfnt.
	font, err := sfnt.Parse(b)
	if err != nil {
		t.Fatalf("font extras invalid: %v", err)
	}
	nume, err := font.Name(nil, sfnt.NameIDFull)
	if err != nil || nume != "Go Bold" {
		t.Fatalf("nume din fontul extras: %q %v", nume, err)
	}
}

func TestDateFisierSimpluSiInexistent(t *testing.T) {
	d := t.TempDir()
	scrie(t, d, "go.ttf", goregular.TTF)
	s := Nou(d)
	b, _, err := s.Date("Go Regular")
	if err != nil || !bytes.Equal(b, goregular.TTF) {
		t.Fatalf("fisierul simplu trebuie intors neschimbat: %v", err)
	}
	// Numele se cauta in lista; o cale nu este acceptata.
	for _, n := range []string{"Arial Inexistent", filepath.Join(d, "go.ttf")} {
		_, _, err := s.Date(n)
		var e *i18n.Eroare
		if !errors.As(err, &e) || e.Cod != "fonturi.inexistent" {
			t.Fatalf("%q: astept fonturi.inexistent, primit %v", n, err)
		}
	}
}

func TestFonturiSistemWindows(t *testing.T) {
	if runtime.GOOS != "windows" {
		t.Skip("doar pe Windows")
	}
	l := Nou().Lista()
	if len(l) == 0 {
		t.Fatal("niciun font Windows gasit")
	}
	m := numeFonturi(l)
	f, ok := m["Arial"]
	if !ok {
		t.Skip("Arial nu este instalat")
	}
	b, _, err := Nou().Date(f.Nume)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := sfnt.Parse(b); err != nil {
		t.Fatalf("Arial nu se poate citi: %v", err)
	}
}

// I125: numele scris cu alte majuscule gaseste fontul si intoarce numele corect.
func TestDateIndiferentDeMajuscule(t *testing.T) {
	d := t.TempDir()
	scrie(t, d, "familie.ttc", colectie(t, goregular.TTF, gobold.TTF))
	s := Nou(d)
	for _, nume := range []string{"go bold", "  GO BOLD ", "Go Bold"} {
		b, f, err := s.Date(nume)
		if err != nil || len(b) == 0 || f.Nume != "Go Bold" {
			t.Fatalf("%q: %v %q", nume, err, f.Nume)
		}
	}
	if _, _, err := s.Date("go bolder"); err == nil {
		t.Fatal("un nume inexistent nu trebuie gasit")
	}
}
