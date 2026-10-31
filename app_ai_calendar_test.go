package main

import (
	"pontix/internal/calcul"
	"reflect"
	"testing"
	"time"
)

func TestZileNormaleExcludSarbatoriInclusivIntreAni(t *testing.T) {
	start, _ := time.Parse(time.DateOnly, "2026-12-31")
	stop, _ := time.Parse(time.DateOnly, "2027-01-04")
	calendar := calcul.Calendar{Sarbatori: map[string]string{"2027-01-01": "Sarbatoare test", "2027-01-02": "Sarbatoare in weekend"}}
	for _, c := range []struct {
		filtru string
		doar   bool
		zile   []string
	}{
		{"lucratoare", false, []string{"2026-12-31", "2027-01-04"}},
		{"", true, []string{"2026-12-31", "2027-01-04"}},
		{"weekend", false, []string{"2027-01-02", "2027-01-03"}},
		{"tot_al_doilea_weekend", false, []string{"2027-01-02", "2027-01-03"}},
		{"toate", false, []string{"2026-12-31", "2027-01-01", "2027-01-02", "2027-01-03", "2027-01-04"}},
	} {
		if got := zileDinInterval(start, stop, c.filtru, c.doar, calendar); !reflect.DeepEqual(got, c.zile) {
			t.Fatal(c.filtru, got, c.zile)
		}
	}
}
