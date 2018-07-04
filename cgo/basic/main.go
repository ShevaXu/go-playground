package main

// #include <stdlib.h>
// #include <stdio.h>
import "C"

// other imports, "C" cannot be in group import
import (
	"fmt"
	"time"
	"unsafe"
)

func grandom() int {
	// if explicitly converse:
	// var r C.long = C.random()
	// return int(r)
	return int(C.random())
}

func gseed(i int) {
	C.srandom(C.uint(i))
}

func gprint(s string) {
	cs := C.CString(s)
	// a common idiom in cgo programs is to defer the free
	// immediately after allocating
	defer C.free(unsafe.Pointer(cs))
	C.fputs(cs, (*C.FILE)(C.stdout))
}

func main() {
	gseed(int(time.Now().Unix()))
	gprint(fmt.Sprintln(grandom()))
}
