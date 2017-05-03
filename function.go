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

package emacs

// FunctionCallContext is the one argument module functions will receive
type FunctionCallContext interface {
	Environment() Environment
	NumberArgs() int
	Arg(int) Value
	StringArg(int) String
	UserPointerArg(int) UserPointer
	GoStringArg(int) (string, error)
}

type emacsCallContext struct {
	env  Environment
	args []Value
	data interface{}
}

func (ctx *emacsCallContext) NumberArgs() int {
	return len(ctx.args)
}

func (ctx *emacsCallContext) Environment() Environment {
	return ctx.env
}

func (ctx *emacsCallContext) Arg(idx int) Value {
	return ctx.args[idx]
}

func (ctx *emacsCallContext) StringArg(idx int) String {
	return ctx.args[idx].AsString()
}

func (ctx *emacsCallContext) UserPointerArg(idx int) UserPointer {
	return ctx.args[idx].AsUserPointer()
}

func (ctx *emacsCallContext) GoStringArg(idx int) (string, error) {
	return ctx.env.GoString(ctx.args[idx])
}

// FunctionType is the type for module functions
type FunctionType func(FunctionCallContext) (Value, error)

type functionEntry struct {
	f     FunctionType
	arity int
	doc   string
	data  interface{}
}

type functionRegistry struct {
	reg Registry
}

func (fr *functionRegistry) Register(fn *functionEntry) int64 {
	return fr.reg.Register(fn)
}

func (fr *functionRegistry) Lookup(idx int64) *functionEntry {
	obj, ok := fr.reg.Lookup(idx)
	if ok {
		return obj.(*functionEntry)
	}
	return nil
}

var funcReg = functionRegistry{
	reg: NewRegistry(),
}
