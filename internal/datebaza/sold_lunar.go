package datebaza

import (
	"context"
	"pontix/internal/bd/acces"
	"pontix/internal/calcul"
	"sort"
	"time"
)

func completeazaSolduri(ctx context.Context, q *acces.Queries, an, luna int, rez *calcul.RezultatPontaj) error {
	norme, e := q.NormeSoldPontaj(ctx)
	if e != nil {
		return e
	}
	init, e := q.ListaInitializariPontaj(ctx)
	if e != nil {
		return e
	}
	prima, e := q.PrimaLunaPontaj(ctx)
	if e != nil {
		return e
	}
	curenta := time.Date(an, time.Month(luna), 1, 0, 0, 0, 0, time.UTC)
	if prima == "" {
		prima = curenta.Format("2006-01")
	}
	for _, v := range init {
		if v.Luna < prima {
			prima = v.Luna
		}
	}
	start, e := time.Parse("2006-01", prima)
	if e != nil {
		return e
	}
	if start.After(curenta) {
		start = curenta
	}
	sarb, e := q.ListaToateSarbatori(ctx)
	if e != nil {
		return e
	}
	libere := map[string]bool{}
	for _, s := range sarb {
		libere[s.Data] = true
	}
	tipuri, e := q.ListaTipuriAbsenta(ctx)
	if e != nil {
		return e
	}
	categorii := map[string]string{}
	for _, t := range tipuri {
		cat := t.Categorie
		switch cat {
		case "ODIHNA":
			cat = "CO"
		case "MEDICAL":
			cat = "CM"
		case "SPECIAL":
			cat = t.Cod
		}
		if t.Cod == "CFS" || t.Cod == "CFP" { // CFP: codul redenumit (D092)
			cat = "CFS"
		}
		categorii[t.Cod] = cat
	}
	config, e := citesteConfiguratiePontaj(ctx, q)
	if e != nil {
		return e
	}
	for _, r := range config.Reguli {
		if r.Sursa == "ZILE_ABSENTA" && (r.Coloana == "CEV" || r.Coloana == "CFS") {
			for _, cod := range r.Coduri {
				categorii[cod] = r.Coloana
			}
		}
	}
	efecte, e := q.EfecteAbsente(ctx)
	if e != nil {
		return e
	}
	for cod, efect := range efecte {
		if efect == "FARA_ORE" {
			categorii[cod] = "CO"
		} else if efect == "NORMA" || efect == "PROGRAM" {
			categorii[cod] = ""
		}
	}
	sold := map[string]int64{}
	inceputCalcul := map[string]bool{}
	definit := map[string]bool{}
	for id := range norme {
		definit[id] = true
	}
	for m := start; !m.After(curenta); m = m.AddDate(0, 1, 0) {
		lunaText := m.Format("2006-01")
		deLa := lunaText + "-01"
		panaLa := m.AddDate(0, 1, -1).Format(time.DateOnly)
		baza := *rez
		if !m.Equal(curenta) {
			baza, _, e = calculeazaPontajBaza(ctx, q, m.Year(), int(m.Month()), nil, nil)
			if e != nil {
				return e
			}
		}
		minute := map[string]int64{}
		complet := map[string]bool{}
		prezent := map[string]bool{}
		for _, z := range baza.Celule {
			prezent[z.AngajatID] = true
			minute[z.AngajatID] += z.Minute
			if !z.Calculabil {
				complet[z.AngajatID] = true
			}
		}
		zile, e := q.ListaAbsenteInterval(ctx, acces.ListaAbsenteIntervalParams{DeLa: deLa, PanaLa: panaLa})
		if e != nil {
			return e
		}
		abs := map[string][]calcul.AbsentaNorma{}
		for _, z := range zile {
			if z.Activ == 1 {
				prezent[z.AngajatID] = true
				abs[z.AngajatID] = append(abs[z.AngajatID], calcul.AbsentaNorma{Data: z.Data, Tip: categorii[z.Cod], Aprobata: z.StareAprobare == "APROBATA"})
			}
		}
		valori := map[string]map[string]calcul.ValoarePontaj{}
		if m.Equal(curenta) {
			for _, t := range rez.Totaluri {
				valori[t.AngajatID] = t.Valori
			}
		}
		for id, norma := range norme {
			v, areInit := init[id]
			if areInit && lunaText < v.Luna {
				continue
			}
			if areInit && lunaText == v.Luna {
				sold[id] = v.Minute
				inceputCalcul[id] = true
			}
			if !inceputCalcul[id] && !prezent[id] && !m.Equal(curenta) {
				continue
			}
			inceputCalcul[id] = true
			n, e := calcul.CalculeazaNormaDatorata(deLa, panaLa, norma, libere, abs[id])
			if e != nil {
				return e
			}
			inceputDefinit := definit[id] && n.Definita
			r := calcul.CalculeazaSoldLunar(sold[id], minute[id], n, definit[id] && !complet[id])
			definit[id] = r.Definit
			if m.Equal(curenta) {
				v := valori[id]
				if v == nil {
					v = map[string]calcul.ValoarePontaj{}
					rez.Totaluri = append(rez.Totaluri, calcul.TotalPontaj{AngajatID: id, Valori: v})
				}
				v["ore_lucratoare"] = calcul.ValoarePontaj{Sutimi: r.Lucratoare * 100, Unitate: "ORE", Definita: n.Definita}
				v["L_INCEPUT"] = calcul.ValoarePontaj{Sutimi: r.Inceput * 100, Unitate: "ORE", Definita: inceputDefinit}
				v["L_FIN"] = calcul.ValoarePontaj{Sutimi: r.Final * 100, Unitate: "ORE", Definita: r.Definit}
			}
			sold[id] = r.Final
		}
	}
	sort.Slice(rez.Totaluri, func(i, j int) bool { return rez.Totaluri[i].AngajatID < rez.Totaluri[j].AngajatID })
	return nil
}
