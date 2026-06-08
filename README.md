# Shrew

A Scheme-flavoured Lisp interpreter written in Go, used as a learning project for compilers, programming languages, and functional programming.

![Shrew](./shrew.png)

### Resources
- [Root of LISP](http://www.paulgraham.com/rootsoflisp.html)
- [Writing an Interpreter in Go](https://interpreterbook.com/)
- [Go Lexer](https://golang.org/src/text/template/parse/lex.go)
- [Rob Pike Lisp 1.5](https://github.com/robpike/lisp)
- [McCarthy Paper](http://www-formal.stanford.edu/jmc/recursive/recursive.html)
- [Little Schemer Files](https://github.com/pkrumins/the-little-schemer)
- [MAL (Make-A-Lisp)](https://github.com/kanaka/mal/blob/master/process/guide.md#step1)

## Roadmap

See [shrew-roadmap.md](../joel.holmes.haus/shrew-roadmap.md) for the full evolution plan covering:
- **Track 1** — Core Scheme completeness (R5RS/R7RS parity)
- **Track 2** — Clojure-ish Go interop (embed Shrew in any Go service)
- **Track 3** — WebAssembly VM (browser REPL via go-app)
- **Track 4** — Logic programming (miniKanren / Reasoned Schemer)

## Current Features

- Lists, vectors, hash maps, atoms
- Tail-call optimized eval loop
- `lambda` / `λ` with multi-body support
- `define` with function shorthand: `(define (f x) body...)`
- `let`, `let*`, `letrec`, `letrec*`, named `let`
- `begin`, `do`, `when`, `unless`
- `cond` with `else`, `if`
- `and`, `or` (short-circuit special forms)
- `quote`, `quasiquote`, `unquote`, `splice-unquote`
- `defmacro`, `macroexpand`
- `try`/`catch*` error handling
- `#t`/`#f` booleans, `#\char` character literals
- Float numbers (`3.14`, `-2.5`)
- `load-file` for loading `.scm` files
- REPL with readline history and multi-line support

### Stdlib Highlights

| Category | Functions |
|---|---|
| Math | `+` `-` `*` `/` `abs` `max` `min` `modulo` `remainder` `quotient` `floor` `ceiling` `round` `truncate` `sqrt` `expt` `exact->inexact` `inexact->exact` |
| String | `string-length` `string-append` `substring` `string-ref` `string-upcase` `string-downcase` `string-contains` `string-split` `string-join` `string-trim` `string->number` `number->string` `string->list` `list->string` |
| List | `car` `cdr` `cons` `list` `append` `reverse` `map` `filter` `for-each` `fold-left` `fold-right` `reduce` `any` `every` `flatten` `list-ref` `list-tail` `sort` `apply` |
| Vector | `vector` `make-vector` `vector-ref` `vector-set!` `vector-length` `vector->list` `list->vector` |
| Predicate | `null?` `pair?` `list?` `atom?` `number?` `integer?` `float?` `boolean?` `string?` `symbol?` `keyword?` `procedure?` `vector?` `zero?` `even?` `odd?` `positive?` `negative?` |
| IO | `display` `newline` `write` `print` `println` `slurp` |
| Misc | `gensym` `error` `type-of` `str` `count` `length` `read-string` |

## Building

```bash
make build   # builds to bin/shrew
make test    # runs all tests
```

## Running

```bash
./bin/shrew           # interactive REPL
```

# Limitations / Todos

- [ ] `call/cc` (call-with-current-continuation)
- [ ] `define-syntax` / `syntax-rules` (hygienic macros)
- [ ] Port/IO system (`open-input-file`, `open-output-file`, `read`)
- [ ] Tail-call optimization audit for all paths (named `let`, deep `cond`)
- [ ] Go interop bridge (`(require 'http)`, embedding API)
- [ ] WebAssembly target
- [ ] miniKanren logic programming

# Little Schemer Progress

- [x] Chapter 1
- [x] Chapter 2
- [x] Chapter 3
- [x] Chapter 4
- [x] Chapter 5
- [x] Chapter 6
- [x] Chapter 7
- [x] Chapter 8
- [x] Chapter 9
- [x] Chapter 10
