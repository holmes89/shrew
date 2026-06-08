package main

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/chzyer/readline"
	"github.com/holmes89/shrew/core"
	. "github.com/holmes89/shrew/env"
	"github.com/holmes89/shrew/lexer"
	. "github.com/holmes89/shrew/repl"
	. "github.com/holmes89/shrew/types"
)

var repl_env = DefaultEnv()

func init() {
	for k, v := range core.NS {
		repl_env.Set(k, Func{Fn: v})
	}
	repl_env.Set(Symbol{Val: "eval"}, Func{
		Fn: func(a []Expression) (Expression, error) {
			return Eval(a, repl_env)
		},
	})
	repl_env.Set(Symbol{Val: "load-file"}, Func{
		Fn: func(a []Expression) (Expression, error) {
			if len(a) != 1 {
				return nil, errors.New("load-file arity mismatch expected: 1")
			}
			b, err := os.ReadFile(a[0].(string))
			if err != nil {
				return nil, errors.New("unable to read file")
			}
			buf := bytes.NewBuffer(b)
			scanner := bufio.NewScanner(buf)
			var exp Expression
			for scanner.Scan() {
				exp, err = lexer.Read(bytes.NewBuffer(scanner.Bytes()))
				if err != nil {
					return nil, err
				}
				exp, err = Eval(exp, repl_env)
				if err != nil {
					return nil, err
				}
			}
			return exp, err
		},
	})
	repl_env.Set(Symbol{Val: "*ARGV*"}, List{})
}

var (
	defaultPrompt  = "shrew=> "
	continuePrompt = "... "
)

func main() {

	rl, err := readline.NewEx(&readline.Config{
		Prompt:                 defaultPrompt,
		HistoryFile:            "/tmp/shrew",
		DisableAutoSaveHistory: true,
	})
	rl.SetVimMode(false)

	if err != nil {
		panic(err)
	}
	defer rl.Close()

	var cmds []string
	for {
		line, err := rl.Readline()
		if err != nil {
			break
		}
		line = strings.TrimSpace(line)
		if len(line) == 0 {
			continue
		}
		cmds = append(cmds, line)
		cmd := strings.Join(cmds, " ")

		lparenCount := strings.Count(cmd, "(")
		rparenCount := strings.Count(cmd, ")")

		if lparenCount > rparenCount {
			rl.SetPrompt(continuePrompt)
			continue
		}
		if rparenCount > lparenCount {
			rl.SetPrompt(defaultPrompt)
			rl.SaveHistory(cmd)
			cmds = cmds[:0]
			fmt.Printf("Error: mismatch paren\n")
			continue
		}

		res, err := Repl(cmd, repl_env)
		rl.SetPrompt(defaultPrompt)
		rl.SaveHistory(cmd)
		cmds = cmds[:0]
		if err != nil {
			if err.Error() == "<empty line>" {
				continue
			}
			fmt.Printf("Error: %v\n", err)

			continue
		}
		fmt.Printf("%v\n", res)
	}
}
