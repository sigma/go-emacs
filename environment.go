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

package emacs

/*
#include "include/wrapper.h"
*/
import "C"
import (
	"errors"
	"unsafe"
)

// Environment provides primitives for emacs modules
type Environment interface {
	MakeGlobalRef(Value) Value
	FreeGlobalRef(Value)
	NonLocalExitCheck() error
	GoString(Value) (string, error)
	String(string) String
	GoBool(Value) bool
	Bool(bool) Value
	GoInt(Value) int64
	Int(int64) Int
	GoFloat(Value) float64
	Float(float64) Float
	MakeFunction(FunctionType, int, string, interface{}) Function
	MakeUserPointer(interface{}) UserPointer
	ResolveUserPointer(UserPointer) (interface{}, bool)
	VecSize(Vector) int
	VecSet(Vector, int, Value)
	VecGet(Vector, int) Value

	// additional helpers
	StdLib() StdLib
	RegisterFunction(string, FunctionType, int, string, interface{}) Function
	ProvideFeature(string)
}

type emacsEnv struct {
	env    *C.emacs_env
	stdlib *emacsLib
}

func (e *emacsEnv) intern(s string) C.emacs_value {
	str := C.CString(s)
	defer C.free(unsafe.Pointer(str))
	return C.Intern(e.env, str)
}

func (e *emacsEnv) eq(a, b C.emacs_value) bool {
	return bool(C.Eq(e.env, a, b))
}

func newStdLib(e *emacsEnv) *emacsLib {
	n := baseValue{
		env: e,
		val: e.intern("nil"),
	}

	t := baseValue{
		env: e,
		val: e.intern("t"),
	}

	return &emacsLib{
		env: e,

		nilValue: n,
		tValue:   t,
	}
}

func (e *emacsEnv) StdLib() StdLib {
	if e.stdlib == nil {
		e.stdlib = newStdLib(e)
	}
	return e.stdlib
}

func (e *emacsEnv) MakeGlobalRef(ref Value) Value {
	return baseValue{
		env: e,
		val: C.MakeGlobalRef(e.env, ref.getVal()),
	}
}

func (e *emacsEnv) FreeGlobalRef(ref Value) {
	C.FreeGlobalRef(e.env, ref.getVal())
}

func (e *emacsEnv) NonLocalExitCheck() error {
	code := C.NonLocalExitCheck(e.env)
	defer C.NonLocalExitClear(e.env)
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

func (e *emacsEnv) GoString(v Value) (string, error) {
	size := C.StringSize(e.env, v.getVal())
	buffer := C.CopyString(e.env, v.getVal(), size)
	defer C.free(unsafe.Pointer(buffer))

	if err := e.NonLocalExitCheck(); err != nil {
		return "", err
	}

	s := C.GoString(buffer)
	return s, nil
}

func (e *emacsEnv) stringVal(s string) C.emacs_value {
	valStr := C.CString(s)
	defer C.free(unsafe.Pointer(valStr))

	return C.MakeString(e.env, valStr, C.int(len(s)))
}

func (e *emacsEnv) String(s string) String {
	return stringValue{
		baseValue{
			env: e,
			val: e.stringVal(s),
		},
	}
}

func (e *emacsEnv) GoBool(v Value) bool {
	return bool(C.IsNotNil(e.env, v.getVal()))
}

func (e *emacsEnv) Bool(b bool) Value {
	stdlib := e.StdLib()
	if b {
		return stdlib.T()
	}
	return stdlib.Nil()
}

func (e *emacsEnv) GoInt(v Value) int64 {
	return int64(C.ExtractInteger(e.env, v.getVal()))
}

func (e *emacsEnv) Int(i int64) Int {
	return intValue{
		baseValue{
			env: e,
			val: C.MakeInteger(e.env, C.intmax_t(i)),
		},
	}
}

func (e *emacsEnv) GoFloat(v Value) float64 {
	return float64(C.ExtractFloat(e.env, v.getVal()))
}

func (e *emacsEnv) Float(i float64) Float {
	return floatValue{
		baseValue{
			env: e,
			val: C.MakeFloat(e.env, C.double(i)),
		},
	}
}

func (e *emacsEnv) MakeFunction(f FunctionType, arity int, doc string, data interface{}) Function {
	cArity := C.int(arity)
	idx := funcReg.Register(&functionEntry{
		f:     f,
		arity: arity,
		doc:   doc,
		data:  data,
	})

	docStr := C.CString(doc)
	defer C.free(unsafe.Pointer(docStr))

	return functionValue{
		baseValue{
			env: e,
			val: C.MakeFunction(e.env, cArity, cArity,
				docStr, C.ptrdiff_t(idx)),
		},
	}
}

func (e *emacsEnv) MakeUserPointer(obj interface{}) UserPointer {
	val := &simplePointerEntry{
		obj: obj,
	}
	idx := ptrReg.Register(val)

	return userPointerValue{
		baseValue{
			env: e,
			val: C.MakeUserPointer(e.env, C.ptrdiff_t(idx)),
		},
	}
}

func (e *emacsEnv) ResolveUserPointer(ptr UserPointer) (interface{}, bool) {
	val := ptr.getVal()
	p := ptrReg.Lookup(int64(C.GetUserPointer(e.env, val)))
	if p != nil {
		return p.underlyingObject(), true
	}
	return nil, false
}

func (e *emacsEnv) VecSize(vec Vector) int {
	return int(C.VecSize(e.env, vec.getVal()))
}

func (e *emacsEnv) VecSet(vec Vector, idx int, val Value) {
	C.VecSet(e.env, vec.getVal(), C.ptrdiff_t(idx), val.getVal())
}

func (e *emacsEnv) VecGet(vec Vector, idx int) Value {
	return baseValue{
		val: C.VecGet(e.env, vec.getVal(), C.ptrdiff_t(idx)),
	}
}

func (e *emacsEnv) RegisterFunction(name string, f FunctionType, arity int, doc string, data interface{}) Function {
	stdlib := e.StdLib()
	function := e.MakeFunction(f, arity, doc, data)
	sym := stdlib.Intern(name)
	stdlib.Fset(sym, function)
	return function
}

func (e *emacsEnv) ProvideFeature(name string) {
	stdlib := e.StdLib()
	sym := stdlib.Intern(name)
	stdlib.Provide(sym)
}

var _ Environment = (*emacsEnv)(nil)
