package calcul

import (
	"math"
	"testing"
	"time"
)

func TestFormulaRaportDateLipsaSiOperatii(t *testing.T) {
	a, b := 154., 147.
	v := map[string]*float64{"ore": &a, "norma": &b, "lipsa": nil}
	for text, asteptat := range map[string]float64{"ore - norma": 7, "ROUND(10.0/3,2)": 3.33, "MAX(ore,norma)": 154, "ABS(-7)": 7, "MIN(ore,norma)": 147} {
		r, e := EvalueazaFormulaRaport(text, v)
		if e != nil || r == nil || math.Abs(*r-asteptat) > .00001 {
			t.Fatal(text, r, e)
		}
	}
	r, e := EvalueazaFormulaRaport("ore + lipsa", v)
	if e != nil || r != nil {
		t.Fatal("date lipsa", r, e)
	}
	for _, s := range []string{"ore/0", "camp_lipsa", "os.Exit(1)", "ore[0]", "ROUND(1,8)", "f()", "1e309"} {
		if _, e := EvalueazaFormulaRaport(s, v); e == nil {
			t.Fatal("expresie acceptata", s)
		}
	}
}
func TestGardaEligibilaFaraDublareSiIntreLuni(t *testing.T) {
	p := func(s string) time.Time {
		x, e := time.Parse("2006-01-02 15:04", s)
		if e != nil {
			t.Fatal(e)
		}
		return x
	}
	lib := map[string]bool{"2026-10-31": true}
	a, z := p("2026-10-31 07:00"), p("2026-11-01 07:00")
	oct := MinuteGardaEligibile(a, z, p("2026-10-01 00:00"), p("2026-11-01 00:00"), true, true, lib)
	nov := MinuteGardaEligibile(a, z, p("2026-11-01 00:00"), p("2026-12-01 00:00"), true, true, lib)
	if oct != 17*60 || nov != 7*60 {
		t.Fatal(oct, nov)
	}
	if v := MinuteGardaEligibile(a, z, a, z, false, false, lib); v != 0 {
		t.Fatal(v)
	}
}
