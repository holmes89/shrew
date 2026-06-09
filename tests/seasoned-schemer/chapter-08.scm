; Chapter 8 — streams and generators (without call/cc)
; call/cc chapters (9-12) are deferred

(define make-generator
  (lambda (from)
    (let ((current from))
      (lambda ()
        (let ((val current))
          (set! current (+ current 1))
          val)))))

; Lazy streams via thunks — represented as (head . thunk-for-tail)
; using (list head thunk) so cdr is accessible via car/cdr
(define stream-cons (lambda (x thunk) (list x thunk)))
(define stream-car  car)
(define stream-cdr  (lambda (s) ((car (cdr s)))))
(define stream-null? null?)

(define integers-from
  (lambda (n)
    (stream-cons n (lambda () (integers-from (+ n 1))))))

(define stream-take
  (lambda (s n)
    (if (or (stream-null? s) (= n 0))
        '()
        (cons (stream-car s)
              (stream-take (stream-cdr s) (- n 1))))))

(define stream-filter
  (lambda (pred s)
    (cond
      ((stream-null? s) '())
      ((pred (stream-car s))
       (stream-cons (stream-car s)
                    (lambda () (stream-filter pred (stream-cdr s)))))
      (else (stream-filter pred (stream-cdr s))))))
