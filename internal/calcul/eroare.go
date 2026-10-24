package calcul

import "pontix/internal/i18n"

// Eroare are un cod stabil (cheia din cataloagele i18n) și parametrii lui,
// ca interfața să o poată afișa în limba utilizatorului.
type Eroare struct {
	Cod   string
	Param map[string]string
}

// Error întoarce textul în limba implicită, pentru jurnale.
func (e *Eroare) Error() string { return e.Text(i18n.Implicita) }

func (e *Eroare) Text(l i18n.Limba) string { return i18n.T(l, e.Cod, e.Param) }
