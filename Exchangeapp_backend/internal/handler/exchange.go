package handler

import (
	"exchangeapp/internal/model"
	"exchangeapp/internal/service"
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
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, rate)
}

func (h *ExchangeHandler) GetAll(ctx *gin.Context) {
	rates, err := h.exchangeSvc.GetExchangeRates()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, rates)
}
