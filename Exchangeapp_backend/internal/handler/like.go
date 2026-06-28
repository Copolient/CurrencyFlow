package handler

import (
	"exchangeapp/internal/service"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type LikeHandler struct {
	likeSvc *service.LikeService
}

func NewLikeHandler(likeSvc *service.LikeService) *LikeHandler {
	return &LikeHandler{likeSvc: likeSvc}
}

func (h *LikeHandler) Like(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		return
	}

	articleID := c.Param("id")
	if articleID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "article id is required"})
		return
	}

	liked, err := h.likeSvc.LikeArticle(c.Request.Context(), articleID, userID)
	if err != nil {
		log.Printf("LikeArticle error: %v", err)
		genericError(c, http.StatusInternalServerError, "failed to like article")
		return
	}

	if !liked {
		c.JSON(http.StatusOK, gin.H{"message": "already liked"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "liked"})
}

func (h *LikeHandler) GetLikes(c *gin.Context) {
	articleID := c.Param("id")
	if articleID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "article id is required"})
		return
	}

	count, err := h.likeSvc.GetArticleLikes(c.Request.Context(), articleID)
	if err != nil {
		log.Printf("GetArticleLikes error: %v", err)
		genericError(c, http.StatusInternalServerError, "failed to get likes")
		return
	}

	c.JSON(http.StatusOK, gin.H{"likes": count})
}
