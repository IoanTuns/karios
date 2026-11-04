package acces

import "context"

// ActualizeazaSectiunePozitie foloseste revizia deja crescuta in aceeasi tranzactie.
func (q *Queries) ActualizeazaSectiunePozitie(ctx context.Context, id, sectiune string, revizie int64) error {
	_, err := q.db.ExecContext(ctx, `UPDATE planificare SET sectiune=? WHERE id=? AND revizie=?`, sectiune, id, revizie)
	return err
}
