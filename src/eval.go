package main

type Scope struct {
	token *Token
	index int
}

type Word struct {
	token   *Token
	builtin Proc
}

type EvalState struct {
	scopes                Stack[*Scope]
	err                   error
	root                  *Token
	words                 map[string]Word
	values                Stack[Value]
	lastPrintedWasNewline bool
}

func (state *EvalState) pushScope(token *Token) {
	state.scopes.Push(&Scope{token, 0})
}

func newEvalState(parseState ParseState, defaultWords map[string]Word) EvalState {
	state := EvalState{
		scopes:                Stack[*Scope]{},
		err:                   parseState.err,
		root:                  parseState.root,
		words:                 defaultWords,
		values:                Stack[Value]{},
		lastPrintedWasNewline: true,
	}
	state.pushScope(state.root)
	builtins := append(Builtins, GeneratedBuiltins...)
	for _, builtin := range builtins {
		state.words[builtin.name] = Word{builtin: builtin.proc}
	}
	return state
}

func (state *EvalState) currentToken() *Token {
	scope, ok := state.scopes.Peek()
	if !ok {
		return nil
	} else if scope.index >= len(scope.token.children) {
		return nil
	}
	return &scope.token.children[scope.index]
}

func (state *EvalState) step() {
	scope, ok := state.scopes.Peek()
	if !ok {
		return
	} else if scope.index >= len(scope.token.children) {
		state.scopes.Pop()
		return
	}
	token := scope.token.children[scope.index]
	switch token.kind {
	case TokenNumber:
		state.values.Push(token.value)
		scope.index++
		return
	case TokenString:
		state.values.Push(token.value)
		scope.index++
		return
	case TokenWord:
		word, ok := state.words[token.value.text]
		if !ok {
			state.Error("undefined word: `%v`", token.value.text)
			return
		}
		if word.token != nil {
			state.pushScope(word.token)
		} else if word.builtin != nil {
			if !word.builtin(state) {
				state.Error("builtin failed: `%v`", token.value.text)
				return
			}
		} else {
			state.Error("malformed word entry: `%v`", token.value.text)
			return
		}
		scope.index++
		return
	case TokenDef:
		state.words[token.value.text] = Word{token: &token}
		scope.index++
		return
	case TokenLoop:
		val, ok := state.values.Pop()
		if !ok {
			state.Error("empty stack")
			return
		}
		if val.kind != ValueNumber {
			state.Error("loop cond should be number, got `%v`", val.kind)
			return
		}
		if val.number == 0 {
			scope.index++
		} else {
			state.pushScope(&token)
		}
		return
	}
}

func eval(parseState ParseState, defaultWords map[string]Word) (state EvalState) {
	state = newEvalState(parseState, defaultWords)
	if state.err != nil {
		return
	}
	for {
		state.step()
		if state.err != nil {
			return
		}
		if state.scopes.Len() < 1 {
			return
		}
	}
}
