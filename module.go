/* module.go - Go wrapper for Emacs module API.

Copyright (C) 2016 Yann Hodique <yann.hodique@gmail.com>.

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or (at
your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program.  If not, see <http://www.gnu.org/licenses/>.  */

package goemacs

/*
#include "include/wrapper.h"
*/
import "C"
import (
	"sync"
	"unsafe"
)

var initFuncs = make([]func(*Environment), 0)

type Environment struct {
	// FIXME: for some reason, struct_emacs_env doesn't compile
	env    *C.struct_emacs_env_25
	stdlib *StdLib
}

func Register(f func(*Environment)) {
	initFuncs = append(initFuncs, f)
}

//export emacs_module_init
func emacs_module_init(e *C.struct_emacs_runtime) C.int {
	env := Environment{
		env: C.GetEnvironment(e),
	}

	for _, f := range initFuncs {
		f(&env)
	}
	return 0
}

type StdLib struct {
	env         *Environment
	messageFunc C.emacs_value
	fsetFunc    C.emacs_value
	Nil         Value
}

func (e *Environment) StdLib() *StdLib {
	if e.stdlib == nil {
		e.stdlib = &StdLib{
			env:         e,
			messageFunc: C.Intern(e.env, C.CString("message")),
			fsetFunc:    C.Intern(e.env, C.CString("fset")),
			Nil:         Value{C.Intern(e.env, C.CString("nil"))},
		}
	}
	return e.stdlib
}

type Value struct {
	val C.emacs_value
}

type Symbol struct {
	Value
}

type String struct {
	Value
}

type Function struct {
	Value
}

func (v Value) ToString() string {
	// FIXME: actually compute string...
	return ""
}

func (e *Environment) String(s string) String {
	return String{
		Value{
			C.MakeString(e.env, C.CString(s), C.int(len(s))),
		},
	}
}

func (stdlib *StdLib) Message(s string) {
	str := stdlib.env.String(s)
	C.Funcall(stdlib.env.env, stdlib.messageFunc, 1, &str.val)
}

func (stdlib *StdLib) Intern(s string) Symbol {
	return Symbol{
		Value{
			C.Intern(stdlib.env.env, C.CString(s)),
		},
	}
}

func (stdlib *StdLib) Fset(sym Symbol, f Function) {
	args := [2]C.emacs_value{sym.val, f.val}
	C.Funcall(stdlib.env.env, stdlib.fsetFunc, 2, &args[0])
}

type FunctionType func(*Environment, int, []Value, interface{}) Value

type FunctionEntry struct {
	f     FunctionType
	arity int
	doc   string
	data  interface{}
}

// from http://stackoverflow.com/questions/37157379/passing-function-pointer-to-the-c-code-using-cgo
var mu sync.Mutex
var index int
var fns = make(map[int]*FunctionEntry)

func register(fn *FunctionEntry) int {
	mu.Lock()
	defer mu.Unlock()
	index++
	for fns[index] != nil {
		index++
	}
	fns[index] = fn
	return index
}

func lookup(i int) *FunctionEntry {
	mu.Lock()
	defer mu.Unlock()
	return fns[i]
}

func unregister(i int) {
	mu.Lock()
	defer mu.Unlock()
	delete(fns, i)
}

func (e *Environment) MakeFunction(f FunctionType, arity int, doc string) Function {
	cArity := C.int(arity)
	idx := register(&FunctionEntry{
		f:     f,
		arity: arity,
		doc:   doc,
	})

	return Function{
		Value{
			C.MakeFunction(e.env, cArity, cArity,
				C.CString(doc), C.ptrdiff_t(idx)),
		},
	}
}

//export emacs_call_function
func emacs_call_function(
	//FIXME: emacs_env_25 shouldn't be used
	env *C.struct_emacs_env_25, nargs C.ptrdiff_t,
	args *C.emacs_value, idx C.ptrdiff_t) C.emacs_value {

	n := int(nargs)
	pargs := (*[1 << 30]C.emacs_value)(unsafe.Pointer(args))
	arguments := make([]Value, n)
	for i := 0; i < n; i++ {
		arguments[i] = Value{
			val: pargs[i],
		}
	}
	entry := lookup(int(idx))
	return entry.f(
		&Environment{
			env: env,
		},
		n,
		arguments,
		entry.data,
	).val
}
