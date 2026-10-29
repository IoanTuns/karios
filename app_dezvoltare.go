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

	"pontix/internal/diagnoza"
	"pontix/internal/i18n"
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
	NumeRepo           string                      `json:"nume_repo"`
	NumeAplicatie      string                      `json:"nume_aplicatie"`
	Versiune           string                      `json:"versiune"`
	Commit             string                      `json:"commit"`
	Arhitectura        string                      `json:"arhitectura"`
	VersiuneGo         string                      `json:"versiune_go"`
	VersiuneWails      string                      `json:"versiune_wails"`
	CaleExecutabil     string                      `json:"cale_executabil"`
	ChecksumExe        string                      `json:"checksum_exe"`
	Autor              string                      `json:"autor"`
	Licenta            string                      `json:"licenta"`
	Confidentialitate  string                      `json:"confidentialitate"`
	CerinteMinime      CerinteSistem               `json:"cerinte_minime"`
	DetaliiLLM         InfoDetaliatLLM             `json:"detalii_llm"`
	StatisticiDiagnoza diagnoza.StatisticiDiagnoza `json:"statistici_diagnoza"`
}

func calculeazaChecksumExe(limba i18n.Limba) (string, string) {
	cale, err := os.Executable()
	if err != nil {
		return "", i18n.T(limba, "dezvoltare.indisponibil", nil)
	}
	caleLower := strings.ToLower(cale)
	if strings.HasSuffix(caleLower, ".test.exe") || strings.Contains(caleLower, "\\go-build") || strings.Contains(caleLower, "/go-build") {
		return cale, i18n.T(limba, "dezvoltare.mod_test", nil)
	}
	f, err := os.Open(cale)
	if err != nil {
		return cale, i18n.T(limba, "dezvoltare.indisponibil_acces", nil)
	}
	defer f.Close()

	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return cale, i18n.T(limba, "dezvoltare.indisponibil_citire", nil)
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
// limba alege limba valorilor descriptive (F03-T10); etichetele sunt în interfață.
func (a *App) DetaliiDezvoltare(limbaInterfata string) DetaliiDezvoltare {
	limba := i18n.Limba(limbaInterfata)
	if limba != i18n.Romana && limba != i18n.Engleza {
		limba = i18n.Implicita
	}
	caleExe, shaExe := calculeazaChecksumExe(limba)
	folderAI := a.folderAI()
	stareAI := a.stareMediuAI()

	dimModel := i18n.T(limba, "dezvoltare.nedetectat", nil)
	numeModel := i18n.T(limba, "dezvoltare.fara_model", nil)
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

	stareText := i18n.T(limba, "dezvoltare.neinstalat", nil)
	if stareAI.Instalat {
		if stareAI.Activ {
			stareText = i18n.T(limba, "dezvoltare.activ", nil)
		} else {
			stareText = i18n.T(limba, "dezvoltare.oprit", nil)
		}
	}

	motorText := i18n.T(limba, "dezvoltare.motor_descriere", nil)
	if !stareAI.Instalat || stareAI.CaleBinar == "" {
		motorText = i18n.T(limba, "dezvoltare.motor_lipsa", nil)
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
		Licenta:           "Apache License 2.0 (Apache-2.0)",
		Confidentialitate: i18n.T(limba, "dezvoltare.confidentialitate", nil),
		CerinteMinime: CerinteSistem{
			SistemOperare: i18n.T(limba, "dezvoltare.cerinta_os", nil),
			Procesor:      i18n.T(limba, "dezvoltare.cerinta_cpu", nil),
			RAMMinim:      i18n.T(limba, "dezvoltare.cerinta_ram_minim", nil),
			RAMRecomandat: i18n.T(limba, "dezvoltare.cerinta_ram_recomandat", nil),
			SpatiuDisc:    i18n.T(limba, "dezvoltare.cerinta_disc", nil),
			Dependinte:    i18n.T(limba, "dezvoltare.cerinta_dependinte", nil),
			Retea:         i18n.T(limba, "dezvoltare.cerinta_retea", nil),
		},
		DetaliiLLM: InfoDetaliatLLM{
			NumeModel:        numeModel,
			Format:           i18n.T(limba, "dezvoltare.format_gguf", nil),
			MotorInferenta:   motorText,
			FereastraContext: 8192,
			PortLocal:        stareAI.Port,
			Stare:            stareText,
			DimensiuneModel:  dimModel,
			IzolareRetea:     i18n.T(limba, "dezvoltare.izolare_descriere", nil),
		},
		StatisticiDiagnoza: diagnoza.ObtineManager().Statistici(),
	}
}
