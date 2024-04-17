package compiler

import (
	"fmt"
	"monkey/ast"
	"monkey/code"
	"monkey/object"
	"sort"
)

type Compiler struct {
	constants   []object.Object
	symbolTable *SymbolTable
	scopes      []CompilationScope
	scopeIndex  int
}

type Bytecode struct {
	Instructions code.Instructions
	Constants    []object.Object
}

type EmittedInstruction struct {
	Opcode   code.Opcode
	Position int
}

type CompilationScope struct {
	instructions           code.Instructions
	lastInstruction        EmittedInstruction
	penultimateInstruction EmittedInstruction
}

func New() *Compiler {
	mainScope := CompilationScope{
		instructions: code.Instructions{},
	}
	return &Compiler{
		constants:   []object.Object{},
		symbolTable: NewSymbolTable(),
		scopes:      []CompilationScope{mainScope},
		scopeIndex:  0,
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
		if symbol.Scope == GlobalScope {
			compiler.emit(code.OpSetGlobal, symbol.Index)
		} else {
			compiler.emit(code.OpSetLocal, symbol.Index)
		}

	case *ast.BlockStatement:
		for _, s := range node.Statements {
			err := compiler.Compile(s)
			if err != nil {
				return err
			}
		}

	case *ast.ReturnStatement:
		err := compiler.Compile(node.ReturnValue)
		if err != nil {
			return err
		}

		compiler.emit(code.OpReturnValue)

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

		if compiler.lastInstructionIs(code.OpPop) {
			compiler.removeLastPop()
		}

		// Emit bogus jump value to backpatch later
		jumpPos := compiler.emit(code.OpJump, 9999)

		afterConsequencePos := len(compiler.currentInstructions())
		compiler.changeOperand(jumpNotTruthyPos, afterConsequencePos)

		if node.Alternative == nil {
			compiler.emit(code.OpNull)
		} else {
			err = compiler.Compile(node.Alternative)
			if err != nil {
				return err
			}

			if compiler.lastInstructionIs(code.OpPop) {
				compiler.removeLastPop()
			}
		}

		afterAlternativePos := len(compiler.currentInstructions())
		compiler.changeOperand(jumpPos, afterAlternativePos)

	case *ast.IndexExpression:
		err := compiler.Compile(node.Left)
		if err != nil {
			return err
		}

		err = compiler.Compile(node.Index)
		if err != nil {
			return err
		}

		compiler.emit(code.OpIndex)

	case *ast.Identifier:
		symbol, ok := compiler.symbolTable.Resolve(node.Value)
		if !ok {
			return fmt.Errorf("undefined variable %s", node.Value)
		}

		if symbol.Scope == GlobalScope {
			compiler.emit(code.OpGetGlobal, symbol.Index)
		} else {
			compiler.emit(code.OpGetLocal, symbol.Index)
		}

	case *ast.IntegerLiteral:
		integer := &object.Integer{Value: node.Value}
		compiler.emit(code.OpConstant, compiler.addConstant(integer))

	case *ast.StringLiteral:
		str := &object.String{Value: node.Value}
		compiler.emit(code.OpConstant, compiler.addConstant(str))

	case *ast.BooleanLiteral:
		if node.Value {
			compiler.emit(code.OpTrue)
		} else {
			compiler.emit(code.OpFalse)
		}

	case *ast.ArrayLiteral:
		for _, el := range node.Elements {
			err := compiler.Compile(el)
			if err != nil {
				return err
			}
		}
		compiler.emit(code.OpArray, len(node.Elements))

	case *ast.HashLiteral:
		keys := []ast.Expression{}
		for key := range node.Pairs {
			keys = append(keys, key)
		}
		// Go doesn't guarantee ordering when iterating through map.
		// This is not a required sort, but will keep tests from breaking randomly
		sort.Slice(keys, func(i, j int) bool {
			return keys[i].String() < keys[j].String()
		})

		for _, k := range keys {
			err := compiler.Compile(k)
			if err != nil {
				return err
			}
			err = compiler.Compile(node.Pairs[k])
			if err != nil {
				return err
			}
		}

		compiler.emit(code.OpHash, len(node.Pairs)*2)

	case *ast.FunctionLiteral:
		compiler.enterScope()

		for _, param := range node.Parameters {
			compiler.symbolTable.Define(param.Value)
		}

		err := compiler.Compile(node.Body)
		if err != nil {
			return err
		}

		if compiler.lastInstructionIs(code.OpPop) {
			compiler.replaceLastPopWithReturn()
		}
		if !compiler.lastInstructionIs(code.OpReturnValue) {
			compiler.emit(code.OpReturn)
		}

		numLocals := compiler.symbolTable.numDefinitions
		instructions := compiler.leaveScope()
		compiledFn := &object.CompiledFunction{Instructions: instructions, NumLocals: numLocals}

		compiler.emit(code.OpConstant, compiler.addConstant(compiledFn))

	case *ast.CallExpression:
		err := compiler.Compile(node.Function)
		if err != nil {
			return err
		}

		for _, arg := range node.Arguments {
			err := compiler.Compile(arg)
			if err != nil {
				return err
			}
		}

		compiler.emit(code.OpCall, len(node.Arguments))
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

func (compiler *Compiler) currentScope() *CompilationScope {
	return &compiler.scopes[compiler.scopeIndex]
}

func (compiler *Compiler) currentInstructions() code.Instructions {
	return compiler.currentScope().instructions
}

func (compiler *Compiler) setLastInstruction(opcode code.Opcode, position int) {
	currentScope := compiler.currentScope()
	penultimate := currentScope.lastInstruction
	last := EmittedInstruction{Opcode: opcode, Position: position}

	currentScope.lastInstruction = last
	currentScope.penultimateInstruction = penultimate
}

func (compiler *Compiler) addInstruction(ins []byte) int {
	posNewInstruction := len(compiler.currentInstructions())
	updatedInstructions := append(compiler.currentInstructions(), ins...)
	compiler.scopes[compiler.scopeIndex].instructions = updatedInstructions
	return posNewInstruction
}

func (compiler *Compiler) lastInstructionIs(opcode code.Opcode) bool {
	currentScope := compiler.currentScope()
	if len(currentScope.instructions) == 0 {
		return false
	}

	return currentScope.lastInstruction.Opcode == opcode
}

func (compiler *Compiler) removeLastPop() {
	currentScope := compiler.currentScope()
	last := currentScope.lastInstruction
	previous := currentScope.penultimateInstruction

	oldIns := compiler.currentInstructions()
	newIns := oldIns[:last.Position]

	currentScope.instructions = newIns
	currentScope.lastInstruction = previous
}

func (compiler *Compiler) replaceInstruction(pos int, newInstruction []byte) {
	instructions := compiler.currentInstructions()
	for i := 0; i < len(newInstruction); i++ {
		instructions[pos+i] = newInstruction[i]
	}
}

func (compiler *Compiler) replaceLastPopWithReturn() {
	currentScope := compiler.currentScope()
	lastPos := currentScope.lastInstruction.Position
	compiler.replaceInstruction(lastPos, code.Make(code.OpReturnValue))

	currentScope.lastInstruction.Opcode = code.OpReturnValue
}

func (compiler *Compiler) changeOperand(opPos int, operand int) {
	opcode := code.Opcode(compiler.currentInstructions()[opPos])
	newInstruction := code.Make(opcode, operand)
	compiler.replaceInstruction(opPos, newInstruction)

}

func (compiler *Compiler) Bytecode() *Bytecode {
	return &Bytecode{
		Instructions: compiler.currentInstructions(),
		Constants:    compiler.constants,
	}
}

func (compiler *Compiler) enterScope() {
	scope := CompilationScope{
		instructions:           code.Instructions{},
		lastInstruction:        EmittedInstruction{},
		penultimateInstruction: EmittedInstruction{},
	}

	compiler.scopes = append(compiler.scopes, scope)
	compiler.scopeIndex++
	compiler.symbolTable = NewEnclosedSymbolTable(compiler.symbolTable)
}

func (compiler *Compiler) leaveScope() code.Instructions {
	instructions := compiler.currentInstructions()

	compiler.scopes = compiler.scopes[:len(compiler.scopes)-1]
	compiler.scopeIndex--

	compiler.symbolTable = compiler.symbolTable.Outer

	return instructions
}
