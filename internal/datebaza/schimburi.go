package datebaza

import (
	"context"
	"database/sql"
	_ "embed"
	"encoding/json"
	"errors"
	"strings"

	"pontix/internal/bd/acces"
	"pontix/internal/i18n"
)

type Schimb struct {
	TipID        string    `json:"tip_id"`
	RegimID      string    `json:"regim_id"`
	ID           string    `json:"id"`
	Cod          string    `json:"cod"`
	Denumire     i18n.Text `json:"denumire"`
	OraStart     int64     `json:"ora_start"`
	DurataMinute int64     `json:"durata_minute"`
	PauzaMinute  int64     `json:"pauza_minute"`
	Tip          string    `json:"tip"`
	Regim        string    `json:"regim"`
	Categorie    string    `json:"categorie"`
	Activ        bool      `json:"activ"`
}

//go:embed initiale/schimburi.json
var schimburiInitiale []byte

func (s *Serviciu) CompleteazaSchimburiInitiale(ctx context.Context, categorie ...string) (int64, error) {
	cat := "LAB"
	if len(categorie) > 0 && strings.TrimSpace(categorie[0]) != "" {
		cat = strings.ToUpper(strings.TrimSpace(categorie[0]))
	}
	if cat != "LAB" {
		return 0, eroare("schimburi.categorie_invalida")
	}
	var initiale []Schimb
	if err := json.Unmarshal(schimburiInitiale, &initiale); err != nil {
		return 0, traduceEroare(err)
	}
	var adaugate int64
	err := s.modificaTx(ctx, func(tx *sql.Tx, q *acces.Queries) error {
		var activ int
		if err := tx.QueryRowContext(ctx, "SELECT activ FROM categorii_schimburi WHERE id=?", cat).Scan(&activ); err != nil {
			return err
		}
		if activ != 1 {
			return eroare("nomenclator.inactiv")
		}
		for _, schimb := range initiale {
			if schimb.Categorie == "" {
				schimb.Categorie = "LAB"
			}
			if schimb.Categorie != cat {
				continue
			}
			// Codul poate lipsi după redenumirea unui șablon deja salvat.
			schimb.ID = idNou()
			text, err := json.Marshal(schimb.Denumire)
			if err != nil {
				return err
			}
			n, err := q.AdaugaSchimbInitial(ctx, acces.AdaugaSchimbInitialParams{ID: schimb.ID, Cod: schimb.Cod, Denumire: string(text), OraStart: schimb.OraStart, DurataMinute: schimb.DurataMinute, Tip: schimb.Tip, Regim: schimb.Regim, Categorie: schimb.Categorie, CreatLa: moment(), ModificatLa: moment()})
			if err != nil {
				return err
			}
			if n == 0 {
				continue
			}
			adaugate += n
			if err := jurnal(ctx, q, "tipuri_schimb", schimb.ID, nil, schimb); err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return 0, err
	}
	return adaugate, nil
}

func (s *Serviciu) ListaSchimburi(ctx context.Context) ([]Schimb, error) {
	randuri, err := acces.New(s.db).ListaSchimburi(ctx)
	if err != nil {
		return nil, traduceEroare(err)
	}
	clasificari, err := acces.New(s.db).ListaClasificariSchimburi(ctx)
	if err != nil {
		return nil, traduceEroare(err)
	}
	rez := []Schimb{}
	for _, r := range randuri {
		var denumire i18n.Text
		if err := json.Unmarshal([]byte(r.Denumire), &denumire); err != nil {
			return nil, traduceEroare(err)
		}
		c, exista := clasificari[r.ID]
		if !exista {
			c = acces.ClasificareSchimb{TipID: r.Tip, RegimID: r.Regim}
		}
		if c.RegimID == "GENERAL" {
			c.RegimID = "LUCRATOARE"
		}
		rez = append(rez, Schimb{TipID: c.TipID, RegimID: c.RegimID, ID: r.ID, Cod: r.Cod, Denumire: denumire, OraStart: r.OraStart, DurataMinute: r.DurataMinute, PauzaMinute: r.PauzaMinute, Tip: r.Tip, Regim: r.Regim, Categorie: r.Categorie, Activ: r.Activ == 1})
	}
	return rez, nil
}

func (s *Serviciu) SalveazaSchimb(ctx context.Context, schimb Schimb) error {
	schimb.Cod = strings.ToUpper(strings.TrimSpace(schimb.Cod))
	if schimb.Cod == "" || len(schimb.Cod) > 16 || strings.IndexFunc(schimb.Cod, func(r rune) bool {
		return !((r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') || r == '-' || r == '_')
	}) >= 0 {
		return eroare("schimburi.cod_invalid")
	}
	denumire := i18n.Text{}
	for l, text := range schimb.Denumire {
		if text = strings.TrimSpace(text); text != "" {
			denumire[l] = text
		}
	}
	if len(denumire) == 0 {
		return eroare("sarbatori.denumire_obligatorie")
	}
	schimb.Denumire = denumire
	if schimb.OraStart < 0 || schimb.OraStart > 1439 || schimb.DurataMinute < 1 || schimb.DurataMinute > 2880 || schimb.PauzaMinute < 0 || schimb.PauzaMinute > schimb.DurataMinute {
		return eroare("schimburi.interval_invalid")
	}
	switch schimb.Tip {
	case "NORMAL", "TURA", "NOAPTE", "GARDA":
	default:
		return eroare("schimburi.tip_invalid")
	}
	return s.modificaTx(ctx, func(tx *sql.Tx, q *acces.Queries) error {
		if _, err := q.CautaCodAbsenta(ctx, schimb.Cod); err == nil {
			return eroare("absente.cod_duplicat")
		} else if !errors.Is(err, sql.ErrNoRows) {
			return err
		}
		if schimb.Regim == "" || schimb.Regim == "GENERAL" {
			schimb.Regim = "LUCRATOARE"
		}
		switch schimb.Regim {
		case "GENERAL", "LUCRATOARE", "NELUCRATOARE", "INTRARE_ZI_LIBERA":
		default:
			return eroare("schimburi.regim_invalid")
		}
		schimb.Categorie = strings.ToUpper(strings.TrimSpace(schimb.Categorie))
		if schimb.Categorie == "" {
			schimb.Categorie = "CUSTOM"
		}
		var activ int
		if err := tx.QueryRowContext(ctx, "SELECT activ FROM categorii_schimburi WHERE id=?", schimb.Categorie).Scan(&activ); err != nil {
			return eroare("schimburi.categorie_invalida")
		}
		if activ != 1 {
			vechi, err := q.CitesteSchimb(ctx, schimb.ID)
			if err != nil || vechi.Categorie != schimb.Categorie {
				return eroare("nomenclator.inactiv")
			}
		}

		idExistent, err := q.CitesteCodSchimb(ctx, schimb.Cod)
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			return err
		}
		if err == nil && idExistent != schimb.ID {
			return eroare("schimburi.cod_duplicat")
		}
		text, err := json.Marshal(denumire)
		if err != nil {
			return err
		}
		clasificari, err := q.ListaClasificariSchimburi(ctx)
		if err != nil {
			return err
		}
		anterioara := clasificari[schimb.ID]
		if schimb.ID != "" && anterioara.TipID == "" {
			if v, e := q.CitesteSchimb(ctx, schimb.ID); e == nil {
				anterioara = acces.ClasificareSchimb{TipID: v.Tip, RegimID: v.Regim}
			}
		}
		if anterioara.RegimID == "GENERAL" {
			anterioara.RegimID = "LUCRATOARE"
		}
		if schimb.TipID == "" {
			schimb.TipID = schimb.Tip
		}
		if schimb.RegimID == "" {
			schimb.RegimID = schimb.Regim
		}
		for _, c := range []struct {
			tabela, id, vechi string
			dest              *string
		}{
			{"tipuri_program", schimb.TipID, anterioara.TipID, &schimb.Tip},
			{"regimuri_program", schimb.RegimID, anterioara.RegimID, &schimb.Regim},
		} {
			var activ int
			if err := tx.QueryRowContext(ctx, "SELECT comportament,activ FROM "+c.tabela+" WHERE id=?", c.id).Scan(c.dest, &activ); err != nil {
				return eroare("nomenclator.invalid")
			}
			if activ != 1 && c.id != c.vechi {
				return eroare("nomenclator.inactiv")
			}
		}
		acum := moment()
		var vechi any
		if schimb.ID == "" {
			schimb.ID = idNou()
			err = q.AdaugaSchimb(ctx, acces.AdaugaSchimbParams{ID: schimb.ID, Cod: schimb.Cod, Denumire: string(text), OraStart: schimb.OraStart, DurataMinute: schimb.DurataMinute, PauzaMinute: schimb.PauzaMinute, Tip: schimb.Tip, Regim: schimb.Regim, Categorie: schimb.Categorie, Activ: numarBool(schimb.Activ), CreatLa: acum, ModificatLa: acum})
		} else {
			r, citireErr := q.CitesteSchimb(ctx, schimb.ID)
			if errors.Is(citireErr, sql.ErrNoRows) {
				return eroare("schimburi.inexistent")
			}
			if citireErr != nil {
				return citireErr
			}
			vechi = r
			_, err = q.ModificaSchimb(ctx, acces.ModificaSchimbParams{ID: schimb.ID, Cod: schimb.Cod, Denumire: string(text), OraStart: schimb.OraStart, DurataMinute: schimb.DurataMinute, PauzaMinute: schimb.PauzaMinute, Tip: schimb.Tip, Regim: schimb.Regim, Categorie: schimb.Categorie, Activ: numarBool(schimb.Activ), ModificatLa: acum})
		}
		if err != nil {
			return err
		}
		if _, err := tx.ExecContext(ctx, "INSERT INTO clasificari_schimburi(schimb_id,tip_id,regim_id) VALUES(?,?,?) ON CONFLICT(schimb_id) DO UPDATE SET tip_id=excluded.tip_id,regim_id=excluded.regim_id", schimb.ID, schimb.TipID, schimb.RegimID); err != nil {
			return err
		}
		if err := jurnal(ctx, q, "clasificari_schimburi", schimb.ID, anterioara, acces.ClasificareSchimb{TipID: schimb.TipID, RegimID: schimb.RegimID}); err != nil {
			return err
		}
		return jurnal(ctx, q, "tipuri_schimb", schimb.ID, vechi, schimb)
	})
}

func (s *Serviciu) SchimbaActivSchimb(ctx context.Context, id string, activ bool) error {
	return s.modifica(ctx, func(q *acces.Queries) error {
		vechi, err := q.CitesteSchimb(ctx, id)
		if errors.Is(err, sql.ErrNoRows) {
			return eroare("schimburi.inexistent")
		}
		if err != nil {
			return err
		}
		if _, err := q.SchimbaActivSchimb(ctx, acces.SchimbaActivSchimbParams{ID: id, Activ: numarBool(activ), ModificatLa: moment()}); err != nil {
			return err
		}
		nou, err := q.CitesteSchimb(ctx, id)
		if err != nil {
			return err
		}
		return jurnal(ctx, q, "tipuri_schimb", id, vechi, nou)
	})
}
