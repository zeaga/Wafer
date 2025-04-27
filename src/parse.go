package main

import (
	"strconv"
)

type ValueKind int

const (
	ValueNumber ValueKind = iota
	ValueText
)

func (kind ValueKind) String() string {
	switch kind {
	case ValueNumber:
		return "number"
	case ValueText:
		return "text"
	}
	return "unknown"
}

type Value struct {
	kind   ValueKind
	number float64
	text   string
}

type TokenKind int

const (
	TokenRoot TokenKind = iota
	TokenNumber
	TokenString
	TokenWord
	TokenDef
	TokenLoop
)

func (kind TokenKind) String() string {
	switch kind {
	case TokenRoot:
		return "root"
	case TokenNumber:
		return "number"
	case TokenString:
		return "string"
	case TokenWord:
		return "word"
	case TokenDef:
		return "definition"
	case TokenLoop:
		return "loop"
	}
	return "unknown"
}

type Token struct {
	kind     TokenKind
	parent   *Token
	children []Token
	value    Value
	file     string
	line     int
	col      int
}

type ParseState struct {
	lexemes []Lexeme
	index   int
	file    string
	line    int
	col     int
	err     error
	scopes  Stack[*Token]
	root    *Token
}

func newParseState(lexState LexState) ParseState {
	root := Token{kind: TokenRoot}
	return ParseState{
		lexemes: lexState.lexemes,
		index:   0,
		file:    lexState.file,
		line:    -1,
		col:     -1,
		err:     lexState.err,
		scopes:  Stack[*Token]{},
		root:    &root,
	}
}

func (state *ParseState) addToken(kind TokenKind) *Token {
	token := Token{
		kind: kind,
		file: state.file,
		line: state.line,
		col:  state.col,
	}
	top, ok := state.scopes.Peek()
	if ok {
		token.parent = top
		top.children = append(top.children, token)
		return &top.children[len(top.children)-1]
	} else {
		token.parent = state.root
		state.root.children = append(state.root.children, token)
		return &state.root.children[len(state.root.children)-1]
	}
}

func (state *ParseState) handleNumber() {
	lexeme := state.lexemes[state.index]
	val, err := strconv.ParseFloat(lexeme.text, 64)
	if err != nil {
		state.Error("malformed number `%v`", lexeme.text)
		return
	}
	token := state.addToken(TokenNumber)
	token.value = Value{kind: ValueNumber, number: val}
	state.index++
}

func (state *ParseState) handleDefBegin() {
	state.index++ // move past ':'
	if state.index >= len(state.lexemes) {
		state.Error("expected word after ':', got eof")
		return
	}

	word := state.lexemes[state.index]
	if word.kind != LexemeWord {
		word.Error("expected word after ':', got `%v`", word.kind)
		return
	}

	token := state.addToken(TokenDef)
	token.value = Value{kind: ValueText, text: word.text}
	state.scopes.Push(token)
	state.index++
}

func (state *ParseState) handleDefEnd() {
	top, ok := state.scopes.Pop()
	if !ok {
		state.Error("unexpected end of definition")
		return
	} else if top.kind != TokenDef {
		state.Error("expected end of definition, got `%v`", top.kind)
		return
	}
	state.index++
}

func (state *ParseState) handleLoopEnd() {
	top, ok := state.scopes.Pop()
	if !ok {
		state.Error("unexpected end of loop")
		return
	} else if top.kind != TokenLoop {
		state.Error("expected end of loop, got `%v`", top.kind)
		return
	}
	state.index++
}

func (state *ParseState) step() {
	lexeme := state.lexemes[state.index]
	state.line = lexeme.line
	state.col = lexeme.col
	switch lexeme.kind {
	case LexemeNumber:
		state.handleNumber()
	case LexemeString:
		state.addToken(TokenString).value = Value{kind: ValueText, text: lexeme.text}
		state.index++
	case LexemeWord:
		state.addToken(TokenWord).value = Value{kind: ValueText, text: lexeme.text}
		state.index++
	case LexemeDefBegin:
		state.handleDefBegin()
	case LexemeDefEnd:
		state.handleDefEnd()
	case LexemeLoopBegin:
		state.scopes.Push(state.addToken(TokenLoop))
		state.index++
	case LexemeLoopEnd:
		state.handleLoopEnd()
	default:
		state.Error("unexpected lexeme in parsing stage: `%v`", lexeme.text)
	}
}

func parse(lexState LexState) (state ParseState) {
	state = newParseState(lexState)
	if state.err != nil {
		return
	}
	for state.index < len(state.lexemes) && state.err == nil {
		state.step()
	}
	return
}
