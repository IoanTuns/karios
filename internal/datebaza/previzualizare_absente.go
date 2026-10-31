package datebaza

import (
	"context"
	"encoding/json"
	"strings"
)

type cheieSimulareAbsente struct{}
type EfectAbsenta struct {
	Sursa   string         `json:"sursa"`
	ID      string         `json:"id"`
	Inainte map[string]any `json:"inainte"`
	Dupa    map[string]any `json:"dupa"`
}

// PrevizualizeazaAbsente executa aceleasi validari si operatii intr-o tranzactie anulata.
// Se returneaza exclusiv campuri de pontaj; documentele si metadatele personale nu ies.
func (s *Serviciu) PrevizualizeazaAbsente(ctx context.Context, zile []ZiAbsenta, an, luna int, celule []CelulaSchita) ([]EfectAbsenta, error) {
	efecte := []EfectAbsenta{}
	ctx = context.WithValue(ctx, cheieSimulareAbsente{}, &efecte)
	copie := append([]ZiAbsenta(nil), zile...)
	var err error
	if len(celule) > 0 {
		err = s.InlocuiestePontajCuAbsente(ctx, an, luna, copie, celule)
	} else {
		err = s.SalveazaAbsente(ctx, copie)
	}
	if err != nil {
		return nil, err
	}
	return efecte, nil
}
func campuriEfect(v any) map[string]any {
	r := map[string]any{}
	b, _ := json.Marshal(v)
	var m map[string]any
	_ = json.Unmarshal(b, &m)
	for k, v := range m {
		n := strings.ToLower(strings.ReplaceAll(k, "_", ""))
		switch n {
		case "angajatid", "data", "valoare", "cod", "schimbid", "activ", "incepututc", "sfarsitutc", "stareaprobare", "start", "sfarsit", "stersla":
			r[n] = v
		}
	}
	return r
}
