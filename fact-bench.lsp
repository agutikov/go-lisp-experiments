
(defun ffact (n) (if (<= n 1) 1 (* n (ffact (- n 1)))))

(define fact (lambda (n) (if (<= n 1) 1 (* n (fact (- n 1))))))


(defun cmp (x y) (if (= x y) "OK" "FALSE"))



(set! enable-print-elapsed t)
(set! enable-trace t)



(define X (__fact 100)) ; bultin, non-recursive

(cmp X (__fact_r 100)) ; builtin, recursive

(cmp X (fact 10)) ; check cmp itself

(cmp X (fact 100)) ; define + lambda

(cmp X (ffact 100)) ; defun is faster than define + lambda if it is recursive






()

