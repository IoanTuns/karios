package main

import (
	"net/url"
	"os"
	"path/filepath"
	"pontix/internal/ajutor"
	"pontix/internal/i18n"
	"pontix/internal/versiune"
	"strings"
)

func (a *App) ArticoleAjutor(limba string) []ajutor.Articol { return ajutor.Lista(i18n.Limba(limba)) }

// Se accepta numai ID-uri din catalog, niciodata cai furnizate de interfata.
func caleManualAjutor(radacina, id string) (string, *MesajEroare) {
	if _, ok := ajutor.Gaseste(id, i18n.Romana); !ok {
		return "", &MesajEroare{Cod: "ajutor.articol_lipsa"}
	}
	cale := filepath.Join(filepath.Dir(radacina), "manual-utilizare-ro-"+versiune.Aplicatie+".html")
	if info, err := os.Stat(cale); err != nil || info.IsDir() {
		return "", &MesajEroare{Cod: "ajutor.manual_lipsa"}
	}
	cale, err := filepath.Abs(cale)
	if err != nil {
		return "", &MesajEroare{Cod: "ajutor.manual_lipsa"}
	}
	c := filepath.ToSlash(cale)
	if !strings.HasPrefix(c, "/") {
		c = "/" + c
	}
	return (&url.URL{Scheme: "file", Path: c, Fragment: "ajutor-" + id}).String(), nil
}
func (a *App) DeschideManualAjutor(id string) *MesajEroare {
	cale, err := caleManualAjutor(a.radacina, id)
	if err != nil {
		return err
	}
	if deschidePDFAplicatie(cale) != nil {
		return &MesajEroare{Cod: "ajutor.manual_lipsa"}
	}
	return nil
}
