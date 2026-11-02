package calcul

import (
	"go/ast"
	"go/parser"
	"go/token"
	"math"
	"strconv"
	"strings"
	"time"
)

// EvalueazaFormulaRaport propaga datele lipsa; nu executa cod Go.
func EvalueazaFormulaRaport(text string, valori map[string]*float64) (*float64, error) {
	if len(text) > 2048 {
		return nil, &Eroare{Cod: "rapoarte.formula_invalida"}
	}
	a, e := parser.ParseExpr(text)
	if e != nil {
		return nil, &Eroare{Cod: "rapoarte.formula_invalida"}
	}
	nr := 0
	var eval func(ast.Expr) (*float64, error)
	gresit := func() (*float64, error) { return nil, &Eroare{Cod: "rapoarte.formula_invalida"} }
	eval = func(a ast.Expr) (*float64, error) {
		nr++
		if nr > 256 {
			return gresit()
		}
		var v float64
		switch a := a.(type) {
		case *ast.BasicLit:
			if a.Kind != token.INT && a.Kind != token.FLOAT {
				return gresit()
			}
			var e error
			v, e = strconv.ParseFloat(a.Value, 64)
			if e != nil {
				return gresit()
			}
		case *ast.Ident:
			p, ok := valori[a.Name]
			if !ok {
				return gresit()
			}
			return p, nil
		case *ast.ParenExpr:
			return eval(a.X)
		case *ast.UnaryExpr:
			p, e := eval(a.X)
			if e != nil || p == nil {
				return p, e
			}
			v = *p
			if a.Op == token.SUB {
				v = -v
			} else if a.Op != token.ADD {
				return gresit()
			}
		case *ast.BinaryExpr:
			p, e := eval(a.X)
			if e != nil {
				return nil, e
			}
			q, e := eval(a.Y)
			if e != nil {
				return nil, e
			}
			if p == nil || q == nil {
				return nil, nil
			}
			switch a.Op {
			case token.ADD:
				v = *p + *q
			case token.SUB:
				v = *p - *q
			case token.MUL:
				v = *p * *q
			case token.QUO:
				if *q == 0 {
					return gresit()
				}
				v = *p / *q
			default:
				return gresit()
			}
		case *ast.CallExpr:
			n, ok := a.Fun.(*ast.Ident)
			if !ok {
				return gresit()
			}
			args := []float64{}
			for _, x := range a.Args {
				p, e := eval(x)
				if e != nil {
					return nil, e
				}
				if p == nil {
					return nil, nil
				}
				args = append(args, *p)
			}
			switch strings.ToUpper(n.Name) {
			case "ABS":
				if len(args) != 1 {
					return gresit()
				}
				v = math.Abs(args[0])
			case "MIN", "MAX":
				if len(args) != 2 {
					return gresit()
				}
				if strings.ToUpper(n.Name) == "MIN" {
					v = math.Min(args[0], args[1])
				} else {
					v = math.Max(args[0], args[1])
				}
			case "ROUND":
				if len(args) != 2 || args[1] < 0 || args[1] > 6 || math.Trunc(args[1]) != args[1] {
					return gresit()
				}
				f := math.Pow10(int(args[1]))
				v = math.Round(args[0]*f) / f
			default:
				return gresit()
			}
		default:
			return gresit()
		}
		if math.IsInf(v, 0) || math.IsNaN(v) {
			return gresit()
		}
		return &v, nil
	}
	return eval(a)
}

// MinuteGardaEligibile intersecteaza fiecare zi o singura data, apoi perioada.
func MinuteGardaEligibile(start, sfarsit, deLa, panaLa time.Time, weekend, sarbatori bool, libere map[string]bool) int64 {
	if start.Before(deLa) {
		start = deLa
	}
	if sfarsit.After(panaLa) {
		sfarsit = panaLa
	}
	var n int64
	for cur := start; cur.Before(sfarsit); {
		zi := time.Date(cur.Year(), cur.Month(), cur.Day(), 0, 0, 0, 0, cur.Location())
		lim := zi.AddDate(0, 0, 1)
		if lim.After(sfarsit) {
			lim = sfarsit
		}
		if weekend && (cur.Weekday() == time.Saturday || cur.Weekday() == time.Sunday) || sarbatori && libere[cur.Format(time.DateOnly)] {
			n += int64(lim.Sub(cur) / time.Minute)
		}
		cur = lim
	}
	return n
}
