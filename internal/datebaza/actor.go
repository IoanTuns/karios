package datebaza

import "context"

// ActorLocal e autorul implicit din jurnal până la conturi (F15).
const ActorLocal = "local"

// SufixAsistentAI marchează în jurnal modificările aplicate de asistentul AI după
// confirmarea utilizatorului (F14-T11): „local/asistent_ai”.
const SufixAsistentAI = "/asistent_ai"

type cheieActor struct{}

// CuActor întoarce un context în care modificările se notează în jurnal cu autorul dat.
func CuActor(ctx context.Context, actor string) context.Context {
	return context.WithValue(ctx, cheieActor{}, actor)
}

// CuAsistentAI marchează modificările ca aplicate de asistent, confirmate de autorul curent.
func CuAsistentAI(ctx context.Context) context.Context {
	return CuActor(ctx, Actor(ctx)+SufixAsistentAI)
}

// Actor întoarce autorul din context sau ActorLocal.
func Actor(ctx context.Context) string {
	if a, ok := ctx.Value(cheieActor{}).(string); ok && a != "" {
		return a
	}
	return ActorLocal
}
