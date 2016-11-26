/* wrapper.c - Go wrapper for Emacs module API.

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

#include "include/wrapper.h"

emacs_value emacs_function_wrapper(emacs_env* env, ptrdiff_t nargs,
                                   emacs_value args[], void* data) {
  int idx = (ptrdiff_t)data;
  return emacs_call_function(env, nargs, args, (ptrdiff_t)data);
}
