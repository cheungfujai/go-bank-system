package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"simplebank/db/util"

	"github.com/gin-gonic/gin"
)

const (
	URL = "https://api.apilayer.com/exchangerates_data/convert?to="
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

type exchangeRequest struct {
	ToCurrency   string `form:"to" binding:"required"`
	FromCurrency string `form:"from" binding:"required"`
	Amount       string `form:"amount" binding:"required"`
}

type exchangeResponse struct {
	From string  `json:"from"`
	To   string  `json:"to"`
	Rate float64 `json:"rate"`
}

func (server *Server) getExchangeRate(ctx *gin.Context) {
	var exchangeReq exchangeRequest
	if err := ctx.ShouldBindQuery(&exchangeReq); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	exchangeRate, err := server.cache.GetExchangeCache(exchangeReq.FromCurrency, exchangeReq.ToCurrency)
	if err == nil {
		resp := exchangeResponse{
			From: exchangeReq.FromCurrency,
			To:   exchangeReq.ToCurrency,
			Rate: exchangeRate,
		}
		ctx.JSON(http.StatusOK, resp)
	}

	response, err := requestExchange(server.config, exchangeReq)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
	}
	server.cache.SetExchangeCache(exchangeReq.FromCurrency, exchangeReq.ToCurrency, response.Rate)

	ctx.JSON(http.StatusOK, response)
}

func requestExchange(config util.Config, exchangeReq exchangeRequest) (exchangeResponse, error) {

	url := buildApiUrl(exchangeReq.Amount, exchangeReq.ToCurrency, exchangeReq.FromCurrency)
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)

	req.Header.Set("apikey", config.ApiKey)

	if err != nil {
		return exchangeResponse{}, fmt.Errorf("api error %v", err)
	}

	res, err := client.Do(req)
	if err != nil {
		return exchangeResponse{}, fmt.Errorf("api error %v", err)
	}

	if res.Body != nil {
		defer res.Body.Close()
	}

	body, err := ioutil.ReadAll(res.Body)
	var responseObject apiResponse
	json.Unmarshal(body, &responseObject)
	if err != nil {
		return exchangeResponse{}, fmt.Errorf("api error - cannot ummarshal json %v", err)
	}
	if res.StatusCode >= 400 {
		return exchangeResponse{}, fmt.Errorf("api error - status code larger than 400 with body %v", responseObject)
	}

	log.Printf("%v" + string(body))
	return exchangeResponse{
		From: exchangeReq.FromCurrency,
		To:   exchangeReq.ToCurrency,
		Rate: responseObject.ExchangeInfo.Rate,
	}, nil
}

func buildApiUrl(amount string, toCurrency string, fromCurrency string) string {
	result := URL + toCurrency + "&from=" +
		fromCurrency + "&amount=" + amount

	return result
}
