package ai

import (
	"crypto/sha256"
	"encoding/hex"
	"io"
	"os"
	"strings"

	"pontix/internal/i18n"
)

// serverePermise sunt sumele SHA-256 ale binarelor llama-server livrate cu Pontix (D056, SEC-5).
// Aplicația pornește doar un server din listă: cine poate scrie în karios/ai/ nu poate face
// aplicația să-i ruleze propriul program. Un binar nou (altă versiune llama.cpp) se adaugă aici,
// în același commit cu actualizarea lui.
var serverePermise = map[string]string{
	"87A2E1C6F9196E5F64FD1D26D6C1807291AA256E3BDDF6FD9F25A5FC6C7DF1FD": "llama-server livrat cu 0.87.2",
}

// HashServerCompilat e suma binarului compilat de `build.ps1 -CuCompilareAI`, pusă la build cu
// -ldflags "-X pontix/internal/ai.HashServerCompilat=…" (D056).
var HashServerCompilat string

// ErrServerNecunoscut: binarul din karios/ai/ nu e unul livrat cu aplicația.
var ErrServerNecunoscut = &i18n.Eroare{Cod: "ai.server_necunoscut"}

// verificaServer calculează suma binarului și o caută printre cele permise.
func verificaServer(cale string) error {
	f, err := os.Open(cale)
	if err != nil {
		return err
	}
	defer f.Close()
	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return err
	}
	suma := strings.ToUpper(hex.EncodeToString(h.Sum(nil)))
	if _, ok := serverePermise[suma]; ok || (HashServerCompilat != "" && strings.EqualFold(suma, HashServerCompilat)) {
		return nil
	}
	return ErrServerNecunoscut
}
