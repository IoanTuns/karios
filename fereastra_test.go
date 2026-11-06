package main

import "testing"

func TestFereastraIncapeInZonaUtila(t *testing.T) {
	for _, c := range []struct{ l, h, dpi int }{{1920, 1040, 96}, {1366, 728, 96}, {1024, 728, 96}, {1280, 680, 144}, {1920, 1040, 192}, {800, 560, 96}, {3840, 2080, 192}} {
		l, h, ml, mh := dimensiuniFereastra(c.l, c.h, c.dpi)
		if l > 1280 || h > 800 || ml > l || mh > h || l*c.dpi/96 > c.l-64*c.dpi/96 || h*c.dpi/96 > c.h-64*c.dpi/96 {
			t.Fatalf("zona %+v: %d %d minim %d %d", c, l, h, ml, mh)
		}
	}
	l, h, _, _ := dimensiuniFereastra(1920, 1040, 96)
	if l != 1280 || h != 800 {
		t.Fatal(l, h)
	}
	l, h, ml, mh := dimensiuniFereastra(0, 0, 0)
	if l != 800 || h != 600 || ml != 640 || mh != 400 {
		t.Fatal(l, h, ml, mh)
	}
}
