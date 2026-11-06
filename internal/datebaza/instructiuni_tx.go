package datebaza

import (
	"context"
	"database/sql"
)

// txPregatit refoloseste instructiunile pregatite in aceeasi tranzactie (I066, F06-T13).
// Salvarea unei luni intregi executa aceleasi INSERT-uri de mii de ori; fara cache,
// fiecare apel pregatea din nou textul SQL, iar parsarea ocupa o parte mare din timp.
//
// Numai ExecContext si QueryRowContext folosesc cache-ul: rezultatul lor se consuma
// imediat. QueryContext ramane direct, ca o instructiune sa nu fie refolosita cat
// randurile unui apel anterior sunt inca deschise. Instructiunile pregatite pe tranzactie
// se inchid automat la Commit sau Rollback.
type txPregatit struct {
	tx    *sql.Tx
	cache map[string]*sql.Stmt
}

func nouTxPregatit(tx *sql.Tx) *txPregatit {
	return &txPregatit{tx: tx, cache: map[string]*sql.Stmt{}}
}

func (t *txPregatit) instructiune(ctx context.Context, query string) (*sql.Stmt, error) {
	if s, ok := t.cache[query]; ok {
		return s, nil
	}
	s, err := t.tx.PrepareContext(ctx, query)
	if err != nil {
		return nil, err
	}
	t.cache[query] = s
	return s, nil
}

func (t *txPregatit) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	s, err := t.instructiune(ctx, query)
	if err != nil {
		return nil, err
	}
	return s.ExecContext(ctx, args...)
}

func (t *txPregatit) PrepareContext(ctx context.Context, query string) (*sql.Stmt, error) {
	return t.tx.PrepareContext(ctx, query)
}

func (t *txPregatit) QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	return t.tx.QueryContext(ctx, query, args...)
}

func (t *txPregatit) QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row {
	s, err := t.instructiune(ctx, query)
	if err != nil {
		// Eroarea de pregatire ajunge la Scan, ca la apelul direct.
		return t.tx.QueryRowContext(ctx, query, args...)
	}
	return s.QueryRowContext(ctx, args...)
}
