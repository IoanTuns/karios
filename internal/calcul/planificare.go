package calcul

import "time"

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
