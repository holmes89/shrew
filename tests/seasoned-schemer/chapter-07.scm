; Chapter 7 — Y combinator and higher-order patterns

; Y combinator (applicative-order / Z combinator for strict evaluation)
(define Y
  (lambda (f)
    ((lambda (x) (f (lambda (v) ((x x) v))))
     (lambda (x) (f (lambda (v) ((x x) v)))))))

(define factorial-gen
  (lambda (self)
    (lambda (n)
      (if (= n 0) 1 (* n (self (- n 1)))))))

; Memoization using set!
(define make-memoize
  (lambda (f)
    (let ((cache '()))
      (lambda (x)
        (let ((cached (assoc x cache)))
          (if cached
              (car (cdr cached))
              (let ((result (f x)))
                (set! cache (cons (list x result) cache))
                result)))))))

; Church numerals
(define zero  (lambda (f) (lambda (x) x)))
(define succ  (lambda (n) (lambda (f) (lambda (x) (f ((n f) x))))))
(define church->int (lambda (n) ((n (lambda (x) (+ x 1))) 0)))
