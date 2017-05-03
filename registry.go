/* registry.go - Go wrapper for Emacs module API.

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

import "sync"

// Registry is a key-value store for arbitrary objects
type Registry interface {
	Register(interface{}) int64
	Lookup(int64) (interface{}, bool)
	Unregister(int64)
}

type registry struct {
	mu      sync.Mutex
	index   int64
	objects map[int64]interface{}
}

// NewRegistry creates a new Registry instance
func NewRegistry() Registry {
	return &registry{
		objects: make(map[int64]interface{}),
	}
}

func (reg *registry) Register(object interface{}) int64 {
	reg.mu.Lock()
	defer reg.mu.Unlock()
	reg.index++
	for reg.objects[reg.index] != nil {
		reg.index++
	}
	reg.objects[reg.index] = object
	return reg.index
}

func (reg *registry) Lookup(i int64) (interface{}, bool) {
	reg.mu.Lock()
	defer reg.mu.Unlock()
	obj, ok := reg.objects[i]
	return obj, ok
}

func (reg *registry) Unregister(i int64) {
	reg.mu.Lock()
	defer reg.mu.Unlock()
	delete(reg.objects, i)
}
