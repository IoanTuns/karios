package ai

import (
	"strings"
	"testing"
)

func contextConcediu() ContextPontix {
	info := contextGrup()
	info.Angajati = append(info.Angajati, InfoAngajat{Marca: "DEMO-001", Nume: "Marin", Prenume: "Daria"})
	for i := range info.Angajati {
		info.Angajati[i].ID = "id-" + info.Angajati[i].Marca
	}
	return info
}

func TestCerereConcediu(t *testing.T) {
	cazuri := []struct {
		mesaj, marca, cod, deLa, panaLa string
	}{
		{"Maria Daria merge in concediu in octombrie 26 pana pe 31", "DEMO-001", "CO", "2026-10-26", "2026-10-31"},
		{"Pop Ana are concediu medical din 28 octombrie pana pe 3 noiembrie", "M1", "CM", "2026-10-28", "2026-11-03"},
		{"DEMO-001 pleaca in concediu de odihna pe 12", "DEMO-001", "CO", "2026-10-12", "2026-10-12"},
	}
	for _, c := range cazuri {
		r, ok := AnalizeazaCerereConcediu(c.mesaj, contextConcediu())
		if !ok || len(r.Propuneri) != 1 {
			t.Errorf("%q: %+v", c.mesaj, r)
			continue
		}
		p := r.Propuneri[0]
		if p.Tip != TipDeschideConcediu || p.Parametri["marca"] != c.marca || p.Parametri["cod"] != c.cod ||
			p.Parametri["data_inceput"] != c.deLa || p.Parametri["data_sfarsit"] != c.panaLa {
			t.Errorf("%q: %+v", c.mesaj, p.Parametri)
		}
	}
}

func TestCerereConcediuIncompletaSauIntrebare(t *testing.T) {
	info := contextConcediu()
	// Întrebările rămân la lista concediilor.
	for _, m := range []string{"Cine este in concediu in octombrie?", "Ce concedii are Pop Ana"} {
		if _, ok := AnalizeazaCerereConcediu(m, info); ok {
			t.Errorf("întrebare tratată ca cerere: %q", m)
		}
	}
	if _, ok := AnalizeazaInterogareSinteza("Cine este in concediu luna asta?", info); !ok {
		t.Error("întrebarea despre concedii nu mai e recunoscută")
	}
	if _, ok := AnalizeazaInterogareSinteza("Maria Daria merge in concediu in octombrie 26 pana pe 31", info); ok {
		t.Error("cererea de concediu a ajuns la lista concediilor")
	}
	// Comenzile incomplete primesc o întrebare de clarificare, fără propunere.
	for _, c := range []struct{ mesaj, contine string }{
		{"Ana merge in concediu pe 5", "Pentru ce angajat"},
		{"Rus Ema merge in concediu", "Pentru ce zile"},
	} {
		r, ok := AnalizeazaCerereConcediu(c.mesaj, info)
		if !ok || len(r.Propuneri) != 0 || !strings.Contains(r.Mesaj, c.contine) {
			t.Errorf("%q: %+v", c.mesaj, r)
		}
	}
}
