package parser

import (
	"fmt"
	"monkey/ast"
	"monkey/lexer"
	"monkey/token"
	"strconv"
)

const (
	_ int = iota
	LOWEST
	EQUALS      // ==
	LESSGREATER // < or >
	SUM         // +
	PRODUCT     // *
	PREFIX      // -n or !n
	CALL        // function(x)
)

var precedenceMap = map[token.TokenType]int{
	token.EQ:       EQUALS,
	token.NOT_EQ:   EQUALS,
	token.LT:       LESSGREATER,
	token.GT:       LESSGREATER,
	token.PLUS:     SUM,
	token.MINUS:    SUM,
	token.ASTERISK: PRODUCT,
	token.SLASH:    PRODUCT,
	token.LPAREN:   CALL,
}

func getPrecedence(tokenType token.TokenType) int {
	if precedence, ok := precedenceMap[tokenType]; ok {
		return precedence
	}
	return LOWEST
}

type Parser struct {
	lexer  *lexer.Lexer
	errors []string

	curToken  token.Token
	peekToken token.Token

	prefixParseFns map[token.TokenType]prefixParseFn
	infixParseFns  map[token.TokenType]infixParseFn
}

type (
	prefixParseFn func() ast.Expression
	infixParseFn  func(expression ast.Expression) ast.Expression
)

func New(lexer *lexer.Lexer) *Parser {
	parser := &Parser{lexer: lexer, errors: []string{}}

	// Set curToken and peekToken
	parser.nextToken()
	parser.nextToken()

	parser.prefixParseFns = make(map[token.TokenType]prefixParseFn)
	parser.registerPrefixFn(token.IDENT, parser.parseIdentifier)
	parser.registerPrefixFn(token.INT, parser.parseIntegerLiteral)
	parser.registerPrefixFn(token.BANG, parser.parsePrefixExpression)
	parser.registerPrefixFn(token.MINUS, parser.parsePrefixExpression)
	parser.registerPrefixFn(token.TRUE, parser.parseBooleanLiteral)
	parser.registerPrefixFn(token.FALSE, parser.parseBooleanLiteral)
	parser.registerPrefixFn(token.STRING, parser.parseStringLiteral)
	parser.registerPrefixFn(token.LPAREN, parser.parseGroupedExpression)
	parser.registerPrefixFn(token.IF, parser.parseIfExpression)
	parser.registerPrefixFn(token.FUNCTION, parser.parseFunctionLiteral)

	parser.infixParseFns = make(map[token.TokenType]infixParseFn)
	parser.registerInfixFn(token.EQ, parser.parseInfixExpression)
	parser.registerInfixFn(token.NOT_EQ, parser.parseInfixExpression)
	parser.registerInfixFn(token.LT, parser.parseInfixExpression)
	parser.registerInfixFn(token.GT, parser.parseInfixExpression)
	parser.registerInfixFn(token.PLUS, parser.parseInfixExpression)
	parser.registerInfixFn(token.MINUS, parser.parseInfixExpression)
	parser.registerInfixFn(token.ASTERISK, parser.parseInfixExpression)
	parser.registerInfixFn(token.SLASH, parser.parseInfixExpression)
	parser.registerInfixFn(token.LPAREN, parser.parseCallExpression)

	return parser
}

func (parser *Parser) nextToken() {
	parser.curToken = parser.peekToken
	parser.peekToken = parser.lexer.NextToken()
}

func (parser *Parser) registerPrefixFn(tokenType token.TokenType, fn prefixParseFn) {
	parser.prefixParseFns[tokenType] = fn
}

func (parser *Parser) registerInfixFn(tokenType token.TokenType, fn infixParseFn) {
	parser.infixParseFns[tokenType] = fn
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
	case token.RETURN:
		return parser.parseReturnStatement()
	default:
		return parser.parseExpressionStatement()
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

	parser.nextToken()
	statement.Value = parser.parseExpression(LOWEST)

	if parser.peekTokenIs(token.SEMICOLON) {
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

func (parser *Parser) parseReturnStatement() *ast.ReturnStatement {
	statement := &ast.ReturnStatement{Token: parser.curToken}

	parser.nextToken()
	statement.ReturnValue = parser.parseExpression(LOWEST)

	for !parser.curTokenIs(token.SEMICOLON) {
		parser.nextToken()
	}

	return statement
}

func (parser *Parser) parseExpressionStatement() *ast.ExpressionStatement {
	statement := &ast.ExpressionStatement{Token: parser.curToken}

	statement.Expression = parser.parseExpression(LOWEST)
	if parser.peekTokenIs(token.SEMICOLON) {
		parser.nextToken()
	}

	return statement
}

func (parser *Parser) parseExpression(precedence int) ast.Expression {
	prefixFn := parser.prefixParseFns[parser.curToken.Type]
	if prefixFn == nil {
		parser.noPrefixParserFnError(parser.curToken.Type)
		return nil
	}

	leftExp := prefixFn()

	for !parser.peekTokenIs(token.SEMICOLON) && precedence < parser.peekPrecendence() {
		infixFn := parser.infixParseFns[parser.peekToken.Type]
		if infixFn == nil {
			return leftExp
		}

		parser.nextToken()

		leftExp = infixFn(leftExp)
	}

	return leftExp
}

func (parser *Parser) noPrefixParserFnError(tt token.TokenType) {
	msg := fmt.Sprintf("No prefix parser function for %s found", tt)
	parser.errors = append(parser.errors, msg)
}

func (parser *Parser) parseIdentifier() ast.Expression {
	return &ast.Identifier{Token: parser.curToken, Value: parser.curToken.Literal}
}

func (parser *Parser) parseIntegerLiteral() ast.Expression {
	literal := &ast.IntegerLiteral{Token: parser.curToken}

	value, err := strconv.ParseInt(parser.curToken.Literal, 0, 64)
	if err != nil {
		msg := fmt.Sprintf("could not parse %q as integer", parser.curToken.Literal)
		parser.errors = append(parser.errors, msg)
		return nil
	}

	literal.Value = value
	return literal
}

func (parser *Parser) parseBooleanLiteral() ast.Expression {
	return &ast.BooleanLiteral{Token: parser.curToken, Value: parser.curTokenIs(token.TRUE)}
}

func (parser *Parser) parseStringLiteral() ast.Expression {
	return &ast.StringLiteral{Token: parser.curToken, Value: parser.curToken.Literal}
}

func (parser *Parser) parsePrefixExpression() ast.Expression {
	expression := &ast.PrefixExpression{
		Token:    parser.curToken,
		Operator: parser.curToken.Literal,
	}

	parser.nextToken()

	expression.Right = parser.parseExpression(PREFIX)
	return expression
}

func (parser *Parser) parseInfixExpression(left ast.Expression) ast.Expression {
	expression := &ast.InfixExpression{
		Token:    parser.curToken,
		Left:     left,
		Operator: parser.curToken.Literal,
	}

	precedence := parser.curPrecedence()
	parser.nextToken()
	expression.Right = parser.parseExpression(precedence)

	return expression

}

func (parser *Parser) peekPrecendence() int {
	return getPrecedence(parser.peekToken.Type)
}

func (parser *Parser) curPrecedence() int {
	return getPrecedence(parser.curToken.Type)
}

func (parser *Parser) parseGroupedExpression() ast.Expression {
	parser.nextToken()

	expression := parser.parseExpression(LOWEST)
	if !parser.expectPeek(token.RPAREN) {
		return nil
	}

	return expression
}

func (parser *Parser) parseIfExpression() ast.Expression {
	expression := &ast.IfExpression{Token: parser.curToken}

	if !parser.expectPeek(token.LPAREN) {
		return nil
	}

	parser.nextToken()
	expression.Condition = parser.parseExpression(LOWEST)

	if !parser.expectPeek(token.RPAREN) {
		return nil
	}

	parser.nextToken()
	expression.Consequence = parser.parseBlockStatement()

	if parser.peekTokenIs(token.ELSE) {
		parser.nextToken()
		if !parser.expectPeek(token.LBRACE) {
			return nil
		}
		expression.Alternative = parser.parseBlockStatement()
	}

	return expression
}

func (parser *Parser) parseBlockStatement() *ast.BlockStatement {
	block := &ast.BlockStatement{Token: parser.curToken}
	block.Statements = []ast.Statement{}

	parser.nextToken()
	for !parser.curTokenIs(token.RBRACE) && !parser.curTokenIs(token.EOF) {
		stmt := parser.parseStatement()
		if stmt != nil {
			block.Statements = append(block.Statements, stmt)
		}
		parser.nextToken()
	}

	return block
}

func (parser *Parser) parseFunctionLiteral() ast.Expression {
	literal := &ast.FunctionLiteral{Token: parser.curToken}

	if !parser.expectPeek(token.LPAREN) {
		return nil
	}

	literal.Parameters = parser.parseFunctionParameters()

	if !parser.expectPeek(token.LBRACE) {
		return nil
	}

	literal.Body = parser.parseBlockStatement()

	return literal
}

func (parser *Parser) parseFunctionParameters() []*ast.Identifier {
	var parameters []*ast.Identifier

	if parser.peekTokenIs(token.RPAREN) {
		parser.nextToken()
		return parameters
	}

	parser.nextToken()
	identifier := &ast.Identifier{Token: parser.curToken, Value: parser.curToken.Literal}
	parameters = append(parameters, identifier)

	for parser.peekTokenIs(token.COMMA) {
		parser.nextToken()
		parser.nextToken()
		identifier := &ast.Identifier{Token: parser.curToken, Value: parser.curToken.Literal}
		parameters = append(parameters, identifier)
	}

	if !parser.expectPeek(token.RPAREN) {
		return nil
	}

	return parameters
}

func (parser *Parser) parseCallExpression(function ast.Expression) ast.Expression {
	expression := &ast.CallExpression{Token: parser.curToken, Function: function}
	expression.Arguments = parser.parseCallArguments()
	return expression
}

func (parser *Parser) parseCallArguments() []ast.Expression {
	var arguments []ast.Expression

	if parser.peekTokenIs(token.RPAREN) {
		parser.nextToken()
		return arguments
	}

	parser.nextToken()
	expression := parser.parseExpression(LOWEST)
	arguments = append(arguments, expression)

	for parser.peekTokenIs(token.COMMA) {
		parser.nextToken()
		parser.nextToken()
		arguments = append(arguments, parser.parseExpression(LOWEST))
	}

	if !parser.expectPeek(token.RPAREN) {
		return nil
	}

	return arguments
}
