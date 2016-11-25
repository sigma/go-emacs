#include "include/wrapper.h"

emacs_value emacs_function_wrapper(emacs_env* env, ptrdiff_t nargs,
                                   emacs_value args[], void* data) {
 /* FIXME: obviously not correct... */
  return *args;
}
