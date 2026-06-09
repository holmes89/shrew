; Chapter 5 — where/when/unless and let in function bodies

; Build on named let
(define count-atoms
  (lambda (l)
    (let loop ((l l) (acc 0))
      (cond
        ((null? l) acc)
        ((pair? (car l)) (loop (cdr l) (+ acc (count-atoms (car l)))))
        (else (loop (cdr l) (+ acc 1)))))))

(define depth
  (lambda (lst)
    (cond
      ((null? lst) 0)
      ((pair? (car lst))
       (max (+ 1 (depth (car lst)))
            (depth (cdr lst))))
      (else (max 1 (depth (cdr lst)))))))

; Using do loop for list operations
(define list-sum
  (lambda (lst)
    (do ((l lst (cdr l)) (s 0 (+ s (car l))))
        ((null? l) s))))

(define list-product
  (lambda (lst)
    (do ((l lst (cdr l)) (p 1 (* p (car l))))
        ((null? l) p))))

(define iota
  (lambda (n)
    (do ((i (- n 1) (- i 1))
         (acc '() (cons i acc)))
        ((< i 0) acc))))
