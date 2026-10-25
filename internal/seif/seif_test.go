package seif

import (
	"bytes"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"testing"

	"pontix/internal/i18n"
)

const lungimeTest = 64 // ca în bd.LungimeCheie

func cuCheieAplicatie(t *testing.T) {
	t.Helper()
	k := make([]byte, lungimeCheie)
	if _, err := rand.Read(k); err != nil {
		t.Fatal(err)
	}
	vechi := CheieAplicatie
	CheieAplicatie = hex.EncodeToString(k)
	t.Cleanup(func() { CheieAplicatie = vechi })
}

func cod(err error) string {
	var e *i18n.Eroare
	if errors.As(err, &e) {
		return e.Cod
	}
	return ""
}

func TestCreeazaSiDeschide(t *testing.T) {
	cuCheieAplicatie(t)
	cale := filepath.Join(t.TempDir(), "seif")

	s, err := Creeaza(cale, lungimeTest)
	if err != nil {
		t.Fatal(err)
	}
	if len(s.CheieDate) != lungimeTest {
		t.Fatalf("cheia de date are %d octeți", len(s.CheieDate))
	}

	d, err := Deschide(cale)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(d.CheieDate, s.CheieDate) {
		t.Error("cheia redeschisă diferă")
	}

	// Cheia de date nu apare în fișier, nici ca octeți, nici în hex.
	data, _ := os.ReadFile(cale)
	if bytes.Contains(data, s.CheieDate) || bytes.Contains(data, []byte(hex.EncodeToString(s.CheieDate))) {
		t.Error("cheia de date apare necriptată în seif")
	}
	if _, err := os.Stat(cale + ".tmp"); !os.IsNotExist(err) {
		t.Error("fișierul temporar a rămas")
	}

	if _, err := Creeaza(cale, lungimeTest); !errors.Is(err, os.ErrExist) {
		t.Errorf("un seif existent nu trebuie suprascris: %v", err)
	}
}

func TestFaraCheieAplicatie(t *testing.T) {
	vechi := CheieAplicatie
	CheieAplicatie = ""
	defer func() { CheieAplicatie = vechi }()

	if _, err := Creeaza(filepath.Join(t.TempDir(), "seif"), lungimeTest); cod(err) != "seif.fara_cheie_aplicatie" {
		t.Errorf("aștepta seif.fara_cheie_aplicatie, a primit %v", err)
	}
}

func TestAltaCheieAplicatie(t *testing.T) {
	cuCheieAplicatie(t)
	cale := filepath.Join(t.TempDir(), "seif")
	if _, err := Creeaza(cale, lungimeTest); err != nil {
		t.Fatal(err)
	}
	cuCheieAplicatie(t) // alt build, altă cheie
	if _, err := Deschide(cale); cod(err) != "seif.deteriorat" {
		t.Errorf("aștepta seif.deteriorat, a primit %v", err)
	}
}

func TestSecret(t *testing.T) {
	cuCheieAplicatie(t)
	cale := filepath.Join(t.TempDir(), "seif")
	s, err := Creeaza(cale, lungimeTest)
	if err != nil {
		t.Fatal(err)
	}
	if err := s.SeteazaSecret("P:cont1", "parola corectă"); err != nil {
		t.Fatal(err)
	}

	d, err := DeschideCuSecret(cale, "P:cont1", "parola corectă")
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(d.CheieDate, s.CheieDate) {
		t.Error("cheia deschisă cu secretul diferă")
	}
	if _, err := DeschideCuSecret(cale, "P:cont1", "parola greșită"); !errors.Is(err, ErrSecretGresit) {
		t.Errorf("parolă greșită: %v", err)
	}
	if _, err := DeschideCuSecret(cale, "P:inexistent", "parola corectă"); !errors.Is(err, ErrSecretGresit) {
		t.Errorf("copie inexistentă: %v", err)
	}
	if _, err := DeschideCuSecret(cale, IdAplicatie, "orice"); !errors.Is(err, ErrSecretGresit) {
		t.Errorf("copia aplicației nu se deschide cu un secret: %v", err)
	}
}

// O copie mutată sub alt identificator nu trebuie să se deschidă.
func TestCopieMutata(t *testing.T) {
	cuCheieAplicatie(t)
	cale := filepath.Join(t.TempDir(), "seif")
	s, err := Creeaza(cale, lungimeTest)
	if err != nil {
		t.Fatal(err)
	}
	if err := s.SeteazaSecret("P:cont1", "secret"); err != nil {
		t.Fatal(err)
	}

	var f fisier
	data, _ := os.ReadFile(cale)
	if err := json.Unmarshal(data, &f); err != nil {
		t.Fatal(err)
	}
	f.Copii["P:cont2"] = f.Copii["P:cont1"]
	data, _ = json.Marshal(f)
	if err := os.WriteFile(cale, data, 0o600); err != nil {
		t.Fatal(err)
	}
	if _, err := DeschideCuSecret(cale, "P:cont2", "secret"); !errors.Is(err, ErrSecretGresit) {
		t.Errorf("copia mutată s-a deschis: %v", err)
	}
}

func TestFisierModificat(t *testing.T) {
	cuCheieAplicatie(t)
	cale := filepath.Join(t.TempDir(), "seif")
	if _, err := Creeaza(cale, lungimeTest); err != nil {
		t.Fatal(err)
	}
	var f fisier
	data, _ := os.ReadFile(cale)
	if err := json.Unmarshal(data, &f); err != nil {
		t.Fatal(err)
	}
	c := f.Copii[IdAplicatie]
	c.Date[0] ^= 1
	f.Copii[IdAplicatie] = c
	data, _ = json.Marshal(f)
	os.WriteFile(cale, data, 0o600)

	if _, err := Deschide(cale); cod(err) != "seif.deteriorat" {
		t.Errorf("seif modificat: %v", err)
	}
	os.WriteFile(cale, []byte("nu e json"), 0o600)
	if _, err := Deschide(cale); cod(err) != "seif.deteriorat" {
		t.Errorf("seif ilizibil: %v", err)
	}
}

// Etapa 2 (I006): fără copia "A", cheia se deschide doar cu un secret.
func TestEtapa2(t *testing.T) {
	cuCheieAplicatie(t)
	cale := filepath.Join(t.TempDir(), "seif")
	s, err := Creeaza(cale, lungimeTest)
	if err != nil {
		t.Fatal(err)
	}
	if err := s.StergeCopie(IdAplicatie); err == nil {
		t.Fatal("ultima copie nu trebuie să se poată șterge")
	}
	if err := s.SeteazaSecret("P:cont1", "secret"); err != nil {
		t.Fatal(err)
	}
	if err := s.StergeCopie(IdAplicatie); err != nil {
		t.Fatal(err)
	}
	if _, err := Deschide(cale); cod(err) != "seif.deteriorat" {
		t.Errorf("fără copia A, deschiderea cu cheia aplicației trebuie refuzată: %v", err)
	}
	d, err := DeschideCuSecret(cale, "P:cont1", "secret")
	if err != nil || !bytes.Equal(d.CheieDate, s.CheieDate) {
		t.Errorf("deschiderea cu secretul după etapa 2: %v", err)
	}
}

func TestMesajeTraduse(t *testing.T) {
	for _, c := range []string{"seif.fara_cheie_aplicatie", "seif.deteriorat", "seif.secret_gresit"} {
		for _, l := range i18n.Limbi() {
			if i18n.T(l, c, nil) == c {
				t.Errorf("%s: %s nu are traducere", l, c)
			}
		}
	}
}
