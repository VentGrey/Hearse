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
}

type AST struct {
	Instructions []*Instruction
}
