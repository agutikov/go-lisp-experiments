
(set! print-elapsed t)

(define fact (lambda (n) (if (<= n 1) 1 (* n (fact (- n 1))))))

(fact 10)

(fact 100)

(fact 1000)
