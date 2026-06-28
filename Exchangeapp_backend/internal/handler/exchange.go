package handler

import (
	"exchangeapp/internal/model"
	"exchangeapp/internal/service"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ExchangeHandler struct {
	exchangeSvc *service.ExchangeRateService
}

func NewExchangeHandler(exchangeSvc *service.ExchangeRateService) *ExchangeHandler {
	return &ExchangeHandler{exchangeSvc: exchangeSvc}
}

func (h *ExchangeHandler) Create(ctx *gin.Context) {
	var rate model.ExchangeRate
	if err := ctx.ShouldBindJSON(&rate); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.exchangeSvc.CreateExchangeRate(&rate); err != nil {
		log.Printf("CreateExchangeRate error: %v", err)
		genericError(ctx, http.StatusInternalServerError, "failed to create exchange rate")
		return
	}

	ctx.JSON(http.StatusCreated, rate)
}

func (h *ExchangeHandler) GetAll(ctx *gin.Context) {
	rates, err := h.exchangeSvc.GetExchangeRates()
	if err != nil {
		log.Printf("GetExchangeRates error: %v", err)
		genericError(ctx, http.StatusInternalServerError, "failed to fetch exchange rates")
		return
	}

	ctx.JSON(http.StatusOK, rates)
}
