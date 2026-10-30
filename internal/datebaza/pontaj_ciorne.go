package datebaza

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"

	"pontix/internal/bd/acces"
	"pontix/internal/calcul"
)

type CelulaSchita struct {
	AngajatID string `json:"angajat_id"`
	Data      string `json:"data"`
	Valoare   string `json:"valoare"`
	Revizie   string `json:"revizie"`
}

type SchitaLuna struct {
	An           int                       `json:"an"`
	Luna         int                       `json:"luna"`
	Stare        string                    `json:"stare"`
	Celule       []CelulaSchita            `json:"celule"`
	Absente      []ZiAbsenta               `json:"absente"`
	Calcul       calcul.RezultatPontaj     `json:"calcul"`
	Configuratie calcul.ConfiguratiePontaj `json:"configuratie"`
}

func limiteLuna(an, luna int) (string, string, error) {
	if an < 1 || an > 9999 || luna < 1 || luna > 12 {
		return "", "", eroare("pontaj.perioada_invalida")
	}
	return fmt.Sprintf("%04d-%02d-01", an, luna), time.Date(an, time.Month(luna)+1, 0, 0, 0, 0, 0, time.UTC).Format(time.DateOnly), nil
}

func (s *Serviciu) CitesteSchita(ctx context.Context, an, luna int, previzualizare ...[]CelulaSchita) (SchitaLuna, error) {
	deLa, panaLa, err := limiteLuna(an, luna)
	if err != nil {
		return SchitaLuna{}, err
	}
	tx, err := s.db.BeginTx(ctx, &sql.TxOptions{ReadOnly: true})
	if err != nil {
		return SchitaLuna{}, traduceEroare(err)
	}
	defer tx.Rollback()
	q := acces.New(tx)
	randuri, err := q.ListaCiorneLuna(ctx, acces.ListaCiorneLunaParams{DataDeLa: deLa, DataPanaLa: panaLa})
	if err != nil {
		return SchitaLuna{}, traduceEroare(err)
	}
	rez := SchitaLuna{An: an, Luna: luna, Celule: []CelulaSchita{}}
	st, _ := s.CitesteStareLuna(ctx, an, luna)
	rez.Stare = st.Stare
	for _, r := range randuri {
		rez.Celule = append(rez.Celule, CelulaSchita{AngajatID: r.AngajatID, Data: r.Data, Valoare: r.Valoare, Revizie: strconv.FormatInt(r.Revizie, 10)})
	}
	rez.Absente, err = citesteAbsenteInterval(ctx, q, deLa, panaLa)
	if err != nil {
		return SchitaLuna{}, traduceEroare(err)
	}
	var modificari []CelulaSchita
	if len(previzualizare) > 0 {
		modificari = previzualizare[0]
	}
	rez.Calcul, rez.Configuratie, err = CalculeazaPontaj(ctx, q, an, luna, modificari)
	if err != nil {
		return SchitaLuna{}, traduceEroare(err)
	}
	if err := tx.Commit(); err != nil {
		return SchitaLuna{}, traduceEroare(err)
	}
	return rez, nil
}

func (s *Serviciu) SalveazaSchita(ctx context.Context, an, luna int, celule []CelulaSchita) error {
	return s.salveazaSchita(ctx, an, luna, celule, false)
}

// SalveazaSchitaCuPlanificare sincronizează programele speciale la salvarea din grilă.
func (s *Serviciu) SalveazaSchitaCuPlanificare(ctx context.Context, an, luna int, celule []CelulaSchita) error {
	return s.salveazaSchita(ctx, an, luna, celule, true)
}

func (s *Serviciu) salveazaSchita(ctx context.Context, an, luna int, celule []CelulaSchita, cuPlanificare bool) error {
	inchisa, err := s.EsteLunaInchisa(ctx, an, luna)
	if err != nil {
		return err
	}
	if inchisa {
		return eroare("pontaj.luna_inchisa")
	}

	deLa, panaLa, err := limiteLuna(an, luna)
	if err != nil {
		return err
	}
	if len(celule) > 10000 {
		return eroare("pontaj.lot_invalid")
	}
	vazute := map[string]bool{}
	for _, c := range celule {
		if !dataValida(c.Data) || c.Data < deLa || c.Data > panaLa || c.AngajatID == "" || utf8.RuneCountInString(c.Valoare) > 16 || strings.ContainsAny(c.Valoare, "\x00\t\r\n") {
			return eroare("pontaj.celula_invalida")
		}
		cheie := c.AngajatID + "/" + c.Data
		if vazute[cheie] {
			return eroare("pontaj.lot_invalid")
		}
		vazute[cheie] = true
	}
	return s.modificaTx(ctx, func(tx *sql.Tx, q *acces.Queries) error {
		if cuPlanificare {
			stare := "DESCHISA"
			err := tx.QueryRowContext(ctx, `SELECT stare FROM stare_luni WHERE an=? AND luna=?`, an, luna).Scan(&stare)
			if err != nil && err != sql.ErrNoRows {
				return err
			}
			if stare != "DESCHISA" {
				return eroare("pontaj.luna_inchisa")
			}
		}
		modificariPlan := []schimbarePontajPlan{}
		angajati := map[string]bool{}
		categorii := map[string]string{}
		schimburi, err := q.ListaSchimburi(ctx)
		if err != nil {
			return err
		}
		programe := map[string]string{}
		for _, p := range schimburi {
			programe[strings.ToUpper(p.Cod)] = p.ID
		}
		// Absențele și valorile existente ale lunii se citesc o singură dată, nu pe celulă (I066):
		// o lună întreagă are mii de celule.
		absenteActive := map[string]bool{}
		absente, err := q.ListaAbsenteInterval(ctx, acces.ListaAbsenteIntervalParams{DeLa: deLa, PanaLa: panaLa})
		if err != nil {
			return err
		}
		for _, z := range absente {
			if z.Activ == 1 {
				absenteActive[z.AngajatID+"/"+z.Data] = true
			}
		}
		existente := map[string]acces.CitesteCiornaRow{}
		ciorne, err := q.ListaCiorneLuna(ctx, acces.ListaCiorneLunaParams{DataDeLa: deLa, DataPanaLa: panaLa})
		if err != nil {
			return err
		}
		for _, x := range ciorne {
			existente[x.AngajatID+"/"+x.Data] = acces.CitesteCiornaRow{Valoare: x.Valoare, ModificatLa: x.ModificatLa, Revizie: x.Revizie}
		}
		// Programele permise ale unei categorii se citesc o dată pe categorie.
		type permisiuni struct {
			activa   bool
			programe map[string]bool
		}
		categoriiCitite := map[string]permisiuni{}
		for _, c := range celule {
			if absenteActive[c.AngajatID+"/"+c.Data] {
				return eroare("absente.celula_protejata")
			}
			if !angajati[c.AngajatID] {
				a, err := q.CitesteAngajat(ctx, c.AngajatID)
				if err != nil {
					return err
				}
				categorii[c.AngajatID] = textSauGol(a.CategorieProgramID)
				angajati[c.AngajatID] = true
			}
			veche, exista := existente[c.AngajatID+"/"+c.Data]
			if (exista && c.Revizie != strconv.FormatInt(veche.Revizie, 10)) || (!exista && c.Revizie != "") {
				return eroare("pontaj.conflict")
			}
			c.Valoare = strings.TrimSpace(c.Valoare)
			if cuPlanificare {
				modificariPlan = append(modificariPlan, schimbarePontajPlan{Celula: c, Inainte: veche.Valoare})
			}
			if (exista && c.Valoare == veche.Valoare) || (!exista && c.Valoare == "") {
				continue
			}
			if categorieID := categorii[c.AngajatID]; categorieID != "" {
				perm, citita := categoriiCitite[categorieID]
				if !citita {
					categorie, err := q.CitesteCategorieProgram(ctx, categorieID)
					if err != nil {
						return err
					}
					ids, err := q.ListaSchimburiCategorieProgram(ctx, categorieID)
					if err != nil {
						return err
					}
					perm = permisiuni{activa: categorie.Activ == 1, programe: map[string]bool{}}
					for _, id := range ids {
						perm.programe[id] = true
					}
					categoriiCitite[categorieID] = perm
				}
				if programID, programCunoscut := programe[strings.ToUpper(c.Valoare)]; programCunoscut && perm.activa && !perm.programe[programID] {
					return eroare("program.program_nepermis")
				}
			}
			revizie, err := q.SalveazaCiorna(ctx, acces.SalveazaCiornaParams{AngajatID: c.AngajatID, Data: c.Data, Valoare: c.Valoare, ModificatLa: moment()})
			if err != nil {
				return err
			}
			c.Revizie = strconv.FormatInt(revizie, 10)
			var inainte any
			if exista {
				inainte = veche
			}
			if err := jurnal(ctx, q, "pontaj_ciorne", c.AngajatID+"/"+c.Data, inainte, c); err != nil {
				return err
			}
		}
		if cuPlanificare {
			return sincronizeazaPontajPlan(ctx, q, deLa, panaLa, modificariPlan)
		}
		return nil
	})
}
