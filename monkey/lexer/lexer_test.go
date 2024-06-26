package lexer

import (
	"monkey/token"
	"testing"
)

type NextTokenTest struct {
	expectedType    token.TokenType
	expectedLiteral string
}

func verifyNextToken(t *testing.T, input string, tests []NextTokenTest) {
	l := New(input)

	for i, tt := range tests {
		tok := l.NextToken()

		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q",
				i, tt.expectedType, tok.Type)
		}

		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q",
				i, tt.expectedLiteral, tok.Literal)
		}
	}
}

// Chapter 1.2
func TestNextTokenFirstTokens(t *testing.T) {
	input := `=+(){},;`

	tests := []NextTokenTest{
		{token.ASSIGN, "="},
		{token.PLUS, "+"},
		{token.LPAREN, "("},
		{token.RPAREN, ")"},
		{token.LBRACE, "{"},
		{token.RBRACE, "}"},
		{token.COMMA, ","},
		{token.SEMICOLON, ";"},
		{token.EOF, ""},
	}

	verifyNextToken(t, input, tests)
}

// Chapter 1.3
func TestNextTokenAssignmentAndFns(t *testing.T) {
	input := `let five = 5;
	let ten = 10;
	
	let add = fn(x, y) {
	  x + y;
	};

	let result = add(five, ten);
	`

	tests := []NextTokenTest{
		{token.LET, "let"},
		{token.IDENT, "five"},
		{token.ASSIGN, "="},
		{token.INT, "5"},
		{token.SEMICOLON, ";"},
		{token.LET, "let"},
		{token.IDENT, "ten"},
		{token.ASSIGN, "="},
		{token.INT, "10"},
		{token.SEMICOLON, ";"},
		{token.LET, "let"},
		{token.IDENT, "add"},
		{token.ASSIGN, "="},
		{token.FUNCTION, "fn"},
		{token.LPAREN, "("},
		{token.IDENT, "x"},
		{token.COMMA, ","},
		{token.IDENT, "y"},
		{token.RPAREN, ")"},
		{token.LBRACE, "{"},
		{token.IDENT, "x"},
		{token.PLUS, "+"},
		{token.IDENT, "y"},
		{token.SEMICOLON, ";"},
		{token.RBRACE, "}"},
		{token.SEMICOLON, ";"},
		{token.LET, "let"},
		{token.IDENT, "result"},
		{token.ASSIGN, "="},
		{token.IDENT, "add"},
		{token.LPAREN, "("},
		{token.IDENT, "five"},
		{token.COMMA, ","},
		{token.IDENT, "ten"},
		{token.RPAREN, ")"},
		{token.SEMICOLON, ";"},
		{token.EOF, ""},
	}

	verifyNextToken(t, input, tests)
}

// Chapter 1.4
func TestNextTokenOperatorsMoreKeywords(t *testing.T) {
	input := `!-/*5;
	5 < 10 > 5;
	if (5 < 10) {
		return true;
	} else {
		return false;
	}
	
	10 == 10;
	10 != 9;
	`

	tests := []NextTokenTest{
		{token.BANG, "!"},
		{token.MINUS, "-"},
		{token.SLASH, "/"},
		{token.ASTERISK, "*"},
		{token.INT, "5"},
		{token.SEMICOLON, ";"},
		{token.INT, "5"},
		{token.LT, "<"},
		{token.INT, "10"},
		{token.GT, ">"},
		{token.INT, "5"},
		{token.SEMICOLON, ";"},
		{token.IF, "if"},
		{token.LPAREN, "("},
		{token.INT, "5"},
		{token.LT, "<"},
		{token.INT, "10"},
		{token.RPAREN, ")"},
		{token.LBRACE, "{"},
		{token.RETURN, "return"},
		{token.TRUE, "true"},
		{token.SEMICOLON, ";"},
		{token.RBRACE, "}"},
		{token.ELSE, "else"},
		{token.LBRACE, "{"},
		{token.RETURN, "return"},
		{token.FALSE, "false"},
		{token.SEMICOLON, ";"},
		{token.RBRACE, "}"},
		{token.INT, "10"},
		{token.EQ, "=="},
		{token.INT, "10"},
		{token.SEMICOLON, ";"},
		{token.INT, "10"},
		{token.NOT_EQ, "!="},
		{token.INT, "9"},
		{token.SEMICOLON, ";"},
		{token.EOF, ""},
	}

	verifyNextToken(t, input, tests)
}

// Chap 4.1
func TestNextTokenString(t *testing.T) {
	input := `
"foobar"
"foo bar"
`
	tests := []NextTokenTest{
		{token.STRING, "foobar"},
		{token.STRING, "foo bar"},
	}

	verifyNextToken(t, input, tests)
}

// Chap 4.4
func TestNextTokenBrackets(t *testing.T) {
	input := "[1, 2]"
	tests := []NextTokenTest{
		{token.LBRACKET, "["},
		{token.INT, "1"},
		{token.COMMA, ","},
		{token.INT, "2"},
		{token.RBRACKET, "]"},
	}
	verifyNextToken(t, input, tests)
}

// Chap 4.5
func TestNextTokenColonForHash(t *testing.T) {
	input := "{1: 2, 3: 4}"
	tests := []NextTokenTest{
		{token.LBRACE, "{"},
		{token.INT, "1"},
		{token.COLON, ":"},
		{token.INT, "2"},
		{token.COMMA, ","},
		{token.INT, "3"},
		{token.COLON, ":"},
		{token.INT, "4"},
		{token.RBRACE, "}"},
	}
	verifyNextToken(t, input, tests)

}
