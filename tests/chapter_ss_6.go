package repltests

var chapterSS6 = []chaptertest{
	{`(let ((c (make-counter)))
	   (list (c) (c) (c)))`,
		`'(1 2 3)`},
	{`(let ((add (make-adder 10)))
	   (list (add 5) (add 3) (add 2)))`,
		`'(15 18 20)`},
	{`(let ((s (make-stack)))
	   (stack-push s 1)
	   (stack-push s 2)
	   (stack-push s 3)
	   (list (s 'pop) (s 'pop) (s 'empty?)))`,
		`'(3 2 false)`},
}
