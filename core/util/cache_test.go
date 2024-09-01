package util_test

import (
	"context"
	"testing"
	"time"

	"github.com/medama-io/medama/util"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func SetupCacheTest(t *testing.T) (*assert.Assertions, *require.Assertions, context.Context) {
	t.Helper()
	assert := assert.New(t)
	require := require.New(t)
	ctx := context.Background()
	return assert, require, ctx
}

func TestGetSet(t *testing.T) {
	_, require, ctx := SetupCacheTest(t)

	cycle := 100 * time.Millisecond
	c := util.NewCache(ctx, cycle)
	defer c.Close()

	c.Set("sticky", "forever", cycle*2)
	c.Set("hello", "Hello", cycle/2)
	hello, err := c.Get(ctx, "hello")
	require.NoError(err, "Cache miss")

	if hello.(string) != "Hello" {
		t.Log("Cache value mismatch")
		t.FailNow()
	}

	time.Sleep(cycle / 2)

	_, err = c.Get(ctx, "hello")
	require.Error(err, "Cache value not expired")

	time.Sleep(cycle)

	_, err = c.Get(ctx, "404")
	require.Error(err, "Cache value not expired")

	_, err = c.Get(ctx, "sticky")
	require.NoError(err, "Cache value not found")
}

func TestHas(t *testing.T) {
	assert, require, ctx := SetupCacheTest(t)
	c := util.NewCache(ctx, time.Minute)

	c.Set("hello", "Hello", time.Hour)
	ok, err := c.Has(ctx, "hello")
	assert.True(ok, "Cache miss")
	require.NoError(err, "Cache error")

	ok, err = c.Has(ctx, "404")
	assert.False(ok, "Cache hit")
	require.NoError(err, "Cache error")
}

func TestDelete(t *testing.T) {
	_, require, ctx := SetupCacheTest(t)
	c := util.NewCache(ctx, time.Minute)
	c.Set("hello", "Hello", time.Hour)
	_, err := c.Get(ctx, "hello")
	require.NoError(err, "Cache miss")

	c.Delete("hello")

	_, err = c.Get(ctx, "hello")
	require.Error(err, "Cache value not deleted")
}

func TestRange(t *testing.T) {
	assert, _, ctx := SetupCacheTest(t)
	c := util.NewCache(ctx, time.Minute)
	c.Set("hello", "Hello", time.Hour)
	c.Set("world", "World", time.Hour)
	count := 0

	c.Range(ctx, func(_key, _value interface{}) bool {
		count++
		return true
	})

	assert.Equal(2, count, "Cache range mismatch")
}

func TestRangeTimer(t *testing.T) {
	_, _, ctx := SetupCacheTest(t)
	c := util.NewCache(ctx, time.Minute)
	c.Set("message", "Hello", time.Nanosecond)
	c.Set("world", "World", time.Nanosecond)
	time.Sleep(time.Microsecond)

	c.Range(ctx, func(_key, _value interface{}) bool {
		t.Log("Cache range mismatch")
		t.FailNow()
		return true
	})
}
