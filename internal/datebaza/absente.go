package datebaza

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"errors"
	"strconv"
	"strings"
	"time"
	_ "time/tzdata" // Fusurile orare rămân disponibile în executabilul portabil.
	"unicode/utf8"

	"pontix/internal/bd/acces"
	"pontix/internal/i18n"
)

const limitaDocumentAbsenta = 10 * 1024 * 1024

type TipAbsenta struct {
	Cod       string    `json:"cod"`
	Denumire  i18n.Text `json:"denumire"`
	Categorie string    `json:"categorie"`
	Activ     bool      `json:"activ"`
}

type ZiAbsenta struct {
	AngajatID           string `json:"angajat_id"`
	Data                string `json:"data"`
	Cod                 string `json:"cod"`
	Activ               bool   `json:"activ"`
	Revizie             string `json:"revizie"`
	NrInregistrare      string `json:"nr_inregistrare"`
	InlocuitorID        string `json:"inlocuitor_id"`
	StareAprobare       string `json:"stare_aprobare"`
	AprobatDe           string `json:"aprobat_de"`
	DocumentID          string `json:"document_id"`
	DocumentNume        string `json:"document_nume"`
	DocumentMime        string `json:"document_mime"`
	DocumentContinutB64 string `json:"document_continut_b64,omitempty"`
}

type DocumentAbsenta struct {
	ID          string `json:"id"`
	NumeFisier  string `json:"nume_fisier"`
	TipMime     string `json:"tip_mime"`
	ContinutB64 string `json:"continut_b64"`
}

type AbsenteAn struct {
	An     int          `json:"an"`
	Tipuri []TipAbsenta `json:"tipuri"`
	Zile   []ZiAbsenta  `json:"zile"`
}

func valPtr(p *string) string {
	if p == nil {
		return ""
	}
	return *p
}

func ptrSauNil(s string) *string {
	if s == "" {
		return nil
	}
	return &s
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
		rez = append(rez, ZiAbsenta{
			AngajatID:      r.AngajatID,
			Data:           r.Data,
			Cod:            r.Cod,
			Activ:          r.Activ == 1,
			Revizie:        strconv.FormatInt(r.Revizie, 10),
			NrInregistrare: r.NrInregistrare,
			InlocuitorID:   r.InlocuitorID,
			StareAprobare:  r.StareAprobare,
			AprobatDe:      r.AprobatDe,
			DocumentID:     r.DocumentID,
			DocumentNume:   r.DocumentNume,
			DocumentMime:   r.DocumentMime,
		})
	}
	return rez, nil
}

func (s *Serviciu) CitesteDocumentAbsenta(ctx context.Context, id string) (DocumentAbsenta, error) {
	id = strings.TrimSpace(id)
	if id == "" {
		return DocumentAbsenta{}, eroare("absente.document_inexistent")
	}
	r, err := acces.New(s.db).CitesteDocumentAbsenta(ctx, id)
	if errors.Is(err, sql.ErrNoRows) {
		return DocumentAbsenta{}, eroare("absente.document_inexistent")
	}
	if err != nil {
		return DocumentAbsenta{}, traduceEroare(err)
	}
	return DocumentAbsenta{
		ID:          r.ID,
		NumeFisier:  r.NumeFisier,
		TipMime:     r.TipMime,
		ContinutB64: base64.StdEncoding.EncodeToString(r.Continut),
	}, nil
}

func valideazaDocumentNou(nume, mime, b64 string) ([]byte, error) {
	if nume == "" || utf8.RuneCountInString(nume) > 255 || strings.ContainsAny(nume, "\x00\r\n") {
		return nil, eroare("absente.document_invalid")
	}
	continut, err := base64.StdEncoding.DecodeString(strings.TrimSpace(b64))
	if err != nil || len(continut) == 0 || len(continut) > limitaDocumentAbsenta {
		return nil, eroare("absente.document_invalid")
	}
	switch mime {
	case "application/pdf":
		if !bytes.HasPrefix(continut, []byte("%PDF-")) {
			return nil, eroare("absente.document_invalid")
		}
	case "image/png":
		if !bytes.HasPrefix(continut, []byte{0x89, 'P', 'N', 'G', '\r', '\n', 0x1a, '\n'}) {
			return nil, eroare("absente.document_invalid")
		}
	case "image/jpeg":
		if !bytes.HasPrefix(continut, []byte{0xff, 0xd8, 0xff}) {
			return nil, eroare("absente.document_invalid")
		}
	default:
		return nil, eroare("absente.document_invalid")
	}
	return continut, nil
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
	for i := range zile {
		z := &zile[i]
		k := z.AngajatID + "/" + z.Data
		if z.AngajatID == "" || !dataValida(z.Data) || z.Cod == "" || vazute[k] {
			return eroare("absente.lot_invalid")
		}
		vazute[k] = true
		z.NrInregistrare = strings.TrimSpace(z.NrInregistrare)
		if utf8.RuneCountInString(z.NrInregistrare) > 64 || strings.ContainsAny(z.NrInregistrare, "\x00\r\n\t") {
			return eroare("absente.nr_inregistrare_invalid")
		}
		z.InlocuitorID = strings.TrimSpace(z.InlocuitorID)
		if z.InlocuitorID != "" && z.InlocuitorID == z.AngajatID {
			return eroare("absente.inlocuitor_invalid")
		}
		z.StareAprobare = strings.ToUpper(strings.TrimSpace(z.StareAprobare))
		if z.StareAprobare != "" {
			switch z.StareAprobare {
			case "IN_ASTEPTARE", "APROBATA", "RESPINSA":
			default:
				return eroare("absente.aprobare_invalida")
			}
		}
		z.AprobatDe = strings.TrimSpace(z.AprobatDe)
		if utf8.RuneCountInString(z.AprobatDe) > 120 || strings.ContainsAny(z.AprobatDe, "\x00\r\n\t") {
			return eroare("absente.aprobat_de_invalid")
		}
		z.DocumentID = strings.TrimSpace(z.DocumentID)
		z.DocumentNume = strings.TrimSpace(z.DocumentNume)
		z.DocumentMime = strings.ToLower(strings.TrimSpace(z.DocumentMime))
	}
	local, err := time.LoadLocation("Europe/Bucharest")
	if err != nil {
		return traduceEroare(err)
	}
	return s.modifica(ctx, func(q *acces.Queries) error {
		documenteNoi := map[string]acces.ExistaDocumentAbsentaRow{}
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
			if !z.Activ && z.NrInregistrare == "" && z.InlocuitorID == "" && z.StareAprobare == "" && z.AprobatDe == "" && z.DocumentID == "" && z.DocumentContinutB64 == "" {
				z.NrInregistrare = veche.NrInregistrare
				z.InlocuitorID = valPtr(veche.InlocuitorID)
				z.StareAprobare = veche.StareAprobare
				z.AprobatDe = veche.AprobatDe
				z.DocumentID = valPtr(veche.DocumentID)
			}
			if z.StareAprobare == "" {
				z.StareAprobare = "IN_ASTEPTARE"
			}
			if z.InlocuitorID != "" {
				inlocuitor, err := q.CitesteAngajat(ctx, z.InlocuitorID)
				if errors.Is(err, sql.ErrNoRows) {
					return eroare("absente.inlocuitor_invalid")
				}
				if err != nil {
					return err
				}
				if z.Activ && inlocuitor.Activ == 0 && !(exista && valPtr(veche.InlocuitorID) == z.InlocuitorID) {
					return eroare("absente.inlocuitor_invalid")
				}
			}
			if z.DocumentContinutB64 != "" {
				cheieDoc := z.DocumentNume + "\x00" + z.DocumentMime + "\x00" + z.DocumentContinutB64
				if doc, ok := documenteNoi[cheieDoc]; ok {
					z.DocumentID, z.DocumentNume, z.DocumentMime = doc.ID, doc.NumeFisier, doc.TipMime
				} else {
					continut, err := valideazaDocumentNou(z.DocumentNume, z.DocumentMime, z.DocumentContinutB64)
					if err != nil {
						return err
					}
					idDoc := idNou()
					if err := q.AdaugaDocumentAbsenta(ctx, acces.AdaugaDocumentAbsentaParams{
						ID:         idDoc,
						NumeFisier: z.DocumentNume,
						TipMime:    z.DocumentMime,
						Continut:   continut,
						CreatLa:    moment(),
					}); err != nil {
						return err
					}
					doc = acces.ExistaDocumentAbsentaRow{ID: idDoc, NumeFisier: z.DocumentNume, TipMime: z.DocumentMime}
					documenteNoi[cheieDoc] = doc
					z.DocumentID = idDoc
				}
				z.DocumentContinutB64 = ""
			} else if z.DocumentID != "" {
				doc, err := q.ExistaDocumentAbsenta(ctx, z.DocumentID)
				if errors.Is(err, sql.ErrNoRows) {
					return eroare("absente.document_invalid")
				}
				if err != nil {
					return err
				}
				z.DocumentNume, z.DocumentMime = doc.NumeFisier, doc.TipMime
			} else {
				z.DocumentNume, z.DocumentMime = "", ""
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
			if exista && veche.Cod == z.Cod && veche.Activ == numarBool(z.Activ) &&
				veche.NrInregistrare == z.NrInregistrare && valPtr(veche.InlocuitorID) == z.InlocuitorID &&
				veche.StareAprobare == z.StareAprobare && veche.AprobatDe == z.AprobatDe &&
				valPtr(veche.DocumentID) == z.DocumentID {
				continue
			}
			acum := moment()
			revizie, err := q.SalveazaZiAbsenta(ctx, acces.SalveazaZiAbsentaParams{
				AngajatID:      z.AngajatID,
				Data:           z.Data,
				Cod:            z.Cod,
				Activ:          numarBool(z.Activ),
				NrInregistrare: z.NrInregistrare,
				InlocuitorID:   ptrSauNil(z.InlocuitorID),
				StareAprobare:  z.StareAprobare,
				AprobatDe:      z.AprobatDe,
				DocumentID:     ptrSauNil(z.DocumentID),
				CreatLa:        acum,
				ModificatLa:    acum,
			})
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
		zilePeData := map[string][]acces.ListaAbsenteIntervalRow{}
		for _, z := range zile {
			stare := z.StareAprobare
			if stare == "" {
				stare = "IN_ASTEPTARE"
			}
			if !z.Activ || stare == "RESPINSA" {
				continue
			}
			randuri, ok := zilePeData[z.Data]
			if !ok {
				var err error
				randuri, err = q.ListaAbsenteInterval(ctx, acces.ListaAbsenteIntervalParams{DeLa: z.Data, PanaLa: z.Data})
				if err != nil {
					return err
				}
				zilePeData[z.Data] = randuri
			}
			for _, r := range randuri {
				if r.Activ == 0 || r.StareAprobare == "RESPINSA" {
					continue
				}
				if z.InlocuitorID != "" && r.AngajatID == z.InlocuitorID {
					return eroare("absente.inlocuitor_absent")
				}
				if r.InlocuitorID == z.AngajatID {
					return eroare("absente.titular_este_inlocuitor")
				}
			}
		}
		return nil
	})
}
