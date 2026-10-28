package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
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
	stare := ai.DetecteazaMediu(folder, ai.PortImplicit)
	if !stare.Instalat || stare.Activ || len(stare.ModeleDisponibile) == 0 {
		return
	}
	mgr := a.asiguratManagerAI()
	_ = mgr.Porneste(context.Background(), "")
}

func (a *App) asiguratManagerAI() *ai.Manager {
	if a.aiManager == nil {
		a.aiManager = ai.NouManager(a.folderAI(), ai.PortImplicit)
	}
	return a.aiManager
}

// StareAI raportează disponibilitatea motorului și starea procesului local.
func (a *App) StareAI() RaspunsStareAI {
	folder := a.folderAI()
	stare := ai.DetecteazaMediu(folder, ai.PortImplicit)
	return RaspunsStareAI{Stare: stare}
}

// PornesteAI lansează motorul local în fundal dacă fișierele sunt prezente.
func (a *App) PornesteAI(numeModel string) RaspunsStareAI {
	mgr := a.asiguratManagerAI()
	err := mgr.Porneste(a.contextLucru(), numeModel)
	stare := ai.DetecteazaMediu(a.folderAI(), mgr.Port())
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
	stare := ai.DetecteazaMediu(a.folderAI(), mgr.Port())
	return RaspunsStareAI{Stare: stare}
}

// InterogheazaAI trimite mesajul utilizatorului către asistentul local.
func (a *App) InterogheazaAI(cerere CerereInterogareAI) RaspunsInterogareAI {
	// Colectează contextul de lucru (angajați și schimburi)
	date, err := datebaza.Nou(a.db).Citeste(a.contextLucru())
	if err != nil {
		return RaspunsInterogareAI{Eroare: eroareInterfata(err)}
	}

	schimburi, _ := datebaza.Nou(a.db).ListaSchimburi(a.contextLucru())

	contextPontix := ai.ContextPontix{
		Limba: "ro",
		An:    cerere.An,
		Luna:  cerere.Luna,
	}

	hartaAngajati := make(map[string]ai.InfoAngajat)
	for _, ang := range date.Angajati {
		if ang.Activ {
			infoAng := ai.InfoAngajat{
				ID:      ang.ID,
				Nume:    ang.Nume,
				Prenume: ang.Prenume,
				Marca:   ang.Marca,
				Functie: ang.Functie,
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
	absenteAn, errAbs := datebaza.Nou(a.db).CitesteAbsente(a.contextLucru(), cerere.An)
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
	planLuna, errPlan := datebaza.Nou(a.db).CitestePlanificare(a.contextLucru(), cerere.An, cerere.Luna)
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
	schita, errSchita := datebaza.Nou(a.db).CitesteSchita(a.contextLucru(), cerere.An, cerere.Luna)
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

	// Interogările analitice deterministe (concedii, disponibilitate weekend, acoperire, statistici)
	// se execută instant, fără a bloca sau porni procesul AI dacă acesta este oprit.
	if len(cerere.Istoric) > 0 {
		ultimulMesaj := cerere.Istoric[len(cerere.Istoric)-1].Continut
		if sinteza, ok := ai.AnalizeazaInterogareSinteza(ultimulMesaj, contextPontix); ok {
			return RaspunsInterogareAI{
				Raspuns: &ai.RaspunsAsistent{
					Mesaj: sinteza,
				},
			}
		}
	}

	mgr := a.asiguratManagerAI()
	if !mgr.EsteActiv() {
		// Încearcă pornirea automată dacă este instalat
		if err := mgr.Porneste(a.contextLucru(), ""); err != nil {
			return RaspunsInterogareAI{
				Eroare: eroareInterfata(fmt.Errorf("ai_inactiv: %w", err)),
			}
		}
	}

	client := ai.NouClient(mgr.Port())
	raspuns, err := client.Interogheaza(a.contextLucru(), cerere.Istoric, contextPontix)
	if err != nil {
		return RaspunsInterogareAI{
			Eroare: &MesajEroare{
				Cod:   "ai.eroare_generare",
				Param: map[string]string{"detalii": err.Error()},
			},
		}
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

	// Găsim angajatul după marcă sau nume (tolerant la formatare)
	var angajatID string
	marcaNorm := strings.ToLower(strings.TrimSpace(marca))
	marcaFaraCratima := strings.ReplaceAll(marcaNorm, "-", "")
	for _, ang := range date.Angajati {
		aNorm := strings.ToLower(strings.TrimSpace(ang.Marca))
		aFaraCratima := strings.ReplaceAll(aNorm, "-", "")
		numeComplet := strings.ToLower(strings.TrimSpace(ang.Nume + " " + ang.Prenume))
		prenumeComplet := strings.ToLower(strings.TrimSpace(ang.Prenume + " " + ang.Nume))
		if aNorm == marcaNorm || aFaraCratima == marcaFaraCratima ||
			numeComplet == marcaNorm || prenumeComplet == marcaNorm ||
			strings.Contains(numeComplet, marcaNorm) || strings.Contains(prenumeComplet, marcaNorm) ||
			(len(ang.Nume) >= 3 && strings.Contains(marcaNorm, strings.ToLower(ang.Nume))) {
			angajatID = ang.ID
			break
		}
	}
	if angajatID == "" {
		return RaspunsAplicaPropunereAI{
			Eroare: &MesajEroare{
				Cod:   "ai.eroare_generare",
				Param: map[string]string{"detalii": fmt.Sprintf("Nu a fost găsit angajatul cu marca sau numele \"%s\".", marca)},
			},
		}
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
				Eroare: &MesajEroare{
					Cod:   "ai.eroare_generare",
					Param: map[string]string{"detalii": "Interval de date invalid."},
				},
			}
		}

		// Referința pentru calculul weekendurilor alternate este prima sâmbătă din sau după t1
		refSambata := t1
		for refSambata.Weekday() != time.Saturday && !refSambata.After(t2) {
			refSambata = refSambata.AddDate(0, 0, 1)
		}

		for cur := t1; !cur.After(t2); cur = cur.AddDate(0, 0, 1) {
			ziSapt := cur.Weekday()
			esteWeekend := (ziSapt == time.Saturday || ziSapt == time.Sunday)

			if filtruZile == "tot_al_doilea_weekend" {
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
			} else if filtruZile == "weekend" {
				if esteWeekend {
					zile = append(zile, cur.Format(time.DateOnly))
				}
			} else if doarLucratoare || filtruZile == "lucratoare" {
				if !esteWeekend {
					zile = append(zile, cur.Format(time.DateOnly))
				}
			} else {
				zile = append(zile, cur.Format(time.DateOnly))
			}
		}
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
				Eroare: &MesajEroare{
					Cod:   "ai.eroare_generare",
					Param: map[string]string{"detalii": fmt.Sprintf("Data invalidă: %s", dStr)},
				},
			}
		}
		zile = append(zile, dStr)
	default:
		return RaspunsAplicaPropunereAI{
			Eroare: &MesajEroare{
				Cod:   "ai.eroare_generare",
				Param: map[string]string{"detalii": fmt.Sprintf("Tip acțiune necunoscut: %s", p.Tip)},
			},
		}
	}

	if len(zile) == 0 {
		return RaspunsAplicaPropunereAI{
			Eroare: &MesajEroare{
				Cod:   "ai.eroare_generare",
				Param: map[string]string{"detalii": "Nu au fost identificate zile eligibile în interval."},
			},
		}
	}

	tPrima, _ := time.Parse(time.DateOnly, zile[0])
	an := cerere.An
	luna := cerere.Luna
	if an <= 0 {
		an = tPrima.Year()
	}
	if luna <= 0 {
		luna = int(tPrima.Month())
	}

	// Citim schița existentă pentru a obține reviziile celulelor curente
	schita, err := datebaza.Nou(a.db).CitesteSchita(a.contextLucru(), an, luna)
	if err != nil {
		return RaspunsAplicaPropunereAI{Eroare: eroareInterfata(err)}
	}
	revizii := make(map[string]string)
	for _, c := range schita.Celule {
		if c.AngajatID == angajatID {
			revizii[c.Data] = c.Revizie
		}
	}

	// Construim celulele pentru pontaj
	var celule []datebaza.CelulaSchita
	for _, zi := range zile {
		celule = append(celule, datebaza.CelulaSchita{
			AngajatID: angajatID,
			Data:      zi,
			Valoare:   codSchimb,
			Revizie:   revizii[zi],
		})
	}

	// Salvăm în pontaj
	if err := datebaza.Nou(a.db).SalveazaSchita(a.contextLucru(), an, luna, celule); err != nil {
		return RaspunsAplicaPropunereAI{Eroare: eroareInterfata(err)}
	}

	// Dacă schimbul este și de tip planificare (gardă sau după-amiază), îl actualizăm și în tabelul de planificare
	schimburi, _ := datebaza.Nou(a.db).ListaSchimburi(a.contextLucru())
	var schimbGasit *datebaza.Schimb
	for i := range schimburi {
		if strings.EqualFold(schimburi[i].Cod, codSchimb) {
			schimbGasit = &schimburi[i]
			break
		}
	}

	if schimbGasit != nil {
		pProg := calcul.ProgramPrecompletare{
			ID:           schimbGasit.ID,
			Cod:          schimbGasit.Cod,
			OraStart:     schimbGasit.OraStart,
			DurataMinute: schimbGasit.DurataMinute,
			Tip:          schimbGasit.Tip,
			Regim:        schimbGasit.Regim,
		}
		if calcul.ProgramDePlanificat(pProg) {
			sectiune := "GARZI"
			if calcul.ProgramInSectiune(pProg, "DUPA_AMIAZA") {
				sectiune = "DUPA_AMIAZA"
			}

			pozExistente := make(map[string]datebaza.PozitiePlanificare)
			if plan, err := datebaza.Nou(a.db).CitestePlanificare(a.contextLucru(), an, luna); err == nil {
				for _, poz := range plan.Pozitii {
					cheie := fmt.Sprintf("%s|%s|%s", poz.AngajatID, poz.Data, poz.Sectiune)
					pozExistente[cheie] = poz
				}
			}

			var pozitii []datebaza.PozitiePlanificare
			for _, zi := range zile {
				cheie := fmt.Sprintf("%s|%s|%s", angajatID, zi, sectiune)
				id := ""
				revizie := ""
				if ex, ok := pozExistente[cheie]; ok {
					id = ex.ID
					revizie = ex.Revizie
				}
				pozitii = append(pozitii, datebaza.PozitiePlanificare{
					ID:        id,
					AngajatID: angajatID,
					Data:      zi,
					SchimbID:  schimbGasit.ID,
					Sectiune:  sectiune,
					Activ:     codSchimb != "",
					Revizie:   revizie,
				})
			}
			_ = datebaza.Nou(a.db).SalveazaPlanificare(a.contextLucru(), an, luna, pozitii)
		}
	}

	return RaspunsAplicaPropunereAI{Succes: true}
}
