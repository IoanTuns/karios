package datebaza

import (
	"context"
	"encoding/json"
	"strconv"
	"strings"

	"pontix/internal/bd/acces"
	"pontix/internal/i18n"
)

type CategorieProgram struct {
	ID               string    `json:"id"`
	Denumire         i18n.Text `json:"denumire"`
	SchimburiID      []string  `json:"schimburi_id"`
	SchimbImplicitID string    `json:"schimb_implicit_id"`
	Activ            bool      `json:"activ"`
	Revizie          string    `json:"revizie"`
}

func citesteCategoriiProgram(ctx context.Context, q *acces.Queries) ([]CategorieProgram, error) {
	randuri, err := q.ListaCategoriiProgram(ctx)
	if err != nil {
		return nil, err
	}
	categorii := []CategorieProgram{}
	for _, r := range randuri {
		var denumire i18n.Text
		if err := json.Unmarshal([]byte(r.Denumire), &denumire); err != nil {
			return nil, err
		}
		ids, err := q.ListaSchimburiCategorieProgram(ctx, r.ID)
		if err != nil {
			return nil, err
		}
		categorii = append(categorii, CategorieProgram{ID: r.ID, Denumire: denumire, SchimburiID: ids,
			SchimbImplicitID: textSauGol(r.SchimbImplicitID), Activ: r.Activ == 1, Revizie: strconv.FormatInt(r.Revizie, 10)})
	}
	return categorii, nil
}

func (s *Serviciu) SalveazaCategorieProgram(ctx context.Context, c CategorieProgram) error {
	denumire := i18n.Text{}
	for limba, text := range c.Denumire {
		if text = strings.TrimSpace(text); text != "" {
			denumire[limba] = text
		}
	}
	if len(denumire) == 0 || len(c.SchimburiID) == 0 || len(c.SchimburiID) > 200 {
		return eroare("program.categorie_invalida")
	}
	c.Denumire = denumire
	return s.modifica(ctx, func(q *acces.Queries) error {
		var veche any
		if c.ID != "" {
			r, err := q.CitesteCategorieProgram(ctx, c.ID)
			if err != nil {
				return err
			}
			if c.Revizie != strconv.FormatInt(r.Revizie, 10) {
				return eroare("program.conflict")
			}
			ids, err := q.ListaSchimburiCategorieProgram(ctx, c.ID)
			if err != nil {
				return err
			}
			veche = struct {
				Categorie acces.CategoriiProgram
				Schimburi []string
			}{r, ids}
		} else if c.Revizie != "" {
			return eroare("program.conflict")
		}
		vazute := map[string]bool{}
		for _, id := range c.SchimburiID {
			if id == "" || vazute[id] {
				return eroare("program.categorie_invalida")
			}
			vazute[id] = true
			if _, err := q.CitesteSchimb(ctx, id); err != nil {
				return eroare("program.schimb_invalid")
			}
		}
		if c.SchimbImplicitID != "" && !vazute[c.SchimbImplicitID] {
			return eroare("program.implicit_invalid")
		}
		date, err := json.Marshal(c.Denumire)
		if err != nil {
			return err
		}
		if c.ID == "" {
			c.ID = idNou()
			err = q.AdaugaCategorieProgram(ctx, acces.AdaugaCategorieProgramParams{ID: c.ID, Denumire: string(date), SchimbImplicitID: textOptional(c.SchimbImplicitID), Activ: numarBool(c.Activ), CreatLa: moment(), ModificatLa: moment()})
			c.Revizie = "1"
		} else {
			revizie, _ := strconv.ParseInt(c.Revizie, 10, 64)
			n, e := q.ModificaCategorieProgram(ctx, acces.ModificaCategorieProgramParams{ID: c.ID, Denumire: string(date), SchimbImplicitID: textOptional(c.SchimbImplicitID), Activ: numarBool(c.Activ), ModificatLa: moment(), Revizie: revizie})
			err = e
			if err == nil && n != 1 {
				return eroare("program.conflict")
			}
			c.Revizie = strconv.FormatInt(revizie+1, 10)
		}
		if err != nil {
			return err
		}
		if err := q.GolesteSchimburiCategorieProgram(ctx, c.ID); err != nil {
			return err
		}
		for _, id := range c.SchimburiID {
			if err := q.AdaugaSchimbCategorieProgram(ctx, acces.AdaugaSchimbCategorieProgramParams{CategorieID: c.ID, SchimbID: id}); err != nil {
				return err
			}
		}
		return jurnal(ctx, q, "categorii_program", c.ID, veche, c)
	})
}
