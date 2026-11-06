package bd

import (
	"context"
	"database/sql"
	"io/fs"
	"strconv"
	"strings"
	"time"

	"pontix/internal/i18n"
	"pontix/internal/versiune"
)

// Protectia la downgrade si istoricul migrarilor (D009 punctul 3, F06-T06).

// versiuneMaximaMigrari intoarce numarul celei mai noi migrari incluse in .exe.
func versiuneMaximaMigrari(surse fs.FS) (int64, error) {
	intrari, err := fs.ReadDir(surse, ".")
	if err != nil {
		return 0, err
	}
	var maxim int64
	for _, e := range intrari {
		nume := e.Name()
		if !strings.HasSuffix(nume, ".sql") {
			continue
		}
		i := strings.IndexByte(nume, '_')
		if i <= 0 {
			continue
		}
		if n, err := strconv.ParseInt(nume[:i], 10, 64); err == nil && n > maxim {
			maxim = n
		}
	}
	return maxim, nil
}

// verificaCompatibilitate refuza o baza migrata de o versiune mai noua a
// aplicatiei: o versiune veche nu cunoaste schema si ar putea strica datele.
func verificaCompatibilitate(versiuneBaza int64, surse fs.FS) error {
	maxim, err := versiuneMaximaMigrari(surse)
	if err != nil {
		return err
	}
	if maxim > 0 && versiuneBaza > maxim {
		return &i18n.Eroare{Cod: "bd.baza_mai_noua", Param: map[string]string{
			"baza": strconv.FormatInt(versiuneBaza, 10), "cunoscuta": strconv.FormatInt(maxim, 10), "aplicatie": versiune.Aplicatie}}
	}
	return nil
}

// noteazaMigrari scrie in istoric migrarile aplicate acum, cu versiunea
// aplicatiei. Bazele de dinainte de migrarea 0034 nu au tabelul: se sare.
func noteazaMigrari(ctx context.Context, db *sql.DB, versiuni []int64) error {
	if len(versiuni) == 0 {
		return nil
	}
	var exista int
	if err := db.QueryRowContext(ctx, `SELECT count(*) FROM sqlite_schema WHERE type='table' AND name='istoric_migrari'`).Scan(&exista); err != nil || exista == 0 {
		return err
	}
	acum := time.Now().UTC().Format(time.RFC3339)
	for _, v := range versiuni {
		if _, err := db.ExecContext(ctx, `INSERT INTO istoric_migrari(versiune,aplicatie,aplicata_la) VALUES(?,?,?)`, v, versiune.Aplicatie, acum); err != nil {
			return err
		}
	}
	return nil
}
