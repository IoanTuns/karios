package ai

import "sync/atomic"

// UtilizareContext: tokenurile raportate de motor pentru ultima cerere trimisă
// modelului (I080, F14-T35). Se afișează numai ce raportează motorul în câmpul
// `usage`; fără el nu se estimează nimic.
type UtilizareContext struct {
	Prompt  int `json:"prompt"`
	Raspuns int `json:"raspuns"`
	Total   int `json:"total"`
	Limita  int `json:"limita"`
}

var ultimaUtilizare atomic.Pointer[UtilizareContext]

// UltimaUtilizare întoarce consumul ultimei cereri reușite sau nil dacă motorul
// nu l-a raportat încă.
func UltimaUtilizare() *UtilizareContext {
	if u := ultimaUtilizare.Load(); u != nil {
		copie := *u
		return &copie
	}
	return nil
}

// usageOpenAI este câmpul `usage` din răspunsul /v1/chat/completions.
type usageOpenAI struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

func inregistreazaUtilizare(u *usageOpenAI) {
	if u == nil || u.PromptTokens <= 0 {
		return
	}
	total := u.TotalTokens
	if total <= 0 {
		total = u.PromptTokens + u.CompletionTokens
	}
	ultimaUtilizare.Store(&UtilizareContext{Prompt: u.PromptTokens, Raspuns: u.CompletionTokens, Total: total, Limita: ContextConfigurat})
}
