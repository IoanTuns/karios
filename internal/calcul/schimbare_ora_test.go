package calcul

import (
	"testing"
	"time"
)

// F04-T08 (D102): turele din program in noptile cu schimbarea orei.
func TestSchimbareOraTureDeNoapte(t *testing.T) {
	p := []ProgramPontaj{{Cod: "N12", OraStart: 19 * 60, DurataMinute: 720}}
	cazuri := []struct {
		nume, data, luna string
		mod              string
		atribuire        string
		total            string
		celule           []int64
	}{
		// Ultima duminica din octombrie 2026: 25 octombrie, 04:00 devine 03:00.
		{"toamna nominal", "2026-10-24", "2026-10", SchimbareOraNominala, AtribuireInceput, "12", []int64{720}},
		{"toamna efectiv", "2026-10-24", "2026-10", SchimbareOraEfectiva, AtribuireInceput, "13", []int64{780}},
		{"toamna efectiv pe zile", "2026-10-24", "2026-10", SchimbareOraEfectiva, AtribuireCalendar, "13", []int64{300, 480}},
		{"toamna gol = nominal", "2026-10-24", "2026-10", "", AtribuireInceput, "12", []int64{720}},
		// Ultima duminica din martie 2027: 28 martie, 03:00 devine 04:00.
		{"primavara nominal", "2027-03-27", "2027-03", SchimbareOraNominala, AtribuireCalendar, "12", []int64{300, 420}},
		{"primavara efectiv", "2027-03-27", "2027-03", SchimbareOraEfectiva, AtribuireInceput, "11", []int64{660}},
		{"primavara efectiv pe zile", "2027-03-27", "2027-03", SchimbareOraEfectiva, AtribuireCalendar, "11", []int64{300, 360}},
		// O noapte obisnuita nu se schimba.
		{"noapte obisnuita efectiv", "2026-10-10", "2026-10", SchimbareOraEfectiva, AtribuireCalendar, "12", []int64{300, 420}},
	}
	for _, caz := range cazuri {
		t.Run(caz.nume, func(t *testing.T) {
			c := ConfiguratiePontajInitiala()
			c.SchimbareOra, c.Atribuire = caz.mod, caz.atribuire
			d, _ := time.Parse("2006-01", caz.luna)
			r, err := CalculeazaPontaj(caz.luna+"-01", d.AddDate(0, 1, -1).Format(time.DateOnly), c, p, []IntrarePontaj{{AngajatID: "a", Data: caz.data, Valoare: "N12"}}, nil)
			if err != nil {
				t.Fatal(err)
			}
			if v := TextValoarePontaj(r.Totaluri[0].Valori["T"]); v != caz.total {
				t.Errorf("T=%s; asteptat %s", v, caz.total)
			}
			if len(r.Celule) != len(caz.celule) {
				t.Fatalf("celule %+v; asteptat %v", r.Celule, caz.celule)
			}
			for i, m := range caz.celule {
				if r.Celule[i].Minute != m {
					t.Errorf("celula %s: %d; asteptat %d", r.Celule[i].Data, r.Celule[i].Minute, m)
				}
			}
		})
	}
}

// ZNL se imparte dupa data locala: tura de sambata spre duminica ramane in weekend.
func TestSchimbareOraZNL(t *testing.T) {
	c := ConfiguratiePontajInitiala()
	c.SchimbareOra, c.Atribuire = SchimbareOraEfectiva, AtribuireCalendar
	r, err := CalculeazaPontaj("2026-10-01", "2026-10-31", c, []ProgramPontaj{{Cod: "N12", OraStart: 19 * 60, DurataMinute: 720}}, []IntrarePontaj{{AngajatID: "a", Data: "2026-10-24", Valoare: "N12"}}, nil)
	if err != nil {
		t.Fatal(err)
	}
	if v := TextValoarePontaj(r.Totaluri[0].Valori["ZNL"]); v != "13" {
		t.Fatalf("ZNL=%s; asteptat 13", v)
	}
}

func TestMomentLocalOreLipsaSiRepetate(t *testing.T) {
	zi := func(s string) time.Time { d, _ := time.Parse(time.DateOnly, s); return d }
	cazuri := []struct {
		zi     string
		minute int64
		utc    string
	}{
		{"2026-10-25", 3*60 + 30, "2026-10-25T00:30:00Z"}, // ora repetata: prima trecere (ora de vara, +3)
		{"2026-10-25", 4 * 60, "2026-10-25T02:00:00Z"},    // dupa schimbare (+2)
		{"2027-03-28", 3*60 + 30, "2027-03-28T01:00:00Z"}, // ora lipsa: momentul schimbarii, 04:00 ora de vara
		{"2027-03-28", 2*60 + 59, "2027-03-28T00:59:00Z"}, // inainte de schimbare (+2)
		{"2026-10-24", 19*60 + 720, "2026-10-25T05:00:00Z"},
	}
	for _, caz := range cazuri {
		if m := MomentLocal(zi(caz.zi), caz.minute).UTC().Format(time.RFC3339); m != caz.utc {
			t.Errorf("%s +%d: %s; asteptat %s", caz.zi, caz.minute, m, caz.utc)
		}
	}
}

func TestSchimbareOraValidare(t *testing.T) {
	c := ConfiguratiePontajInitiala()
	if c.SchimbareOra != SchimbareOraNominala {
		t.Fatal(c.SchimbareOra)
	}
	c.SchimbareOra = "ALTA"
	if ValideazaConfiguratiePontaj(c) == nil {
		t.Fatal("valoare necunoscuta acceptata")
	}
}

func TestLunaCuSchimbareOra(t *testing.T) {
	for luna := 1; luna <= 12; luna++ {
		if obtinut, asteptat := LunaCuSchimbareOra(2026, luna), luna == 3 || luna == 10; obtinut != asteptat {
			t.Errorf("luna %d: %v", luna, obtinut)
		}
	}
}
