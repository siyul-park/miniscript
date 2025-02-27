package parser

import (
	"strings"
	"testing"

	"github.com/siyul-park/minijs/internal/ast"
	"github.com/siyul-park/minijs/internal/lexer"
	"github.com/siyul-park/minijs/internal/token"

	"github.com/stretchr/testify/assert"
)

func TestParser_Parse(t *testing.T) {
	tests := []struct {
		source  string
		program *ast.Program
	}{
		{"", ast.NewProgram()},
		{";", ast.NewProgram(ast.NewEmptyStatement())},
		{
			"{ 1; 2; }",
			ast.NewProgram(
				ast.NewBlockStatement(
					ast.NewExpressionStatement(
						ast.NewNumberLiteral(token.New(token.NUMBER, "1"), 1),
					),
					ast.NewExpressionStatement(
						ast.NewNumberLiteral(token.New(token.NUMBER, "2"), 2),
					),
				),
			),
		},
		{
			"a + b; c + d",
			ast.NewProgram(
				ast.NewExpressionStatement(
					ast.NewInfixExpression(
						token.New(token.PLUS, "+"),
						ast.NewIdentifierLiteral(token.New(token.IDENTIFIER, "a"), "a"),
						ast.NewIdentifierLiteral(token.New(token.IDENTIFIER, "b"), "b"),
					),
				),
				ast.NewExpressionStatement(
					ast.NewInfixExpression(
						token.New(token.PLUS, "+"),
						ast.NewIdentifierLiteral(token.New(token.IDENTIFIER, "c"), "c"),
						ast.NewIdentifierLiteral(token.New(token.IDENTIFIER, "d"), "d"),
					),
				),
			),
		},
		{
			"null",
			ast.NewProgram(
				ast.NewExpressionStatement(
					ast.NewNullLiteral(token.New(token.NULL, "null")),
				),
			),
		},
		{
			"undefined",
			ast.NewProgram(
				ast.NewExpressionStatement(
					ast.NewUndefinedLiteral(token.New(token.UNDEFINED, "undefined")),
				),
			),
		},
		{
			"123",
			ast.NewProgram(
				ast.NewExpressionStatement(
					ast.NewNumberLiteral(token.New(token.NUMBER, "123"), 123),
				),
			),
		},
		{
			"1.23",
			ast.NewProgram(
				ast.NewExpressionStatement(
					ast.NewNumberLiteral(token.New(token.NUMBER, "1.23"), 1.23),
				),
			),
		},
		{
			"0b01",
			ast.NewProgram(
				ast.NewExpressionStatement(
					ast.NewNumberLiteral(token.New(token.NUMBER, "0b01"), 0b01),
				),
			),
		},
		{
			"0o01",
			ast.NewProgram(
				ast.NewExpressionStatement(
					ast.NewNumberLiteral(token.New(token.NUMBER, "0o01"), 0o01),
				),
			),
		},
		{
			"0x01",
			ast.NewProgram(
				ast.NewExpressionStatement(
					ast.NewNumberLiteral(token.New(token.NUMBER, "0x01"), 0x01),
				),
			),
		},
		{
			"true",
			ast.NewProgram(
				ast.NewExpressionStatement(
					ast.NewBoolLiteral(token.New(token.TRUE, "true"), true),
				),
			),
		},
		{
			"foo",
			ast.NewProgram(
				ast.NewExpressionStatement(
					ast.NewIdentifierLiteral(token.New(token.IDENTIFIER, "foo"), "foo"),
				),
			),
		},
		{
			`"hello"`,
			ast.NewProgram(
				ast.NewExpressionStatement(
					ast.NewStringLiteral(token.New(token.STRING, "hello"), "hello"),
				),
			),
		},
		{
			"-1",
			ast.NewProgram(
				ast.NewExpressionStatement(
					ast.NewPrefixExpression(
						token.New(token.MINUS, "-"),
						ast.NewNumberLiteral(token.New(token.NUMBER, "1"), 1),
					),
				),
			),
		},
		{
			"a + b",
			ast.NewProgram(
				ast.NewExpressionStatement(
					ast.NewInfixExpression(
						token.New(token.PLUS, "+"),
						ast.NewIdentifierLiteral(token.New(token.IDENTIFIER, "a"), "a"),
						ast.NewIdentifierLiteral(token.New(token.IDENTIFIER, "b"), "b"),
					),
				),
			),
		},
		{
			"a + b + c",
			ast.NewProgram(
				ast.NewExpressionStatement(
					ast.NewInfixExpression(
						token.New(token.PLUS, "+"),
						ast.NewInfixExpression(
							token.New(token.PLUS, "+"),
							ast.NewIdentifierLiteral(token.New(token.IDENTIFIER, "a"), "a"),
							ast.NewIdentifierLiteral(token.New(token.IDENTIFIER, "b"), "b"),
						),
						ast.NewIdentifierLiteral(token.New(token.IDENTIFIER, "c"), "c"),
					),
				),
			),
		},
		{
			"a * b + c",
			ast.NewProgram(
				ast.NewExpressionStatement(
					ast.NewInfixExpression(
						token.New(token.PLUS, "+"),
						ast.NewInfixExpression(
							token.New(token.MULTIPLY, "*"),
							ast.NewIdentifierLiteral(token.New(token.IDENTIFIER, "a"), "a"),
							ast.NewIdentifierLiteral(token.New(token.IDENTIFIER, "b"), "b"),
						),
						ast.NewIdentifierLiteral(token.New(token.IDENTIFIER, "c"), "c"),
					),
				),
			),
		},
		{
			"a = b",
			ast.NewProgram(
				ast.NewExpressionStatement(
					ast.NewAssignmentExpression(
						token.New(token.ASSIGN, "="),
						ast.NewIdentifierLiteral(token.New(token.IDENTIFIER, "a"), "a"),
						ast.NewIdentifierLiteral(token.New(token.IDENTIFIER, "b"), "b"),
					),
				),
			),
		},
		{
			"var a = b",
			ast.NewProgram(
				ast.NewVariableStatement(
					token.New(token.VAR, "var"),
					ast.NewAssignmentExpression(
						token.New(token.ASSIGN, "="),
						ast.NewIdentifierLiteral(token.New(token.IDENTIFIER, "a"), "a"),
						ast.NewIdentifierLiteral(token.New(token.IDENTIFIER, "b"), "b"),
					),
				),
			),
		},
		{
			"var a = b, c = d",
			ast.NewProgram(
				ast.NewVariableStatement(
					token.New(token.VAR, "var"),
					ast.NewAssignmentExpression(
						token.New(token.ASSIGN, "="),
						ast.NewIdentifierLiteral(token.New(token.IDENTIFIER, "a"), "a"),
						ast.NewIdentifierLiteral(token.New(token.IDENTIFIER, "b"), "b"),
					),
					ast.NewAssignmentExpression(
						token.New(token.ASSIGN, "="),
						ast.NewIdentifierLiteral(token.New(token.IDENTIFIER, "c"), "c"),
						ast.NewIdentifierLiteral(token.New(token.IDENTIFIER, "d"), "d"),
					),
				),
			),
		},
	}

	for _, tt := range tests {
		t.Run(tt.source, func(t *testing.T) {
			l := lexer.New(strings.NewReader(tt.source))
			p := New(l)
			program, err := p.Parse()
			assert.NoError(t, err)
			assert.Equal(t, tt.program, program)
		})
	}
}
