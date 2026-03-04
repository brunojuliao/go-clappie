//go:build windows

package filestore

import (
	"os"
	"syscall"
	"unsafe"
)

var (
	modkernel32    = syscall.NewLazyDLL("kernel32.dll")
	procLockFileEx = modkernel32.NewProc("LockFileEx")
	procUnlockFile = modkernel32.NewProc("UnlockFileEx")
)

const (
	lockfileExclusiveLock = 0x02
	lockfileFailImmediately = 0x01
)

func lockFile(f *os.File) error {
	h := syscall.Handle(f.Fd())
	ol := new(syscall.Overlapped)
	r1, _, err := procLockFileEx.Call(
		uintptr(h),
		uintptr(lockfileExclusiveLock),
		0,
		1, 0,
		uintptr(unsafe.Pointer(ol)),
	)
	if r1 == 0 {
		return err
	}
	return nil
}

func tryLockFile(f *os.File) bool {
	h := syscall.Handle(f.Fd())
	ol := new(syscall.Overlapped)
	r1, _, _ := procLockFileEx.Call(
		uintptr(h),
		uintptr(lockfileExclusiveLock|lockfileFailImmediately),
		0,
		1, 0,
		uintptr(unsafe.Pointer(ol)),
	)
	return r1 != 0
}

func unlockFile(f *os.File) {
	h := syscall.Handle(f.Fd())
	ol := new(syscall.Overlapped)
	procUnlockFile.Call(
		uintptr(h),
		0,
		1, 0,
		uintptr(unsafe.Pointer(ol)),
	)
}
