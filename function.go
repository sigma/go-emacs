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

type FunctionCallContext struct {
	env  Environment
	args []Value
	data interface{}
}

func (ctx *FunctionCallContext) Environment() Environment {
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

type functionRegistry struct {
	reg Registry
}

func (fr *functionRegistry) Register(fn *FunctionEntry) int64 {
	return fr.reg.Register(fn)
}

func (fr *functionRegistry) Lookup(idx int64) *FunctionEntry {
	obj := fr.reg.Lookup(idx)
	return obj.(*FunctionEntry)
}

var funcReg = functionRegistry{
	reg: NewRegistry(),
}
