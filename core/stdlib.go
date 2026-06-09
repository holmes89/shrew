package core

import (
	"errors"
	"fmt"
	"math"
	"os"
	"sort"
	"strconv"
	"strings"
	"text/scanner"

	"github.com/holmes89/shrew/lexer"
	. "github.com/holmes89/shrew/types"
)

func init() {
	// Math
	NS[makeSymbol("+")] = addNum
	NS[makeSymbol("-")] = subNum
	NS[makeSymbol("*")] = mulNum
	NS[makeSymbol("/")] = divNum
	NS[makeSymbol("abs")] = absNum
	NS[makeSymbol("max")] = maxNum
	NS[makeSymbol("min")] = minNum
	NS[makeSymbol("modulo")] = modulo
	NS[makeSymbol("remainder")] = remainder
	NS[makeSymbol("quotient")] = quotient
	NS[makeSymbol("floor")] = floorFn
	NS[makeSymbol("ceiling")] = ceilingFn
	NS[makeSymbol("round")] = roundFn
	NS[makeSymbol("truncate")] = truncateFn
	NS[makeSymbol("sqrt")] = sqrtFn
	NS[makeSymbol("exact->inexact")] = exactToInexact
	NS[makeSymbol("inexact->exact")] = inexactToExact
	NS[makeSymbol("number->string")] = numberToString
	NS[makeSymbol("string->number")] = stringToNumber
	NS[makeSymbol("integer?")] = integerQ
	NS[makeSymbol("float?")] = floatQ
	NS[makeSymbol("exact?")] = integerQ // ints are exact
	NS[makeSymbol("inexact?")] = floatQ // floats are inexact

	// Boolean
	NS[makeSymbol("boolean?")] = booleanQ

	// Characters
	NS[makeSymbol("char?")] = charQ
	NS[makeSymbol("char->integer")] = charToInteger
	NS[makeSymbol("integer->char")] = integerToChar
	NS[makeSymbol("char=?")] = charEq
	NS[makeSymbol("char<?")] = charLt
	NS[makeSymbol("char>?")] = charGt
	NS[makeSymbol("char-alphabetic?")] = charAlphabetic
	NS[makeSymbol("char-numeric?")] = charNumeric
	NS[makeSymbol("char-whitespace?")] = charWhitespace
	NS[makeSymbol("char-upcase")] = charUpcase
	NS[makeSymbol("char-downcase")] = charDowncase

	// Strings
	NS[makeSymbol("string?")] = stringQ
	NS[makeSymbol("string=?")] = stringEq
	NS[makeSymbol("string<?")] = stringLt
	NS[makeSymbol("string>?")] = stringGt
	NS[makeSymbol("string-length")] = stringLength
	NS[makeSymbol("string-append")] = stringAppend
	NS[makeSymbol("substring")] = substring
	NS[makeSymbol("string-ref")] = stringRef
	NS[makeSymbol("string-contains")] = stringContains
	NS[makeSymbol("string-upcase")] = stringUpcase
	NS[makeSymbol("string-downcase")] = stringDowncase
	NS[makeSymbol("string-split")] = stringSplit
	NS[makeSymbol("string->list")] = stringToList
	NS[makeSymbol("list->string")] = listToString
	NS[makeSymbol("string->symbol")] = stringToSymbol
	NS[makeSymbol("symbol->string")] = symbolToString
	NS[makeSymbol("number->string")] = numberToString
	NS[makeSymbol("string->number")] = stringToNumber
	NS[makeSymbol("string-join")] = stringJoin
	NS[makeSymbol("string-trim")] = stringTrim

	// Higher-order list functions
	NS[makeSymbol("map")] = mapFn
	NS[makeSymbol("filter")] = filterFn
	NS[makeSymbol("for-each")] = forEachFn
	NS[makeSymbol("fold-left")] = foldLeft
	NS[makeSymbol("fold-right")] = foldRight
	NS[makeSymbol("foldl")] = foldLeft
	NS[makeSymbol("foldr")] = foldRight
	NS[makeSymbol("reduce")] = reduceFn
	NS[makeSymbol("append")] = appendFn
	NS[makeSymbol("reverse")] = reverseFn
	NS[makeSymbol("list-ref")] = listRef
	NS[makeSymbol("list-tail")] = listTail
	NS[makeSymbol("list?")] = listQ
	NS[makeSymbol("sort")] = sortFn
	NS[makeSymbol("any")] = anyFn
	NS[makeSymbol("every")] = everyFn
	NS[makeSymbol("flatten")] = flattenFn

	// Predicates
	NS[makeSymbol("procedure?")] = procedureQ
	NS[makeSymbol("symbol?")] = symbolQ
	NS[makeSymbol("keyword?")] = keywordQ
	NS[makeSymbol("nil?")] = nilQ

	// Vectors
	NS[makeSymbol("vector")] = vectorFn
	NS[makeSymbol("vector?")] = vectorQ
	NS[makeSymbol("make-vector")] = makeVector
	NS[makeSymbol("vector-ref")] = vectorRef
	NS[makeSymbol("vector-set!")] = vectorSet
	NS[makeSymbol("vector-length")] = vectorLength
	NS[makeSymbol("vector->list")] = vectorToList
	NS[makeSymbol("list->vector")] = listToVector
	NS[makeSymbol("vector-fill!")] = vectorFill

	// IO
	NS[makeSymbol("display")] = display
	NS[makeSymbol("newline")] = newline
	NS[makeSymbol("print")] = printFn
	NS[makeSymbol("println")] = printlnFn
	NS[makeSymbol("write")] = writeFn

	// Atoms / state
	NS[makeSymbol("atom")] = newAtom
	NS[makeSymbol("deref")] = derefAtom
	NS[makeSymbol("reset!")] = resetAtom
	NS[makeSymbol("swap!")] = swapAtom

	// Misc
	NS[makeSymbol("gensym")] = gensymFn
	NS[makeSymbol("error")] = errorFn
	NS[makeSymbol("not=")] = notEqual
	NS[makeSymbol("positive?")] = positiveQ
	NS[makeSymbol("negative?")] = negativeQ
	NS[makeSymbol("odd?")] = oddQ
	NS[makeSymbol("type-of")] = typeOf

	// Association lists
	NS[makeSymbol("assoc")] = assocFn
	NS[makeSymbol("assq")] = assqFn
	NS[makeSymbol("assv")] = assvFn

	// Membership
	NS[makeSymbol("member")] = memberFn
	NS[makeSymbol("memq")] = memqFn
	NS[makeSymbol("memv")] = memvFn

	// Multiple values
	NS[makeSymbol("values")] = valuesFn
	NS[makeSymbol("call-with-values")] = callWithValuesFn

	// Ports / IO
	NS[makeSymbol("open-input-file")] = openInputFile
	NS[makeSymbol("open-input-string")] = openInputString
	NS[makeSymbol("close-input-port")] = closeInputPort
	NS[makeSymbol("read")] = readFn
	NS[makeSymbol("read-char")] = readCharFn
	NS[makeSymbol("peek-char")] = peekCharFn
	NS[makeSymbol("char-ready?")] = charReadyFn
	NS[makeSymbol("eof-object?")] = eofObjectQ
	NS[makeSymbol("eof-object")] = func(a []Expression) (Expression, error) { return EOF{}, nil }
	NS[makeSymbol("input-port?")] = inputPortQ
}

// ── Math ─────────────────────────────────────────────────────────────────────

func toFloat(e Expression) (float64, error) {
	f, ok := ToFloat(e)
	if !ok {
		return 0, fmt.Errorf("expected number, got %T", e)
	}
	return f, nil
}

func addNum(a []Expression) (Expression, error) {
	var hasFloat bool
	var isum int
	var fsum float64
	for _, e := range a {
		switch v := e.(type) {
		case int:
			isum += v
			fsum += float64(v)
		case float64:
			hasFloat = true
			fsum += v
		default:
			return nil, fmt.Errorf("expected number, got %T", e)
		}
	}
	if hasFloat {
		return fsum, nil
	}
	return isum, nil
}

func subNum(a []Expression) (Expression, error) {
	if len(a) == 0 {
		return nil, errors.New("arity mismatch")
	}
	var hasFloat bool
	switch a[0].(type) {
	case float64:
		hasFloat = true
	}
	if len(a) == 1 {
		if hasFloat {
			return -a[0].(float64), nil
		}
		return -a[0].(int), nil
	}
	fres, ok := ToFloat(a[0])
	if !ok {
		return nil, fmt.Errorf("expected number, got %T", a[0])
	}
	ires := a[0].(int)
	if hasFloat {
		ires = 0
	}
	for _, e := range a[1:] {
		switch v := e.(type) {
		case int:
			ires -= v
			fres -= float64(v)
		case float64:
			hasFloat = true
			fres -= v
		default:
			return nil, fmt.Errorf("expected number, got %T", e)
		}
	}
	if hasFloat {
		return fres, nil
	}
	return ires, nil
}

func mulNum(a []Expression) (Expression, error) {
	var hasFloat bool
	ires := 1
	fres := 1.0
	for _, e := range a {
		switch v := e.(type) {
		case int:
			ires *= v
			fres *= float64(v)
		case float64:
			hasFloat = true
			fres *= v
		default:
			return nil, fmt.Errorf("expected number, got %T", e)
		}
	}
	if hasFloat {
		return fres, nil
	}
	return ires, nil
}

func divNum(a []Expression) (Expression, error) {
	if len(a) == 0 {
		return nil, errors.New("arity mismatch")
	}
	fres, ok := ToFloat(a[0])
	if !ok {
		return nil, fmt.Errorf("expected number, got %T", a[0])
	}
	_, isFloat := a[0].(float64)
	for _, e := range a[1:] {
		switch v := e.(type) {
		case int:
			if v == 0 {
				return nil, errors.New("division by zero")
			}
			fres /= float64(v)
		case float64:
			if v == 0 {
				return nil, errors.New("division by zero")
			}
			isFloat = true
			fres /= v
		default:
			return nil, fmt.Errorf("expected number, got %T", e)
		}
	}
	if isFloat {
		return fres, nil
	}
	return int(fres), nil
}

func absNum(a []Expression) (Expression, error) {
	if err := assertArgNum(a, 1); err != nil {
		return nil, err
	}
	switch v := a[0].(type) {
	case int:
		if v < 0 {
			return -v, nil
		}
		return v, nil
	case float64:
		return math.Abs(v), nil
	}
	return nil, fmt.Errorf("expected number, got %T", a[0])
}

func maxNum(a []Expression) (Expression, error) {
	if len(a) == 0 {
		return nil, errors.New("max requires at least 1 argument")
	}
	fmax, ok := ToFloat(a[0])
	if !ok {
		return nil, fmt.Errorf("expected number, got %T", a[0])
	}
	_, hasFloat := a[0].(float64)
	for _, e := range a[1:] {
		f, ok := ToFloat(e)
		if !ok {
			return nil, fmt.Errorf("expected number, got %T", e)
		}
		if _, ok := e.(float64); ok {
			hasFloat = true
		}
		if f > fmax {
			fmax = f
		}
	}
	if hasFloat {
		return fmax, nil
	}
	return int(fmax), nil
}

func minNum(a []Expression) (Expression, error) {
	if len(a) == 0 {
		return nil, errors.New("min requires at least 1 argument")
	}
	fmin, ok := ToFloat(a[0])
	if !ok {
		return nil, fmt.Errorf("expected number, got %T", a[0])
	}
	_, hasFloat := a[0].(float64)
	for _, e := range a[1:] {
		f, ok := ToFloat(e)
		if !ok {
			return nil, fmt.Errorf("expected number, got %T", e)
		}
		if _, ok := e.(float64); ok {
			hasFloat = true
		}
		if f < fmin {
			fmin = f
		}
	}
	if hasFloat {
		return fmin, nil
	}
	return int(fmin), nil
}

func modulo(a []Expression) (Expression, error) {
	if err := assertArgNum(a, 2); err != nil {
		return nil, err
	}
	x, ok1 := ToInt(a[0])
	y, ok2 := ToInt(a[1])
	if !ok1 || !ok2 {
		return nil, errors.New("modulo requires integers")
	}
	if y == 0 {
		return nil, errors.New("modulo by zero")
	}
	r := x % y
	if r != 0 && (r < 0) != (y < 0) {
		r += y
	}
	return r, nil
}

func remainder(a []Expression) (Expression, error) {
	if err := assertArgNum(a, 2); err != nil {
		return nil, err
	}
	x, ok1 := ToInt(a[0])
	y, ok2 := ToInt(a[1])
	if !ok1 || !ok2 {
		return nil, errors.New("remainder requires integers")
	}
	if y == 0 {
		return nil, errors.New("remainder by zero")
	}
	return x % y, nil
}

func quotient(a []Expression) (Expression, error) {
	if err := assertArgNum(a, 2); err != nil {
		return nil, err
	}
	x, ok1 := ToInt(a[0])
	y, ok2 := ToInt(a[1])
	if !ok1 || !ok2 {
		return nil, errors.New("quotient requires integers")
	}
	if y == 0 {
		return nil, errors.New("quotient by zero")
	}
	return x / y, nil
}

func floorFn(a []Expression) (Expression, error) {
	if err := assertArgNum(a, 1); err != nil {
		return nil, err
	}
	f, err := toFloat(a[0])
	if err != nil {
		return nil, err
	}
	return math.Floor(f), nil
}

func ceilingFn(a []Expression) (Expression, error) {
	if err := assertArgNum(a, 1); err != nil {
		return nil, err
	}
	f, err := toFloat(a[0])
	if err != nil {
		return nil, err
	}
	return math.Ceil(f), nil
}

func roundFn(a []Expression) (Expression, error) {
	if err := assertArgNum(a, 1); err != nil {
		return nil, err
	}
	f, err := toFloat(a[0])
	if err != nil {
		return nil, err
	}
	return math.Round(f), nil
}

func truncateFn(a []Expression) (Expression, error) {
	if err := assertArgNum(a, 1); err != nil {
		return nil, err
	}
	f, err := toFloat(a[0])
	if err != nil {
		return nil, err
	}
	return math.Trunc(f), nil
}

func sqrtFn(a []Expression) (Expression, error) {
	if err := assertArgNum(a, 1); err != nil {
		return nil, err
	}
	f, err := toFloat(a[0])
	if err != nil {
		return nil, err
	}
	return math.Sqrt(f), nil
}

func exactToInexact(a []Expression) (Expression, error) {
	if err := assertArgNum(a, 1); err != nil {
		return nil, err
	}
	f, ok := ToFloat(a[0])
	if !ok {
		return nil, fmt.Errorf("expected number, got %T", a[0])
	}
	return f, nil
}

func inexactToExact(a []Expression) (Expression, error) {
	if err := assertArgNum(a, 1); err != nil {
		return nil, err
	}
	i, ok := ToInt(a[0])
	if !ok {
		return nil, fmt.Errorf("expected number, got %T", a[0])
	}
	return i, nil
}

func numberToString(a []Expression) (Expression, error) {
	if err := assertArgNum(a, 1); err != nil {
		return nil, err
	}
	switch v := a[0].(type) {
	case int:
		return strconv.Itoa(v), nil
	case float64:
		return strconv.FormatFloat(v, 'g', -1, 64), nil
	}
	return nil, fmt.Errorf("expected number, got %T", a[0])
}

func stringToNumber(a []Expression) (Expression, error) {
	if err := assertArgNum(a, 1); err != nil {
		return nil, err
	}
	s, ok := a[0].(string)
	if !ok {
		return nil, errors.New("string->number requires a string")
	}
	if i, err := strconv.Atoi(s); err == nil {
		return i, nil
	}
	if f, err := strconv.ParseFloat(s, 64); err == nil {
		return f, nil
	}
	return false, nil // R7RS: return #f on failure
}

func integerQ(a []Expression) (Expression, error) {
	if err := assertArgNum(a, 1); err != nil {
		return nil, err
	}
	return Int_Q(a[0]), nil
}

func floatQ(a []Expression) (Expression, error) {
	if err := assertArgNum(a, 1); err != nil {
		return nil, err
	}
	return Float_Q(a[0]), nil
}

func positiveQ(a []Expression) (Expression, error) {
	if err := assertArgNum(a, 1); err != nil {
		return nil, err
	}
	f, ok := ToFloat(a[0])
	if !ok {
		return nil, fmt.Errorf("expected number, got %T", a[0])
	}
	return f > 0, nil
}

func negativeQ(a []Expression) (Expression, error) {
	if err := assertArgNum(a, 1); err != nil {
		return nil, err
	}
	f, ok := ToFloat(a[0])
	if !ok {
		return nil, fmt.Errorf("expected number, got %T", a[0])
	}
	return f < 0, nil
}

func oddQ(a []Expression) (Expression, error) {
	if err := assertArgNum(a, 1); err != nil {
		return nil, err
	}
	i, ok := ToInt(a[0])
	if !ok {
		return nil, fmt.Errorf("expected integer, got %T", a[0])
	}
	return i%2 != 0, nil
}

func notEqual(a []Expression) (Expression, error) {
	if err := assertArgNum(a, 2); err != nil {
		return nil, err
	}
	return !Equal_Q(a[0], a[1]), nil
}

// ── Boolean ───────────────────────────────────────────────────────────────────

func booleanQ(a []Expression) (Expression, error) {
	if err := assertArgNum(a, 1); err != nil {
		return nil, err
	}
	_, ok := a[0].(bool)
	return ok, nil
}

// ── Characters ────────────────────────────────────────────────────────────────

func charQ(a []Expression) (Expression, error) {
	if err := assertArgNum(a, 1); err != nil {
		return nil, err
	}
	return Char_Q(a[0]), nil
}

func charToInteger(a []Expression) (Expression, error) {
	if err := assertArgNum(a, 1); err != nil {
		return nil, err
	}
	r, ok := a[0].(rune)
	if !ok {
		return nil, errors.New("char->integer requires a character")
	}
	return int(r), nil
}

func integerToChar(a []Expression) (Expression, error) {
	if err := assertArgNum(a, 1); err != nil {
		return nil, err
	}
	i, ok := ToInt(a[0])
	if !ok {
		return nil, errors.New("integer->char requires an integer")
	}
	return rune(i), nil
}

func charEq(a []Expression) (Expression, error) {
	if err := assertArgNum(a, 2); err != nil {
		return nil, err
	}
	c1, ok1 := a[0].(rune)
	c2, ok2 := a[1].(rune)
	if !ok1 || !ok2 {
		return nil, errors.New("char=? requires characters")
	}
	return c1 == c2, nil
}

func charLt(a []Expression) (Expression, error) {
	if err := assertArgNum(a, 2); err != nil {
		return nil, err
	}
	c1, ok1 := a[0].(rune)
	c2, ok2 := a[1].(rune)
	if !ok1 || !ok2 {
		return nil, errors.New("char<? requires characters")
	}
	return c1 < c2, nil
}

func charGt(a []Expression) (Expression, error) {
	if err := assertArgNum(a, 2); err != nil {
		return nil, err
	}
	c1, ok1 := a[0].(rune)
	c2, ok2 := a[1].(rune)
	if !ok1 || !ok2 {
		return nil, errors.New("char>? requires characters")
	}
	return c1 > c2, nil
}

func charAlphabetic(a []Expression) (Expression, error) {
	if err := assertArgNum(a, 1); err != nil {
		return nil, err
	}
	r, ok := a[0].(rune)
	if !ok {
		return nil, errors.New("char-alphabetic? requires a character")
	}
	return (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z'), nil
}

func charNumeric(a []Expression) (Expression, error) {
	if err := assertArgNum(a, 1); err != nil {
		return nil, err
	}
	r, ok := a[0].(rune)
	if !ok {
		return nil, errors.New("char-numeric? requires a character")
	}
	return r >= '0' && r <= '9', nil
}

func charWhitespace(a []Expression) (Expression, error) {
	if err := assertArgNum(a, 1); err != nil {
		return nil, err
	}
	r, ok := a[0].(rune)
	if !ok {
		return nil, errors.New("char-whitespace? requires a character")
	}
	return r == ' ' || r == '\t' || r == '\n' || r == '\r', nil
}

func charUpcase(a []Expression) (Expression, error) {
	if err := assertArgNum(a, 1); err != nil {
		return nil, err
	}
	r, ok := a[0].(rune)
	if !ok {
		return nil, errors.New("char-upcase requires a character")
	}
	if r >= 'a' && r <= 'z' {
		return r - 32, nil
	}
	return r, nil
}

func charDowncase(a []Expression) (Expression, error) {
	if err := assertArgNum(a, 1); err != nil {
		return nil, err
	}
	r, ok := a[0].(rune)
	if !ok {
		return nil, errors.New("char-downcase requires a character")
	}
	if r >= 'A' && r <= 'Z' {
		return r + 32, nil
	}
	return r, nil
}

// ── Strings ───────────────────────────────────────────────────────────────────

func stringQ(a []Expression) (Expression, error) {
	if err := assertArgNum(a, 1); err != nil {
		return nil, err
	}
	_, ok := a[0].(string)
	return ok, nil
}

func stringEq(a []Expression) (Expression, error) {
	if err := assertArgNum(a, 2); err != nil {
		return nil, err
	}
	s1, ok1 := a[0].(string)
	s2, ok2 := a[1].(string)
	if !ok1 || !ok2 {
		return nil, errors.New("string=? requires strings")
	}
	return s1 == s2, nil
}

func stringLt(a []Expression) (Expression, error) {
	if err := assertArgNum(a, 2); err != nil {
		return nil, err
	}
	s1, ok1 := a[0].(string)
	s2, ok2 := a[1].(string)
	if !ok1 || !ok2 {
		return nil, errors.New("string<? requires strings")
	}
	return s1 < s2, nil
}

func stringGt(a []Expression) (Expression, error) {
	if err := assertArgNum(a, 2); err != nil {
		return nil, err
	}
	s1, ok1 := a[0].(string)
	s2, ok2 := a[1].(string)
	if !ok1 || !ok2 {
		return nil, errors.New("string>? requires strings")
	}
	return s1 > s2, nil
}

func stringLength(a []Expression) (Expression, error) {
	if err := assertArgNum(a, 1); err != nil {
		return nil, err
	}
	s, ok := a[0].(string)
	if !ok {
		return nil, errors.New("string-length requires a string")
	}
	return len([]rune(s)), nil
}

func stringAppend(a []Expression) (Expression, error) {
	var sb strings.Builder
	for _, e := range a {
		s, ok := e.(string)
		if !ok {
			return nil, fmt.Errorf("string-append requires strings, got %T", e)
		}
		sb.WriteString(s)
	}
	return sb.String(), nil
}

func substring(a []Expression) (Expression, error) {
	if len(a) < 2 || len(a) > 3 {
		return nil, errors.New("substring requires 2 or 3 arguments")
	}
	s, ok := a[0].(string)
	if !ok {
		return nil, errors.New("substring: first argument must be a string")
	}
	runes := []rune(s)
	start, ok2 := ToInt(a[1])
	if !ok2 {
		return nil, errors.New("substring: start must be an integer")
	}
	end := len(runes)
	if len(a) == 3 {
		end, ok2 = ToInt(a[2])
		if !ok2 {
			return nil, errors.New("substring: end must be an integer")
		}
	}
	if start < 0 || end > len(runes) || start > end {
		return nil, errors.New("substring: index out of range")
	}
	return string(runes[start:end]), nil
}

func stringRef(a []Expression) (Expression, error) {
	if err := assertArgNum(a, 2); err != nil {
		return nil, err
	}
	s, ok := a[0].(string)
	if !ok {
		return nil, errors.New("string-ref: first argument must be a string")
	}
	idx, ok2 := ToInt(a[1])
	if !ok2 {
		return nil, errors.New("string-ref: index must be an integer")
	}
	runes := []rune(s)
	if idx < 0 || idx >= len(runes) {
		return nil, errors.New("string-ref: index out of range")
	}
	return runes[idx], nil
}

func stringContains(a []Expression) (Expression, error) {
	if err := assertArgNum(a, 2); err != nil {
		return nil, err
	}
	s, ok1 := a[0].(string)
	sub, ok2 := a[1].(string)
	if !ok1 || !ok2 {
		return nil, errors.New("string-contains requires strings")
	}
	return strings.Contains(s, sub), nil
}

func stringUpcase(a []Expression) (Expression, error) {
	if err := assertArgNum(a, 1); err != nil {
		return nil, err
	}
	s, ok := a[0].(string)
	if !ok {
		return nil, errors.New("string-upcase requires a string")
	}
	return strings.ToUpper(s), nil
}

func stringDowncase(a []Expression) (Expression, error) {
	if err := assertArgNum(a, 1); err != nil {
		return nil, err
	}
	s, ok := a[0].(string)
	if !ok {
		return nil, errors.New("string-downcase requires a string")
	}
	return strings.ToLower(s), nil
}

func stringSplit(a []Expression) (Expression, error) {
	if len(a) < 1 || len(a) > 2 {
		return nil, errors.New("string-split requires 1 or 2 arguments")
	}
	s, ok := a[0].(string)
	if !ok {
		return nil, errors.New("string-split: first argument must be a string")
	}
	sep := " "
	if len(a) == 2 {
		sep, ok = a[1].(string)
		if !ok {
			return nil, errors.New("string-split: separator must be a string")
		}
	}
	parts := strings.Split(s, sep)
	exprs := make([]Expression, len(parts))
	for i, p := range parts {
		exprs[i] = p
	}
	return List{Val: exprs}, nil
}

func stringToList(a []Expression) (Expression, error) {
	if err := assertArgNum(a, 1); err != nil {
		return nil, err
	}
	s, ok := a[0].(string)
	if !ok {
		return nil, errors.New("string->list requires a string")
	}
	runes := []rune(s)
	exprs := make([]Expression, len(runes))
	for i, r := range runes {
		exprs[i] = r
	}
	return List{Val: exprs}, nil
}

func listToString(a []Expression) (Expression, error) {
	if err := assertArgNum(a, 1); err != nil {
		return nil, err
	}
	slc, err := GetSlice(a[0])
	if err != nil {
		return nil, err
	}
	var sb strings.Builder
	for _, e := range slc {
		r, ok := e.(rune)
		if !ok {
			return nil, fmt.Errorf("list->string: expected character, got %T", e)
		}
		sb.WriteRune(r)
	}
	return sb.String(), nil
}

func stringToSymbol(a []Expression) (Expression, error) {
	if err := assertArgNum(a, 1); err != nil {
		return nil, err
	}
	s, ok := a[0].(string)
	if !ok {
		return nil, errors.New("string->symbol requires a string")
	}
	return Symbol{Val: s}, nil
}

func symbolToString(a []Expression) (Expression, error) {
	if err := assertArgNum(a, 1); err != nil {
		return nil, err
	}
	s, ok := a[0].(Symbol)
	if !ok {
		return nil, errors.New("symbol->string requires a symbol")
	}
	return s.Val, nil
}

func stringJoin(a []Expression) (Expression, error) {
	if len(a) < 1 || len(a) > 2 {
		return nil, errors.New("string-join requires 1 or 2 arguments")
	}
	slc, err := GetSlice(a[0])
	if err != nil {
		return nil, err
	}
	sep := ""
	if len(a) == 2 {
		var ok bool
		sep, ok = a[1].(string)
		if !ok {
			return nil, errors.New("string-join: separator must be a string")
		}
	}
	parts := make([]string, len(slc))
	for i, e := range slc {
		s, ok := e.(string)
		if !ok {
			return nil, fmt.Errorf("string-join: expected string, got %T", e)
		}
		parts[i] = s
	}
	return strings.Join(parts, sep), nil
}

func stringTrim(a []Expression) (Expression, error) {
	if err := assertArgNum(a, 1); err != nil {
		return nil, err
	}
	s, ok := a[0].(string)
	if !ok {
		return nil, errors.New("string-trim requires a string")
	}
	return strings.TrimSpace(s), nil
}

// ── Higher-order list functions ───────────────────────────────────────────────

func mapFn(a []Expression) (Expression, error) {
	if len(a) < 2 {
		return nil, errors.New("map requires at least 2 arguments")
	}
	f := a[0]
	results := []Expression{}
	if len(a) == 2 {
		args, e := GetSlice(a[1])
		if e != nil {
			return nil, e
		}
		for _, arg := range args {
			res, e := Apply(f, []Expression{arg})
			if e != nil {
				return nil, e
			}
			results = append(results, res)
		}
	} else {
		// multi-list map
		slices := make([][]Expression, len(a)-1)
		minLen := math.MaxInt32
		for i, lst := range a[1:] {
			slc, e := GetSlice(lst)
			if e != nil {
				return nil, e
			}
			slices[i] = slc
			if len(slc) < minLen {
				minLen = len(slc)
			}
		}
		for i := 0; i < minLen; i++ {
			args := make([]Expression, len(slices))
			for j, slc := range slices {
				args[j] = slc[i]
			}
			res, e := Apply(f, args)
			if e != nil {
				return nil, e
			}
			results = append(results, res)
		}
	}
	return List{Val: results}, nil
}

func filterFn(a []Expression) (Expression, error) {
	if err := assertArgNum(a, 2); err != nil {
		return nil, err
	}
	f := a[0]
	args, e := GetSlice(a[1])
	if e != nil {
		return nil, e
	}
	results := []Expression{}
	for _, arg := range args {
		res, e := Apply(f, []Expression{arg})
		if e != nil {
			return nil, e
		}
		if res != nil && res != false {
			results = append(results, arg)
		}
	}
	return List{Val: results}, nil
}

func forEachFn(a []Expression) (Expression, error) {
	if err := assertArgNum(a, 2); err != nil {
		return nil, err
	}
	f := a[0]
	args, e := GetSlice(a[1])
	if e != nil {
		return nil, e
	}
	for _, arg := range args {
		if _, e := Apply(f, []Expression{arg}); e != nil {
			return nil, e
		}
	}
	return nil, nil
}

func foldLeft(a []Expression) (Expression, error) {
	if err := assertArgNum(a, 3); err != nil {
		return nil, err
	}
	f := a[0]
	acc := a[1]
	slc, e := GetSlice(a[2])
	if e != nil {
		return nil, e
	}
	for _, x := range slc {
		acc, e = Apply(f, []Expression{acc, x})
		if e != nil {
			return nil, e
		}
	}
	return acc, nil
}

func foldRight(a []Expression) (Expression, error) {
	if err := assertArgNum(a, 3); err != nil {
		return nil, err
	}
	f := a[0]
	acc := a[1]
	slc, e := GetSlice(a[2])
	if e != nil {
		return nil, e
	}
	for i := len(slc) - 1; i >= 0; i-- {
		acc, e = Apply(f, []Expression{slc[i], acc})
		if e != nil {
			return nil, e
		}
	}
	return acc, nil
}

func reduceFn(a []Expression) (Expression, error) {
	if len(a) < 2 {
		return nil, errors.New("reduce requires at least 2 arguments")
	}
	f := a[0]
	slc, e := GetSlice(a[len(a)-1])
	if e != nil {
		return nil, e
	}
	if len(slc) == 0 {
		if len(a) == 3 {
			return a[1], nil // identity element
		}
		return nil, errors.New("reduce: empty list with no identity")
	}
	acc := slc[0]
	for _, x := range slc[1:] {
		acc, e = Apply(f, []Expression{acc, x})
		if e != nil {
			return nil, e
		}
	}
	return acc, nil
}

func appendFn(a []Expression) (Expression, error) {
	result := []Expression{}
	for i, arg := range a {
		if i == len(a)-1 {
			// Last arg: if it's a list append its elements; if atom cons it
			switch v := arg.(type) {
			case List:
				result = append(result, v.Val...)
			default:
				if len(result) == 0 {
					return arg, nil
				}
				// improper list - just append as element
				result = append(result, arg)
			}
		} else {
			slc, e := GetSlice(arg)
			if e != nil {
				return nil, e
			}
			result = append(result, slc...)
		}
	}
	return List{Val: result}, nil
}

func reverseFn(a []Expression) (Expression, error) {
	if err := assertArgNum(a, 1); err != nil {
		return nil, err
	}
	slc, e := GetSlice(a[0])
	if e != nil {
		return nil, e
	}
	rev := make([]Expression, len(slc))
	for i, v := range slc {
		rev[len(slc)-1-i] = v
	}
	return List{Val: rev}, nil
}

func listRef(a []Expression) (Expression, error) {
	if err := assertArgNum(a, 2); err != nil {
		return nil, err
	}
	slc, e := GetSlice(a[0])
	if e != nil {
		return nil, e
	}
	idx, ok := ToInt(a[1])
	if !ok {
		return nil, errors.New("list-ref: index must be an integer")
	}
	if idx < 0 || idx >= len(slc) {
		return nil, errors.New("list-ref: index out of range")
	}
	return slc[idx], nil
}

func listTail(a []Expression) (Expression, error) {
	if err := assertArgNum(a, 2); err != nil {
		return nil, err
	}
	slc, e := GetSlice(a[0])
	if e != nil {
		return nil, e
	}
	idx, ok := ToInt(a[1])
	if !ok {
		return nil, errors.New("list-tail: index must be an integer")
	}
	if idx < 0 || idx > len(slc) {
		return nil, errors.New("list-tail: index out of range")
	}
	return List{Val: slc[idx:]}, nil
}

func listQ(a []Expression) (Expression, error) {
	if err := assertArgNum(a, 1); err != nil {
		return nil, err
	}
	return List_Q(a[0]), nil
}

func sortFn(a []Expression) (Expression, error) {
	if err := assertArgNum(a, 2); err != nil {
		return nil, err
	}
	slc, e := GetSlice(a[1])
	if e != nil {
		return nil, e
	}
	f := a[0]
	cp := make([]Expression, len(slc))
	copy(cp, slc)
	var sortErr error
	sort.SliceStable(cp, func(i, j int) bool {
		if sortErr != nil {
			return false
		}
		res, err := Apply(f, []Expression{cp[i], cp[j]})
		if err != nil {
			sortErr = err
			return false
		}
		return res != nil && res != false
	})
	if sortErr != nil {
		return nil, sortErr
	}
	return List{Val: cp}, nil
}

func anyFn(a []Expression) (Expression, error) {
	if err := assertArgNum(a, 2); err != nil {
		return nil, err
	}
	f := a[0]
	slc, e := GetSlice(a[1])
	if e != nil {
		return nil, e
	}
	for _, x := range slc {
		res, e := Apply(f, []Expression{x})
		if e != nil {
			return nil, e
		}
		if res != nil && res != false {
			return true, nil
		}
	}
	return false, nil
}

func everyFn(a []Expression) (Expression, error) {
	if err := assertArgNum(a, 2); err != nil {
		return nil, err
	}
	f := a[0]
	slc, e := GetSlice(a[1])
	if e != nil {
		return nil, e
	}
	for _, x := range slc {
		res, e := Apply(f, []Expression{x})
		if e != nil {
			return nil, e
		}
		if res == nil || res == false {
			return false, nil
		}
	}
	return true, nil
}

func flattenFn(a []Expression) (Expression, error) {
	if err := assertArgNum(a, 1); err != nil {
		return nil, err
	}
	var flat func(Expression) []Expression
	flat = func(e Expression) []Expression {
		if List_Q(e) {
			result := []Expression{}
			for _, x := range e.(List).Val {
				result = append(result, flat(x)...)
			}
			return result
		}
		return []Expression{e}
	}
	return List{Val: flat(a[0])}, nil
}

// ── Predicates ────────────────────────────────────────────────────────────────

func procedureQ(a []Expression) (Expression, error) {
	if err := assertArgNum(a, 1); err != nil {
		return nil, err
	}
	switch a[0].(type) {
	case Func, ExpressionFunc, func([]Expression) (Expression, error):
		return true, nil
	}
	return false, nil
}

func symbolQ(a []Expression) (Expression, error) {
	if err := assertArgNum(a, 1); err != nil {
		return nil, err
	}
	return Symbol_Q(a[0]), nil
}

func keywordQ(a []Expression) (Expression, error) {
	if err := assertArgNum(a, 1); err != nil {
		return nil, err
	}
	return Keyword_Q(a[0]), nil
}

func nilQ(a []Expression) (Expression, error) {
	if err := assertArgNum(a, 1); err != nil {
		return nil, err
	}
	return Nil_Q(a[0]), nil
}

// ── Vectors ───────────────────────────────────────────────────────────────────

func vectorFn(a []Expression) (Expression, error) {
	return Vector{Val: a}, nil
}

func vectorQ(a []Expression) (Expression, error) {
	if err := assertArgNum(a, 1); err != nil {
		return nil, err
	}
	return Vector_Q(a[0]), nil
}

func makeVector(a []Expression) (Expression, error) {
	if len(a) < 1 || len(a) > 2 {
		return nil, errors.New("make-vector requires 1 or 2 arguments")
	}
	n, ok := ToInt(a[0])
	if !ok {
		return nil, errors.New("make-vector: size must be an integer")
	}
	var fill Expression
	if len(a) == 2 {
		fill = a[1]
	}
	v := make([]Expression, n)
	for i := range v {
		v[i] = fill
	}
	return Vector{Val: v}, nil
}

func vectorRef(a []Expression) (Expression, error) {
	if err := assertArgNum(a, 2); err != nil {
		return nil, err
	}
	vec, ok := a[0].(Vector)
	if !ok {
		return nil, errors.New("vector-ref: first argument must be a vector")
	}
	idx, ok2 := ToInt(a[1])
	if !ok2 {
		return nil, errors.New("vector-ref: index must be an integer")
	}
	if idx < 0 || idx >= len(vec.Val) {
		return nil, errors.New("vector-ref: index out of range")
	}
	return vec.Val[idx], nil
}

func vectorSet(a []Expression) (Expression, error) {
	if err := assertArgNum(a, 3); err != nil {
		return nil, err
	}
	vec, ok := a[0].(Vector)
	if !ok {
		return nil, errors.New("vector-set!: first argument must be a vector")
	}
	idx, ok2 := ToInt(a[1])
	if !ok2 {
		return nil, errors.New("vector-set!: index must be an integer")
	}
	if idx < 0 || idx >= len(vec.Val) {
		return nil, errors.New("vector-set!: index out of range")
	}
	vec.Val[idx] = a[2]
	return nil, nil
}

func vectorLength(a []Expression) (Expression, error) {
	if err := assertArgNum(a, 1); err != nil {
		return nil, err
	}
	vec, ok := a[0].(Vector)
	if !ok {
		return nil, errors.New("vector-length: first argument must be a vector")
	}
	return len(vec.Val), nil
}

func vectorToList(a []Expression) (Expression, error) {
	if err := assertArgNum(a, 1); err != nil {
		return nil, err
	}
	vec, ok := a[0].(Vector)
	if !ok {
		return nil, errors.New("vector->list: first argument must be a vector")
	}
	return List{Val: vec.Val}, nil
}

func listToVector(a []Expression) (Expression, error) {
	if err := assertArgNum(a, 1); err != nil {
		return nil, err
	}
	slc, e := GetSlice(a[0])
	if e != nil {
		return nil, e
	}
	return Vector{Val: slc}, nil
}

func vectorFill(a []Expression) (Expression, error) {
	if err := assertArgNum(a, 2); err != nil {
		return nil, err
	}
	vec, ok := a[0].(Vector)
	if !ok {
		return nil, errors.New("vector-fill!: first argument must be a vector")
	}
	for i := range vec.Val {
		vec.Val[i] = a[1]
	}
	return nil, nil
}

// ── IO ────────────────────────────────────────────────────────────────────────

func display(a []Expression) (Expression, error) {
	if len(a) == 0 {
		return nil, errors.New("display requires at least 1 argument")
	}
	fmt.Print(displayStr(a[0]))
	return nil, nil
}

func displayStr(e Expression) string {
	switch v := e.(type) {
	case string:
		return v
	case rune:
		return string(v)
	default:
		return fmt.Sprintf("%v", v)
	}
}

func newline(a []Expression) (Expression, error) {
	fmt.Println()
	return nil, nil
}

func printFn(a []Expression) (Expression, error) {
	parts := make([]string, len(a))
	for i, x := range a {
		parts[i] = displayStr(x)
	}
	fmt.Print(strings.Join(parts, " "))
	return nil, nil
}

func printlnFn(a []Expression) (Expression, error) {
	parts := make([]string, len(a))
	for i, x := range a {
		parts[i] = displayStr(x)
	}
	fmt.Println(strings.Join(parts, " "))
	return nil, nil
}

func writeFn(a []Expression) (Expression, error) {
	if len(a) == 0 {
		return nil, errors.New("write requires at least 1 argument")
	}
	fmt.Printf("%v", a[0])
	return nil, nil
}

// ── Atoms ─────────────────────────────────────────────────────────────────────

func newAtom(a []Expression) (Expression, error) {
	if err := assertArgNum(a, 1); err != nil {
		return nil, err
	}
	return &Atom{Val: a[0]}, nil
}

func derefAtom(a []Expression) (Expression, error) {
	if err := assertArgNum(a, 1); err != nil {
		return nil, err
	}
	atm, ok := a[0].(*Atom)
	if !ok {
		return nil, errors.New("deref: expected atom")
	}
	return atm.Val, nil
}

func resetAtom(a []Expression) (Expression, error) {
	if err := assertArgNum(a, 2); err != nil {
		return nil, err
	}
	atm, ok := a[0].(*Atom)
	if !ok {
		return nil, errors.New("reset!: expected atom")
	}
	atm.Set(a[1])
	return a[1], nil
}

func swapAtom(a []Expression) (Expression, error) {
	if len(a) < 2 {
		return nil, errors.New("swap! requires at least 2 arguments")
	}
	atm, ok := a[0].(*Atom)
	if !ok {
		return nil, errors.New("swap!: expected atom")
	}
	args := []Expression{atm.Val}
	args = append(args, a[2:]...)
	res, e := Apply(a[1], args)
	if e != nil {
		return nil, e
	}
	atm.Set(res)
	return res, nil
}

// ── Misc ──────────────────────────────────────────────────────────────────────

var gensymCounter int

func gensymFn(a []Expression) (Expression, error) {
	gensymCounter++
	prefix := "G"
	if len(a) == 1 {
		if s, ok := a[0].(string); ok {
			prefix = s
		}
	}
	return Symbol{Val: fmt.Sprintf("%s%d", prefix, gensymCounter)}, nil
}

func errorFn(a []Expression) (Expression, error) {
	if len(a) == 0 {
		return nil, errors.New("error")
	}
	parts := make([]string, len(a))
	for i, x := range a {
		parts[i] = fmt.Sprintf("%v", x)
	}
	return nil, errors.New(strings.Join(parts, " "))
}

func typeOf(a []Expression) (Expression, error) {
	if err := assertArgNum(a, 1); err != nil {
		return nil, err
	}
	switch a[0].(type) {
	case int:
		return "integer", nil
	case float64:
		return "float", nil
	case bool:
		return "boolean", nil
	case string:
		return "string", nil
	case rune:
		return "char", nil
	case Symbol:
		return "symbol", nil
	case Keyword:
		return "keyword", nil
	case List:
		return "list", nil
	case Vector:
		return "vector", nil
	case HashMap:
		return "hash-map", nil
	case *Atom:
		return "atom", nil
	case Func, ExpressionFunc:
		return "procedure", nil
	case nil:
		return "nil", nil
	}
	return "unknown", nil
}

// ── Association lists ─────────────────────────────────────────────────────────

func assocFn(a []Expression) (Expression, error) {
	if len(a) != 2 {
		return nil, errors.New("assoc: requires 2 arguments")
	}
	return assocSearch(a[0], a[1], Equal_Q)
}

func assqFn(a []Expression) (Expression, error) {
	if len(a) != 2 {
		return nil, errors.New("assq: requires 2 arguments")
	}
	return assocSearch(a[0], a[1], func(x, y Expression) bool { return x == y })
}

func assvFn(a []Expression) (Expression, error) {
	if len(a) != 2 {
		return nil, errors.New("assv: requires 2 arguments")
	}
	return assocSearch(a[0], a[1], func(x, y Expression) bool { return x == y })
}

func assocSearch(key, lst Expression, cmp func(Expression, Expression) bool) (Expression, error) {
	if lst == nil {
		return false, nil
	}
	items, e := GetSlice(lst)
	if e != nil {
		return nil, errors.New("assoc: second argument must be a list")
	}
	for _, item := range items {
		pair, e := GetSlice(item)
		if e != nil || len(pair) < 1 {
			continue
		}
		if cmp(key, pair[0]) {
			return item, nil
		}
	}
	return false, nil
}

// ── Membership ────────────────────────────────────────────────────────────────

func memberFn(a []Expression) (Expression, error) {
	if len(a) != 2 {
		return nil, errors.New("member: requires 2 arguments")
	}
	return memberSearch(a[0], a[1], Equal_Q)
}

func memqFn(a []Expression) (Expression, error) {
	if len(a) != 2 {
		return nil, errors.New("memq: requires 2 arguments")
	}
	return memberSearch(a[0], a[1], func(x, y Expression) bool { return x == y })
}

func memvFn(a []Expression) (Expression, error) {
	if len(a) != 2 {
		return nil, errors.New("memv: requires 2 arguments")
	}
	return memberSearch(a[0], a[1], func(x, y Expression) bool { return x == y })
}

func memberSearch(x, lst Expression, cmp func(Expression, Expression) bool) (Expression, error) {
	if lst == nil {
		return false, nil
	}
	items, e := GetSlice(lst)
	if e != nil {
		return nil, errors.New("member: second argument must be a list")
	}
	for i, item := range items {
		if cmp(x, item) {
			return List{Val: items[i:]}, nil
		}
	}
	return false, nil
}

// ── Multiple values ───────────────────────────────────────────────────────────

func valuesFn(a []Expression) (Expression, error) {
	if len(a) == 1 {
		return a[0], nil
	}
	return MultipleValues{Vals: a}, nil
}

func callWithValuesFn(a []Expression) (Expression, error) {
	if len(a) != 2 {
		return nil, errors.New("call-with-values: requires 2 arguments")
	}
	produced, e := Apply(a[0], []Expression{})
	if e != nil {
		return nil, e
	}
	var args []Expression
	if mv, ok := produced.(MultipleValues); ok {
		args = mv.Vals
	} else {
		args = []Expression{produced}
	}
	return Apply(a[1], args)
}

// ── Ports / IO ────────────────────────────────────────────────────────────────

func openInputFile(a []Expression) (Expression, error) {
	if len(a) != 1 {
		return nil, errors.New("open-input-file: requires 1 argument")
	}
	name, ok := a[0].(string)
	if !ok {
		return nil, errors.New("open-input-file: argument must be a string")
	}
	f, e := os.Open(name)
	if e != nil {
		return nil, e
	}
	return &Port{R: lexer.NewLexer(f), Name: name}, nil
}

func openInputString(a []Expression) (Expression, error) {
	if len(a) != 1 {
		return nil, errors.New("open-input-string: requires 1 argument")
	}
	s, ok := a[0].(string)
	if !ok {
		return nil, errors.New("open-input-string: argument must be a string")
	}
	return &Port{R: lexer.NewLexer(strings.NewReader(s)), Name: "<string>"}, nil
}

func closeInputPort(a []Expression) (Expression, error) {
	if len(a) != 1 {
		return nil, errors.New("close-input-port: requires 1 argument")
	}
	p, ok := a[0].(*Port)
	if !ok {
		return nil, errors.New("close-input-port: argument must be a port")
	}
	p.Closed = true
	return nil, nil
}

func readFn(a []Expression) (Expression, error) {
	if len(a) != 1 {
		return nil, errors.New("read: requires 1 argument (port)")
	}
	p, ok := a[0].(*Port)
	if !ok {
		return nil, errors.New("read: argument must be a port")
	}
	if p.Closed {
		return EOF{}, nil
	}
	exp, e := p.R.ReadExpr()
	if e != nil {
		if strings.Contains(e.Error(), "<empty line>") || strings.Contains(e.Error(), "EOF") {
			return EOF{}, nil
		}
		return nil, e
	}
	return exp, nil
}

func readCharFn(a []Expression) (Expression, error) {
	if len(a) != 1 {
		return nil, errors.New("read-char: requires 1 argument (port)")
	}
	p, ok := a[0].(*Port)
	if !ok {
		return nil, errors.New("read-char: argument must be a port")
	}
	if p.Closed {
		return EOF{}, nil
	}
	r := p.R.NextRune()
	if r == scanner.EOF {
		return EOF{}, nil
	}
	return r, nil
}

func peekCharFn(a []Expression) (Expression, error) {
	if len(a) != 1 {
		return nil, errors.New("peek-char: requires 1 argument (port)")
	}
	p, ok := a[0].(*Port)
	if !ok {
		return nil, errors.New("peek-char: argument must be a port")
	}
	if p.Closed {
		return EOF{}, nil
	}
	r := p.R.PeekRune()
	if r == scanner.EOF {
		return EOF{}, nil
	}
	return r, nil
}

func charReadyFn(a []Expression) (Expression, error) {
	if len(a) != 1 {
		return nil, errors.New("char-ready?: requires 1 argument (port)")
	}
	_, ok := a[0].(*Port)
	if !ok {
		return nil, errors.New("char-ready?: argument must be a port")
	}
	return true, nil // always ready for string/file ports
}

func eofObjectQ(a []Expression) (Expression, error) {
	if len(a) != 1 {
		return nil, errors.New("eof-object?: requires 1 argument")
	}
	return EOF_Q(a[0]), nil
}

func inputPortQ(a []Expression) (Expression, error) {
	if len(a) != 1 {
		return nil, errors.New("input-port?: requires 1 argument")
	}
	return Port_Q(a[0]), nil
}
