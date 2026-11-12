package main

import (
	"context"
	"database/sql"
	"sync"

	"pontix/internal/ai"
	"pontix/internal/audit"
	"pontix/internal/conturi"
	"pontix/internal/i18n"
	"pontix/internal/seif"
	"pontix/internal/stocare"
	"pontix/internal/versiune"
)

// App conține metodele apelate din interfața React.
type App struct {
	// doarCitire: alta instanta tine folderul de date (D002 punctul 6, F05-T11).
	doarCitire bool
	detinator  *stocare.InfoBlocare
	// conversia modelelor de raport vechi (F09-T23), o data pe sesiune
	conversieMu        sync.Mutex
	rapoarteConvertite bool
	ctx                context.Context
	db                 *sql.DB
	radacina           string
	raportMu           sync.Mutex
	raportCache        map[string]raportPregatit
	raportOrdine       []string
	aiManager          *ai.Manager
	aiManagerOnce      sync.Once
	aiMu               sync.Mutex
	aiSarcina          *sarcinaAI
	paza               *audit.Paza
	// seif e seiful deschis la pornire sau la deblocare (D054); codRecuperare, cheia de
	// recuperare încă neconfirmată, de afișat administratorului.
	seif          *seif.Seif
	codRecuperare string

	inchidereMu         sync.Mutex
	modificariNesalvate bool
	limbaInterfata      string

	contMu            sync.Mutex
	cont              *conturi.Cont
	autoLoginIncercat bool
}

func NewApp(db *sql.DB) *App {
	return &App{db: db}
}

func (a *App) pornire(ctx context.Context) {
	a.ctx = ctx
	go a.autoPornesteAIDacaEPosibil()
}

func (a *App) oprire(ctx context.Context) {
	// Și dacă pornirea automată creează managerul chiar acum, procesul se oprește.
	_ = a.asiguratManagerAI().Opreste()
	// Contorul din seif reține starea de la închidere (F06-T09); nu și peste un incident nerezolvat.
	if a.db == nil {
		return // datele nu au fost deblocate
	}
	// Doar instanta care scrie actualizeaza contorul din seif (F05-T11).
	if a.seif != nil && !a.paza.Blocata() && !a.doarCitire {
		_ = stocare.ActualizeazaSecventa(context.Background(), a.db, a.seif)
	}
	a.db.Close()
}

// Informatii descrie versiunea care rulează, pentru ecranul „Despre” și subsolul rapoartelor.
type Informatii struct {
	Versiune string `json:"versiune"`
	Commit   string `json:"commit"`
}

func (a *App) Informatii() Informatii {
	return Informatii{Versiune: versiune.Aplicatie, Commit: versiune.Commit}
}

// Catalog întoarce textele interfeței în limba cerută.
func (a *App) Catalog(limba string) map[string]string {
	return i18n.Catalog(i18n.Limba(limba))
}

// TraduceriLipsa întoarce cheile din limba implicită fără traducere în limba cerută (F03-T08),
// pentru ecranul de limbă din Setări. O limbă necunoscută le are pe toate lipsă.
func (a *App) TraduceriLipsa(limba string) []string {
	rez := i18n.Lipsa(i18n.Limba(limba))
	if rez == nil {
		rez = []string{}
	}
	return rez
}

// Limbi întoarce codurile limbilor disponibile, cu limba implicită prima.
func (a *App) Limbi() []string {
	var rez []string
	for _, l := range i18n.Limbi() {
		rez = append(rez, string(l))
	}
	return rez
}
