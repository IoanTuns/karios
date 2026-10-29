// Package ai gestionează integrarea și comunicarea exclusiv locală cu motorul LLM (D042, F14).
package ai

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

// PortImplicit este portul local standard pe care rulează serverul de inferență.
const PortImplicit = 28491

// StareMediu descrie starea componentelor motorului local pe disc și în memorie.
type StareMediu struct {
	ModelActiv        string   `json:"model_activ"`
	ContextConfigurat int      `json:"context_configurat"`
	Instalat          bool     `json:"instalat"`
	CaleFolder        string   `json:"cale_folder"`
	CaleBinar         string   `json:"cale_binar"`
	CaleModel         string   `json:"cale_model"`
	ModeleDisponibile []string `json:"modele_disponibile"`
	Activ             bool     `json:"activ"`
	Port              int      `json:"port"`
	Mesaj             string   `json:"mesaj"`
}

// NumeBinarServer întoarce numele executabilului serverului local în funcție de sistemul de operare.
func NumeBinarServer() string {
	if runtime.GOOS == "windows" {
		return "llama-server.exe"
	}
	return "llama-server"
}

// DetecteazaMediu scanează folderul dedicat și verifică dacă motorul și modelele sunt prezente.
func DetecteazaMediu(folderAI string, port int) StareMediu {
	if port <= 0 {
		port = PortImplicit
	}

	rez := StareMediu{
		CaleFolder: folderAI,
		Port:       port,
	}

	binar := filepath.Join(folderAI, NumeBinarServer())
	if info, err := os.Stat(binar); err == nil && !info.IsDir() {
		rez.CaleBinar = binar
	}

	// Caută fișiere .gguf în folderul AI
	intrari, err := os.ReadDir(folderAI)
	if err == nil {
		for _, e := range intrari {
			if !e.IsDir() && strings.HasSuffix(strings.ToLower(e.Name()), ".gguf") {
				rez.ModeleDisponibile = append(rez.ModeleDisponibile, e.Name())
				if rez.CaleModel == "" {
					rez.CaleModel = filepath.Join(folderAI, e.Name())
				}
			}
		}
	}

	rez.Instalat = rez.CaleBinar != "" && rez.CaleModel != ""
	rez.Activ = VerificaSanatate(port)

	switch {
	case rez.Activ:
		rez.Mesaj = "Motorul local este activ și gata de utilizare."
	case rez.Instalat:
		rez.Mesaj = "Fișierele locale sunt instalate. Motorul poate fi pornit."
	case rez.CaleBinar == "" && len(rez.ModeleDisponibile) == 0:
		rez.Mesaj = fmt.Sprintf("Lipsește executabilul (%s) și modelul (*.gguf) din %s.", NumeBinarServer(), folderAI)
	case rez.CaleBinar == "":
		rez.Mesaj = fmt.Sprintf("Lipsește executabilul (%s) din %s.", NumeBinarServer(), folderAI)
	default:
		rez.Mesaj = fmt.Sprintf("Lipsește fișierul model (*.gguf) din %s.", folderAI)
	}

	return rez
}

// VerificaSanatate verifică dacă serverul local răspunde la endpoint-ul /health.
func VerificaSanatate(port int) bool {
	ctx, cancel := context.WithTimeout(context.Background(), 800*time.Millisecond)
	defer cancel()

	url := fmt.Sprintf("http://127.0.0.1:%d/health", port)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return false
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return false
	}
	defer resp.Body.Close()

	return resp.StatusCode == http.StatusOK
}
