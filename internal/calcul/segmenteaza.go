package calcul

import "time"

const formatOra = "2006-01-02 15:04"

// Segment este o bucată dintr-o tură care cade într-o singură zi calendaristică
// și într-o singură categorie de ore.
type Segment struct {
	Data      time.Time // miezul nopții zilei calendaristice
	Start     time.Time
	Sfarsit   time.Time
	TipZi     TipZi
	Categorie string
	Minute    int
}

// Segmenteaza taie tura [start, sfarsit) la miezul nopții și la limitele
// intervalelor din reguli. Tipul zilei se ia din ziua calendaristică a fiecărui
// segment. Minutele sunt timp real, deci la schimbarea orei o tură de 12 h pe
// ceas are 11 sau 13 h. Pauza nu se scade aici.
func Segmenteaza(start, sfarsit time.Time, reguli SetReguli, cal Calendar) ([]Segment, error) {
	start, sfarsit = start.In(Locatie), sfarsit.In(Locatie)
	if !sfarsit.After(start) {
		return nil, &Eroare{"calcul.tura_inversata", map[string]string{
			"start": start.Format(formatOra), "sfarsit": sfarsit.Format(formatOra)}}
	}

	var segmente []Segment
	for t := start; t.Before(sfarsit); {
		an, luna, zi := t.Date()
		data := time.Date(an, luna, zi, 0, 0, 0, 0, Locatie)
		tip := cal.TipZi(data)
		minut := t.Hour()*60 + t.Minute()

		var limita time.Time
		categorie := ""
		for _, iv := range reguli.intervale(tip) {
			if iv.Sfarsit <= minut {
				continue
			}
			// time.Date normalizează 24:00 la miezul nopții următoare și ține cont de schimbarea orei.
			limita = time.Date(an, luna, zi, 0, iv.Sfarsit, 0, 0, Locatie)
			if limita.After(t) {
				categorie = iv.Categorie
				break
			}
		}
		if categorie == "" {
			return nil, &Eroare{"calcul.regula_lipsa", map[string]string{"tip": string(tip), "ora": t.Format(formatOra)}}
		}
		if limita.After(sfarsit) {
			limita = sfarsit
		}

		segmente = append(segmente, Segment{
			Data:      data,
			Start:     t,
			Sfarsit:   limita,
			TipZi:     tip,
			Categorie: categorie,
			Minute:    int(limita.Sub(t) / time.Minute),
		})
		t = limita
	}
	return segmente, nil
}
