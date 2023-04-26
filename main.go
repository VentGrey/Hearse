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
