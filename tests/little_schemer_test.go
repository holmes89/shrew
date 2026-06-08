package repltests

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/holmes89/shrew/core"
	"github.com/holmes89/shrew/env"
	"github.com/holmes89/shrew/repl"
	. "github.com/holmes89/shrew/types"
	"github.com/stretchr/testify/suite"
)

type LittleSchemerTestSuite struct {
	suite.Suite
	ev EnvType
}

// In order for 'go test' to run this suite, we need to create
// a normal test function and pass our suite to suite.Run
func TestLittleSchemerTestSuite(t *testing.T) {
	suite.Run(t, new(LittleSchemerTestSuite))
}

func (suite *LittleSchemerTestSuite) SetupSuite() {
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

func (suite *LittleSchemerTestSuite) LoadLibrary(i int) {
	b, err := os.ReadFile(fmt.Sprintf("little-schemer/chapter-%d.scm", i))
	if err != nil {
		suite.FailNow(err.Error())
	}
	for _, lib := range strings.Split(string(b), "\n\n") {
		lib = strings.TrimSpace(lib)
		if lib == "" {
			continue
		}
		if _, err := repl.Repl(lib, suite.ev); err != nil {
			if err.Error() == "<empty line>" {
				continue
			}
			fmt.Printf("invalid input: %s\n", lib)
			suite.FailNow(err.Error())
		}
	}

}

func (suite *LittleSchemerTestSuite) TestRepl() {
	tt := []booktests{
		// Chapter 1 is just basic scheme commands, maybe add without repl load?
		{Chapter: 2, Tests: chapter2},
		{Chapter: 3, Tests: chapter3},
		{Chapter: 4, Tests: chapter4},
		{Chapter: 5, Tests: chapter5},
		{Chapter: 6, Tests: chapter6},
		{Chapter: 7, Tests: chapter7},
		{Chapter: 8, Tests: chapter8},
		{Chapter: 9, Tests: chapter9},
		{Chapter: 10, Tests: chapter10},
	}

	for _, t := range tt {
		suite.LoadLibrary(t.Chapter)
		suite.process(t.Tests)
	}
}

func (suite *LittleSchemerTestSuite) process(tests []chaptertest) {
	for _, t := range tests {
		res, err := repl.Repl(t.Command, suite.ev)
		suite.NoError(err)
		suite.Equal(t.Result, res, t.Command)
	}
}

type booktests struct {
	Chapter int
	Tests   []chaptertest
}

type chaptertest struct {
	Command string
	Result  string
}
