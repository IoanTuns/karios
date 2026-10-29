package stocare

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"time"

	"pontix/internal/bd"
	"pontix/internal/i18n"
	"pontix/internal/versiune"
)

// backupZilnic face o copie automată la prima pornire din zi (F05-T26), dacă în acea zi
// (UTC) nu există deja una automată. Intră în rotația celor MaxBackupuri copii automate.
func backupZilnic(ctx context.Context, caleBaza, radacina string, cheie []byte, acum func() time.Time) error {
	dir := filepath.Join(radacina, "backup")
	if err := os.MkdirAll(dir, 0o700); err != nil {
		return err
	}
	zi := acum().UTC().Format("2006-01-02")
	intrari, err := os.ReadDir(dir)
	if err != nil {
		return err
	}
	for _, i := range intrari {
		if strings.HasPrefix(i.Name(), prefixBackup+zi) {
			return nil
		}
	}
	db, err := bd.Deschide(caleBaza, cheie)
	if err != nil {
		return err
	}
	defer db.Close()
	nume := prefixBackup + acum().UTC().Format("2006-01-02-150405") + "-zilnic-v" + versiune.Aplicatie + sufixBackup
	if err := bd.Backup(ctx, db, filepath.Join(dir, nume), cheie); err != nil {
		return &i18n.Eroare{Cod: "stocare.backup_esuat", Param: map[string]string{"detalii": err.Error()}}
	}
	return pastreazaUltimele(dir, MaxBackupuri)
}
