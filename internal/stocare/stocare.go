// Package stocare leagă seiful de baza de date la pornirea aplicației (F05-T12)
// și face backup-ul automat înainte de migrări (F05-T06).
package stocare

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"pontix/internal/audit"
	"pontix/internal/bd"
	"pontix/internal/i18n"
	"pontix/internal/seif"
	"pontix/internal/versiune"
)

const (
	NumeSeif = "seif"
	NumeBaza = "pontix.db"

	// MaxBackupuri: câte copii automate se păstrează (decizie Ioan Tuns, 2026-10-25).
	MaxBackupuri = 7

	prefixBackup       = "pontix-"
	prefixBackupManual = "pontix-manual-"
	sufixBackup        = ".db"
)

// Deschide deschide seiful și baza din folderul de date, verifică sigiliile, aplică
// migrările și redeschide baza cu auditul activ (D044). La prima pornire (fără seif
// și fără bază) le creează pe amândouă. Paza întoarsă spune dacă s-a găsit un incident.
func Deschide(ctx context.Context, radacina string) (*sql.DB, *audit.Paza, error) {
	return deschide(ctx, radacina, bd.Migrari, time.Now)
}

func deschide(ctx context.Context, radacina string, migrari fs.FS, acum func() time.Time) (*sql.DB, *audit.Paza, error) {
	caleSeif := filepath.Join(radacina, NumeSeif)
	caleBaza := filepath.Join(radacina, NumeBaza)

	s, err := seif.Deschide(caleSeif)
	if errors.Is(err, fs.ErrNotExist) {
		// Fără seif, o bază existentă nu mai poate fi citită. Un seif nou ar avea
		// altă cheie și ar ascunde problema, deci se raportează ca seif deteriorat.
		if _, errBaza := os.Stat(caleBaza); errBaza == nil {
			return nil, nil, seif.ErrDeteriorat
		}
		s, err = seif.Creeaza(caleSeif, bd.LungimeCheie)
	}
	if err != nil {
		return nil, nil, err
	}

	_, errBaza := os.Stat(caleBaza)
	bazaNoua := errors.Is(errBaza, fs.ErrNotExist)
	rest, err := aplicaRestaurare(ctx, radacina, s.CheieDate, acum)
	if err != nil {
		return nil, nil, err
	}
	rez, err := verificaSiMigreaza(ctx, caleBaza, radacina, s.CheieDate, migrari, acum)
	if err != nil {
		return nil, nil, err
	}
	if rest != nil && rez.Stare == audit.StareOK {
		if rez, err = noteazaRestaurarea(ctx, caleBaza, s.CheieDate, rest); err != nil {
			return nil, nil, err
		}
	}
	// Copia zilnică nu se face pentru o bază abia creată și nici peste o bază cu incident:
	// rotația ar împinge afară copiile bune.
	if !bazaNoua && rez.Stare != audit.StareIncident {
		if err := backupZilnic(ctx, caleBaza, radacina, s.CheieDate, acum); err != nil {
			return nil, nil, err
		}
	}
	// O copie mai veche, dar validă, a bazei (F06-T09). Restaurarea cerută din aplicație e
	// intenționată și notată în jurnal, deci nu e incident; contorul se rescrie după ea.
	if rest == nil && rez.Stare != audit.StareIncident {
		veche, err := verificaSecventa(ctx, caleBaza, s)
		if err != nil {
			return nil, nil, err
		}
		if veche {
			rez.Stare, rez.CopieVeche = audit.StareIncident, true
		}
	}
	paza := audit.NouaPaza(rez)
	db, err := bd.DeschideCuAudit(caleBaza, s.CheieDate, paza)
	if err != nil {
		return nil, nil, err
	}
	if rez.Stare != audit.StareIncident {
		if err := scrieSecventa(ctx, db, s); err != nil {
			db.Close()
			return nil, nil, err
		}
	}
	return db, paza, nil
}

// noteazaRestaurarea scrie restaurarea în jurnalul bazei restaurate și o resigilează.
func noteazaRestaurarea(ctx context.Context, caleBaza string, cheie []byte, rest *restaurare) (audit.Rezultat, error) {
	db, err := bd.Deschide(caleBaza, cheie)
	if err != nil {
		return audit.Rezultat{}, err
	}
	defer db.Close()
	if err := audit.NoteazaEveniment(ctx, db, "restaurare", nil, rest); err != nil {
		return audit.Rezultat{}, err
	}
	return audit.Verifica(ctx, db)
}

// verificaSiMigreaza lucrează pe o conexiune fără trigger-ele auditului: verifică
// sigiliile, migrează și resigilează. O bază cu incident se migrează, dar nu se
// resigilează, ca migrarea să nu ascundă modificarea (D044 punctul 4).
func verificaSiMigreaza(ctx context.Context, caleBaza, radacina string, cheie []byte, migrari fs.FS, acum func() time.Time) (audit.Rezultat, error) {
	db, err := bd.Deschide(caleBaza, cheie)
	if err != nil {
		return audit.Rezultat{}, err
	}
	defer db.Close()
	rez, err := audit.Verifica(ctx, db)
	if err != nil {
		return rez, err
	}
	deLa, panaLa, err := migreaza(ctx, db, radacina, cheie, migrari, acum)
	if err != nil {
		return rez, err
	}
	if deLa == panaLa || rez.Stare == audit.StareIncident {
		return rez, nil
	}
	if rez.Stare == audit.StareInactiv {
		// Prima sigilare, doar dacă migrarea care aduce sigiliile tocmai s-a aplicat.
		if rez, err = audit.Verifica(ctx, db); err != nil || rez.Stare == audit.StareInactiv {
			return rez, err
		}
		err = audit.Initializeaza(ctx, db)
	} else {
		err = audit.DupaMigrare(ctx, db, deLa, panaLa)
	}
	if err != nil {
		return rez, err
	}
	return audit.Verifica(ctx, db)
}

// migreaza face backup înainte de migrări, dacă baza are deja date (D002 punctul 3),
// apoi aplică migrările și întoarce versiunea dinainte și de după. Dacă backup-ul
// eșuează, migrarea nu se mai face.
func migreaza(ctx context.Context, db *sql.DB, radacina string, cheie []byte, migrari fs.FS, acum func() time.Time) (int64, int64, error) {
	versiuneBaza, inAsteptare, err := bd.StareMigrari(ctx, db, migrari)
	if err != nil {
		return 0, 0, err
	}
	if !inAsteptare {
		return versiuneBaza, versiuneBaza, nil
	}
	if versiuneBaza > 0 {
		dir := filepath.Join(radacina, "backup")
		if err := os.MkdirAll(dir, 0o700); err != nil {
			return 0, 0, err
		}
		nume := prefixBackup + acum().UTC().Format("2006-01-02-150405") + "-v" + versiune.Aplicatie + sufixBackup
		if err := bd.Backup(ctx, db, filepath.Join(dir, nume), cheie); err != nil {
			return 0, 0, &i18n.Eroare{Cod: "stocare.backup_esuat", Param: map[string]string{"detalii": err.Error()}}
		}
		if err := pastreazaUltimele(dir, MaxBackupuri); err != nil {
			return 0, 0, err
		}
	}
	if _, err := bd.Migreaza(ctx, db, migrari); err != nil {
		return 0, 0, err
	}
	versiuneNoua, _, err := bd.StareMigrari(ctx, db, migrari)
	return versiuneBaza, versiuneNoua, err
}

// BackupManual creează la cerere o copie criptată a bazei deschise în folderul
// backup/ din radacina și întoarce calea completă a fișierului creat (F05-T16).
func BackupManual(ctx context.Context, db *sql.DB, radacina string) (string, error) {
	return backupManual(ctx, db, radacina, time.Now)
}

func backupManual(ctx context.Context, db *sql.DB, radacina string, acum func() time.Time) (string, error) {
	s, err := seif.Deschide(filepath.Join(radacina, NumeSeif))
	if err != nil {
		return "", err
	}
	dir := filepath.Join(radacina, "backup")
	if err := os.MkdirAll(dir, 0o700); err != nil {
		return "", err
	}
	baza := prefixBackupManual + acum().UTC().Format("2006-01-02-150405")
	cale := filepath.Join(dir, baza+"-v"+versiune.Aplicatie+sufixBackup)
	for i := 2; ; i++ {
		if _, err := os.Stat(cale); errors.Is(err, fs.ErrNotExist) {
			break
		}
		cale = filepath.Join(dir, fmt.Sprintf("%s-%d-v%s%s", baza, i, versiune.Aplicatie, sufixBackup))
	}
	if err := bd.Backup(ctx, db, cale, s.CheieDate); err != nil {
		return "", &i18n.Eroare{Cod: "stocare.backup_esuat", Param: map[string]string{"detalii": err.Error()}}
	}
	return cale, nil
}

// pastreazaUltimele șterge backup-urile automate mai vechi, astfel încât să rămână
// cel mult n. Numele încep cu data, deci ordinea alfabetică e ordinea în timp.
// Alte fișiere din folder (inclusiv copiile manuale) nu se ating.
func pastreazaUltimele(dir string, n int) error {
	intrari, err := os.ReadDir(dir)
	if err != nil {
		return err
	}
	var backupuri []string
	for _, i := range intrari {
		if !i.IsDir() && strings.HasPrefix(i.Name(), prefixBackup) && !strings.HasPrefix(i.Name(), prefixBackupManual) && strings.HasSuffix(i.Name(), sufixBackup) {
			backupuri = append(backupuri, i.Name())
		}
	}
	sort.Strings(backupuri)
	for len(backupuri) > n {
		if err := os.Remove(filepath.Join(dir, backupuri[0])); err != nil {
			return err
		}
		backupuri = backupuri[1:]
	}
	return nil
}
