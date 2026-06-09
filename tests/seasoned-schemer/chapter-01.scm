; Chapter 1 — revisiting recursion with multi-star variants

(define rember*
  (lambda (a l)
    (cond
      ((null? l) '())
      ((pair? (car l))
       (cons (rember* a (car l))
             (rember* a (cdr l))))
      ((eq? a (car l)) (rember* a (cdr l)))
      (else (cons (car l) (rember* a (cdr l)))))))

(define insertR*
  (lambda (new old l)
    (cond
      ((null? l) '())
      ((pair? (car l))
       (cons (insertR* new old (car l))
             (insertR* new old (cdr l))))
      ((eq? old (car l))
       (cons old (cons new (insertR* new old (cdr l)))))
      (else (cons (car l) (insertR* new old (cdr l)))))))

(define occur*
  (lambda (a l)
    (cond
      ((null? l) 0)
      ((pair? (car l))
       (+ (occur* a (car l))
          (occur* a (cdr l))))
      ((eq? a (car l))
       (+ 1 (occur* a (cdr l))))
      (else (occur* a (cdr l))))))

(define subst*
  (lambda (new old l)
    (cond
      ((null? l) '())
      ((pair? (car l))
       (cons (subst* new old (car l))
             (subst* new old (cdr l))))
      ((eq? old (car l))
       (cons new (subst* new old (cdr l))))
      (else (cons (car l) (subst* new old (cdr l)))))))

(define member*
  (lambda (a l)
    (cond
      ((null? l) #f)
      ((pair? (car l))
       (or (member* a (car l))
           (member* a (cdr l))))
      (else (or (eq? a (car l))
                (member* a (cdr l)))))))

(define leftmost
  (lambda (l)
    (cond
      ((pair? (car l)) (leftmost (car l)))
      (else (car l)))))
