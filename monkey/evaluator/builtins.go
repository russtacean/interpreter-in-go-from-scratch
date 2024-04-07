package evaluator

import "monkey/object"

var builtins = map[string]*object.Builtin{
	"len": &object.Builtin{
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments to `len`. got=%d, want=1", len(args))
			}

			switch arg := args[0].(type) {
			case *object.String:
				return &object.Integer{Value: int64(len(arg.Value))}
			case *object.Array:
				return &object.Integer{Value: int64(len(arg.Elements))}
			default:
				return newError("argument to `len` not supported, got %s", args[0].Type())
			}
		},
	},
	"first": &object.Builtin{
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments to `first`. got=%d, want=1", len(args))
			}

			switch args[0].(type) {
			case *object.Array:
				arr := args[0].(*object.Array)
				if len(arr.Elements) == 0 {
					return NULL
				}
				return arr.Elements[0]
			default:
				return newError("argument to `first` not supported, got %s", args[0].Type())
			}
		},
	},
	"last": &object.Builtin{
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments to `last`. got=%d, want=1", len(args))
			}

			switch args[0].(type) {
			case *object.Array:
				arr := args[0].(*object.Array)
				length := len(arr.Elements)
				if length == 0 {
					return NULL
				}
				return arr.Elements[length-1]
			default:
				return newError("argument to `last` not supported, got %s", args[0].Type())
			}
		},
	},
	"rest": &object.Builtin{
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments to `rest`. got=%d, want=1", len(args))
			}

			switch args[0].(type) {
			case *object.Array:
				arr := args[0].(*object.Array)
				length := len(arr.Elements)
				if length == 0 {
					return NULL
				}
				newElements := make([]object.Object, length-1, length-1)
				copy(newElements, arr.Elements[1:])
				return &object.Array{Elements: newElements}
			default:
				return newError("argument to `rest` not supported, got %s", args[0].Type())
			}
		},
	},
	"push": &object.Builtin{
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 2 {
				return newError("wrong number of arguments to `rest`. got=%d, want=1", len(args))
			}

			switch args[0].(type) {
			case *object.Array:
				arr := args[0].(*object.Array)
				length := len(arr.Elements)
				if length == 0 {
					return NULL
				}
				newElements := make([]object.Object, length+1, length+1)
				copy(newElements, arr.Elements)
				newElements[length] = args[1]
				return &object.Array{Elements: newElements}
			default:
				return newError("argument to `rest` not supported, got %s", args[0].Type())
			}
		},
	},
}
