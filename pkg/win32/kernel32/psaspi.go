package kernel32

/*
	#cgo CFLAGS: -DPSAPI_VERSION=1
 	#cgo LDFLAGS: -lpsapi
	#include <windows.h>
	#include <psapi.h>
*/
import "C"

func EnumProcesses() ([]uint32, error) {
	var idBuffer []uint32 = make([]uint32, 1024)

	var bytesNeeded uint32

	if C.EnumProcesses((*C.DWORD)(&idBuffer[0]), C.DWORD(len(idBuffer)), (*C.ulong)(&bytesNeeded)) == 0 {
		return nil, GetLastError()
	}

	if bytesNeeded > 0 {
		idBuffer = make([]uint32, bytesNeeded/4)
		if C.EnumProcesses((*C.DWORD)(&idBuffer[0]), C.DWORD(len(idBuffer)), (*C.ulong)(&bytesNeeded)) == 0 {
			return nil, GetLastError()
		}
	}

	return idBuffer[:bytesNeeded/4], nil
}
