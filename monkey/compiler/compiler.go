package compiler

import (
	"fmt"
	"monkey/ast"
	"monkey/code"
	"monkey/object"
)

type Compiler struct {
	instructions code.Instructions
	constants    []object.Object
}

func New() *Compiler {
	return &Compiler{
		instructions: code.Instructions{},
		constants:    []object.Object{},
	}
}

func (compiler *Compiler) Compile(node ast.Node) error {
	switch node := node.(type) {
	case *ast.Program:
		for _, s := range node.Statements {
			err := compiler.Compile(s)
			if err != nil {
				return err
			}
		}

	case *ast.ExpressionStatement:
		err := compiler.Compile(node.Expression)
		if err != nil {
			return err
		}
		compiler.emit(code.OpPop)

	case *ast.PrefixExpression:
		err := compiler.Compile(node.Right)
		if err != nil {
			return err
		}

		switch node.Operator {
		case "!":
			compiler.emit(code.OpBang)
		case "-":
			compiler.emit(code.OpMinus)
		default:
			return fmt.Errorf("unknown operator %s", node.Operator)
		}

	case *ast.InfixExpression:
		if node.Operator == "<" { // Special case to reduce instruction set
			err := compiler.Compile(node.Right)
			if err != nil {
				return err
			}

			err = compiler.Compile(node.Left)
			if err != nil {
				return err
			}
			compiler.emit(code.OpGreaterThan)
			return nil
		}

		err := compiler.Compile(node.Left)
		if err != nil {
			return err
		}

		err = compiler.Compile(node.Right)
		if err != nil {
			return err
		}

		switch node.Operator {
		case "+":
			compiler.emit(code.OpAdd)
		case "-":
			compiler.emit(code.OpSub)
		case "*":
			compiler.emit(code.OpMul)
		case "/":
			compiler.emit(code.OpDiv)
		case ">":
			compiler.emit(code.OpGreaterThan)
		case "==":
			compiler.emit(code.OpEqual)
		case "!=":
			compiler.emit(code.OpNotEqual)
		default:
			return fmt.Errorf("unknown operator %s", node.Operator)
		}

	case *ast.IntegerLiteral:
		integer := &object.Integer{Value: node.Value}
		compiler.emit(code.OpConstant, compiler.addConstant(integer))

	case *ast.BooleanLiteral:
		if node.Value {
			compiler.emit(code.OpTrue)
		} else {
			compiler.emit(code.OpFalse)
		}
	}

	return nil
}

func (compiler *Compiler) addConstant(obj object.Object) int {
	compiler.constants = append(compiler.constants, obj)
	return len(compiler.constants) - 1
}

func (compiler *Compiler) emit(op code.Opcode, operands ...int) int {
	ins := code.Make(op, operands...)
	pos := compiler.addInstruction(ins)
	return pos
}

func (compiler *Compiler) addInstruction(ins []byte) int {
	posNewInstruction := len(compiler.instructions)
	compiler.instructions = append(compiler.instructions, ins...)
	return posNewInstruction
}

func (compiler *Compiler) Bytecode() *Bytecode {
	return &Bytecode{
		Instructions: compiler.instructions,
		Constants:    compiler.constants,
	}
}

type Bytecode struct {
	Instructions code.Instructions
	Constants    []object.Object
}
