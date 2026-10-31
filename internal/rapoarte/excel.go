package rapoarte

import "github.com/xuri/excelize/v2"

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
	b, err := f.WriteToBuffer()
	if err != nil {
		return nil, err
	}
	return b.Bytes(), nil
}
