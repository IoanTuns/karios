package calcul

import "testing"

func p64(v int64) *int64 { return &v }

// I010 (D101): norma zilnica se ia pe zile; schimbarea din 15 martie conteaza din 15 martie.
func TestNormaDatorataPeZileCuIstoric(t *testing.T) {
	// Martie 2027: 1 martie e luni; 23 de zile lucratoare, fara sarbatori.
	sch := []SchimbareContract{{DeLa: "2027-03-15", NormaMinute: p64(240)}}
	n, err := CalculeazaNormaDatorata("2027-03-01", "2027-03-31", p64(480), nil, nil, sch...)
	if err != nil {
		t.Fatal(err)
	}
	// 1-12 martie: 10 zile x 8 h; 15-31 martie: 13 zile x 4 h.
	if n.ZileDatorate != 23 || n.Minute != 10*480+13*240 {
		t.Fatalf("norma pe zile: %+v", n)
	}
	// Concediul aprobat nu are norma, nici dupa schimbare.
	n, _ = CalculeazaNormaDatorata("2027-03-01", "2027-03-31", p64(480), nil, []AbsentaNorma{{Data: "2027-03-16", Tip: "CO", Aprobata: true}}, sch...)
	if n.ZileDatorate != 22 || n.Minute != 10*480+12*240 {
		t.Fatalf("concediu dupa schimbare: %+v", n)
	}
	// Fara norma zilnica (garda, lunara), istoricul nu creeaza una.
	n, _ = CalculeazaNormaDatorata("2027-03-01", "2027-03-31", nil, nil, nil, sch...)
	if n.Definita || n.Minute != 0 {
		t.Fatalf("norma nedefinita: %+v", n)
	}
	// O schimbare doar de concediu nu schimba norma.
	n, _ = CalculeazaNormaDatorata("2027-03-01", "2027-03-31", p64(480), nil, nil, SchimbareContract{DeLa: "2027-03-10", ConcediuZileAn: p64(25)})
	if n.Minute != 23*480 {
		t.Fatalf("schimbare doar de concediu: %+v", n)
	}
	if _, err := CalculeazaNormaDatorata("2027-03-01", "2027-03-31", p64(480), nil, nil, SchimbareContract{DeLa: "2027-03-10", NormaMinute: p64(0)}); err == nil {
		t.Fatal("norma 0 din istoric acceptata")
	}
}

// I010 (D101): dreptul la concediu e proportional cu zilele calendaristice, rotunjit.
func TestDreptConcediuProportional(t *testing.T) {
	cazuri := []struct {
		nume    string
		an      int
		initial *int64
		sch     []SchimbareContract
		astept  *int64
	}{
		{"fara istoric", 2027, p64(21), nil, p64(21)},
		{"de la 1 iulie", 2027, p64(21), []SchimbareContract{{DeLa: "2027-07-01", ConcediuZileAn: p64(25)}}, p64(23)}, // 21*181/365+25*184/365 = 23,0
		{"schimbare in anul anterior", 2027, p64(21), []SchimbareContract{{DeLa: "2026-05-01", ConcediuZileAn: p64(25)}}, p64(25)},
		{"schimbare in anul urmator", 2027, p64(21), []SchimbareContract{{DeLa: "2028-01-01", ConcediuZileAn: p64(25)}}, p64(21)},
		{"de la 1 ianuarie", 2027, p64(21), []SchimbareContract{{DeLa: "2027-01-01", ConcediuZileAn: p64(30)}}, p64(30)},
		{"doua schimbari", 2027, p64(20), []SchimbareContract{{DeLa: "2027-10-01", ConcediuZileAn: p64(30)}, {DeLa: "2027-04-01", ConcediuZileAn: p64(24)}}, nil}, // asteptarea e calculata mai jos
		{"an bisect", 2028, p64(21), []SchimbareContract{{DeLa: "2028-07-01", ConcediuZileAn: p64(25)}}, p64(23)},
		{"initial necompletat", 2027, nil, []SchimbareContract{{DeLa: "2027-07-01", ConcediuZileAn: p64(25)}}, nil},
		{"necompletat, fara istoric", 2027, nil, nil, nil},
		{"schimbare doar de norma", 2027, p64(21), []SchimbareContract{{DeLa: "2027-07-01", NormaMinute: p64(240)}}, p64(21)},
	}
	// (20*90 + 24*183 + 30*92) / 365 = (1800 + 4392 + 2760) / 365 = 8952/365 = 24,53 -> 25
	cazuri[5].astept = p64(25)
	for _, c := range cazuri {
		r := DreptConcediuAn(c.an, c.initial, c.sch)
		if (r == nil) != (c.astept == nil) || (r != nil && *r != *c.astept) {
			t.Errorf("%s: primit %v, astept %v", c.nume, deref(r), deref(c.astept))
		}
	}
}

func deref(p *int64) any {
	if p == nil {
		return nil
	}
	return *p
}

func TestNormaLaData(t *testing.T) {
	sch := []SchimbareContract{{DeLa: "2027-06-01", NormaMinute: p64(360)}, {DeLa: "2027-03-01", NormaMinute: p64(240)}}
	for data, astept := range map[string]int64{"2027-02-28": 480, "2027-03-01": 240, "2027-05-31": 240, "2027-06-01": 360, "2030-01-01": 360} {
		if n := NormaLaData(p64(480), sch, data); n == nil || *n != astept {
			t.Errorf("%s: %v, astept %d", data, deref(n), astept)
		}
	}
}
