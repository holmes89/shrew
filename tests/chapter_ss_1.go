package repltests

var chapterSS1 = []chaptertest{
	{`(rember* 'cup '((coffee) cup ((tea) cup) (and (hick)) cup))`,
		`'((coffee) ((tea)) (and (hick)))`},
	{`(rember* 'sauce '(((tomato sauce)) ((bean) sauce) (and ((flying)) sauce)))`,
		`'(((tomato)) ((bean)) (and ((flying))))`},
	{`(insertR* 'roast 'chuck '((how much (wood)) could ((a (wood) chuck)) (((chuck))) (if (a) ((wood chuck))) could chuck wood))`,
		`'((how much (wood)) could ((a (wood) chuck roast)) (((chuck roast))) (if (a) ((wood chuck roast))) could chuck roast wood)`},
	{`(occur* 'banana '((banana) (split ((banana) ice) (cream (banana)) sherbet)))`,
		"3"},
	{`(subst* 'orange 'banana '((banana) (split ((banana) ice) (cream (banana)) sherbet)))`,
		`'((orange) (split ((orange) ice) (cream (orange)) sherbet))`},
	{`(member* 'chips '((potato) (chips ((with) fish) (chips))))`,
		"true"},
	{`(leftmost '((potato) (chips ((with) fish) (chips))))`,
		"'potato"},
}
