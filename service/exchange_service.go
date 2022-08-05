package exchangeService

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"simplebank/db/util"
)

const (
	URL = "https://api.apilayer.com/exchangerates_data/convert"
)

type apiExchangeInfo struct {
	Timestamp int     `json:"timestamp"`
	Rate      float64 `json:"rate"`
}
type apiResponse struct {
	Status       bool            `json:"success"`
	Result       float64         `json:"result"`
	ExchangeInfo apiExchangeInfo `json:"info"`
	Error        apiError
}

type apiError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type ExchangeApiRequest struct {
	ToCurrency   string  `form:"to" binding:"required"`
	FromCurrency string  `form:"from" binding:"required"`
	Amount       float64 `form:"amount" binding:"required"`
}

type ExchangeApiResponse struct {
	From      string  `json:"from"`
	To        string  `json:"to"`
	Rate      float64 `json:"rate"`
	Converted float64 `json:"converted"`
}

type ExchangeServiceInterface interface {
	RequestExchange(exchangeReq ExchangeApiRequest) (ExchangeApiResponse, error)
}

type ExchangeService struct {
	config util.Config
}

func NewExchangeService(config util.Config) *ExchangeService {
	return &ExchangeService{
		config: config,
	}
}

func (service *ExchangeService) RequestExchange(exchangeReq ExchangeApiRequest) (ExchangeApiResponse, error) {

	url := buildApiUrl(exchangeReq.Amount, exchangeReq.ToCurrency, exchangeReq.FromCurrency)
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)

	req.Header.Set("apikey", service.config.ApiKey)

	if err != nil {
		return ExchangeApiResponse{}, fmt.Errorf("api error %v", err)
	}

	res, err := client.Do(req)
	if err != nil {
		return ExchangeApiResponse{}, fmt.Errorf("api error %v", err)
	}

	if res.Body != nil {
		defer res.Body.Close()
	}

	body, err := ioutil.ReadAll(res.Body)
	var responseObject apiResponse
	json.Unmarshal(body, &responseObject)
	if err != nil {
		return ExchangeApiResponse{}, fmt.Errorf("api error - cannot ummarshal json %v", err)
	}
	if res.StatusCode >= 400 {
		return ExchangeApiResponse{}, fmt.Errorf("api error - status code larger than 400 with body %v", string(body))
	}

	log.Printf("%v" + string(body))
	converted, err := util.RoundFloatAtDecimal(responseObject.Result, 3)
	if err != nil {
		return ExchangeApiResponse{}, err
	}
	return ExchangeApiResponse{
		From:      exchangeReq.FromCurrency,
		To:        exchangeReq.ToCurrency,
		Rate:      responseObject.ExchangeInfo.Rate,
		Converted: converted,
	}, nil
}

func buildApiUrl(amount float64, toCurrency string, fromCurrency string) string {
	result := URL + "?to=" + toCurrency + "&from=" +
		fromCurrency + "&amount=" + fmt.Sprintf("%f", amount)
	return result
}
