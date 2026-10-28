package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"pontix/internal/ai"
	"pontix/internal/versiune"
)

// CerinteSistem descrie cerințele minime și recomandate pentru rularea Pontix.
type CerinteSistem struct {
	SistemOperare string `json:"sistem_operare"`
	Procesor      string `json:"procesor"`
	RAMMinim      string `json:"ram_minim"`
	RAMRecomandat string `json:"ram_recomandat"`
	SpatiuDisc    string `json:"spatiu_disc"`
	Dependinte    string `json:"dependinte"`
	Retea         string `json:"retea"`
}

// InfoDetaliatLLM descrie caracteristicile motorului și modelului local AI.
type InfoDetaliatLLM struct {
	NumeModel        string `json:"nume_model"`
	Format           string `json:"format"`
	MotorInferenta   string `json:"motor_inferenta"`
	FereastraContext int    `json:"fereastra_context"`
	PortLocal        int    `json:"port_local"`
	Stare            string `json:"stare"`
	DimensiuneModel  string `json:"dimensiune_model"`
	IzolareRetea     string `json:"izolare_retea"`
}

// DetaliiDezvoltare conține metadatele complete despre repo, build, integritate și cerințe.
type DetaliiDezvoltare struct {
	NumeRepo          string          `json:"nume_repo"`
	NumeAplicatie     string          `json:"nume_aplicatie"`
	Versiune          string          `json:"versiune"`
	Commit            string          `json:"commit"`
	Arhitectura       string          `json:"arhitectura"`
	VersiuneGo        string          `json:"versiune_go"`
	VersiuneWails     string          `json:"versiune_wails"`
	CaleExecutabil    string          `json:"cale_executabil"`
	ChecksumExe       string          `json:"checksum_exe"`
	Autor             string          `json:"autor"`
	Licenta           string          `json:"licenta"`
	Confidentialitate string          `json:"confidentialitate"`
	CerinteMinime     CerinteSistem   `json:"cerinte_minime"`
	DetaliiLLM        InfoDetaliatLLM `json:"detalii_llm"`
}

func calculeazaChecksumExe() (string, string) {
	cale, err := os.Executable()
	if err != nil {
		return "", "indisponibil"
	}
	caleLower := strings.ToLower(cale)
	if strings.HasSuffix(caleLower, ".test.exe") || strings.Contains(caleLower, "\\go-build") || strings.Contains(caleLower, "/go-build") {
		return cale, "Mod test / dezvoltare (fără binar final de lansare)"
	}
	f, err := os.Open(cale)
	if err != nil {
		return cale, "indisponibil (eroare acces)"
	}
	defer f.Close()

	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return cale, "indisponibil (eroare citire)"
	}
	return cale, strings.ToUpper(hex.EncodeToString(h.Sum(nil)))
}

func formateazaOcteti(octeti int64) string {
	const unitate = 1024
	if octeti < unitate {
		return fmt.Sprintf("%d B", octeti)
	}
	div, exp := int64(unitate), 0
	for n := octeti / unitate; n >= unitate; n /= unitate {
		div *= unitate
		exp++
	}
	val := float64(octeti) / float64(div)
	litere := []string{"KB", "MB", "GB", "TB"}
	if exp < len(litere) {
		return fmt.Sprintf("%.2f %s", val, litere[exp])
	}
	return fmt.Sprintf("%.2f GB", val)
}

// DetaliiDezvoltare întoarce raportul tehnic complet de dezvoltare, compilare, integritate și mediu AI.
func (a *App) DetaliiDezvoltare() DetaliiDezvoltare {
	caleExe, shaExe := calculeazaChecksumExe()
	folderAI := a.folderAI()
	stareAI := ai.DetecteazaMediu(folderAI, ai.PortImplicit)

	dimModel := "Nedetectat"
	numeModel := "Niciun model instalat"
	if stareAI.CaleModel != "" {
		numeModel = filepath.Base(stareAI.CaleModel)
		if fi, err := os.Stat(stareAI.CaleModel); err == nil {
			dimModel = formateazaOcteti(fi.Size())
		}
	} else if len(stareAI.ModeleDisponibile) > 0 {
		numeModel = stareAI.ModeleDisponibile[0]
		if fi, err := os.Stat(filepath.Join(folderAI, numeModel)); err == nil {
			dimModel = formateazaOcteti(fi.Size())
		}
	}

	stareText := "Neinstalat"
	if stareAI.Instalat {
		if stareAI.Activ {
			stareText = "Activ și conectat"
		} else {
			stareText = "Instalat (oprit)"
		}
	}

	motorText := "llama-server (llama.cpp C++ nativ cu AVX2)"
	if !stareAI.Instalat || stareAI.CaleBinar == "" {
		motorText = "Neinstalat în karios/ai/"
	}

	return DetaliiDezvoltare{
		NumeRepo:          "karios",
		NumeAplicatie:     "Pontix",
		Versiune:          versiune.Aplicatie,
		Commit:            versiune.Commit,
		Arhitectura:       fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH),
		VersiuneGo:        runtime.Version(),
		VersiuneWails:     "v2.17.0",
		CaleExecutabil:    caleExe,
		ChecksumExe:       shaExe,
		Autor:             "Ioan Tuns",
		Licenta:           "Proprietară / Drepturi rezervate © 2026 Ioan Tuns",
		Confidentialitate: "100% Offline (Date medicale și personale păstrate exclusiv local, fără cloud/telemetrie)",
		CerinteMinime: CerinteSistem{
			SistemOperare: "Windows 10 (versiunea 1809 sau mai nouă) / Windows 11 (64-biți)",
			Procesor:      "Procesor x64 pe 64-biți (recomandat quad-core cu suport AVX2 pentru AI)",
			RAMMinim:      "4 GB RAM (modul de bază) / 8 GB RAM (cu Asistentul AI local)",
			RAMRecomandat: "16 GB RAM pentru comutare și inferență fluidă",
			SpatiuDisc:    "100 MB aplicația de bază / ~3.5 GB cu motorul C++ și modelul AI local",
			Dependinte:    "Microsoft Edge WebView2 Runtime (inclus nativ în Windows 10/11)",
			Retea:         "Nu este necesară conexiune la internet (0 KB/s trafic exterior)",
		},
		DetaliiLLM: InfoDetaliatLLM{
			NumeModel:        numeModel,
			Format:           "GGUF (cuantizare 4-bit optimizată CPU/RAM)",
			MotorInferenta:   motorText,
			FereastraContext: 8192,
			PortLocal:        stareAI.Port,
			Stare:            stareText,
			DimensiuneModel:  dimModel,
			IzolareRetea:     "127.0.0.1 (Strict local, zero conexiuni exterioare, conformitate GDPR)",
		},
	}
}
