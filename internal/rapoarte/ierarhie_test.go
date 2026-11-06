package rapoarte

import (
	"context"
	"strings"
	"testing"
)

func TestMarcajIerarhieInToateSursele(t *testing.T) {
	m := pregateste(t)
	dateGarziDetaliat(t, m)
	for _, sql := range []string{
		`INSERT INTO pontaj_ciorne(angajat_id,data,valoare,modificat_la) VALUES('g','2026-10-06','G24','x')`,
		`INSERT INTO absente_zile(angajat_id,data,cod,activ,creat_la,modificat_la) VALUES('g','2026-10-07','CO',1,'x','x')`,
	} {
		if _, err := m.db.Exec(sql); err != nil {
			t.Fatal(err)
		}
	}
	for _, limba := range []string{"ro", "en"} {
		for _, id := range []string{"raport-pontaj", "raport-pontaj-detaliat", "raport-fisa", "raport-absente", "raport-garzi"} {
			d, err := m.s.Pregateste(context.Background(), Cerere{ModelID: id, An: 2026, Luna: 10, Limba: limba, AngajatID: "g"}, momentFix)
			if err != nil {
				t.Fatal(id, err)
			}
			for _, r := range d.Randuri {
				v := map[string]string{}
				for i, c := range d.Coloane {
					v[c.Cod] = r.Valori[i]
				}
				if v["marca"] == "GARDA-R01" && (!strings.HasPrefix(v["angajat"], "Alt nume Ana (") || !strings.HasSuffix(v["angajat"], ")")) {
					t.Fatal(id, v)
				}
			}
			if id == "raport-fisa" && !strings.Contains(strings.Join(d.Antet, " "), "Alt nume Ana (") {
				t.Fatal(id, d.Antet)
			}
			if id != "raport-fisa" {
				gasit := false
				for _, r := range d.Randuri {
					for _, v := range r.Valori {
						if strings.HasPrefix(v, "Alt nume Ana (") {
							gasit = true
						}
					}
				}
				if !gasit {
					t.Fatal(id, "marcaj lipsa")
				}
			}
		}
	}
}
