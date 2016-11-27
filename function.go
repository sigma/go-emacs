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

import "sync"

type FunctionCallContext struct {
	env  *Environment
	args []Value
	data interface{}
}

func (ctx *FunctionCallContext) Environment() *Environment {
	return ctx.env
}

func (ctx *FunctionCallContext) StdLib() *StdLib {
	return ctx.env.StdLib()
}

func (ctx *FunctionCallContext) StringArg(idx int) (string, error) {
	return ctx.env.GoString(ctx.args[idx])
}

type FunctionType func(*FunctionCallContext) Value

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
