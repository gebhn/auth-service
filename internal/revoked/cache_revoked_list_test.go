package revoked

import (
	"context"
	"log"
	"os"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	"github.com/gebhn/auth-service/api/pb"
	"github.com/gebhn/auth-service/internal/cache"
	"github.com/gebhn/auth-service/internal/config"
	"github.com/stretchr/testify/assert"
)

type mockList struct {
	c *cacheRevokedList
	m *miniredis.Miniredis
}

var testList *mockList

var globalContext context.Context

func TestMain(m *testing.M) {
	var cancelFunc context.CancelFunc

	mr, err := miniredis.Run()
	if err != nil {
		log.Fatal(err)
	}
	defer mr.Close()

	cache := cache.NewRedisCache(mr.Addr(), "")
	defer cache.Close()

	testList = &mockList{
		c: NewCacheRevokedList(cache),
		m: mr,
	}

	globalContext, cancelFunc = context.WithCancel(context.Background())
	defer cancelFunc()

	os.Exit(m.Run())
}

func TestCreate_Success(t *testing.T) {
	defer testList.m.FlushAll()

	err := testList.c.Create(globalContext, "testJti", pb.TokenKind_TOKEN_KIND_ACCESS, config.GetTokenDuration(pb.TokenKind_TOKEN_KIND_ACCESS))
	assert.NoError(t, err)
}

func TestCreate_Invalid(t *testing.T) {
	defer testList.m.FlushAll()

	t.Run("Invalid JTI", func(t *testing.T) {
		err := testList.c.Create(globalContext, "", pb.TokenKind_TOKEN_KIND_ACCESS, config.GetTokenDuration(pb.TokenKind_TOKEN_KIND_ACCESS))
		assert.Error(t, err)
		assert.Contains(t, err.Error(), ErrInvalidKey.Error())
	})
	t.Run("Invalid Duration", func(t *testing.T) {
		err := testList.c.Create(globalContext, "testJti", pb.TokenKind_TOKEN_KIND_ACCESS, time.Second)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), ErrInvalidDuration.Error())
	})
}

func TestCreate_Fail(t *testing.T) {
	defer testList.m.FlushAll()
	testList.m.SetError("failed to create")
	defer testList.m.SetError("")

	err := testList.c.Create(globalContext, "testJti", pb.TokenKind_TOKEN_KIND_ACCESS, config.GetTokenDuration(pb.TokenKind_TOKEN_KIND_ACCESS))
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to create")
}

func TestFind_Success(t *testing.T) {
	defer testList.m.FlushAll()

	err := testList.c.Create(globalContext, "testJti", pb.TokenKind_TOKEN_KIND_ACCESS, config.GetTokenDuration(pb.TokenKind_TOKEN_KIND_ACCESS))
	assert.NoError(t, err)

	ok, err := testList.c.Find(globalContext, "testJti")
	assert.NoError(t, err)
	assert.True(t, ok)
}

func TestFind_Invalid(t *testing.T) {
	defer testList.m.FlushAll()

	err := testList.c.Create(globalContext, "testJti", pb.TokenKind_TOKEN_KIND_ACCESS, config.GetTokenDuration(pb.TokenKind_TOKEN_KIND_ACCESS))
	assert.NoError(t, err)

	ok, err := testList.c.Find(globalContext, "")
	assert.Error(t, err)
	assert.False(t, ok)
	assert.Contains(t, err.Error(), ErrInvalidKey.Error())
}

func TestFind_Fail(t *testing.T) {
	defer testList.m.FlushAll()
	testList.m.SetError("failed to find")
	defer testList.m.SetError("")

	ok, err := testList.c.Find(globalContext, "testJti")
	assert.Error(t, err)
	assert.False(t, ok)
	assert.Contains(t, err.Error(), "failed to find")
}
