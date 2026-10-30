package ai

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestServerNecunoscutNuPorneste(t *testing.T) {
	dir := t.TempDir()
	for _, f := range []string{NumeBinarServer(), "model.gguf"} {
		if err := os.WriteFile(filepath.Join(dir, f), []byte("alt program"), 0o700); err != nil {
			t.Fatal(err)
		}
	}
	m := NouManager(dir, 0)
	if err := m.Porneste(context.Background(), ""); !errors.Is(err, ErrServerNecunoscut) {
		t.Fatalf("server necunoscut: %v", err)
	}
	if m.EsteActiv() {
		t.Fatal("serverul necunoscut a pornit")
	}
}

func TestServerCompilatAcceptatPrinHash(t *testing.T) {
	cale := filepath.Join(t.TempDir(), NumeBinarServer())
	if err := os.WriteFile(cale, []byte("binar compilat"), 0o700); err != nil {
		t.Fatal(err)
	}
	if err := verificaServer(cale); !errors.Is(err, ErrServerNecunoscut) {
		t.Fatalf("fără hash: %v", err)
	}
	vechi := HashServerCompilat
	t.Cleanup(func() { HashServerCompilat = vechi })
	HashServerCompilat = "00" + strings.Repeat("ab", 31)
	if err := verificaServer(cale); !errors.Is(err, ErrServerNecunoscut) {
		t.Fatal("un hash greșit a fost acceptat")
	}
	suma := sha256.Sum256([]byte("binar compilat"))
	HashServerCompilat = hex.EncodeToString(suma[:]) // litere mici: comparația nu ține cont de ele
	if err := verificaServer(cale); err != nil {
		t.Fatalf("binarul compilat: %v", err)
	}
}
