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

import (
	"fmt"
	"log"

	emacs "github.com/sigma/go-emacs"
	_ "github.com/sigma/go-emacs/gpl-compatible"
)

func init() {
	emacs.Register(initModule)
}

func initModule(env emacs.Environment) {
	log.Println("module initialization started")

	stdlib := env.StdLib()
	stdlib.Message("hello from go init")

	log.Println("creating native function")
	helloFunc := env.MakeFunction(Hello, 1, "hello", nil)

	log.Println("creating symbol")
	helloSym := stdlib.Intern("hello")

	log.Println("calling function")
	stdlib.Funcall(helloFunc, env.String("function"))

	log.Println("calling symbol before it's bound")
	_, err := stdlib.Funcall(helloSym, env.String("symbol"))
	if err != nil {
		fmt.Println(err)
	}

	log.Println("binding symbol to function")
	stdlib.Fset(helloSym, helloFunc)

	log.Println("calling symbol after it's bound")
	stdlib.Funcall(helloSym, env.String("symbol"))

	stdlib.Provide(helloSym)
	log.Println("module initialization complete")
}

// Hello is a sample function that calls "message"
func Hello(ctx emacs.FunctionCallContext) (emacs.Value, error) {
	stdlib := ctx.Environment().StdLib()

	// we're guaranteed to be called with 1 argument
	s, err := ctx.GoStringArg(0)
	if err != nil {
		return stdlib.Nil(), err
	}

	messages := make(chan string)
	go func() { messages <- s }()

	stdlib.Message(fmt.Sprintf("Hello %s!", <-messages))
	return stdlib.T(), nil
}

func main() {}
