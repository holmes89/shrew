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
continueLoop:
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
		case "set!":
			if listLen < 3 {
				return nil, errors.New("set!: requires symbol and value")
			}
			sym, ok := a1.(Symbol)
			if !ok {
				return nil, errors.New("set!: first argument must be a symbol")
			}
			frame := env.Find(sym)
			if frame == nil {
				return nil, fmt.Errorf("set!: '%s' not found", sym.Val)
			}
			val, e := eval(a2, env)
			if e != nil {
				return nil, e
			}
			return frame.Set(sym, val), nil
		case "define":
			// Support (define (f args...) body) shorthand
			if List_Q(a1) {
				fnList, _ := GetSlice(a1)
				if len(fnList) == 0 {
					return nil, errors.New("define: empty function signature")
				}
				name := fnList[0].(Symbol)
				params := List{Val: fnList[1:]}
				body := list.Val[2:]
				var bodyExpr Expression
				if len(body) == 1 {
					bodyExpr = body[0]
				} else {
					bodyExpr = List{Val: append([]Expression{Symbol{Val: "begin"}}, body...)}
				}
				fn := ExpressionFunc{
					Eval:    eval,
					Exp:     bodyExpr,
					Env:     env,
					Params:  params,
					IsMacro: false,
					GenEnv:  NewEnv,
				}
				return env.Set(name, fn), nil
			}
			res, e := eval(a2, env)
			if e != nil {
				return nil, e
			}
			return env.Set(a1.(Symbol), res), nil
		case "begin":
			if listLen <= 1 {
				return nil, nil
			}
			for _, expr := range list.Val[1 : listLen-1] {
				if _, e := eval(expr, env); e != nil {
					return nil, e
				}
			}
			ast = list.Val[listLen-1]
		case "do":
			// R5RS iterative do: (do ((var init step) ...) (test result ...) body ...)
			if listLen < 3 {
				return nil, errors.New("do: requires bindings and test clause")
			}
			// Parse variable bindings
			bindList, e := GetSlice(a1)
			if e != nil {
				return nil, errors.New("do: first argument must be a list of bindings")
			}
			type doVar struct {
				sym  Symbol
				step Expression // nil means no step — var keeps its value
			}
			vars := make([]doVar, 0, len(bindList))
			do_env, e := NewEnv(env, nil, nil)
			if e != nil {
				return nil, e
			}
			for _, b := range bindList {
				parts, e := GetSlice(b)
				if e != nil || len(parts) < 2 {
					return nil, errors.New("do: malformed binding")
				}
				sym, ok := parts[0].(Symbol)
				if !ok {
					return nil, errors.New("do: binding name must be a symbol")
				}
				initVal, e := eval(parts[1], env)
				if e != nil {
					return nil, e
				}
				do_env.Set(sym, initVal)
				var step Expression
				if len(parts) >= 3 {
					step = parts[2]
				}
				vars = append(vars, doVar{sym, step})
			}
			// Parse test clause
			testClause, e := GetSlice(a2)
			if e != nil || len(testClause) == 0 {
				return nil, errors.New("do: test clause must be a non-empty list")
			}
			testExpr := testClause[0]
			resultExprs := testClause[1:]
			// Iterate
			for {
				testVal, e := eval(testExpr, do_env)
				if e != nil {
					return nil, e
				}
				if testVal != nil && testVal != false {
					// Test passed — evaluate result expressions
					if len(resultExprs) == 0 {
						return nil, nil
					}
					for _, r := range resultExprs[:len(resultExprs)-1] {
						if _, e := eval(r, do_env); e != nil {
							return nil, e
						}
					}
					return eval(resultExprs[len(resultExprs)-1], do_env)
				}
				// Evaluate body forms for side effects
				for _, body := range list.Val[3:] {
					if _, e := eval(body, do_env); e != nil {
						return nil, e
					}
				}
				// Compute all step values simultaneously, then assign
				newVals := make([]Expression, len(vars))
				for i, v := range vars {
					if v.step != nil {
						newVals[i], e = eval(v.step, do_env)
						if e != nil {
							return nil, e
						}
					} else {
						newVals[i], e = do_env.Get(v.sym)
						if e != nil {
							return nil, e
						}
					}
				}
				for i, v := range vars {
					do_env.Set(v.sym, newVals[i])
				}
			}
		case "let":
			// Support named let: (let name ((var init) ...) body...)
			if Symbol_Q(a1) {
				name := a1.(Symbol)
				bindings, e := GetSlice(a2)
				if e != nil {
					return nil, e
				}
				let_env, e := NewEnv(env, nil, nil)
				if e != nil {
					return nil, e
				}
				params := []Expression{}
				vals := []Expression{}
				for i := 0; i < len(bindings); i++ {
					pair, e := GetSlice(bindings[i])
					if e != nil || len(pair) != 2 {
						return nil, errors.New("named let: malformed binding")
					}
					params = append(params, pair[0])
					v, e := eval(pair[1], env)
					if e != nil {
						return nil, e
					}
					vals = append(vals, v)
				}
				body := list.Val[3:]
				var bodyExpr Expression
				if len(body) == 1 {
					bodyExpr = body[0]
				} else {
					bodyExpr = List{Val: append([]Expression{Symbol{Val: "begin"}}, body...)}
				}
				fn := ExpressionFunc{
					Eval:    eval,
					Exp:     bodyExpr,
					Env:     let_env,
					Params:  List{Val: params},
					IsMacro: false,
					GenEnv:  NewEnv,
				}
				let_env.Set(name, fn)
				ast, e = NewEnv(let_env, List{Val: params}, List{Val: vals})
				if e != nil {
					return nil, e
				}
				// Tail call: invoke fn with initial vals
				env, e = NewEnv(let_env, List{Val: params}, List{Val: vals})
				if e != nil {
					return nil, e
				}
				ast = bodyExpr
				continue
			}
			// Regular let: (let ((var val) ...) body...)
			let_env, e := NewEnv(env, nil, nil)
			if e != nil {
				return nil, e
			}
			bindings, e := GetSlice(a1)
			if e != nil {
				return nil, e
			}
			for _, b := range bindings {
				pair, e := GetSlice(b)
				if e != nil || len(pair) != 2 {
					return nil, errors.New("let: malformed binding")
				}
				val, e := eval(pair[1], env)
				if e != nil {
					return nil, e
				}
				let_env.Set(pair[0].(Symbol), val)
			}
			for _, expr := range list.Val[2 : listLen-1] {
				if _, e := eval(expr, let_env); e != nil {
					return nil, e
				}
			}
			ast = list.Val[listLen-1]
			env = let_env
		case "let*":
			let_env, e := NewEnv(env, nil, nil)
			if e != nil {
				return nil, e
			}
			bindings, e := GetSlice(a1)
			if e != nil {
				return nil, e
			}
			for _, b := range bindings {
				pair, e := GetSlice(b)
				if e != nil || len(pair) != 2 {
					return nil, errors.New("let*: malformed binding")
				}
				sym, ok := pair[0].(Symbol)
				if !ok {
					return nil, errors.New("let*: binding name must be a symbol")
				}
				val, e := eval(pair[1], let_env)
				if e != nil {
					return nil, e
				}
				let_env.Set(sym, val)
			}
			for _, expr := range list.Val[2 : listLen-1] {
				if _, e := eval(expr, let_env); e != nil {
					return nil, e
				}
			}
			ast = list.Val[listLen-1]
			env = let_env
		case "letrec":
			fallthrough
		case "letrec*":
			// Bind all names to nil first, then evaluate and set
			let_env, e := NewEnv(env, nil, nil)
			if e != nil {
				return nil, e
			}
			bindings, e := GetSlice(a1)
			if e != nil {
				return nil, e
			}
			names := make([]Symbol, 0, len(bindings))
			inits := make([]Expression, 0, len(bindings))
			for _, b := range bindings {
				pair, e := GetSlice(b)
				if e != nil || len(pair) != 2 {
					return nil, errors.New("letrec: malformed binding")
				}
				sym, ok := pair[0].(Symbol)
				if !ok {
					return nil, errors.New("letrec: binding name must be a symbol")
				}
				let_env.Set(sym, nil)
				names = append(names, sym)
				inits = append(inits, pair[1])
			}
			for i, init := range inits {
				val, e := eval(init, let_env)
				if e != nil {
					return nil, e
				}
				let_env.Set(names[i], val)
			}
			for _, expr := range list.Val[2 : listLen-1] {
				if _, e := eval(expr, let_env); e != nil {
					return nil, e
				}
			}
			ast = list.Val[listLen-1]
			env = let_env
		case "when":
			cond, e := eval(a1, env)
			if e != nil {
				return nil, e
			}
			if cond == nil || cond == false {
				return nil, nil
			}
			for _, expr := range list.Val[2 : listLen-1] {
				if _, e := eval(expr, env); e != nil {
					return nil, e
				}
			}
			ast = list.Val[listLen-1]
		case "unless":
			cond, e := eval(a1, env)
			if e != nil {
				return nil, e
			}
			if cond != nil && cond != false {
				return nil, nil
			}
			for _, expr := range list.Val[2 : listLen-1] {
				if _, e := eval(expr, env); e != nil {
					return nil, e
				}
			}
			ast = list.Val[listLen-1]
		case "and":
			var res Expression = true
			for _, expr := range list.Val[1:] {
				v, e := eval(expr, env)
				if e != nil {
					return nil, e
				}
				if v == nil || v == false {
					return false, nil
				}
				res = v
			}
			return res, nil
		case "or":
			for _, expr := range list.Val[1:] {
				v, e := eval(expr, env)
				if e != nil {
					return nil, e
				}
				if v != nil && v != false {
					return v, nil
				}
			}
			return false, nil
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
		case "cond":
			for _, c := range list.Val[1:] {
				clause, ok := c.(List)
				if !ok {
					return nil, errors.New("cond: clause must be a list")
				}
				if len(clause.Val) == 0 {
					return nil, errors.New("cond: empty clause")
				}
				test := clause.Val[0]
				// else clause
				if sym, ok := test.(Symbol); ok && sym.Val == "else" {
					if len(clause.Val) == 1 {
						return nil, nil
					}
					for _, expr := range clause.Val[1 : len(clause.Val)-1] {
						if _, e := eval(expr, env); e != nil {
							return nil, e
						}
					}
					ast = clause.Val[len(clause.Val)-1]
					goto continueLoop
				}
				res, e := eval(test, env)
				if e != nil {
					return nil, e
				}
				if res != nil && res != false {
					if len(clause.Val) == 1 {
						return res, nil
					}
					for _, expr := range clause.Val[1 : len(clause.Val)-1] {
						if _, e := eval(expr, env); e != nil {
							return nil, e
						}
					}
					ast = clause.Val[len(clause.Val)-1]
					goto continueLoop
				}
			}
			return nil, nil
		case "λ":
			fallthrough
		case "lambda":
			var bodyExpr Expression
			if listLen == 3 {
				bodyExpr = a2
			} else if listLen > 3 {
				bodyExpr = List{Val: append([]Expression{Symbol{Val: "begin"}}, list.Val[2:]...)}
			} else {
				return nil, errors.New("lambda: missing body")
			}
			fn := ExpressionFunc{
				Eval:    eval,
				Exp:     bodyExpr,
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
	if mv, ok := exp.(MultipleValues); ok {
		parts := make([]string, len(mv.Vals))
		for i, v := range mv.Vals {
			s, e := print(v)
			if e != nil {
				return "", e
			}
			parts[i] = s
		}
		return strings.Join(parts, "\n"), nil
	}
	if r, ok := exp.(rune); ok {
		return string(r), nil
	}
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
