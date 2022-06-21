package kernel32

/*
	#include <memoryapi.h>
*/
import "C"

import (
	"unsafe"

	"github.com/warrenulrich/win32-go/pkg/win32"
)

type AllocType uint32

const (
	/*
		MEM_COMMIT Allocates memory charges (from the overall size of memory and the paging files on disk) for the specified reserved memory pages.
		The function also guarantees that when the caller later initially accesses the memory, the contents will be zero.
		Actual physical pages are not allocated unless/until the virtual addresses are actually accessed.
		To reserve and commit pages in one step, call VirtualAlloc with MEM_COMMIT | MEM_RESERVE.

		Attempting to commit a specific address range by specifying MEM_COMMIT without MEM_RESERVE and a non-NULL lpAddress
		fails unless the entire range has already been reserved. The resulting error code is ERROR_INVALID_ADDRESS.

		An attempt to commit a page that is already committed does not cause the function to fail.
		This means that you can commit pages without first determining the current commitment state of each page.

		If lpAddress specifies an address within an enclave, flAllocationType must be MEM_COMMIT.

		For more information, see: https://docs.microsoft.com/en-us/windows/win32/api/memoryapi/nf-memoryapi-virtualalloc
	*/
	MEM_COMMIT AllocType = 0x00001000

	/*
		MEM_RESERVE Reserves a range of the process's virtual address space without allocating
		any actual physical storage in memory or in the paging file on disk.
		You can commit reserved pages in subsequent calls to the VirtualAlloc function.
		To reserve and commit pages in one step, call VirtualAlloc with MEM_COMMIT | MEM_RESERVE.

		Other memory allocation functions, such as malloc and LocalAlloc,
		cannot use a reserved range of memory until it is released.

		For more information, see: https://docs.microsoft.com/en-us/windows/win32/api/memoryapi/nf-memoryapi-virtualalloc
	*/
	MEM_RESERVE AllocType = 0x00002000

	/*
		MEM_RESET Indicates that data in the memory range specified by lpAddress and dwSize is no longer of interest.
		The pages should not be read from or written to the paging file. However,
		the memory block will be used again later, so it should not be decommitted.
		This value cannot be used with any other value.

		Using this value does not guarantee that the range operated on with MEM_RESET will contain zeros.
		If you want the range to contain zeros, decommit the memory and then recommit it.

		When you specify MEM_RESET, the VirtualAlloc function ignores the value of flProtect.
		However, you must still set flProtect to a valid protection value, such as PAGE_NOACCESS.

		VirtualAlloc returns an error if you use MEM_RESET and the range of memory is mapped to a file.
		A shared view is only acceptable if it is mapped to a paging file.

		For more information, see: https://docs.microsoft.com/en-us/windows/win32/api/memoryapi/nf-memoryapi-virtualalloc
	*/
	MEM_RESET AllocType = 0x00080000

	/*
		MEM_RESET_UNDO should only be called on an address range to which MEM_RESET was successfully applied earlier.
		It indicates that the data in the specified memory range specified by lpAddress and dwSize is of interest
		to the caller and attempts to reverse the effects of MEM_RESET.
		If the function succeeds, that means all data in the specified address range is intact.
		If the function fails, at least some of the data in the address range has been replaced with zeroes.

		This value cannot be used with any other value. If MEM_RESET_UNDO is called on an address range which was not MEM_RESET earlier,
		the behavior is undefined. When you specify MEM_RESET, the VirtualAlloc function ignores the value of flProtect.
		However, you must still set flProtect to a valid protection value, such as PAGE_NOACCESS.

		Windows Server 2008 R2, Windows 7, Windows Server 2008, Windows Vista, Windows Server 2003 and Windows XP:
		The MEM_RESET_UNDO flag is not supported until Windows 8 and Windows Server 2012.

		For more information, see: https://docs.microsoft.com/en-us/windows/win32/api/memoryapi/nf-memoryapi-virtualalloc
	*/
	MEM_RESET_UNDO AllocType = 0x1000000

	/*
		MEM_LARGE_PAGES Allocates memory using large page support.

		The size and alignment must be a multiple of the large-page minimum.
		To obtain this value, use the GetLargePageMinimum function.

		If you specify this value, you must also specify MEM_RESERVE and MEM_COMMIT.

		For more information, see: https://docs.microsoft.com/en-us/windows/win32/api/memoryapi/nf-memoryapi-virtualalloc
	*/
	MEM_LARGE_PAGES AllocType = 0x20000000

	/*
		MEM_PHYSICAL Reserves an address range that can be used to map Address Windowing Extensions (AWE) pages.

		This value must be used with MEM_RESERVE and no other values.

		For more information, see: https://docs.microsoft.com/en-us/windows/win32/api/memoryapi/nf-memoryapi-virtualalloc
	*/
	MEM_PHYSICAL AllocType = 0x00400000

	/*
		MEM_TOP_DOWN Allocates memory at the highest possible address.
		This can be slower than regular allocations, especially when there are many allocations.

		For more information, see: https://docs.microsoft.com/en-us/windows/win32/api/memoryapi/nf-memoryapi-virtualalloc
	*/
	MEM_TOP_DOWN AllocType = 0x100000

	/*
		MEM_WRITE_WATCH Causes the system to track pages that are written to in the allocated region.
		If you specify this value, you must also specify MEM_RESERVE.

		To retrieve the addresses of the pages that have been written to since the region was
		allocated or the write-tracking state was reset, call the GetWriteWatch function.
		To reset the write-tracking state, call GetWriteWatch or ResetWriteWatch.
		The write-tracking feature remains enabled for the memory region until the region is freed.

		For more information, see: https://docs.microsoft.com/en-us/windows/win32/api/memoryapi/nf-memoryapi-virtualalloc
	*/
	MEM_WRITE_WATCH AllocType = 0x00200000
)

type PageAccess uint32

const (
	/*
		PAGE_EXECUTE enables execute access to the committed region of pages.
		An attempt to write to the committed region results in an access violation.
		This flag is not supported by the CreateFileMapping function.
	*/
	PAGE_EXECUTE PageAccess = 0x10

	/*
		PAGE_EXECUTE_READ Enables execute or read-only access to
		the committed region of pages. An attempt to write to
		the committed region results in an access violation.

		Windows Server 2003 and Windows XP: This attribute
		is not supported by the CreateFileMapping function
		until Windows XP with SP2 and Windows Server 2003 with SP1.
	*/
	PAGE_EXECUTE_READ PageAccess = 0x20

	/*
		PAGE_EXECUTE_READWRITE enables execute, read-only, or read/write
		access to the committed region of pages.

		Windows Server 2003 and Windows XP: This attribute is
		not supported by the CreateFileMapping function until
		Windows XP with SP2 and Windows Server 2003 with SP1.
	*/
	PAGE_EXECUTE_READWRITE PageAccess = 0x40

	/*
		PAGE_EXECUTE_WRITECOPY enables execute, read-only, or copy-on-write
		access to a mapped view of a file mapping object. An attempt to
		write to a committed copy-on-write page results in a private
		copy of the page being made for the process. The private page is
		marked as PAGE_EXECUTE_READWRITE, and the change is written to the new page.

		This flag is not supported by the VirtualAlloc or VirtualAllocEx functions.
		Windows Vista, Windows Server 2003 and Windows XP:
		This attribute is not supported by the CreateFileMapping function until
		Windows Vista with SP1 and Windows Server 2008.
	*/
	PAGE_EXECUTE_WRITECOPY PageAccess = 0x80

	/*
		PAGE_NOACCESS disables all access to the committed region of pages.
		An attempt to read from, write to, or execute the committed
		region results in an access violation.

		This flag is not supported by the CreateFileMapping function.
	*/
	PAGE_NOACCESS PageAccess = 0x01

	/*
		PAGE_READONLY enables read-only access to the committed region of pages.
		An attempt to write to the committed region results in an access violation.
		If Data Execution Prevention is enabled, an attempt to execute code
		in the committed region results in an access violation.
	*/
	PAGE_READONLY PageAccess = 0x02

	/*
		PAGE_READWRITE enables read-only or read/write access to the
		committed region of pages. If Data Execution Prevention is
		enabled, attempting to execute code in the committed
		region results in an access violation.
	*/
	PAGE_READWRITE PageAccess = 0x04

	/*
		PAGE_WRITECOPY enables read-only or copy-on-write access to a
		mapped view of a file mapping object. An attempt to write to a
		committed copy-on-write page results in a private copy of the
		page being made for the process. The private page is marked as
		PAGE_READWRITE, and the change is written to the new page.
		If Data Execution Prevention is enabled, attempting to execute
		code in the committed region results in an access violation.

		This flag is not supported by the VirtualAlloc or VirtualAllocEx functions.
	*/
	PAGE_WRITECOPY PageAccess = 0x08

	/*
		PAGE_TARGETS_INVALID sets all locations in the pages as invalid targets for CFG.
		Used along with any execute page protection like PAGE_EXECUTE,
		PAGE_EXECUTE_READ, PAGE_EXECUTE_READWRITE and PAGE_EXECUTE_WRITECOPY.
		Any indirect call to locations in those pages will fail CFG checks and
		the process will be terminated. The default behavior for executable
		pages allocated is to be marked valid call targets for CFG.

		This flag is not supported by the VirtualProtect or CreateFileMapping functions.
	*/
	PAGE_TARGETS_INVALID PageAccess = 0x40000000

	/*
		PAGE_TARGETS_NO_UPDATE pages in the region will not have their CFG information
		updated while the protection changes for VirtualProtect. For example, if the
		pages in the region was allocated using PAGE_TARGETS_INVALID, then the invalid
		information will be maintained while the page protection changes.
		This flag is only valid when the protection changes to an executable type like
		PAGE_EXECUTE, PAGE_EXECUTE_READ, PAGE_EXECUTE_READWRITE and PAGE_EXECUTE_WRITECOPY.
		The default behavior for VirtualProtect protection change to executable
		is to mark all locations as valid call targets for CFG.
	*/
	PAGE_TARGETS_NO_UPDATE PageAccess = 0x40000000

	/*
		PAGE_GUARD pages in the region become guard pages. Any attempt to access a
		guard page causes the system to raise a STATUS_GUARD_PAGE_VIOLATION exception
		and turn off the guard page status. Guard pages thus act as a one-time access
		alarm. For more information, see Creating Guard Pages. When an access attempt
		leads the system to turn off guard page status, the underlying page protection takes over.
		If a guard page exception occurs during a system service, the service typically
		returns a failure status indicator. This value cannot be used with PAGE_NOACCESS.

		This flag is not supported by the CreateFileMapping function.
	*/
	PAGE_GUARD PageAccess = 0x100

	/*
		PAGE_NOCACHE Sets all pages to be non-cachable. Applications should not use
		this attribute except when explicitly required for a device. Using the interlocked
		functions with memory that is mapped with SEC_NOCACHE can result in an EXCEPTION_ILLEGAL_INSTRUCTION exception.
		The PAGE_NOCACHE flag cannot be used with the PAGE_GUARD, PAGE_NOACCESS, or PAGE_WRITECOMBINE flags.
		The PAGE_NOCACHE flag can be used only when allocating private memory with the VirtualAlloc,
		VirtualAllocEx, or VirtualAllocExNuma functions. To enable non-cached memory
		access for shared memory, specify the SEC_NOCACHE flag when calling the CreateFileMapping function.
	*/
	PAGE_NOCACHE PageAccess = 0x200

	/*
		PAGE_WRITECOMBINE Sets all pages to be write-combined.
		Applications should not use this attribute except when
		explicitly required for a device. Using the interlocked
		functions with memory that is mapped as write-combined
		can result in an EXCEPTION_ILLEGAL_INSTRUCTION exception.

		The PAGE_WRITECOMBINE flag cannot be specified with the
		PAGE_NOACCESS, PAGE_GUARD, and PAGE_NOCACHE flags.
		The PAGE_WRITECOMBINE flag can be used only when allocating
		private memory with the VirtualAlloc, VirtualAllocEx,
		or VirtualAllocExNuma functions. To enable write-combined
		memory access for shared memory, specify the SEC_WRITECOMBINE
		flag when calling the CreateFileMapping function.
	*/
	PAGE_WRITECOMBINE PageAccess = 0x400

	/*
		PAGE_ENCLAVE_DECOMMIT indicates that the page will be
		protected to prevent further use in an enclave.
		This flag must not be combined with any other flags.
		This flag is only valid for SGX2 enclaves.
	*/
	PAGE_ENCLAVE_DECOMMIT PageAccess = 0x80000000

	/*
		PAGE_ENCLAVE_THREAD_CONTROL The page contains a thread control structure (TCS).
	*/
	PAGE_ENCLAVE_THREAD_CONTROL PageAccess = 0x80000000

	/*
		PAGE_ENCLAVE_UNVALIDATED The page contents that you supply
		are excluded from measurement with the EEXTEND instruction
		of the Intel SGX programming model.
	*/
	PAGE_ENCLAVE_UNVALIDATED PageAccess = 0x80000000
)

/*
	VirtualAlloc reserves, commits, or changes the state of a region of
	pages in the virtual address space of the calling process.
	Memory allocated by this function is automatically initialized to zero.

	To allocate memory in the address space of another process, use the VirtualAllocEx function

	If the function succeeds, the return value is the base address of the allocated region of pages.

	For more information, see: https://docs.microsoft.com/en-us/windows/win32/api/memoryapi/nf-memoryapi-virtualalloc
*/
func VirtualAlloc(addr uintptr, size uintptr, allocType AllocType, flProtect PageAccess) (uintptr, error) {
	var baseAddr uintptr
	if C.VirtualAlloc(C.LPVOID(addr), C.SIZE_T(size), C.DWORD(allocType), C.DWORD(flProtect)) == nil {
		return baseAddr, GetLastError()
	}

	return baseAddr, nil
}

/*
	VirtualAllocEx reserves, commits, or changes the state of a region
	of memory within the virtual address space of a specified process.
	The function initializes the memory it allocates to zero.

	If the function succeeds, the return value is the base address of the allocated region of pages.

	For more information, see: https://docs.microsoft.com/en-us/windows/win32/api/memoryapi/nf-memoryapi-virtualallocex
*/
func VirtualAllocEx(ph win32.Handle, baseAddr uintptr, size uintptr, allocType AllocType, flProtect PageAccess) (uintptr, error) {
	if C.VirtualAllocEx(C.HANDLE(unsafe.Pointer(ph)), C.LPVOID(baseAddr), C.SIZE_T(size), C.DWORD(allocType), C.DWORD(flProtect)) == nil {
		return baseAddr, GetLastError()
	}

	return baseAddr, nil
}

/*
	ReadProcessMemory copies the data in the specified
	address range from the address space of the specified
	process into the specified buffer of the current process.
	Any process that has a handle with PROCESS_VM_READ access
	can call the function.

	The entire area to be read must be accessible,
	and if it is not accessible, the function fails.

	For more information, see: https://docs.microsoft.com/en-us/windows/win32/api/memoryapi/nf-memoryapi-readprocessmemory
*/
func ReadProcessMemory[T any](ph win32.Handle, baseAddr uintptr, buf *T, size uintptr) (uint, error) {
	var bytesRead uint
	if C.ReadProcessMemory(C.HANDLE(unsafe.Pointer(ph)), C.LPCVOID(baseAddr), (C.LPVOID)(buf), C.ulonglong(size), (*C.ulonglong)(unsafe.Pointer(&bytesRead))) == 0 {
		return bytesRead, GetLastError()
	}

	return bytesRead, nil
}

/*
	WriteProcessMemory copies the data from the specified buffer
	in the current process to the address range of the specified process.
	Any process that has a handle with PROCESS_VM_WRITE and PROCESS_VM_OPERATION
	access to the process to be written to can call the function.
	Typically but not always, the process with address
	space that is being written to is being debugged.

	The entire area to be written to must be accessible,
	and if it is not accessible, the function fails.

	For more information, see: https://docs.microsoft.com/en-us/windows/win32/api/memoryapi/nf-memoryapi-writeprocessmemory
*/
func WriteProcessMemory[T any](ph win32.Handle, baseAddr uintptr, src *T, size uintptr) (uint, error) {
	var bytesWritten uint
	if C.WriteProcessMemory(C.HANDLE(unsafe.Pointer(ph)), C.LPVOID(baseAddr), (C.LPCVOID)(src), C.ulonglong(size), (*C.ulonglong)(unsafe.Pointer(&bytesWritten))) == 0 {
		return bytesWritten, GetLastError()
	}

	return bytesWritten, nil
}
