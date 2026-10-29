package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"pontix/internal/ai"
	"pontix/internal/calcul"
	"pontix/internal/datebaza"
	"pontix/internal/i18n"
	"pontix/internal/portabil"
)

type RaspunsStareAI struct {
	Stare  ai.StareMediu `json:"stare"`
	Eroare *MesajEroare  `json:"eroare"`
}

type CerereInterogareAI struct {
	Istoric []ai.MesajConversatie `json:"istoric"`
	An      int                   `json:"an"`
	Luna    int                   `json:"luna"`
	// Actiuni: propunerile din conversație, cu starea lor (F14-T17).
	Actiuni []string `json:"actiuni"`
}

type RaspunsInterogareAI struct {
	Raspuns *ai.RaspunsAsistent `json:"raspuns"`
	Eroare  *MesajEroare        `json:"eroare"`
}

type CerereAplicaPropunereAI struct {
	An        int                 `json:"an"`
	Luna      int                 `json:"luna"`
	Propunere ai.PropunereActiune `json:"propunere"`
}

type RaspunsAplicaPropunereAI struct {
	Succes bool         `json:"succes"`
	Eroare *MesajEroare `json:"eroare,omitempty"`
}

func (a *App) folderAI() string {
	radacina := a.radacina
	if radacina == "" {
		foldere, err := portabil.Pregateste()
		if err == nil {
			return foldere.AI
		}
	}
	if filepath.Base(radacina) == portabil.NumeFolderDate {
		return filepath.Join(radacina, "ai")
	}
	return filepath.Join(radacina, portabil.NumeFolderDate, "ai")
}

const fisierFaraAutoPornire = "fara-autostart"

// EsteAutoPornireAI verifică dacă pornirea automată a AI este activă (implicit true dacă nu există fișierul de dezactivare).
func (a *App) EsteAutoPornireAI() bool {
	folder := a.folderAI()
	caleFisier := filepath.Join(folder, fisierFaraAutoPornire)
	_, err := os.Stat(caleFisier)
	return os.IsNotExist(err)
}

// SeteazaAutoPornireAI activează sau dezactivează pornirea automată a motorului AI la pornirea aplicației.
func (a *App) SeteazaAutoPornireAI(activa bool) bool {
	folder := a.folderAI()
	caleFisier := filepath.Join(folder, fisierFaraAutoPornire)
	if activa {
		_ = os.Remove(caleFisier)
	} else {
		_ = os.WriteFile(caleFisier, []byte("fara-autostart"), 0644)
	}
	return a.EsteAutoPornireAI()
}

// autoPornesteAIDacaEPosibil pornește motorul în fundal la lansarea aplicației dacă fișierele sunt prezente și opțiunea e activă.
func (a *App) autoPornesteAIDacaEPosibil() {
	if !a.EsteAutoPornireAI() {
		return
	}
	folder := a.folderAI()
	stare := ai.DetecteazaMediu(folder, 0)
	mgr := a.asiguratManagerAI()
	if !stare.Instalat || mgr.EsteActiv() || len(stare.ModeleDisponibile) == 0 {
		return
	}
	_ = mgr.Porneste(context.Background(), "")
}

// asiguratManagerAI creează managerul o singură dată (F14-T14, SEC-6): pornirea automată
// rulează în paralel cu apelurile din interfață, iar două managere ar lansa două procese.
func (a *App) asiguratManagerAI() *ai.Manager {
	a.aiManagerOnce.Do(func() {
		a.aiManager = ai.NouManager(a.folderAI(), ai.PortImplicit)
	})
	return a.aiManager
}

// stareMediuAI: fișierele de pe disc și starea procesului propriu. Un alt program care
// răspunde pe un port nu contează ca motor activ (F14-T12).
func (a *App) stareMediuAI() ai.StareMediu {
	mgr := a.asiguratManagerAI()
	stare := ai.DetecteazaMediu(a.folderAI(), mgr.Port())
	stare.Activ = mgr.EsteActiv()
	stare.ContextConfigurat = ai.ContextConfigurat
	if stare.Activ {
		stare.ModelActiv = mgr.ModelActiv()
	}
	if stare.Activ {
		stare.Mesaj = "Motorul local este activ și gata de utilizare."
	} else if stare.Instalat {
		stare.Mesaj = "Fișierele locale sunt instalate. Motorul poate fi pornit."
	}
	return stare
}

// StareAI raportează disponibilitatea motorului și starea procesului local.
func (a *App) StareAI() RaspunsStareAI {
	return RaspunsStareAI{Stare: a.stareMediuAI()}
}

// PornesteAI lansează motorul local în fundal dacă fișierele sunt prezente.
func (a *App) PornesteAI(numeModel string) RaspunsStareAI {
	mgr := a.asiguratManagerAI()
	err := mgr.Porneste(a.contextLucru(), numeModel)
	stare := a.stareMediuAI()
	if err != nil {
		return RaspunsStareAI{
			Stare:  stare,
			Eroare: eroareInterfata(fmt.Errorf("eroare_pornire_ai: %w", err)),
		}
	}
	return RaspunsStareAI{Stare: stare}
}

// OpresteAI oprește procesul copil al serverului AI local.
func (a *App) OpresteAI() RaspunsStareAI {
	mgr := a.asiguratManagerAI()
	_ = mgr.Opreste()
	stare := a.stareMediuAI()
	return RaspunsStareAI{Stare: stare}
}

// InterogheazaAI trimite mesajul utilizatorului către asistentul local și așteaptă răspunsul.
// Interfața folosește PornesteSarcinaAI, care nu blochează.
func (a *App) InterogheazaAI(cerere CerereInterogareAI) RaspunsInterogareAI {
	return a.interogheazaAI(a.contextLucru(), cerere)
}

func (a *App) interogheazaAI(ctx context.Context, cerere CerereInterogareAI) RaspunsInterogareAI {
	// Colectează contextul de lucru (angajați și schimburi)
	date, err := datebaza.Nou(a.db).Citeste(ctx)
	if err != nil {
		return RaspunsInterogareAI{Eroare: eroareInterfata(err)}
	}

	schimburi, _ := datebaza.Nou(a.db).ListaSchimburi(ctx)

	contextPontix := ai.ContextPontix{
		Limba:             "ro",
		An:                cerere.An,
		Luna:              cerere.Luna,
		ActiuniAnterioare: cerere.Actiuni,
		CuPontaj: func(an, luna int) map[string]bool {
			rez := map[string]bool{}
			if s, err := datebaza.Nou(a.db).CitesteSchita(ctx, an, luna); err == nil {
				for _, c := range s.Celule {
					if strings.TrimSpace(c.Valoare) != "" {
						rez[c.AngajatID] = true
					}
				}
			}
			return rez
		},
	}

	numeSectii := make(map[string]string)
	for _, s := range date.Sectii {
		numeSectii[s.ID] = s.Denumire.In(i18n.Romana)
	}
	numeCategorii := make(map[string]string)
	for _, c := range date.CategoriiProgram {
		numeCategorii[c.ID] = c.Denumire.In(i18n.Romana)
	}
	numeDepartamente := make(map[string]string)
	for _, d := range date.Departamente {
		numeDepartamente[d.ID] = d.Denumire.In(i18n.Romana)
	}
	hartaAngajati := make(map[string]ai.InfoAngajat)
	for _, ang := range date.Angajati {
		if ang.Activ {
			infoAng := ai.InfoAngajat{
				ID:          ang.ID,
				Nume:        ang.Nume,
				Prenume:     ang.Prenume,
				Marca:       ang.Marca,
				Functie:     ang.Functie,
				Sectie:      numeSectii[ang.SectieID],
				Categorie:   numeCategorii[ang.CategorieProgramID],
				Departament: numeDepartamente[ang.DepartamentID],
			}
			contextPontix.Angajati = append(contextPontix.Angajati, infoAng)
			hartaAngajati[ang.ID] = infoAng
		}
	}

	for _, s := range schimburi {
		if s.Activ {
			interval := ""
			if s.DurataMinute > 0 {
				startH := s.OraStart / 60
				startM := s.OraStart % 60
				endMin := (s.OraStart + s.DurataMinute) % 1440
				endH := endMin / 60
				endM := endMin % 60
				interval = fmt.Sprintf("%02d:%02d–%02d:%02d", startH, startM, endH, endM)
			}
			contextPontix.Schimburi = append(contextPontix.Schimburi, ai.InfoSchimb{
				Cod:      s.Cod,
				Denumire: s.Denumire.In(i18n.Romana),
				Ore:      int(s.DurataMinute / 60),
				Interval: interval,
			})
		}
	}

	// Colectăm concediile și absențele din luna selectată
	absenteAn, errAbs := datebaza.Nou(a.db).CitesteAbsente(ctx, cerere.An)
	prefixLuna := fmt.Sprintf("%04d-%02d-", cerere.An, cerere.Luna)
	if errAbs == nil {
		for _, za := range absenteAn.Zile {
			if za.Activ && strings.HasPrefix(za.Data, prefixLuna) {
				ang := hartaAngajati[za.AngajatID]
				nume := strings.TrimSpace(ang.Nume + " " + ang.Prenume)
				contextPontix.Absente = append(contextPontix.Absente, ai.InfoAbsenta{
					AngajatID: za.AngajatID,
					Marca:     ang.Marca,
					Nume:      nume,
					Data:      za.Data,
					Cod:       za.Cod,
				})
			}
		}
	}

	// Colectăm turele și gărțile planificate în luna selectată
	planLuna, errPlan := datebaza.Nou(a.db).CitestePlanificare(ctx, cerere.An, cerere.Luna)
	hartaPlan := make(map[string]bool)
	if errPlan == nil {
		for _, p := range planLuna.Pozitii {
			if p.Activ {
				ang := hartaAngajati[p.AngajatID]
				nume := strings.TrimSpace(ang.Nume + " " + ang.Prenume)
				codSchimb := p.SchimbID
				oreSchimb := 0
				for _, s := range schimburi {
					if s.ID == p.SchimbID {
						codSchimb = s.Cod
						oreSchimb = int(s.DurataMinute / 60)
						break
					}
				}
				contextPontix.TurePlanificate = append(contextPontix.TurePlanificate, ai.InfoPozitiePlanificare{
					AngajatID: p.AngajatID,
					Marca:     ang.Marca,
					Nume:      nume,
					Data:      p.Data,
					CodSchimb: codSchimb,
					Ore:       oreSchimb,
				})
				hartaPlan[p.AngajatID+"/"+p.Data] = true
			}
		}
	}

	// Adăugăm și valorile completate direct în ciorna de pontaj
	schita, errSchita := datebaza.Nou(a.db).CitesteSchita(ctx, cerere.An, cerere.Luna)
	if errSchita == nil {
		for _, c := range schita.Celule {
			val := strings.TrimSpace(c.Valoare)
			if val == "" || hartaPlan[c.AngajatID+"/"+c.Data] {
				continue
			}
			ang := hartaAngajati[c.AngajatID]
			nume := strings.TrimSpace(ang.Nume + " " + ang.Prenume)
			ore := 0
			for _, s := range schimburi {
				if strings.EqualFold(s.Cod, val) {
					ore = int(s.DurataMinute / 60)
					break
				}
			}
			contextPontix.TurePlanificate = append(contextPontix.TurePlanificate, ai.InfoPozitiePlanificare{
				AngajatID: c.AngajatID,
				Marca:     ang.Marca,
				Nume:      nume,
				Data:      c.Data,
				CodSchimb: val,
				Ore:       ore,
			})
			hartaPlan[c.AngajatID+"/"+c.Data] = true
		}
	}

	// O singură poartă alege calea mesajului (F14-T27): concediu, planificare ghidată, întrebări
	// la care răspunde codul; doar restul ajunge la model. Rutele din cod nu pornesc motorul AI,
	// iar planificarea ghidată nu trimite pontajul echipei la model (F14-T22).
	if len(cerere.Istoric) > 0 {
		if d := ai.Ruteaza(cerere.Istoric[len(cerere.Istoric)-1].Continut, contextPontix); d.Ruta != ai.RutaModel {
			return RaspunsInterogareAI{Raspuns: &d.Raspuns}
		}
	}

	mgr := a.asiguratManagerAI()
	if !mgr.EsteActiv() {
		// Încearcă pornirea automată dacă este instalat
		if err := mgr.Porneste(ctx, ""); err != nil {
			return RaspunsInterogareAI{
				Eroare: eroareInterfata(fmt.Errorf("ai_inactiv: %w", err)),
			}
		}
	}

	client := mgr.Client()
	mesajCurent := ""
	if len(cerere.Istoric) > 0 {
		mesajCurent = cerere.Istoric[len(cerere.Istoric)-1].Continut
	}
	raspuns, err := client.Interogheaza(ctx, cerere.Istoric, ai.ContextPentruCerere(mesajCurent, contextPontix))
	if err != nil {
		if e := eroareInterfata(err); e.Cod != "date.eroare" {
			return RaspunsInterogareAI{Eroare: e}
		}
		return RaspunsInterogareAI{
			Eroare: &MesajEroare{
				Cod:   "ai.eroare_generare",
				Param: map[string]string{"detalii": err.Error()},
			},
		}
	}

	// Nici modelul nu aplică direct o selecție: utilizatorul revizuiește intenția.
	if len(raspuns.Propuneri) > 0 {
		// Ștergerile păstrează confirmarea existentă; formularul nou este pentru planificare.
		stergeri := true
		for _, p := range raspuns.Propuneri {
			stergeri = stergeri && strings.HasPrefix(p.Tip, "sterge_")
		}
		if stergeri {
			return RaspunsInterogareAI{Raspuns: &raspuns}
		}
		intentie := ai.IntentiePlan{AnDeConfirmat: true, MesajOriginal: mesajCurent, Selectie: ai.SelectiePlan{Marci: []string{}}}
		if len(raspuns.Propuneri) == 1 {
			p := raspuns.Propuneri[0]
			if strings.HasPrefix(p.Tip, "planifica_") {
				intentie.CodSchimb, _ = p.Parametri["cod_schimb"].(string)
				intentie.DeLa, _ = p.Parametri["data_inceput"].(string)
				intentie.PanaLa, _ = p.Parametri["data_sfarsit"].(string)
				if d, ok := p.Parametri["data"].(string); ok {
					intentie.DeLa, intentie.PanaLa = d, d
				}
			}
		}
		intentie.CompleteazaClarificari()
		raspuns = ai.RaspunsAsistent{Intentie: &intentie}
	}
	return RaspunsInterogareAI{Raspuns: &raspuns}
}

// AplicaPropunereAI execută acțiunea aprobată de utilizator (salvează în pontaj și, dacă este cazul, în planificare).
func (a *App) AplicaPropunereAI(cerere CerereAplicaPropunereAI) RaspunsAplicaPropunereAI {
	p := cerere.Propunere
	marca, _ := p.Parametri["marca"].(string)
	codSchimb, _ := p.Parametri["cod_schimb"].(string)
	if p.Tip == "sterge_tura" || p.Tip == "sterge_interval" {
		codSchimb = ""
	}

	date, err := datebaza.Nou(a.db).Citeste(a.contextLucru())
	if err != nil {
		return RaspunsAplicaPropunereAI{Eroare: eroareInterfata(err)}
	}
	if p.Tip == ai.TipPlanificaGrup {
		return a.aplicaPropunereGrup(p, date)
	}

	// Angajatul se identifică exact (F14-T13, SEC-4): după marcă sau după numele complet,
	// doar dacă e unic. O potrivire parțială („Ana” în „Ioana”) nu mai e acceptată.
	angajatID, e := identificaAngajat(marca, date.Angajati)
	if e != nil {
		return RaspunsAplicaPropunereAI{Eroare: e}
	}

	// Determinăm lista de date calendaristice vizate
	var zile []string
	switch p.Tip {
	case "planifica_interval", "sterge_interval":
		d1Str, _ := p.Parametri["data_inceput"].(string)
		d2Str, _ := p.Parametri["data_sfarsit"].(string)
		doarLucratoare, _ := p.Parametri["doar_lucratoare"].(bool)
		filtruZile, _ := p.Parametri["filtru_zile"].(string)

		t1, err1 := time.Parse(time.DateOnly, d1Str)
		t2, err2 := time.Parse(time.DateOnly, d2Str)
		if err1 != nil || err2 != nil || t2.Before(t1) {
			return RaspunsAplicaPropunereAI{
				Eroare: &MesajEroare{Cod: "ai.interval_invalid"},
			}
		}

		zile = zileDinInterval(t1, t2, filtruZile, doarLucratoare)
	case "planifica_zile", "sterge_zile":
		if rawZile, ok := p.Parametri["zile"].([]any); ok {
			for _, rz := range rawZile {
				if s, ok := rz.(string); ok && s != "" {
					zile = append(zile, s)
				}
			}
		} else if strZile, ok := p.Parametri["zile"].([]string); ok {
			zile = append(zile, strZile...)
		}
	case "planifica_tura", "sterge_tura":
		dStr, _ := p.Parametri["data"].(string)
		if _, err := time.Parse(time.DateOnly, dStr); err != nil {
			return RaspunsAplicaPropunereAI{
				Eroare: &MesajEroare{Cod: "ai.data_invalida", Param: map[string]string{"data": dStr}},
			}
		}
		zile = append(zile, dStr)
	default:
		return RaspunsAplicaPropunereAI{
			Eroare: &MesajEroare{Cod: "ai.tip_necunoscut", Param: map[string]string{"tip": p.Tip}},
		}
	}

	if len(zile) == 0 {
		return RaspunsAplicaPropunereAI{
			Eroare: &MesajEroare{Cod: "ai.fara_zile"},
		}
	}

	// Luna se ia din datele propunerii, nu din ecran (F14-T10); un interval peste
	// sfârșitul lunii se salvează pe fiecare lună.
	for _, l := range peLuni(zile) {
		if _, err := time.Parse(time.DateOnly, l.zile[0]); err != nil {
			return RaspunsAplicaPropunereAI{Eroare: &MesajEroare{Cod: "ai.data_invalida", Param: map[string]string{"data": l.zile[0]}}}
		}
		if e := a.salveazaPropunere(l.an, l.luna, map[string][]string{angajatID: l.zile}, codSchimb); e != nil {
			return RaspunsAplicaPropunereAI{Eroare: e}
		}
	}
	return RaspunsAplicaPropunereAI{Succes: true}
}

type zileLuna struct {
	an, luna int
	zile     []string
}

// peLuni grupează datele AAAA-LL-ZZ pe luni, în ordine.
func peLuni(zile []string) []zileLuna {
	var rez []zileLuna
	index := map[string]int{}
	for _, z := range zile {
		t, err := time.Parse(time.DateOnly, z)
		if err != nil {
			rez = append(rez, zileLuna{zile: []string{z}})
			continue
		}
		cheie := t.Format("2006-01")
		i, ok := index[cheie]
		if !ok {
			i = len(rez)
			index[cheie] = i
			rez = append(rez, zileLuna{an: t.Year(), luna: int(t.Month())})
		}
		rez[i].zile = append(rez[i].zile, z)
	}
	return rez
}

// aplicaPropunereGrup aplică același schimb mai multor angajați (F14-T08). Angajații se
// identifică doar după marca exactă, iar zilele cu concediu sau absență activă se sar.
func (a *App) aplicaPropunereGrup(p ai.PropunereActiune, date datebaza.Date) RaspunsAplicaPropunereAI {
	codSchimb, _ := p.Parametri["cod_schimb"].(string)
	filtru, _ := p.Parametri["filtru_zile"].(string)
	d1, _ := p.Parametri["data_inceput"].(string)
	d2, _ := p.Parametri["data_sfarsit"].(string)
	t1, err1 := time.Parse(time.DateOnly, d1)
	t2, err2 := time.Parse(time.DateOnly, d2)
	if err1 != nil || err2 != nil || t2.Before(t1) || t1.Year() != t2.Year() || t1.Month() != t2.Month() || codSchimb == "" {
		return RaspunsAplicaPropunereAI{Eroare: &MesajEroare{Cod: "ai.interval_invalid"}}
	}
	var marci []string
	switch v := p.Parametri["marci"].(type) {
	case []any:
		for _, m := range v {
			if s, ok := m.(string); ok {
				marci = append(marci, s)
			}
		}
	case []string:
		marci = v
	}
	dupaMarca := make(map[string]string)
	for _, ang := range date.Angajati {
		if ang.Activ {
			dupaMarca[strings.ToLower(ang.Marca)] = ang.ID
		}
	}
	absente := make(map[string]bool)
	if abs, err := datebaza.Nou(a.db).CitesteAbsente(a.contextLucru(), t1.Year()); err == nil {
		for _, z := range abs.Zile {
			if z.Activ {
				absente[z.AngajatID+"/"+z.Data] = true
			}
		}
	}
	zile := zileDinInterval(t1, t2, filtru, false)
	// Implicit nu se suprascrie nimic (F14-T16): zilele deja completate în pontaj se sar.
	ocupate := make(map[string]bool)
	if doarGoale, ok := p.Parametri["doar_goale"].(bool); !ok || doarGoale {
		if s, err := datebaza.Nou(a.db).CitesteSchita(a.contextLucru(), t1.Year(), int(t1.Month())); err == nil {
			for _, c := range s.Celule {
				if strings.TrimSpace(c.Valoare) != "" {
					ocupate[c.AngajatID+"/"+c.Data] = true
				}
			}
		}
	}
	tinte := make(map[string][]string)
	for _, m := range marci {
		id, ok := dupaMarca[strings.ToLower(strings.TrimSpace(m))]
		if !ok {
			return RaspunsAplicaPropunereAI{Eroare: &MesajEroare{Cod: "ai.angajat_negasit", Param: map[string]string{"cautat": m}}}
		}
		for _, zi := range zile {
			if !absente[id+"/"+zi] && !ocupate[id+"/"+zi] {
				tinte[id] = append(tinte[id], zi)
			}
		}
	}
	if len(tinte) == 0 {
		return RaspunsAplicaPropunereAI{Eroare: &MesajEroare{Cod: "ai.fara_zile"}}
	}
	if e := a.salveazaPropunere(t1.Year(), int(t1.Month()), tinte, codSchimb); e != nil {
		return RaspunsAplicaPropunereAI{Eroare: e}
	}
	return RaspunsAplicaPropunereAI{Succes: true}
}

// zileDinInterval întoarce zilele dintre t1 și t2 (inclusiv) care trec de filtru:
// "weekend", "tot_al_doilea_weekend", "lucratoare" (sau doarLucratoare) ori toate.
func zileDinInterval(t1, t2 time.Time, filtruZile string, doarLucratoare bool) []string {
	var zile []string
	// Referința pentru calculul weekendurilor alternate este prima sâmbătă din sau după t1
	refSambata := t1
	for refSambata.Weekday() != time.Saturday && !refSambata.After(t2) {
		refSambata = refSambata.AddDate(0, 0, 1)
	}
	for cur := t1; !cur.After(t2); cur = cur.AddDate(0, 0, 1) {
		ziSapt := cur.Weekday()
		esteWeekend := ziSapt == time.Saturday || ziSapt == time.Sunday
		switch {
		case filtruZile == "tot_al_doilea_weekend":
			if esteWeekend {
				sambataCurenta := cur
				if ziSapt == time.Sunday {
					sambataCurenta = cur.AddDate(0, 0, -1)
				}
				diferentaZile := int(sambataCurenta.Sub(refSambata).Hours() / 24)
				if diferentaZile >= 0 && (diferentaZile/7)%2 == 0 {
					zile = append(zile, cur.Format(time.DateOnly))
				}
			}
		case filtruZile == "weekend":
			if esteWeekend {
				zile = append(zile, cur.Format(time.DateOnly))
			}
		case doarLucratoare || filtruZile == "lucratoare":
			if !esteWeekend {
				zile = append(zile, cur.Format(time.DateOnly))
			}
		default:
			zile = append(zile, cur.Format(time.DateOnly))
		}
	}
	return zile
}

// salveazaPropunere scrie codul schimbului în pontaj pentru zilele fiecărui angajat, într-o
// singură salvare, și, pentru schimburile de planificat (gărzi, după-amiază), în planificare.
func (a *App) salveazaPropunere(an, luna int, tinte map[string][]string, codSchimb string) *MesajEroare {
	ctx := datebaza.CuAsistentAI(a.contextLucru()) // autorul din jurnal: „local/asistent_ai” (F14-T11)
	schita, err := datebaza.Nou(a.db).CitesteSchita(ctx, an, luna)
	if err != nil {
		return eroareInterfata(err)
	}
	revizii := make(map[string]string)
	for _, c := range schita.Celule {
		revizii[c.AngajatID+"/"+c.Data] = c.Revizie
	}
	ids := make([]string, 0, len(tinte))
	for id := range tinte {
		ids = append(ids, id)
	}
	sort.Strings(ids)

	var celule []datebaza.CelulaSchita
	for _, id := range ids {
		for _, zi := range tinte[id] {
			celule = append(celule, datebaza.CelulaSchita{AngajatID: id, Data: zi, Valoare: codSchimb, Revizie: revizii[id+"/"+zi]})
		}
	}
	if err := datebaza.Nou(a.db).SalveazaSchita(ctx, an, luna, celule); err != nil {
		return eroareInterfata(err)
	}

	// Dacă schimbul este și de tip planificare (gardă sau după-amiază), îl actualizăm și în tabelul de planificare
	schimburi, _ := datebaza.Nou(a.db).ListaSchimburi(ctx)
	var schimbGasit *datebaza.Schimb
	for i := range schimburi {
		if strings.EqualFold(schimburi[i].Cod, codSchimb) {
			schimbGasit = &schimburi[i]
			break
		}
	}
	if schimbGasit == nil {
		return nil
	}
	pProg := calcul.ProgramPrecompletare{
		ID:           schimbGasit.ID,
		Cod:          schimbGasit.Cod,
		OraStart:     schimbGasit.OraStart,
		DurataMinute: schimbGasit.DurataMinute,
		Tip:          schimbGasit.Tip,
		Regim:        schimbGasit.Regim,
	}
	if !calcul.ProgramDePlanificat(pProg) {
		return nil
	}
	sectiune := "GARZI"
	if calcul.ProgramInSectiune(pProg, "DUPA_AMIAZA") {
		sectiune = "DUPA_AMIAZA"
	}
	pozExistente := make(map[string]datebaza.PozitiePlanificare)
	if plan, err := datebaza.Nou(a.db).CitestePlanificare(ctx, an, luna); err == nil {
		for _, poz := range plan.Pozitii {
			pozExistente[poz.AngajatID+"|"+poz.Data+"|"+poz.Sectiune] = poz
		}
	}
	var pozitii []datebaza.PozitiePlanificare
	for _, id := range ids {
		for _, zi := range tinte[id] {
			ex := pozExistente[id+"|"+zi+"|"+sectiune]
			pozitii = append(pozitii, datebaza.PozitiePlanificare{
				ID:        ex.ID,
				AngajatID: id,
				Data:      zi,
				SchimbID:  schimbGasit.ID,
				Sectiune:  sectiune,
				Activ:     codSchimb != "",
				Revizie:   ex.Revizie,
			})
		}
	}
	_ = datebaza.Nou(a.db).SalveazaPlanificare(ctx, an, luna, pozitii)
	return nil
}

// identificaAngajat găsește un angajat activ după marca exactă (cu sau fără cratime) sau,
// dacă nu, după numele complet („Nume Prenume” ori „Prenume Nume”), fără diacritice și
// majuscule. Mai mulți candidați sau niciunul dau eroare.
func identificaAngajat(text string, angajati []datebaza.Angajat) (string, *MesajEroare) {
	norm := func(s string) string {
		return strings.Join(strings.Fields(faraDiacriticeAI.Replace(strings.ToLower(s))), " ")
	}
	cautat := norm(text)
	if cautat == "" {
		return "", &MesajEroare{Cod: "ai.angajat_negasit", Param: map[string]string{"cautat": text}}
	}
	faraCratime := strings.ReplaceAll(strings.ReplaceAll(cautat, "-", ""), " ", "")
	var dupaMarca, dupaNume []datebaza.Angajat
	for _, a := range angajati {
		if !a.Activ {
			continue
		}
		m := norm(a.Marca)
		if m == cautat || strings.ReplaceAll(strings.ReplaceAll(m, "-", ""), " ", "") == faraCratime {
			dupaMarca = append(dupaMarca, a)
		}
		if norm(a.Nume+" "+a.Prenume) == cautat || norm(a.Prenume+" "+a.Nume) == cautat {
			dupaNume = append(dupaNume, a)
		}
	}
	gasiti := dupaMarca
	if len(gasiti) == 0 {
		gasiti = dupaNume
	}
	switch len(gasiti) {
	case 1:
		return gasiti[0].ID, nil
	case 0:
		return "", &MesajEroare{Cod: "ai.angajat_negasit", Param: map[string]string{"cautat": text}}
	default:
		nume := make([]string, len(gasiti))
		for i, a := range gasiti {
			nume[i] = strings.TrimSpace(a.Nume+" "+a.Prenume) + " (" + a.Marca + ")"
		}
		return "", &MesajEroare{Cod: "ai.angajat_ambiguu", Param: map[string]string{"cautat": text, "candidati": strings.Join(nume, ", ")}}
	}
}

var faraDiacriticeAI = strings.NewReplacer("ă", "a", "â", "a", "î", "i", "ș", "s", "ş", "s", "ț", "t", "ţ", "t")
