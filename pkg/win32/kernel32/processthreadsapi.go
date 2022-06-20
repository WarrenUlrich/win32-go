package kernel32

/*
	#include <processthreadsapi.h>
*/
import "C"

import (
	"unsafe"

	"github.com/warrenulrich/win32-go/pkg/win32"
)

type ProcessAccess uint32

const (
	/*
		DELETE is the right to delete the object.
	*/
	DELETE ProcessAccess = 0x00010000

	/*
		SYNCHRONIZE is the right to use the object for synchronization.
		This enables a thread to wait until the object is in the signaled state.
		Some object types do not support this access right.
	*/
	SYNCHRONIZE ProcessAccess = 0x00100000

	/*
		READ_CONTROL is the right to read the information in the object's security descriptor,
		not including the information in the system access control list (SACL).
	*/
	READ_CONTROL ProcessAccess = 0x00020000

	/*
		WRITE_DAC is the right to modify the discretionary access control list (DACL) in the object's security descriptor.
	*/
	WRITE_DAC ProcessAccess = 0x00040000

	/*
		WRITE_OWNER is the right to change the owner in the object's security descriptor.
	*/
	WRITE_OWNER ProcessAccess = 0x00080000

	/*
		STANDARD_RIGHTS_ALL combines DELETE, READ_CONTROL, WRITE_DAC, WRITE_OWNER, and SYNCHRONIZE access.
	*/
	STANDARD_RIGHTS_ALL ProcessAccess = DELETE | READ_CONTROL | WRITE_DAC | WRITE_OWNER | SYNCHRONIZE

	/*
		STANDARD_RIGHTS_EXECUTE is currently defined to equal READ_CONTROL.
	*/
	STANDARD_RIGHTS_EXECUTE ProcessAccess = READ_CONTROL

	/*
		STANDARD_RIGHTS_READ is currently defined to equal READ_CONTROL.
	*/
	STANDARD_RIGHTS_READ ProcessAccess = READ_CONTROL

	/*
		STANDARD_RIGHTS_REQUIRED combines DELETE, READ_CONTROL, WRITE_DAC, and WRITE_OWNER access.
	*/
	STANDARD_RIGHTS_REQUIRED ProcessAccess = DELETE | READ_CONTROL | WRITE_DAC | WRITE_OWNER

	/*
		STANDARD_RIGHTS_WRITE is currently defined to equal READ_CONTROL.
	*/
	STANDARD_RIGHTS_WRITE ProcessAccess = READ_CONTROL

	/*
		PROCESS_ALL_ACCESS is all possible access rights for a process object.Windows Server 2003 and Windows XP:
		The size of the PROCESS_ALL_ACCESS flag increased on Windows Server 2008 and Windows Vista.
		If an application compiled for Windows Server 2008 and Windows Vista is run on Windows Server 2003 or Windows XP,
		the PROCESS_ALL_ACCESS flag is too large and the function specifying this flag fails with ERROR_ACCESS_DENIED.
		To avoid this problem, specify the minimum set of access rights required for the operation.
		If PROCESS_ALL_ACCESS must be used, set _WIN32_WINNT to the minimum operating system targeted by your application
		(for example, #define _WIN32_WINNT _WIN32_WINNT_WINXP).
		For more information, see Using the Windows Headers.
	*/
	PROCESS_ALL_ACCESS ProcessAccess = STANDARD_RIGHTS_REQUIRED | SYNCHRONIZE | 0xFFF

	/*
		PROCESS_CREATE_PROCESS is required to use this process as the parent process with PROC_THREAD_ATTRIBUTE_PARENT_PROCESS.
	*/
	PROCESS_CREATE_PROCESS ProcessAccess = 0x0080

	/*
		PROCESS_CREATE_THREAD is required to create a thread in the process.
	*/
	PROCESS_CREATE_THREAD ProcessAccess = 0x0002

	/*
		PROCESS_DUP_HANDLE is required to duplicate a handle using DuplicateHandle.
	*/
	PROCESS_DUP_HANDLE ProcessAccess = 0x0040

	/*
		PROCESS_QUERY_INFORMATION is required to retrieve certain information about a process,
		 such as its token, exit code, and priority class (see OpenProcessToken).
	*/
	PROCESS_QUERY_INFORMATION ProcessAccess = 0x0400

	/*
		PROCESS_QUERY_LIMITED_INFORMATION is required to retrieve certain information about a process
		(see GetExitCodeProcess, GetPriorityClass, IsProcessInJob, QueryFullProcessImageName).
		A handle that has the PROCESS_QUERY_INFORMATION access right is automatically
		granted PROCESS_QUERY_LIMITED_INFORMATION.Windows Server 2003 and Windows XP:
		This access right is not supported.
	*/
	PROCESS_QUERY_LIMITED_INFORMATION ProcessAccess = 0x1000

	/*
		PROCESS_SET_INFORMATION is required to set certain information about a process, such as its priority class (see SetPriorityClass).
	*/
	PROCESS_SET_INFORMATION ProcessAccess = 0x0200

	/*
		PROCESS_SET_QUOTA is required to set memory limits using SetProcessWorkingSetSize.
	*/
	PROCESS_SET_QUOTA ProcessAccess = 0x0100

	/*
		PROCESS_SUSPEND_RESUME is required to suspend or resume a process.
	*/
	PROCESS_SUSPEND_RESUME ProcessAccess = 0x080

	/*
		PROCESS_TERMINATE is required to terminate a process using TerminateProcess.
	*/
	PROCESS_TERMINATE ProcessAccess = 0x0001

	/*
		PROCESS_VM_OPERATION is required to perform an operation on the address space of a process (see VirtualProtectEx and WriteProcessMemory).
	*/
	PROCESS_VM_OPERATION ProcessAccess = 0x0008

	/*
		PROCESS_VM_READ is required to read memory in a process using ReadProcessMemory.
	*/
	PROCESS_VM_READ ProcessAccess = 0x0010

	/*
		PROCESS_VM_WRITE is required to write to memory in a process using WriteProcessMemory.
	*/
	PROCESS_VM_WRITE ProcessAccess = 0x0020
)

/*
	OpenProcess opens an existing local process object.
*/
func OpenProcess(desiredAccess ProcessAccess, inheritHandle bool, processId uint32) (win32.Handle, error) {
	var inherit C.BOOL
	if inheritHandle {
		inherit = 1
	}

	handle := C.OpenProcess(C.DWORD(desiredAccess), inherit, C.DWORD(processId))
	if handle == nil {
		return 0, GetLastError()
	}

	return win32.Handle(unsafe.Pointer(handle)), nil
}

/*
	TerminateProcess terminates the specified process and all of its threads.
*/
func TerminateProcess(process win32.Handle, exitCode uint32) error {
	if C.TerminateProcess(C.HANDLE(unsafe.Pointer(process)), C.UINT(exitCode)) == 0 {
		return GetLastError()
	}

	return nil
}
