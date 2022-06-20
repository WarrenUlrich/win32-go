package kernel32

/*
	#include <windows.h>
	#include <tlhelp32.h>
*/
import "C"

import (
	"github.com/warrenulrich/win32-go/pkg/win32"
	"unsafe"
)

type ThFlags uint32

/*
	ProcessEntry32 describes an entry from a list of the processes
	residing in the system address space when a snapshot was taken.
*/
type ProcessEntry32 struct {

	/*
		Size is the size of the structure, in bytes.
		Before calling the Process32First function,
		set this member to unsafe.SizeOf(ProcessEntry32).
		If you do not initialize Size, Process32First fails.
	*/
	Size uint32

	/*
		Usage is no longer used and is always set to zero.
	*/
	Usage uint32

	/*
		ProcessID is the process identifier.
	*/
	ProcessID uint32

	/*
		DefaultHeapID is no longer used and is always set to zero.
	*/
	DefaultHeapID uintptr

	/*
		ModuleID is no longer used and is always set to zero.
	*/
	ModuleID uint32

	/*
		Threads is the number of execution threads started by the process.
	*/
	Threads uint32

	/*
		ParentProcessID is the identifier of the process that created this process (its parent process).
	*/
	ParentProcessID uint32
	/*
		PriorityClass is the base priority of any threads created by this process.
	*/
	PriorityClass uint32

	/*
		Flags is no longer used and is always set to zero.
	*/
	Flags uint32

	/*
		The name of the executable file for the process.
		To retrieve the full path to the executable file,
		call the Module32First function and check the
		szExePath member of the MODULEENTRY32 structure that is returned.
		However, if the calling process is a 32-bit process,
		you must call the QueryFullProcessImageName function to
		retrieve the full path of the executable file for a 64-bit process.
	*/
	ExeFile [260]byte
}

/*
	ExeFileString creates a string from pe32.ExeFile.
*/
func (pe32 ProcessEntry32) ExeFileString() string {
	var i int
	for i = 0; i < 260; i++ {
		if pe32.ExeFile[i] == 0 {
			break
		}
	}

	return string(pe32.ExeFile[:i])
}

const (
	/*
		TH32CS_INHERIT indicates that the snapshot handle is to be inheritable.
	*/
	TH32CS_INHERIT ThFlags = 0x80000000

	/*
		TH32CS_SNAPHEAPLIST includes all heaps of the process specified
		in th32ProcessID in the snapshot.
		To enumerate the heaps, see Heap32ListFirst.
	*/
	TH32CS_SNAPHEAPLIST ThFlags = 0x00000001

	/*
		TH32CS_SNAPMODULE includes all modules of the process specified
		in th32ProcessID in the snapshot. To enumerate the modules, see Module32First.
		If the function fails with ERROR_BAD_LENGTH, retry the function until it succeeds.

		64-bit Windows:  Using this flag in a 32-bit process includes the 32-bit modules
		of the process specified in th32ProcessID, while using it in a 64-bit process
		includes the 64-bit modules. To include the 32-bit modules of the process specified
		in th32ProcessID from a 64-bit process, use the TH32CS_SNAPMODULE32 flag.
	*/
	TH32CS_SNAPMODULE ThFlags = 0x00000008

	/*
		TH32CS_SNAPMODULE32 includes all 32-bit modules of the process specified
		in th32ProcessID in the snapshot when called from a 64-bit process.
		This flag can be combined with TH32CS_SNAPMODULE or TH32CS_SNAPALL.
		If the function fails with ERROR_BAD_LENGTH, retry the function until it succeeds.
	*/
	TH32CS_SNAPMODULE32 ThFlags = 0x00000010

	/*
		TH32CS_SNAPPROCESS includes all processes in the system in the snapshot.
		To enumerate the processes, see Process32First.
	*/
	TH32CS_SNAPPROCESS ThFlags = 0x00000002

	/*
		TH32CS_SNAPTHREAD includes all threads in the system in the snapshot.
		To enumerate the threads, see Thread32First.

		To identify the threads that belong to a specific process,
		compare its process identifier to the th32OwnerProcessID member
		of the THREADENTRY32 structure when enumerating the threads.
	*/
	TH32CS_SNAPTHREAD ThFlags = 0x00000004

	/*
		TH32CS_SNAPALL includes all processes and threads in the system,
		plus the heaps and modules of the process specified in th32ProcessID.
		Equivalent to specifying the TH32CS_SNAPHEAPLIST, TH32CS_SNAPMODULE,
		TH32CS_SNAPPROCESS, and TH32CS_SNAPTHREAD values combined using an OR operation ('|').
	*/
	TH32CS_SNAPALL ThFlags = TH32CS_SNAPHEAPLIST | TH32CS_SNAPMODULE | TH32CS_SNAPMODULE32 | TH32CS_SNAPPROCESS | TH32CS_SNAPTHREAD
)

/*
	CreateToolhelp32Snapshot takes a snapshot of the specified processes,
	as well as the heaps, modules, and threads used by these processes.

	The snapshot taken by this function is examined by the
	other tool help functions to provide their results.
	Access to the snapshot is read only.
	The snapshot handle acts as an object handle and is
	subject to the same rules regarding which processes
	and threads it is valid in.

	To enumerate the heap or module states for all processes,
	specify TH32CS_SNAPALL and set th32ProcessID to zero.
	Then, for each additional process in the snapshot,
	call CreateToolhelp32Snapshot again,
	specifying its process identifier and the
	TH32CS_SNAPHEAPLIST or TH32_SNAPMODULE value.

	When taking snapshots that include heaps and modules for a
	process other than the current process, the CreateToolhelp32Snapshot
	function can fail or return incorrect information for a variety of reasons.
	For example, if the loader data table in the target process is corrupted or
	not initialized, or if the module list changes during the function call as
	a result of DLLs being loaded or unloaded, the function might fail with
	ERROR_BAD_LENGTH or other error code. Ensure that the target process was not
	started in a suspended state, and try calling the function again.
	If the function fails with ERROR_BAD_LENGTH when called with TH32CS_SNAPMODULE
	or TH32CS_SNAPMODULE32, call the function again until it succeeds.

	The TH32CS_SNAPMODULE and TH32CS_SNAPMODULE32 flags do not retrieve handles for
	modules that were loaded with the LOAD_LIBRARY_AS_DATAFILE or similar flags.
	For more information, see LoadLibraryEx.

	If the function succeeds, it returns an open handle to the specified snapshot.

	For more information, see: https://docs.microsoft.com/en-us/windows/win32/api/tlhelp32/nf-tlhelp32-createtoolhelp32snapshot
*/
func CreateToolhelp32Snapshot(flags ThFlags, pid uint32) (win32.Handle, error) {
	snapshot := C.CreateToolhelp32Snapshot(C.DWORD(flags), C.DWORD(pid))
	if unsafe.Pointer(snapshot) == C.INVALID_HANDLE_VALUE {
		return win32.Handle(unsafe.Pointer(snapshot)), GetLastError()
	}

	return win32.Handle(unsafe.Pointer(snapshot)), nil
}

/*
	Process32First retrieves information about the first process encountered in a system snapshot.

	Returns TRUE if the first entry of the process list has been copied to the buffer or FALSE otherwise.
	The ERROR_NO_MORE_FILES error value is returned by the GetLastError function if no processes exist
	or the snapshot does not contain process information.

	The calling application must set the Size member of ProcessEntry32 to the size, in bytes, of the structure.

	For more info, see: https://docs.microsoft.com/en-us/windows/win32/api/tlhelp32/nf-tlhelp32-process32first
*/
func Process32First(snapShot win32.Handle, pe *ProcessEntry32) (bool, error) {
	if C.Process32First(C.HANDLE(unsafe.Pointer(snapShot)), (*C.PROCESSENTRY32)(unsafe.Pointer(pe))) == 0 {
		return false, GetLastError()
	}

	return true, nil
}

/*
	Process32Next retrieves information about the next process recorded in a system snapshot.

	To retrieve information about the first process recorded in a snapshot, use the Process32First function.

	For more info, see: https://docs.microsoft.com/en-us/windows/win32/api/tlhelp32/nf-tlhelp32-process32next
*/
func Process32Next(snapShot win32.Handle, pe *ProcessEntry32) (bool, error) {
	if C.Process32Next(C.HANDLE(unsafe.Pointer(snapShot)), (*C.PROCESSENTRY32)(unsafe.Pointer(pe))) == 0 {
		return false, GetLastError()
	}

	return true, nil
}
