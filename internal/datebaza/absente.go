package datebaza

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"strconv"
	"strings"
	"time"
	_ "time/tzdata" // Fusurile orare rămân disponibile în executabilul portabil.

	"pontix/internal/bd/acces"
	"pontix/internal/i18n"
)

type TipAbsenta struct {
	Cod       string    `json:"cod"`
	Denumire  i18n.Text `json:"denumire"`
	Categorie string    `json:"categorie"`
	Activ     bool      `json:"activ"`
}

type ZiAbsenta struct {
	AngajatID string `json:"angajat_id"`
	Data      string `json:"data"`
	Cod       string `json:"cod"`
	Activ     bool   `json:"activ"`
	Revizie   string `json:"revizie"`
}

type AbsenteAn struct {
	An     int          `json:"an"`
	Tipuri []TipAbsenta `json:"tipuri"`
	Zile   []ZiAbsenta  `json:"zile"`
}

func (s *Serviciu) CitesteAbsente(ctx context.Context, an int) (AbsenteAn, error) {
	deLa, _, err := limiteLuna(an, 1)
	if err != nil {
		return AbsenteAn{}, err
	}
	_, panaLa, _ := limiteLuna(an, 12)
	tx, err := s.db.BeginTx(ctx, &sql.TxOptions{ReadOnly: true})
	if err != nil {
		return AbsenteAn{}, traduceEroare(err)
	}
	defer tx.Rollback()
	q := acces.New(tx)
	rez := AbsenteAn{An: an, Tipuri: []TipAbsenta{}, Zile: []ZiAbsenta{}}
	tipuri, err := q.ListaTipuriAbsenta(ctx)
	if err != nil {
		return AbsenteAn{}, traduceEroare(err)
	}
	for _, r := range tipuri {
		var denumire i18n.Text
		if err := json.Unmarshal([]byte(r.Denumire), &denumire); err != nil {
			return AbsenteAn{}, traduceEroare(err)
		}
		rez.Tipuri = append(rez.Tipuri, TipAbsenta{Cod: r.Cod, Denumire: denumire, Categorie: r.Categorie, Activ: r.Activ == 1})
	}
	rez.Zile, err = citesteAbsenteInterval(ctx, q, deLa, panaLa)
	if err != nil {
		return AbsenteAn{}, traduceEroare(err)
	}
	if err := tx.Commit(); err != nil {
		return AbsenteAn{}, traduceEroare(err)
	}
	return rez, nil
}

func citesteAbsenteInterval(ctx context.Context, q *acces.Queries, deLa, panaLa string) ([]ZiAbsenta, error) {
	randuri, err := q.ListaAbsenteInterval(ctx, acces.ListaAbsenteIntervalParams{DeLa: deLa, PanaLa: panaLa})
	if err != nil {
		return nil, err
	}
	rez := []ZiAbsenta{}
	for _, r := range randuri {
		rez = append(rez, ZiAbsenta{AngajatID: r.AngajatID, Data: r.Data, Cod: r.Cod, Activ: r.Activ == 1, Revizie: strconv.FormatInt(r.Revizie, 10)})
	}
	return rez, nil
}

func (s *Serviciu) SalveazaTipAbsenta(ctx context.Context, initial string, tip TipAbsenta) error {
	if initial == "" {
		tip.Cod = strings.ToUpper(strings.TrimSpace(tip.Cod))
	}
	if tip.Cod == "" || len(tip.Cod) > 16 || strings.IndexFunc(tip.Cod, func(r rune) bool {
		return !((r >= 'A' && r <= 'Z') || (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == '_' || r == '-')
	}) >= 0 || (initial != "" && initial != tip.Cod) {
		return eroare("absente.cod_invalid")
	}
	switch tip.Categorie {
	case "ODIHNA", "MEDICAL", "SPECIAL":
	default:
		return eroare("absente.categorie_invalida")
	}
	curata := i18n.Text{}
	for limba, text := range tip.Denumire {
		if text = strings.TrimSpace(text); text != "" {
			curata[limba] = text
		}
	}
	if len(curata) == 0 {
		return eroare("sarbatori.denumire_obligatorie")
	}
	tip.Denumire = curata
	return s.modifica(ctx, func(q *acces.Queries) error {
		vechi, err := q.CitesteTipAbsenta(ctx, tip.Cod)
		exista := err == nil
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			return err
		}
		if initial != "" && !exista {
			return eroare("date.inexistent")
		}
		if initial == "" {
			if _, err := q.CautaCodAbsenta(ctx, tip.Cod); err == nil {
				return eroare("absente.cod_duplicat")
			} else if !errors.Is(err, sql.ErrNoRows) {
				return err
			}
		}
		if _, err := q.CitesteCodSchimb(ctx, tip.Cod); err == nil {
			return eroare("absente.cod_duplicat")
		} else if !errors.Is(err, sql.ErrNoRows) {
			return err
		}
		if exista && vechi.Categorie != tip.Categorie {
			folosit, err := q.TipAbsentaFolosit(ctx, tip.Cod)
			if err != nil {
				return err
			}
			if folosit != nil && *folosit {
				return eroare("absente.categorie_folosita")
			}
		}
		text, err := json.Marshal(tip.Denumire)
		if err != nil {
			return err
		}
		acum := moment()
		if err := q.SalveazaTipAbsenta(ctx, acces.SalveazaTipAbsentaParams{Cod: tip.Cod, Denumire: string(text), Categorie: tip.Categorie, Activ: numarBool(tip.Activ), CreatLa: acum, ModificatLa: acum}); err != nil {
			return err
		}
		var inainte any
		if exista {
			inainte = vechi
		}
		return jurnal(ctx, q, "coduri_absenta", tip.Cod, inainte, tip)
	})
}

func (s *Serviciu) SalveazaAbsente(ctx context.Context, zile []ZiAbsenta) error {
	if len(zile) == 0 || len(zile) > 10000 {
		return eroare("absente.lot_invalid")
	}
	vazute := map[string]bool{}
	for _, z := range zile {
		k := z.AngajatID + "/" + z.Data
		if z.AngajatID == "" || !dataValida(z.Data) || z.Cod == "" || vazute[k] {
			return eroare("absente.lot_invalid")
		}
		vazute[k] = true
	}
	local, err := time.LoadLocation("Europe/Bucharest")
	if err != nil {
		return traduceEroare(err)
	}
	return s.modifica(ctx, func(q *acces.Queries) error {
		for _, z := range zile {
			if _, err := q.CitesteAngajat(ctx, z.AngajatID); err != nil {
				return err
			}
			veche, err := q.CitesteZiAbsenta(ctx, acces.CitesteZiAbsentaParams{AngajatID: z.AngajatID, Data: z.Data})
			exista := err == nil
			if err != nil && !errors.Is(err, sql.ErrNoRows) {
				return err
			}
			if (exista && z.Revizie != strconv.FormatInt(veche.Revizie, 10)) || (!exista && z.Revizie != "") {
				return eroare("absente.conflict")
			}
			if !z.Activ && (!exista || veche.Cod != z.Cod) {
				return eroare("absente.lot_invalid")
			}
			tip, err := q.CitesteTipAbsenta(ctx, z.Cod)
			if errors.Is(err, sql.ErrNoRows) {
				return eroare("absente.tip_invalid")
			}
			if err != nil {
				return err
			}
			if z.Activ {
				if tip.Activ == 0 && !(exista && veche.Activ == 1 && veche.Cod == z.Cod) {
					return eroare("absente.tip_invalid")
				}
				ciorna, err := q.CitesteCiorna(ctx, acces.CitesteCiornaParams{AngajatID: z.AngajatID, Data: z.Data})
				if err != nil && !errors.Is(err, sql.ErrNoRows) {
					return err
				}
				if err == nil && strings.TrimSpace(ciorna.Valoare) != "" {
					return eroare("absente.pontaj_existent")
				}
				inceput, _ := time.ParseInLocation(time.DateOnly, z.Data, local)
				pontaj, err := q.PontajExistentPeZi(ctx, acces.PontajExistentPeZiParams{AngajatID: z.AngajatID, Data: z.Data, InceputUtc: inceput.UTC().Format(time.RFC3339), SfarsitUtc: inceput.AddDate(0, 0, 1).UTC().Format(time.RFC3339)})
				if err != nil {
					return err
				}
				if pontaj {
					return eroare("absente.pontaj_existent")
				}
			}
			if exista && veche.Cod == z.Cod && veche.Activ == numarBool(z.Activ) {
				continue
			}
			acum := moment()
			revizie, err := q.SalveazaZiAbsenta(ctx, acces.SalveazaZiAbsentaParams{AngajatID: z.AngajatID, Data: z.Data, Cod: z.Cod, Activ: numarBool(z.Activ), CreatLa: acum, ModificatLa: acum})
			if err != nil {
				return err
			}
			z.Revizie = strconv.FormatInt(revizie, 10)
			var inainte any
			if exista {
				inainte = veche
			}
			if err := jurnal(ctx, q, "absente_zile", z.AngajatID+"/"+z.Data, inainte, z); err != nil {
				return err
			}
		}
		return nil
	})
}
