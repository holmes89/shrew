package repltests

var chapterSS8 = []chaptertest{
	{`(let ((g (make-generator 1)))
	   (list (g) (g) (g)))`,
		`'(1 2 3)`},
	{`(stream-take (integers-from 1) 5)`,
		`'(1 2 3 4 5)`},
	{`(stream-take
	   (stream-filter odd? (integers-from 1))
	   5)`,
		`'(1 3 5 7 9)`},
}
