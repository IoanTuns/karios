package calcul

import (
	"fmt"
	"sort"
	"strings"
	"time"
)

// RegulaConfigurabila reprezintă o regulă configurată liber de utilizator (D041, I061).
type RegulaConfigurabila struct {
	ID              string `json:"id"`
	Cod             string `json:"cod"`
	Denumire        string `json:"denumire"`
	Descriere       string `json:"descriere"`
	TipRegula       string `json:"tip_regula"`       // REPAUS_MINIM, PLAFON_LUNAR, SUCCESIUNE_INTERZISA
	Severitate      string `json:"severitate"`       // INFO, AVERTISMENT, BLOCANT
	CoduriSchimb    string `json:"coduri_schimb"`    // ex: "G24" sau "G24,Z7"
	ValoareNumerica int    `json:"valoare_numerica"` // ex: 24 (ore repaus) sau 6 (maxim gărzi/lună)
	Activ           bool   `json:"activ"`
}

// IncalcareRegula conține detaliile unei reguli încălcate în planificare sau pontaj.
type IncalcareRegula struct {
	RegulaID   string `json:"regula_id"`
	RegulaCod  string `json:"regula_cod"`
	Severitate string `json:"severitate"`
	AngajatID  string `json:"angajat_id"`
	Nume       string `json:"nume"`
	Data       string `json:"data"` // AAAA-LL-ZZ sau lună
	Mesaj      string `json:"mesaj"`
}

// PozitiePentruRegula conține datele minime necesare evaluării regulilor.
type PozitiePentruRegula struct {
	AngajatID string `json:"angajat_id"`
	Nume      string `json:"nume"`
	Data      string `json:"data"`
	CodSchimb string `json:"cod_schimb"`
	OraStart  int    `json:"ora_start"`     // minute de la miezul nopții
	DurataMin int    `json:"durata_minute"` // durata schimbului în minute
	Activ     bool   `json:"activ"`
}

// EvalueazaReguli verifică o listă de reguli active peste turele planificate/efectuate.
func EvalueazaReguli(reguli []RegulaConfigurabila, pozitii []PozitiePentruRegula) []IncalcareRegula {
	var incalcari []IncalcareRegula
	if len(reguli) == 0 || len(pozitii) == 0 {
		return incalcari
	}

	// Grupăm pe angajat
	peAngajat := make(map[string][]PozitiePentruRegula)
	for _, p := range pozitii {
		if p.Activ && p.CodSchimb != "" {
			peAngajat[p.AngajatID] = append(peAngajat[p.AngajatID], p)
		}
	}

	// Sortăm turele fiecărui angajat cronologic
	for angID := range peAngajat {
		sort.Slice(peAngajat[angID], func(i, j int) bool {
			return peAngajat[angID][i].Data < peAngajat[angID][j].Data
		})
	}

	for _, reg := range reguli {
		if !reg.Activ {
			continue
		}

		coduri := strings.Split(strings.ToUpper(reg.CoduriSchimb), ",")
		for i := range coduri {
			coduri[i] = strings.TrimSpace(coduri[i])
		}

		switch reg.TipRegula {
		case "REPAUS_MINIM":
			repausOre := reg.ValoareNumerica
			if repausOre <= 0 {
				repausOre = 24
			}
			codVizat := coduri[0]

			for _, lista := range peAngajat {
				for idx, tura := range lista {
					if !strings.EqualFold(tura.CodSchimb, codVizat) {
						continue
					}
					dCurenta, err1 := time.Parse(time.DateOnly, tura.Data)
					if err1 != nil {
						continue
					}
					startMin := tura.OraStart
					durataMin := tura.DurataMin
					if durataMin <= 0 && strings.EqualFold(tura.CodSchimb, "G24") {
						durataMin = 1440
					}
					sfarsitTura := dCurenta.Add(time.Duration(startMin+durataMin) * time.Minute)

					for urmIdx := idx + 1; urmIdx < len(lista); urmIdx++ {
						urmTura := lista[urmIdx]
						dUrm, err2 := time.Parse(time.DateOnly, urmTura.Data)
						if err2 != nil {
							continue
						}
						startUrm := dUrm.Add(time.Duration(urmTura.OraStart) * time.Minute)
						if startUrm.Before(sfarsitTura) {
							continue
						}
						oreRepaus := startUrm.Sub(sfarsitTura).Hours()
						if oreRepaus < float64(repausOre) {
							incalcari = append(incalcari, IncalcareRegula{
								RegulaID:   reg.ID,
								RegulaCod:  reg.Cod,
								Severitate: reg.Severitate,
								AngajatID:  tura.AngajatID,
								Nume:       tura.Nume,
								Data:       tura.Data,
								Mesaj: fmt.Sprintf(
									"Repaus insuficient după %s pe %s: are alocat %s pe %s (doar %.0fh repaus, minim necesar: %dh).",
									tura.CodSchimb, tura.Data, urmTura.CodSchimb, urmTura.Data, oreRepaus, repausOre,
								),
							})
						}
						break
					}
				}
			}

		case "PLAFON_LUNAR":
			maxTure := reg.ValoareNumerica
			for _, lista := range peAngajat {
				numar := 0
				var numeAngajat string
				for _, tura := range lista {
					numeAngajat = tura.Nume
					for _, c := range coduri {
						if strings.EqualFold(tura.CodSchimb, c) {
							numar++
							break
						}
					}
				}
				if numar > maxTure && len(lista) > 0 {
					incalcari = append(incalcari, IncalcareRegula{
						RegulaID:   reg.ID,
						RegulaCod:  reg.Cod,
						Severitate: reg.Severitate,
						AngajatID:  lista[0].AngajatID,
						Nume:       numeAngajat,
						Data:       lista[0].Data[:7],
						Mesaj: fmt.Sprintf(
							"Plafon depășit pentru %s: %d ture planificate în luna %s (maxim permis: %d).",
							reg.CoduriSchimb, numar, lista[0].Data[:7], maxTure,
						),
					})
				}
			}

		case "SUCCESIUNE_INTERZISA":
			if len(coduri) >= 2 {
				t1Cod, t2Cod := coduri[0], coduri[1]
				for _, lista := range peAngajat {
					for idx := 0; idx < len(lista)-1; idx++ {
						t1 := lista[idx]
						t2 := lista[idx+1]
						if strings.EqualFold(t1.CodSchimb, t1Cod) && strings.EqualFold(t2.CodSchimb, t2Cod) {
							d1, _ := time.Parse(time.DateOnly, t1.Data)
							d2, _ := time.Parse(time.DateOnly, t2.Data)
							if d2.Sub(d1).Hours() == 24 {
								incalcari = append(incalcari, IncalcareRegula{
									RegulaID:   reg.ID,
									RegulaCod:  reg.Cod,
									Severitate: reg.Severitate,
									AngajatID:  t1.AngajatID,
									Nume:       t1.Nume,
									Data:       t2.Data,
									Mesaj: fmt.Sprintf(
										"Succesiune interzisă: tura %s din %s este urmată imediat de %s pe %s.",
										t1.CodSchimb, t1.Data, t2.CodSchimb, t2.Data,
									),
								})
							}
						}
					}
				}
			}
		}
	}

	return incalcari
}
