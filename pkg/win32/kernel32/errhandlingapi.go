package kernel32

/*
	#include <Errhandlingapi.h>
*/
import "C"

import (
	"strconv"
)

type ErrorCode uint32

func (e ErrorCode) Error() string {
	return "win32 error code: " + strconv.Itoa(int(e))
}

/*
	GetLastError retrieves the calling thread's last-error code value.
	The last-error code is maintained on a per-thread basis.
	Multiple threads do not overwrite each other's last-error code.

	For more infomtation, see: https://docs.microsoft.com/en-us/windows/win32/api/errhandlingapi/nf-errhandlingapi-getlasterror
*/
func GetLastError() error {
	err := C.GetLastError()
	if err == 0 {
		return nil
	}

	return ErrorCode(err)
}