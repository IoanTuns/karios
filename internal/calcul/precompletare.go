package calcul

import (
	"strings"
	"time"
)

type AngajatPrecompletare struct {
	ID, CategorieID string
	NormaMinute     int64
	Activ           bool
}
type CategoriePrecompletare struct {
	ID, ImplicitID string
	Activ          bool
}
type ProgramPrecompletare struct {
	ID, Cod, Tip, Regim                 string
	OraStart, DurataMinute, PauzaMinute int64
	Activ                               bool
}
type CelulaPrecompletare struct {
	AngajatID string `json:"angajat_id"`
	Data      string `json:"data"`
	Valoare   string `json:"valoare"`
}
type PlanPrecompletare struct {
	Celule []CelulaPrecompletare `json:"celule"`
	Omise  map[string]int        `json:"omise"`
}

// PlanificaPrecompletare propune numai celule goale, fără rotații sau reguli salariale.
// Datele și intervalele sunt nominale; ocupate include și zilele pontajelor efective.
func PlanificaPrecompletare(deLa, panaLa string, angajati []AngajatPrecompletare, categorii []CategoriePrecompletare, programe []ProgramPrecompletare, intrari []IntrarePontaj, libere, absente, ocupate map[string]bool) (PlanPrecompletare, error) {
	rez := PlanPrecompletare{Celule: []CelulaPrecompletare{}, Omise: map[string]int{}}
	inceput, e1 := time.Parse(time.DateOnly, deLa)
	sfarsit, e2 := time.Parse(time.DateOnly, panaLa)
	if e1 != nil || e2 != nil || sfarsit.Before(inceput) || sfarsit.Sub(inceput) >= 31*24*time.Hour {
		return rez, &Eroare{Cod: "pontaj.perioada_invalida"}
	}
	profil := map[string]CategoriePrecompletare{}
	for _, c := range categorii {
		profil[c.ID] = c
	}
	schimb := map[string]ProgramPrecompletare{}
	coduri := map[string]ProgramPrecompletare{}
	for _, p := range programe {
		schimb[p.ID] = p
		coduri[strings.ToUpper(p.Cod)] = p
	}
	existente := map[string]bool{}
	neclare := map[string]bool{}
	garzi := map[string]bool{}
	type interval struct{ inceput, sfarsit time.Time }
	intervale := map[string][]interval{}
	for _, in := range intrari {
		if strings.TrimSpace(in.Valoare) == "" {
			continue
		}
		existente[in.AngajatID+"/"+in.Data] = true
		d, e := time.Parse(time.DateOnly, in.Data)
		if e != nil {
			return rez, &Eroare{Cod: "pontaj.celula_invalida"}
		}
		p, cunoscut := coduri[strings.ToUpper(strings.TrimSpace(in.Valoare))]
		if cunoscut && p.Tip == "GARDA" && p.DurataMinute >= 1440 {
			garzi[in.AngajatID+"/"+d.AddDate(0, 0, 1).Format(time.DateOnly)] = true
		}
		if cunoscut && p.OraStart >= 0 && p.OraStart < 1440 && p.DurataMinute > 0 && p.DurataMinute <= 2880 {
			start := d.Add(time.Duration(p.OraStart) * time.Minute)
			intervale[in.AngajatID] = append(intervale[in.AngajatID], interval{start, start.Add(time.Duration(p.DurataMinute) * time.Minute)})
		} else if _, manual := OreManuale(in.Valoare); !cunoscut && manual {
			// Orele manuale aparțin propriei zile, fără un interval de tură dedus.
			neclare[in.AngajatID+"/"+in.Data] = true
		} else {
			for i := 0; i <= 2; i++ {
				neclare[in.AngajatID+"/"+d.AddDate(0, 0, i).Format(time.DateOnly)] = true
			}
		}
	}
	for _, a := range angajati {
		c, areCategorie := profil[a.CategorieID]
		p, areProgram := schimb[c.ImplicitID]
		for d := inceput; !d.After(sfarsit); d = d.AddDate(0, 0, 1) {
			data := d.Format(time.DateOnly)
			k := a.ID + "/" + data
			motiv := ""
			switch {
			case !a.Activ:
				motiv = "inactiv"
			case d.Weekday() == time.Saturday || d.Weekday() == time.Sunday:
				motiv = "weekend"
			case libere[data]:
				motiv = "libera"
			case absente[k]:
				motiv = "absenta"
			case existente[k] || ocupate[k]:
				motiv = "existent"
			case !areCategorie:
				motiv = "fara_categorie"
			case !c.Activ:
				motiv = "categorie_inactiva"
			case c.ImplicitID == "":
				motiv = "fara_implicit"
			case !areProgram || !p.Activ:
				motiv = "schimb_inactiv"
			case p.OraStart < 0 || p.OraStart >= 1440 || p.DurataMinute <= 0 || p.DurataMinute > 2880 || p.PauzaMinute != 0:
				motiv = "program_incomplet"
			case p.Regim == "NELUCRATOARE" || p.Regim == "INTRARE_ZI_LIBERA" || p.Tip == "GARDA":
				motiv = "alegere_manuala"
			case p.DurataMinute != a.NormaMinute:
				motiv = "norma_diferita"
			case garzi[k]:
				motiv = "dupa_garda"
			}
			if motiv == "" {
				start := d.Add(time.Duration(p.OraStart) * time.Minute)
				stop := start.Add(time.Duration(p.DurataMinute) * time.Minute)
				for zi := d; zi.Before(stop); zi = zi.AddDate(0, 0, 1) {
					kz := a.ID + "/" + zi.Format(time.DateOnly)
					if absente[kz] {
						motiv = "conflict_absenta"
						break
					}
					if ocupate[kz] || neclare[kz] {
						motiv = "suprapunere"
						break
					}
					if zi.After(d) && (libere[zi.Format(time.DateOnly)] || zi.Weekday() == time.Saturday || zi.Weekday() == time.Sunday) {
						motiv = "trecere_zi_libera"
						break
					}
				}
				if motiv == "" {
					for _, i := range intervale[a.ID] {
						if start.Before(i.sfarsit) && stop.After(i.inceput) {
							motiv = "suprapunere"
							break
						}
					}
				}
				if motiv == "" {
					rez.Celule = append(rez.Celule, CelulaPrecompletare{AngajatID: a.ID, Data: data, Valoare: p.Cod})
					intervale[a.ID] = append(intervale[a.ID], interval{start, stop})
				}
			}
			if motiv != "" {
				rez.Omise[motiv]++
			}
		}
	}
	return rez, nil
}
