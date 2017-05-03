/* types.go - Go wrapper for Emacs module API.

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

// Value wraps an emacs value
type Value interface {
	getVal() C.emacs_value
	AsString() String
	AsInt() Int
	AsFloat() Float
	AsSymbol() Symbol
	AsFunction() Function
	AsVector() Vector
	AsList() List
	AsUserPointer() UserPointer
}

type baseValue struct {
	val C.emacs_value
	env Environment
}

func (v baseValue) getVal() C.emacs_value {
	return v.val
}

func (v baseValue) AsString() String           { return stringValue{v} }
func (v baseValue) AsInt() Int                 { return intValue{v} }
func (v baseValue) AsFloat() Float             { return floatValue{v} }
func (v baseValue) AsSymbol() Symbol           { return symbolValue{v} }
func (v baseValue) AsFunction() Function       { return functionValue{v} }
func (v baseValue) AsVector() Vector           { return vectorValue{v} }
func (v baseValue) AsList() List               { return listValue{v} }
func (v baseValue) AsUserPointer() UserPointer { return userPointerValue{v} }

// String wraps an emacs string
type String interface {
	Value
	String() string
}

type stringValue struct {
	baseValue
}

func (s stringValue) String() string {
	res, _ := s.env.GoString(s)
	return res
}

// Int wraps an emacs integer
type Int interface {
	Value
}

type intValue struct {
	baseValue
}

// Float wraps an emacs float number
type Float interface {
	Value
}

type floatValue struct {
	baseValue
}

// Callable wraps a callable object: function or symbol
type Callable interface {
	Value
	Callable() bool
}

// Symbol wraps an emacs symbol
type Symbol interface {
	Callable
}

type symbolValue struct {
	baseValue
}

func (s symbolValue) Callable() bool {
	return s.env.StdLib().Fboundp(s)
}

// Function wraps a function object
type Function interface {
	Callable
}

type functionValue struct {
	baseValue
}

func (f functionValue) Callable() bool {
	return true
}

// UserPointer represents a module-created pointer
type UserPointer interface {
	Value
}

type userPointerValue struct {
	baseValue
}

// Vector wraps an emacs vector
type Vector interface {
	Value
	Size() int
	Get(int) Value
	Set(int, Value)
}

type vectorValue struct {
	baseValue
}

func (v vectorValue) Size() int {
	return v.env.VecSize(v)
}

func (v vectorValue) Set(idx int, item Value) {
	v.env.VecSet(v, idx, item)
}

func (v vectorValue) Get(idx int) Value {
	return v.env.VecGet(v, idx)
}

// List wraps an emacs list
type List interface {
	Value
}

type listValue struct {
	baseValue
}
