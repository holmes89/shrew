; Chapter 2 — mirrors and multirember with let/letrec

(define multirember
  (lambda (a lat)
    (cond
      ((null? lat) '())
      ((eq? a (car lat)) (multirember a (cdr lat)))
      (else (cons (car lat) (multirember a (cdr lat)))))))

(define multirember&co
  (lambda (a lat col)
    (cond
      ((null? lat)
       (col '() '()))
      ((eq? (car lat) a)
       (multirember&co a (cdr lat)
                       (lambda (newlat seen)
                         (col newlat (cons (car lat) seen)))))
      (else
       (multirember&co a (cdr lat)
                       (lambda (newlat seen)
                         (col (cons (car lat) newlat) seen)))))))

(define multiinsertLR
  (lambda (new oldL oldR lat)
    (cond
      ((null? lat) '())
      ((eq? (car lat) oldL)
       (cons new (cons oldL (multiinsertLR new oldL oldR (cdr lat)))))
      ((eq? (car lat) oldR)
       (cons oldR (cons new (multiinsertLR new oldL oldR (cdr lat)))))
      (else
       (cons (car lat) (multiinsertLR new oldL oldR (cdr lat)))))))

(define multiinsertLR&co
  (lambda (new oldL oldR lat col)
    (cond
      ((null? lat)
       (col '() 0 0))
      ((eq? (car lat) oldL)
       (multiinsertLR&co new oldL oldR (cdr lat)
                         (lambda (newlat cL cR)
                           (col (cons new (cons oldL newlat))
                                (+ 1 cL) cR))))
      ((eq? (car lat) oldR)
       (multiinsertLR&co new oldL oldR (cdr lat)
                         (lambda (newlat cL cR)
                           (col (cons oldR (cons new newlat))
                                cL (+ 1 cR)))))
      (else
       (multiinsertLR&co new oldL oldR (cdr lat)
                         (lambda (newlat cL cR)
                           (col (cons (car lat) newlat) cL cR)))))))
