package calcul

import (
	"strings"
	"time"
)

// ProgramDePlanificat identifică schimburile gestionate în secțiunea separată.
func ProgramDePlanificat(p ProgramPrecompletare) bool {
	return p.Tip == "GARDA" || p.Tip == "NOAPTE" || p.Regim == "NELUCRATOARE" || p.Regim == "INTRARE_ZI_LIBERA" || (p.OraStart >= 720 && p.OraStart < 1320)
}

func ProgramInSectiune(p ProgramPrecompletare, sectiune string) bool {
	if sectiune == "GARZI" {
		return p.Tip == "GARDA" || p.Tip == "NOAPTE" || p.Regim == "NELUCRATOARE" || p.Regim == "INTRARE_ZI_LIBERA"
	}
	return sectiune == "DUPA_AMIAZA" && p.Tip != "GARDA" && p.Tip != "NOAPTE" && p.Regim != "NELUCRATOARE" && p.Regim != "INTRARE_ZI_LIBERA" && p.OraStart >= 720 && p.OraStart < 1320
}

type IntervalPlanificare struct {
	AngajatID, Cheie string
	Inceput, Sfarsit time.Time
}

// VerificaTuraPlanificata verifică intervale nominale; nu stabilește tarife sau repaus legal.
func VerificaTuraPlanificata(angajatID, data string, p ProgramPrecompletare, absente map[string]bool, intervale []IntervalPlanificare, ignora map[string]bool) string {
	d, err := time.Parse(time.DateOnly, data)
	if err != nil || p.OraStart < 0 || p.OraStart >= 1440 || p.DurataMinute <= 0 || p.DurataMinute > 2880 || p.PauzaMinute != 0 {
		return "program_incomplet"
	}
	start := d.Add(time.Duration(p.OraStart) * time.Minute)
	stop := start.Add(time.Duration(p.DurataMinute) * time.Minute)
	for zi := d; zi.Before(stop); zi = zi.AddDate(0, 0, 1) {
		if absente[angajatID+"/"+zi.Format(time.DateOnly)] {
			return "absenta"
		}
	}
	for _, i := range intervale {
		if i.AngajatID == angajatID && !ignora[i.Cheie] && start.Before(i.Sfarsit) && stop.After(i.Inceput) {
			return "suprapunere"
		}
	}
	return ""
}

type PozitiePlanificareSimple struct {
	ID        string `json:"id"`
	AngajatID string `json:"angajat_id"`
	Data      string `json:"data"`
	SchimbID  string `json:"schimb_id"`
	CodSchimb string `json:"cod_schimb"`
	Activ     bool   `json:"activ"`
}

type ZiAcoperire struct {
	Data               string              `json:"data"`
	ZiuaSaptamanii     int                 `json:"ziua_saptamanii"`
	EsteWeekend        bool                `json:"este_weekend"`
	EsteSarbatoare     bool                `json:"este_sarbatoare"`
	DenumireSarbatoare string              `json:"denumire_sarbatoare"`
	Posturi            map[string][]string `json:"posturi"`
	Lipsuri            []string            `json:"lipsuri"`
	Complet            bool                `json:"complet"`
}

type RaportAcoperire struct {
	An                    int           `json:"an"`
	Luna                  int           `json:"luna"`
	TotalZileNelucratoare int           `json:"total_zile_nelucratoare"`
	TotalPosturiNecesar   int           `json:"total_posturi_necesar"`
	TotalPosturiAcoperite int           `json:"total_posturi_acoperite"`
	TotalPosturiLipsa     int           `json:"total_posturi_lipsa"`
	Zile                  []ZiAcoperire `json:"zile"`
}

// AnalizeazaAcoperire verifică prezența personalului pe posturile obligatorii în zilele nelucrătoare.
func AnalizeazaAcoperire(an, luna int, sarbatori map[string]string, pozitii []PozitiePlanificareSimple, posturiObligatorii []string) RaportAcoperire {
	if len(posturiObligatorii) == 0 {
		posturiObligatorii = []string{"G24", "Z12", "N12", "IN9"}
	}
	rap := RaportAcoperire{
		An:   an,
		Luna: luna,
		Zile: []ZiAcoperire{},
	}
	primaZi := time.Date(an, time.Month(luna), 1, 12, 0, 0, 0, time.UTC)
	numarZile := time.Date(an, time.Month(luna)+1, 0, 12, 0, 0, 0, time.UTC).Day()

	hartaPozitii := map[string]map[string][]string{} // data -> codSchimb -> []angajatID
	for _, p := range pozitii {
		if !p.Activ {
			continue
		}
		if hartaPozitii[p.Data] == nil {
			hartaPozitii[p.Data] = map[string][]string{}
		}
		cod := p.CodSchimb
		hartaPozitii[p.Data][cod] = append(hartaPozitii[p.Data][cod], p.AngajatID)
	}

	for zi := 1; zi <= numarZile; zi++ {
		d := primaZi.AddDate(0, 0, zi-1)
		dataStr := d.Format(time.DateOnly)
		weekday := d.Weekday()
		esteWeekend := weekday == time.Saturday || weekday == time.Sunday
		denumireSarbatoare, esteSarbatoare := sarbatori[dataStr]

		if !esteWeekend && !esteSarbatoare {
			continue
		}
		rap.TotalZileNelucratoare++

		posturiZi := map[string][]string{}
		var lipsuri []string
		for _, post := range posturiObligatorii {
			rap.TotalPosturiNecesar++
			angajati := hartaPozitii[dataStr][post]
			posturiZi[post] = angajati
			if len(angajati) > 0 {
				rap.TotalPosturiAcoperite++
			} else {
				rap.TotalPosturiLipsa++
				lipsuri = append(lipsuri, post)
			}
		}

		rap.Zile = append(rap.Zile, ZiAcoperire{
			Data:               dataStr,
			ZiuaSaptamanii:     int(weekday),
			EsteWeekend:        esteWeekend,
			EsteSarbatoare:     esteSarbatoare,
			DenumireSarbatoare: denumireSarbatoare,
			Posturi:            posturiZi,
			Lipsuri:            lipsuri,
			Complet:            len(lipsuri) == 0,
		})
	}

	return rap
}

type AvertismentEchitate struct {
	AngajatID string `json:"angajat_id"`
	Data      string `json:"data"`
	CodSchimb string `json:"cod_schimb"`
	Tip       string `json:"tip"` // CONSECUTIV, ROTATIE_CRACIUN_REVELION, CONFLICT_CRACIUN_PASTI
	Mesaj     string `json:"mesaj"`
}

func EsteCraciun(data string) bool {
	return strings.HasSuffix(data, "-12-25") || strings.HasSuffix(data, "-12-26")
}

func EsteRevelion(data string) bool {
	return strings.HasSuffix(data, "-12-31") || strings.HasSuffix(data, "-01-01") || strings.HasSuffix(data, "-01-02")
}

func EstePasti(data string, sarbatori map[string]string) bool {
	den := strings.ToLower(sarbatori[data])
	return strings.Contains(den, "pașt") || strings.Contains(den, "past") || strings.Contains(den, "vinerea mare")
}

// AnalizeazaEchitate verifică rotația sărbătorilor majore și marcajele pentru ture consecutive.
func AnalizeazaEchitate(an, luna int, pozitii []PozitiePlanificareSimple, istoricAnTrecut, istoricAnViitor []PozitiePlanificareSimple, sarbatori map[string]string) []AvertismentEchitate {
	var avertismente []AvertismentEchitate

	pozitiiPeAngajat := map[string][]PozitiePlanificareSimple{}
	for _, p := range pozitii {
		if p.Activ {
			pozitiiPeAngajat[p.AngajatID] = append(pozitiiPeAngajat[p.AngajatID], p)
		}
	}

	craciunAnTrecut := map[string]bool{} // angajatID -> a fost de craciun
	for _, p := range istoricAnTrecut {
		if p.Activ && EsteCraciun(p.Data) {
			craciunAnTrecut[p.AngajatID] = true
		}
	}

	pastiAnViitor := map[string]string{} // angajatID -> codSchimb
	for _, p := range istoricAnViitor {
		if p.Activ && EstePasti(p.Data, sarbatori) {
			pastiAnViitor[p.AngajatID] = p.CodSchimb
		}
	}

	for angajatID, lista := range pozitiiPeAngajat {
		zilePeData := map[string]PozitiePlanificareSimple{}
		for _, p := range lista {
			zilePeData[p.Data] = p
		}

		// Verificare ture consecutive
		for _, p := range lista {
			d, err := time.Parse(time.DateOnly, p.Data)
			if err != nil {
				continue
			}
			ziUrmatoare := d.AddDate(0, 0, 1).Format(time.DateOnly)
			if urm, exista := zilePeData[ziUrmatoare]; exista {
				if p.CodSchimb == "G24" {
					avertismente = append(avertismente, AvertismentEchitate{
						AngajatID: angajatID,
						Data:      p.Data,
						CodSchimb: p.CodSchimb,
						Tip:       "CONSECUTIV",
						Mesaj:     "Garda de 24h este urmată imediat de o altă tură planificată a doua zi fără repaus.",
					})
				} else {
					avertismente = append(avertismente, AvertismentEchitate{
						AngajatID: angajatID,
						Data:      p.Data,
						CodSchimb: p.CodSchimb,
						Tip:       "CONSECUTIV",
						Mesaj:     "Ture planificate în zile consecutive (" + p.CodSchimb + " și " + urm.CodSchimb + ").",
					})
				}
			}

			// Rotație Crăciun / Revelion
			if EsteRevelion(p.Data) && craciunAnTrecut[angajatID] {
				avertismente = append(avertismente, AvertismentEchitate{
					AngajatID: angajatID,
					Data:      p.Data,
					CodSchimb: p.CodSchimb,
					Tip:       "ROTATIE_CRACIUN_REVELION",
					Mesaj:     "Respectă rotația de sărbători: a fost planificat de Crăciun anul trecut, este planificat de Revelion anul acesta.",
				})
			}

			// Corelare Crăciun curent cu Paști anul viitor
			if EsteCraciun(p.Data) {
				if codPasti, exista := pastiAnViitor[angajatID]; exista && codPasti == p.CodSchimb {
					avertismente = append(avertismente, AvertismentEchitate{
						AngajatID: angajatID,
						Data:      p.Data,
						CodSchimb: p.CodSchimb,
						Tip:       "CONFLICT_CRACIUN_PASTI",
						Mesaj:     "Atenție la echitate: angajatul este planificat pe același schimb (" + p.CodSchimb + ") atât de Crăciun anul acesta, cât și de Paști anul viitor.",
					})
				}
			}
		}
	}

	return avertismente
}
