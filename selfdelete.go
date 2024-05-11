/*
   License: MIT

   Credit to @jonasLyk for the discovery of this method and LloydLabs for the initial C PoC code.

   References:
       - https://github.com/LloydLabs/delete-self-poc
       - https://twitter.com/jonasLyk/status/1350401461985955840
*/

package selfdelete

import (
	"os"
	"unsafe"

	"golang.org/x/sys/windows"
)

type FILE_RENAME_INFO struct {
	ReplaceIfExists bool
	RootDirectory   windows.Handle
	FileNameLength  uint32
	FileName        [1]uint16
}

type FILE_DISPOSITION_INFO struct {
	DeleteFile bool
}

func getHandle(path string) (windows.Handle, error) {

	ntPath, _ := windows.NewNTUnicodeString(`\??\` + path)

	var handle windows.Handle

	var oa windows.OBJECT_ATTRIBUTES
	oa.Length = uint32(unsafe.Sizeof(oa))
	oa.ObjectName = ntPath
	oa.Attributes = windows.OBJ_CASE_INSENSITIVE

	var iosb windows.IO_STATUS_BLOCK

	var allocationSize int64 = 0

	err := windows.NtCreateFile(&handle,
		windows.SYNCHRONIZE|windows.DELETE|windows.GENERIC_READ,
		&oa,
		&iosb,
		&allocationSize,
		windows.FILE_ATTRIBUTE_NORMAL,
		windows.FILE_SHARE_READ,
		windows.FILE_OPEN_IF,
		windows.FILE_SYNCHRONOUS_IO_NONALERT,
		0,
		0)

	if err != nil {
		return handle, err
	}

	return handle, err
}

func renameFileInformation(handle windows.Handle) error {

	dataStream, err := windows.UTF16FromString(":bbq")

	if err != nil {
		return err
	}

	pDataStream := &dataStream[0]

	var iosb windows.IO_STATUS_BLOCK

	var fileRenameInfo FILE_RENAME_INFO

	fileRenameInfo.FileNameLength = uint32(unsafe.Sizeof(pDataStream))

	windows.NewLazyDLL("ntdll.dll").NewProc("RtlCopyMemory").Call(
		uintptr(unsafe.Pointer(&fileRenameInfo.FileName[0])),
		uintptr(unsafe.Pointer(pDataStream)),
		unsafe.Sizeof(pDataStream),
	)

	if err != nil {
		return err
	}

	err = windows.NtSetInformationFile(handle,
		&iosb,
		(*byte)(unsafe.Pointer(&fileRenameInfo)),
		uint32(unsafe.Sizeof(fileRenameInfo)+unsafe.Sizeof(pDataStream)),
		windows.FileRenameInformation)

	if err != nil {
		return err
	}

	return nil
}

func deleteFileInformation(handle windows.Handle) error {

	var fileDispositionInfo = FILE_DISPOSITION_INFO{DeleteFile: true}

	var iosb windows.IO_STATUS_BLOCK

	err := windows.NtSetInformationFile(handle,
		&iosb,
		(*byte)(unsafe.Pointer(&fileDispositionInfo)),
		uint32(unsafe.Sizeof(fileDispositionInfo)),
		windows.FileDispositionInformation)

	if err != nil {
		return err
	}

	return nil
}

func SelfDelete() error {
	var handle windows.Handle

	fullFilePath, err := os.Executable()

	if err != nil {
		return err
	}

	handle, err = getHandle(fullFilePath)

	if err != nil {
		windows.CloseHandle(handle)
		return err
	}

	err = renameFileInformation(handle)

	if err != nil {
		windows.CloseHandle(handle)
		return err
	}

	windows.CloseHandle(handle)

	handle, err = getHandle(fullFilePath)

	if err != nil {
		return err
	}

	err = deleteFileInformation(handle)

	if err != nil {
		windows.CloseHandle(handle)
		return err
	}

	windows.CloseHandle(handle)

	return nil
}
