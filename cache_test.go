package cache_test

import (
	"testing"
	"time"

	"github.com/peterlabuschagne/cache"
	"github.com/stretchr/testify/require"
)

type mock struct {
	val int
}

func TestCache_Get(t *testing.T) {
	expireAfter := time.Second * 1
	c := cache.New[mock](expireAfter)

	val1 := 123
	t.Run("empty cache - insert new record and lookup", func(t *testing.T) {
		m, err := c.Get(func() (mock, error) {
			return mock{val: val1}, nil
		})
		require.Nil(t, err)
		require.Equal(t, val1, m.val)
	})

	val2 := 456
	t.Run("cache exists - record not expired yet", func(t *testing.T) {
		m, err := c.Get(func() (mock, error) {
			return mock{val: val2}, nil
		})
		require.Nil(t, err)
		require.Equal(t, val1, m.val)
	})

	time.Sleep(time.Second * 1)
	val3 := 567
	t.Run("cache exists - refresh as record is expired", func(t *testing.T) {
		m, err := c.Get(func() (mock, error) {
			return mock{val: val3}, nil
		})
		require.Nil(t, err)
		require.Equal(t, val3, m.val)
	})
}

func TestCache_Clear(t *testing.T) {
	expireAfter := time.Second * 1
	c := cache.New[mock](expireAfter)

	t.Run("clear cache - record should be empty", func(t *testing.T) {
		val := 123
		m, err := c.Get(func() (mock, error) {
			return mock{val: val}, nil
		})
		require.Nil(t, err)
		require.Equal(t, val, m.val)

		c.Clear()
		m, err = c.Get(func() (mock, error) {
			return mock{}, nil
		})
		require.Nil(t, err)
		require.Equal(t, 0, m.val)
	})
}
