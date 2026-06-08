package repltests

import (
	"testing"

	"github.com/holmes89/shrew/core"
	"github.com/holmes89/shrew/env"
	. "github.com/holmes89/shrew/repl"
	. "github.com/holmes89/shrew/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newTestEnv() EnvType {
	ev := env.DefaultEnv()
	for k, v := range core.NS {
		ev.Set(k, Func{Fn: v})
	}
	ev.Set(Symbol{Val: "eval"}, Func{
		Fn: func(a []Expression) (Expression, error) {
			return Eval(a, ev)
		},
	})
	return ev
}

func evalStr(t *testing.T, ev EnvType, expr string) Expression {
	t.Helper()
	res, err := Repl(expr, ev)
	require.NoError(t, err, "evaluating: %s", expr)
	return res
}

func TestPhase1Floats(t *testing.T) {
	ev := newTestEnv()
	assert.Equal(t, "3.14", evalStr(t, ev, "3.14"))
	assert.Equal(t, "4", evalStr(t, ev, "(+ 1.5 2.5)"))
	assert.Equal(t, "3", evalStr(t, ev, "(+ 1 2.0)"))
	assert.Equal(t, "6.28", evalStr(t, ev, "(* 2 3.14)"))
	assert.Equal(t, "3", evalStr(t, ev, "(floor 3.7)"))
	assert.Equal(t, "4", evalStr(t, ev, "(ceiling 3.2)"))
	assert.Equal(t, "2", evalStr(t, ev, "(sqrt 4.0)"))
	assert.Equal(t, "5", evalStr(t, ev, "(abs -5)"))
	assert.Equal(t, "3", evalStr(t, ev, "(max 1 2 3)"))
	assert.Equal(t, "1", evalStr(t, ev, "(min 1 2 3)"))
	assert.Equal(t, "1", evalStr(t, ev, "(modulo 10 3)"))
	assert.Equal(t, "2", evalStr(t, ev, "(remainder 11 3)"))
	assert.Equal(t, "3", evalStr(t, ev, "(quotient 10 3)"))
	assert.Equal(t, "true", evalStr(t, ev, "(integer? 42)"))
	assert.Equal(t, "true", evalStr(t, ev, "(float? 3.14)"))
	assert.Equal(t, "false", evalStr(t, ev, "(integer? 3.14)"))
	assert.Equal(t, "42", evalStr(t, ev, "(inexact->exact 42.9)"))
	assert.Equal(t, "3.14", evalStr(t, ev, "(exact->inexact 3.14)"))
	assert.Equal(t, "true", evalStr(t, ev, "(positive? 3)"))
	assert.Equal(t, "true", evalStr(t, ev, "(negative? -1)"))
	assert.Equal(t, "true", evalStr(t, ev, "(odd? 3)"))
	assert.Equal(t, "false", evalStr(t, ev, "(odd? 4)"))
}

func TestPhase1Strings(t *testing.T) {
	ev := newTestEnv()
	assert.Equal(t, "5", evalStr(t, ev, `(string-length "hello")`))
	assert.Equal(t, "foobar", evalStr(t, ev, `(string-append "foo" "bar")`))
	assert.Equal(t, "el", evalStr(t, ev, `(substring "hello" 1 3)`))
	assert.Equal(t, "42", evalStr(t, ev, `(string->number "42")`))
	assert.Equal(t, "42", evalStr(t, ev, `(number->string 42)`))
	assert.Equal(t, "HELLO", evalStr(t, ev, `(string-upcase "hello")`))
	assert.Equal(t, "hello", evalStr(t, ev, `(string-downcase "HELLO")`))
	assert.Equal(t, "true", evalStr(t, ev, `(string? "hello")`))
	assert.Equal(t, "true", evalStr(t, ev, `(string=? "abc" "abc")`))
	assert.Equal(t, "true", evalStr(t, ev, `(string<? "abc" "abd")`))
	assert.Equal(t, "true", evalStr(t, ev, `(string-contains "hello world" "world")`))
	assert.Equal(t, "hello", evalStr(t, ev, `(string-trim "  hello  ")`))
}

func TestPhase1HigherOrder(t *testing.T) {
	ev := newTestEnv()
	assert.Equal(t, "'(1 4 9 16)", evalStr(t, ev, "(map (lambda (x) (* x x)) '(1 2 3 4))"))
	assert.Equal(t, "'(1 3 5)", evalStr(t, ev, "(filter odd? '(1 2 3 4 5))"))
	assert.Equal(t, "15", evalStr(t, ev, "(fold-left + 0 '(1 2 3 4 5))"))
	assert.Equal(t, "15", evalStr(t, ev, "(fold-right + 0 '(1 2 3 4 5))"))
	assert.Equal(t, "15", evalStr(t, ev, "(reduce + '(1 2 3 4 5))"))
	assert.Equal(t, "'(3 2 1)", evalStr(t, ev, "(reverse '(1 2 3))"))
	assert.Equal(t, "'(1 2 3 4)", evalStr(t, ev, "(append '(1 2) '(3 4))"))
	assert.Equal(t, "2", evalStr(t, ev, "(list-ref '(1 2 3) 1)"))
	assert.Equal(t, "'(3)", evalStr(t, ev, "(list-tail '(1 2 3) 2)"))
	assert.Equal(t, "true", evalStr(t, ev, "(any odd? '(2 4 5 6))"))
	assert.Equal(t, "false", evalStr(t, ev, "(every odd? '(1 3 4))"))
	assert.Equal(t, "'(1 2 3 4)", evalStr(t, ev, "(flatten '(1 (2 (3)) 4))"))
}

func TestPhase1Vectors(t *testing.T) {
	ev := newTestEnv()
	assert.Equal(t, "20", evalStr(t, ev, "(vector-ref (vector 10 20 30) 1)"))
	assert.Equal(t, "5", evalStr(t, ev, "(vector-length (make-vector 5 0))"))
	assert.Equal(t, "'(10 20 30)", evalStr(t, ev, "(vector->list (vector 10 20 30))"))
	assert.Equal(t, "true", evalStr(t, ev, "(vector? (vector 1 2 3))"))
	assert.Equal(t, "false", evalStr(t, ev, "(vector? '(1 2 3))"))
}

func TestPhase1EvalForms(t *testing.T) {
	ev := newTestEnv()
	assert.Equal(t, "3", evalStr(t, ev, "(begin 1 2 3)"))
	assert.Equal(t, "8", evalStr(t, ev, "(let ((x 5) (y 3)) (+ x y))"))
	assert.Equal(t, "8", evalStr(t, ev, "(let* ((x 5) (y (+ x 3))) y)"))
	assert.Equal(t, "42", evalStr(t, ev, "(when #t 42)"))
	assert.Equal(t, "<nil>", evalStr(t, ev, "(when #f 42)"))
	assert.Equal(t, "99", evalStr(t, ev, "(unless #f 99)"))
	assert.Equal(t, "<nil>", evalStr(t, ev, "(unless #t 99)"))
	// letrec mutual recursion
	res := evalStr(t, ev, `(letrec ((my-even? (lambda (n) (if (= n 0) #t (my-odd? (- n 1)))))
	                          (my-odd? (lambda (n) (if (= n 0) #f (my-even? (- n 1))))))
	               (my-even? 10))`)
	assert.Equal(t, "true", res)
	// define shorthand
	evalStr(t, ev, "(define (square x) (* x x))")
	assert.Equal(t, "25", evalStr(t, ev, "(square 5)"))
	// multi-body lambda
	assert.Equal(t, "6", evalStr(t, ev, "((lambda (x) (define y (* x 2)) (+ y x)) 2)"))
	// cond with else
	assert.Equal(t, "big", evalStr(t, ev, `(cond ((< 10 5) "small") (else "big"))`))
	// and / or special forms
	assert.Equal(t, "3", evalStr(t, ev, "(and 1 2 3)"))
	assert.Equal(t, "false", evalStr(t, ev, "(and 1 #f 3)"))
	assert.Equal(t, "1", evalStr(t, ev, "(or #f 1 2)"))
	assert.Equal(t, "false", evalStr(t, ev, "(or #f #f)"))
}

func TestPhase1Predicates(t *testing.T) {
	ev := newTestEnv()
	assert.Equal(t, "true", evalStr(t, ev, "(boolean? #t)"))
	assert.Equal(t, "true", evalStr(t, ev, "(boolean? #f)"))
	assert.Equal(t, "false", evalStr(t, ev, `(boolean? "hello")`))
	assert.Equal(t, "true", evalStr(t, ev, `(string? "hello")`))
	assert.Equal(t, "true", evalStr(t, ev, "(procedure? car)"))
	assert.Equal(t, "true", evalStr(t, ev, "(procedure? (lambda (x) x))"))
	assert.Equal(t, "true", evalStr(t, ev, "(symbol? 'foo)"))
}
