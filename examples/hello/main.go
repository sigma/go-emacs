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

	"github.com/sigma/goemacs"
)

func init() {
	goemacs.Register(initModule)
}

func initModule(env *goemacs.Environment) {
	stdlib := env.StdLib()
	stdlib.Message("hello from go")

	stdlib.Fset(stdlib.Intern("hello"), env.MakeFunction(Hello, 1, "hello"))
}

func Hello(env *goemacs.Environment, nargs int,
	args []goemacs.Value, _ interface{}) goemacs.Value {
	stdlib := env.StdLib()
	if nargs != 1 {
		// TODO: display error message
		return stdlib.Nil
	}
	s := env.GoString(args[0])
	stdlib.Message(fmt.Sprintf("Hello %s!", s))

	return stdlib.Nil
}

func main() {}
