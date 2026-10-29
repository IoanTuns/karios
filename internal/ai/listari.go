package ai

import (
	"fmt"
	"sort"
	"strings"
)

var expresiiListare = []string{"listeaza", "lista angajati", "lista cu angajati", "lista personal", "arata angajatii", "arata-mi angajatii", "arata personalul", "care sunt angajatii", "ce angajati", "cati angajati", "numar angajati", "numarul angajatilor"}

// esteListareAngajati recunoaște cererile de listare sau numărare a angajaților (F14-T09).
func esteListareAngajati(m string) bool {
	for _, e := range expresiiListare {
		if strings.Contains(m, e) {
			return strings.Contains(m, "angajat") || strings.Contains(m, "personal") || strings.Contains(m, "salariat") || strings.Contains(m, "oameni")
		}
	}
	return false
}

func esteListareSchimburi(m string) bool {
	return (strings.Contains(m, "listeaza") || strings.Contains(m, "ce schimburi") || strings.Contains(m, "care sunt schimburile") || strings.Contains(m, "lista schimburi")) &&
		strings.Contains(m, "schimb")
}

// sintetizeazaListaAngajati întoarce angajații activi, grupați pe secții; cu „secția X” sau
// „categoria Y” în cerere, doar acel grup.
func sintetizeazaListaAngajati(m string, info ContextPontix) string {
	angajati := info.Angajati
	eticheta := "Angajați activi"
	if g, ok := gasesteGrup(m, info); ok {
		angajati = g.angajati
		eticheta = "Angajați activi · " + g.eticheta
	}
	if len(angajati) == 0 {
		return "Nu există angajați activi pentru această cerere."
	}
	pe := map[string][]InfoAngajat{}
	for _, a := range angajati {
		pe[a.Sectie] = append(pe[a.Sectie], a)
	}
	sectii := make([]string, 0, len(pe))
	for s := range pe {
		sectii = append(sectii, s)
	}
	sort.Strings(sectii)
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("👥 **%s: %d**\n", eticheta, len(angajati)))
	for _, s := range sectii {
		lista := pe[s]
		sort.Slice(lista, func(i, j int) bool {
			return normalizeaza(lista[i].Nume+" "+lista[i].Prenume) < normalizeaza(lista[j].Nume+" "+lista[j].Prenume)
		})
		titlu := s
		if titlu == "" {
			titlu = "Fără secție"
		}
		sb.WriteString(fmt.Sprintf("\n**%s** (%d)\n", titlu, len(lista)))
		for _, a := range lista {
			linie := fmt.Sprintf("- %s %s (%s)", a.Nume, a.Prenume, a.Marca)
			if a.Functie != "" {
				linie += " · " + a.Functie
			}
			if a.Categorie != "" {
				linie += " · " + a.Categorie
			}
			sb.WriteString(linie + "\n")
		}
	}
	return sb.String()
}

func sintetizeazaListaSchimburi(info ContextPontix) string {
	if len(info.Schimburi) == 0 {
		return "Nu există schimburi active."
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("📋 **Schimburi active: %d**\n", len(info.Schimburi)))
	for _, s := range info.Schimburi {
		linie := fmt.Sprintf("- %s · %s", s.Cod, s.Denumire)
		if s.Interval != "" {
			linie += " · " + s.Interval
		}
		if s.Ore > 0 {
			linie += fmt.Sprintf(" · %d ore", s.Ore)
		}
		sb.WriteString(linie + "\n")
	}
	return sb.String()
}
