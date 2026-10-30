package seif

import (
	"crypto/rand"
	"crypto/subtle"
	"errors"
	"os"
	"sort"
	"strings"

	"golang.org/x/crypto/argon2"

	"pontix/internal/i18n"
)

// Copiile seifului fără cheie în cod (D054):
//
//   - "W:<SID>": cheia de date protejată cu DPAPI pentru utilizatorul Windows SID, pe acest
//     calculator; deschide datele fără parolă;
//   - "P:<cont>" / "R:<cont>": parola, respectiv răspunsurile unui cont (Argon2id);
//   - "K": cheia de recuperare generată de aplicație (Argon2id);
//   - "P0:<cont>" / "R0:<cont>": copiile unui cont dinaintea trecerii D054, legate de cheia de date
//     veche; dovedesc doar parola, până la prima autentificare după trecere.
const (
	PrefixWindows         = "W:"
	PrefixParola          = "P:"
	PrefixRaspunsuri      = "R:"
	PrefixParolaVeche     = "P0:"
	PrefixRaspunsuriVechi = "R0:"
	IdRecuperare          = "K"
)

var (
	// ErrBlocat: nicio copie nu se poate deschide fără un secret (lipsește W pentru acest utilizator).
	ErrBlocat = &i18n.Eroare{Cod: "seif.blocat"}
	// ErrRecuperareGresita: cheia de recuperare nu deschide copia K.
	ErrRecuperareGresita   = &i18n.Eroare{Cod: "seif.recuperare_gresita"}
	errWindowsIndisponibil = errors.New("seif: protecția datelor din Windows nu e disponibilă")
)

// Incarca citește seiful fără să-l deschidă; cheia de date rămâne goală.
func Incarca(cale string) (*Seif, error) {
	return citeste(cale)
}

// CreeazaNou generează o cheie de date nouă și un seif fără copia A (D054). Seiful se salvează
// după ce primește prima copie (W pentru utilizatorul curent, dacă DPAPI e disponibil).
func CreeazaNou(cale string, lungimeCheieDate int) (*Seif, error) {
	if _, err := os.Stat(cale); err == nil {
		return nil, os.ErrExist
	}
	s := &Seif{cale: cale, f: fisier{Versiune: versiuneFormat, Copii: map[string]copie{}}}
	s.CheieDate = make([]byte, lungimeCheieDate)
	if _, err := rand.Read(s.CheieDate); err != nil {
		return nil, err
	}
	return s, nil
}

// Cale întoarce fișierul seifului.
func (s *Seif) Cale() string { return s.cale }

// MutaLa schimbă fișierul în care se salvează seiful (trecerea D054 scrie întâi alături).
func (s *Seif) MutaLa(cale string) { s.cale = cale }

// AreCopia spune dacă seiful are copia id.
func (s *Seif) AreCopia(id string) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	_, ok := s.f.Copii[id]
	return ok
}

// CopiiCuPrefix întoarce, sortate, copiile ale căror id-uri încep cu prefix.
func (s *Seif) CopiiCuPrefix(prefix string) []string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	var rez []string
	for id := range s.f.Copii {
		if strings.HasPrefix(id, prefix) {
			rez = append(rez, id)
		}
	}
	sort.Strings(rez)
	return rez
}

// DeschideWindows deschide cheia de date din copia W a utilizatorului Windows curent.
func (s *Seif) DeschideWindows() error {
	sid, err := sidCurent()
	if err != nil {
		return ErrBlocat
	}
	s.mu.RLock()
	c, ok := s.f.Copii[PrefixWindows+sid]
	s.mu.RUnlock()
	if !ok {
		return ErrBlocat
	}
	cheie, err := deprotejeaza(c.Date)
	if err != nil {
		return ErrBlocat
	}
	s.CheieDate = cheie
	return nil
}

// AdaugaWindows scrie copia W pentru utilizatorul Windows curent și salvează seiful.
func (s *Seif) AdaugaWindows() error {
	sid, err := sidCurent()
	if err != nil {
		return err
	}
	blob, err := protejeaza(s.CheieDate)
	if err != nil {
		return err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.f.Copii[PrefixWindows+sid] = copie{Date: blob}
	return s.salveaza()
}

// deschideCopie încearcă secretul pe copia id și întoarce cheia de date, fără să o seteze.
func (s *Seif) deschideCopie(id, secret string) ([]byte, bool) {
	s.mu.RLock()
	c, ok := s.f.Copii[id]
	s.mu.RUnlock()
	if !ok || len(c.Sare) == 0 {
		return nil, false
	}
	k := argon2.IDKey([]byte(secret), c.Sare, c.Timp, c.Memorie, c.Fire, lungimeCheie)
	aad := id
	if c.IdOriginal != "" {
		aad = c.IdOriginal
	}
	cheie, err := despacheteaza(k, aad, c)
	return cheie, err == nil
}

// DeschideCuParola încearcă parola pe copiile P ale tuturor conturilor (sunt puține) și întoarce
// contul care a deschis datele. Copiile P0 nu deschid datele după trecerea D054.
func (s *Seif) DeschideCuParola(parola string) (string, error) {
	for _, id := range s.CopiiCuPrefix(PrefixParola) {
		if cheie, ok := s.deschideCopie(id, parola); ok {
			s.CheieDate = cheie
			return strings.TrimPrefix(id, PrefixParola), nil
		}
	}
	return "", ErrSecretGresit
}

// VerificaSecret spune dacă secretul deschide copia id și dă exact cheia de date curentă.
// Pentru copiile vechi P0/R0 ajunge ca secretul să le deschidă: ele țin cheia dinaintea trecerii.
func (s *Seif) VerificaSecret(id, secret string) bool {
	cheie, ok := s.deschideCopie(id, secret)
	if !ok {
		return false
	}
	if strings.HasPrefix(id, PrefixParolaVeche) || strings.HasPrefix(id, PrefixRaspunsuriVechi) {
		return true
	}
	return subtle.ConstantTimeCompare(cheie, s.CheieDate) == 1
}

// alfabetRecuperare nu are caractere ușor de confundat (0/O, 1/I/L).
const alfabetRecuperare = "ABCDEFGHJKMNPQRSTUVWXYZ23456789"

// NormalizeazaRecuperare păstrează doar literele și cifrele, cu majuscule.
func NormalizeazaRecuperare(cod string) string {
	var b strings.Builder
	for _, r := range strings.ToUpper(cod) {
		if strings.ContainsRune(alfabetRecuperare, r) {
			b.WriteRune(r)
		}
	}
	return b.String()
}

func codRecuperareNou() (string, error) {
	var b [25]byte
	if _, err := rand.Read(b[:]); err != nil {
		return "", err
	}
	var sb strings.Builder
	for i, x := range b {
		if i > 0 && i%5 == 0 {
			sb.WriteByte('-')
		}
		sb.WriteByte(alfabetRecuperare[int(x)%len(alfabetRecuperare)])
	}
	return sb.String(), nil
}

// SeteazaRecuperare generează o cheie de recuperare nouă, înlocuiește copia K și marchează
// cheia ca neconfirmată (D054 punctul 3). Codul se întoarce o singură dată, pentru afișare.
func (s *Seif) SeteazaRecuperare() (string, error) {
	cod, err := codRecuperareNou()
	if err != nil {
		return "", err
	}
	s.mu.Lock()
	s.f.RecuperareConfirmata = false
	s.mu.Unlock()
	if err := s.SeteazaSecret(IdRecuperare, NormalizeazaRecuperare(cod)); err != nil {
		return "", err
	}
	return cod, nil
}

// DeschideCuRecuperare deschide cheia de date cu cheia de recuperare.
func (s *Seif) DeschideCuRecuperare(cod string) error {
	cheie, ok := s.deschideCopie(IdRecuperare, NormalizeazaRecuperare(cod))
	if !ok {
		return ErrRecuperareGresita
	}
	s.CheieDate = cheie
	return nil
}

// ConfirmaRecuperare verifică tastarea cheii de recuperare și o marchează confirmată.
func (s *Seif) ConfirmaRecuperare(cod string) error {
	if !s.VerificaSecret(IdRecuperare, NormalizeazaRecuperare(cod)) {
		return ErrRecuperareGresita
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.f.RecuperareConfirmata = true
	return s.salveaza()
}

// RecuperareConfirmata spune dacă administratorul a confirmat cheia de recuperare.
func (s *Seif) RecuperareConfirmata() bool {
	s.mu.RLock()
	confirmata := s.f.RecuperareConfirmata
	s.mu.RUnlock()
	return confirmata && s.AreCopia(IdRecuperare)
}

// AreCopiaAplicatiei spune dacă seiful e încă în formatul dinaintea D054 (copia A).
func (s *Seif) AreCopiaAplicatiei() bool { return s.AreCopia(IdAplicatie) }

// RedenumesteCopie mută copia de, cu datele ei, sub id-ul nou (fără salvare). Folosită de trecerea
// D054 pentru P → P0, unde id-ul autentificat intern nu se schimbă: P0 se verifică cu id-ul vechi.
func (s *Seif) RedenumesteCopie(de, la string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if c, ok := s.f.Copii[de]; ok {
		delete(s.f.Copii, de)
		c.IdOriginal = de
		s.f.Copii[la] = c
	}
}

// Salveaza scrie seiful pe disc.
func (s *Seif) Salveaza() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.salveaza()
}

// PreiaCopiiVechi copiază copiile conturilor din seiful dinaintea trecerii D054 ca P0/R0:
// ele dovedesc parola și răspunsurile până la prima autentificare, dar țin cheia de date veche.
func (s *Seif) PreiaCopiiVechi(din *Seif) {
	din.mu.RLock()
	defer din.mu.RUnlock()
	s.mu.Lock()
	defer s.mu.Unlock()
	for id, c := range din.f.Copii {
		nou := ""
		switch {
		case strings.HasPrefix(id, PrefixParola):
			nou = PrefixParolaVeche + strings.TrimPrefix(id, PrefixParola)
		case strings.HasPrefix(id, PrefixRaspunsuri):
			nou = PrefixRaspunsuriVechi + strings.TrimPrefix(id, PrefixRaspunsuri)
		default:
			continue
		}
		c.IdOriginal = id
		s.f.Copii[nou] = c
	}
}
