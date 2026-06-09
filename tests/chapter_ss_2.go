package repltests

var chapterSS2 = []chaptertest{
	{`(multirember 'tuna '(shrimp salad tuna salad and tuna))`,
		`'(shrimp salad salad and)`},
	{`(multirember 'tuna '())`,
		`'()`},
	{`(multirember&co 'tuna '(strawberries tuna and swordfish) list)`,
		`'((strawberries and swordfish) (tuna))`},
	{`(multiinsertLR 'new 'oldL 'oldR '(oldL a oldR b oldL c oldR))`,
		`'(new oldL a oldR new b new oldL c oldR new)`},
	{`(multiinsertLR&co 'new 'oldL 'oldR '(oldL a oldR) list)`,
		`'((new oldL a oldR new) 1 1)`},
}
