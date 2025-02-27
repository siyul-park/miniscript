package interpreter

import (
	"math"
	"testing"

	"github.com/siyul-park/minijs/internal/bytecode"
	"github.com/stretchr/testify/assert"
)

func TestInterpreter_Execute(t *testing.T) {
	tests := []struct {
		instructions []bytecode.Instruction
		literals     []string
		stack        []Value
	}{
		{
			instructions: []bytecode.Instruction{
				bytecode.New(bytecode.NOP),
			},
		},
		{
			instructions: []bytecode.Instruction{
				bytecode.New(bytecode.POP),
			},
		},
		{
			instructions: []bytecode.Instruction{
				bytecode.New(bytecode.I32LOAD, 1),
				bytecode.New(bytecode.SLTSTORE, 1),
			},
		},
		{
			instructions: []bytecode.Instruction{
				bytecode.New(bytecode.I32LOAD, 1),
				bytecode.New(bytecode.SLTSTORE, 1),
				bytecode.New(bytecode.SLTLOAD, 1),
			},
			stack: []Value{Int32(1)},
		},
		{
			instructions: []bytecode.Instruction{
				bytecode.New(bytecode.UNDEFLOAD),
			},
			stack: []Value{Undefined{}},
		},
		{
			instructions: []bytecode.Instruction{
				bytecode.New(bytecode.UNDEFLOAD),
				bytecode.New(bytecode.UNDEFTOF64),
			},
			//stack: []Value{Float64(math.NaN())},
		},
		{
			instructions: []bytecode.Instruction{
				bytecode.New(bytecode.UNDEFLOAD),
				bytecode.New(bytecode.UNDEFTOSTR),
			},
			stack: []Value{String("undefined")},
		},
		{
			instructions: []bytecode.Instruction{
				bytecode.New(bytecode.NULLLOAD),
			},
			stack: []Value{Null{}},
		},
		{
			instructions: []bytecode.Instruction{
				bytecode.New(bytecode.NULLLOAD),
				bytecode.New(bytecode.NULLTOI32),
			},
			stack: []Value{Int32(0)},
		},
		{
			instructions: []bytecode.Instruction{
				bytecode.New(bytecode.NULLLOAD),
				bytecode.New(bytecode.NULLTOSTR),
			},
			stack: []Value{String("null")},
		},
		{
			instructions: []bytecode.Instruction{
				bytecode.New(bytecode.BOOLLOAD, 1),
			},
			stack: []Value{Bool(1)},
		},
		{
			instructions: []bytecode.Instruction{
				bytecode.New(bytecode.BOOLLOAD, 1),
				bytecode.New(bytecode.BOOLTOI32),
			},
			stack: []Value{Int32(1)},
		},
		{
			instructions: []bytecode.Instruction{
				bytecode.New(bytecode.BOOLLOAD, 1),
				bytecode.New(bytecode.BOOLTOSTR),
			},
			stack: []Value{String("true")},
		},
		{
			instructions: []bytecode.Instruction{
				bytecode.New(bytecode.I32LOAD, 1),
			},
			stack: []Value{Int32(1)},
		},
		{
			instructions: []bytecode.Instruction{
				bytecode.New(bytecode.I32LOAD, 1),
				bytecode.New(bytecode.I32LOAD, 2),
				bytecode.New(bytecode.I32ADD),
			},
			stack: []Value{Int32(3)},
		},
		{
			instructions: []bytecode.Instruction{
				bytecode.New(bytecode.I32LOAD, 1),
				bytecode.New(bytecode.I32LOAD, 2),
				bytecode.New(bytecode.I32SUB),
			},
			stack: []Value{Int32(-1)},
		},
		{
			instructions: []bytecode.Instruction{
				bytecode.New(bytecode.I32LOAD, 1),
				bytecode.New(bytecode.I32LOAD, 2),
				bytecode.New(bytecode.I32MUL),
			},
			stack: []Value{Int32(2)},
		},
		{
			instructions: []bytecode.Instruction{
				bytecode.New(bytecode.I32LOAD, 6),
				bytecode.New(bytecode.I32LOAD, 2),
				bytecode.New(bytecode.I32DIV),
			},
			stack: []Value{Int32(3)},
		},
		{
			instructions: []bytecode.Instruction{
				bytecode.New(bytecode.I32LOAD, 7),
				bytecode.New(bytecode.I32LOAD, 3),
				bytecode.New(bytecode.I32MOD),
			},
			stack: []Value{Int32(1)},
		},
		{
			instructions: []bytecode.Instruction{
				bytecode.New(bytecode.I32LOAD, 5),
				bytecode.New(bytecode.I32LOAD, 1),
				bytecode.New(bytecode.I32DIV),
			},
			stack: []Value{Int32(5)},
		},
		{
			instructions: []bytecode.Instruction{
				bytecode.New(bytecode.I32LOAD, 1),
				bytecode.New(bytecode.I32TOBOOL),
			},
			stack: []Value{Bool(1)},
		},
		{
			instructions: []bytecode.Instruction{
				bytecode.New(bytecode.I32LOAD, 5),
				bytecode.New(bytecode.I32TOF64),
			},
			stack: []Value{Float64(5)},
		},
		{
			instructions: []bytecode.Instruction{
				bytecode.New(bytecode.I32LOAD, 42),
				bytecode.New(bytecode.I32TOSTR),
			},
			stack: []Value{String("42")},
		},
		{
			instructions: []bytecode.Instruction{
				bytecode.New(bytecode.F64LOAD, math.Float64bits(1)),
			},
			stack: []Value{Float64(1)},
		},
		{
			instructions: []bytecode.Instruction{
				bytecode.New(bytecode.F64LOAD, math.Float64bits(1)),
				bytecode.New(bytecode.F64LOAD, math.Float64bits(2)),
				bytecode.New(bytecode.F64ADD),
			},
			stack: []Value{Float64(3)},
		},
		{
			instructions: []bytecode.Instruction{
				bytecode.New(bytecode.F64LOAD, math.Float64bits(1)),
				bytecode.New(bytecode.F64LOAD, math.Float64bits(2)),
				bytecode.New(bytecode.F64SUB),
			},
			stack: []Value{Float64(-1)},
		},
		{
			instructions: []bytecode.Instruction{
				bytecode.New(bytecode.F64LOAD, math.Float64bits(1)),
				bytecode.New(bytecode.F64LOAD, math.Float64bits(2)),
				bytecode.New(bytecode.F64MUL),
			},
			stack: []Value{Float64(2)},
		},
		{
			instructions: []bytecode.Instruction{
				bytecode.New(bytecode.F64LOAD, math.Float64bits(1)),
				bytecode.New(bytecode.F64LOAD, math.Float64bits(2)),
				bytecode.New(bytecode.F64DIV),
			},
			stack: []Value{Float64(0.5)},
		},
		{
			instructions: []bytecode.Instruction{
				bytecode.New(bytecode.F64LOAD, math.Float64bits(1)),
				bytecode.New(bytecode.F64LOAD, math.Float64bits(2)),
				bytecode.New(bytecode.F64MOD),
			},
			stack: []Value{Float64(1)},
		},
		{
			instructions: []bytecode.Instruction{
				bytecode.New(bytecode.F64LOAD, math.Float64bits(3.7)),
				bytecode.New(bytecode.F64TOI32),
			},
			stack: []Value{Int32(3)},
		},
		{
			instructions: []bytecode.Instruction{
				bytecode.New(bytecode.F64LOAD, math.Float64bits(1)),
				bytecode.New(bytecode.F64TOSTR),
			},
			stack: []Value{String("1")},
		},
		{
			instructions: []bytecode.Instruction{
				bytecode.New(bytecode.STRLOAD, 0, 3),
			},
			literals: []string{"abc"},
			stack:    []Value{String("abc")},
		},
		{
			instructions: []bytecode.Instruction{
				bytecode.New(bytecode.STRLOAD, 0, 3),
				bytecode.New(bytecode.STRLOAD, 0, 3),
				bytecode.New(bytecode.STRADD),
			},
			literals: []string{"abc"},
			stack:    []Value{String("abcabc")},
		},
		{
			instructions: []bytecode.Instruction{
				bytecode.New(bytecode.STRLOAD, 0, 3),
				bytecode.New(bytecode.STRTOI32),
			},
			literals: []string{"123"},
			stack:    []Value{Int32(123)},
		},
		{
			instructions: []bytecode.Instruction{
				bytecode.New(bytecode.STRLOAD, 0, 1),
				bytecode.New(bytecode.STRTOF64),
			},
			literals: []string{"1"},
			stack:    []Value{Float64(1)},
		},
	}

	for _, tt := range tests {
		var code bytecode.Bytecode
		code.Emit(tt.instructions...)
		for _, c := range tt.literals {
			code.Store([]byte(c + "\x00"))
		}

		t.Run(code.String(), func(t *testing.T) {
			interpreter := New()

			err := interpreter.Execute(code)
			assert.NoError(t, err)

			for _, val := range tt.stack {
				assert.Equal(t, val, interpreter.Pop())
			}
		})
	}
}

func BenchmarkInterpreter_Execute(b *testing.B) {
	tests := []struct {
		instructions []bytecode.Instruction
		literals     []string
	}{
		{
			instructions: []bytecode.Instruction{
				bytecode.New(bytecode.NOP),
			},
		},
		{
			instructions: []bytecode.Instruction{
				bytecode.New(bytecode.POP),
			},
		},
		{
			instructions: []bytecode.Instruction{
				bytecode.New(bytecode.I32LOAD, 1),
				bytecode.New(bytecode.SLTSTORE, 1),
			},
		},
		{
			instructions: []bytecode.Instruction{
				bytecode.New(bytecode.I32LOAD, 1),
				bytecode.New(bytecode.SLTSTORE, 1),
				bytecode.New(bytecode.SLTLOAD, 1),
			},
		},
		{
			instructions: []bytecode.Instruction{
				bytecode.New(bytecode.UNDEFLOAD),
				bytecode.New(bytecode.POP),
			},
		},
		{
			instructions: []bytecode.Instruction{
				bytecode.New(bytecode.UNDEFLOAD),
				bytecode.New(bytecode.UNDEFTOF64),
				bytecode.New(bytecode.POP),
			},
		},
		{
			instructions: []bytecode.Instruction{
				bytecode.New(bytecode.UNDEFLOAD),
				bytecode.New(bytecode.UNDEFTOSTR),
				bytecode.New(bytecode.POP),
			},
		},
		{
			instructions: []bytecode.Instruction{
				bytecode.New(bytecode.NULLLOAD),
				bytecode.New(bytecode.POP),
			},
		},
		{
			instructions: []bytecode.Instruction{
				bytecode.New(bytecode.NULLLOAD),
				bytecode.New(bytecode.NULLTOI32),
				bytecode.New(bytecode.POP),
			},
		},
		{
			instructions: []bytecode.Instruction{
				bytecode.New(bytecode.NULLLOAD),
				bytecode.New(bytecode.NULLTOSTR),
				bytecode.New(bytecode.POP),
			},
		},
		{
			instructions: []bytecode.Instruction{
				bytecode.New(bytecode.BOOLLOAD, 1),
				bytecode.New(bytecode.POP),
			},
		},
		{
			instructions: []bytecode.Instruction{
				bytecode.New(bytecode.BOOLLOAD, 1),
				bytecode.New(bytecode.BOOLTOI32),
				bytecode.New(bytecode.POP),
			},
		},
		{
			instructions: []bytecode.Instruction{
				bytecode.New(bytecode.BOOLLOAD, 1),
				bytecode.New(bytecode.BOOLTOSTR),
				bytecode.New(bytecode.POP),
			},
		},
		{
			instructions: []bytecode.Instruction{
				bytecode.New(bytecode.I32LOAD, 1),
				bytecode.New(bytecode.POP),
			},
		},
		{
			instructions: []bytecode.Instruction{
				bytecode.New(bytecode.I32LOAD, 1),
				bytecode.New(bytecode.I32LOAD, 2),
				bytecode.New(bytecode.I32ADD),
				bytecode.New(bytecode.POP),
			},
		},
		{
			instructions: []bytecode.Instruction{
				bytecode.New(bytecode.I32LOAD, 1),
				bytecode.New(bytecode.I32LOAD, 2),
				bytecode.New(bytecode.I32SUB),
				bytecode.New(bytecode.POP),
			},
		},
		{
			instructions: []bytecode.Instruction{
				bytecode.New(bytecode.I32LOAD, 1),
				bytecode.New(bytecode.I32LOAD, 2),
				bytecode.New(bytecode.I32MUL),
				bytecode.New(bytecode.POP),
			},
		},
		{
			instructions: []bytecode.Instruction{
				bytecode.New(bytecode.I32LOAD, 6),
				bytecode.New(bytecode.I32LOAD, 2),
				bytecode.New(bytecode.I32DIV),
				bytecode.New(bytecode.POP),
			},
		},
		{
			instructions: []bytecode.Instruction{
				bytecode.New(bytecode.I32LOAD, 7),
				bytecode.New(bytecode.I32LOAD, 3),
				bytecode.New(bytecode.I32MOD),
				bytecode.New(bytecode.POP),
			},
		},
		{
			instructions: []bytecode.Instruction{
				bytecode.New(bytecode.I32LOAD, 5),
				bytecode.New(bytecode.I32LOAD, 1),
				bytecode.New(bytecode.I32DIV),
				bytecode.New(bytecode.POP),
			},
		},
		{
			instructions: []bytecode.Instruction{
				bytecode.New(bytecode.I32LOAD, 1),
				bytecode.New(bytecode.I32TOBOOL),
			},
		},
		{
			instructions: []bytecode.Instruction{
				bytecode.New(bytecode.I32LOAD, 5),
				bytecode.New(bytecode.I32TOF64),
				bytecode.New(bytecode.POP),
			},
		},
		{
			instructions: []bytecode.Instruction{
				bytecode.New(bytecode.I32LOAD, 42),
				bytecode.New(bytecode.I32TOSTR),
				bytecode.New(bytecode.POP),
			},
		},
		{
			instructions: []bytecode.Instruction{
				bytecode.New(bytecode.F64LOAD, math.Float64bits(1)),
				bytecode.New(bytecode.POP),
			},
		},
		{
			instructions: []bytecode.Instruction{
				bytecode.New(bytecode.F64LOAD, math.Float64bits(1)),
				bytecode.New(bytecode.F64LOAD, math.Float64bits(2)),
				bytecode.New(bytecode.F64ADD),
				bytecode.New(bytecode.POP),
			},
		},
		{
			instructions: []bytecode.Instruction{
				bytecode.New(bytecode.F64LOAD, math.Float64bits(1)),
				bytecode.New(bytecode.F64LOAD, math.Float64bits(2)),
				bytecode.New(bytecode.F64SUB),
				bytecode.New(bytecode.POP),
			},
		},
		{
			instructions: []bytecode.Instruction{
				bytecode.New(bytecode.F64LOAD, math.Float64bits(1)),
				bytecode.New(bytecode.F64LOAD, math.Float64bits(2)),
				bytecode.New(bytecode.F64MUL),
				bytecode.New(bytecode.POP),
			},
		},
		{
			instructions: []bytecode.Instruction{
				bytecode.New(bytecode.F64LOAD, math.Float64bits(1)),
				bytecode.New(bytecode.F64LOAD, math.Float64bits(2)),
				bytecode.New(bytecode.F64DIV),
				bytecode.New(bytecode.POP),
			},
		},
		{
			instructions: []bytecode.Instruction{
				bytecode.New(bytecode.F64LOAD, math.Float64bits(1)),
				bytecode.New(bytecode.F64LOAD, math.Float64bits(2)),
				bytecode.New(bytecode.F64MOD),
				bytecode.New(bytecode.POP),
			},
		},
		{
			instructions: []bytecode.Instruction{
				bytecode.New(bytecode.F64LOAD, math.Float64bits(3.7)),
				bytecode.New(bytecode.F64TOI32),
				bytecode.New(bytecode.POP),
			},
		},
		{
			instructions: []bytecode.Instruction{
				bytecode.New(bytecode.F64LOAD, math.Float64bits(1)),
				bytecode.New(bytecode.F64TOSTR),
				bytecode.New(bytecode.POP),
			},
		},
		{
			instructions: []bytecode.Instruction{
				bytecode.New(bytecode.STRLOAD, 0, 3),
				bytecode.New(bytecode.POP),
			},
			literals: []string{"abc"},
		},
		{
			instructions: []bytecode.Instruction{
				bytecode.New(bytecode.STRLOAD, 0, 3),
				bytecode.New(bytecode.STRLOAD, 0, 3),
				bytecode.New(bytecode.STRADD),
				bytecode.New(bytecode.POP),
			},
			literals: []string{"abc"},
		},
		{
			instructions: []bytecode.Instruction{
				bytecode.New(bytecode.STRLOAD, 0, 3),
				bytecode.New(bytecode.STRTOI32),
				bytecode.New(bytecode.POP),
			},
			literals: []string{"123"},
		},
		{
			instructions: []bytecode.Instruction{
				bytecode.New(bytecode.STRLOAD, 0, 1),
				bytecode.New(bytecode.STRTOF64),
				bytecode.New(bytecode.POP),
			},
			literals: []string{"1"},
		},
	}

	for _, tt := range tests {
		var code bytecode.Bytecode
		code.Emit(tt.instructions...)
		for _, c := range tt.literals {
			code.Store([]byte(c + "\x00"))
		}

		b.Run(code.String(), func(b *testing.B) {
			interpreter := New()
			b.ResetTimer()

			for i := 0; i < b.N; i++ {
				err := interpreter.Execute(code)
				assert.NoError(b, err)
			}
		})
	}
}
