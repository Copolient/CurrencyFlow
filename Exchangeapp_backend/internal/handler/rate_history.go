package handler

import (
	"exchangeapp/internal/service"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type RateHistoryHandler struct {
	rateHistorySvc *service.RateHistoryService
}

func NewRateHistoryHandler(rateHistorySvc *service.RateHistoryService) *RateHistoryHandler {
	return &RateHistoryHandler{rateHistorySvc: rateHistorySvc}
}

func (h *RateHistoryHandler) GetHistory(ctx *gin.Context) {
	from := ctx.Query("from")
	to := ctx.Query("to")
	rangeStr := ctx.DefaultQuery("range", "1M")

	if from == "" || to == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "from and to query params are required"})
		return
	}

	fromCode, valid := sanitizeCurrencyCode(from)
	if !valid {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid from currency code"})
		return
	}
	toCode, valid := sanitizeCurrencyCode(to)
	if !valid {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid to currency code"})
		return
	}

	histories, err := h.rateHistorySvc.GetHistoryByPair(ctx.Request.Context(), fromCode, toCode, rangeStr)
	if err != nil {
		log.Printf("GetHistoryByPair error: %v", err)
		genericError(ctx, http.StatusInternalServerError, "failed to fetch rate history")
		return
	}

	ctx.JSON(http.StatusOK, histories)
}

func (h *RateHistoryHandler) GetLatest(ctx *gin.Context) {
	histories, err := h.rateHistorySvc.GetLatestByAllPairs(ctx.Request.Context())
	if err != nil {
		log.Printf("GetLatestByAllPairs error: %v", err)
		genericError(ctx, http.StatusInternalServerError, "failed to fetch latest rates")
		return
	}

	ctx.JSON(http.StatusOK, histories)
}
