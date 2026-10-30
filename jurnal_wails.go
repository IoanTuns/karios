package main

import (
	"encoding/json"
	"errors"
	"strings"

	"pontix/internal/diagnoza"
)

// jurnalWails primește jurnalul Wails. Wails prinde panicile din metodele legate ale App, dar le
// scrie doar în jurnalul lui, împreună cu argumentele apelului, adică datele introduse de
// utilizator (SEC-8). Aici, erorile și panicile ajung în jurnalul de diagnoză cu numele metodei și
// fără argumente; restul mesajelor Wails se ignoră.
type jurnalWails struct{}

func (jurnalWails) Print(string)   {}
func (jurnalWails) Trace(string)   {}
func (jurnalWails) Debug(string)   {}
func (jurnalWails) Info(string)    {}
func (jurnalWails) Warning(string) {}
func (jurnalWails) Error(mesaj string) {
	sursa, text := descrieMesajWails(mesaj)
	diagnoza.ObtineManager().InregistreazaEroare(sursa, errors.New(text))
}
func (jurnalWails) Fatal(mesaj string) {
	sursa, text := descrieMesajWails(mesaj)
	diagnoza.ObtineManager().InregistreazaEroare(sursa, errors.New("FATAL: "+text))
}

// descrieMesajWails desparte „process message error: C{"name":"main.App.X","args":[…]} -> eroare”
// în metoda apelată și eroare; argumentele nu se păstrează.
func descrieMesajWails(mesaj string) (sursa, text string) {
	const prefix = "process message error: "
	if !strings.HasPrefix(mesaj, prefix) {
		return "wails", mesaj
	}
	apel, eroare, gasit := strings.Cut(strings.TrimPrefix(mesaj, prefix), " -> ")
	if !gasit {
		return "wails", "eroare la procesarea unui mesaj"
	}
	var c struct {
		Name string `json:"name"`
	}
	if strings.HasPrefix(apel, "C") && json.Unmarshal([]byte(apel[1:]), &c) == nil && c.Name != "" {
		return "wails/" + c.Name, eroare
	}
	return "wails", eroare
}
