package pm

import (
	"unsafe"

	"github.com/kencyke/pm/syscall/linux"
)

// Some codes originally taken from https://stackoverflow.com/questions/32223562/how-to-convert-uintptr-to-byte-in-golang

const sizeOfUintPtr = unsafe.Sizeof(uintptr(0))

func uintptrToBytes(u uintptr) []byte {
	return (*[sizeOfUintPtr]byte)(unsafe.Pointer(&u))[:]
}

func uintptrToIovecBase(u uintptr) *byte {
	return &uintptrToBytes(u)[0]
}

// ReadAddress copies memory from another process into an already allocated byte buffer.
func ReadAddress(pid int, address uintptr, buffer []byte) error {
	var length uint64
	if len(buffer) > 0 {
		length = uint64(len(buffer))
	}
	liov := make([]linux.Iovec, 0)
	liov = append(liov, linux.Iovec{Base: &buffer[0], Len: length})
	riov := make([]linux.Iovec, 0)
	riov = append(riov, linux.Iovec{Base: uintptrToIovecBase(address), Len: length})
	_, err := linux.ProcessVMReadv(pid, liov, 1, riov, 1, 0)

	return err
}

// CopyAddress copies memory from another process by allocating memory for you.
func CopyAddress(pid int, address uintptr, length uint64) (data []byte, err error) {
	data = make([]byte, length)
	liov := make([]linux.Iovec, 0)
	liov = append(liov, linux.Iovec{Base: &data[0], Len: length})
	riov := make([]linux.Iovec, 0)
	riov = append(riov, linux.Iovec{Base: uintptrToIovecBase(address), Len: length})
	_, err = linux.ProcessVMReadv(pid, liov, 1, riov, 1, 0)

	return data, err
}
