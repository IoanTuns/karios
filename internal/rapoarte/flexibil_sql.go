package rapoarte

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/ncruces/go-sqlite3"
	"github.com/ncruces/go-sqlite3/driver"
	"strings"
	"time"
)

type bazaRaport struct {
	parametri map[string]any
	sarbatori map[string]bool
	db        *sql.DB
	schema    []TabelRaport
	variabile map[string]string
	// standard: documentele surselor standard folosite de SQL (D090).
	standard []documentStandard
}

func nouaBazaRaport() (*bazaRaport, error) {
	d, e := driver.Open(":memory:")
	if e != nil {
		return nil, e
	}
	d.SetMaxOpenConns(1)
	return &bazaRaport{db: d}, nil
}
func (b *bazaRaport) adauga(ctx context.Context, nume string, campuri []string, randuri [][]any) error {
	cols := []string{}
	marc := []string{}
	for _, c := range campuri {
		if !codFlexibil.MatchString(c) {
			return eroare("rapoarte.flex_invalid")
		}
		cols = append(cols, `"`+c+`"`)
		marc = append(marc, "?")
	}
	if _, e := b.db.ExecContext(ctx, "CREATE TABLE "+nume+" ("+strings.Join(cols, ",")+")"); e != nil {
		return e
	}
	b.schema = append(b.schema, TabelRaport{Cod: nume, Campuri: campuri})
	tx, e := b.db.BeginTx(ctx, nil)
	if e != nil {
		return e
	}
	defer tx.Rollback()
	s, e := tx.PrepareContext(ctx, "INSERT INTO "+nume+" VALUES("+strings.Join(marc, ",")+")")
	if e != nil {
		return e
	}
	defer s.Close()
	for _, r := range randuri {
		if _, e := s.ExecContext(ctx, r...); e != nil {
			return e
		}
	}
	return tx.Commit()
}
func (b *bazaRaport) query(ctx context.Context, text string, param map[string]any, scalar ...bool) (DateFlexibile, error) {
	rez := DateFlexibile{Schema: b.schema, Coloane: []string{}, Randuri: [][]any{}}
	if len(text) > 32768 || strings.TrimSpace(text) == "" {
		return rez, eroare("rapoarte.sql_invalid")
	}
	ctx, stop := context.WithTimeout(ctx, 2*time.Second)
	defer stop()
	conn, e := b.db.Conn(ctx)
	if e != nil {
		return rez, e
	}
	defer conn.Close()
	e = conn.Raw(func(raw any) error {
		c := raw.(driver.Conn).Raw()
		c.Limit(sqlite3.LIMIT_LENGTH, 256*1024)
		c.Limit(sqlite3.LIMIT_SQL_LENGTH, 32768)
		c.Limit(sqlite3.LIMIT_EXPR_DEPTH, 80)
		c.SetInterrupt(ctx)
		defer c.SetInterrupt(context.Background())
		if e := c.SetAuthorizer(func(a sqlite3.AuthorizerActionCode, n1, n2, sch, inner string) sqlite3.AuthorizerReturnCode {
			switch a {
			case sqlite3.AUTH_SELECT, sqlite3.AUTH_RECURSIVE:
				return sqlite3.AUTH_OK
			case sqlite3.AUTH_READ:
				for _, t := range b.schema {
					if t.Cod == n1 {
						return sqlite3.AUTH_OK
					}
				}
			case sqlite3.AUTH_FUNCTION:
				if strings.Contains("|sum|count|avg|min|max|round|abs|coalesce|ifnull|nullif|case|lower|upper|length|substr|substring|strftime|date|time|datetime|julianday|printf|cast|typeof|replace|trim|group_concat|total|iif|like|", "|"+strings.ToLower(n2)+"|") {
					return sqlite3.AUTH_OK
				}
			}
			return sqlite3.AUTH_DENY
		}); e != nil {
			return e
		}
		defer c.SetAuthorizer(nil)
		s, tail, e := c.Prepare(text)
		if e != nil || s == nil {
			return eroare("rapoarte.sql_invalid")
		}
		defer s.Close()
		if strings.TrimSpace(tail) != "" || !s.ReadOnly() || s.ColumnCount() > 80 {
			return eroare("rapoarte.sql_invalid")
		}
		for i := 1; i <= s.BindCount(); i++ {
			nom := strings.TrimLeft(s.BindName(i), ":@$?")
			v, ok := param[nom]
			if !ok {
				return eroare("rapoarte.sql_parametru")
			}
			switch v := v.(type) {
			case string:
				e = s.BindText(i, v)
			case int:
				e = s.BindInt64(i, int64(v))
			case int64:
				e = s.BindInt64(i, v)
			case float64:
				e = s.BindFloat(i, v)
			case nil:
				e = s.BindNull(i)
			default:
				return eroare("rapoarte.sql_parametru")
			}
			if e != nil {
				return e
			}
		}
		vaz := map[string]bool{}
		for i := 0; i < s.ColumnCount(); i++ {
			nom := s.ColumnName(i)
			if (len(scalar) == 0 || !scalar[0]) && (!codFlexibil.MatchString(nom) || vaz[strings.ToLower(nom)]) {
				return eroare("rapoarte.sql_coloane")
			}
			vaz[strings.ToLower(nom)] = true
			rez.Coloane = append(rez.Coloane, nom)
		}
		total := 0
		for s.Step() {
			if len(rez.Randuri) >= 10000 {
				return eroare("rapoarte.prea_mare")
			}
			r := []any{}
			for i := 0; i < s.ColumnCount(); i++ {
				var v any
				switch s.ColumnType(i) {
				case sqlite3.NULL:
					v = nil
				case sqlite3.INTEGER:
					v = s.ColumnInt64(i)
				case sqlite3.FLOAT:
					v = s.ColumnFloat(i)
				case sqlite3.TEXT:
					v = s.ColumnText(i)
				default:
					return eroare("rapoarte.sql_invalid")
				}
				total += len(fmt.Sprint(v))
				if total > 8*1024*1024 {
					return eroare("rapoarte.prea_mare")
				}
				r = append(r, v)
			}
			rez.Randuri = append(rez.Randuri, r)
		}
		if s.Err() != nil {
			return eroare("rapoarte.sql_invalid")
		}
		return nil
	})
	if e != nil {
		return rez, e
	}
	return rez, nil
}
