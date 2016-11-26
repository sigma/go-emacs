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
		messageStr := C.CString("message")
		defer C.free(unsafe.Pointer(messageStr))

		fsetStr := C.CString("fset")
		defer C.free(unsafe.Pointer(fsetStr))

		nilStr := C.CString("nil")
		defer C.free(unsafe.Pointer(nilStr))

		e.stdlib = &StdLib{
			env:         e,
			messageFunc: C.Intern(e.env, messageStr),
			fsetFunc:    C.Intern(e.env, fsetStr),
			Nil:         Value{C.Intern(e.env, nilStr)},
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

func (e *Environment) GoString(v Value) string {
	size := C.StringSize(e.env, v.val)
	buffer := C.CopyString(e.env, v.val, size)
	s := C.GoStringN(buffer, C.int(size))
	C.free(unsafe.Pointer(buffer))
	return s
}

func (e *Environment) String(s string) String {
	valStr := C.CString(s)
	defer C.free(unsafe.Pointer(valStr))

	return String{
		Value{
			C.MakeString(e.env, valStr, C.int(len(s))),
		},
	}
}

func (stdlib *StdLib) Message(s string) {
	str := stdlib.env.String(s)
	C.Funcall(stdlib.env.env, stdlib.messageFunc, 1, &str.val)
}

func (stdlib *StdLib) Intern(s string) Symbol {
	valStr := C.CString(s)
	defer C.free(unsafe.Pointer(valStr))

	return Symbol{
		Value{
			C.Intern(stdlib.env.env, valStr),
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

	docStr := C.CString(doc)
	defer C.free(unsafe.Pointer(docStr))

	return Function{
		Value{
			C.MakeFunction(e.env, cArity, cArity,
				docStr, C.ptrdiff_t(idx)),
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
