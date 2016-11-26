#include "include/wrapper.h"

emacs_value emacs_function_wrapper(emacs_env* env, ptrdiff_t nargs,
                                   emacs_value args[], void* data) {
  int idx = (ptrdiff_t)data;
  return emacs_call_function(env, nargs, args, (ptrdiff_t)data);
}
