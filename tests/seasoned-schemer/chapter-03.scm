; Chapter 3 — looking at arguments / letrec

(define even?
  (lambda (n)
    (= (remainder n 2) 0)))

; even? and odd? with mutual letrec
(define evens-and-odds
  (letrec
    ((my-even? (lambda (n) (if (= n 0) #t (my-odd?  (- n 1)))))
     (my-odd?  (lambda (n) (if (= n 0) #f (my-even? (- n 1))))))
    (lambda (n)
      (list (my-even? n) (my-odd? n)))))

; flatten using letrec — process car into acc first, then cdr, then reverse
(define flatten-deep
  (lambda (l)
    (letrec
      ((go (lambda (l acc)
             (cond
               ((null? l) acc)
               ((pair? (car l))
                (go (cdr l) (go (car l) acc)))
               (else (go (cdr l) (cons (car l) acc)))))))
      (reverse (go l '())))))
