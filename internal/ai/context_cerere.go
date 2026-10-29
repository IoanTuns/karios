package ai

import "strings"

// ContextPentruCerere păstrează datele relevante, cu un buget conservator de
// caractere. Nu estimează tokenuri și nu pretinde că datele omise sunt complete.
// Rezolvarea selecțiilor și analizele exacte folosesc contextul complet în Go.
func ContextPentruCerere(mesaj string, info ContextPontix) ContextPontix {
	m := normalizeaza(mesaj)
	rez := info
	rez.Domeniu = "ajutor"
	switch {
	case strings.Contains(m, "conced") || strings.Contains(m, "absen"):
		rez.Domeniu = "absente"
	case strings.Contains(m, "pontaj") || strings.Contains(m, "ore") || strings.Contains(m, "statistic"):
		rez.Domeniu = "pontaj"
	case strings.Contains(m, "planific") || strings.Contains(m, "tura") || strings.Contains(m, "schimb") || strings.Contains(m, "lucre"):
		rez.Domeniu = "planificare"
	}
	rez.Angajati, rez.Absente, rez.TurePlanificate = nil, nil, nil
	rez.Schimburi = nil
	bugetPrograme := 2000
	for _, p := range info.Schimburi {
		cost := len(p.Cod+p.Denumire+p.Interval) + 80
		if cost > bugetPrograme {
			rez.ContextIncomplet = true
			continue
		}
		bugetPrograme -= cost
		rez.Schimburi = append(rez.Schimburi, p)
	}
	rez.ActiuniAnterioare = nil // Intențiile și previzualizările rămân în interfață, nu sunt instrucțiuni noi.
	selectati := map[string]bool{}
	for _, a := range info.Angajati {
		if contineExpresie(m, a.Marca) || contineExpresie(m, a.Nume+" "+a.Prenume) || contineExpresie(m, a.Prenume+" "+a.Nume) {
			selectati[a.ID] = true
		}
	}
	buget := 6000
	for _, a := range info.Angajati {
		if !selectati[a.ID] {
			continue
		}
		cost := len(a.ID+a.Nume+a.Prenume+a.Marca+a.Functie+a.Sectie+a.Categorie+a.Departament) + 100
		if cost > buget {
			rez.ContextIncomplet = true
			continue
		}
		buget -= cost
		rez.Angajati = append(rez.Angajati, a)
	}
	if len(rez.Angajati) < len(info.Angajati) {
		rez.ContextIncomplet = true
	}
	if rez.Domeniu == "absente" || rez.Domeniu == "pontaj" || rez.Domeniu == "planificare" {
		for _, a := range info.Absente {
			if !selectati[a.AngajatID] {
				continue
			}
			cost := len(a.AngajatID+a.Marca+a.Nume+a.Data+a.Cod) + 80
			if cost > buget {
				rez.ContextIncomplet = true
				continue
			}
			buget -= cost
			rez.Absente = append(rez.Absente, a)
		}
	}
	if rez.Domeniu == "pontaj" || rez.Domeniu == "planificare" {
		for _, p := range info.TurePlanificate {
			if !selectati[p.AngajatID] {
				continue
			}
			cost := len(p.AngajatID+p.Marca+p.Nume+p.Data+p.CodSchimb) + 80
			if cost > buget {
				rez.ContextIncomplet = true
				continue
			}
			buget -= cost
			rez.TurePlanificate = append(rez.TurePlanificate, p)
		}
	}
	return rez
}
