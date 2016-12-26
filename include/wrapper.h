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

#include <stdlib.h>
#include "emacs-module.h"

extern int emacsModuleInit (struct emacs_runtime *ert);

extern emacs_value emacsFunctionWrapper(emacs_env* env, ptrdiff_t nargs,
                                        emacs_value args[], void* data);

extern emacs_value emacsCallFunction (emacs_env* env, ptrdiff_t nargs,
                                        emacs_value *args, ptrdiff_t idx);

extern void emacsPointerWrapper(void* ptr);

extern void emacsFinalizeFunction (ptrdiff_t idx);

static inline emacs_env * GetEnvironment(struct emacs_runtime *ert) {
  return ert->get_environment(ert);
}

static inline emacs_value MakeGlobalRef(emacs_env *env, emacs_value ref) {
  return env->make_global_ref(env, ref);
}

static inline void FreeGlobalRef(emacs_env *env, emacs_value ref) {
  return env->free_global_ref(env, ref);
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

static inline ptrdiff_t StringSize(emacs_env *env, emacs_value value) {
  ptrdiff_t size;
  env->copy_string_contents(env, value, NULL, &size);
  return size;
}

static inline char* CopyString(emacs_env *env, emacs_value value, ptrdiff_t size) {
  char* buffer = malloc(size * sizeof(char));
  env->copy_string_contents(env, value, buffer, &size);
  return buffer;
}

static inline bool IsNotNil(emacs_env *env, emacs_value value) {
  return env->is_not_nil(env, value);
}

static inline intmax_t ExtractInteger(emacs_env *env, emacs_value value) {
  return env->extract_integer(env, value);
}

static inline emacs_value MakeInteger(emacs_env *env, intmax_t value) {
  return env->make_integer(env, value);
}

static inline double ExtractFloat(emacs_env *env, emacs_value value) {
  return env->extract_float(env, value);
}

static inline emacs_value MakeFloat(emacs_env *env, double value) {
  return env->make_float(env, value);
}

static inline bool Eq(emacs_env *env, emacs_value a, emacs_value b) {
  return env->eq(env, a, b);
}

static inline emacs_value MakeFunction(emacs_env *env, int min_arity, int max_arity,
                                       const char* documentation, ptrdiff_t idx) {
  return env->make_function(env, min_arity, max_arity,
                            &emacsFunctionWrapper,
                            documentation, (void*)idx);
}

static inline enum emacs_funcall_exit NonLocalExitCheck(emacs_env *env) {
  return env->non_local_exit_check(env);
}

static inline void NonLocalExitClear(emacs_env *env) {
  env->non_local_exit_clear(env);
}

static inline void NonLocalExitSignal(emacs_env *env, emacs_value symbol, emacs_value data) {
  env->non_local_exit_signal(env, symbol, data);
}

static inline void NonLocalExitThrow(emacs_env *env, emacs_value symbol, emacs_value data) {
  env->non_local_exit_throw(env, symbol, data);
}

static inline emacs_value MakeUserPointer(emacs_env *env, ptrdiff_t idx) {
  return env->make_user_ptr(env, &emacsPointerWrapper, (void*)idx);
}

static inline ptrdiff_t GetUserPointer(emacs_env *env, emacs_value uptr) {
  return (ptrdiff_t)env->get_user_ptr(env, uptr);
}

static inline void SetUserPointer(emacs_env *env, emacs_value uptr, ptrdiff_t idx) {
  env->set_user_ptr(env, uptr, (void*)idx);
}

static inline emacs_value VecGet(emacs_env *env, emacs_value vec, ptrdiff_t idx) {
  return env->vec_get(env, vec, idx);
}

static inline void VecSet(emacs_env *env, emacs_value vec, ptrdiff_t idx, emacs_value val) {
  env->vec_set(env, vec, idx, val);
}

static inline ptrdiff_t VecSize(emacs_env *env, emacs_value vec) {
  return env->vec_size(env, vec);
}

#endif /* GOEMACS_WRAPPER_H */
