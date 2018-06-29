package tensority

// #cgo windows,386 CFLAGS: -I.
// #cgo windows,386 LDFLAGS: -L./lib/ -l:cSimdTs_win32.o -lstdc++ -lgomp -lpthread
// #cgo windows,amd64 CFLAGS: -I.
// #cgo windows,amd64 LDFLAGS: -L./lib/ -l:cSimdTs_win64.o -lstdc++ -lgomp
// #cgo linux,386 CFLAGS: -I.
// #cgo linux,386 LDFLAGS: -L./lib/ -l:cSimdTs_linux32.o -lstdc++ -lgomp -lpthread
// #cgo linux,amd64 CFLAGS: -I.
// #cgo linux,amd64 LDFLAGS: -L./lib/ -l:cSimdTs_linux64.o -lstdc++ -lgomp -lpthread
// #cgo darwin,amd64 CFLAGS: -I. -I/usr/local/opt/llvm/include
// #cgo darwin,amd64 LDFLAGS: -L./lib/ -lcSimdTs_darwin64.o -lstdc++ -lomp -L/usr/local/opt/llvm/lib
// #include "./lib/cSimdTs.h"
import "C"

import (
	"runtime"
	"unsafe"

	"github.com/golang/groupcache/lru"

	"github.com/bytom/crypto/sha3pool"
	"github.com/bytom/protocol/bc"
   	// cfg "github.com/bytom/config"

	


	log "github.com/sirupsen/logrus"
)

const maxAIHashCached = 64

func legacyAlgorithm(hash, seed *bc.Hash) *bc.Hash {
	cache := calcSeedCache(seed.Bytes())
	data := mulMatrix(hash.Bytes(), cache)
	return hashMatrix(data)
}

func cgoAlgorithm(blockHeader, seed *bc.Hash) *bc.Hash {
	bhBytes := blockHeader.Bytes()
	sdBytes := seed.Bytes()

	// Get the array pointers from the corresponding slices
	bhPtr := (*C.uchar)(unsafe.Pointer(&bhBytes[0]))
	seedPtr := (*C.uchar)(unsafe.Pointer(&sdBytes[0]))

	resPtr := C.SimdTs(bhPtr, seedPtr)

	res := bc.NewHash(*(*[32]byte)(unsafe.Pointer(resPtr)))
	return &res
}

func algorithm(hash, seed *bc.Hash) *bc.Hash {
	if (runtime.GOOS == "windows" || runtime.GOOS == "linux" || (runtime.GOOS == "darwin" && runtime.GOARCH == "amd64")) /*&& cfg.Config.Simd.Enable*/ {
		log.Info("hhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhh")
		return cgoAlgorithm(hash, seed)
	} else {
		return legacyAlgorithm(hash, seed)
	}
}

func calcCacheKey(hash, seed *bc.Hash) *bc.Hash {
	var b32 [32]byte
	sha3pool.Sum256(b32[:], append(hash.Bytes(), seed.Bytes()...))
	key := bc.NewHash(b32)
	return &key
}

// Cache is create for cache the tensority result
type Cache struct {
	lruCache *lru.Cache
}

// NewCache create a cache struct
func NewCache() *Cache {
	return &Cache{lruCache: lru.New(maxAIHashCached)}
}

// AddCache is used for add tensority calculate result
func (a *Cache) AddCache(hash, seed, result *bc.Hash) {
	key := calcCacheKey(hash, seed)
	a.lruCache.Add(*key, result)
}

// Hash is the real entry for call tensority algorithm
func (a *Cache) Hash(hash, seed *bc.Hash) *bc.Hash {
	key := calcCacheKey(hash, seed)
	if v, ok := a.lruCache.Get(*key); ok {
		return v.(*bc.Hash)
	}
	return algorithm(hash, seed)
}

// AIHash is created for let different package share same cache
var AIHash = NewCache()
