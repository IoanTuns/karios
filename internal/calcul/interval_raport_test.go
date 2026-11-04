package calcul

import (
	"fmt"
	"testing"
	"time"
)

func TestIntervalRaportGrupuriLaAmbeleMargini(t *testing.T) {
	libere := map[string]bool{"2026-11-30": true, "2026-12-01": true}
	for _, c := range []struct {
		luna       int
		deLa, pana string
	}{{10, "2026-10-01", "2026-11-02"}, {11, "2026-10-31", "2026-12-02"}, {12, "2026-11-28", "2027-01-01"}, {9, "2026-09-01", "2026-10-01"}} {
		t.Run(fmt.Sprint(c.luna), func(t *testing.T) {
			a := time.Date(2026, time.Month(c.luna), 1, 0, 0, 0, 0, time.UTC)
			x, y, e := IntervalRaportZileLibere(a, a.AddDate(0, 1, 0), libere)
			if e != nil || x.Format(time.DateOnly) != c.deLa || y.Format(time.DateOnly) != c.pana {
				t.Fatal(x, y, e)
			}
		})
	}
}
func TestIntervalRaportNuSarePesteZiLucratoare(t *testing.T) {
	a := time.Date(2026, 9, 1, 0, 0, 0, 0, time.UTC)
	x, y, e := IntervalRaportZileLibere(a, a.AddDate(0, 1, 0), map[string]bool{"2026-10-02": true})
	if e != nil || !x.Equal(a) || !y.Equal(a.AddDate(0, 1, 0)) {
		t.Fatal(x, y, e)
	}
	if _, _, e := IntervalRaportZileLibere(a, a, nil); e == nil {
		t.Fatal("interval invalid acceptat")
	}
	libere := map[string]bool{}
	for d := a.AddDate(-2, 0, 0); d.Before(a.AddDate(2, 0, 0)); d = d.AddDate(0, 0, 1) {
		libere[d.Format(time.DateOnly)] = true
	}
	if _, _, e := IntervalRaportZileLibere(a, a.AddDate(0, 1, 0), libere); e == nil {
		t.Fatal("calendar excesiv acceptat")
	}
}
