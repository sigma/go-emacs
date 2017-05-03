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

package emacs

/*
#include "include/wrapper.h"
*/
import "C"
import "fmt"

// StdLib exposes high-level emacs functions
type StdLib interface {
	Funcall(f Callable, args ...Value) (Value, error)
	Eq(a, b Value) bool
	Equal(a, b Value) bool
	Intern(s string) Symbol
	Fset(sym Symbol, f Function)
	Fboundp(sym Symbol) bool
	Provide(sym Symbol)
	Message(s string)
	List(items ...Value) List
	Nil() Value
	T() Value
}

type emacsLib struct {
	env *emacsEnv

	nilValue Value
	tValue   Value
}

func (stdlib *emacsLib) Nil() Value {
	return stdlib.nilValue
}

func (stdlib *emacsLib) T() Value {
	return stdlib.tValue
}

func (stdlib *emacsLib) Funcall(f Callable, args ...Value) (Value, error) {
	cargs := make([]C.emacs_value, len(args))
	for i := 0; i < len(args); i++ {
		cargs[i] = args[i].getVal()
	}
	ret := C.Funcall(stdlib.env.env, f.getVal(), C.int(len(args)), &cargs[0])

	if stdlib.env.NonLocalExitCheck() != nil {
		return stdlib.Nil(), fmt.Errorf("symbol is not a function")
	}
	return baseValue{
		env: stdlib.env,
		val: ret,
	}, nil
}

func (stdlib *emacsLib) Eq(a, b Value) bool {
	return stdlib.env.eq(a.getVal(), b.getVal())
}

func (stdlib *emacsLib) Equal(a, b Value) bool {
	equal := stdlib.Intern("equal")
	res, _ := stdlib.Funcall(equal, a, b)
	return !stdlib.Eq(res, stdlib.Nil())
}

func (stdlib *emacsLib) Intern(s string) Symbol {
	return symbolValue{
		baseValue: baseValue{
			env: stdlib.env,
			val: stdlib.env.intern(s),
		},
	}
}

func (stdlib *emacsLib) Fset(sym Symbol, f Function) {
	fset := stdlib.Intern("fset")
	stdlib.Funcall(fset, sym, f)
}

func (stdlib *emacsLib) Fboundp(sym Symbol) bool {
	fboundp := stdlib.Intern("fboundp")
	res, _ := stdlib.Funcall(fboundp, sym)
	return stdlib.env.GoBool(res)
}

func (stdlib *emacsLib) Provide(sym Symbol) {
	provide := stdlib.Intern("provide")
	stdlib.Funcall(provide, sym)
}

func (stdlib *emacsLib) Message(s string) {
	message := stdlib.Intern("message")
	stdlib.Funcall(message, stdlib.env.String(s))
}

func (stdlib *emacsLib) List(items ...Value) List {
	if len(items) == 0 {
		return stdlib.Nil()
	}
	list := stdlib.Intern("list")
	res, _ := stdlib.Funcall(list, items...)
	return res.AsList()
}

var _ StdLib = (*emacsLib)(nil)
