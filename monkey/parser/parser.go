package parser

import (
	"fmt"
	"monkey/ast"
	"monkey/lexer"
	"monkey/token"
)

type Parser struct {
	lexer     *lexer.Lexer
	errors    []string
	curToken  token.Token
	peekToken token.Token
}

func New(lexer *lexer.Lexer) *Parser {
	parser := &Parser{lexer: lexer, errors: []string{}}

	// Set curToken and peekToken
	parser.nextToken()
	parser.nextToken()
	return parser
}

func (parser *Parser) nextToken() {
	parser.curToken = parser.peekToken
	parser.peekToken = parser.lexer.NextToken()
}

func (parser *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{}
	program.Statements = []ast.Statement{}

	for !parser.curTokenIs(token.EOF) {
		statement := parser.parseStatement()
		if statement != nil {
			program.Statements = append(program.Statements, statement)
		}
		parser.nextToken()
	}

	return program
}

func (parser *Parser) parseStatement() ast.Statement {
	switch parser.curToken.Type {
	case token.LET:
		return parser.parseLetStatement()
	default:
		return nil
	}
}

func (parser *Parser) parseLetStatement() *ast.LetStatement {
	statement := &ast.LetStatement{Token: parser.curToken}

	if !parser.expectPeek(token.IDENT) {
		return nil
	}

	statement.Name = &ast.Identifier{Token: parser.curToken, Value: parser.curToken.Literal}

	if !parser.expectPeek(token.ASSIGN) {
		return nil
	}

	// TODO: Parse expressions
	for !parser.curTokenIs(token.SEMICOLON) {
		parser.nextToken()
	}

	return statement
}

func (parser *Parser) curTokenIs(tt token.TokenType) bool {
	return parser.curToken.Type == tt
}

func (parser *Parser) peekTokenIs(tt token.TokenType) bool {
	return parser.peekToken.Type == tt
}

func (parser *Parser) expectPeek(tt token.TokenType) bool {
	if parser.peekTokenIs(tt) {
		parser.nextToken()
		return true
	}
	parser.peekError(tt)
	return false
}

func (parser *Parser) Errors() []string {
	return parser.errors
}

func (parser *Parser) peekError(tt token.TokenType) {
	msg := fmt.Sprintf("expected next token to be %s, got %s instead", tt, parser.peekToken.Type)
	parser.errors = append(parser.errors, msg)
}
