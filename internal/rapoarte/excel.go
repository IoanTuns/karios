package rapoarte

import (
	"fmt"
	"github.com/xuri/excelize/v2"
)

// GenereazaExcel exporta exact documentul previzualizat, inclusiv filtrele si grupele.
func GenereazaExcel(d Document) ([]byte, error) {
	f := excelize.NewFile()
	defer f.Close()
	foaie := f.GetSheetName(0)
	stil, err := f.NewStyle(&excelize.Style{Font: &excelize.Font{Bold: true}, Fill: excelize.Fill{Type: "pattern", Color: []string{"E7EFEA"}, Pattern: 1}})
	if err != nil {
		return nil, err
	}
	scrie := func(rand int, valori []string, antet bool) error {
		for i, v := range valori {
			celula, err := excelize.CoordinatesToCellName(i+1, rand)
			if err != nil {
				return err
			}
			if err := f.SetCellStr(foaie, celula, v); err != nil {
				return err
			}
			if antet {
				if err := f.SetCellStyle(foaie, celula, celula, stil); err != nil {
					return err
				}
			}
		}
		return nil
	}
	if err := scrie(1, []string{d.Titlu}, true); err != nil {
		return nil, err
	}
	if err := scrie(2, []string{d.Nota}, false); err != nil {
		return nil, err
	}
	antet := []string{}
	for _, c := range d.Coloane {
		antet = append(antet, c.Titlu)
	}
	if err := scrie(4, antet, true); err != nil {
		return nil, err
	}
	for i, r := range d.Randuri {
		if err := scrie(i+5, r.Valori, r.Grup); err != nil {
			return nil, err
		}
	}
	for i, l := range d.Legenda {
		if err := scrie(len(d.Randuri)+6+i, []string{l}, false); err != nil {
			return nil, err
		}
	}
	if err := scrie(len(d.Randuri)+6+len(d.Legenda), []string{d.Metadate}, false); err != nil {
		return nil, err
	}
	// Stilurile se aplica exclusiv pe documentul pregatit/exportat.
	stiluriExcel := map[string]int{}
	aplica := func(celula string, font float64, bold bool, chei ...string) error {
		activ := bold && len(chei) > 0 && chei[0] == "tabel"
		for _, k := range chei {
			_, ok := d.Stiluri[k]
			activ = activ || ok
		}
		if !activ {
			return nil
		}
		st := StilPentru(d, font, "L", bold, chei...)
		aliniere := map[string]string{"L": "left", "C": "center", "R": "right"}[st.Aliniere]
		vechi, e := f.GetCellStyle(foaie, celula)
		if e != nil {
			return e
		}
		cheie := fmt.Sprintf("%d/%g/%t/%t/%s", vechi, st.FontPT, st.Aldin, st.Cursiv, aliniere)
		id, ok := stiluriExcel[cheie]
		if !ok {
			v, e := f.GetStyle(vechi)
			if e != nil {
				return e
			}
			if v.Font == nil {
				v.Font = &excelize.Font{}
			}
			v.Font.Bold = st.Aldin
			v.Font.Italic = st.Cursiv
			v.Font.Size = st.FontPT
			if v.Alignment == nil {
				v.Alignment = &excelize.Alignment{}
			}
			v.Alignment.Horizontal = aliniere
			v.Alignment.WrapText = true
			id, e = f.NewStyle(v)
			if e != nil {
				return e
			}
			stiluriExcel[cheie] = id
		}
		return f.SetCellStyle(foaie, celula, celula, id)
	}
	if err := aplica("A1", 11, true, "titlu"); err != nil {
		return nil, err
	}
	for i, c := range d.Coloane {
		cel, _ := excelize.CoordinatesToCellName(i+1, 4)
		if err := aplica(cel, 11, true, "antet_tabel", "eticheta:"+c.Cod); err != nil {
			return nil, err
		}
	}
	for ri, r := range d.Randuri {
		for ci, c := range d.Coloane {
			cel, _ := excelize.CoordinatesToCellName(ci+1, ri+5)
			if err := aplica(cel, 11, r.Aldin || r.Grup, "tabel", "coloana:"+c.Cod, fmt.Sprintf("celula:%d:%s", ri, c.Cod)); err != nil {
				return nil, err
			}
		}
	}
	if d.Personalizat || len(d.Stiluri) > 0 {
		row := len(d.Randuri) + 8 + len(d.Legenda)
		if err := scrie(row, []string{d.Introducere}, false); err != nil {
			return nil, err
		}
		if err := aplica(fmt.Sprintf("A%d", row), 11, false, "introducere"); err != nil {
			return nil, err
		}
		row++
		for _, bloc := range []struct {
			prefix string
			texte  []string
		}{{"antet", d.Antet}, {"subsol", d.Subsol}} {
			if err := scrie(row, bloc.texte, false); err != nil {
				return nil, err
			}
			for i := range bloc.texte {
				cel, _ := excelize.CoordinatesToCellName(i+1, row)
				if err := aplica(cel, 11, false, bloc.prefix+"_"+[]string{"stanga", "centru", "dreapta"}[i]); err != nil {
					return nil, err
				}
			}
			row++
		}
	}
	b, err := f.WriteToBuffer()
	if err != nil {
		return nil, err
	}
	return b.Bytes(), nil
}
