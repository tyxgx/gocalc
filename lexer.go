package main

import (
	"fmt"
	"os"
	"strconv"
)

func isDigit(c byte) bool {
	return c <= '9' && c >= '0'
}

const (
	// Number
	TOKEN_NUM = iota

	// Operators
	TOKEN_PLUS
	TOKEN_MINUS
	TOKEN_MULT
	TOKEN_DIV

	// Error
	TOKEN_ERROR

	// EOF
	TOKEN_EOF
)

type Optional struct {
	None int
	Some float64
}

func (opt Optional) hasValue() bool {
	if opt.None == 0 {
		return true
	}
	return false
}

type Token struct {
	tokenStart_   int
	tokenEnd_     int
	tokenType_    int
	tokenContent_ string
	tokenValue_   Optional
}

func (token Token) Dump() {
	fmt.Printf("\ntype: %s; starts: %d; ends: %d; content: %s", func() string {
		switch token.tokenType_ {
		case 0:
			return "Number"
		case 1:
			return "Operator: PLUS"
		case 2:
			return "Operator: MINUS"
		case 3:
			return "Operator: MULTIPLY"
		case 4:
			return "Operator: DIVIDE"
		case 5:
			return "ERROR"
		case 6:
			return "EOF"
		default:
			return "Unknown"
		}
	}(), token.tokenStart_, token.tokenEnd_, token.tokenContent_)
}

type Lexer struct {
	src_            string
	startPos_       int
	curPos_         int
	tokenList_      []Token
	hadError_       bool
	hadErrorAtPrev_ bool
	hadErrorAtCurr_ bool
}

func (lexer *Lexer) globalLogError(errMsg string, eof bool) {
	if !lexer.hadErrorAtCurr_ && lexer.hadErrorAtPrev_ {
		if !eof {
			lexer.curPos_--
		}

		lexer.logError(errMsg)
		lexer.hadErrorAtPrev_ = false
		lexer.hadErrorAtCurr_ = false

		lexer.startPos_ = lexer.curPos_
	}
}

func (lexer *Lexer) Init(src string) {
	lexer.src_ = src
	lexer.startPos_ = 0
	lexer.curPos_ = 0
	lexer.tokenList_ = []Token{}
}

func (lexer *Lexer) advance() byte {
	lexer.curPos_ = lexer.curPos_ + 1
	return lexer.src_[lexer.curPos_-1]
}

func (lexer *Lexer) atEnd() bool {
	return lexer.curPos_ >= len(lexer.src_)
}

func (lexer *Lexer) peek() byte {
	if lexer.atEnd() {
		return 0
	}
	return lexer.src_[lexer.curPos_]
}

func (lexer *Lexer) peekNext() byte {
	if lexer.curPos_+1 >= len(lexer.src_) {
		return 0
	}
	return lexer.src_[lexer.curPos_+1]
}

func (lexer *Lexer) LexOnetoken() {
	if !lexer.hadErrorAtCurr_ && !lexer.hadErrorAtPrev_ {
		lexer.startPos_ = lexer.curPos_
	}

	var c byte = lexer.advance()
	switch c {
	case '+':
		lexer.hadErrorAtCurr_ = false
		lexer.globalLogError("Illegal or unrecognised token", false)
		lexer.addToken(TOKEN_PLUS, "+", Optional{None: 1, Some: 0})
		break
	case '-':
		lexer.hadErrorAtCurr_ = false
		lexer.globalLogError("Illegal or unrecognised token", false)
		lexer.addToken(TOKEN_MINUS, "-", Optional{None: 1, Some: 0})
		break
	case '/':
		lexer.hadErrorAtCurr_ = false
		lexer.globalLogError("Illegal or unrecognised token", false)
		lexer.addToken(TOKEN_DIV, "/", Optional{None: 1, Some: 0})
		break
	case '*':
		lexer.hadErrorAtCurr_ = false
		lexer.globalLogError("Illegal or unrecognised token", false)
		lexer.addToken(TOKEN_MULT, "*", Optional{None: 1, Some: 0})
		break
	case ' ':
	case '\t':
		break
	default:
		if isDigit(c) {
			lexer.hadErrorAtCurr_ = false
			lexer.globalLogError("Illegal or unrecognised token", false)
			lexer.lexNumber()
			break
		}

		if (lexer.hadErrorAtPrev_ && !lexer.hadErrorAtCurr_) || (lexer.hadErrorAtPrev_ && lexer.atEnd()) {
			lexer.logError("Illegal or unrecognised token")
		}
		lexer.hadErrorAtCurr_ = true
		lexer.hadErrorAtPrev_ = true
	}
	return
}

func (lexer *Lexer) lexNumber() {
	for isDigit(lexer.peek()) {
		_ = lexer.advance()
	}

	if lexer.peek() == '.' && isDigit(lexer.peekNext()) {
		_ = lexer.advance()
		for isDigit(lexer.peek()) {
			_ = lexer.advance()
		}
	}

	content := lexer.src_[lexer.startPos_:lexer.curPos_]
	value, _ := strconv.ParseFloat(content, 64)

	lexer.addToken(TOKEN_NUM, content, Optional{None: 0, Some: value})
}

// define some colors
const colorRed = "\033[0;31m"
const colorNone = "\033[0m"

func (lexer *Lexer) logError(errMsg string) {
	lexer.hadError_ = true
	fmt.Fprintf(os.Stderr, "%serror: %s%s\n", colorRed, colorNone, errMsg)
	fmt.Fprintf(os.Stderr, "| %s\n", lexer.src_)
	fmt.Fprintf(os.Stderr, "| ")
	i := 0
	for i < lexer.startPos_ {
		fmt.Fprintf(os.Stderr, " ")
		i += 1
	}
	i += 1
	fmt.Fprintf(os.Stderr, "^")
	for i < lexer.curPos_ {
		fmt.Fprintf(os.Stderr, "~")
		i += 1
	}
	fmt.Fprintf(os.Stderr, "\n")
}

func (lexer *Lexer) addToken(tt int, tokenContent string, tokenValue Optional) {
	var token Token = Token{
		tokenStart_:   lexer.startPos_,
		tokenEnd_:     lexer.curPos_,
		tokenType_:    tt,
		tokenContent_: tokenContent,
		tokenValue_:   tokenValue,
	}
	lexer.tokenList_ = append(lexer.tokenList_, token)
}

func (lexer *Lexer) Lex() {
	for !lexer.atEnd() {
		lexer.LexOnetoken()
	}
	lexer.addToken(TOKEN_EOF, "EOF", Optional{None: 1, Some: 0})
}

func (lexer *Lexer) GetTokenList() []Token {
	return lexer.tokenList_
}
