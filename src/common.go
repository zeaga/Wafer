package main

import "fmt"

func isDigit(char byte) bool {
	return char >= '0' && char <= '9'
}

func isWordChar(char byte) bool {
	return (char >= '!' && char < '\\') ||
		(char > '\\' && char <= '~')
}

func isWhitespace(char byte) bool {
	return char == ' ' || char == '\t' || char == '\r'
}

func boolToFloat(from bool) float64 {
	if from {
		return 1
	}
	return 0
}

func floatToBool(to float64) bool {
	return to != 0
}

func (state *EvalState) printv(value Value) {
	str := value.text
	if value.kind == ValueNumber {
		str = fmt.Sprint(value.number)
	}
	state.lastPrintedWasNewline = str[len(str)-1] == '\n'
	fmt.Print(str)
}

func (state *LexState) Error(format string, args ...any) bool {
	msg := fmt.Sprintf(format, args...)
	state.err = fmt.Errorf("%s:%d:%d: %s", state.file, state.line+1, state.index-state.lastLineIdx+1, msg)
	return true
}

func (state *ParseState) Error(format string, args ...any) bool {
	msg := fmt.Sprintf(format, args...)
	state.err = fmt.Errorf("%s:%d:%d: %s", state.file, state.line+1, state.col+1, msg)
	return true
}

func (state *EvalState) Error(format string, args ...any) bool {
	token := state.currentToken()
	msg := fmt.Sprintf(format, args...)
	state.err = fmt.Errorf("%s:%d:%d: %s", token.file, token.line+1, token.col+1, msg)
	return true
}

func (lexeme *Lexeme) Error(format string, args ...any) bool {
	msg := fmt.Sprintf(format, args...)
	lexeme.state.err = fmt.Errorf("%s:%d:%d: %s", lexeme.file, lexeme.line+1, lexeme.col+1, msg)
	return true
}
