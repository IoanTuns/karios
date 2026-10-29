package datebaza

import "sort"

type SumarPlanAI struct {
	Cheie     string `json:"cheie"`
	Nume      string `json:"nume"`
	Selectati int    `json:"selectati"`
	Afectati  int    `json:"afectati"`
	Zile      int    `json:"zile"`
	Noi       int    `json:"noi"`
	Inlocuiri int    `json:"inlocuiri"`
	Pastrate  int    `json:"pastrate"`
	Blocate   int    `json:"blocate"`
	Identice  int    `json:"identice"`
}

func rezumaPlanAI(p PlanAI) ([]SumarPlanAI, []SumarPlanAI) {
	luni, persoane := map[string]*SumarPlanAI{}, map[string]*SumarPlanAI{}
	afectati, zile := map[string]map[string]bool{}, map[string]map[string]bool{}
	for _, a := range p.Persoane {
		persoane[a.ID] = &SumarPlanAI{Cheie: a.Marca, Nume: a.Nume, Selectati: 1}
		zile[a.ID] = map[string]bool{}
	}
	for _, l := range p.Luni {
		luni[l] = &SumarPlanAI{Cheie: l, Nume: l}
		afectati[l] = map[string]bool{}
		zile[l] = map[string]bool{}
	}
	for _, r := range p.Randuri {
		l := r.Data[:7]
		zile[l][r.Data] = true
		zile[r.AngajatID][r.Data] = true
		modifica := r.Stare == "noua" || r.Stare == "inlocuire"
		if modifica {
			afectati[l][r.AngajatID] = true
			persoane[r.AngajatID].Afectati = 1
		}
		for _, s := range []*SumarPlanAI{luni[l], persoane[r.AngajatID]} {
			switch r.Stare {
			case "noua":
				s.Noi++
			case "inlocuire":
				s.Inlocuiri++
			case "ocupata":
				s.Pastrate++
			case "identica":
				s.Identice++
			default:
				s.Blocate++
			}
		}
	}
	rl, rp := []SumarPlanAI{}, []SumarPlanAI{}
	for l, s := range luni {
		s.Selectati = len(p.Persoane)
		s.Afectati = len(afectati[l])
		s.Zile = len(zile[l])
		rl = append(rl, *s)
	}
	for id, s := range persoane {
		s.Zile = len(zile[id])
		rp = append(rp, *s)
	}
	sort.Slice(rl, func(i, j int) bool { return rl[i].Cheie < rl[j].Cheie })
	sort.Slice(rp, func(i, j int) bool { return rp[i].Cheie < rp[j].Cheie })
	return rl, rp
}
