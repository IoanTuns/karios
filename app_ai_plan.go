package main

import (
	"pontix/internal/ai"
	"pontix/internal/datebaza"
)

type RaspunsPlanAI struct {
	Plan   *datebaza.PlanAI `json:"plan,omitempty"`
	Eroare *MesajEroare     `json:"eroare,omitempty"`
}

func (a *App) PrevizualizeazaPlanAI(intentie ai.IntentiePlan) RaspunsPlanAI {
	plan, err := datebaza.Nou(a.db).PregatestePlanAI(a.contextLucru(), intentie)
	if err != nil {
		return RaspunsPlanAI{Eroare: eroareInterfata(err)}
	}
	return RaspunsPlanAI{Plan: &plan}
}

func (a *App) ConfirmaPlanAI(intentie ai.IntentiePlan, token string) RaspunsPlanAI {
	ctx := datebaza.CuAsistentAI(a.contextLucru())
	plan, err := datebaza.Nou(a.db).AplicaPlanAI(ctx, intentie, token)
	if err != nil {
		return RaspunsPlanAI{Eroare: eroareInterfata(err)}
	}
	return RaspunsPlanAI{Plan: &plan}
}
