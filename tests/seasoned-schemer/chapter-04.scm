; Chapter 4 — association tables (lookup tables)

(define lookup-in-entry
  (lambda (name entry entry-f)
    (cond
      ((null? (car entry)) (entry-f name))
      ((eq? name (car (car entry)))
       (car (cdr entry)))
      (else
       (lookup-in-entry
         name
         (list (cdr (car entry))
               (cdr (cdr entry)))
         entry-f)))))

(define extend-table
  (lambda (entry table)
    (cons entry table)))

(define lookup-in-table
  (lambda (name table table-f)
    (cond
      ((null? table) (table-f name))
      (else
       (lookup-in-entry
         name
         (car table)
         (lambda (name)
           (lookup-in-table name (cdr table) table-f)))))))

; Association-list based simple environment
(define make-env
  (lambda ()
    '()))

(define env-lookup
  (lambda (key env)
    (let ((pair (assoc key env)))
      (if pair (car (cdr pair)) #f))))

(define env-extend
  (lambda (key val env)
    (cons (list key val) env)))
