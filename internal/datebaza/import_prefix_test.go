package datebaza

import (
	"context"
	"testing"
)

// csvImport construieste un fisier CSV minimal pentru import.
func csvImport(randuri ...string) []byte {
	s := "Nume;Prenume;Prefix;Marca\n"
	for _, r := range randuri {
		s += r + "\n"
	}
	return []byte(s)
}

func randDupaMarca(t *testing.T, rez RezultatAnalizaImportAngajati, marca string) RandImportAngajat {
	t.Helper()
	for _, r := range rez.Randuri {
		if r.Marca == marca {
			return r
		}
	}
	t.Fatalf("randul %s lipseste din analiza", marca)
	return RandImportAngajat{}
}

func prefixRO(a Angajat) string {
	if a.Prefix == nil {
		return ""
	}
	return a.Prefix["ro"]
}

// F07-T41 (D089): coloana Prefix din sablon se potriveste cu nomenclatorul.
func TestImportPrefixDinSablon(t *testing.T) {
	m := pregateste(t)
	ctx := context.Background()
	for _, format := range []string{"csv", "xlsx"} {
		continut, nume, err := GenereazaSablonAngajati(format)
		if err != nil {
			t.Fatal(err)
		}
		rez, err := m.serviciu.AnalizeazaImportAngajati(ctx, continut, nume)
		if err != nil {
			t.Fatal(err)
		}
		if rez.TotalValide != 3 || randDupaMarca(t, rez, "M001").PrefixImport != "Dr." {
			t.Fatalf("%s: analiza %+v", format, rez)
		}
	}
	continut, _, _ := GenereazaSablonAngajati("csv")
	rez, _ := m.serviciu.AnalizeazaImportAngajati(ctx, continut, "a.csv")
	if _, err := m.serviciu.ExecutaImportAngajati(ctx, rez.Randuri, false); err != nil {
		t.Fatal(err)
	}
	a := dupaMarca(t, m.serviciu)
	if prefixRO(a["M001"]) != "Dr." || prefixRO(a["M002"]) != "" || prefixRO(a["M003"]) != "Biol." {
		t.Fatalf("prefixe salvate: %q %q %q", prefixRO(a["M001"]), prefixRO(a["M002"]), prefixRO(a["M003"]))
	}
}

// Potrivirea ignora majusculele, spatiile si punctul; un prefix necunoscut
// se raporteaza cu cod traductibil si nu se creeaza.
func TestImportPrefixPotrivireSiNecunoscut(t *testing.T) {
	m := pregateste(t)
	ctx := context.Background()
	rez, err := m.serviciu.AnalizeazaImportAngajati(ctx, csvImport("A;B; dr ;I-01", "C;D;BIOL;I-02", "E;F;Prof.;I-03"), "p.csv")
	if err != nil {
		t.Fatal(err)
	}
	if !randDupaMarca(t, rez, "I-01").Valid || !randDupaMarca(t, rez, "I-02").Valid {
		t.Fatalf("variantele de scriere trebuiau recunoscute: %+v", rez.Randuri)
	}
	r := randDupaMarca(t, rez, "I-03")
	if r.Valid || r.CodEroare != "import.prefix_necunoscut" || r.ParamEroare["prefix"] != "Prof." || r.MesajEroare == "" {
		t.Fatalf("prefix necunoscut: %+v", r)
	}
	inainte := len(citeste(t, m.serviciu).Prefixe)
	if _, err := m.serviciu.ExecutaImportAngajati(ctx, rez.Randuri, false); err != nil {
		t.Fatal(err)
	}
	a := dupaMarca(t, m.serviciu)
	if prefixRO(a["I-01"]) != "Dr." || prefixRO(a["I-02"]) != "Biol." {
		t.Fatalf("prefixe: %+v %+v", a["I-01"], a["I-02"])
	}
	if _, ok := a["I-03"]; ok || len(citeste(t, m.serviciu).Prefixe) != inainte {
		t.Fatal("randul invalid a fost importat sau prefixul a fost creat")
	}
	// Un rand trimis ca valid, dar cu prefix necunoscut, e refuzat la executie.
	r.Valid = true
	verificaCod(t, func() error {
		_, err := m.serviciu.ExecutaImportAngajati(ctx, []RandImportAngajat{r}, false)
		return err
	}(), "import.prefix_necunoscut")
	if _, ok := dupaMarca(t, m.serviciu)["I-03"]; ok {
		t.Fatal("executia refuzata a salvat angajatul")
	}
}

// Celula goala pastreaza prefixul existent la actualizare; o valoare noua il
// inlocuieste. Un prefix inactiv e acceptat doar la cine il are deja.
func TestImportPrefixActualizareSiInactiv(t *testing.T) {
	m := pregateste(t)
	ctx := context.Background()
	importa := func(actualizeaza bool, randuri ...string) RezultatAnalizaImportAngajati {
		t.Helper()
		rez, err := m.serviciu.AnalizeazaImportAngajati(ctx, csvImport(randuri...), "p.csv")
		if err != nil {
			t.Fatal(err)
		}
		if _, err := m.serviciu.ExecutaImportAngajati(ctx, rez.Randuri, actualizeaza); err != nil {
			t.Fatal(err)
		}
		return rez
	}
	importa(false, "A;B;Biol.;U-01")
	importa(true, "A;B;;U-01")
	if p := prefixRO(dupaMarca(t, m.serviciu)["U-01"]); p != "Biol." {
		t.Fatalf("celula goala a schimbat prefixul: %q", p)
	}
	biol := elementTest(t, m, "prefixe", "Biol.")
	biol.Activ = false
	if err := m.serviciu.SalveazaElementNomenclator(ctx, "prefixe", biol); err != nil {
		t.Fatal(err)
	}
	rez := importa(true, "A2;B;Biol.;U-01", "C;D;Biol.;U-02")
	if !randDupaMarca(t, rez, "U-01").Valid {
		t.Fatal("prefixul inactiv existent trebuia acceptat")
	}
	if r := randDupaMarca(t, rez, "U-02"); r.Valid || r.CodEroare != "import.prefix_inactiv" {
		t.Fatalf("prefix inactiv nou: %+v", r)
	}
	importa(true, "A3;B;Dr.;U-01")
	if p := prefixRO(dupaMarca(t, m.serviciu)["U-01"]); p != "Dr." {
		t.Fatalf("prefixul nou nu a inlocuit: %q", p)
	}
}
