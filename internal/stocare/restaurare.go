package stocare

import (
	"context"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"pontix/internal/audit"
	"pontix/internal/bd"
	"pontix/internal/i18n"
	"pontix/internal/versiune"
)

// Restaurarea unui backup se face la pornire (F05-T25): baza e deschisă cât rulează
// aplicația. Aplicația verifică întâi copia și scrie numele ei în fișierul de programare;
// la pornire, baza curentă se salvează ca backup de siguranță și copia îi ia locul.
const (
	numeProgramare      = "restaurare"
	prefixInainteRestau = prefixBackupManual + "inainte-de-restaurare-" // nu intră în rotația celor 7
)

// Backup descrie o copie din folderul backup/.
type Backup struct {
	Nume       string `json:"nume"`
	Dimensiune int64  `json:"dimensiune"`
	Modificat  string `json:"modificat"`
	Manual     bool   `json:"manual"`
}

// ListaBackupuri întoarce copiile din backup/, cele mai noi primele.
func ListaBackupuri(radacina string) ([]Backup, error) {
	intrari, err := os.ReadDir(filepath.Join(radacina, "backup"))
	if errors.Is(err, fs.ErrNotExist) {
		return []Backup{}, nil
	}
	if err != nil {
		return nil, err
	}
	rez := []Backup{}
	for _, i := range intrari {
		if i.IsDir() || !strings.HasPrefix(i.Name(), prefixBackup) || !strings.HasSuffix(i.Name(), sufixBackup) {
			continue
		}
		info, err := i.Info()
		if err != nil {
			continue
		}
		rez = append(rez, Backup{Nume: i.Name(), Dimensiune: info.Size(), Modificat: info.ModTime().UTC().Format(time.RFC3339),
			Manual: strings.HasPrefix(i.Name(), prefixBackupManual)})
	}
	sort.Slice(rez, func(a, b int) bool { return rez[a].Modificat > rez[b].Modificat })
	return rez, nil
}

// caleBackup acceptă doar un nume de fișier din backup/, fără căi.
func caleBackup(radacina, nume string) (string, error) {
	if nume == "" || nume != filepath.Base(nume) || !strings.HasPrefix(nume, prefixBackup) || !strings.HasSuffix(nume, sufixBackup) {
		return "", &i18n.Eroare{Cod: "stocare.backup_inexistent"}
	}
	cale := filepath.Join(radacina, "backup", nume)
	if info, err := os.Stat(cale); err != nil || info.IsDir() {
		return "", &i18n.Eroare{Cod: "stocare.backup_inexistent"}
	}
	return cale, nil
}

func copiaza(sursa, tinta string) error {
	in, err := os.Open(sursa)
	if err != nil {
		return err
	}
	defer in.Close()
	out, err := os.OpenFile(tinta, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0o600)
	if err != nil {
		return err
	}
	if _, err := io.Copy(out, in); err != nil {
		out.Close()
		return err
	}
	if err := out.Sync(); err != nil {
		out.Close()
		return err
	}
	return out.Close()
}

// VerificaBackup deschide o copie temporară a backup-ului cu cheia de date și îi verifică
// sigiliile (D044), fără să atingă fișierul original.
func VerificaBackup(ctx context.Context, radacina, nume string, cheie []byte) (audit.Rezultat, error) {
	cale, err := caleBackup(radacina, nume)
	if err != nil {
		return audit.Rezultat{}, err
	}
	dir, err := os.MkdirTemp(radacina, "verificare-backup-")
	if err != nil {
		return audit.Rezultat{}, err
	}
	defer os.RemoveAll(dir)
	copie := filepath.Join(dir, "copie.db")
	if err := copiaza(cale, copie); err != nil {
		return audit.Rezultat{}, err
	}
	db, err := bd.Deschide(copie, cheie)
	if err != nil {
		return audit.Rezultat{}, err
	}
	defer db.Close()
	return audit.Verifica(ctx, db)
}

// ProgrameazaRestaurare notează backup-ul care va înlocui baza la următoarea pornire.
func ProgrameazaRestaurare(radacina, nume string) error {
	if _, err := caleBackup(radacina, nume); err != nil {
		return err
	}
	return os.WriteFile(filepath.Join(radacina, numeProgramare), []byte(nume), 0o600)
}

// AnuleazaRestaurare șterge programarea, dacă există.
func AnuleazaRestaurare(radacina string) error {
	err := os.Remove(filepath.Join(radacina, numeProgramare))
	if errors.Is(err, fs.ErrNotExist) {
		return nil
	}
	return err
}

// RestaurareProgramata întoarce numele backup-ului programat sau "".
func RestaurareProgramata(radacina string) string {
	b, err := os.ReadFile(filepath.Join(radacina, numeProgramare))
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(b))
}

// restaurare descrie o restaurare aplicată la pornire, pentru jurnal.
type restaurare struct {
	Backup         string `json:"backup"`
	CopieSiguranta string `json:"copie_siguranta"`
}

// aplicaRestaurare rulează la pornire, cu baza închisă: salvează baza curentă ca backup de
// siguranță și pune copia programată în locul ei. Programarea se șterge oricum.
func aplicaRestaurare(ctx context.Context, radacina string, cheie []byte, acum func() time.Time) (*restaurare, error) {
	nume := RestaurareProgramata(radacina)
	if nume == "" {
		return nil, nil
	}
	defer AnuleazaRestaurare(radacina)
	sursa, err := caleBackup(radacina, nume)
	if err != nil {
		return nil, err
	}
	caleBaza := filepath.Join(radacina, NumeBaza)
	siguranta := prefixInainteRestau + acum().UTC().Format("2006-01-02-150405") + "-v" + versiune.Aplicatie + sufixBackup
	if _, err := os.Stat(caleBaza); err == nil {
		db, err := bd.Deschide(caleBaza, cheie)
		if err != nil {
			return nil, err
		}
		err = bd.Backup(ctx, db, filepath.Join(radacina, "backup", siguranta), cheie)
		db.Close()
		if err != nil {
			return nil, &i18n.Eroare{Cod: "stocare.backup_esuat", Param: map[string]string{"detalii": err.Error()}}
		}
	}
	temporar := caleBaza + ".restaurare"
	if err := copiaza(sursa, temporar); err != nil {
		return nil, err
	}
	for _, sufix := range []string{"-wal", "-shm"} {
		if err := os.Remove(caleBaza + sufix); err != nil && !errors.Is(err, fs.ErrNotExist) {
			return nil, fmt.Errorf("restaurare: %w", err)
		}
	}
	if err := os.Rename(temporar, caleBaza); err != nil {
		return nil, err
	}
	return &restaurare{Backup: nume, CopieSiguranta: siguranta}, nil
}
