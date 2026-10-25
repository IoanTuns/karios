package bd

import (
	"context"
	"database/sql"
	"io/fs"

	"github.com/pressly/goose/v3"
)

// Proba 0.20.0 avea migrările 0001–0011 și 0013. Permitem completarea
// numai a acelui istoric cunoscut; celelalte goluri păstrează refuzul goose.
func furnizorMigrari(ctx context.Context, db *sql.DB, surse fs.FS) (*goose.Provider, error) {
	var optiuni []goose.ProviderOption
	_, cereri := fs.Stat(surse, "0012_cereri_concediu.sql")
	_, rapoarte := fs.Stat(surse, "0013_definitii_rapoarte.sql")
	if cereri == nil && rapoarte == nil {
		var tabele int
		if err := db.QueryRowContext(ctx, `SELECT count(*) FROM sqlite_schema WHERE type='table' AND name IN ('goose_db_version','definitii_rapoarte')`).Scan(&tabele); err != nil {
			return nil, err
		}
		if tabele == 2 {
			var numar, maxim, areCereri int
			if err := db.QueryRowContext(ctx, `SELECT count(DISTINCT version_id), coalesce(max(version_id),0), coalesce(sum(CASE WHEN version_id=12 THEN 1 ELSE 0 END),0) FROM goose_db_version WHERE is_applied=1 AND version_id>0`).Scan(&numar, &maxim, &areCereri); err != nil {
				return nil, err
			}
			if numar == 12 && maxim == 13 && areCereri == 0 {
				optiuni = append(optiuni, goose.WithAllowOutofOrder(true))
			}
		}
	}
	return goose.NewProvider(goose.DialectSQLite3, db, surse, optiuni...)
}
