package exchangeService

import (
	"simplebank/db/util"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRequestExchangeRate(t *testing.T) {

	randomAmount := float64(util.RandomInt(1, 1000))
	// t.Run("OK", func(t *testing.T) {
	// 	controller := gomock.NewController(t)
	// 	defer controller.Finish()

	// 	store := mockdb.NewMockStore(controller)
	// 	server := newTestServer(t, store)

	// 	from := util.EUR
	// 	to := util.HKD
	// 	exchangeResponse, err := requestExchange(server.config, exchangeRequest{
	// 		FromCurrency: from,
	// 		ToCurrency:   to,
	// 		Amount:       randomAmount,
	// 	})
	// 	require.NoError(t, err)
	// 	require.Equal(t, exchangeResponse.From, from)
	// 	require.Equal(t, exchangeResponse.To, to)
	// 	require.True(t, exchangeResponse.Rate == float64(exchangeResponse.Rate))
	// })

	t.Run("api not provided", func(t *testing.T) {
		config := util.Config{
			ApiKey: "fake",
		}
		service := NewExchangeService(config)
		from := util.EUR
		to := util.HKD
		_, err := service.RequestExchange(ExchangeApiRequest{
			FromCurrency: from,
			ToCurrency:   to,
			Amount:       randomAmount,
		})

		require.ErrorContains(t, err, "api error - status code larger than 400")
	})
}
