package main

import (
	"fmt"
	"os"
)

// Define the AST
type Node struct {
	token_ Token
	rhs_   *Node
	lhs_   *Node
}

// expr    := term
// term    := factor [ ("+" | "-") factor ]*
// factor  := unary [ ("/" | "*") unary ]*
// unary   := [ ("-") unary ] | primary
// primary := NUM [ | '(' expr ')']

// 1 + 1 / 2 2
//     +
//	  / \
//   1   /
//      1 2

type Parser struct {
	src_       string
	curPos_    int
	tokenList_ []Token
	evalStack_ []float64
}

func isOperator(oper int) bool {
	switch oper {
	case TOKEN_DIV, TOKEN_MINUS, TOKEN_PLUS, TOKEN_MULT:
		return true
	default:
		return false
	}
}

func (parser *Parser) InitParser(src string, tokenList []Token) {
	parser.src_ = src
	parser.curPos_ = 0
	parser.tokenList_ = tokenList
	parser.evalStack_ = []float64{}
}

func (parser *Parser) atEof() bool {
	return parser.peek().tokenType_ == TOKEN_EOF
}

func (parser *Parser) previous(offset int) Token {
	return parser.tokenList_[parser.curPos_-offset]
}

func (parser *Parser) peek() Token {
	return parser.tokenList_[parser.curPos_]
}

func (parser *Parser) advance() Token {
	if !parser.atEof() {
		parser.curPos_++
	}
	return parser.previous(1)
}

func (parser *Parser) tokenMatches(expected int) bool {
	if parser.atEof() {
		return false
	}
	res := parser.peek().tokenType_ == expected
	return res
}

func (parser *Parser) match(expected ...int) bool {
	for _, tt := range expected {
		if parser.tokenMatches(tt) {
			parser.advance()
			return true
		}
	}
	return false
}

func (parser *Parser) consume(expected int, errorMsg string) bool {
	if parser.tokenMatches(expected) {
		_ = parser.advance()
		return true
	}

	parser.logError(errorMsg)
	return false
}

func (parser *Parser) expression() *Node {
	return parser.parseTerm()
}

func (parser *Parser) parseTerm() *Node {
	expr := parser.parseFactor()
	if expr == nil {
		return nil
	}

	for parser.match(TOKEN_PLUS, TOKEN_MINUS) {
		token := parser.previous(1)
		rhs := parser.parseFactor()
		if rhs == nil {
			parser.logError("expected expression")
			return nil
		}
		expr = &Node{
			token_: token,
			rhs_:   rhs,
			lhs_:   expr,
		}
	}

	if parser.curPos_ != 0 && parser.tokenMatches(TOKEN_NUM) && !isOperator(parser.previous(1).tokenType_) {
		parser.logError("syntax error")
		return nil
	}

	return expr
}

func (parser *Parser) parseFactor() *Node {
	expr := parser.parseUnary()
	if expr == nil {
		return nil
	}

	for parser.match(TOKEN_MULT, TOKEN_DIV) {
		token := parser.previous(1)
		rhs := parser.parseUnary()
		if rhs == nil {
			return nil
		}
		expr = &Node{
			token_: token,
			rhs_:   rhs,
			lhs_:   expr,
		}
	}

	if parser.curPos_ != 0 && parser.tokenMatches(TOKEN_NUM) && !isOperator(parser.previous(1).tokenType_) {
		parser.logError("syntax error")
		return nil
	}

	return expr
}

func (parser *Parser) parseUnary() *Node {
	if parser.match(TOKEN_MINUS) {
		token := parser.previous(1)
		rhs := parser.parseUnary()
		if rhs == nil {
			parser.logError("expected expression")
			return nil
		}
		dummyToken := Token{
			tokenStart_:   parser.peek().tokenStart_,
			tokenEnd_:     parser.peek().tokenStart_,
			tokenType_:    TOKEN_NUM,
			tokenContent_: "0",
			tokenValue_:   Optional{None: 0, Some: 0.00},
		}
		dummyNode := &Node{
			token_: dummyToken,
			lhs_:   nil,
			rhs_:   nil,
		}
		node := &Node{
			token_: token,
			rhs_:   rhs,
			lhs_:   dummyNode,
		}
		return node
	}

	if parser.curPos_ != 0 && parser.tokenMatches(TOKEN_NUM) && !isOperator(parser.previous(1).tokenType_) {
		parser.logError("syntax error")
		return nil
	}

	return parser.parsePrimary()
}

func (parser *Parser) parsePrimary() *Node {
	if parser.match(TOKEN_NUM) {
		token := parser.previous(1)
		node := &Node{
			token_: token,
			lhs_:   nil,
			rhs_:   nil,
		}
		return node
	}
	parser.logError("expected expression. Did you put this token accidentally?")
	return nil
}

func (tree *Node) eval() float64 {
	switch tree.token_.tokenType_ {
	case TOKEN_PLUS:
		return tree.lhs_.eval() + tree.rhs_.eval()
	case TOKEN_MINUS:
		return tree.lhs_.eval() - tree.rhs_.eval()
	case TOKEN_DIV:
		return tree.lhs_.eval() / tree.rhs_.eval()
	case TOKEN_MULT:
		return tree.lhs_.eval() * tree.rhs_.eval()
	case TOKEN_NUM:
		return tree.token_.tokenValue_.Some
	}
	return -69.6969
}

func (parser *Parser) logError(errMsg string) {
	fmt.Fprintf(os.Stderr, "%serror: %s%s\n", colorRed, colorNone, errMsg)
	fmt.Fprintf(os.Stderr, "| %s\n", parser.src_)
	fmt.Fprintf(os.Stderr, "| ")
	i := 0
	for i < parser.peek().tokenStart_ {
		fmt.Fprintf(os.Stderr, " ")
		i += 1
	}
	i += 1
	fmt.Fprintf(os.Stderr, "^")
	for i < parser.peek().tokenEnd_ {
		fmt.Fprintf(os.Stderr, "~")
		i += 1
	}
	fmt.Fprintf(os.Stderr, "\n")
}
