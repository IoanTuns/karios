package ai

import "testing"

func TestStergereaDinCod(t *testing.T) {
	info := contextConcediu()
	for _, c := range []struct {
		mesaj, tip, marca, data, deLa, panaLa string
	}{
		{"Sterge garda lui Pop Ana din 15 octombrie", "sterge_tura", "M1", "2026-10-15", "", ""},
		{"Anulează tura lui Ion Dan de pe 20", "sterge_tura", "M2", "2026-10-20", "", ""},
		{"Sterge programul lui Rus Ema intre 10 si 12 octombrie", "sterge_interval", "M3", "", "2026-10-10", "2026-10-12"},
		{"Scoate-l pe M2 din tura de pe 3 noiembrie", "", "", "", "", ""}, // verbul nu e primul cuvânt întreg
	} {
		r, ok := AnalizeazaStergere(c.mesaj, info)
		if c.tip == "" {
			if ok {
				t.Fatalf("%q: recunoscut %+v", c.mesaj, r)
			}
			continue
		}
		if !ok || len(r.Propuneri) != 1 {
			t.Fatalf("%q: %v %+v", c.mesaj, ok, r)
		}
		p := r.Propuneri[0]
		text := func(k string) string { s, _ := p.Parametri[k].(string); return s }
		if p.Tip != c.tip || text("marca") != c.marca || text("data") != c.data || text("data_inceput") != c.deLa || text("data_sfarsit") != c.panaLa || !p.NecesitaConfirmare || p.Descriere == "" {
			t.Fatalf("%q: %+v", c.mesaj, p)
		}
	}
}

func TestPolitetea(t *testing.T) {
	for _, m := range []string{"Multumesc!", "mulțumesc frumos", "Salut", "ok 👍"} {
		if _, ok := RaspunsPolitete(m); !ok {
			t.Fatalf("%q nu e recunoscut", m)
		}
	}
	for _, m := range []string{"Multumesc, acum sterge garda lui Pop Ana", "da, pune-l de noapte"} {
		if _, ok := RaspunsPolitete(m); ok {
			t.Fatalf("%q tratat ca politețe", m)
		}
	}
}
