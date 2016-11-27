/* main.go - Example for goemacs API

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

package main

// int plugin_is_GPL_compatible;
import "C"

import (
	"fmt"

	emacs "github.com/sigma/goemacs"
)

func init() {
	emacs.Register(initModule)
}

func initModule(env *emacs.Environment) {
	stdlib := env.StdLib()
	stdlib.Message("hello from go init")

	helloFunc := env.MakeFunction(Hello, 1, "hello")
	helloSym := stdlib.Intern("hello")
	stdlib.Fset(helloSym, helloFunc)

	stdlib.Funcall(helloFunc, env.String("function"))
	stdlib.Funcall(helloSym, env.String("symbol"))
}

func Hello(env *emacs.Environment, nargs int,
	args []emacs.Value, _ interface{}) emacs.Value {
	stdlib := env.StdLib()

	// we're guaranteed to be called with 1 argument
	s, err := env.GoString(args[0])
	if err != nil {
		return stdlib.Nil
	}

	stdlib.Message(fmt.Sprintf("Hello %s!", s))
	return stdlib.T
}

func main() {}
