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
import (
	"fmt"
	"unsafe"
)

type StdLib struct {
	env         *Environment
	messageFunc C.emacs_value
	fsetFunc    C.emacs_value
	fboundpFunc C.emacs_value
	Nil         Value
	T           Value
}

func newStdLib(e *Environment) *StdLib {
	n := baseValue{
		env: e,
		val: e.intern("nil"),
	}

	t := baseValue{
		env: e,
		val: e.intern("t"),
	}

	return &StdLib{
		env: e,

		messageFunc: e.intern("message"),
		fsetFunc:    e.intern("fset"),
		fboundpFunc: e.intern("fboundp"),

		Nil: n,
		T:   t,
	}
}

func (stdlib *StdLib) Message(s string) {
	str := stdlib.env.String(s).getVal()
	C.Funcall(stdlib.env.env, stdlib.messageFunc, 1, &str)
}

func (stdlib *StdLib) Funcall(f Callable, args ...Value) (Value, error) {
	if !f.Callable() {
		return stdlib.Nil, fmt.Errorf("symbol is not a function")
	}
	cargs := make([]C.emacs_value, len(args))
	for i := 0; i < len(args); i++ {
		cargs[i] = args[i].getVal()
	}
	return baseValue{
		env: stdlib.env,
		val: C.Funcall(stdlib.env.env, f.getVal(), C.int(len(args)), &cargs[0]),
	}, nil
}

func (stdlib *StdLib) Intern(s string) Symbol {
	valStr := C.CString(s)
	defer C.free(unsafe.Pointer(valStr))

	return symbolValue{
		baseValue: baseValue{
			env: stdlib.env,
			val: C.Intern(stdlib.env.env, valStr),
		},
	}
}

func (stdlib *StdLib) Fset(sym Symbol, f Function) {
	args := [2]C.emacs_value{sym.getVal(), f.getVal()}
	C.Funcall(stdlib.env.env, stdlib.fsetFunc, 2, &args[0])
}

func (stdlib *StdLib) Fboundp(sym Symbol) bool {
	symbol := sym.getVal()
	val := baseValue{
		env: stdlib.env,
		val: C.Funcall(stdlib.env.env, stdlib.fboundpFunc, 1, &symbol),
	}
	return stdlib.env.GoBool(val)
}
