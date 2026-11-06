package rapoarte

import (
	"context"
	"encoding/json"
	"math"
	"pontix/internal/bd/acces"
	"pontix/internal/datebaza"
	"pontix/internal/i18n"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"
)

type StilText struct {
	FontPT   float64 `json:"font_pt,omitempty"`
	Aliniere string  `json:"aliniere,omitempty"`
	Aldin    *bool   `json:"aldin,omitempty"`
	Cursiv   *bool   `json:"cursiv,omitempty"`
}
type PrezentareRaport struct {
	Stiluri  map[string]StilText  `json:"stiluri,omitempty"`
	Etichete map[string]i18n.Text `json:"etichete,omitempty"`
}
type StilEfectiv struct {
	FontPT        float64
	Aliniere      string
	Aldin, Cursiv bool
}

func (s StilEfectiv) FontPDF() string {
	v := ""
	if s.Aldin {
		v += "B"
	}
	if s.Cursiv {
		v += "I"
	}
	return v
}
func completeazaStil(s StilEfectiv, v StilText) StilEfectiv {
	if v.FontPT != 0 {
		s.FontPT = v.FontPT
	}
	if v.Aliniere != "" {
		s.Aliniere = v.Aliniere
	}
	if v.Aldin != nil {
		s.Aldin = *v.Aldin
	}
	if v.Cursiv != nil {
		s.Cursiv = *v.Cursiv
	}
	return s
}
func StilPentru(d Document, font float64, align string, bold bool, chei ...string) StilEfectiv {
	s := StilEfectiv{FontPT: font, Aliniere: align, Aldin: bold}
	for _, k := range chei {
		if v, ok := d.Stiluri[k]; ok {
			if v.FontPT != 0 {
				s.FontPT = v.FontPT
			}
			if v.Aliniere != "" {
				s.Aliniere = v.Aliniere
			}
			if v.Aldin != nil {
				s.Aldin = *v.Aldin
			}
			if v.Cursiv != nil {
				s.Cursiv = *v.Cursiv
			}
		}
	}
	return s
}
func cheieStilModel(k string) bool {
	if contine([]string{"titlu", "introducere", "antet_stanga", "antet_centru", "antet_dreapta", "subsol_stanga", "subsol_centru", "subsol_dreapta", "tabel", "antet_tabel"}, k) {
		return true
	}
	for _, p := range []string{"coloana:", "eticheta:"} {
		if strings.HasPrefix(k, p) {
			v := strings.TrimPrefix(k, p)
			return len(v) > 0 && len(v) <= 100 && !strings.ContainsAny(v, ":\n\r\x00")
		}
	}
	return false
}
func valideazaStil(s StilText) error {
	if math.IsNaN(s.FontPT) || math.IsInf(s.FontPT, 0) || (s.FontPT != 0 && (s.FontPT < 6 || s.FontPT > 40)) || !contine([]string{"", "L", "C", "R"}, s.Aliniere) {
		return eroare("rapoarte.stil_invalid")
	}
	return nil
}
func valideazaPrezentare(p *PrezentareRaport) error {
	if p == nil {
		return nil
	}
	if len(p.Stiluri) > 200 || len(p.Etichete) > 100 {
		return eroare("rapoarte.stil_invalid")
	}
	for k, s := range p.Stiluri {
		if !cheieStilModel(k) {
			return eroare("rapoarte.stil_invalid")
		}
		if e := valideazaStil(s); e != nil {
			return e
		}
	}
	for k, v := range p.Etichete {
		if len(k) == 0 || len(k) > 100 {
			return eroare("rapoarte.stil_invalid")
		}
		for _, t := range v {
			if utf8.RuneCountInString(t) > 200 || strings.ContainsRune(t, 0) {
				return eroare("rapoarte.text_invalid")
			}
		}
	}
	return nil
}
func aplicaPrezentare(d *Document, p *PrezentareRaport) {
	if p == nil {
		return
	}
	d.Stiluri = map[string]StilText{}
	for k, v := range p.Stiluri {
		d.Stiluri[k] = v
	}
	if len(d.Stiluri) == 0 {
		d.Stiluri = nil
	}
	for i, c := range d.Coloane {
		if v, ok := p.Etichete[c.Cod]; ok && v.In(d.Limba) != "" {
			d.Coloane[i].Titlu = v.In(d.Limba)
		}
	}
}

type EditareTextRaport struct {
	Cheie string    `json:"cheie"`
	Text  string    `json:"text"`
	Stil  *StilText `json:"stil,omitempty"`
}

func tintaText(d *Document, k string) (*string, error) {
	switch k {
	case "titlu":
		return &d.Titlu, nil
	case "introducere":
		return &d.Introducere, nil
	}
	for _, p := range []string{"antet_", "subsol_"} {
		if strings.HasPrefix(k, p) {
			for i, z := range []string{"stanga", "centru", "dreapta"} {
				if k == p+z {
					if p == "antet_" && len(d.Antet) > i {
						return &d.Antet[i], nil
					}
					if p == "subsol_" && len(d.Subsol) > i {
						return &d.Subsol[i], nil
					}
				}
			}
			return nil, eroare("rapoarte.editare_invalida")
		}
	}
	parti := strings.Split(k, ":")
	if len(parti) == 2 && parti[0] == "eticheta" {
		for i, c := range d.Coloane {
			if c.Cod == parti[1] {
				return &d.Coloane[i].Titlu, nil
			}
		}
	}
	if len(parti) == 3 && parti[0] == "celula" {
		r, e := strconv.Atoi(parti[1])
		if e == nil && r >= 0 && r < len(d.Randuri) {
			for i, c := range d.Coloane {
				if c.Cod == parti[2] && i < len(d.Randuri[r].Valori) {
					return &d.Randuri[r].Valori[i], nil
				}
			}
		}
	}
	return nil, eroare("rapoarte.editare_invalida")
}

// EditeazaDocument lucreaza pe o copie profunda; tokenul original ramane neschimbat.
func EditeazaDocument(d Document, editari []EditareTextRaport, actor string, acum time.Time) (Document, error) {
	if len(editari) == 0 || len(editari) > 1000 {
		return Document{}, eroare("rapoarte.editare_invalida")
	}
	b, e := json.Marshal(d)
	if e != nil {
		return Document{}, e
	}
	var n Document
	if e = json.Unmarshal(b, &n); e != nil {
		return Document{}, e
	}
	n.Limba = d.Limba
	n.GeneratLa = d.GeneratLa
	for _, v := range d.valoriRanduri {
		c := map[string]string{}
		for k, x := range v {
			c[k] = x
		}
		n.valoriRanduri = append(n.valoriRanduri, c)
	}
	vaz := map[string]bool{}
	total := 0
	for _, v := range editari {
		total += len(v.Text)
		if total > 2*1024*1024 || utf8.RuneCountInString(v.Text) > 4000 || strings.ContainsRune(v.Text, 0) || vaz[v.Cheie] {
			return Document{}, eroare("rapoarte.editare_invalida")
		}
		vaz[v.Cheie] = true
		p, e := tintaText(&n, v.Cheie)
		if e != nil {
			return Document{}, e
		}
		*p = v.Text
		if v.Stil != nil {
			if e := valideazaStil(*v.Stil); e != nil {
				return Document{}, e
			}
			if n.Stiluri == nil {
				n.Stiluri = map[string]StilText{}
			}
			n.Stiluri[v.Cheie] = *v.Stil
		}
	}
	for ri, r := range n.Randuri {
		for ci, c := range n.Coloane {
			if ri < len(n.valoriRanduri) {
				n.valoriRanduri[ri][c.Cod] = r.Valori[ci]
			}
			if ri == 0 {
				if n.VariabileRaport == nil {
					n.VariabileRaport = map[string]string{}
				}
				n.VariabileRaport[c.Cod] = r.Valori[ci]
			}
		}
	}
	if n.Flexibil != nil {
		if e := valideazaVariabileFlexibile(n); e != nil {
			return Document{}, e
		}
	}
	n.Personalizat = true
	local, errLocal := time.LoadLocation("Europe/Bucharest")
	if errLocal != nil {
		local = time.UTC
	}
	mesaj := i18n.T(n.Limba, "rapoarte.personalizat_marcaj", map[string]string{"autor": actor, "moment": acum.In(local).Format("02.01.2006 15:04")})
	n.Metadate += " | " + mesaj
	n.Nota += " | " + i18n.T(n.Limba, "rapoarte.copie_personalizata", nil)
	return n, nil
}
func (s *Serviciu) JurnalPersonalizare(ctx context.Context, token string, inainte, dupa Document, editari []EditareTextRaport, acum time.Time, cereri ...Cerere) error {
	type valoare struct {
		Cheie, Text string
		Stil        *StilText `json:",omitempty"`
	}
	vechi, noi := []valoare{}, []valoare{}
	for _, v := range editari {
		p, e := tintaText(&inainte, v.Cheie)
		if e != nil {
			return e
		}
		stil, ok := inainte.Stiluri[v.Cheie]
		var sp *StilText
		if ok {
			vStil := stil
			sp = &vStil
		}
		vechi = append(vechi, valoare{v.Cheie, *p, sp})
		p, e = tintaText(&dupa, v.Cheie)
		if e != nil {
			return e
		}
		stil, ok = dupa.Stiluri[v.Cheie]
		sp = nil
		if ok {
			vStil := stil
			sp = &vStil
		}
		noi = append(noi, valoare{v.Cheie, *p, sp})
	}
	type instantaneu struct {
		Cerere       Cerere
		Sursa        string
		RevizieModel string
		TitluInitial string
		GeneratLa    string
		Metadate     string
		Modificari   []valoare
	}
	var cerere Cerere
	if len(cereri) > 0 {
		cerere = cereri[0]
	}
	vb, e := json.Marshal(instantaneu{cerere, inainte.Sursa, inainte.RevizieModel, inainte.Titlu, inainte.GeneratLa.Format(time.RFC3339), inainte.Metadate, vechi})
	if e != nil {
		return e
	}
	nb, e := json.Marshal(instantaneu{cerere, dupa.Sursa, dupa.RevizieModel, inainte.Titlu, dupa.GeneratLa.Format(time.RFC3339), dupa.Metadate, noi})
	if e != nil {
		return e
	}
	vs, ns := string(vb), string(nb)
	tx, e := s.db.BeginTx(ctx, nil)
	if e != nil {
		return traduce(e)
	}
	defer tx.Rollback()
	e = acces.New(tx).AdaugaJurnal(ctx, acces.AdaugaJurnalParams{Tabel: "rapoarte_personalizate", RandID: token, ValoareVeche: &vs, ValoareNoua: &ns, Cand: acum.UTC().Format(time.RFC3339Nano), Cine: datebaza.Actor(ctx)})
	if e != nil {
		return traduce(e)
	}
	return traduce(tx.Commit())
}
