package handler

import (
	"exchangeapp/internal/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type FollowHandler struct {
	followSvc *service.FollowService
}

func NewFollowHandler(followSvc *service.FollowService) *FollowHandler {
	return &FollowHandler{followSvc: followSvc}
}

func (h *FollowHandler) Follow(c *gin.Context) {
	followerID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}

	followeeIDStr := c.Param("id")
	followeeID, err := strconv.ParseUint(followeeIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}

	if err := h.followSvc.Follow(followerID.(uint), uint(followeeID)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "followed"})
}

func (h *FollowHandler) Unfollow(c *gin.Context) {
	followerID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}

	followeeIDStr := c.Param("id")
	followeeID, err := strconv.ParseUint(followeeIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}

	if err := h.followSvc.Unfollow(followerID.(uint), uint(followeeID)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "unfollowed"})
}

func (h *FollowHandler) IsFollowing(c *gin.Context) {
	followerID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}

	followeeIDStr := c.Param("id")
	followeeID, err := strconv.ParseUint(followeeIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}

	isFollowing, err := h.followSvc.IsFollowing(followerID.(uint), uint(followeeID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"following": isFollowing})
}
