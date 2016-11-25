/* wrapper.h - Thin wrapper to expose a more convenient Emacs inferface.

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

#ifndef GOEMACS_WRAPPER_H
#define GOEMACS_WRAPPER_H

#include "emacs-module.h"

static inline emacs_env * GetEnvironment(struct emacs_runtime *ert) {
  return ert->get_environment(ert);
}

static inline emacs_value Intern(emacs_env *env, const char* name) {
  return env->intern(env, name);
}

static inline emacs_value Funcall(emacs_env *env, emacs_value func, int nargs, emacs_value args[]) {
  return env->funcall(env, func, nargs, args);
}

static inline emacs_value MakeString(emacs_env *env, const char* contents, int length) {
  return env->make_string(env, contents, length);
}

extern emacs_value emacs_function_wrapper(emacs_env* env, ptrdiff_t nargs,
                                          emacs_value args[], void* data);

static inline emacs_value MakeFunction(emacs_env *env, int min_arity, int max_arity,
                                       const char* documentation, void *data) {
  return env->make_function(env, min_arity, max_arity,
                            &emacs_function_wrapper,
                            documentation, data);
}

#endif /* GOEMACS_WRAPPER_H */
