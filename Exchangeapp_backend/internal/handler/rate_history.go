package handler

import (
	"exchangeapp/internal/service"
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

	histories, err := h.rateHistorySvc.GetHistoryByPair(ctx.Request.Context(), from, to, rangeStr)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, histories)
}

func (h *RateHistoryHandler) GetLatest(ctx *gin.Context) {
	histories, err := h.rateHistorySvc.GetLatestByAllPairs(ctx.Request.Context())
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, histories)
}
