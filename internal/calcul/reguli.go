package calcul

import "time"

type TipZi string

const (
	Lucratoare TipZi = "LUCRATOARE"
	Sambata    TipZi = "SAMBATA"
	Duminica   TipZi = "DUMINICA"
	Sarbatoare TipZi = "SARBATOARE"
)

var TipuriZi = []TipZi{Lucratoare, Sambata, Duminica, Sarbatoare}

// Calendar stabilește tipul fiecărei zile calendaristice.
type Calendar struct {
	Sarbatori map[string]string // "2006-01-02" → denumire
}

// TipZi întoarce tipul zilei. Sărbătoarea are prioritate față de weekend.
func (c Calendar) TipZi(zi time.Time) TipZi {
	if _, ok := c.Sarbatori[zi.Format(time.DateOnly)]; ok {
		return Sarbatoare
	}
	switch zi.Weekday() {
	case time.Saturday:
		return Sambata
	case time.Sunday:
		return Duminica
	}
	return Lucratoare
}
