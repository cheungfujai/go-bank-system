package bankcache

import (
	"fmt"
	"math"
	"simplebank/db/util"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCreateExchangeMap(t *testing.T) {
	bankCache := NewCache()
	bankCache.createExchangeMap()

	x, found := bankCache.Cache.Get(ExchangeRateCacheKey)
	require.True(t, found)
	require.NotNil(t, x.(*ExchangeMap))
}

func TestGetExchangeCacheNotFound(t *testing.T) {
	bankCache := NewCache()
	bankCache.createExchangeMap()

	from := util.EUR
	to := util.HKD

	_, err := bankCache.GetExchangeCache(from, to)
	require.Error(t, err, "map with key "+from+"not found")
}

func TestGetSetExchangeCache(t *testing.T) {
	bankCache := NewCache()
	bankCache.createExchangeMap()

	from := util.EUR
	to := util.HKD
	rate := 8.123
	bankCache.SetExchangeCache(from, to, rate)

	resultRate, err := bankCache.GetExchangeCache(from, to)
	require.NoError(t, err)
	require.Equal(t, rate, resultRate)

	resultRate, err = bankCache.GetExchangeCache(to, from)
	require.NoError(t, err)
	fmt.Println(resultRate)
	require.Equal(t, resultRate, math.Round(1/rate*1000)/1000)
}
