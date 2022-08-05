package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	mockdb "simplebank/db/mock"
	"simplebank/db/util"
	exchangeService "simplebank/service"
	"simplebank/token"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

type SpyExchangeService struct {
}

func (service *SpyExchangeService) RequestExchange(exchangeReq exchangeService.ExchangeApiRequest) (exchangeService.ExchangeApiResponse, error) {
	rate := util.RandomExchangeRate()
	converted, err := util.RoundFloatAtDecimal(exchangeReq.Amount*rate, 3)
	if err != nil {
		return exchangeService.ExchangeApiResponse{}, err
	}
	return exchangeService.ExchangeApiResponse{
		From:      exchangeReq.FromCurrency,
		To:        exchangeReq.ToCurrency,
		Rate:      rate,
		Converted: converted,
	}, nil
}

func buildUrl(request exchangeService.ExchangeApiRequest) string {
	return fmt.Sprintf(
		"/exchange?from=%s&to=%s&amount=%f",
		request.FromCurrency, request.ToCurrency, request.Amount)
}

func getResponse(recorder *httptest.ResponseRecorder) exchangeService.ExchangeApiResponse {
	respBody, _ := ioutil.ReadAll(recorder.Body)
	var body exchangeService.ExchangeApiResponse
	json.Unmarshal(respBody, &body)
	return body
}
func TestRequestExchangeRate(t *testing.T) {
	user, _ := randomUser(t)

	testCases := []struct {
		name          string
		requestParam  exchangeService.ExchangeApiRequest
		setupAuth     func(t *testing.T, request *http.Request, tokenMaker token.Maker)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "Happy Case",
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuth(t, request, tokenMaker, authTypeBearer, user.Username, time.Minute)
			},
			requestParam: exchangeService.ExchangeApiRequest{
				FromCurrency: util.USD,
				ToCurrency:   util.EUR,
				Amount:       util.RandomFloat64(),
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name: "Bad Request",
			requestParam: exchangeService.ExchangeApiRequest{
				FromCurrency: util.USD,
				ToCurrency:   util.EUR,
				Amount:       0.0,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuth(t, request, tokenMaker, authTypeBearer, user.Username, time.Minute)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
	}

	for i := range testCases {
		currTestCase := testCases[i]

		t.Run(currTestCase.name, func(t *testing.T) {
			mockController := gomock.NewController(t)
			defer mockController.Finish()

			store := mockdb.NewMockStore(mockController)

			server := newTestServer(t, store)
			recorder := httptest.NewRecorder()

			url := buildUrl(currTestCase.requestParam)
			request, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			currTestCase.setupAuth(t, request, server.tokenMaker)

			server.router.ServeHTTP(recorder, request)
			currTestCase.checkResponse(t, recorder)

		})
	}

	t.Run("Should first request same as second request", func(t *testing.T) {
		mockController := gomock.NewController(t)
		defer mockController.Finish()

		store := mockdb.NewMockStore(mockController)

		server := newTestServer(t, store)
		recorder := httptest.NewRecorder()

		url := buildUrl(testCases[0].requestParam)
		request, err := http.NewRequest(http.MethodGet, url, nil)
		require.NoError(t, err)

		testCases[0].setupAuth(t, request, server.tokenMaker)

		server.router.ServeHTTP(recorder, request)
		require.Equal(t, http.StatusOK, recorder.Code)
		// first request from external api
		firstResponseBody := getResponse(recorder)

		request, err = http.NewRequest(http.MethodGet, url, nil)
		require.NoError(t, err)
		server.router.ServeHTTP(recorder, request)
		require.Equal(t, http.StatusOK, recorder.Code)
		// second request from cache
		secondResponseBody := getResponse(recorder)

		require.Equal(t, firstResponseBody.From, secondResponseBody.From)
		require.Equal(t, firstResponseBody.To, secondResponseBody.To)
		require.Equal(t, firstResponseBody.Rate, secondResponseBody.Rate)
		require.Equal(t, firstResponseBody.Converted, secondResponseBody.Converted)
	})
}
