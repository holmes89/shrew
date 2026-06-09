package repltests

var chapterSS7 = []chaptertest{
	{`((Y factorial-gen) 5)`, "120"},
	{`((Y factorial-gen) 0)`, "1"},
	{`(let ((double (make-memoize (lambda (n) (* n 2)))))
	   (list (double 3) (double 5) (double 3)))`,
		`'(6 10 6)`},
	{`(church->int zero)`, "0"},
	{`(church->int (succ (succ (succ zero))))`, "3"},
}
