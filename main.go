package main

import (
	"bufio"
	"fmt"
	"os"
)

func debug() {
	fmt.Println("Inside driver")
}

func main() {
	fmt.Fprintf(os.Stderr, "GoCalc 0.0.1. Type \"exit\" or press Ctrl-D to exit the prompt.")
	for true {
		buf := bufio.NewScanner(os.Stdin)
		fmt.Fprintf(os.Stderr, "\n>>> ")
		success := buf.Scan()

		if buf.Text() == "exit" {
			return
		}

		// Handle EOF input (in case of ^D)
		if err := buf.Err(); err == nil && !success {
			return
		}

		lexer := Lexer{}
		lexer.Init(buf.Text())
		lexer.Lex()

		if lexer.hadError_ {
			continue
		}

		parser := Parser{}
		parser.InitParser(buf.Text(), lexer.tokenList_)
		node := parser.expression()
		if node == nil {
			continue
		}
		fmt.Fprintf(os.Stderr, "%f", node.eval())
	}

}
