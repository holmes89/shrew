package repltests

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/holmes89/shrew/core"
	"github.com/holmes89/shrew/env"
	"github.com/holmes89/shrew/lexer"
	"github.com/holmes89/shrew/repl"
	. "github.com/holmes89/shrew/types"
	"github.com/stretchr/testify/suite"
)

type SeasonedSchemerTestSuite struct {
	suite.Suite
	ev EnvType
}

func TestSeasonedSchemerTestSuite(t *testing.T) {
	suite.Run(t, new(SeasonedSchemerTestSuite))
}

func (suite *SeasonedSchemerTestSuite) SetupSuite() {
	suite.ev = env.DefaultEnv()
	for k, v := range core.NS {
		suite.ev.Set(k, Func{Fn: v})
	}
	suite.ev.Set(Symbol{Val: "eval"}, Func{
		Fn: func(a []Expression) (Expression, error) {
			return repl.Eval(a, suite.ev)
		},
	})
}

func (suite *SeasonedSchemerTestSuite) LoadChapter(i int) {
	b, err := os.ReadFile(fmt.Sprintf("seasoned-schemer/chapter-%02d.scm", i))
	if err != nil {
		suite.FailNow(err.Error())
	}
	l := lexer.NewLexer(strings.NewReader(string(b)))
	for {
		expr, err := l.ReadExpr()
		if err != nil {
			if strings.Contains(err.Error(), "<empty line>") || strings.Contains(err.Error(), "EOF") {
				break
			}
			suite.FailNow(fmt.Sprintf("chapter %02d parse error: %v", i, err))
		}
		if expr == nil {
			continue
		}
		if _, err := repl.Eval(expr, suite.ev); err != nil {
			suite.FailNow(fmt.Sprintf("chapter %02d eval error: %v", i, err))
		}
	}
}

func (suite *SeasonedSchemerTestSuite) TestRepl() {
	tt := []struct {
		Chapter int
		Tests   []chaptertest
	}{
		{1, chapterSS1},
		{2, chapterSS2},
		{3, chapterSS3},
		{4, chapterSS4},
		{5, chapterSS5},
		{6, chapterSS6},
		{7, chapterSS7},
		{8, chapterSS8},
	}

	for _, t := range tt {
		suite.LoadChapter(t.Chapter)
		suite.processSS(t.Tests)
	}
}

func (suite *SeasonedSchemerTestSuite) processSS(tests []chaptertest) {
	for _, t := range tests {
		res, err := repl.Repl(t.Command, suite.ev)
		suite.NoError(err)
		suite.Equal(t.Result, res, t.Command)
	}
}
