package kernel32

/*
	#include <handleapi.h>
*/
import "C"

import (
	"unsafe"

	"github.com/warrenulrich/win32-go/pkg/win32"
)

func CloseHandle(handle win32.Handle) error {
	if C.CloseHandle(C.HANDLE(unsafe.Pointer(handle))) == 0 {
		return GetLastError()
	}
	return nil
}
