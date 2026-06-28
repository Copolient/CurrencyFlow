package handler

import (
	"exchangeapp/internal/service"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type AlertHandler struct {
	alertSvc *service.AlertService
}

func NewAlertHandler(alertSvc *service.AlertService) *AlertHandler {
	return &AlertHandler{alertSvc: alertSvc}
}

type alertRequest struct {
	FromCurrency string  `json:"fromCurrency" binding:"required,len=3"`
	ToCurrency   string  `json:"toCurrency" binding:"required,len=3"`
	TargetRate   float64 `json:"targetRate" binding:"required,gt=0"`
	Direction    string  `json:"direction" binding:"required,oneof=above below"`
}

func (h *AlertHandler) Create(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		return
	}

	var req alertRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.alertSvc.CreateAlert(userID, req.FromCurrency, req.ToCurrency, req.TargetRate, req.Direction); err != nil {
		log.Printf("CreateAlert error: %v", err)
		genericError(c, http.StatusInternalServerError, "failed to create alert")
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "alert created"})
}

func (h *AlertHandler) GetByUser(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		return
	}

	alerts, err := h.alertSvc.GetUserAlerts(userID)
	if err != nil {
		log.Printf("GetUserAlerts error: %v", err)
		genericError(c, http.StatusInternalServerError, "failed to fetch alerts")
		return
	}

	c.JSON(http.StatusOK, alerts)
}

func (h *AlertHandler) Delete(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		return
	}

	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid alert id"})
		return
	}

	if err := h.alertSvc.DeleteAlert(uint(id), userID); err != nil {
		log.Printf("DeleteAlert error: %v", err)
		genericError(c, http.StatusInternalServerError, "failed to delete alert")
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "alert deleted"})
}
