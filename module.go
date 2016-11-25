/* module.go - Go wrapper for Emacs module API.

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

// #include "include/emacs-module.h"
// static inline emacs_env * GetEnvironment(struct emacs_runtime *ert) { return ert->get_environment(ert); }
// static inline emacs_value Intern(emacs_env *env, const char* name) { return env->intern(env, name); }
// static inline emacs_value Funcall(emacs_env *env, emacs_value func, int nargs, emacs_value args[]) { return env->funcall(env, func, nargs, args); }
// static inline emacs_value MakeString(emacs_env *env, const char* contents, int length) { return env->make_string(env, contents, length); }
import "C"

var initFuncs = make([]func(*Environment), 0)

type Environment struct {
	// FIXME: for some reason, struct_emacs_env doesn't compile
	env *C.struct_emacs_env_25
}

func Register(f func(*Environment)) {
	initFuncs = append(initFuncs, f)
}

//export emacs_module_init
func emacs_module_init(e *C.struct_emacs_runtime) C.int {
	env := Environment{
		env: C.GetEnvironment(e),
	}

	for _, f := range initFuncs {
		f(&env)
	}
	return 0
}

type StdLib struct {
	env         *Environment
	messageFunc C.emacs_value
}

func (e *Environment) StdLib() *StdLib {
	return &StdLib{
		env:         e,
		messageFunc: C.Intern(e.env, C.CString("message")),
	}
}

func (stdlib *StdLib) Message(s string) {
	str := C.MakeString(stdlib.env.env, C.CString(s), C.int(len(s)))
	C.Funcall(stdlib.env.env, stdlib.messageFunc, 1, &str)
}
