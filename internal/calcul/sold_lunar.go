package calcul

// SoldLunar aplica formula confirmata, exclusiv in minute.
type SoldLunar struct {
	Inceput, Lucratoare, Final int64
	Definit                    bool
}

func CalculeazaSoldLunar(inceput, lucrate int64, norma NormaDatorata, complet bool) SoldLunar {
	return SoldLunar{Inceput: inceput, Lucratoare: norma.Minute, Final: inceput + lucrate - norma.Minute, Definit: norma.Definita && complet}
}
