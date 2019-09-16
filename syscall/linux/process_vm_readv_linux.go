package linux

import (
	"unsafe"

	"golang.org/x/sys/unix"
)

// Iovec is a structure that contains the starting address and the number of bytes.
type Iovec unix.Iovec

var (
	errEINVAL error = unix.EINVAL
	errEFAULT error = unix.EFAULT
	errENOMEM error = unix.ENOMEM
	errEPERM  error = unix.EPERM
	errESRCH  error = unix.ESRCH
)

func errnoErr(e unix.Errno) error {
	switch e {
	case 0:
		return nil
	case unix.EINVAL:
		return errEINVAL
	case unix.EFAULT:
		return errEFAULT
	case unix.ENOMEM:
		return errENOMEM
	case unix.EPERM:
		return errEPERM
	case unix.ESRCH:
		return errESRCH
	}
	return e
}

// ProcessVMReadv transfers data from the remote process to the local process.
func ProcessVMReadv(pid int, liov []Iovec, liovcnt uint64, riov []Iovec, riovcnt uint64, flags uint64) (size int, err error) {
	var lp, rp unsafe.Pointer
	if len(liov) > 0 {
		lp = unsafe.Pointer(&liov[0])
	}
	if len(riov) > 0 {
		rp = unsafe.Pointer(&riov[0])
	}
	r0, _, e1 := unix.Syscall6(
		unix.SYS_PROCESS_VM_READV,
		uintptr(pid),
		uintptr(lp),
		uintptr(liovcnt),
		uintptr(rp),
		uintptr(riovcnt),
		uintptr(flags),
	)
	size = int(r0)
	if e1 != 0 {
		err = errnoErr(e1)
	}
	return
}
