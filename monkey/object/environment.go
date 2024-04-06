package object

import (
	"bytes"
	"fmt"
)

type Environment struct {
	store map[string]Object
	outer *Environment
}

func NewEnvironment() *Environment {
	store := make(map[string]Object)
	return &Environment{store: store, outer: nil}
}

func NewEnclosedEnvironment(outer *Environment) *Environment {
	env := NewEnvironment()
	env.outer = outer
	return env
}

func (e *Environment) Get(name string) (Object, bool) {
	obj, ok := e.store[name]
	if !ok && e.outer != nil {
		obj, ok = e.outer.Get(name)
	}
	return obj, ok
}

func (e *Environment) Set(name string, value Object) Object {
	e.store[name] = value
	return value
}

func (e *Environment) String() string {
	var out bytes.Buffer

	out.WriteString("inner: {\n")
	for k, v := range e.store {
		out.WriteString(fmt.Sprintf("    %s: %s,\n", k, v))
	}
	out.WriteString("}")

	if e.outer != nil {
		out.WriteString("\nouter: {\n")
		for k, v := range e.outer.store {
			out.WriteString(fmt.Sprintf("    %s: %s, \n", k, v.Inspect()))
		}
		out.WriteString("}")
	}

	return out.String()
}
