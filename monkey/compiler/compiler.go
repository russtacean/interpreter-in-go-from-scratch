package compiler

import (
	"fmt"
	"monkey/ast"
	"monkey/code"
	"monkey/object"
)

type Compiler struct {
	instructions           code.Instructions
	constants              []object.Object
	lastInstruction        EmittedInstruction
	penultimateInstruction EmittedInstruction
	symbolTable            *SymbolTable
}

type Bytecode struct {
	Instructions code.Instructions
	Constants    []object.Object
}

type EmittedInstruction struct {
	Opcode   code.Opcode
	Position int
}

func New() *Compiler {
	return &Compiler{
		instructions:           code.Instructions{},
		constants:              []object.Object{},
		lastInstruction:        EmittedInstruction{},
		penultimateInstruction: EmittedInstruction{},
		symbolTable:            NewSymbolTable(),
	}
}

func NewWithState(symbolTable *SymbolTable, constants []object.Object) *Compiler {
	compiler := New()
	compiler.symbolTable = symbolTable
	compiler.constants = constants
	return compiler
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

	case *ast.LetStatement:
		err := compiler.Compile(node.Value)
		if err != nil {
			return err
		}
		symbol := compiler.symbolTable.Define(node.Name.Value)
		compiler.emit(code.OpSetGlobal, symbol.Index)

	case *ast.BlockStatement:
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

	case *ast.IfExpression:
		err := compiler.Compile(node.Condition)
		if err != nil {
			return err
		}

		// Emit OpJumpNotTruthy with bogus value, we will backpatch later
		jumpNotTruthyPos := compiler.emit(code.OpJumpNotTruthy, 9999)

		err = compiler.Compile(node.Consequence)
		if err != nil {
			return err
		}

		if compiler.isLastInstructionPop() {
			compiler.removeLastPop()
		}

		// Emit bogus jump value to backpatch later
		jumpPos := compiler.emit(code.OpJump, 9999)

		afterConsequencePos := len(compiler.instructions)
		compiler.changeOperand(jumpNotTruthyPos, afterConsequencePos)

		if node.Alternative == nil {
			compiler.emit(code.OpNull)
		} else {
			err = compiler.Compile(node.Alternative)
			if err != nil {
				return err
			}

			if compiler.isLastInstructionPop() {
				compiler.removeLastPop()
			}
		}

		afterAlternativePos := len(compiler.instructions)
		compiler.changeOperand(jumpPos, afterAlternativePos)

	case *ast.Identifier:
		symbol, ok := compiler.symbolTable.Resolve(node.Value)
		if !ok {
			return fmt.Errorf("undefined variable %s", node.Value)
		}
		compiler.emit(code.OpGetGlobal, symbol.Index)

	case *ast.IntegerLiteral:
		integer := &object.Integer{Value: node.Value}
		compiler.emit(code.OpConstant, compiler.addConstant(integer))

	case *ast.StringLiteral:
		string := &object.String{Value: node.Value}
		compiler.emit(code.OpConstant, compiler.addConstant(string))

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

func (compiler *Compiler) emit(opcode code.Opcode, operands ...int) int {
	ins := code.Make(opcode, operands...)
	pos := compiler.addInstruction(ins)
	compiler.setLastInstruction(opcode, pos)
	return pos
}

func (compiler *Compiler) setLastInstruction(opcode code.Opcode, position int) {
	penultimate := compiler.lastInstruction
	last := EmittedInstruction{Opcode: opcode, Position: position}

	compiler.lastInstruction = last
	compiler.penultimateInstruction = penultimate
}

func (compiler *Compiler) addInstruction(ins []byte) int {
	posNewInstruction := len(compiler.instructions)
	compiler.instructions = append(compiler.instructions, ins...)
	return posNewInstruction
}

func (compiler *Compiler) isLastInstructionPop() bool {
	return compiler.lastInstruction.Opcode == code.OpPop
}

func (compiler *Compiler) removeLastPop() {
	compiler.instructions = compiler.instructions[:compiler.lastInstruction.Position]
	compiler.lastInstruction = compiler.penultimateInstruction
}

func (compiler *Compiler) replaceInstruction(pos int, newInstruction []byte) {
	for i := 0; i < len(newInstruction); i++ {
		compiler.instructions[pos+i] = newInstruction[i]
	}
}

func (compiler *Compiler) changeOperand(opPos int, operand int) {
	opcode := code.Opcode(compiler.instructions[opPos])
	newInstruction := code.Make(opcode, operand)
	compiler.replaceInstruction(opPos, newInstruction)

}

func (compiler *Compiler) Bytecode() *Bytecode {
	return &Bytecode{
		Instructions: compiler.instructions,
		Constants:    compiler.constants,
	}
}
