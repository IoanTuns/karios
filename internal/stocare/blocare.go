package stocare

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"io/fs"
	"os"
	"os/user"
	"path/filepath"
	"time"

	"pontix/internal/audit"
	"pontix/internal/bd"
	"pontix/internal/i18n"
	"pontix/internal/seif"
)

// O singura instanta care scrie (D002 punctul 6, F05-T11): karios/pontix.lock
// ramane blocat cat ruleaza aplicatia si spune cine o are deschisa. A doua
// instanta porneste doar in citire.

// NumeBlocare este fisierul de blocare din folderul de date.
const NumeBlocare = "pontix.lock"

// InfoBlocare spune cine are aplicatia deschisa.
type InfoBlocare struct {
	Calculator string `json:"calculator"`
	Utilizator string `json:"utilizator"`
	DeLa       string `json:"de_la"`
}

// blocareActiva tine fisierul deschis (si blocat) pana la iesirea din proces.
var blocareActiva *os.File

func infoCurenta(acum time.Time) InfoBlocare {
	calc, _ := os.Hostname()
	util := os.Getenv("USERNAME")
	if u, err := user.Current(); err == nil && u.Username != "" {
		util = u.Username
	}
	return InfoBlocare{Calculator: calc, Utilizator: util, DeLa: acum.Format(time.RFC3339)}
}

// blocheaza incearca blocarea fisierului. Daca reuseste, scrie cine il tine si
// intoarce fisierul deschis; altfel citeste si intoarce detinatorul.
func blocheaza(cale string, acum time.Time) (*os.File, *InfoBlocare, error) {
	f, err := os.OpenFile(cale, os.O_RDWR|os.O_CREATE, 0o600)
	if err != nil {
		return nil, nil, err
	}
	ok, err := incearcaBlocarea(f)
	if err != nil {
		f.Close()
		return nil, nil, err
	}
	if !ok {
		defer f.Close()
		var info InfoBlocare
		b, _ := io.ReadAll(io.LimitReader(f, 4096))
		_ = json.Unmarshal(b, &info)
		return nil, &info, nil
	}
	info := infoCurenta(acum)
	b, _ := json.Marshal(info)
	if err := f.Truncate(0); err == nil {
		_, _ = f.WriteAt(b, 0)
		_ = f.Sync()
	}
	return f, nil, nil
}

// DeschideCuBlocare este pornirea aplicatiei: ia blocarea folderului de date si
// deschide totul normal; daca alta instanta o tine, deschide doar in citire.
func DeschideCuBlocare(ctx context.Context, radacina string) (*Deschidere, error) {
	f, detinator, err := blocheaza(filepath.Join(radacina, NumeBlocare), time.Now())
	if err != nil {
		return nil, err
	}
	if detinator == nil {
		blocareActiva = f
		return Deschide(ctx, radacina)
	}
	return deschideDoarCitire(ctx, radacina, detinator)
}

// deschideDoarCitire nu scrie nimic la pornire: fara seif nou, oglinda,
// trecere, restaurare, migrari, copii sau contor. Baza se deschide cu paza in
// modul doar citire; autentificarea ramane posibila (operatie de sistem).
func deschideDoarCitire(ctx context.Context, radacina string, detinator *InfoBlocare) (*Deschidere, error) {
	s, err := seif.Incarca(filepath.Join(radacina, NumeSeif))
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) || errors.Is(err, seif.ErrDeteriorat) {
			return nil, &i18n.Eroare{Cod: "stocare.alta_instanta_pregatire"}
		}
		return nil, err
	}
	if s.AreCopiaAplicatiei() {
		return nil, &i18n.Eroare{Cod: "stocare.alta_instanta_pregatire"}
	}
	if err := s.DeschideWindows(); err != nil {
		return &Deschidere{Seif: s, Blocat: true, DoarCitire: true, Detinator: detinator}, nil
	}
	return continuaDoarCitire(ctx, radacina, s, detinator)
}

func continuaDoarCitire(ctx context.Context, radacina string, s *seif.Seif, detinator *InfoBlocare) (*Deschidere, error) {
	caleBaza := filepath.Join(radacina, NumeBaza)
	if _, err := os.Stat(caleBaza); err != nil {
		return nil, &i18n.Eroare{Cod: "stocare.alta_instanta_pregatire"}
	}
	db, err := bd.Deschide(caleBaza, s.CheieDate)
	if err != nil {
		return nil, err
	}
	_, inAsteptare, err := bd.StareMigrari(ctx, db, bd.Migrari)
	if err == nil && inAsteptare {
		// Baza trebuie actualizata, dar o actualizeaza doar instanta care o tine deschisa.
		err = &i18n.Eroare{Cod: "stocare.alta_instanta_actualizare"}
	}
	var rez audit.Rezultat
	if err == nil {
		rez, err = audit.Verifica(ctx, db)
	}
	db.Close()
	if err != nil {
		return nil, err
	}
	paza := audit.NouaPaza(rez)
	paza.DoarCitire(true)
	lucru, err := bd.DeschideCuAudit(caleBaza, s.CheieDate, paza)
	if err != nil {
		return nil, err
	}
	return &Deschidere{DB: lucru, Paza: paza, Seif: s, DoarCitire: true, Detinator: detinator}, nil
}

// DeblocheazaDoarCitire continua pornirea doar in citire dupa parola sau cheia
// de recuperare; nu adauga copia W in seif (ar fi o scriere).
func DeblocheazaDoarCitire(ctx context.Context, radacina string, s *seif.Seif, detinator *InfoBlocare) (*Deschidere, error) {
	if len(s.CheieDate) == 0 {
		return nil, seif.ErrBlocat
	}
	return continuaDoarCitire(ctx, radacina, s, detinator)
}
