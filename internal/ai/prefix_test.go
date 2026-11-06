package ai

import (
	"strings"
	"testing"
)

// contextPrefix: Pop Ana are prefixul Dr. (D089, F14-T33).
func contextPrefix() ContextPontix {
	info := contextConcediu()
	for i := range info.Angajati {
		if info.Angajati[i].Marca == "M1" {
			info.Angajati[i].Prefix = "Dr."
		}
	}
	return info
}

// Căutarea găsește același angajat cu sau fără prefix, cu sau fără punct.
func TestCautareNumeCuSauFaraPrefix(t *testing.T) {
	info := contextPrefix()
	for _, text := range []string{"Pop Ana", "Dr. Pop Ana", "dr Pop Ana", "DR. Ana", "Ana", "Dr. Ana Pop"} {
		a, ok := gasesteAngajatInfo(text, info.Angajati)
		if !ok || a.Marca != "M1" {
			t.Errorf("%q: %+v %v", text, a, ok)
		}
	}
	// Prefixul singur nu alege pe nimeni.
	if a, ok := gasesteAngajatInfo("Dr.", info.Angajati); ok {
		t.Errorf("prefixul singur a ales %+v", a)
	}
}

func TestConcediuCuPrefixInMesaj(t *testing.T) {
	for _, m := range []string{"Dr. Pop Ana merge in concediu din 26 pana pe 31 octombrie", "Pop Ana merge in concediu din 26 pana pe 31 octombrie"} {
		r, ok := AnalizeazaCerereConcediu(m, contextPrefix())
		if !ok || len(r.Propuneri) != 1 || r.Propuneri[0].Parametri["marca"] != "M1" {
			t.Fatalf("%q: %+v", m, r)
		}
		if !strings.HasPrefix(r.Propuneri[0].Parametri["interpretare"].(string), "Dr. Pop Ana (M1)") {
			t.Errorf("interpretarea nu afișează prefixul: %v", r.Propuneri[0].Parametri["interpretare"])
		}
	}
}

// Răspunsurile afișează prefixul, ca listele aplicației; cei fără prefix rămân neschimbați.
func TestRaspunsuriAfiseazaPrefixul(t *testing.T) {
	info := contextPrefix()
	r, ok := AnalizeazaInterogareSinteza("Listeaza toti angajati", info)
	if !ok || !strings.Contains(r, "Dr. Pop Ana (M1)") || !strings.Contains(r, "- Rus Ema (M3)") {
		t.Fatalf("listare:\n%s", r)
	}
	if s := formateazaAngajat("M1", info.Angajati); s != "Dr. Pop Ana (marca M1)" {
		t.Errorf("formateazaAngajat: %q", s)
	}
	// Contextul trimis modelului păstrează angajatul numit cu prefix.
	c := ContextPentruCerere("Ce program are Dr. Pop Ana", info)
	gasit := false
	for _, a := range c.Angajati {
		gasit = gasit || a.Marca == "M1"
	}
	if !gasit {
		t.Error("angajatul numit cu prefix lipsește din context")
	}
}
