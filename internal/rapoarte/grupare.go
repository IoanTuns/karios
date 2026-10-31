package rapoarte

import (
	"context"
	"database/sql"
	"golang.org/x/text/collate"
	"golang.org/x/text/language"
	"pontix/internal/bd/acces"
	"pontix/internal/calcul"
	"pontix/internal/i18n"
	"sort"
)

type profilRaport struct {
	Grup, Functie, Contract, NumeGrup string
	NumeContract, NormaContract       string
	Ordine                            int
	Conducere, Garda                  bool
}

func profileRapoarte(ctx context.Context, tx *sql.Tx, limba i18n.Limba) (map[string]profilRaport, map[string]string, error) {
	rows, err := tx.QueryContext(ctx, `SELECT a.id,coalesce(f.grup_id,''),coalesce(af.functie_id,''),coalesce(ac.contract_id,a.tip_contract),coalesce(g.denumire,'{}'),coalesce(g.ordine,100001),coalesce(f.conducere,0),CASE WHEN ag.principal_id IS NOT NULL OR coalesce(tc.garda,0)=1 THEN 1 ELSE 0 END,coalesce(ag.principal_id,''),coalesce(tc.denumire,'{}'),coalesce(tc.normare,CASE a.tip_contract WHEN 'PARTIALA' THEN 'INDIVIDUALA' ELSE 'GLOBALA' END),CASE coalesce(tc.normare,CASE a.tip_contract WHEN 'PARTIALA' THEN 'INDIVIDUALA' ELSE 'GLOBALA' END) WHEN 'GLOBALA' THEN (SELECT cast(valoare AS INTEGER) FROM setari WHERE cheie='norma_intreaga_minute') WHEN 'INDIVIDUALA' THEN coalesce(a.norma_partiala_minute,0) ELSE coalesce(tc.norma_minute,0) END FROM angajati a LEFT JOIN angajati_functii af ON af.angajat_id=a.id LEFT JOIN functii f ON f.id=af.functie_id LEFT JOIN grupuri_profesionale g ON g.id=f.grup_id LEFT JOIN angajati_contracte ac ON ac.angajat_id=a.id LEFT JOIN tipuri_contract tc ON tc.id=ac.contract_id LEFT JOIN angajati_garzi ag ON ag.angajat_id=a.id`)
	if err != nil {
		return nil, nil, err
	}
	defer rows.Close()
	rez := map[string]profilRaport{}
	legaturi := map[string]string{}
	for rows.Next() {
		var id, nume, principal, numeContract, modNormare string
		var minute int64
		var p profilRaport
		var sef, garda int
		if err := rows.Scan(&id, &p.Grup, &p.Functie, &p.Contract, &nume, &p.Ordine, &sef, &garda, &principal, &numeContract, &modNormare, &minute); err != nil {
			return nil, nil, err
		}
		p.NumeGrup = textJSON(nume, limba)
		p.NumeContract = textJSON(numeContract, limba)
		if p.NumeContract == "" {
			p.NumeContract = p.Contract
		}
		if modNormare == "FARA" {
			p.NormaContract = i18n.T(limba, "structura.normare_FARA", nil)
		} else {
			p.NormaContract = calcul.TextValoarePontaj(calcul.ValoarePontaj{Sutimi: minute * 100, Unitate: "ORE", Definita: true}) + " h"
			if modNormare == "LUNARA" {
				p.NormaContract += " / " + i18n.T(limba, "structura.luna", nil)
			}
		}
		p.Conducere = sef == 1
		p.Garda = garda == 1
		rez[id] = p
		if principal != "" {
			legaturi[id] = principal
		}
	}
	for id, principal := range legaturi {
		p := rez[id]
		if p.Functie == "" {
			f := rez[principal]
			p.Functie = f.Functie
			p.Grup = f.Grup
			p.NumeGrup = f.NumeGrup
			p.Ordine = f.Ordine
			p.Conducere = f.Conducere
			rez[id] = p
		}
	}
	return rez, legaturi, rows.Err()
}
func (p profilRaport) corespunde(c Cerere) bool {
	return (c.GrupID == "" || p.Grup == c.GrupID) && (c.FunctieID == "" || p.Functie == c.FunctieID) && (c.ContractID == "" || p.Contract == c.ContractID) && (c.Pozitie == "" || p.Conducere == (c.Pozitie == "CONDUCERE")) && (c.Garda == "" || p.Garda == (c.Garda == "DA"))
}
func ordoneazaProfesional(oameni []acces.ListaAngajatiRow, profile map[string]profilRaport, legaturi map[string]string) {
	index := map[string]acces.ListaAngajatiRow{}
	for _, a := range oameni {
		index[a.ID] = a
	}
	cheie := func(id string) string {
		if legaturi[id] != "" {
			return legaturi[id]
		}
		return id
	}
	col := collate.New(language.Romanian)
	sort.SliceStable(oameni, func(i, j int) bool {
		a, b := oameni[i], oameni[j]
		pa, pb := cheie(a.ID), cheie(b.ID)
		if pa == pb {
			if a.ID == b.ID {
				return false
			}
			if a.ID == pa {
				return true
			}
			if b.ID == pb {
				return false
			}
			return a.Marca < b.Marca
		}
		fa, fb := profile[pa], profile[pb]
		if fa.Ordine != fb.Ordine {
			return fa.Ordine < fb.Ordine
		}
		if fa.Grup != fb.Grup {
			return col.CompareString(fa.NumeGrup, fb.NumeGrup) < 0
		}
		if fa.Conducere != fb.Conducere {
			return fa.Conducere
		}
		na, nb := index[pa], index[pb]
		if cmp := col.CompareString(na.Nume+" "+na.Prenume, nb.Nume+" "+nb.Prenume); cmp != 0 {
			return cmp < 0
		}
		return na.Marca < nb.Marca
	})
}
