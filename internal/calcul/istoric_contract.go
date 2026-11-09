package calcul

import (
	"sort"
	"time"
)

// SchimbareContract este o schimbare a normei zilnice si/sau a dreptului anual la concediu,
// valabila de la o data pana la urmatoarea schimbare (I010, D101). Un camp nil ramane neschimbat.
// Valorile din fisa angajatului sunt cele de dinaintea primei schimbari.
type SchimbareContract struct {
	DeLa           string
	NormaMinute    *int64
	ConcediuZileAn *int64
}

func sorteazaSchimbari(s []SchimbareContract) []SchimbareContract {
	r := append([]SchimbareContract(nil), s...)
	sort.SliceStable(r, func(i, j int) bool { return r[i].DeLa < r[j].DeLa })
	return r
}

// NormaLaData intoarce norma zilnica valabila in ziua data. O norma nedefinita (nil) ramane
// nedefinita: contractele fara norma zilnica (garda, lunara, fara) nu primesc norma din istoric.
func NormaLaData(initiala *int64, schimbari []SchimbareContract, data string) *int64 {
	if initiala == nil {
		return nil
	}
	rez := initiala
	for _, s := range sorteazaSchimbari(schimbari) {
		if s.DeLa > data {
			break
		}
		if s.NormaMinute != nil {
			rez = s.NormaMinute
		}
	}
	return rez
}

func concediuLaData(initial *int64, schimbari []SchimbareContract, data string) *int64 {
	rez := initial
	for _, s := range schimbari {
		if s.DeLa > data {
			break
		}
		if s.ConcediuZileAn != nil {
			rez = s.ConcediuZileAn
		}
	}
	return rez
}

// DreptConcediuAn calculeaza dreptul la concediu pe anul an, proportional cu zilele calendaristice
// in care a fost valabila fiecare valoare, rotunjit la cel mai apropiat intreg (0,5 in sus).
// Regula a fost aleasa de Ioan Tuns (2026-11-09). Daca o parte a anului nu are valoare, rezultatul
// este nedefinit (nil): lipsa datelor nu se trateaza ca zero zile.
func DreptConcediuAn(an int, initial *int64, schimbari []SchimbareContract) *int64 {
	s := sorteazaSchimbari(schimbari)
	inceput := time.Date(an, 1, 1, 0, 0, 0, 0, time.UTC)
	sfarsit := time.Date(an+1, 1, 1, 0, 0, 0, 0, time.UTC)
	zileAn := int64(sfarsit.Sub(inceput).Hours() / 24)
	limite := []time.Time{inceput}
	for _, x := range s {
		if x.ConcediuZileAn == nil {
			continue
		}
		d, err := time.Parse(time.DateOnly, x.DeLa)
		if err == nil && d.After(inceput) && d.Before(sfarsit) {
			limite = append(limite, d)
		}
	}
	limite = append(limite, sfarsit)
	var numarator int64
	for i := 0; i+1 < len(limite); i++ {
		if !limite[i].Before(limite[i+1]) {
			continue
		}
		v := concediuLaData(initial, s, limite[i].Format(time.DateOnly))
		if v == nil {
			return nil
		}
		numarator += *v * int64(limite[i+1].Sub(limite[i]).Hours()/24)
	}
	rez := (2*numarator + zileAn) / (2 * zileAn)
	return &rez
}
