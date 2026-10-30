package ai

import (
	"strings"

	"pontix/internal/ajutor"
	"pontix/internal/i18n"
)

// AjutorUtilizare identifica intrebari de utilizare, nu cereri de modificare.
func AjutorUtilizare(mesaj, limba string) (RaspunsAsistent, bool) {
	m := strings.TrimSpace(ajutor.Normalizat(mesaj))
	intrebari := []string{"cum ", "cum?", "cum sa ", "cum pot ", "cum fac ", "unde gasesc", "unde pot", "unde este", "how ", "where can", "where do", "how-to", "ajutor ", "ghid ", "manual "}
	este := m == "manual" || m == "ghid" || m == "cum" || m == "manual interactiv"
	for _, prefix := range intrebari {
		este = este || strings.HasPrefix(m, prefix)
	}
	if !este {
		return RaspunsAsistent{}, false
	}
	l := i18n.Limba(limba)
	articole := ajutor.Cauta(mesaj, l)
	if m == "manual" || m == "ghid" || m == "manual interactiv" {
		articole = nil
		for _, id := range []string{"ai.foloseste", "pontaj.completeaza", "absente.inregistreaza"} {
			if articol, ok := ajutor.Gaseste(id, l); ok {
				articole = append(articole, articol)
			}
		}
	}
	r := RaspunsAsistent{ArticoleAjutor: []string{}}
	if len(articole) == 0 {
		r.Mesaj = i18n.T(l, "ajutor.fara_acoperire", nil)
		return r, true
	}
	r.VersiuneAjutor = articole[0].Versiune
	texte := []string{}
	for _, a := range articole {
		r.ArticoleAjutor = append(r.ArticoleAjutor, a.ID)
		pasi := []string{a.Titlu}
		for i, p := range a.Pasi {
			pasi = append(pasi, []string{"1. ", "2. ", "3. "}[i]+p)
		}
		pasi = append(pasi, i18n.T(l, "ajutor.sursa", map[string]string{"id": a.ID, "versiune": a.Versiune}))
		texte = append(texte, strings.Join(pasi, "\n"))
	}
	r.Mesaj = strings.Join(texte, "\n\n")
	return r, true
}
