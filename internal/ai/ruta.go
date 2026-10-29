package ai

// Ruta spune pe ce cale merge un mesaj al utilizatorului (F14-T27).
type Ruta string

const (
	RutaConcediu    Ruta = "concediu"    // cerere de concediu: deschide formularul precompletat (D052)
	RutaPlanificare Ruta = "planificare" // planificare ghidată, previzualizată și confirmată (F14-T22)
	RutaSinteza     Ruta = "sinteza"     // întrebare la care răspunde codul (concedii, listări, acoperire…)
	RutaModel       Ruta = "model"       // tot restul ajunge la modelul local
)

// Decizie e ruta aleasă și, pentru rutele rezolvate din cod, răspunsul gata făcut.
type Decizie struct {
	Ruta    Ruta
	Raspuns RaspunsAsistent
}

// Ruteaza e singura poartă prin care trece un mesaj: alege calea într-o ordine fixă, de la
// comenzile cele mai precise la model. Căile nu se mai contrazic (de ex. o cerere de concediu
// nu mai primește lista concediilor), iar setul de comenzi reale verifică exact această decizie.
func Ruteaza(mesaj string, info ContextPontix) Decizie {
	if r, ok := AnalizeazaCerereConcediu(mesaj, info); ok {
		return Decizie{Ruta: RutaConcediu, Raspuns: r}
	}
	if i, ok := InterpreteazaPlan(mesaj, info); ok {
		return Decizie{Ruta: RutaPlanificare, Raspuns: RaspunsAsistent{Intentie: &i}}
	}
	if s, ok := AnalizeazaInterogareSinteza(mesaj, info); ok {
		return Decizie{Ruta: RutaSinteza, Raspuns: RaspunsAsistent{Mesaj: s}}
	}
	return Decizie{Ruta: RutaModel}
}
