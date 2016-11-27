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

	log "github.com/Sirupsen/logrus"
	emacs "github.com/sigma/goemacs"
)

func init() {
	emacs.Register(initModule)
}

func initModule(env *emacs.Environment) {
	log.Info("module initialization started")

	stdlib := env.StdLib()
	stdlib.Message("hello from go init")

	log.Info("creating native function")
	helloFunc := env.MakeFunction(Hello, 1, "hello", nil)

	log.Info("creating symbol")
	helloSym := stdlib.Intern("hello")

	log.Info("calling function")
	stdlib.Funcall(helloFunc, env.String("function"))

	log.Info("calling symbol before it's bound")
	_, err := stdlib.Funcall(helloSym, env.String("symbol"))
	if err != nil {
		fmt.Println(err)
	}

	log.Info("binding symbol to function")
	stdlib.Fset(helloSym, helloFunc)

	log.Info("calling symbol after it's bound")
	stdlib.Funcall(helloSym, env.String("symbol"))

	log.Info("module initialization complete")
}

func Hello(ctx *emacs.FunctionCallContext) emacs.Value {
	stdlib := ctx.StdLib()

	// we're guaranteed to be called with 1 argument
	s, err := ctx.StringArg(0)
	if err != nil {
		return stdlib.Nil
	}

	messages := make(chan string)
	go func() { messages <- s }()

	stdlib.Message(fmt.Sprintf("Hello %s!", <-messages))
	return stdlib.T
}

func main() {}
