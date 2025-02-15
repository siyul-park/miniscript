package interpreter

import (
	"encoding/binary"
	"fmt"
	"math"
	"strconv"

	"github.com/siyul-park/minijs/bytecode"
)

type Interpreter struct {
	stack  []byte
	frames []*Frame
	sp     int
	fp     int
}

func New() *Interpreter {
	return &Interpreter{
		stack:  make([]byte, 64),
		frames: make([]*Frame, 64),
	}
}

func (i *Interpreter) Top() any {
	return i.decode(i.top())
}

func (i *Interpreter) Execute(code bytecode.Bytecode) error {
	frame := NewFrame(code, 0)
	insns := frame.Instructions()
	consts := frame.Constants()

	i.exec(frame)
	defer i.exit()

	for frame.ip < len(insns)-1 {
		frame.ip++

		ip := frame.ip
		opcode := bytecode.Opcode(insns[ip])

		switch opcode {
		case bytecode.NOP:
		case bytecode.POP:
			i.pop()
		case bytecode.I32LOAD:
			val := int32(binary.BigEndian.Uint32(insns[frame.ip+1:]))
			i.push32(INT32, uint32(val))
			frame.ip += 4
		case bytecode.I32MUL:
			_, val2 := i.pop32()
			_, val1 := i.pop32()
			i1 := int32(val1)
			i2 := int32(val2)
			i.push32(INT32, uint32(i1*i2))
		case bytecode.I32ADD:
			_, val2 := i.pop32()
			_, val1 := i.pop32()
			i1 := int32(val1)
			i2 := int32(val2)
			i.push32(INT32, uint32(i1+i2))
		case bytecode.I32SUB:
			_, val2 := i.pop32()
			_, val1 := i.pop32()
			i1 := int32(val1)
			i2 := int32(val2)
			i.push32(INT32, uint32(i1-i2))
		case bytecode.I32DIV:
			_, val2 := i.pop32()
			_, val1 := i.pop32()
			i1 := int32(val1)
			i2 := int32(val2)
			i.push32(INT32, uint32(i1/i2))
		case bytecode.I32MOD:
			_, val2 := i.pop32()
			_, val1 := i.pop32()
			i1 := int32(val1)
			i2 := int32(val2)
			i.push32(INT32, uint32(i1%i2))
		case bytecode.I322F64:
			_, val := i.pop32()
			i.push64(FLOAT64, math.Float64bits(float64(val)))
		case bytecode.I322C:
			_, val := i.pop32()
			i.push(STRING, []byte(strconv.Itoa(int(val))))
		case bytecode.F64LOAD:
			val := binary.BigEndian.Uint64(insns[frame.ip+1:])
			i.push64(FLOAT64, val)
			frame.ip += 8
		case bytecode.F64ADD:
			_, val2 := i.pop64()
			_, val1 := i.pop64()
			f1 := math.Float64frombits(val1)
			f2 := math.Float64frombits(val2)
			i.push64(FLOAT64, math.Float64bits(f1+f2))
		case bytecode.F64SUB:
			_, val2 := i.pop64()
			_, val1 := i.pop64()
			f1 := math.Float64frombits(val1)
			f2 := math.Float64frombits(val2)
			i.push64(FLOAT64, math.Float64bits(f1-f2))
		case bytecode.F64MUL:
			_, val2 := i.pop64()
			_, val1 := i.pop64()
			f1 := math.Float64frombits(val1)
			f2 := math.Float64frombits(val2)
			i.push64(FLOAT64, math.Float64bits(f1*f2))
		case bytecode.F64DIV:
			_, val2 := i.pop64()
			_, val1 := i.pop64()
			f1 := math.Float64frombits(val1)
			f2 := math.Float64frombits(val2)
			i.push64(FLOAT64, math.Float64bits(f1/f2))
		case bytecode.F64MOD:
			_, val2 := i.pop64()
			_, val1 := i.pop64()
			f1 := math.Float64frombits(val1)
			f2 := math.Float64frombits(val2)
			i.push64(FLOAT64, math.Float64bits(math.Mod(f1, f2)))
		case bytecode.F64I32:
			_, val := i.pop64()
			f := math.Float64frombits(val)
			i.push32(INT32, uint32(int32(f)))
		case bytecode.F642C:
			_, val := i.pop64()
			f := math.Float64frombits(val)
			i.push(STRING, []byte(strconv.FormatFloat(f, 'f', -1, 64)))
		case bytecode.CLOAD:
			offset := int(binary.BigEndian.Uint32(insns[frame.ip+1:]))
			size := int(binary.BigEndian.Uint32(insns[frame.ip+5:]))
			i.push(STRING, consts[offset:offset+size])
			frame.ip += 8
		case bytecode.CADD:
			_, val2 := i.pop()
			_, val1 := i.pop()
			i.push(STRING, append(val1, val2...))
		case bytecode.C2F64:
			_, val := i.pop()
			f, err := strconv.ParseFloat(string(val), 64)
			if err != nil {
				f = math.NaN()
			}
			i.push64(FLOAT64, math.Float64bits(f))
		case bytecode.C2I32:
			_, val := i.pop()
			n, err := strconv.Atoi(string(val))
			if err != nil {
				n = 0
			}
			i.push32(INT32, uint32(n))
		default:
			typ := bytecode.TypeOf(opcode)
			if typ == nil {
				return fmt.Errorf("unknown opcode: %v", opcode)
			}
			return fmt.Errorf("unknown opcode: %v", typ.Mnemonic)
		}

		frame = i.frame()
		insns = frame.Instructions()
		consts = frame.Constants()
	}
	return nil
}

func (i *Interpreter) frame() *Frame {
	return i.frames[i.fp-1]
}

func (i *Interpreter) exec(f *Frame) {
	if len(i.frames) <= i.fp {
		frames := make([]*Frame, len(i.frames)*2)
		copy(frames, i.frames)
		i.frames = frames
	}

	i.frames[i.fp] = f
	i.fp++
}

func (i *Interpreter) exit() {
	if i.fp == 0 {
		return
	}
	i.fp--
}

func (i *Interpreter) decode(kind Kind, val []byte) any {
	switch kind {
	case INT32:
		v := binary.BigEndian.Uint32(val)
		return int32(v)
	case FLOAT64:
		v := binary.BigEndian.Uint64(val)
		return math.Float64frombits(v)
	case STRING:
		return string(val)
	default:
		return nil
	}
}

func (i *Interpreter) push(kind Kind, val []byte) {
	size := len(val)
	i.resize(i.sp + size + 9)
	copy(i.stack[i.sp:], val)
	binary.BigEndian.PutUint64(i.stack[i.sp+size:], uint64(size))
	i.stack[i.sp+size+8] = byte(KIND & kind)
	i.sp += size + 9
}

func (i *Interpreter) pop() (Kind, []byte) {
	if i.sp == 0 {
		return 0, nil
	}

	mark := i.stack[i.sp-1]
	i.sp -= 1

	var size int
	if mark&PRIMITIVE == PRIMITIVE {
		size = int(mark & SIZE)
	} else {
		size = int(binary.BigEndian.Uint64(i.stack[i.sp-8 : i.sp]))
		i.sp -= 8
	}

	val := i.stack[i.sp-size : i.sp]
	i.sp -= size
	return Kind(mark & KIND), val
}

func (i *Interpreter) top() (Kind, []byte) {
	if i.sp == 0 {
		return 0, nil
	}

	sp := i.sp - 1
	mark := i.stack[sp]

	var size int
	if mark&PRIMITIVE == PRIMITIVE {
		size = int(mark & SIZE)
	} else {
		size = int(binary.BigEndian.Uint64(i.stack[sp-8 : sp]))
		sp -= 8
	}

	val := i.stack[sp-size : sp]
	return Kind(mark & KIND), val
}

func (i *Interpreter) push32(kind Kind, val uint32) {
	i.resize(i.sp + 4 + 1)
	binary.BigEndian.PutUint32(i.stack[i.sp:], val)
	i.stack[i.sp+4] = byte(PRIMITIVE | KIND&kind | SIZE&4)
	i.sp += 5
}

func (i *Interpreter) pop32() (Kind, uint32) {
	if i.sp == 0 {
		return 0, 0
	}
	mark := i.stack[i.sp-1]
	i.sp -= 1
	val := binary.BigEndian.Uint32(i.stack[i.sp-4 : i.sp])
	i.sp -= 4
	return Kind(mark & KIND), val
}

func (i *Interpreter) push64(kind Kind, val uint64) {
	i.resize(i.sp + 8 + 1)
	binary.BigEndian.PutUint64(i.stack[i.sp:], val)
	i.stack[i.sp+8] = byte(PRIMITIVE | KIND&kind | SIZE&8)
	i.sp += 9
}

func (i *Interpreter) pop64() (Kind, uint64) {
	if i.sp == 0 {
		return 0, 0
	}
	mark := i.stack[i.sp-1]
	i.sp -= 1
	val := binary.BigEndian.Uint64(i.stack[i.sp-8 : i.sp])
	i.sp -= 8
	return Kind(mark & KIND), val
}

func (i *Interpreter) resize(size int) {
	if len(i.stack) < size {
		stack := make([]byte, size*2)
		copy(stack, i.stack)
		i.stack = stack
	}
}
