package stocare

import (
	"context"
	"database/sql"
	"path/filepath"
	"testing"
	"time"

	"pontix/internal/audit"
	"pontix/internal/bd"
	"pontix/internal/seif"
)

// deschideTest pornește ca aplicația, cu migrările reale și ceasul sistemului.
func deschideTest(ctx context.Context, dir string) (*sql.DB, *audit.Paza, error) {
	return deschide(ctx, dir, bd.Migrari, time.Now)
}

// seifTest deschide seiful folderului cu copia W a utilizatorului curent (D054).
func seifTest(t *testing.T, dir string) *seif.Seif {
	t.Helper()
	s, err := seif.Incarca(filepath.Join(dir, NumeSeif))
	if err != nil {
		t.Fatal(err)
	}
	if err := s.DeschideWindows(); err != nil {
		t.Fatal(err)
	}
	return s
}
