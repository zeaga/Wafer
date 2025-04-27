package main

type LexemeKind int

const (
	LexemeNumber LexemeKind = iota
	LexemeString
	LexemeWord
	LexemeDefBegin
	LexemeDefEnd
	LexemeLoopBegin
	LexemeLoopEnd
)

func (kind LexemeKind) String() string {
	switch kind {
	case LexemeNumber:
		return "number"
	case LexemeString:
		return "string"
	case LexemeWord:
		return "word"
	case LexemeDefBegin:
		return ":"
	case LexemeDefEnd:
		return ";"
	case LexemeLoopBegin:
		return "{"
	case LexemeLoopEnd:
		return "}"
	}
	return "unknown"
}

type Lexeme struct {
	state *LexState
	kind  LexemeKind
	text  string
	file  string
	line  int
	col   int
}

type LexState struct {
	file        string
	script      string
	index       int
	line        int
	lastLineIdx int
	err         error
	lexemes     []Lexeme
}

func newLexState(file string, script string) LexState {
	return LexState{
		file:        file,
		script:      script,
		index:       0,
		line:        0,
		lastLineIdx: 0,
		err:         nil,
		lexemes:     make([]Lexeme, 0),
	}
}

func (state *LexState) addLexeme(kind LexemeKind, text string, start int) *Lexeme {
	lexeme := Lexeme{
		state: state,
		kind:  kind,
		text:  text,
		file:  state.file,
		line:  state.line,
		col:   start - state.lastLineIdx,
	}
	state.lexemes = append(state.lexemes, lexeme)
	return &lexeme
}

func (state *LexState) handleSingleChar() bool {
	lexeme := LexemeKind(-1)
	c := state.script[state.index]
	switch c {
	case ':':
		lexeme = LexemeDefBegin
	case ';':
		lexeme = LexemeDefEnd
	case '{':
		lexeme = LexemeLoopBegin
	case '}':
		lexeme = LexemeLoopEnd
	}
	if lexeme >= 0 {
		state.addLexeme(lexeme, string(c), state.index)
		state.index++
		return true
	}
	return false
}

func (state *LexState) handleString() bool {
	c := state.script[state.index]
	if c != '"' {
		return false
	}
	start := state.index - 1
	state.index++
	text := ""
	for state.index < len(state.script) {
		c = state.script[state.index]
		if c == '"' {
			lexeme := state.addLexeme(LexemeString, text, start)
			state.index++
			if state.index < len(state.script) && !isWhitespace(state.script[state.index]) && state.script[state.index] != '\n' {
				return lexeme.Error("expected whitespace after string")
			}
			return true
		} else if c == '\n' {
			state.index = start
			return state.Error("unexpected newline in string")
		} else if c == '\\' {
			state.index++
			if state.index >= len(state.script) {
				state.index = start
				return state.Error("unexpected eof after escape character")
			}
			switch state.script[state.index] {
			case 't':
				c = '\t'
			case 'r':
				c = '\r'
			case 'n':
				c = '\n'
			case '"':
				c = '"'
			case '\\':
				c = '\\'
			default:
				state.index = start
				return state.Error("invalid escape sequence `\\%v`", c)
			}
		}
		text += string(c)
		state.index++
	}
	state.index = start
	return state.Error("unexpected eof in string")
}

func (state *LexState) handleNumber() bool {
	c := state.script[state.index]
	// make sure it starts with -, +, or a digit
	if (c != '-' && c != '+') && !isDigit(c) {
		return false
	}
	// if it is -/+, make sure the next char exists
	if (c == '-' || c == '+') && state.index+1 >= len(state.script) {
		return false
	}
	start := state.index
	if c == '-' || c == '+' {
		// no need to bounds check here because of the one earlier
		state.index++
		c = state.script[state.index]
		// since it is -/+, make sure the next char is digital
		if !isDigit(c) {
			// whoops, let's backpedal. this is a word.
			state.index = start
			return false
		}
	}
	for state.index < len(state.script) && isDigit(state.script[state.index]) {
		state.index++
	}
	if state.index < len(state.script) && state.script[state.index] == '.' {
		state.index++
		hasSecondPart := false
		for state.index < len(state.script) && isDigit(state.script[state.index]) {
			state.index++
			hasSecondPart = true
		}
		if !hasSecondPart {
			num := state.script[start:state.index]
			state.index = start
			return state.Error("malformed number `%v`", num)
		}
	}
	if !isWhitespace(state.script[state.index]) && state.script[state.index] != '\n' {
		c := state.script[state.index]
		state.index = start - 1
		return state.Error("expected whitespace after number, got `%v`", c)
	}
	state.addLexeme(LexemeNumber, state.script[start:state.index], start)
	return true
}

func (state *LexState) handleWord() bool {
	c := state.script[state.index]
	if !isWordChar(c) {
		return false
	}
	start := state.index
	state.index++
	for state.index < len(state.script) && isWordChar(state.script[state.index]) {
		state.index++
	}
	state.addLexeme(LexemeWord, state.script[start:state.index], start)
	return true
}

func (state *LexState) step() {
	if state.index >= len(state.script) {
		return
	}
	c := state.script[state.index]
	if c == '\n' { // Handle newline
		state.lastLineIdx = state.index
		state.line++
		state.index++
		return
	}
	if c == '#' { // Handle/skip comments
		for state.index < len(state.script) && state.script[state.index] != '\n' {
			state.index++
		}
		return
	}
	if isWhitespace(c) { // Handle whitespace
		for state.index < len(state.script) && isWhitespace(state.script[state.index]) {
			state.index++
		}
		return
	}
	if state.handleSingleChar() {
		return
	}
	if state.handleString() {
		return
	}
	if state.handleNumber() {
		return
	}
	if state.handleWord() {
		return
	}
	state.Error("unexpected character `%v`", c)
}

func lex(file, script string) (state LexState) {
	state = newLexState(file, script)
	for state.index < len(state.script) && state.err == nil {
		state.step()
	}
	return
}
