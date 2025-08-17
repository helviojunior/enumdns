//go:build windows
// +build windows

package disk

import (
	"fmt"
	"os"
	"path/filepath"
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows"
)

var (
	kernel32 = windows.NewLazySystemDLL("kernel32.dll")

	// GetDiskFreeSpaceEx - https://msdn.microsoft.com/en-us/library/windows/desktop/aa364937(v=vs.85).aspx
	// Retrieves information about the amount of space that is available on a disk volume,
	// which is the total amount of space, the total amount of free space, and the total
	// amount of free space available to the user that is associated with the calling thread.
	GetDiskFreeSpaceEx = kernel32.NewProc("GetDiskFreeSpaceExW")

	// GetDiskFreeSpace - https://msdn.microsoft.com/en-us/library/windows/desktop/aa364935(v=vs.85).aspx
	// Retrieves information about the specified disk, including the amount of free space on the disk.
	GetDiskFreeSpace = kernel32.NewProc("GetDiskFreeSpaceW")

	// GetVolumeInformation provides windows drive volume information.
	GetVolumeInformation = kernel32.NewProc("GetVolumeInformationW")
)

// GetInfo returns total and free bytes available in a directory, e.g. `C:\`.
// It returns free space available to the user (including quota limitations)
//
// https://msdn.microsoft.com/en-us/library/windows/desktop/aa364937(v=vs.85).aspx
func GetInfo(path string, _ bool) (info Info, err error) {
	// Stat to know if the path exists.
	if _, err = os.Stat(path); err != nil {
		return Info{}, err
	}

	lpFreeBytesAvailable := int64(0)
	lpTotalNumberOfBytes := int64(0)
	lpTotalNumberOfFreeBytes := int64(0)

	// Extract values safely
	// BOOL WINAPI GetDiskFreeSpaceEx(
	// _In_opt_  LPCTSTR         lpDirectoryName,
	// _Out_opt_ PULARGE_INTEGER lpFreeBytesAvailable,
	// _Out_opt_ PULARGE_INTEGER lpTotalNumberOfBytes,
	// _Out_opt_ PULARGE_INTEGER lpTotalNumberOfFreeBytes
	// );
	_, _, _ = GetDiskFreeSpaceEx.Call(uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(path))),
		uintptr(unsafe.Pointer(&lpFreeBytesAvailable)),
		uintptr(unsafe.Pointer(&lpTotalNumberOfBytes)),
		uintptr(unsafe.Pointer(&lpTotalNumberOfFreeBytes)))

	if uint64(lpTotalNumberOfFreeBytes) > uint64(lpTotalNumberOfBytes) {
		return info, fmt.Errorf("detected free space (%d) > total drive space (%d), fs corruption at (%s). please run 'fsck'",
			uint64(lpTotalNumberOfFreeBytes), uint64(lpTotalNumberOfBytes), path)
	}

	info = Info{
		Total:  uint64(lpTotalNumberOfBytes),
		Free:   uint64(lpTotalNumberOfFreeBytes),
		Used:   uint64(lpTotalNumberOfBytes) - uint64(lpTotalNumberOfFreeBytes),
		FSType: getFSType(path),
	}

	// Return values of GetDiskFreeSpace()
	lpSectorsPerCluster := uint32(0)
	lpBytesPerSector := uint32(0)
	lpNumberOfFreeClusters := uint32(0)
	lpTotalNumberOfClusters := uint32(0)

	// Extract values safely
	// BOOL WINAPI GetDiskFreeSpace(
	//   _In_  LPCTSTR lpRootPathName,
	//   _Out_ LPDWORD lpSectorsPerCluster,
	//   _Out_ LPDWORD lpBytesPerSector,
	//   _Out_ LPDWORD lpNumberOfFreeClusters,
	//   _Out_ LPDWORD lpTotalNumberOfClusters
	// );
	_, _, _ = GetDiskFreeSpace.Call(uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(path))),
		uintptr(unsafe.Pointer(&lpSectorsPerCluster)),
		uintptr(unsafe.Pointer(&lpBytesPerSector)),
		uintptr(unsafe.Pointer(&lpNumberOfFreeClusters)),
		uintptr(unsafe.Pointer(&lpTotalNumberOfClusters)))

	info.Files = uint64(lpTotalNumberOfClusters)
	info.Ffree = uint64(lpNumberOfFreeClusters)

	return info, nil
}

// getFSType returns the filesystem type of the underlying mounted filesystem
func getFSType(path string) string {
	volumeNameSize, nFileSystemNameSize := uint32(260), uint32(260)
	var lpVolumeSerialNumber uint32
	var lpFileSystemFlags, lpMaximumComponentLength uint32
	var lpFileSystemNameBuffer, volumeName [260]uint16
	ps := syscall.StringToUTF16Ptr(filepath.VolumeName(path))

	// Extract values safely
	// BOOL WINAPI GetVolumeInformation(
	// _In_opt_  LPCTSTR lpRootPathName,
	// _Out_opt_ LPTSTR  lpVolumeNameBuffer,
	// _In_      DWORD   nVolumeNameSize,
	// _Out_opt_ LPDWORD lpVolumeSerialNumber,
	// _Out_opt_ LPDWORD lpMaximumComponentLength,
	// _Out_opt_ LPDWORD lpFileSystemFlags,
	// _Out_opt_ LPTSTR  lpFileSystemNameBuffer,
	// _In_      DWORD   nFileSystemNameSize
	// );

	_, _, _ = GetVolumeInformation.Call(uintptr(unsafe.Pointer(ps)),
		uintptr(unsafe.Pointer(&volumeName)),
		uintptr(volumeNameSize),
		uintptr(unsafe.Pointer(&lpVolumeSerialNumber)),
		uintptr(unsafe.Pointer(&lpMaximumComponentLength)),
		uintptr(unsafe.Pointer(&lpFileSystemFlags)),
		uintptr(unsafe.Pointer(&lpFileSystemNameBuffer)),
		uintptr(nFileSystemNameSize))

	return syscall.UTF16ToString(lpFileSystemNameBuffer[:])
}
