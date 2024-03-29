package repl

import (
	"errors"
	"fmt"
	"strings"

	"github.com/holmes89/shrew/lexer"
	"github.com/holmes89/shrew/types"
	. "github.com/holmes89/shrew/types"

	. "github.com/holmes89/shrew/env"
)

// read
func read(str string) (Expression, error) {
	return lexer.Read(strings.NewReader(str))
}

// eval
func eval(ast Expression, env EnvType) (Expression, error) {
	var e error
	for {
		list, ok := ast.(List)

		if !ok {
			return eval_ast(ast, env)
		}
		// apply list
		ast, e = macroexpand(ast, env)
		if e != nil {
			return nil, e
		}

		list, ok = ast.(List)
		if !ok {
			return eval_ast(ast, env)
		}

		listLen := len(list.Val)
		if listLen == 0 {
			return ast, nil
		}

		a0 := list.Val[0]
		var a1 Expression = nil
		var a2 Expression = nil

		if listLen > 1 {
			a1 = list.Val[1]
		}

		if listLen > 2 {
			a2 = list.Val[2]
		}

		a0sym := "__<*fn*>__"
		if Symbol_Q(a0) {
			a0sym = a0.(Symbol).Val
		}

		switch a0sym {
		case "define":
			res, e := eval(a2, env)
			if e != nil {
				return nil, e
			}
			return env.Set(a1.(Symbol), res), nil
		case "let*":
			let_env, e := NewEnv(env, nil, nil)
			if e != nil {
				return nil, e
			}
			arr1, e := GetSlice(a1)
			if e != nil {
				return nil, e
			}
			for i := 0; i < len(arr1); i += 2 {
				if !Symbol_Q(arr1[i]) {
					return nil, errors.New("non-symbol bind value")
				}
				exp, e := eval(arr1[i+1], let_env)
				if e != nil {
					return nil, e
				}
				let_env.Set(arr1[i].(Symbol), exp)
			}
			ast = a2
			env = let_env
		case "do":
			el, e := eval_ast(List{
				Val: list.Val[1:],
			}, env)
			if e != nil {
				return nil, e
			}
			lst := el.(List).Val
			if len(lst) == 1 {
				return nil, nil
			}
			ast = lst[len(lst)-1]
		case "defmacro":
			fn, e := eval(a2, env)
			fn = fn.(ExpressionFunc).SetMacro()
			if e != nil {
				return nil, e
			}
			return env.Set(a1.(Symbol), fn), nil
		case "macroexpand":
			return macroexpand(a1, env)
		case "if":
			cond, e := eval(a1, env)
			if e != nil {
				return nil, e
			}
			if cond == nil || cond == false {
				if len(list.Val) >= 4 {
					ast = list.Val[3]
				} else {
					return nil, nil
				}
			} else {
				ast = a2
			}
		case "else":
			return eval(a1, env)
		case "cond":
			for _, c := range list.Val[1:] {
				cond, ok := c.(List)
				if !ok {
					return nil, errors.New("does not evaluate")
				}

				exp := cond.Val[0]
				var res Expression
				var e error
				if _, ok := exp.(List); ok {
					res, e = eval(exp, env)
					if e != nil {
						return nil, e
					}
					if res == true {
						return eval(cond.Val[1], env)
					}
					continue
				}
				//assume else?
				return eval(cond, env)
			}

			return nil, nil
		case "λ":
			fallthrough
		case "lambda":
			fn := ExpressionFunc{
				Eval:    eval,
				Exp:     a2,
				Env:     env,
				Params:  a1,
				IsMacro: false,
				GenEnv:  NewEnv,
				Meta:    nil,
			}

			return fn, nil
		case "quote":
			return a1, nil
		case "quasiquoteexpand":
			return quasiquote(a1), nil
		case "quasiquote":
			ast = quasiquote(a1)
		case "try":
			var exc Expression
			exp, e := eval(a1, env)
			if e == nil {
				return exp, nil
			} else {
				if a2 != nil && List_Q(a2) {
					a2s, _ := GetSlice(a2)
					if Symbol_Q(a2s[0]) && (a2s[0].(Symbol).Val == "catch*") {
						switch e.(type) {
						case ExpressionError:
							exc = e.(ExpressionError).Obj
						default:
							exc = e.Error()
						}
						binds := NewList(a2s[1])
						new_env, e := NewEnv(env, binds, NewList(exc))
						if e != nil {
							return nil, e
						}
						exp, e = eval(a2s[2], new_env)
						if e == nil {
							return exp, nil
						}
					}
				}
				return nil, e
			}
		default:
			el, e := eval_ast(ast, env)
			if e != nil {
				return nil, e
			}
			f := el.(List).Val[0]
			if ExpressionFunc_Q(f) {
				fn := f.(ExpressionFunc)
				ast = fn.Exp
				env, e = NewEnv(fn.Env, fn.Params, List{Val: el.(List).Val[1:]})
				if e != nil {
					return nil, e
				}
			} else {
				fn, ok := f.(Func)
				if !ok {
					return nil, fmt.Errorf("attempt to call non-function: %v", f)
				}
				return fn.Fn(el.(List).Val[1:])
			}
		}
	}

}

func eval_ast(ast Expression, env EnvType) (Expression, error) {
	switch {
	case Symbol_Q(ast):
		return env.Get(ast.(Symbol))
	case List_Q(ast):
		var lst []Expression
		l := ast.(List).Val
		for _, a := range l {
			exp, err := eval(a, env)
			if err != nil {
				return nil, err
			}
			lst = append(lst, exp)
		}
		return List{Val: lst}, nil
	case Vector_Q(ast):
		var lst []Expression
		l := ast.(Vector).Val
		for _, a := range l {
			exp, err := eval(a, env)
			if err != nil {
				return nil, err
			}
			lst = append(lst, exp)
		}
		return Vector{Val: lst}, nil
	case HashMap_Q(ast):
		m := ast.(HashMap)
		new_hm := HashMap{Val: map[Keyword]Expression{}}
		for k, v := range m.Val {
			ke, e1 := eval(k, env)
			if e1 != nil {
				return nil, e1
			}
			if _, ok := ke.(Keyword); !ok {
				return nil, errors.New("non Keyword hash-map key")
			}
			kv, e2 := eval(v, env)
			if e2 != nil {
				return nil, e2
			}
			new_hm.Val[ke.(Keyword)] = kv
		}
		return new_hm, nil
	default:
		return ast, nil
	}
}

func starts_with(xs []Expression, sym string) bool {
	if len(xs) <= 0 {
		return false
	}
	s, ok := xs[0].(Symbol)
	if !ok {
		return false
	}

	return s.Val == sym
}

func qq_loop(xs []Expression) Expression {
	acc := NewList()
	for i := len(xs) - 1; 0 <= i; i -= 1 {
		elt := xs[i]
		switch e := elt.(type) {
		case List:
			if starts_with(e.Val, "splice-unquote") {
				acc = NewList(Symbol{Val: "concat"}, e.Val[1], acc)
				continue
			}
		default:
		}
		acc = NewList(Symbol{Val: "cons"}, quasiquote(elt), acc)
	}
	return acc
}

func quasiquote(ast Expression) Expression {
	switch a := ast.(type) {
	case Vector:
		return NewList(Symbol{Val: "vec"}, qq_loop(a.Val))
	case HashMap, Symbol:
		return NewList(Symbol{Val: "quote"}, ast)
	case List:
		if starts_with(a.Val, "unquote") {
			return a.Val[1]
		}
		return qq_loop(a.Val)
	default:
		return ast
	}
}

func is_macro_call(ast Expression, env EnvType) bool {
	if List_Q(ast) {
		slc, _ := GetSlice(ast)
		if len(slc) == 0 {
			return false
		}
		a0 := slc[0]
		sym, ok := a0.(Symbol)
		if ok && env.Find(sym) != nil {
			exp, e := env.Get(sym)
			if e != nil {
				return false
			}
			if ExpressionFunc_Q(exp) {
				return exp.(ExpressionFunc).GetMacro()
			}
		}
	}
	return false
}

func macroexpand(ast Expression, env EnvType) (Expression, error) {
	var exp Expression
	var e error
	for is_macro_call(ast, env) {
		slc, _ := GetSlice(ast)
		a0 := slc[0]
		exp, e = env.Get(a0.(Symbol))
		if e != nil {
			return nil, e
		}
		fn := exp.(ExpressionFunc)
		ast, e = Apply(fn, slc[1:])
		if e != nil {
			return nil, e
		}
	}
	return ast, nil
}

// print
func print(exp Expression) (string, error) {
	if List_Q(exp) || Symbol_Q(exp) {
		return fmt.Sprintf("'%v", exp), nil
	}
	return fmt.Sprintf("%v", exp), nil
}

func Eval(exp Expression, repl_env types.EnvType) (Expression, error) {
	return eval(exp, repl_env)
}

func Repl(str string, repl_env EnvType) (Expression, error) {
	var exp Expression
	var res string
	var e error
	if exp, e = read(str); e != nil {
		return nil, e
	}
	if exp, e = eval(exp, repl_env); e != nil {
		return nil, e
	}
	if res, e = print(exp); e != nil {
		return nil, e
	}
	return res, nil
}
