package repltests

var chapterSS3 = []chaptertest{
	{`(evens-and-odds 4)`, `'(true false)`},
	{`(evens-and-odds 7)`, `'(false true)`},
	{`(flatten-deep '(1 (2 (3 4) 5) (6)))`, `'(1 2 3 4 5 6)`},
	{`(flatten-deep '((a b) (c (d (e)))))`, `'(a b c d e)`},
}
