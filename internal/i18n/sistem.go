package i18n

import "strings"

// LimbaSistem întoarce limba Windows, dacă există în cataloage, altfel limba
// implicită. Se folosește înainte ca setările să poată fi citite (D008).
func LimbaSistem() Limba {
	return alegeLimba(limbiPreferate())
}

// alegeLimba ia prima limbă preferată (de forma "en-US") care are catalog.
func alegeLimba(preferate []string) Limba {
	for _, p := range preferate {
		cod, _, _ := strings.Cut(strings.ToLower(p), "-")
		if _, ok := cataloage[Limba(cod)]; ok {
			return Limba(cod)
		}
	}
	return Implicita
}
