/* errors.go - Go wrapper for Emacs module API.

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

type nonLocalExit interface {
	error
	Symbol() Symbol
	Value() Value
}

type signal interface {
	nonLocalExit
	isSignal() bool
}

type throw interface {
	nonLocalExit
	isThrow() bool
}

type nonLocalExitImpl struct {
	symbol Symbol
	value  Value
}

func (e *nonLocalExitImpl) Error() string {
	return "emacs non-local exit"
}

func (e *nonLocalExitImpl) Symbol() Symbol {
	return e.symbol
}

func (e *nonLocalExitImpl) Value() Value {
	return e.value
}

var _ nonLocalExit = (*nonLocalExitImpl)(nil)

type signalImpl struct {
	nonLocalExitImpl
}

func (s *signalImpl) isSignal() bool { return true }

var _ signal = (*signalImpl)(nil)

type throwImpl struct {
	nonLocalExitImpl
}

func (s *throwImpl) isThrow() bool { return true }

var _ throw = (*throwImpl)(nil)
