package compiler

import (
	"github.com/siyul-park/minijs/ast"
	"github.com/siyul-park/minijs/bytecode"
	"github.com/siyul-park/minijs/token"
	"github.com/stretchr/testify/assert"
	"math"
	"testing"
)

func TestCompiler_Compile(t *testing.T) {
	tests := []struct {
		node         ast.Node
		instructions []bytecode.Instruction
	}{
		{
			node: &ast.NumberLiteral{
				Token: token.NewToken(token.DECIMAL, `1234567890`),
				Value: 1234567890,
			},
			instructions: []bytecode.Instruction{
				bytecode.New(bytecode.F64LOAD, math.Float64bits(1234567890)),
			},
		},
		{
			node: &ast.Program{Nodes: []ast.Node{
				&ast.PrefixExpression{
					Token: token.NewToken(token.MINUS, "-"),
					Right: &ast.NumberLiteral{Token: token.NewToken(token.DECIMAL, `1234567890`), Value: 1234567890},
				},
			}},
			instructions: []bytecode.Instruction{
				bytecode.New(bytecode.F64LOAD, math.Float64bits(1234567890)),
				bytecode.New(bytecode.F64LOAD, math.Float64bits(-1)),
				bytecode.New(bytecode.F64MUL),
			},
		},
	}

	for _, tt := range tests {
		var code bytecode.Bytecode
		code.Append(tt.instructions...)

		t.Run(code.String(), func(t *testing.T) {
			compiler := New(tt.node)

			result, err := compiler.Compile()
			assert.NoError(t, err)
			assert.Equal(t, code.String(), result.String())
		})
	}
}
