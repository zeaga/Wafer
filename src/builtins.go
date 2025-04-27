package main

import "os"

type Proc func(state *EvalState) bool

type Builtin struct {
	category string
	name     string
	inputs   string
	outputs  string
	proc     Proc
}

var Builtins = []Builtin{
	{category: "io", name: "runstring", inputs: "1s", outputs: "0", proc: func(state *EvalState) bool {
		script, ok := state.pop1s()
		if !ok {
			return false
		}
		token := state.currentToken()
		if token == nil {
			return false
		}
		lexState := lex(token.file, script)
		if lexState.err != nil {
			return false
		}
		parseState := parse(lexState)
		if parseState.err != nil {
			return false
		}
		state.scopes.Push(&Scope{parseState.root, 0})
		return true
	}},
	{category: "io", name: "loadfile", inputs: "1s", outputs: "1s", proc: func(state *EvalState) bool {
		filename, ok := state.pop1s()
		if !ok {
			return false
		}
		file, err := os.ReadFile(filename)
		if err != nil {
			return false
		}
		return state.push1s(string(file))
	}},
	{category: "io", name: "runfile", inputs: "1s", outputs: "0", proc: func(state *EvalState) bool {
		filename, ok := state.pop1s()
		if !ok {
			return false
		}
		file, err := os.ReadFile(filename)
		if err != nil {
			return false
		}
		lexState := lex(filename, string(file))
		if lexState.err != nil {
			return false
		}
		parseState := parse(lexState)
		if parseState.err != nil {
			return false
		}
		state.scopes.Push(&Scope{parseState.root, 0})
		return true
	}},
}
