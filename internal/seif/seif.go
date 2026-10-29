// Package seif păstrează cheia de date a bazei criptate (D008 punctul 2).
//
// Cheia de date e aleatoare și se generează o singură dată. Fișierul karios/seif
// păstrează mai multe copii ale ei, fiecare criptată (AES-256-GCM) cu altă cheie:
//
//   - copia "A": cu cheia aplicației, injectată la build (etapa 1);
//   - copii deschise cu un secret (parola sau răspunsurile unui cont, F15),
//     cu cheia derivată prin Argon2id.
//
// Trecerea la etapa 2 (I006) înseamnă doar ștergerea copiei "A".
package seif

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"os"

	"golang.org/x/crypto/argon2"

	"pontix/internal/i18n"
)

// CheieAplicatie este cheia aplicației, în hex (32 de octeți). Deocamdată stă în cod
// (D011): datele nu se pot citi fără executabil, dar oricine are codul sursă le poate
// decripta. Mutarea ei în afara codului: I008. Poate fi înlocuită la build cu
// -ldflags "-X pontix/internal/seif.CheieAplicatie=…".
//
// NU SE SCHIMBĂ: altă cheie face ilizibile toate folderele karios/ existente.
var CheieAplicatie = "c9f3d6df617f8a72e9a7a8f8ee044173f0e6fb083526e7fbb12f1b649679ebbe"

// IdAplicatie e copia deschisă cu cheia aplicației.
const IdAplicatie = "A"

const versiuneFormat = 1

// Argon2id: parametrii se păstrează în fiecare copie, ca să poată fi măriți ulterior.
const (
	argonTimp    = 3
	argonMemorie = 64 * 1024 // KiB
	argonFire    = 4
	lungimeCheie = 32
	lungimeSare  = 16
)

type fisier struct {
	Versiune int              `json:"versiune"`
	Copii    map[string]copie `json:"copii"`
	// Secventa e ultimul număr din jurnalul bazei, autentificat de cine îl scrie (F06-T09):
	// o bază înlocuită cu o copie mai veche are mai puține intrări decât spune seiful.
	Secventa *Secventa `json:"secventa,omitempty"`
}

// Secventa e un contor cu codul lui de autentificare; seiful doar îl păstrează.
type Secventa struct {
	Numar int64  `json:"numar"`
	Mac   []byte `json:"mac"`
}

// Secventa întoarce contorul păstrat sau nil.
func (s *Seif) Secventa() *Secventa { return s.f.Secventa }

// SeteazaSecventa înlocuiește contorul și salvează seiful.
func (s *Seif) SeteazaSecventa(x Secventa) error {
	s.f.Secventa = &x
	return s.salveaza()
}

type copie struct {
	Sare    []byte `json:"sare,omitempty"` // doar pentru copiile deschise cu un secret
	Timp    uint32 `json:"timp,omitempty"`
	Memorie uint32 `json:"memorie,omitempty"`
	Fire    uint8  `json:"fire,omitempty"`
	Nonce   []byte `json:"nonce"`
	Date    []byte `json:"date"`
}

// Seif e fișierul deschis; CheieDate e cheia bazei de date.
type Seif struct {
	cale      string
	f         fisier
	CheieDate []byte
}

var (
	errFaraCheieAplicatie = &i18n.Eroare{Cod: "seif.fara_cheie_aplicatie"}
	// ErrDeteriorat: seiful lipsește din copia cerută, e modificat sau aparține altei chei.
	ErrDeteriorat = &i18n.Eroare{Cod: "seif.deteriorat"}
	// ErrSecretGresit: secretul nu deschide copia (parolă sau răspunsuri greșite).
	ErrSecretGresit = &i18n.Eroare{Cod: "seif.secret_gresit"}
)

func cheieAplicatie() ([]byte, error) {
	c, err := hex.DecodeString(CheieAplicatie)
	if err != nil || len(c) != lungimeCheie {
		return nil, errFaraCheieAplicatie
	}
	return c, nil
}

// Creeaza generează o cheie de date nouă, de lungimea cerută, și scrie seiful cu copia "A".
// Nu suprascrie un seif existent.
func Creeaza(cale string, lungimeCheieDate int) (*Seif, error) {
	ka, err := cheieAplicatie()
	if err != nil {
		return nil, err
	}
	if _, err := os.Stat(cale); err == nil {
		return nil, os.ErrExist
	}
	s := &Seif{cale: cale, f: fisier{Versiune: versiuneFormat, Copii: map[string]copie{}}}
	s.CheieDate = make([]byte, lungimeCheieDate)
	if _, err := rand.Read(s.CheieDate); err != nil {
		return nil, err
	}
	if s.f.Copii[IdAplicatie], err = impacheteaza(ka, IdAplicatie, s.CheieDate); err != nil {
		return nil, err
	}
	return s, s.salveaza()
}

// Deschide citește seiful și scoate cheia de date din copia "A".
func Deschide(cale string) (*Seif, error) {
	ka, err := cheieAplicatie()
	if err != nil {
		return nil, err
	}
	s, err := citeste(cale)
	if err != nil {
		return nil, err
	}
	c, ok := s.f.Copii[IdAplicatie]
	if !ok {
		return nil, ErrDeteriorat
	}
	if s.CheieDate, err = despacheteaza(ka, IdAplicatie, c); err != nil {
		return nil, ErrDeteriorat
	}
	return s, nil
}

// DeschideCuSecret citește seiful și scoate cheia de date din copia id, cu secretul dat.
func DeschideCuSecret(cale, id, secret string) (*Seif, error) {
	s, err := citeste(cale)
	if err != nil {
		return nil, err
	}
	c, ok := s.f.Copii[id]
	if !ok || len(c.Sare) == 0 {
		return nil, ErrSecretGresit
	}
	k := argon2.IDKey([]byte(secret), c.Sare, c.Timp, c.Memorie, c.Fire, lungimeCheie)
	if s.CheieDate, err = despacheteaza(k, id, c); err != nil {
		return nil, ErrSecretGresit
	}
	return s, nil
}

// SeteazaSecret adaugă sau înlocuiește copia id, deschisă cu secretul dat, și salvează seiful.
func (s *Seif) SeteazaSecret(id, secret string) error {
	if id == IdAplicatie {
		return errors.New("seif: copia aplicației nu se deschide cu un secret")
	}
	sare := make([]byte, lungimeSare)
	if _, err := rand.Read(sare); err != nil {
		return err
	}
	k := argon2.IDKey([]byte(secret), sare, argonTimp, argonMemorie, argonFire, lungimeCheie)
	c, err := impacheteaza(k, id, s.CheieDate)
	if err != nil {
		return err
	}
	c.Sare, c.Timp, c.Memorie, c.Fire = sare, argonTimp, argonMemorie, argonFire
	s.f.Copii[id] = c
	return s.salveaza()
}

// StergeCopie elimină copia id și salvează seiful. Ultima copie nu se poate șterge.
func (s *Seif) StergeCopie(id string) error {
	if _, ok := s.f.Copii[id]; !ok {
		return nil
	}
	if len(s.f.Copii) == 1 {
		return errors.New("seif: ultima copie a cheii nu se poate șterge")
	}
	delete(s.f.Copii, id)
	return s.salveaza()
}

// Copii întoarce identificatorii copiilor din seif.
func (s *Seif) Copii() []string {
	var rez []string
	for id := range s.f.Copii {
		rez = append(rez, id)
	}
	return rez
}

// impacheteaza criptează cheia de date; id-ul copiei intră în datele autentificate,
// ca o copie să nu poată fi mutată sub alt id.
func impacheteaza(cheie []byte, id string, cheieDate []byte) (copie, error) {
	gcm, err := nouGCM(cheie)
	if err != nil {
		return copie{}, err
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err := rand.Read(nonce); err != nil {
		return copie{}, err
	}
	return copie{Nonce: nonce, Date: gcm.Seal(nil, nonce, cheieDate, []byte(id))}, nil
}

func despacheteaza(cheie []byte, id string, c copie) ([]byte, error) {
	gcm, err := nouGCM(cheie)
	if err != nil {
		return nil, err
	}
	return gcm.Open(nil, c.Nonce, c.Date, []byte(id))
}

func nouGCM(cheie []byte) (cipher.AEAD, error) {
	b, err := aes.NewCipher(cheie)
	if err != nil {
		return nil, err
	}
	return cipher.NewGCM(b)
}

func citeste(cale string) (*Seif, error) {
	data, err := os.ReadFile(cale)
	if err != nil {
		return nil, err
	}
	s := &Seif{cale: cale}
	if err := json.Unmarshal(data, &s.f); err != nil || s.f.Versiune != versiuneFormat || len(s.f.Copii) == 0 {
		return nil, ErrDeteriorat
	}
	return s, nil
}

// salveaza scrie seiful printr-un fișier temporar și o redenumire, ca o întrerupere
// în timpul scrierii să nu lase un seif pe jumătate scris.
func (s *Seif) salveaza() error {
	data, err := json.Marshal(s.f)
	if err != nil {
		return err
	}
	tmp := s.cale + ".tmp"
	f, err := os.OpenFile(tmp, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0o600)
	if err != nil {
		return err
	}
	if _, err := f.Write(data); err != nil {
		f.Close()
		return err
	}
	if err := f.Sync(); err != nil {
		f.Close()
		return err
	}
	if err := f.Close(); err != nil {
		return err
	}
	return os.Rename(tmp, s.cale)
}
