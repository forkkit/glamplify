package cache

import (
	cachego "github.com/patrickmn/go-cache"
	"time"
)

type Cache struct {
	cache *cachego.Cache
}

func New() *Cache {
	c := cachego.New(5*time.Minute, 1*time.Minute)

	return &Cache {
		cache: c,
	}
}

func (c Cache) Get(key string) (interface{}, bool) {
	return c.cache.Get(key)
}

func (c *Cache) Set(key string, val interface{}, duration time.Duration) {
	c.cache.Set(key, val, duration)
}

