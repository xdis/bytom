package main

// // #cgo windows,386 CFLAGS: -I.
// // #cgo windows,386 LDFLAGS: -L./lib/ -l:cSimdTs_win32.o -lstdc++ -lgomp -lpthread
// // #cgo windows,amd64 CFLAGS: -I.
// // #cgo windows,amd64 LDFLAGS: -L./lib/ -l:cSimdTs_win64.o -lstdc++ -lgomp -lpthread
// // #cgo linux,386 CFLAGS: -I.
// // #cgo linux,386 LDFLAGS: -L./lib/ -l:cSimdTs_linux32.o -lstdc++ -lgomp -lpthread
// // #cgo linux,amd64 CFLAGS: -I.
// // #cgo linux,amd64 LDFLAGS: -L./lib/ -l:cSimdTs_linux64.o -lstdc++ -lgomp -lpthread
// // #cgo darwin,amd64 CFLAGS: -I. -I/usr/local/opt/llvm/include
// // #cgo darwin,amd64 LDFLAGS: -L./lib/ -lcSimdTs_darwin64.o -lstdc++ -lomp -L/usr/local/opt/llvm/lib
// // #include "./lib/cSimdTs.h"
import "C"

import (
    "unsafe"

    "github.com/bytom/protocol/bc"
)

var bh *bc.Hash
var seed *bc.Hash
var res *bc.Hash

func cgoAlgorithm() {
    bhBytes := bh.Bytes()
    sdBytes := seed.Bytes()

    // Get the array pointers from the corresponding slices
    bhPtr := (*C.uchar)(unsafe.Pointer(&bhBytes[0]))
    seedPtr := (*C.uchar)(unsafe.Pointer(&sdBytes[0]))

    resPtr := C.SimdTs(bhPtr, seedPtr)

    resHsh := bc.NewHash(*(*[32]byte)(unsafe.Pointer(resPtr)))
    res := &resHsh
}
