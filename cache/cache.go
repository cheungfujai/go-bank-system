package bankcache

import (
	"fmt"
	"math"
	"time"

	"github.com/patrickmn/go-cache"
)

const (
	ExchangeRateCacheKey = "ExchangeRate"
)

type BankCache struct {
	Cache *cache.Cache
}

// type BankCacheInterface interface {
// 	CreateExchangeMap()
// 	SetExchangeCache(pair ExchangeRatePair)
// 	GetExchangeCache(From string, To string) (float64, error)
// }

func NewCache() BankCache {
	bankCache := BankCache{
		Cache: cache.New(8*60*time.Minute, 8*60*time.Minute),
	}
	bankCache.createExchangeMap()
	return bankCache
}

type MyStruct struct {
	Value  string `json:"value"`
	Value2 int    `json:"value2"`
}

type ExchangeMap struct {
	ExchangeMap map[string]map[string]float64
}

func (bankCache *BankCache) createExchangeMap() {
	exchangeMap := ExchangeMap{
		ExchangeMap: make(map[string]map[string]float64),
	}
	bankCache.Cache.Set(ExchangeRateCacheKey, &exchangeMap, cache.NoExpiration)
}

func (bankCache *BankCache) SetExchangeCache(From string, To string, ExchangeRate float64) {
	if x, found := bankCache.Cache.Get(ExchangeRateCacheKey); found {
		exchangeMap := x.(*ExchangeMap)

		exchangeMap.ExchangeMap[From][To] = ExchangeRate
		exchangeMap.ExchangeMap[To][From] = math.Round(1/ExchangeRate*1000) / 1000
	}
}

func (bankCache *BankCache) GetExchangeCache(From string, To string) (float64, error) {
	var exchangeMap *ExchangeMap
	if x, found := bankCache.Cache.Get(ExchangeRateCacheKey); found {
		exchangeMap = x.(*ExchangeMap)
		return 0, fmt.Errorf("no cache with name " + ExchangeRateCacheKey + " exist")
	}

	if fromMap, found := exchangeMap.ExchangeMap[From]; found {
		return fromMap[To], nil
	}
	return 0, fmt.Errorf("")
}
