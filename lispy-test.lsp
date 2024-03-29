
(define circle-area (lambda (r) (* pi (* r r))))
(circle-area 3)
;28.274333877

(define fact (lambda (n) (if (<= n 1) 1 (* n (fact (- n 1))))))
(fact 10)
;3628800

(fact 100)
;9332621544394415268169923885626670049071596826438162146859296389521759999322991
;5608941463976156518286253697920827223758251185210916864000000000000000000000000

(circle-area (fact 10))
;4.1369087198e+13

(define first car)
(define rest cdr)
(define count (lambda (item L) (if L (+ (equal? item (first L)) (count item (rest L))) 0)))
(count 0 (list 0 1 2 3 0 0))
;3

(count (quote the) (quote (the more the merrier the bigger the better)))
;4

(define twice (lambda (x) (* 2 x)))
(twice 5)
;10

(define repeat (lambda (f) (lambda (x) (f (f x)))))
((repeat twice) 10)
;40

((repeat (repeat twice)) 10)
;160

((repeat (repeat (repeat twice))) 10)
;2560

((repeat (repeat (repeat (repeat twice)))) 10)
;655360

(pow 2 16)
;65536.0

(define fib (lambda (n) (if (< n 2) 1 (+ (fib (- n 1)) (fib (- n 2))))))
(define range (lambda (a b) (if (= a b) (quote ()) (cons a (range (+ a 1) b)))))
(range 0 10)
;(0 1 2 3 4 5 6 7 8 9)
(map fib (range 0 10))
;(1 1 2 3 5 8 13 21 34 55)
(map fib (range 0 27))
