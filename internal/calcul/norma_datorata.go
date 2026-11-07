package calcul

import "time"

// AbsentaNorma foloseste incadrarea explicita de raport, nu denumirea sau codul liber al absentei.
type AbsentaNorma struct {
	Data     string
	Tip      string
	Aprobata bool
}

type NormaDatorata struct {
	ZileLucratoare int
	ZileConcediu   int
	ZileDatorate   int
	Minute         int64
	Definita       bool
}

// CalculeazaNormaDatorata calculeaza norma contractelor zilnice pe calendar luni-vineri.
// Liberele calendaristice sunt excluse o singura data; compensarile fara regula nu sunt deduse.
// Norma nil reprezinta absenta unei norme zilnice prestabilite, nu o obligatie de zero ore.
func CalculeazaNormaDatorata(deLa, panaLa string, normaZilnica *int64, libere map[string]bool, absente []AbsentaNorma) (NormaDatorata, error) {
	var rez NormaDatorata
	prima, e1 := time.Parse(time.DateOnly, deLa)
	ultima, e2 := time.Parse(time.DateOnly, panaLa)
	if e1 != nil || e2 != nil || ultima.Before(prima) || ultima.Sub(prima) > 366*24*time.Hour {
		return rez, &Eroare{Cod: "pontaj.perioada_invalida"}
	}
	if normaZilnica != nil && (*normaZilnica <= 0 || *normaZilnica > 1440) {
		return rez, &Eroare{Cod: "date.norma_invalida"}
	}
	concedii := map[string]bool{}
	for _, a := range absente {
		if !a.Aprobata {
			continue
		}
		switch a.Tip {
		case "CO", "CEV", "CFS", "CM", "CVE", "CFP": // CVE, CFP: codurile redenumite (D092)
			if _, err := time.Parse(time.DateOnly, a.Data); err != nil {
				return rez, &Eroare{Cod: "pontaj.celula_invalida"}
			}
			concedii[a.Data] = true
		}
	}
	for d := prima; !d.After(ultima); d = d.AddDate(0, 0, 1) {
		data := d.Format(time.DateOnly)
		if d.Weekday() == time.Saturday || d.Weekday() == time.Sunday || libere[data] {
			continue
		}
		rez.ZileLucratoare++
		if concedii[data] {
			rez.ZileConcediu++
		}
	}
	rez.ZileDatorate = rez.ZileLucratoare - rez.ZileConcediu
	if normaZilnica != nil {
		rez.Definita = true
		rez.Minute = int64(rez.ZileDatorate) * (*normaZilnica)
	}
	return rez, nil
}
