package acces

import "context"

// PauzaProgram: cum se scade pauza unui program (F04-T05, D107).
type PauzaProgram struct {
	Mod string
	Ora *int64
}

// ListaPauzeProgram intoarce modul pauzei pe id-ul programului.
func (q *Queries) ListaPauzeProgram(ctx context.Context) (map[string]PauzaProgram, error) {
	rows, err := q.db.QueryContext(ctx, "SELECT id,mod_pauza,pauza_ora FROM tipuri_schimb")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	rez := map[string]PauzaProgram{}
	for rows.Next() {
		var id string
		var p PauzaProgram
		if err := rows.Scan(&id, &p.Mod, &p.Ora); err != nil {
			return nil, err
		}
		rez[id] = p
	}
	return rez, rows.Err()
}

// SeteazaPauzaProgram scrie modul pauzei unui program.
func (q *Queries) SeteazaPauzaProgram(ctx context.Context, id, mod string, ora *int64) error {
	_, err := q.db.ExecContext(ctx, "UPDATE tipuri_schimb SET mod_pauza=?, pauza_ora=? WHERE id=?", mod, ora, id)
	return err
}
