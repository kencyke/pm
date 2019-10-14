package main

import (
	"fmt"
	"os"
	"bufio"
	"unsafe"
)

type Person struct {
	Name string
	Age  int
}

func main() {
	var p *Person

	p = &Person{
		Name: "Rob",
		Age:  63,
	}

	fmt.Printf("caller pid: %d\n", os.Getppid())
	fmt.Printf("data: %+v\n", p)
	fmt.Printf("address: %p\n", p)
	fmt.Printf("uintptr: 0x%x\n", uintptr(unsafe.Pointer(p)))
	fmt.Printf("size: %d\n", unsafe.Sizeof(p))

	// Wait to exit until stdin is closed.
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
}
