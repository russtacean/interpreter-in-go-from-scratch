package compiler

type SymbolScope string

const (
	GlobalScope   SymbolScope = "GLOBAL"
	LocalScope    SymbolScope = "LOCAL"
	BuiltinScope  SymbolScope = "BUILTIN"
	FreeScope     SymbolScope = "FREE"
	FunctionScope SymbolScope = "FUNCTION"
)

type Symbol struct {
	Name  string
	Scope SymbolScope
	Index int
}

type SymbolTable struct {
	Outer          *SymbolTable // Similar to outer env in tree walking interpreter
	store          map[string]Symbol
	numDefinitions int
	FreeSymbols    []Symbol
}

func NewSymbolTable() *SymbolTable {
	store := make(map[string]Symbol)
	free := []Symbol{}
	return &SymbolTable{store: store, FreeSymbols: free}
}

func NewEnclosedSymbolTable(outer *SymbolTable) *SymbolTable {
	s := NewSymbolTable()
	s.Outer = outer
	return s
}

func (symbolTable *SymbolTable) Define(name string) Symbol {
	symbol := Symbol{Name: name, Index: symbolTable.numDefinitions}
	if symbolTable.Outer == nil {
		symbol.Scope = GlobalScope
	} else {
		symbol.Scope = LocalScope
	}

	symbolTable.store[name] = symbol
	symbolTable.numDefinitions++
	return symbol
}

func (symbolTable *SymbolTable) defineFree(original Symbol) Symbol {
	symbolTable.FreeSymbols = append(symbolTable.FreeSymbols, original)

	symbol := Symbol{Name: original.Name, Index: len(symbolTable.FreeSymbols) - 1}
	symbol.Scope = FreeScope

	symbolTable.store[original.Name] = symbol
	return symbol
}

func (symbolTable *SymbolTable) DefineFunctionName(name string) Symbol {
	symbol := Symbol{Name: name, Index: 0, Scope: FunctionScope}
	symbolTable.store[name] = symbol
	return symbol
}

func (symbolTable *SymbolTable) DefineBuiltin(index int, name string) Symbol {
	symbol := Symbol{Name: name, Index: index, Scope: BuiltinScope}
	symbolTable.store[name] = symbol
	return symbol
}

func (symbolTable *SymbolTable) Resolve(name string) (Symbol, bool) {
	symbol, ok := symbolTable.store[name]
	if !ok && symbolTable.Outer != nil {
		symbol, ok = symbolTable.Outer.Resolve(name)
		if !ok {
			return symbol, ok
		}

		if symbol.Scope == GlobalScope || symbol.Scope == BuiltinScope {
			return symbol, ok
		}

		free := symbolTable.defineFree(symbol)
		return free, true
	}
	return symbol, ok
}
