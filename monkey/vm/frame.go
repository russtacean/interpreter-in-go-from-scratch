package vm

import (
	"monkey/code"
	"monkey/object"
)

type Frame struct {
	closure     *object.Closure
	ip          int
	basePointer int // Points to bottom of stack of current call frame
}

func NewFrame(closure *object.Closure, basePointer int) *Frame {
	return &Frame{closure: closure, ip: -1, basePointer: basePointer}
}

func (frame *Frame) Instructions() code.Instructions {
	return frame.closure.Fn.Instructions
}
