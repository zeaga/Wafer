package main

import (
	"fmt"
	"os"
	"path/filepath"
)

const APP_NAME = "wafer"

func printUsage() {
	exeName := APP_NAME
	if exePath, err := os.Executable(); err == nil {
		exeName = filepath.Base(exePath)
	}
	fmt.Printf("Usage: %v <filename>\n", exeName)
}

func loadStdLib() (words map[string]Word, err error) {
	lexState := lex("stdlib", STDLIB)
	if lexState.err != nil {
		err = lexState.err
		return
	}

	parseState := parse(lexState)
	if parseState.err != nil {
		err = parseState.err
		return
	}
	evalState := eval(parseState, make(map[string]Word))
	if !evalState.lastPrintedWasNewline {
		fmt.Print("\n")
	}
	if evalState.err != nil {
		err = evalState.err
	}
	words = evalState.words
	return
}

func main() {
	if len(os.Args) < 2 {
		printUsage()
		return
	}

	filename := os.Args[1]
	data, err := os.ReadFile(filename)
	if err != nil {
		fmt.Printf("failed to read file: %v\n", err)
		return
	}
	contents := string(data)

	words, err := loadStdLib()
	if err != nil {
		fmt.Println(err)
	}

	lexState := lex(filename, contents)
	if lexState.err != nil {
		fmt.Println(lexState.err)
		return
	}

	parseState := parse(lexState)
	if parseState.err != nil {
		fmt.Println(parseState.err)
		return
	}

	evalState := eval(parseState, words)
	if !evalState.lastPrintedWasNewline {
		fmt.Print("\n")
	}
	if evalState.err != nil {
		token := evalState.currentToken()
		if token == nil {
			fmt.Printf("%v:?:?:E: %v\n", filename, evalState.err)
			return
		}
		fmt.Println(evalState.err)
		return
	}
}
