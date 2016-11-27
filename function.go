/* function.go - Go wrapper for Emacs module API.

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

//export emacs_call_function
func emacs_call_function(
	//FIXME: emacs_env_25 shouldn't be used
	env *C.struct_emacs_env_25, nargs C.ptrdiff_t,
	args *C.emacs_value, idx C.ptrdiff_t) C.emacs_value {

	n := int(nargs)
	pargs := (*[1 << 30]C.emacs_value)(unsafe.Pointer(args))
	arguments := make([]Value, n)
	for i := 0; i < n; i++ {
		arguments[i] = baseValue{
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
	).getVal()
}
