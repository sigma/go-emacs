/* environment.go - Go wrapper for Emacs module API.

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
	"errors"
	"unsafe"
)

type Environment struct {
	// FIXME: for some reason, struct_emacs_env doesn't compile
	env    *C.struct_emacs_env_25
	stdlib *StdLib
}

func (e *Environment) StdLib() *StdLib {
	if e.stdlib == nil {
		e.stdlib = newStdLib(e)
	}
	return e.stdlib
}

func (e *Environment) MakeGlobalRef(ref Value) Value {
	return Value{
		C.MakeGlobalRef(e.env, ref.val),
	}
}

func (e *Environment) FreeGlobalRef(ref Value) {
	C.FreeGlobalRef(e.env, ref.val)
}

func (e *Environment) NonLocalExitCheck() error {
	code := C.NonLocalExitCheck(e.env)
	if code == C.emacs_funcall_exit_return {
		return nil
	}
	switch code {
	case C.emacs_funcall_exit_signal:
		return errors.New("signal")
	case C.emacs_funcall_exit_throw:
		return errors.New("throw")
	default:
		return nil
	}
}

func (e *Environment) GoString(v Value) (string, error) {
	size := C.StringSize(e.env, v.val)
	buffer := C.CopyString(e.env, v.val, size)
	defer C.free(unsafe.Pointer(buffer))

	if err := e.NonLocalExitCheck(); err != nil {
		return "", err
	}

	s := C.GoStringN(buffer, C.int(size))
	return s, nil
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
