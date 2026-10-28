package calcul

import (
	"fmt"
	"strings"
	"time"
)

type ZiComparativ struct {
	Data        string `json:"data"`
	Zi          int    `json:"zi"`
	PlanValoare string `json:"plan_valoare"`
	PlanOre     int    `json:"plan_ore"`
	RealValoare string `json:"real_valoare"`
	RealOre     int    `json:"real_ore"`
	Stare       string `json:"stare"` // "IDENTIC", "MODIFICAT", "NEPLANIFICAT", "NEEFECTUAT", "ABSENTA", "LIBER"
}

type RandComparativAngajat struct {
	AngajatID           string         `json:"angajat_id"`
	Nume                string         `json:"nume"`
	Prenume             string         `json:"prenume"`
	Marca               string         `json:"marca"`
	Functie             string         `json:"functie"`
	TotalOrePlanificate int            `json:"total_ore_planificate"`
	TotalOreRealizate   int            `json:"total_ore_realizate"`
	DiferentaOre        int            `json:"diferenta_ore"`
	ZileModificate      int            `json:"zile_modificate"`
	Zile                []ZiComparativ `json:"zile"`
}

type RaportPlanificatVsRealizat struct {
	An                  int                     `json:"an"`
	Luna                int                     `json:"luna"`
	TotalPlanificatOre  int                     `json:"total_planificat_ore"`
	TotalRealizatOre    int                     `json:"total_realizat_ore"`
	DiferentaGlobalaOre int                     `json:"diferenta_globala_ore"`
	TotalDiscrepante    int                     `json:"total_discrepante"`
	Randuri             []RandComparativAngajat `json:"randuri"`
}

type AngajatComparativ struct {
	ID      string
	Nume    string
	Prenume string
	Marca   string
	Functie string
	Activ   bool
}

type PozitieComparativ struct {
	AngajatID string
	Data      string
	CodSchimb string
	Ore       int
	Activ     bool
}

type CelulaComparativ struct {
	AngajatID string
	Data      string
	Valoare   string
	Ore       int
}

type AbsentaComparativ struct {
	AngajatID string
	Data      string
	Cod       string
}

// CalculeazaPlanificatVsRealizat generează analiza comparativă între tura planificată și datele reale din pontaj.
func CalculeazaPlanificatVsRealizat(
	an, luna int,
	angajati []AngajatComparativ,
	pozitii []PozitieComparativ,
	celule []CelulaComparativ,
	absente []AbsentaComparativ,
) RaportPlanificatVsRealizat {
	numarZile := time.Date(an, time.Month(luna)+1, 0, 12, 0, 0, 0, time.UTC).Day()

	hartaPlan := make(map[string]PozitieComparativ)
	for _, p := range pozitii {
		if p.Activ {
			hartaPlan[p.AngajatID+"/"+p.Data] = p
		}
	}

	hartaReal := make(map[string]CelulaComparativ)
	for _, c := range celule {
		if strings.TrimSpace(c.Valoare) != "" {
			hartaReal[c.AngajatID+"/"+c.Data] = c
		}
	}

	hartaAbsente := make(map[string]string)
	for _, a := range absente {
		hartaAbsente[a.AngajatID+"/"+a.Data] = a.Cod
	}

	rap := RaportPlanificatVsRealizat{
		An:      an,
		Luna:    luna,
		Randuri: make([]RandComparativAngajat, 0, len(angajati)),
	}

	for _, a := range angajati {
		if !a.Activ {
			continue
		}
		rand := RandComparativAngajat{
			AngajatID: a.ID,
			Nume:      a.Nume,
			Prenume:   a.Prenume,
			Marca:     a.Marca,
			Functie:   a.Functie,
			Zile:      make([]ZiComparativ, 0, numarZile),
		}

		for zi := 1; zi <= numarZile; zi++ {
			data := fmt.Sprintf("%04d-%02d-%02d", an, luna, zi)
			k := a.ID + "/" + data

			p, arePlan := hartaPlan[k]
			r, areReal := hartaReal[k]
			absCod, areAbs := hartaAbsente[k]

			zc := ZiComparativ{
				Data: data,
				Zi:   zi,
			}

			if arePlan {
				zc.PlanValoare = p.CodSchimb
				zc.PlanOre = p.Ore
			}
			if areReal {
				zc.RealValoare = r.Valoare
				zc.RealOre = r.Ore
			}

			if areAbs {
				zc.Stare = "ABSENTA"
				zc.RealValoare = absCod
				if arePlan {
					rand.ZileModificate++
					rap.TotalDiscrepante++
				}
			} else if !arePlan && !areReal {
				zc.Stare = "LIBER"
			} else if arePlan && areReal {
				if strings.EqualFold(zc.PlanValoare, zc.RealValoare) {
					zc.Stare = "IDENTIC"
				} else {
					zc.Stare = "MODIFICAT"
					rand.ZileModificate++
					rap.TotalDiscrepante++
				}
			} else if arePlan && !areReal {
				zc.Stare = "NEEFECTUAT"
				rand.ZileModificate++
				rap.TotalDiscrepante++
			} else { // !arePlan && areReal
				zc.Stare = "NEPLANIFICAT"
				rand.ZileModificate++
				rap.TotalDiscrepante++
			}

			rand.TotalOrePlanificate += zc.PlanOre
			rand.TotalOreRealizate += zc.RealOre
			rand.Zile = append(rand.Zile, zc)
		}

		rand.DiferentaOre = rand.TotalOreRealizate - rand.TotalOrePlanificate
		rap.TotalPlanificatOre += rand.TotalOrePlanificate
		rap.TotalRealizatOre += rand.TotalOreRealizate

		rap.Randuri = append(rap.Randuri, rand)
	}

	rap.DiferentaGlobalaOre = rap.TotalRealizatOre - rap.TotalPlanificatOre
	return rap
}
