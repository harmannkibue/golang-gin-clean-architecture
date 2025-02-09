package httpCaching

import (
	"fmt"
	"github.com/patrickmn/go-cache"
	"time"
)

type ICache interface {
	GetCacheValue() (string, error)
	SetCacheValue(cacheValue string) error
}

type CacheParams struct {
	Key       string        `json:"key"`
	ExpiresIn time.Duration `json:"expiresIn"`
	Cache     *cache.Cache  `json:"cache"`
}

// New Returns a new instance of the cache -.
func New(Key string, expires time.Duration) ICache {
	// Initializing a cache instance -.
	c := cache.New(expires, expires)
	return &CacheParams{Key: Key, ExpiresIn: expires, Cache: c}
}

// GetCacheValue Fetching the cache token value -.
func (p *CacheParams) GetCacheValue() (string, error) {

	// Get a value from the cache -.
	value, found := p.Cache.Get(p.Key)

	if found {
		return value.(string), nil
	} else {
		return "", fmt.Errorf("key not found in cache")
	}
}

// SetCacheValue Set the cache value -.
func (p *CacheParams) SetCacheValue(cacheValue string) error {
	// Set a value in the cache -.
	p.Cache.Set(p.Key, cacheValue, p.ExpiresIn)

	return nil
}
