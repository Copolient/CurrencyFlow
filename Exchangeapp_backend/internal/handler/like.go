package handler

import (
	"exchangeapp/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type LikeHandler struct {
	likeSvc *service.LikeService
}

func NewLikeHandler(likeSvc *service.LikeService) *LikeHandler {
	return &LikeHandler{likeSvc: likeSvc}
}

func (h *LikeHandler) Like(ctx *gin.Context) {
	articleID := ctx.Param("id")

	if err := h.likeSvc.LikeArticle(ctx.Request.Context(), articleID); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Successfully increase likes"})
}

func (h *LikeHandler) GetLikes(ctx *gin.Context) {
	articleID := ctx.Param("id")

	likes, err := h.likeSvc.GetArticleLikes(ctx.Request.Context(), articleID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"likes": likes})
}
