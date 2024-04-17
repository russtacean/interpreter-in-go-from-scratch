package vm

import (
	"monkey/code"
	"monkey/object"
)

type Frame struct {
	fn          *object.CompiledFunction
	ip          int
	basePointer int // Points to bottom of stack of current call frame
}

func NewFrame(fn *object.CompiledFunction, basePointer int) *Frame {
	return &Frame{fn: fn, ip: -1, basePointer: basePointer}
}

func (frame *Frame) Instructions() code.Instructions {
	return frame.fn.Instructions
}
