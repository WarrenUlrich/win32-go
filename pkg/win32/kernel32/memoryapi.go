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

/*
	VirtualAlloc reserves, commits, or changes the state of a region of
	pages in the virtual address space of the calling process.
	Memory allocated by this function is automatically initialized to zero.

	To allocate memory in the address space of another process, use the VirtualAllocEx function

	If the function succeeds, the return value is the base address of the allocated region of pages.
	
	For more information, see: https://docs.microsoft.com/en-us/windows/win32/api/memoryapi/nf-memoryapi-virtualalloc
*/
func VirtualAlloc(addr uintptr, size uintptr, allocType uint32, flProtect uint32) (uintptr, error) {
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
func VirtualAllocEx(ph win32.Handle, baseAddr uintptr, size uintptr, allocType AllocType, flProtect uint32) (uintptr, error) {
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
func ReadProcessMemory(ph win32.Handle, baseAddr uintptr, dest unsafe.Pointer, size uintptr) (uint, error) {
	var bytesRead uint
	if C.ReadProcessMemory(C.HANDLE(unsafe.Pointer(ph)), C.LPCVOID(baseAddr), (C.LPVOID)(dest), C.ulonglong(size), (*C.ulonglong)(unsafe.Pointer(&bytesRead))) == 0 {
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
func WriteProcessMemory(ph win32.Handle, baseAddr uintptr, src unsafe.Pointer, size uintptr) (uint, error) {
	var bytesWritten uint
	if C.WriteProcessMemory(C.HANDLE(unsafe.Pointer(ph)), C.LPVOID(baseAddr), (C.LPCVOID)(src), C.ulonglong(size), (*C.ulonglong)(unsafe.Pointer(&bytesWritten))) == 0 {
		return bytesWritten, GetLastError()
	}

	return bytesWritten, nil
}