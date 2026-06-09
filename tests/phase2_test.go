package repltests

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPhase2AssocMember(t *testing.T) {
	ev := newTestEnv()

	// assoc — structural equality
	assert.Equal(t, "'(b 2)", evalStr(t, ev, "(assoc 'b '((a 1) (b 2) (c 3)))"))
	assert.Equal(t, "'(a 1)", evalStr(t, ev, "(assoc 'a '((a 1) (b 2)))"))
	assert.Equal(t, "false", evalStr(t, ev, "(assoc 'd '((a 1) (b 2)))"))
	assert.Equal(t, "false", evalStr(t, ev, "(assoc 'x '())"))

	// assq / assv — identity comparison
	assert.Equal(t, "'(b 2)", evalStr(t, ev, "(assq 'b '((a 1) (b 2) (c 3)))"))
	assert.Equal(t, "false", evalStr(t, ev, "(assq 'd '((a 1) (b 2)))"))
	assert.Equal(t, "'(2 two)", evalStr(t, ev, "(assv 2 '((1 one) (2 two) (3 three)))"))

	// member — returns sublist or #f
	assert.Equal(t, "'(2 3)", evalStr(t, ev, "(member 2 '(1 2 3))"))
	assert.Equal(t, "'(3)", evalStr(t, ev, "(member 3 '(1 2 3))"))
	assert.Equal(t, "false", evalStr(t, ev, "(member 4 '(1 2 3))"))
	assert.Equal(t, "false", evalStr(t, ev, "(member 1 '())"))

	// memq / memv
	assert.Equal(t, "'(b c)", evalStr(t, ev, "(memq 'b '(a b c))"))
	assert.Equal(t, "false", evalStr(t, ev, "(memq 'd '(a b c))"))
	assert.Equal(t, "'(2 3)", evalStr(t, ev, "(memv 2 '(1 2 3))"))
}

func TestPhase2DoLoop(t *testing.T) {
	ev := newTestEnv()

	// basic counting do loop
	assert.Equal(t, "'done", evalStr(t, ev, "(do ((i 0 (+ i 1))) ((= i 5) 'done))"))

	// accumulate a sum
	assert.Equal(t, "10", evalStr(t, ev, "(do ((i 1 (+ i 1)) (s 0 (+ s i))) ((> i 4) s))"))

	// do with body side effects (use atom to observe)
	evalStr(t, ev, "(define counter (atom 0))")
	evalStr(t, ev, "(do ((i 0 (+ i 1))) ((= i 3)) (reset! counter (+ (deref counter) 1)))")
	assert.Equal(t, "3", evalStr(t, ev, "(deref counter)"))

	// do with no bindings
	assert.Equal(t, "'done", evalStr(t, ev, "(do () (#t 'done))"))

	// do building a list in reverse (common pattern)
	assert.Equal(t, "'(4 3 2 1 0)", evalStr(t, ev, `
		(do ((i 0 (+ i 1))
		     (acc '() (cons i acc)))
		    ((= i 5) acc))`))

	// simultaneous update — both vars see old values during step computation
	assert.Equal(t, "'(1 0)", evalStr(t, ev, `
		(do ((x 0 y)
		     (y 1 x))
		    ((= x 1) (list x y)))`))
}

func TestPhase2Values(t *testing.T) {
	ev := newTestEnv()

	// single value is transparent
	assert.Equal(t, "42", evalStr(t, ev, "(values 42)"))

	// call-with-values: producer returns multiple, consumer receives them
	assert.Equal(t, "3", evalStr(t, ev, "(call-with-values (lambda () (values 1 2)) +)"))
	assert.Equal(t, "'(1 2)", evalStr(t, ev, "(call-with-values (lambda () (values 1 2)) list)"))
	assert.Equal(t, "hello", evalStr(t, ev, `(call-with-values (lambda () "hello") (lambda (x) x))`))
}

func TestPhase2Ports(t *testing.T) {
	ev := newTestEnv()

	// open-input-string and read
	evalStr(t, ev, "(define p (open-input-string \"(+ 1 2)\"))")
	assert.Equal(t, "'(+ 1 2)", evalStr(t, ev, "(read p)"))

	// read multiple expressions from a port
	evalStr(t, ev, "(define p2 (open-input-string \"42 hello\"))")
	assert.Equal(t, "42", evalStr(t, ev, "(read p2)"))
	assert.Equal(t, "'hello", evalStr(t, ev, "(read p2)"))

	// eof detection
	evalStr(t, ev, "(define p3 (open-input-string \"\"))")
	assert.Equal(t, "true", evalStr(t, ev, "(eof-object? (read p3))"))
	assert.Equal(t, "true", evalStr(t, ev, "(eof-object? (eof-object))"))
	assert.Equal(t, "false", evalStr(t, ev, "(eof-object? 42)"))

	// read-char
	evalStr(t, ev, "(define p4 (open-input-string \"abc\"))")
	assert.Equal(t, "a", evalStr(t, ev, "(read-char p4)"))
	assert.Equal(t, "b", evalStr(t, ev, "(read-char p4)"))

	// peek-char doesn't consume
	evalStr(t, ev, "(define p5 (open-input-string \"xy\"))")
	assert.Equal(t, "x", evalStr(t, ev, "(peek-char p5)"))
	assert.Equal(t, "x", evalStr(t, ev, "(read-char p5)"))
	assert.Equal(t, "y", evalStr(t, ev, "(read-char p5)"))

	// input-port? predicate
	assert.Equal(t, "true", evalStr(t, ev, "(input-port? (open-input-string \"test\"))"))
	assert.Equal(t, "false", evalStr(t, ev, "(input-port? 42)"))
}

// TCO stress tests — these will blow the stack if tail calls aren't handled right
func TestPhase2TCO(t *testing.T) {
	ev := newTestEnv()

	// Named let with 1,000,000 iterations
	assert.Equal(t, "1000000", evalStr(t, ev, `
		(let loop ((i 0))
		  (if (= i 1000000)
		      i
		      (loop (+ i 1))))`))

	// Direct recursion via define
	evalStr(t, ev, `(define (count-down n)
		  (if (= n 0) 'done (count-down (- n 1))))`)
	assert.Equal(t, "'done", evalStr(t, ev, "(count-down 1000000)"))

	// Tail call in cond
	evalStr(t, ev, `(define (cond-loop n)
		  (cond ((= n 0) 'done)
		        (else (cond-loop (- n 1)))))`)
	assert.Equal(t, "'done", evalStr(t, ev, "(cond-loop 1000000)"))

	// Tail call in letrec
	evalStr(t, ev, `(define letrec-loop
		  (letrec ((f (lambda (n) (if (= n 0) 'done (f (- n 1))))))
		    f))`)
	assert.Equal(t, "'done", evalStr(t, ev, "(letrec-loop 1000000)"))
}
