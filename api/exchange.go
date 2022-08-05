package api

import (
	"fmt"
	"log"
	"net/http"
	"simplebank/db/util"
	exchangeService "simplebank/service"

	"github.com/gin-gonic/gin"
)

func (server *Server) getExchangeRate(ctx *gin.Context) {
	var exchangeReq exchangeService.ExchangeApiRequest
	if err := ctx.ShouldBindQuery(&exchangeReq); err != nil {
		fmt.Println(err)
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	log.Println("getting exchange rate from cache")
	exchangeRate, err := server.cache.GetExchangeCache(exchangeReq.FromCurrency, exchangeReq.ToCurrency)
	fmt.Println(err)

	if err == nil {
		converted, err := util.RoundFloatAtDecimal(exchangeReq.Amount*exchangeRate, 3)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}
		resp := exchangeService.ExchangeApiResponse{
			From:      exchangeReq.FromCurrency,
			To:        exchangeReq.ToCurrency,
			Rate:      exchangeRate,
			Converted: converted,
		}
		ctx.JSON(http.StatusOK, resp)
		return
	}

	log.Println("exchange rate not exist in cache, calling from api ...")
	response, err := server.exchangeService.RequestExchange(exchangeReq)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	server.cache.SetExchangeCache(exchangeReq.FromCurrency, exchangeReq.ToCurrency, response.Rate)
	ctx.JSON(http.StatusOK, response)
}
