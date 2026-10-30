package ai

import (
	"fmt"
	"strings"
)

// Ștergerile simple și mesajele de politețe se rezolvă din cod (F14-T30): pe CPU, modelul
// avea nevoie de 20–30 s pentru o ștergere și de peste 2 minute pentru un „Mulțumesc!”.
var (
	verbeStergere = []string{"sterge", "stergeti", "anuleaza", "anulati", "elimina", "scoate", "scoateti"}
	// Ce se poate șterge: ture și programe. Concediile au formularul lor.
	tinteStergere = []string{"tur", "gard", "garz", "program", "schimb", "pontaj", "serviciu", "noapte"}

	multumiri = map[string]bool{"multumesc": true, "multumim": true, "mersi": true, "ms": true, "mulțumesc": true,
		"thanks": true, "thank you": true, "multumesc frumos": true, "mersi frumos": true, "ok multumesc": true, "ok mersi": true}
	saluturi = map[string]bool{"salut": true, "buna": true, "buna ziua": true, "hello": true, "hi": true, "hey": true, "neata": true, "servus": true}
	acorduri = map[string]bool{"ok": true, "bine": true, "super": true, "perfect": true, "da": true, "am inteles": true, "in regula": true}
)

// RaspunsPolitete răspunde la un mesaj care e doar o mulțumire, un salut sau o confirmare.
func RaspunsPolitete(mesaj string) (string, bool) {
	m := strings.Trim(normalizeaza(mesaj), " !.?,;:🙂😊👍")
	switch {
	case multumiri[m]:
		return "Cu plăcere! Spune-mi dacă mai e ceva de planificat sau de verificat.", true
	case saluturi[m]:
		return "Bună! Pot planifica ture, șterge programe, deschide cereri de concediu și răspunde la întrebări despre pontaj. Scrie „ajutor” pentru exemple.", true
	case acorduri[m]:
		return "Bine. Spune-mi ce vrei să fac mai departe.", true
	}
	return "", false
}

// esteComandaStergere: un verb de ștergere la începutul mesajului și o tură sau un program ca țintă.
func esteComandaStergere(mesaj, m string) bool {
	if esteIntrebare(mesaj, m) || esteToti(m) || contineOricare(m, cuvinteConcediu) {
		return false
	}
	cuvinte := strings.Fields(m)
	if len(cuvinte) == 0 {
		return false
	}
	verb, tinta := false, false
	for _, v := range verbeStergere {
		verb = verb || cuvinte[0] == v
	}
	// Țintele sunt rădăcini („gard” → garda, gărzile), deci se potrivesc la începutul cuvântului.
	for _, c := range cuvinte[1:] {
		for _, r := range tinteStergere {
			tinta = tinta || strings.HasPrefix(c, r)
		}
	}
	return verb && tinta
}

// AnalizeazaStergere recunoaște „Șterge garda lui Pop Ana din 15” sau „Anulează programul lui
// Rus Ema între 10 și 12 octombrie” și propune ștergerea, pe care utilizatorul o confirmă, ca pentru
// propunerile modelului. Fără persoană sau fără zi clară, mesajul merge mai departe la model.
func AnalizeazaStergere(mesaj string, info ContextPontix) (RaspunsAsistent, bool) {
	m := normalizeaza(mesaj)
	if !esteComandaStergere(mesaj, m) {
		return RaspunsAsistent{}, false
	}
	a, mai := gasesteAngajatInMesaj(mesaj, m, info.Angajati)
	if len(mai) > 1 {
		var nume []string
		for _, x := range mai {
			nume = append(nume, fmt.Sprintf("%s %s (%s)", x.Nume, x.Prenume, x.Marca))
		}
		return RaspunsAsistent{Mesaj: "Am găsit mai multe persoane potrivite: " + strings.Join(nume, ", ") + ". Scrie numele complet sau marca."}, true
	}
	if a.ID == "" {
		return RaspunsAsistent{}, false
	}
	deLa, panaLa, ok := intervalDin(m, info)
	if !ok {
		return RaspunsAsistent{}, false
	}
	tip, args := "sterge_tura", map[string]any{"marca": a.Marca, "data": deLa.Format("2006-01-02")}
	if !panaLa.Equal(deLa) {
		tip, args = "sterge_interval", map[string]any{"marca": a.Marca, "data_inceput": deLa.Format("2006-01-02"), "data_sfarsit": panaLa.Format("2006-01-02")}
	}
	return RaspunsAsistent{Propuneri: []PropunereActiune{{
		ID:                 "cod-" + tip + "-" + a.Marca + "-" + deLa.Format("20060102"),
		Tip:                tip,
		Descriere:          construiesteDescrierePropunere(tip, args, info),
		Parametri:          args,
		NecesitaConfirmare: true,
	}}}, true
}
