package calcul

import (
	"time"
	_ "time/tzdata" // Windows nu are baza de fusuri orare folosită de Go; o includem în .exe
)

// Locatie este fusul orar în care se fac toate calculele.
var Locatie = func() *time.Location {
	loc, err := time.LoadLocation("Europe/Bucharest")
	if err != nil {
		panic(err)
	}
	return loc
}()

// MomentLocal intoarce momentul in care ceasul din Europe/Bucharest arata ziua zi (00:00,
// data calendaristica citita din zi) plus minute minute (F04-T08, D102). Minutele pot trece
// de miezul noptii (de exemplu 19:00 + 12 h = 07:00 a doua zi).
//   - Ora care se repeta toamna (03:00-04:00 apare de doua ori): se ia prima trecere.
//   - Ora care lipseste primavara (03:00-04:00 nu exista): se ia primul moment care exista
//     dupa ea, adica momentul schimbarii (04:00, ora de vara).
func MomentLocal(zi time.Time, minute int64) time.Time {
	perete := time.Date(zi.Year(), zi.Month(), zi.Day(), 0, 0, 0, 0, time.UTC).Add(time.Duration(minute) * time.Minute)
	acelasiCeas := func(t time.Time) bool {
		l := t.In(Locatie)
		return time.Date(l.Year(), l.Month(), l.Day(), l.Hour(), l.Minute(), l.Second(), l.Nanosecond(), time.UTC).Equal(perete)
	}
	// Decalajele de dinainte si de dupa o eventuala schimbare (in Romania schimbarile sunt la luni distanta).
	_, inainte := perete.Add(-24 * time.Hour).In(Locatie).Zone()
	_, dupa := perete.Add(24 * time.Hour).In(Locatie).Zone()
	var gasit *time.Time
	for _, dec := range []int{inainte, dupa} {
		c := perete.Add(-time.Duration(dec) * time.Second)
		if acelasiCeas(c) && (gasit == nil || c.Before(*gasit)) {
			c := c
			gasit = &c
		}
	}
	if gasit != nil {
		return gasit.In(Locatie)
	}
	// Ora lipseste: cu decalajul de dinainte, momentul cade dupa schimbare; inceputul zonei
	// in care cade este chiar momentul schimbarii.
	start, _ := perete.Add(-time.Duration(inainte) * time.Second).In(Locatie).ZoneBounds()
	return start
}

// LunaCuSchimbareOra spune daca in luna data ora se schimba in Europe/Bucharest (F04-T08).
func LunaCuSchimbareOra(an, luna int) bool {
	_, la1 := time.Date(an, time.Month(luna), 1, 12, 0, 0, 0, Locatie).Zone()
	_, laFinal := time.Date(an, time.Month(luna)+1, 1, 12, 0, 0, 0, Locatie).Zone()
	return la1 != laFinal
}
