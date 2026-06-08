package core

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/holmes89/shrew/lexer"
	. "github.com/holmes89/shrew/types"
)

var NS = map[Symbol]func(a []Expression) (Expression, error){
	makeSymbol("+"):           add,
	makeSymbol("-"):           sub,
	makeSymbol("*"):           mul,
	makeSymbol("/"):           div,
	makeSymbol("^"):           expt,
	makeSymbol("expt"):        expt,
	makeSymbol("and"):         and,
	makeSymbol("or"):          or,
	makeSymbol("not"):         not,
	makeSymbol("apply"):       apply,
	makeSymbol("car"):         first,
	makeSymbol("cdr"):         rest,
	makeSymbol("cons"):        cons,
	makeSymbol("="):           equal,
	makeSymbol("eq?"):         equal,
	makeSymbol(">"):           gt,
	makeSymbol(">="):          gte,
	makeSymbol("<"):           lt,
	makeSymbol("<="):          lte,
	makeSymbol("null?"):       null,
	makeSymbol("atom?"):       atom,
	makeSymbol("number?"):     number,
	makeSymbol("pair?"):       pair,
	makeSymbol("or?"):         or,
	makeSymbol("str"):         str,
	makeSymbol("slurp"):       slurp,
	makeSymbol("read-string"): read_string,
	makeSymbol("count"):       count,
	makeSymbol("zero?"):       zero,
	makeSymbol("even?"):       even,
	makeSymbol("length"):      length,
	makeSymbol("list"):        func(a []Expression) (Expression, error) { return List{Val: a}, nil },
}

func makeSymbol(text string) Symbol {
	return Symbol{
		Val: text,
	}
}

// Errors/Exceptions
// func throw(a []Expression) (Expression, error) {
// 	return nil, MalError{a[0]}
// }

func fn_q(a []Expression) (Expression, error) {
	switch f := a[0].(type) {
	case ExpressionFunc:
		return !f.GetMacro(), nil
	case Func:
		return true, nil
	case func([]Expression) (Expression, error):
		return true, nil
	default:
		return false, nil
	}
}

func read_string(a []Expression) (Expression, error) {
	if e := assertArgNum(a, 1); e != nil {
		return nil, e
	}
	return lexer.Read(strings.NewReader(a[0].(string)))
}

// Logic

func and(a []Expression) (Expression, error) {
	if e := assertArgNum(a, 2); e != nil {
		return nil, e
	}
	return a[0].(bool) && a[1].(bool), nil
}

func or(a []Expression) (Expression, error) {
	if e := assertArgNum(a, 2); e != nil {
		return nil, e
	}
	return a[0].(bool) || a[1].(bool), nil
}

func not(a []Expression) (Expression, error) {
	if e := assertArgNum(a, 1); e != nil {
		return nil, e
	}
	b, ok := a[0].(bool)
	if !ok {
		return false, nil
	}
	return !b, nil
}

func equal(a []Expression) (Expression, error) {
	if e := assertArgNum(a, 2); e != nil {
		return nil, e
	}
	return Equal_Q(a[0], a[1]), nil
}

// Math
func add(a []Expression) (Expression, error) {
	var res int
	for _, e := range a {
		n, ok := e.(int)
		if !ok {
			return nil, errors.New("expected number")
		}
		res += n
	}
	return res, nil
}

func sub(a []Expression) (Expression, error) {
	if len(a) == 0 {
		return nil, errors.New("arity mismatch")
	}
	if len(a) == 1 {
		return -1 * a[0].(int), nil
	}
	res := a[0].(int)
	for _, e := range a[1:] {
		n, ok := e.(int)
		if !ok {
			return nil, errors.New("expected number")
		}
		res -= n
	}
	return res, nil
}

func mul(a []Expression) (Expression, error) {
	res := 1
	for _, e := range a {
		n, ok := e.(int)
		if !ok {
			return nil, errors.New("expected number")
		}
		res *= n
	}
	return res, nil
}

func div(a []Expression) (Expression, error) {
	res := 1
	for _, e := range a {
		n, ok := e.(int)
		if !ok {
			return nil, errors.New("expected number")
		}
		if n == 0 {
			return nil, errors.New("divide by zero")
		}
		res /= n
	}
	return res, nil
}

// TODO this is messy
func expt(a []Expression) (Expression, error) {
	if e := assertArgNum(a, 2); e != nil {
		return nil, e
	}
	b := a[0].(int)
	s := b
	pow := a[1].(int)
	for p := pow - 1; p > 0; p-- {
		s *= b
	}
	return s, nil
}

func assertArgNum(a []Expression, n int) error {
	if len(a) != n {
		return errors.New("wrong number of arguments")
	}
	return nil
}

func gt(a []Expression) (Expression, error) {
	if e := assertArgNum(a, 2); e != nil {
		return nil, e
	}
	return a[0].(int) > a[1].(int), nil
}

func gte(a []Expression) (Expression, error) {
	if e := assertArgNum(a, 2); e != nil {
		return nil, e
	}
	return a[0].(int) >= a[1].(int), nil
}

func lt(a []Expression) (Expression, error) {
	if e := assertArgNum(a, 2); e != nil {
		return nil, e
	}
	return a[0].(int) < a[1].(int), nil
}

func lte(a []Expression) (Expression, error) {
	if e := assertArgNum(a, 2); e != nil {
		return nil, e
	}
	return a[0].(int) <= a[1].(int), nil
}

// String functions

// func pr_str(a []Expression) (Expression, error) {
// 	return printer.Pr_list(a, true, "", "", " "), nil
// }

func str(a []Expression) (Expression, error) {
	sarray := []string{}
	for _, b := range a {
		sarray = append(sarray, fmt.Sprintf("%+v", b))
	}
	return strings.Join(sarray, ""), nil
}

// func prn(a []Expression) (Expression, error) {
// 	fmt.Println(printer.Pr_list(a, true, "", "", " "))
// 	return nil, nil
// }

// func println(a []Expression) (Expression, error) {
// 	fmt.Println(printer.Pr_list(a, false, "", "", " "))
// 	return nil, nil
// }

func slurp(a []Expression) (Expression, error) {
	b, e := os.ReadFile(a[0].(string))
	if e != nil {
		return nil, e
	}
	return string(b), nil
}

// Number functions
func time_ms(a []Expression) (Expression, error) {
	return int(time.Now().UnixNano() / int64(time.Millisecond)), nil
}

// Hash Map functions
func copy_hash_map(hm HashMap) HashMap {
	new_hm := HashMap{Val: map[Keyword]Expression{}}
	for k, v := range hm.Val {
		new_hm.Val[k] = v
	}
	return new_hm
}

func assoc(a []Expression) (Expression, error) {
	if len(a) < 3 {
		return nil, errors.New("assoc requires at least 3 arguments")
	}
	if len(a)%2 != 1 {
		return nil, errors.New("assoc requires odd number of arguments")
	}
	if !HashMap_Q(a[0]) {
		return nil, errors.New("assoc called on non-hash map")
	}
	new_hm := copy_hash_map(a[0].(HashMap))
	for i := 1; i < len(a); i += 2 {
		key := a[i]
		if !String_Q(key) {
			return nil, errors.New("assoc called with non-string key")
		}
		new_hm.Val[key.(Keyword)] = a[i+1]
	}
	return new_hm, nil
}

func dissoc(a []Expression) (Expression, error) {
	if len(a) < 2 {
		return nil, errors.New("dissoc requires at least 3 arguments")
	}
	if !HashMap_Q(a[0]) {
		return nil, errors.New("dissoc called on non-hash map")
	}
	new_hm := copy_hash_map(a[0].(HashMap))
	for i := 1; i < len(a); i += 1 {
		key := a[i]
		if !String_Q(key) {
			return nil, errors.New("dissoc called with non-string key")
		}
		delete(new_hm.Val, key.(Keyword))
	}
	return new_hm, nil
}

func get(a []Expression) (Expression, error) {
	if Nil_Q(a[0]) {
		return nil, nil
	}
	if !HashMap_Q(a[0]) {
		return nil, errors.New("get called on non-hash map")
	}
	if !String_Q(a[1]) {
		return nil, errors.New("get called with non-string key")
	}
	return a[0].(HashMap).Val[a[1].(Keyword)], nil
}

func contains_Q(hm Expression, key Expression) (Expression, error) {
	if Nil_Q(hm) {
		return false, nil
	}
	if !HashMap_Q(hm) {
		return nil, errors.New("get called on non-hash map")
	}
	if !String_Q(key) {
		return nil, errors.New("get called with non-string key")
	}
	_, ok := hm.(HashMap).Val[key.(Keyword)]
	return ok, nil
}

func keys(a []Expression) (Expression, error) {
	if !HashMap_Q(a[0]) {
		return nil, errors.New("keys called on non-hash map")
	}
	slc := []Expression{}
	for k := range a[0].(HashMap).Val {
		slc = append(slc, k)
	}
	return List{Val: slc}, nil
}

func vals(a []Expression) (Expression, error) {
	if !HashMap_Q(a[0]) {
		return nil, errors.New("keys called on non-hash map")
	}
	slc := []Expression{}
	for _, v := range a[0].(HashMap).Val {
		slc = append(slc, v)
	}
	return List{Val: slc}, nil
}

// Sequence functions

func cons(a []Expression) (Expression, error) {
	val := a[0]
	lst, e := GetSlice(a[1])
	if e != nil {
		return nil, e
	}
	return List{Val: append([]Expression{val}, lst...)}, nil
}

func concat(a []Expression) (Expression, error) {
	if len(a) == 0 {
		return List{}, nil
	}
	slc1, e := GetSlice(a[0])
	if e != nil {
		return nil, e
	}
	for i := 1; i < len(a); i += 1 {
		slc2, e := GetSlice(a[i])
		if e != nil {
			return nil, e
		}
		slc1 = append(slc1, slc2...)
	}
	return List{Val: slc1}, nil
}

func vec(a []Expression) (Expression, error) {
	switch obj := a[0].(type) {
	case Vector:
		return obj, nil
	case List:
		return Vector{Val: obj.Val}, nil
	default:
		return nil, errors.New("vec: expects a sequence")
	}
}

func nth(a []Expression) (Expression, error) {
	slc, e := GetSlice(a[0])
	if e != nil {
		return nil, e
	}
	idx := a[1].(int)
	if idx < len(slc) {
		return slc[idx], nil
	} else {
		return nil, errors.New("nth: index out of range")
	}
}

func first(a []Expression) (Expression, error) {
	if len(a) == 0 {
		return nil, nil
	}
	if a[0] == nil {
		return nil, nil
	}
	slc, e := GetSlice(a[0])
	if e != nil {
		return nil, e
	}
	if len(slc) == 0 {
		return nil, nil
	}
	return slc[0], nil
}

func rest(a []Expression) (Expression, error) {
	if a[0] == nil {
		return List{}, nil
	}
	slc, e := GetSlice(a[0])
	if e != nil {
		return nil, e
	}
	if len(slc) == 0 {
		return List{}, nil
	}
	return List{Val: slc[1:]}, nil
}

func empty_Q(a []Expression) (Expression, error) {
	switch obj := a[0].(type) {
	case List:
		return len(obj.Val) == 0, nil
	case Vector:
		return len(obj.Val) == 0, nil
	case nil:
		return true, nil
	default:
		return nil, errors.New("empty? called on non-sequence")
	}
}

func count(a []Expression) (Expression, error) {
	switch obj := a[0].(type) {
	case List:
		return len(obj.Val), nil
	case Vector:
		return len(obj.Val), nil
	case map[string]Expression:
		return len(obj), nil
	case nil:
		return 0, nil
	default:
		return nil, errors.New("count called on non-sequence")
	}
}

func apply(a []Expression) (Expression, error) {
	if len(a) < 2 {
		return nil, errors.New("apply requires at least 2 args")
	}
	f := a[0]
	args := []Expression{}
	for _, b := range a[1 : len(a)-1] {
		args = append(args, b)
	}
	last, e := GetSlice(a[len(a)-1])
	if e != nil {
		return nil, e
	}
	args = append(args, last...)
	return Apply(f, args)
}

func do_map(a []Expression) (Expression, error) {
	f := a[0]
	results := []Expression{}
	args, e := GetSlice(a[1])
	if e != nil {
		return nil, e
	}
	for _, arg := range args {
		res, e := Apply(f, []Expression{arg})
		results = append(results, res)
		if e != nil {
			return nil, e
		}
	}
	return List{Val: results}, nil
}

func conj(a []Expression) (Expression, error) {
	if len(a) < 2 {
		return nil, errors.New("conj requires at least 2 arguments")
	}
	switch seq := a[0].(type) {
	case List:
		new_slc := []Expression{}
		for i := len(a) - 1; i > 0; i -= 1 {
			new_slc = append(new_slc, a[i])
		}
		return List{Val: append(new_slc, seq.Val...)}, nil
	case Vector:
		new_slc := seq.Val
		for _, x := range a[1:] {
			new_slc = append(new_slc, x)
		}
		return Vector{Val: new_slc}, nil
	}

	if !HashMap_Q(a[0]) {
		return nil, errors.New("dissoc called on non-hash map")
	}
	new_hm := copy_hash_map(a[0].(HashMap))
	for i := 1; i < len(a); i += 1 {
		key := a[i]
		if !String_Q(key) {
			return nil, errors.New("dissoc called with non-string key")
		}
		delete(new_hm.Val, key.(Keyword))
	}
	return new_hm, nil
}

func seq(a []Expression) (Expression, error) {
	if a[0] == nil {
		return nil, nil
	}
	switch arg := a[0].(type) {
	case List:
		if len(arg.Val) == 0 {
			return nil, nil
		}
		return arg, nil
	case Vector:
		if len(arg.Val) == 0 {
			return nil, nil
		}
		return List{Val: arg.Val}, nil
	case string:
		if len(arg) == 0 {
			return nil, nil
		}
		new_slc := []Expression{}
		for _, ch := range strings.Split(arg, "") {
			new_slc = append(new_slc, ch)
		}
		return List{Val: new_slc}, nil
	}
	return nil, errors.New("seq requires string or list or vector or nil")
}

// Metadata functions
func with_meta(a []Expression) (Expression, error) {
	obj := a[0]
	m := a[1]
	switch tobj := obj.(type) {
	case List:
		return List{Val: tobj.Val, Meta: m}, nil
	case Vector:
		return Vector{Val: tobj.Val, Meta: m}, nil
	case HashMap:
		return HashMap{Val: tobj.Val, Meta: m}, nil
	case Func:
		return Func{Fn: tobj.Fn, Meta: m}, nil
	case ExpressionFunc:
		fn := tobj
		fn.Meta = m
		return fn, nil
	default:
		return nil, errors.New("with-meta not supported on type")
	}
}

func meta(a []Expression) (Expression, error) {
	obj := a[0]
	switch tobj := obj.(type) {
	case List:
		return tobj.Meta, nil
	case Vector:
		return tobj.Meta, nil
	case HashMap:
		return tobj.Meta, nil
	case Func:
		return tobj.Meta, nil
	case ExpressionFunc:
		return tobj.Meta, nil
	default:
		return nil, errors.New("meta not supported on type")
	}
}

// Atom functions
func deref(a []Expression) (Expression, error) {
	if !Atom_Q(a[0]) {
		return nil, errors.New("deref called with non-atom")
	}
	return a[0].(*Atom).Val, nil
}

func reset_BANG(a []Expression) (Expression, error) {
	if !Atom_Q(a[0]) {
		return nil, errors.New("reset! called with non-atom")
	}
	a[0].(*Atom).Set(a[1])
	return a[1], nil
}

func swap_BANG(a []Expression) (Expression, error) {
	if !Atom_Q(a[0]) {
		return nil, errors.New("swap! called with non-atom")
	}
	atm := a[0].(*Atom)
	args := []Expression{atm.Val}
	f := a[1]
	args = append(args, a[2:]...)
	res, e := Apply(f, args)
	if e != nil {
		return nil, e
	}
	atm.Set(res)
	return res, nil
}

func null(a []Expression) (Expression, error) {
	if len(a) != 1 {
		return nil, fmt.Errorf("wrong number of arguments (%d instead of 1)", len(a))
	}
	list, ok := a[0].(List)
	if ok {
		return len(list.Val) == 0, nil
	}
	return Nil_Q(a[0]), nil
}

func atom(a []Expression) (Expression, error) {
	if len(a) != 1 {
		return nil, fmt.Errorf("wrong number of arguments (%d instead of 1)", len(a))
	}
	return !List_Q(a[0]), nil // HACK
}

func zero(a []Expression) (Expression, error) {
	if len(a) != 1 {
		return nil, fmt.Errorf("wrong number of arguments (%d instead of 1)", len(a))
	}
	return a[0] == 0, nil // HACK
}

func even(a []Expression) (Expression, error) {
	if len(a) != 1 {
		return nil, fmt.Errorf("wrong number of arguments (%d instead of 1)", len(a))
	}
	i, ok := a[0].(int)
	if !ok {
		return nil, errors.New("expected integer")
	}
	return i%2 == 0, nil
}

func length(a []Expression) (Expression, error) {
	if len(a) != 1 {
		return nil, fmt.Errorf("wrong number of arguments (%d instead of 1)", len(a))
	}
	l, ok := a[0].(List)
	if !ok {
		return nil, errors.New("list required")
	}
	return len(l.Val), nil // HACK
}

func pair(a []Expression) (Expression, error) {
	if len(a) != 1 {
		return nil, fmt.Errorf("wrong number of arguments (%d instead of 1)", len(a))
	}
	return List_Q(a[0]), nil
}

func number(a []Expression) (Expression, error) {
	if len(a) != 1 {
		return nil, fmt.Errorf("wrong number of arguments (%d instead of 1)", len(a))
	}
	return Number_Q(a[0]), nil
}
