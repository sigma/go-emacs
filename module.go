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

package emacs

/*
#include "include/wrapper.h"
*/
import "C"
import "unsafe"

var initFuncs = make([]func(Environment), 0)

// InitFunction is the type for functions to be called at module initialization
type InitFunction func(Environment)

// Register makes sure a function will be called at module initialization
func Register(init InitFunction) {
	initFuncs = append(initFuncs, init)
}

//export emacsModuleInit
func emacsModuleInit(e *C.struct_emacs_runtime) C.int {
	env := &emacsEnv{
		env: C.GetEnvironment(e),
	}

	for _, f := range initFuncs {
		f(env)
	}
	return 0
}

//export emacsCallFunction
func emacsCallFunction(
	env *C.emacs_env, nargs C.ptrdiff_t,
	args *C.emacs_value, idx C.ptrdiff_t) C.emacs_value {

	e := &emacsEnv{
		env: env,
	}

	n := int(nargs)
	pargs := (*[1 << 30]C.emacs_value)(unsafe.Pointer(args))

	arguments := make([]Value, n)
	for i := 0; i < n; i++ {
		arguments[i] = baseValue{
			env: e,
			val: pargs[i],
		}
	}
	entry := funcReg.Lookup(int64(idx))

	res, err := entry.f(
		&emacsCallContext{
			e,
			arguments,
			entry.data,
		},
	)

	if err == nil {
		if res == nil {
			return e.intern("nil")
		}
		return res.getVal()
	}

	if isSignal(err) {
		s := err.(signal)
		C.NonLocalExitSignal(env, s.Symbol().getVal(), s.Value().getVal())
	} else if isThrow(err) {
		t := err.(throw)
		C.NonLocalExitThrow(env, t.Symbol().getVal(), t.Value().getVal())
	} else {
		msg := err.Error()
		C.NonLocalExitThrow(env, e.intern("error"), e.stringVal(msg))
	}
	return e.intern("nil")
}

//export emacsFinalizeFunction
func emacsFinalizeFunction(idx C.ptrdiff_t) {
	index := int64(idx)
	defer ptrReg.Unregister(index)

	entry := ptrReg.Lookup(index)
	entry.Finalize()
}
