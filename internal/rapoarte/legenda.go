package rapoarte

import (
	"pontix/internal/calcul"
	"pontix/internal/i18n"
	"strconv"
	"strings"
)

func legendaColoane(d Document, c calcul.ConfiguratiePontaj, limba i18n.Limba, etichete ...map[string]string) []string {
	t := func(k string) string { return i18n.T(limba, "rapoarte.legenda_"+k, nil) }
	formula := func(r calcul.RegulaColoana) string {
		if r.Sursa == "NEDEFINITA" {
			return t("nedefinita")
		}
		text := i18n.T(limba, "pontaj.sursa_"+r.Sursa, nil)
		if len(r.Coduri) > 0 {
			text += " [" + strings.Join(r.Coduri, ", ") + "]"
		}
		if r.Sursa == "ZILE_ABSENTA" {
			text += "; " + t("zile_inregistrate")
			if r.NumaiAprobate {
				text += "; " + t("aprobate")
			} else {
				text += "; " + t("active")
			}
		} else {
			text += " x " + strconv.FormatFloat(float64(r.FactorSutimi)/100, 'f', -1, 64)
		}
		return text
	}
	rez := []string{}
	zile := false
	for _, col := range d.Coloane {
		cod := col.Cod
		label := col.Titlu
		explicatie := ""
		if len(cod) == 10 && cod[4] == '-' {
			if zile {
				continue
			}
			zile = true
			cod = "zile"
			label = t("zile")
			if d.Sursa == "GRAFIC_GARZI" {
				explicatie = t("zile_garzi")
			} else {
				explicatie = t("zile_ore")
			}
		}
		if d.Sursa == "PONTAJ_DETALIAT" && (cod == "total_L_INCEPUT" || cod == "total_L_FIN" || cod == "total_ore_lucratoare") {
			rez = append(rez, label+": "+i18n.T(limba, "sold.legenda_"+strings.TrimPrefix(cod, "total_"), nil))
			continue
		}
		if strings.HasPrefix(cod, "total_") && d.Sursa == "PONTAJ_DETALIAT" {
			id := strings.TrimPrefix(cod, "total_")
			// F08-T15: coloana veche neconfigurata (CEV...) se explica prin coloana D092 pe care o arata.
			if nou, ok := calcul.EchivalenteRapoarteD092[id]; ok {
				configurata := false
				for _, r := range c.Reguli {
					configurata = configurata || (r.Coloana == id && r.Sursa != "NEDEFINITA")
				}
				if !configurata {
					id = nou
				}
			}
			explicatie = t(id)
			sursa := t("nedefinita")
			for _, r := range c.Reguli {
				if r.Coloana == id {
					sursa = formula(r)
					break
				}
			}
			variante := []string{}
			vazute := map[string]bool{}
			for _, r := range c.Conditionate {
				if r.Activ && r.Regula.Coloana == id {
					v := formula(r.Regula)
					criterii := []string{}
					nume := func(prefix, id string) string {
						if len(etichete) > 0 && etichete[0][prefix+id] != "" {
							return etichete[0][prefix+id]
						}
						return t("profil_configurat")
					}
					for _, f := range []struct{ k, prefix, id string }{{"structura.grup", "grup/", r.GrupID}, {"ui.functie", "functie/", r.FunctieID}, {"ui.contract", "contract/", r.ContractID}} {
						if f.id != "" {
							criterii = append(criterii, i18n.T(limba, f.k, nil)+"="+nume(f.prefix, f.id))
						}
					}
					if r.Pozitie != "" {
						k := "structura.executie"
						if r.Pozitie == "CONDUCERE" {
							k = "structura.conducere"
						}
						criterii = append(criterii, i18n.T(limba, k, nil))
					}
					if r.Garda != "" {
						k := "nu_garda"
						if r.Garda == "DA" {
							k = "garda"
						}
						criterii = append(criterii, t(k))
					}
					if len(criterii) == 0 {
						criterii = append(criterii, t("toate_profilele"))
					}
					v = "(" + strings.Join(criterii, ", ") + ") " + v
					if !vazute[v] {
						vazute[v] = true
						variante = append(variante, v)
					}
				}
			}
			explicatie += " " + t("calcul") + ": " + sursa
			if len(variante) > 0 {
				explicatie += "; " + t("conditionate") + ": " + strings.Join(variante, " / ")
			}
		} else if explicatie == "" {
			explicatie = t(cod)
		}
		rez = append(rez, label+": "+explicatie+".")
	}
	if d.Sursa == "PONTAJ_DETALIAT" {
		rez = append(rez, t("configuratie_curenta"))
	}
	return rez
}
