package types

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
)

// Errors/Exceptions
type ExpressionError struct {
	Obj Expression
}

func (e ExpressionError) Error() string {
	return fmt.Sprintf("%#v", e.Obj)
}

// General Types
type Expression interface{}

type EnvType interface {
	Find(key Symbol) EnvType
	Set(key Symbol, value Expression) Expression
	Get(key Symbol) (Expression, error)
}

// Scalars
func Nil_Q(obj Expression) bool {
	return obj == nil
}

func True_Q(obj Expression) bool {
	b, ok := obj.(bool)
	return ok && b == true
}

func False_Q(obj Expression) bool {
	b, ok := obj.(bool)
	return ok && b == false
}

func Number_Q(obj Expression) bool {
	switch obj.(type) {
	case int, float64:
		return true
	}
	return false
}

func Float_Q(obj Expression) bool {
	_, ok := obj.(float64)
	return ok
}

func Int_Q(obj Expression) bool {
	_, ok := obj.(int)
	return ok
}

func Char_Q(obj Expression) bool {
	_, ok := obj.(rune)
	return ok
}

// ToFloat coerces an int or float64 expression to float64.
func ToFloat(obj Expression) (float64, bool) {
	switch v := obj.(type) {
	case float64:
		return v, true
	case int:
		return float64(v), true
	}
	return 0, false
}

// ToInt coerces a numeric expression to int (truncating floats).
func ToInt(obj Expression) (int, bool) {
	switch v := obj.(type) {
	case int:
		return v, true
	case float64:
		return int(v), true
	}
	return 0, false
}

// Symbols
type Symbol struct {
	Val string
}

func (s Symbol) String() string {
	return s.Val
}

func Symbol_Q(obj Expression) bool {
	_, ok := obj.(Symbol)
	return ok
}

// Keywords
type Keyword string

func NewKeyword(s string) (Expression, error) {
	return Keyword(s), nil
}

func (k Keyword) String() string {
	return fmt.Sprintf(":%s", string(k))
}

func Keyword_Q(obj Expression) bool {
	_, ok := obj.(Keyword)
	return ok
}

// Strings
func String_Q(obj Expression) bool {
	_, ok := obj.(string)
	return ok
}

// Functions
type Func struct {
	Fn   func([]Expression) (Expression, error)
	Meta Expression
}

func (e Func) Func() string {
	return fmt.Sprintf("<function %v>", e)
}

func Func_Q(obj Expression) bool {
	_, ok := obj.(Func)
	return ok
}

type ExpressionFunc struct {
	Eval    func(Expression, EnvType) (Expression, error)
	Exp     Expression
	Env     EnvType
	Params  Expression
	IsMacro bool
	GenEnv  func(EnvType, Expression, Expression) (EnvType, error)
	Meta    Expression
}

func (e ExpressionFunc) String() string {
	return fmt.Sprintf("(λ %v %v)", e.Params, e.Exp)
}

func ExpressionFunc_Q(obj Expression) bool {
	_, ok := obj.(ExpressionFunc)
	return ok
}

func (f ExpressionFunc) SetMacro() Expression {
	f.IsMacro = true
	return f
}

func (f ExpressionFunc) GetMacro() bool {
	return f.IsMacro
}

// Take either a ExpressionFunc or regular function and apply it to the
// arguments
func Apply(f_mt Expression, a []Expression) (Expression, error) {
	switch f := f_mt.(type) {
	case ExpressionFunc:
		env, e := f.GenEnv(f.Env, f.Params, List{a, nil})
		if e != nil {
			return nil, e
		}
		return f.Eval(f.Exp, env)
	case Func:
		return f.Fn(a)
	case func([]Expression) (Expression, error):
		return f(a)
	default:
		return nil, errors.New("Invalid function to Apply")
	}
}

// Lists
type List struct {
	Val  []Expression
	Meta Expression
}

func (l List) String() string {
	str := strings.Builder{}
	str.WriteString(`(`)
	for _, v := range l.Val {
		if str.Len() != 1 {
			str.WriteRune(' ')
		}
		str.WriteString(fmt.Sprintf("%v", v))
	}
	str.WriteRune(')')
	return str.String()
}

func NewList(a ...Expression) Expression {
	return List{a, nil}
}

func List_Q(obj Expression) bool {
	_, ok := obj.(List)
	return ok
}

// Vectors
type Vector struct {
	Val  []Expression
	Meta Expression
}

func (v Vector) String() string {
	str := strings.Builder{}
	str.WriteRune('[')
	for _, vi := range v.Val {
		if str.Len() != 1 {
			str.WriteRune(' ')
		}
		str.WriteString(fmt.Sprintf("%v", vi))
	}
	str.WriteRune(']')
	return str.String()
}

func Vector_Q(obj Expression) bool {
	_, ok := obj.(Vector)
	return ok
}

func GetSlice(seq Expression) ([]Expression, error) {
	switch obj := seq.(type) {
	case List:
		return obj.Val, nil
	case Vector:
		return obj.Val, nil
	default:
		return nil, errors.New("GetSlice called on non-sequence")
	}
}

// Hash Maps
type HashMap struct {
	Val  map[Keyword]Expression
	Meta Expression
}

func (h HashMap) String() string {
	str := strings.Builder{}
	str.WriteRune('{')
	for k, v := range h.Val {
		if str.Len() != 1 {
			str.WriteString(", ")
		}
		str.WriteString(fmt.Sprintf("%s %v", k, v))
	}
	str.WriteRune('}')
	return str.String()
}

func NewHashMap(seq Expression) (Expression, error) {
	lst, e := GetSlice(seq)
	if e != nil {
		return nil, e
	}
	if len(lst)%2 == 1 {
		return nil, errors.New("Odd number of arguments to NewHashMap")
	}
	m := map[Keyword]Expression{}
	for i := 0; i < len(lst); i += 2 {
		str, ok := lst[i].(Keyword)
		if !ok {
			return nil, errors.New("expected hash-map key string")
		}
		m[str] = lst[i+1]
	}
	return HashMap{m, nil}, nil
}

func HashMap_Q(obj Expression) bool {
	_, ok := obj.(HashMap)
	return ok
}

// func Quote(obj Expression) Expression {
// 	Symbol
// }

// Atoms
type Atom struct {
	Val  Expression
	Meta Expression
}

func (a Atom) String() string {
	return fmt.Sprintf("(atom %v)", a.Val)
}

func (a *Atom) Set(val Expression) Expression {
	a.Val = val
	return a
}

func Atom_Q(obj Expression) bool {
	_, ok := obj.(*Atom)
	return ok
}

// General functions
func _obj_type(obj Expression) string {
	if obj == nil {
		return "nil"
	}
	return reflect.TypeOf(obj).Name()
}

func Sequential_Q(seq Expression) bool {
	if seq == nil {
		return false
	}
	return (reflect.TypeOf(seq).Name() == "List") ||
		(reflect.TypeOf(seq).Name() == "Vector")
}

func Equal_Q(a Expression, b Expression) bool {
	ota := reflect.TypeOf(a)
	otb := reflect.TypeOf(b)
	if !((ota == otb) || (Sequential_Q(a) && Sequential_Q(b))) {
		return false
	}

	switch a.(type) {
	case Symbol:
		return a.(Symbol).Val == b.(Symbol).Val
	case List:
		as, _ := GetSlice(a)
		bs, _ := GetSlice(b)
		if len(as) != len(bs) {
			return false
		}
		for i := 0; i < len(as); i += 1 {
			if !Equal_Q(as[i], bs[i]) {
				return false
			}
		}
		return true
	case Vector:
		as, _ := GetSlice(a)
		bs, _ := GetSlice(b)
		if len(as) != len(bs) {
			return false
		}
		for i := 0; i < len(as); i += 1 {
			if !Equal_Q(as[i], bs[i]) {
				return false
			}
		}
		return true
	case HashMap:
		am := a.(HashMap).Val
		bm := b.(HashMap).Val
		if len(am) != len(bm) {
			return false
		}
		for k, v := range am {
			if !Equal_Q(v, bm[k]) {
				return false
			}
		}
		return true
	default:
		return a == b
	}
}
