/* pointers.go - Go wrapper for Emacs module API.

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

type pointerEntry interface {
	underlyingObject() interface{}
	Finalize()
}

type simplePointerEntry struct {
	obj interface{}
}

func (ptr *simplePointerEntry) underlyingObject() interface{} {
	return ptr.obj
}

func (ptr *simplePointerEntry) Finalize() {}

type pointerRegistry struct {
	reg Registry
}

func (pr *pointerRegistry) Register(entry pointerEntry) int64 {
	idx := pr.reg.Register(entry)
	return idx
}

func (pr *pointerRegistry) Lookup(idx int64) pointerEntry {
	obj, ok := pr.reg.Lookup(idx)
	if ok {
		return obj.(pointerEntry)
	}
	return nil
}

func (pr *pointerRegistry) Unregister(idx int64) {
	pr.reg.Unregister(idx)
}

var ptrReg = pointerRegistry{
	reg: NewRegistry(),
}
