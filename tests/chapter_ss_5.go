package repltests

var chapterSS5 = []chaptertest{
	{`(count-atoms '((a b) c (d (e f))))`, "6"},
	{`(count-atoms '())`, "0"},
	{`(depth '((a (b)) c))`, "3"},
	{`(list-sum '(1 2 3 4 5))`, "15"},
	{`(list-product '(1 2 3 4 5))`, "120"},
	{`(iota 5)`, `'(0 1 2 3 4)`},
	{`(iota 0)`, `'()`},
}
