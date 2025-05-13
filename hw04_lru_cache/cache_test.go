package hw04lrucache

import (
	"math/rand"
	"strconv"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCache(t *testing.T) {
	t.Run("empty cache", func(t *testing.T) {
		c := NewCache(10)

		_, ok := c.Get("aaa")
		require.False(t, ok)

		_, ok = c.Get("bbb")
		require.False(t, ok)
	})

	t.Run("simple", func(t *testing.T) {
		c := NewCache(5)

		wasInCache := c.Set("aaa", 100)
		require.False(t, wasInCache)

		wasInCache = c.Set("bbb", 200)
		require.False(t, wasInCache)

		val, ok := c.Get("aaa")
		require.True(t, ok)
		require.Equal(t, 100, val)

		val, ok = c.Get("bbb")
		require.True(t, ok)
		require.Equal(t, 200, val)

		wasInCache = c.Set("aaa", 300)
		require.True(t, wasInCache)

		val, ok = c.Get("aaa")
		require.True(t, ok)
		require.Equal(t, 300, val)

		val, ok = c.Get("ccc")
		require.False(t, ok)
		require.Nil(t, val)
	})

	t.Run("purge logic", func(t *testing.T) {
		c := NewCache(3)

		c.Set("aaa", 100)
		c.Set("bbb", 200)
		c.Set("ccc", 300)
		c.Set("ddd", 400)

		c.Clear()

		var lc = c.(*lruCache)
		require.True(t, len(lc.items) == 0)
		require.True(t, lc.queue.Len() == 0)
		require.True(t, lc.capacity == 3)

	})
	t.Run("getting a place in a place where there is no place", func(t *testing.T) {
		c := NewCache(3)

		c.Set("aaa", 100)
		c.Set("bbb", 200)
		c.Set("ccc", 300)
		c.Set("ddd", 400)

		require.True(t, len(c.(*lruCache).items) == 3)

		val, isPresent := c.Get("aaa")
		require.True(t, isPresent == false)
		require.True(t, val == nil)

		val, isPresent = c.Get("bbb")
		require.True(t, isPresent == true)
		require.True(t, val != nil)
		require.True(t, val == 200)

		val, isPresent = c.Get("ccc")
		require.True(t, isPresent == true)
		require.True(t, val != nil)
		require.True(t, val == 300)

		val, isPresent = c.Get("ddd")
		require.True(t, isPresent == true)
		require.True(t, val != nil)
		require.True(t, val == 400)
	})

	t.Run("getting out most old element", func(t *testing.T) {
		c := NewCache(3)

		c.Set("aaa", 100)
		c.Set("bbb", 200)
		c.Set("ccc", 300)

		c.Get("aaa")
		c.Get("bbb")
		c.Get("ccc")
		c.Get("bbb")
		c.Set("ccc", 333)

		c.Set("ddd", 400)

		require.Equal(t, 3, c.(*lruCache).queue.Len())
		require.Equal(t, 3, len(c.(*lruCache).items))

		val, isPresent := c.Get("aaa")
		require.False(t, isPresent)
		require.True(t, val == nil)

		val, isPresent = c.Get("bbb")
		require.True(t, isPresent == true)
		require.True(t, val != nil)
		require.True(t, val == 200)

		val, isPresent = c.Get("ccc")
		require.True(t, isPresent == true)
		require.True(t, val != nil)
		require.True(t, val == 333)

		val, isPresent = c.Get("ddd")
		require.True(t, isPresent == true)
		require.True(t, val != nil)
		require.True(t, val == 400)

	})

}

func TestCacheMultithreading(t *testing.T) {
	t.Skip() // Remove me if task with asterisk completed.

	c := NewCache(10)
	wg := &sync.WaitGroup{}
	wg.Add(2)

	go func() {
		defer wg.Done()
		for i := range 1_000_000 {
			c.Set(Key(strconv.Itoa(i)), i)
		}
	}()

	go func() {
		defer wg.Done()
		for range 1_000_000 {
			c.Get(Key(strconv.Itoa(rand.Intn(1_000_000))))
		}
	}()

	wg.Wait()
}
