package cache

import (
	"sync"
	"time"
)

// Cache structure contains cached records necessary to perform a distributed cache lookup or entry.
type Cache[T any] struct {
	cache *record[T]
	sync.Mutex
}

// record is a structure that contains the cached values for lookup, cache refreshing, and
// error handling.
type record[T any] struct {
	// value is represents the item to be cached.
	value T

	// createdAt is the time that the cache was created.
	createdAt time.Time

	// cacheFor is the time duration that a cache record will remain in the cache. If the cached record has expired
	// the cache value will only be refreshed on the next call of Get.
	cacheFor time.Duration

	err   error
	ready chan struct{}
}

// Fn type is to receive any function that returns ([]byte, error).
type Fn[T any] func() (T, error)

// Get performs the request on the supplied fn and caches if necessary.
// Get can be used to populate the cache for the first time if the cache record
// doesn't exist yet. If the expiry time has elapsed, a new record will be looked up and the
// cache will be refreshed.
func (c *Cache[T]) Get(fn Fn[T]) (T, error) {
	c.Lock()
	r := c.cache
	elapsed := time.Since(r.createdAt)
	if elapsed > r.cacheFor {
		r = &record[T]{
			ready:    make(chan struct{}),
			cacheFor: r.cacheFor,
		}
		r.createdAt = time.Now()

		c.cache = r
		c.Unlock()

		// Avoid blocking thread if runtime panics.
		defer func() {
			if rec := recover(); rec != nil {
				close(r.ready)
			}
		}()

		r.value, r.err = fn()
		close(r.ready)
	} else {
		c.Unlock()
		<-r.ready
	}

	if r.err != nil {
		return r.value, r.err
	}

	return r.value, nil
}

// Clear clears the cache value and resets the expiry time.
func (c *Cache[T]) Clear() {
	c.Lock()
	r := c.cache
	r.value = *new(T)
	r.cacheFor = 0
	c.Unlock()
}

// New returns an empty Cache type.
func New[T any](cacheFor time.Duration) *Cache[T] {
	return &Cache[T]{
		cache: &record[T]{
			cacheFor: cacheFor,
		},
	}
}

// NewForTesting offers a test-specific cache stub that can be used when the cache is not relevant to the testing scenario
func NewForTesting[T any]() *Cache[T] {
	return &Cache[T]{
		cache: &record[T]{
			cacheFor: 0,
		},
	}
}
