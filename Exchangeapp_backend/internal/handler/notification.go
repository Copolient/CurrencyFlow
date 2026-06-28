package handler

import (
	"exchangeapp/internal/service"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type NotificationHandler struct {
	notifSvc *service.NotificationService
}

func NewNotificationHandler(notifSvc *service.NotificationService) *NotificationHandler {
	return &NotificationHandler{notifSvc: notifSvc}
}

func (h *NotificationHandler) GetAll(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		return
	}

	notifications, err := h.notifSvc.GetNotifications(c.Request.Context(), userID)
	if err != nil {
		log.Printf("GetNotifications error: %v", err)
		genericError(c, http.StatusInternalServerError, "failed to fetch notifications")
		return
	}

	c.JSON(http.StatusOK, notifications)
}

func (h *NotificationHandler) MarkRead(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		return
	}

	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid notification id"})
		return
	}

	if err := h.notifSvc.MarkRead(c.Request.Context(), uint(id), userID); err != nil {
		log.Printf("MarkRead error: %v", err)
		genericError(c, http.StatusInternalServerError, "failed to mark as read")
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "marked as read"})
}

func (h *NotificationHandler) MarkAllRead(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		return
	}

	if err := h.notifSvc.MarkAllRead(c.Request.Context(), userID); err != nil {
		log.Printf("MarkAllRead error: %v", err)
		genericError(c, http.StatusInternalServerError, "failed to mark all as read")
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "all marked as read"})
}

func (h *NotificationHandler) CountUnread(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		return
	}

	count, err := h.notifSvc.CountUnread(c.Request.Context(), userID)
	if err != nil {
		log.Printf("CountUnread error: %v", err)
		genericError(c, http.StatusInternalServerError, "failed to count unread")
		return
	}

	c.JSON(http.StatusOK, gin.H{"count": count})
}
