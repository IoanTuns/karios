package stocare

import (
	"context"
	"strings"
	"testing"
	"time"

	"pontix/internal/bd"
)

func TestBackupZilnicLaPrimaPornireDinZi(t *testing.T) {
	ctx := context.Background()
	dir := t.TempDir()
	la := func(zi, ora int) func() time.Time {
		return func() time.Time { return time.Date(2026, 10, zi, ora, 0, 0, 0, time.UTC) }
	}
	porneste := func(c func() time.Time) {
		t.Helper()
		db, _, err := deschide(ctx, dir, bd.Migrari, c)
		if err != nil {
			t.Fatal(err)
		}
		db.Close()
	}
	zilnice := func() int {
		n := 0
		for _, b := range backupuri(t, dir) {
			if strings.Contains(b, "-zilnic-") {
				n++
			}
		}
		return n
	}
	porneste(la(29, 8))
	if n := zilnice(); n != 0 {
		t.Fatalf("copie zilnică pentru o bază abia creată: %d", n)
	}
	porneste(la(29, 9))
	porneste(la(29, 15))
	if n := zilnice(); n != 1 {
		t.Fatalf("o singură copie în aceeași zi: %d", n)
	}
	porneste(la(30, 7))
	if n := zilnice(); n != 2 {
		t.Fatalf("copie nouă în ziua următoare: %d", n)
	}
	modificaExtern(t, dir, `INSERT INTO setari (cheie, valoare, modificat_la) VALUES ('strecurat', 'da', 'x')`)
	porneste(la(31, 7))
	if n := zilnice(); n != 2 {
		t.Fatalf("copie zilnică peste o bază cu incident: %d", n)
	}
}
