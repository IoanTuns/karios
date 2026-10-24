package i18n

// Eroare are un cod stabil (cheia din cataloage) și parametrii lui, ca
// interfața să o poată afișa în limba utilizatorului (D004).
type Eroare struct {
	Cod   string
	Param map[string]string
}

// Error întoarce textul în limba implicită, pentru jurnale.
func (e *Eroare) Error() string { return e.Text(Implicita) }

func (e *Eroare) Text(l Limba) string { return T(l, e.Cod, e.Param) }
