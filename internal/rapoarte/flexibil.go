package rapoarte

import (
	"context"
	"database/sql"
	"fmt"
	"pontix/internal/calcul"
	"pontix/internal/datebaza"
	"pontix/internal/i18n"
	raportFlexibil "pontix/internal/rapoarte/flexibil"
	"pontix/internal/versiune"
	"strconv"
	"strings"
	"time"
)

func (s *Serviciu) PrevizualizeazaFlexibil(ctx context.Context, d Definitie, c Cerere, acum time.Time) (Document, error) {
	if !esteFlexibil(d.Continut.Sursa) {
		return Document{}, eroare("rapoarte.flex_invalid")
	}
	if e := Valideaza(d); e != nil {
		return Document{}, e
	}
	if c.An < 1 || c.An > 9999 || c.Luna < 1 || c.Luna > 12 || (c.Limba != "ro" && c.Limba != "en") {
		return Document{}, eroare("rapoarte.parametri_invalizi")
	}
	tx, e := s.db.BeginTx(ctx, &sql.TxOptions{ReadOnly: true})
	if e != nil {
		return Document{}, traduce(e)
	}
	defer tx.Rollback()
	return preparaFlexibil(ctx, tx, d, c, acum)
}
func preparaFlexibil(ctx context.Context, tx *sql.Tx, d Definitie, c Cerere, acum time.Time) (Document, error) {
	ctx, stop := context.WithTimeout(ctx, 20*time.Second)
	defer stop()
	f := d.Continut.Flexibil
	b, e := incarcatBazaRaport(ctx, tx, c, f)
	if e != nil {
		return Document{}, traduce(e)
	}
	defer b.db.Close()
	sqlText := f.SQL
	if strings.TrimSpace(sqlText) == "" {
		sqlText = sqlImplicit(d.Continut.Sursa)
	}
	param := parametriFlexibili(c)
	r, e := b.query(ctx, sqlText, param)
	if e != nil {
		return Document{}, traduce(e)
	}
	sqlCols := 0
	for _, v := range f.Calculate {
		if v.Activ && v.Mod == "SQL" {
			sqlCols++
		}
	}
	if sqlCols*len(r.Randuri) > 2000 {
		return Document{}, eroare("rapoarte.prea_mare")
	}
	for _, cc := range f.Calculate {
		if !cc.Activ {
			continue
		}
		for _, nom := range r.Coloane {
			if strings.EqualFold(nom, cc.Cod) {
				return Document{}, eroare("rapoarte.sql_coloane")
			}
		}
		for ri, row := range r.Randuri {
			val := map[string]*float64{}
			pp := map[string]any{}
			for k, v := range param {
				pp[k] = v
			}
			for i, nom := range r.Coloane {
				pp[nom] = row[i]
				switch v := row[i].(type) {
				case int64:
					x := float64(v)
					val[nom] = &x
				case float64:
					x := v
					val[nom] = &x
				case nil:
					val[nom] = nil
				}
			}
			if pp["angajat_id"] == "" {
				delete(pp, "angajat_id")
			}
			var n *float64
			if cc.Mod == "FORMULA" {
				n, e = calcul.EvalueazaFormulaRaport(cc.Expresie, val)
			} else {
				var rr DateFlexibile
				rr, e = b.query(ctx, cc.Expresie, pp, true)
				if e == nil {
					if len(rr.Coloane) != 1 || len(rr.Randuri) > 1 {
						e = eroare("rapoarte.sql_scalar")
					} else if len(rr.Randuri) == 1 {
						switch v := rr.Randuri[0][0].(type) {
						case int64:
							x := float64(v)
							n = &x
						case float64:
							n = &v
						case nil:
						default:
							e = eroare("rapoarte.sql_scalar")
						}
					}
				}
			}
			if e != nil {
				return Document{}, traduce(e)
			}
			if n == nil && cc.LipsaZero {
				x := 0.
				n = &x
			}
			var value any
			if n != nil {
				value = *n
			}
			r.Randuri[ri] = append(row, value)
		}
		r.Coloane = append(r.Coloane, cc.Cod)
	}
	limba := i18n.Limba(c.Limba)
	t := func(k string) string { return i18n.T(limba, "rapoarte."+k, nil) }
	doc := Document{VariabileRaport: map[string]string{"an": strconv.Itoa(c.An), "luna": strconv.Itoa(c.Luna), "perioada": fmt.Sprintf("%04d-%02d", c.An, c.Luna)}, Limba: limba, Flexibil: f, Schema: r.Schema, Titlu: d.Denumire.In(limba), Sursa: d.Continut.Sursa, VersiuneRaport: raportFlexibil.Versiune, RevizieModel: d.Revizie, Pagina: d.Pagina, GeneratLa: acum, Coloane: []Coloana{}, Randuri: []Rand{}, Antet: []string{d.Antet.Stanga.In(limba), d.Antet.Centru.In(limba), d.Antet.Dreapta.In(limba)}, Subsol: []string{d.Subsol.Stanga.In(limba), d.Subsol.Centru.In(limba), d.Subsol.Dreapta.In(limba)}, AntetImagini: d.Antet.Imagini, SubsolImagini: d.Subsol.Imagini, Introducere: d.Continut.Text.In(limba), FaraDate: t("fara_date"), Nota: t("nota_schita"), Paginare: t("paginare"), Metadate: fmt.Sprintf("%s | %s | %s | %s | %s", d.Denumire.In(limba), d.Revizie, versiune.Aplicatie, versiune.Commit, acum.Format(time.RFC3339)+" | "+datebaza.Actor(ctx))}
	for k, v := range b.variabile {
		doc.VariabileRaport[k] = v
	}
	prec := map[string]int{}
	for _, cc := range f.Calculate {
		if cc.Activ {
			prec[cc.Cod] = cc.Zecimale
		}
	}
	for _, nom := range r.Coloane {
		label := strings.ReplaceAll(nom, "_", " ")
		for _, cc := range f.Calculate {
			if cc.Cod == nom {
				label = cc.Denumire.In(limba)
				if label == "" {
					label = nom
				}
			}
		}
		doc.Coloane = append(doc.Coloane, Coloana{Cod: nom, Titlu: label, Pondere: 10})
	}
	for _, row := range r.Randuri {
		v := []string{}
		for i, x := range row {
			var txt string
			switch x := x.(type) {
			case nil:
				txt = "\u2014"
			case float64:
				dec, ok := prec[r.Coloane[i]]
				if !ok {
					dec = 4
				}
				txt = strconv.FormatFloat(x, 'f', dec, 64)
				if !ok {
					txt = strings.TrimRight(strings.TrimRight(txt, "0"), ".")
				}
			default:
				txt = fmt.Sprint(x)
			}
			v = append(v, txt)
		}
		doc.Randuri = append(doc.Randuri, Rand{Valori: v})
	}
	for _, r := range doc.Randuri {
		val := map[string]string{}
		for i, c := range doc.Coloane {
			val[c.Cod] = r.Valori[i]
		}
		doc.valoriRanduri = append(doc.valoriRanduri, val)
	}
	if len(doc.Randuri) > 0 {
		for i, c := range doc.Coloane {
			doc.VariabileRaport[c.Cod] = doc.Randuri[0].Valori[i]
		}
	}
	ascunse := map[string]bool{}
	for _, nom := range f.Ascunse {
		gasit := false
		for _, c := range doc.Coloane {
			gasit = gasit || c.Cod == nom
		}
		if !gasit {
			return Document{}, eroare("rapoarte.sql_coloane")
		}
		ascunse[nom] = true
	}
	vizibile := []Coloana{}
	indici := []int{}
	for i, c := range doc.Coloane {
		if !ascunse[c.Cod] {
			vizibile = append(vizibile, c)
			indici = append(indici, i)
		}
	}
	if len(vizibile) == 0 {
		return Document{}, eroare("rapoarte.sql_coloane")
	}
	for i, r := range doc.Randuri {
		v := []string{}
		for _, j := range indici {
			v = append(v, r.Valori[j])
		}
		doc.Randuri[i].Valori = v
	}
	doc.Coloane = vizibile
	if err := valideazaVariabileFlexibile(doc); err != nil {
		return Document{}, err
	}
	doc.Legenda = append(doc.Legenda, i18n.T(limba, "rapoarte.flex_legenda", nil), "SQL: "+sqlText)
	for _, cc := range f.Calculate {
		if cc.Activ {
			doc.Legenda = append(doc.Legenda, cc.Cod+" ["+cc.Unitate+"]: "+cc.Expresie)
		}
	}
	if d.Continut.AscundeExplicatii {
		doc.Legenda = nil
	}
	return doc, nil
}

func valideazaVariabileFlexibile(d Document) error {
	valori := map[string]bool{"titlu": true, "perioada": true, "generat_la": true, "tabel": true, "an": true, "luna": true, "sectie": true, "departament": true, "compartiment": true, "angajat": true}
	for _, c := range d.Coloane {
		valori[c.Cod] = true
	}
	for k := range d.VariabileRaport {
		valori[k] = true
	}
	texte := []string{d.Introducere}
	if d.Flexibil.ModAspect == "MARKDOWN" {
		texte = append(texte, d.Flexibil.Markdown)
	}
	if d.Flexibil.ModAspect == "LIBER" {
		for _, e := range d.Flexibil.Elemente {
			texte = append(texte, e.Text)
		}
	}
	for _, texteBloc := range [][]string{d.Antet, d.Subsol} {
		texte = append(texte, texteBloc...)
	}
	for _, s := range texte {
		for _, v := range variabila.FindAllStringSubmatch(s, -1) {
			if !valori[strings.TrimSpace(v[1])] {
				return eroare("rapoarte.variabila_invalida")
			}
		}
	}
	return nil
}
