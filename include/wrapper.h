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
