package repltests

var chapterSS4 = []chaptertest{
	{`(let ((env (env-extend 'x 42 (env-extend 'y 7 (make-env)))))
	   (env-lookup 'x env))`,
		"42"},
	{`(let ((env (env-extend 'x 42 (env-extend 'y 7 (make-env)))))
	   (env-lookup 'y env))`,
		"7"},
	{`(let ((env (env-extend 'x 42 (make-env))))
	   (env-lookup 'z env))`,
		"false"},
}
