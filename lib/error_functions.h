/*************************************************************************\
*                  Copyright (C) Michael Kerrisk, 2023.                   *
*                                                                         *
* This program is free software. You may use, modify, and redistribute it *
* under the terms of the GNU Lesser General Public License as published   *
* by the Free Software Foundation, either version 3 or (at your option)   *
* any later version. This program is distributed without any warranty.    *
* See the files COPYING.lgpl-v3 and COPYING.gpl-v3 for details.           *
\*************************************************************************/

/* Listing 3-2 */

/* error_functions.h

   Header file for error_functions.c.
*/
#ifndef ERROR_FUNCTIONS_H
#define ERROR_FUNCTIONS_H

/* Error diagnostic routines */

/* Display error message including 'errno' diagnostic, and
   return to caller */
void errMsg(const char *format, ...);

#ifdef __GNUC__

    /* This macro stops 'gcc -Wall' complaining that "control reaches
       end of non-void function" if we use the following functions to
       terminate main() or some other non-void function. */

#define NORETURN __attribute__ ((__noreturn__))
#else
#define NORETURN
#endif

/* Display error message including 'errno' diagnostic, and
   terminate the process */
void errExit(const char *format, ...) NORETURN ;

/* Display error message including 'errno' diagnostic, and
   terminate the process by calling _exit().

   The relationship between this function and errExit() is analogous
   to that between _exit(2) and exit(3): unlike errExit(), this
   function does not flush stdout and calls _exit(2) to terminate the
   process (rather than exit(3), which would cause exit handlers to be
   invoked).

   These differences make this function especially useful in a library
   function that creates a child process that must then terminate
   because of an error: the child must terminate without flushing
   stdio buffers that were partially filled by the caller and without
   invoking exit handlers that were established by the caller. */
void err_exit(const char *format, ...) NORETURN ;

/* The following function does the same as errExit(), but expects
   the error number in 'errnum' */
void errExitEN(int errnum, const char *format, ...) NORETURN ;

/* Print an error message (without an 'errno' diagnostic) */
void fatal(const char *format, ...) NORETURN ;

/* Print a command usage error message and terminate the process */
void usageErr(const char *format, ...) NORETURN ;

/* Diagnose an error in command-line arguments and
   terminate the process */
void cmdLineErr(const char *format, ...) NORETURN ;

#endif
