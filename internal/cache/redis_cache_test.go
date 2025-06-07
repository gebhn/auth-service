package cache

import (
	"context"
	"log"
	"os"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	"github.com/stretchr/testify/assert"
)

type cacheMock struct {
	c *redisCache
	m *miniredis.Miniredis
}

var testCache *cacheMock

func TestMain(m *testing.M) {
	mr, err := miniredis.Run()
	if err != nil {
		log.Fatal(err)
	}
	defer mr.Close()

	testCache = &cacheMock{
		c: NewRedisCache(mr.Addr(), ""),
		m: mr,
	}

	os.Exit(m.Run())
}

func TestSet_Success(t *testing.T) {
	testCache.m.FastForward(time.Hour * 24)

	_, err := testCache.c.Set(context.Background(), "key", "value", time.Second*5)
	assert.NoError(t, err)
}

func TestSet_Fail(t *testing.T) {
	testCache.m.FastForward(time.Hour * 24)
	testCache.m.SetError("err")
	defer testCache.m.SetError("")

	_, err := testCache.c.Set(context.Background(), "key", "value", time.Second*5)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "err")
}

func TestSet_Invalid(t *testing.T) {
	testCache.m.FastForward(time.Hour * 24)
	tc := []struct {
		key   string
		value string
		exp   time.Duration
		label string
	}{
		{
			key:   "",
			value: "value",
			exp:   time.Second * 5,
			label: "Missing Key",
		},
		{
			key:   "key",
			value: "",
			exp:   time.Second * 5,
			label: "Missing Value",
		},
		{
			key:   "key",
			value: "value",
			exp:   -(time.Second * 5),
			label: "Invalid Expiration",
		},
	}

	for _, c := range tc {
		t.Run(c.label, func(t *testing.T) {
			_, err := testCache.c.Set(context.Background(), c.key, c.value, c.exp)
			assert.Error(t, err)
			assert.Contains(t, err.Error(), ErrInvalidInput.Error())
		})
	}
}

func TestGet_Success(t *testing.T) {
	testCache.m.FastForward(time.Hour * 24)

	var err error
	var v string

	_, err = testCache.c.Set(context.Background(), "key1", "value1", time.Second*5)
	assert.NoError(t, err)

	v, err = testCache.c.Get(context.Background(), "key1")
	assert.NoError(t, err)
	assert.NotEmpty(t, v)
	assert.Contains(t, v, "value1")
}

func TestGet_Fail(t *testing.T) {
	testCache.m.FastForward(time.Hour * 24)

	var err error
	var v string

	_, err = testCache.c.Set(context.Background(), "key1", "value1", time.Second*5)
	assert.NoError(t, err)

	testCache.m.SetError("err")
	defer testCache.m.SetError("")

	v, err = testCache.c.Get(context.Background(), "key1")
	assert.Error(t, err)
	assert.Empty(t, v)
	assert.Contains(t, err.Error(), "err")
}

func TestGet_Invalid(t *testing.T) {
	testCache.m.FastForward(time.Hour * 24)

	var err error

	_, err = testCache.c.Set(context.Background(), "key1", "value1", time.Second*5)
	assert.NoError(t, err)

	_, err = testCache.c.Get(context.Background(), "")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), ErrInvalidInput.Error())
}
