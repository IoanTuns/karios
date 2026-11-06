package calcul

import "time"

// IntervalRaportZileLibere pastreaza grupurile consecutive care ating marginile lunii.
// Limita protejeaza raportarea de calendare configurate cu perioade excesiv de lungi.
func IntervalRaportZileLibere(deLa, panaLa time.Time, sarbatori map[string]bool) (time.Time, time.Time, error) {
	if !deLa.Before(panaLa) {
		return deLa, panaLa, &Eroare{Cod: "rapoarte.interval_invalid"}
	}
	liber := func(d time.Time) bool {
		return d.Weekday() == time.Saturday || d.Weekday() == time.Sunday || sarbatori[d.Format(time.DateOnly)]
	}
	n := 0
	if liber(deLa) {
		for liber(deLa.AddDate(0, 0, -1)) {
			if n >= 366 || deLa.Year() <= 1 {
				return deLa, panaLa, &Eroare{Cod: "rapoarte.interval_invalid"}
			}
			deLa = deLa.AddDate(0, 0, -1)
			n++
		}
	}
	if liber(panaLa.AddDate(0, 0, -1)) {
		for liber(panaLa) {
			if n >= 366 || panaLa.Year() > 9999 {
				return deLa, panaLa, &Eroare{Cod: "rapoarte.interval_invalid"}
			}
			panaLa = panaLa.AddDate(0, 0, 1)
			n++
		}
	}
	return deLa, panaLa, nil
}
