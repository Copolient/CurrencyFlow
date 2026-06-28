package handler

import (
	"exchangeapp/internal/service"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AIAnalystHandler struct {
	aiSvc *service.AIAnalystService
}

func NewAIAnalystHandler(aiSvc *service.AIAnalystService) *AIAnalystHandler {
	return &AIAnalystHandler{aiSvc: aiSvc}
}

type analyzeRequest struct {
	From     string `json:"from" binding:"required,len=3"`
	To       string `json:"to" binding:"required,len=3"`
	Question string `json:"question"`
}

func (h *AIAnalystHandler) Analyze(c *gin.Context) {
	var req analyzeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := h.aiSvc.Analyze(c.Request.Context(), req.From, req.To, req.Question)
	if err != nil {
		log.Printf("AIAnalyze error: %v", err)
		genericError(c, http.StatusInternalServerError, "analysis failed")
		return
	}

	c.JSON(http.StatusOK, result)
}
