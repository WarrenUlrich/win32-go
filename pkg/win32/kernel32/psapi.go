package kernel32

/*
	#cgo CFLAGS: -DPSAPI_VERSION=1
 	#cgo LDFLAGS: -lpsapi
	#include <windows.h>
	#include <psapi.h>
*/
import "C"

/*
	EnumProcesses retrieves the process identifier for each process object in the system.
*/
func EnumProcesses() ([]uint32, error) {
	var buffer []uint32 = make([]uint32, 1024)
	
	var fn func(buf []uint32) (uint32, error)

	fn = func(buf []uint32) (uint32, error) {
		var cb uint32 = uint32(len(buf) * 4)

		var cbNeeded uint32
		if C.EnumProcesses((*C.ulong)(&buf[0]), C.ulong(cb), (*C.ulong)(&cbNeeded)) == 0 {
			return cbNeeded / 4, GetLastError()
		}

		if cbNeeded == cb {
			return fn(make([]uint32, len(buf)*2))
		}

		return cbNeeded / 4, nil
	}

	needed, err := fn(buffer)
	if err != nil {
		return nil, err
	}

	return buffer[:needed], nil
}
