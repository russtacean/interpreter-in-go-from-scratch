package object

import "fmt"

var Builtins = []struct {
	Name    string
	Builtin *Builtin
}{
	{
		"len",
		&Builtin{Fn: func(args ...Object) Object {
			if len(args) != 1 {
				return newError("wrong number of arguments to `len`. got=%d, want=1",
					len(args))
			}

			switch arg := args[0].(type) {
			case *Array:
				return &Integer{Value: int64(len(arg.Elements))}
			case *String:
				return &Integer{Value: int64(len(arg.Value))}
			default:
				return newError("argument to `len` not supported, got %s",
					args[0].Type())
			}
		},
		},
	},
	{
		"first",
		&Builtin{
			Fn: func(args ...Object) Object {
				if len(args) != 1 {
					return newError("wrong number of arguments to `first`. got=%d, want=1", len(args))
				}

				switch args[0].(type) {
				case *Array:
					arr := args[0].(*Array)
					if len(arr.Elements) == 0 {
						return nil
					}
					return arr.Elements[0]
				default:
					return newError("argument to `first` not supported, got %s", args[0].Type())
				}
			},
		},
	},
	{
		"last",
		&Builtin{
			Fn: func(args ...Object) Object {
				if len(args) != 1 {
					return newError("wrong number of arguments to `last`. got=%d, want=1", len(args))
				}

				switch args[0].(type) {
				case *Array:
					arr := args[0].(*Array)
					length := len(arr.Elements)
					if length == 0 {
						return nil
					}
					return arr.Elements[length-1]
				default:
					return newError("argument to `last` not supported, got %s", args[0].Type())
				}

			},
		},
	},
	{"rest",
		&Builtin{
			Fn: func(args ...Object) Object {
				if len(args) != 1 {
					return newError("wrong number of arguments to `rest`. got=%d, want=1", len(args))
				}

				switch args[0].(type) {
				case *Array:
					arr := args[0].(*Array)
					length := len(arr.Elements)
					if length == 0 {
						return nil
					}
					newElements := make([]Object, length-1, length-1)
					copy(newElements, arr.Elements[1:])
					return &Array{Elements: newElements}
				default:
					return newError("argument to `rest` not supported, got %s", args[0].Type())
				}
			},
		},
	},
	{
		"push",
		&Builtin{
			Fn: func(args ...Object) Object {
				if len(args) != 2 {
					return newError("wrong number of arguments to `rest`. got=%d, want=1", len(args))
				}

				switch args[0].(type) {
				case *Array:
					arr := args[0].(*Array)
					length := len(arr.Elements)
					if length == 0 {
						return nil
					}
					newElements := make([]Object, length+1, length+1)
					copy(newElements, arr.Elements)
					newElements[length] = args[1]
					return &Array{Elements: newElements}
				default:
					return newError("argument to `rest` not supported, got %s", args[0].Type())
				}
			},
		},
	},
	{
		"puts",
		&Builtin{
			Fn: func(args ...Object) Object {
				for _, arg := range args {
					fmt.Println(arg.Inspect())
				}

				return nil
			},
		},
	},
}

func newError(format string, a ...interface{}) *Error {
	return &Error{Message: fmt.Sprintf(format, a...)}
}

func GetBuiltinByName(name string) *Builtin {
	for _, def := range Builtins {
		if def.Name == name {
			fmt.Println(def.Name)
			return def.Builtin
		}
	}
	return nil
}
