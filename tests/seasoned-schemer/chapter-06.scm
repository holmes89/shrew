; Chapter 6 — state with set!

(define make-counter
  (lambda ()
    (let ((n 0))
      (lambda ()
        (set! n (+ n 1))
        n))))

(define make-adder
  (lambda (start)
    (let ((total start))
      (lambda (x)
        (set! total (+ total x))
        total))))

; Stack using message-passing with tagged list messages: (push val) or symbol
(define make-stack
  (lambda ()
    (let ((data '()))
      (lambda (msg)
        (cond
          ((eq? msg 'pop)
           (let ((top (car data)))
             (set! data (cdr data))
             top))
          ((eq? msg 'peek) (car data))
          ((eq? msg 'empty?) (null? data))
          ((and (pair? msg) (eq? (car msg) 'push))
           (set! data (cons (car (cdr msg)) data)))
          (else (error "unknown message")))))))

(define stack-push
  (lambda (stack val)
    (stack (list 'push val))))
