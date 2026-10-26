//go:build !windows

package main

import "pontix/internal/i18n"

func deschidePDFAplicatie(string) error { return &i18n.Eroare{Cod: "rapoarte.deschidere_esuata"} }
