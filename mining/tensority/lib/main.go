package main

// #cgo linux windows LDFLAGS: -L. -l:cSimdTs.o -lstdc++ -lgomp -lpthread
// #cgo darwin,amd64 CFLAGS: -I. -I/usr/local/opt/llvm/include
// #cgo darwin,amd64 LDFLAGS: -L. -lcSimdTs_darwin64.o -lstdc++ -lomp -L/usr/local/opt/llvm/lib
// #include "cSimdTs.h"
import "C"

import (
    "unsafe"

    "github.com/bytom/protocol/bc"
)

var BH bc.Hash
var SEED bc.Hash
var RES bc.Hash

func CgoAlgorithm() {
    // type conversion
    bhBytes := BH.Bytes()
    sdBytes := SEED.Bytes()
    bhPtr := (*C.uchar)(unsafe.Pointer(&bhBytes[0]))
    seedPtr := (*C.uchar)(unsafe.Pointer(&sdBytes[0]))
    
    // invoke c func
    resPtr := C.SimdTs(bhPtr, seedPtr)

    // type conversion
    RES = bc.NewHash(*(*[32]byte)(unsafe.Pointer(resPtr)))
}
