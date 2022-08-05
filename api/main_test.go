package api

import (
	"os"
	bankcache "simplebank/cache"
	db "simplebank/db/sqlc"
	"simplebank/db/util"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

func newTestServer(t *testing.T, store db.Store) *Server {
	config, err := util.LoadViberConfig("../")
	require.NoError(t, err)
	cache := bankcache.NewCache()
	exchangeService := &SpyExchangeService{}
	server, err := NewServer(config, store, cache, exchangeService)
	require.NoError(t, err)

	return server
}

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)
	os.Exit(m.Run())
}
