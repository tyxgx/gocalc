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
	for true {
		buf := bufio.NewScanner(os.Stdin)

		fmt.Fprintf(os.Stderr, "\n>>> ")

		_ = buf.Scan()

		if buf.Text() == "exit" {
			return
		}

		lexer := Lexer{}
		lexer.Init(buf.Text())
		lexer.Lex()

		if lexer.hadError_ {
			continue
		}

		tokenList := lexer.tokenList_
		for _, token := range tokenList {
			token.Dump()
		}
	}

}
