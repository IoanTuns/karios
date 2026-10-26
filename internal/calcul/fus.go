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
