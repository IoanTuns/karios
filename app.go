package main

import (
	"context"
	"database/sql"
	"sync"

	"pontix/internal/ai"
	"pontix/internal/audit"
	"pontix/internal/i18n"
	"pontix/internal/versiune"
)

// App conține metodele apelate din interfața React.
type App struct {
	ctx          context.Context
	db           *sql.DB
	radacina     string
	raportMu     sync.Mutex
	raportCache  map[string]raportPregatit
	raportOrdine []string
	aiManager    *ai.Manager
	aiMu         sync.Mutex
	aiSarcina    *sarcinaAI
	paza         *audit.Paza
}

func NewApp(db *sql.DB) *App {
	return &App{db: db}
}

func (a *App) pornire(ctx context.Context) {
	a.ctx = ctx
	go a.autoPornesteAIDacaEPosibil()
}

func (a *App) oprire(ctx context.Context) {
	if a.aiManager != nil {
		_ = a.aiManager.Opreste()
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

// Limbi întoarce codurile limbilor disponibile, cu limba implicită prima.
func (a *App) Limbi() []string {
	var rez []string
	for _, l := range i18n.Limbi() {
		rez = append(rez, string(l))
	}
	return rez
}
