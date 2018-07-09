package tensority

import (
	// "runtime"
	"plugin"

	"github.com/golang/groupcache/lru"

	"github.com/bytom/crypto/sha3pool"
	"github.com/bytom/protocol/bc"
   	// cfg "github.com/bytom/config"

	


	log "github.com/sirupsen/logrus"
)

const maxAIHashCached = 64

var UseSIMD = true

func legacyAlgorithm(bh, seed *bc.Hash) *bc.Hash {
	cache := calcSeedCache(seed.Bytes())
	data := mulMatrix(bh.Bytes(), cache)
	return hashMatrix(data)
}

func cgoAlgorithm(bh, seed *bc.Hash) *bc.Hash {
	p, err := plugin.Open("plugin_name.so")
	if err != nil {
	panic(err)
	}
	v, err := p.Lookup("bh")
	if err != nil {
	panic(err)
	}
	f, err := p.Lookup("F")
	if err != nil {
	panic(err)
	}

	*v.(*int) = 7
	f.(func())() // prints "Hello, number 7"

	return bh
}

func algorithm(bh, seed *bc.Hash) *bc.Hash {
	if false {
		log.Info("hhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhh")
		return cgoAlgorithm(bh, seed)
	} else {
		return legacyAlgorithm(bh, seed)
	}
}

func calcCacheKey(bh, seed *bc.Hash) *bc.Hash {
	var b32 [32]byte
	sha3pool.Sum256(b32[:], append(bh.Bytes(), seed.Bytes()...))
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
func (a *Cache) AddCache(bh, seed, result *bc.Hash) {
	key := calcCacheKey(bh, seed)
	a.lruCache.Add(*key, result)
}

// Hash is the real entry for call tensority algorithm
func (a *Cache) Hash(bh, seed *bc.Hash) *bc.Hash {
	key := calcCacheKey(bh, seed)
	if v, ok := a.lruCache.Get(*key); ok {
		return v.(*bc.Hash)
	}
	return algorithm(bh, seed)
}

// AIHash is created for let different package share same cache
var AIHash = NewCache()
