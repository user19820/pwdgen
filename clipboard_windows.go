package main

import (
	"errors"
	"fmt"
	"syscall"
	"unsafe"
)

// NOTE: this code was not verified by running it,
// only used Claude to code review it. TESTING NEEDED!

const (
	CF_UNICODETEXT = 1
	GMEM_MOVEABLE  = 0x0002
)

var (
	user32           = syscall.MustLoadDLL("user32")
	openClipboard    = user32.MustFindProc("OpenClipboard")
	closeClipboard   = user32.MustFindProc("CloseClipboard")
	emptyClipboard   = user32.MustFindProc("EmptyClipboard")
	setClipboardData = user32.MustFindProc("SetClipboardData")

	kernel  = syscall.MustLoadDLL("kernel32")
	gAlloc  = kernel.MustFindProc("GlobalAlloc")
	gFree   = kernel.MustFindProc("GlobalFree")
	gLock   = kernel.MustFindProc("GlobalLock")
	gUnlock = kernel.MustFindProc("GlobalUnlock")
)

func copyToClipboard(pwd string) error {
	/*
		ASIDE: windows works way different.
		I have to call Win32 APIs, hence syscall.MustLoadDLL's.

		Their return values are not stored in the error, error usually
		doesn't contain helpful values, therefore I opted to discard it.

		When I am handing the generated password to the clipboard, the memory
		is given to Windows, which now says you are no longer in control of this
		memory. Therefore I cannot give my garbage-collector managed Go object
		to Windows, instead I have to manually allocate memory that is in Windows'
		control from the start and then copy my data to the allocated buffer, then
		hand THAT to windows.
	*/
	ocRet, _, _ := openClipboard.Call(0)
	if ocRet == 0 {
		return errors.New("Unable to open system clipboard.")
	}
	defer closeClipboard.Call()

	ecRet, _, _ := emptyClipboard.Call()
	if ecRet == 0 {
		return errors.New("Unable to empty the system clipboard.")
	}

	utf16Text, err := syscall.UTF16FromString(pwd)
	if err != nil {
		return fmt.Errorf("Unable to convert to UTF-16: %w", err)
	}

	allocRet, _, _ := gAlloc.Call(GMEM_MOVEABLE, uintptr(len(utf16Text)*2))
	if allocRet == 0 {
		return errors.New("Unable to allocate memory.")
	}

	lockRet, _, _ := gLock.Call(allocRet)
	if lockRet == 0 {
		gFree.Call(allocRet)
		return errors.New("Unable to lock allocated memory.")
	}
	copy(unsafe.Slice((*uint16)(unsafe.Pointer(lockRet)), len(utf16Text)), utf16Text)

	gUnlock.Call(allocRet)

	scRet, _, _ := setClipboardData.Call(CF_UNICODETEXT, allocRet)

	if scRet == 0 {
		gFree.Call(allocRet)
		return errors.New("Unable to set data to the system clipboard.")
	}

	return nil
}
