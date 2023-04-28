// Hearse is a small compiler + interpreter for Brainfuck written entirely in Go
package main

import (
	"fmt"
	"os"
)

type Op int

type Instruction struct {
	Op 	Op
	Arg int
	Next *Instruction
	Prev *Instruction
	Offset *Instruction
}

type AST struct {
	Instructions []*Instruction
}

const (
	INCREMENT_POINTER Op = iota
	DECREMENT_POINTER
	INCREMENT_VALUE
	DECREMENT_VALUE
	OUTPUT_VALUE
	INPUT_VALUE
	JUMP_FORWARD
	JUMP_BACKWARD
)

func Parse(code string) *AST {
	root := &Instruction{}
	curr := root

	for _, c := range code {
		switch c {
		case '>':
			curr.Next = &Instruction{Op: INCREMENT_POINTER, Prev: curr}
		case '<':
            curr.Next = &Instruction{Op: DECREMENT_POINTER, Prev: curr}
            curr = curr.Next
        case '+':
            curr.Next = &Instruction{Op: INCREMENT_VALUE, Prev: curr}
            curr = curr.Next
        case '-':
            curr.Next = &Instruction{Op: DECREMENT_VALUE, Prev: curr}
            curr = curr.Next
        case '.':
            curr.Next = &Instruction{Op: OUTPUT_VALUE, Prev: curr}
            curr = curr.Next
        case ',':
            curr.Next = &Instruction{Op: INPUT_VALUE, Prev: curr}
            curr = curr.Next
        case '[':
            curr.Next = &Instruction{Op: JUMP_FORWARD, Prev: curr}
            curr = curr.Next
        case ']':
            curr.Next = &Instruction{Op: JUMP_BACKWARD, Prev: curr}
            curr = curr.Next
		}
	}


	root = root.Next

	var stack []*Instruction

	for curr := root; curr != nil; curr = curr.Next {
		switch curr.Op {
		case JUMP_FORWARD:
			stack = append(stack, curr)
		case JUMP_BACKWARD:
			if len(stack) == 0 {
				panic("Unmatched ]!")
			}
			forward := stack[len(stack) - 1]
			stack = stack[:len(stack) - 1]
			curr.Offset = forward
			forward.Offset = curr
		}
	}

	if len(stack) != 0 {
		panic("Unmatched [!")
	}

	ast := &AST{Instructions: make([]*Instruction, 0)}
	for curr := root; curr != nil; curr = curr.Next {
		if curr.Op != JUMP_FORWARD && curr.Op != JUMP_BACKWARD {
			ast.Instructions = append(ast.Instructions, curr)
		}
	}

	return ast
}

func Interpret(ast *AST) {
	tape := make([]byte, 30000)
	ptr := 0

	for _, inst := range ast.Instructions {
		switch inst.Op {
		case INCREMENT_POINTER:
			ptr++
		case DECREMENT_POINTER:
			ptr--
		case INCREMENT_VALUE:
			tape[ptr]++
		case DECREMENT_VALUE:
			tape[ptr]--
		case OUTPUT_VALUE:
			fmt.Printf(string(tape[ptr]))
		case INPUT_VALUE:
			var input byte
			fmt.Scan(&input)
			tape[ptr] = input
		case JUMP_FORWARD:
			if tape[ptr] == 0 {
				inst = inst.Offset
			}
		case JUMP_BACKWARD:
			if tape[ptr] != 0 {
				inst = inst.Offset
			}
		}
	}
}

func Compile(ast *AST) []byte {
    code := make([]byte, 0)
    for _, instr := range ast.Instructions {
        switch instr.Op {
        case INCREMENT_POINTER:
            code = append(code, 0x41) // increment pointer
        case DECREMENT_POINTER:
            code = append(code, 0x42) // decrement pointer
        case INCREMENT_VALUE:
            code = append(code, 0x43) // increment value
        case DECREMENT_VALUE:
            code = append(code, 0x44) // decrement value
        case OUTPUT_VALUE:
            code = append(code, 0x45) // output value
        case INPUT_VALUE:
            code = append(code, 0x46) // input value
        case JUMP_FORWARD:
            code = append(code, 0x47) // jump forward
            code = append(code, make([]byte, 4)...) // reserve space for offset
        case JUMP_BACKWARD:
            code = append(code, 0x48) // jump backward
            code = append(code, make([]byte, 4)...) // reserve space for offset
        }
    }

    // Fill in the jump offset addresses
    for i, instr := range ast.Instructions {
        if instr.Op == JUMP_FORWARD || instr.Op == JUMP_BACKWARD {
            offset := instr.Offset
            if instr.Op == JUMP_FORWARD {
                code[i+1] = byte(offset.Arg >> 24)
                code[i+2] = byte(offset.Arg >> 16)
                code[i+3] = byte(offset.Arg >> 8)
                code[i+4] = byte(offset.Arg)
            } else {
                code[i+1] = byte(instr.Prev.Arg >> 24)
                code[i+2] = byte(instr.Prev.Arg >> 16)
                code[i+3] = byte(instr.Prev.Arg >> 8)
                code[i+4] = byte(instr.Prev.Arg)
            }
        }
    }

    return code
}
