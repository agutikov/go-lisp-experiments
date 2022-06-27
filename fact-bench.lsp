
(define fact (lambda (n) (if (<= n 1) 1 (* n (fact (- n 1))))))

(set! enable-print-elapsed t)
(set! enable-trace t)


;(fact 10)

(__fact 100)
(__fact_r 100)

(fact 100)



;(define run-go-fact (lambda (n) (< (__fact n) 0)))

;(define run-fact (lambda (n) (< (fact n) 0)))



