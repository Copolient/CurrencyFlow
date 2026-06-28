package handler

import (
	"exchangeapp/internal/service"
	"log"
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
	followerID, ok := getUserID(c)
	if !ok {
		return
	}

	followeeIDStr := c.Param("id")
	followeeID, err := strconv.ParseUint(followeeIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}

	if err := h.followSvc.Follow(followerID, uint(followeeID)); err != nil {
		log.Printf("Follow error: %v", err)
		genericError(c, http.StatusInternalServerError, "failed to follow user")
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "followed"})
}

func (h *FollowHandler) Unfollow(c *gin.Context) {
	followerID, ok := getUserID(c)
	if !ok {
		return
	}

	followeeIDStr := c.Param("id")
	followeeID, err := strconv.ParseUint(followeeIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}

	if err := h.followSvc.Unfollow(followerID, uint(followeeID)); err != nil {
		log.Printf("Unfollow error: %v", err)
		genericError(c, http.StatusInternalServerError, "failed to unfollow user")
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "unfollowed"})
}

func (h *FollowHandler) IsFollowing(c *gin.Context) {
	followerID, ok := getUserID(c)
	if !ok {
		return
	}

	followeeIDStr := c.Param("id")
	followeeID, err := strconv.ParseUint(followeeIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}

	isFollowing, err := h.followSvc.IsFollowing(followerID, uint(followeeID))
	if err != nil {
		log.Printf("IsFollowing error: %v", err)
		genericError(c, http.StatusInternalServerError, "failed to check follow status")
		return
	}

	c.JSON(http.StatusOK, gin.H{"following": isFollowing})
}
