/* stdlib.go - Go wrapper for Emacs module API.

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
import "unsafe"

type Value interface {
	getVal() C.emacs_value
}

type baseValue struct {
	val C.emacs_value
	env *Environment
}

func (v baseValue) getVal() C.emacs_value {
	return v.val
}

type String interface {
	Value
	String() string
}

type stringValue struct {
	baseValue
}

func (s stringValue) String() string {
	res, _ := s.env.GoString(s)
	return res
}

type Symbol struct {
	baseValue
}

type Function struct {
	baseValue
}

type StdLib struct {
	env         *Environment
	messageFunc C.emacs_value
	fsetFunc    C.emacs_value
	Nil         Value
}

func newStdLib(e *Environment) *StdLib {
	messageStr := C.CString("message")
	defer C.free(unsafe.Pointer(messageStr))

	fsetStr := C.CString("fset")
	defer C.free(unsafe.Pointer(fsetStr))

	nilStr := C.CString("nil")
	defer C.free(unsafe.Pointer(nilStr))

	return &StdLib{
		env:         e,
		messageFunc: C.Intern(e.env, messageStr),
		fsetFunc:    C.Intern(e.env, fsetStr),
		Nil: baseValue{
			env: e,
			val: C.Intern(e.env, nilStr),
		},
	}
}

func (stdlib *StdLib) Message(s string) {
	str := stdlib.env.String(s).getVal()
	C.Funcall(stdlib.env.env, stdlib.messageFunc, 1, &str)
}

func (stdlib *StdLib) Funcall(f Function, args ...Value) Value {
	cargs := make([]C.emacs_value, len(args))
	for i := 0; i < len(args); i++ {
		cargs[i] = args[i].getVal()
	}
	return baseValue{
		env: stdlib.env,
		val: C.Funcall(stdlib.env.env, f.getVal(), C.int(len(args)), &cargs[0]),
	}
}

func (stdlib *StdLib) Intern(s string) Symbol {
	valStr := C.CString(s)
	defer C.free(unsafe.Pointer(valStr))

	return Symbol{
		baseValue{
			env: stdlib.env,
			val: C.Intern(stdlib.env.env, valStr),
		},
	}
}

func (stdlib *StdLib) Fset(sym Symbol, f Function) {
	args := [2]C.emacs_value{sym.val, f.val}
	C.Funcall(stdlib.env.env, stdlib.fsetFunc, 2, &args[0])
}
