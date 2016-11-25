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
}

func (e *Environment) StdLib() *StdLib {
	if e.stdlib == nil {
		e.stdlib = &StdLib{
			env:         e,
			messageFunc: C.Intern(e.env, C.CString("message")),
			fsetFunc:    C.Intern(e.env, C.CString("fset")),
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

type FunctionType func(*Environment, int, []Value, interface{})

func (e *Environment) MakeFunction(f FunctionType) Function {
	return Function{}
}
