package handler

import (
	"exchangeapp/internal/model"
	"exchangeapp/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ArticleHandler struct {
	articleSvc *service.ArticleService
}

func NewArticleHandler(articleSvc *service.ArticleService) *ArticleHandler {
	return &ArticleHandler{articleSvc: articleSvc}
}

func (h *ArticleHandler) Create(ctx *gin.Context) {
	var article model.Article
	if err := ctx.ShouldBindJSON(&article); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.articleSvc.CreateArticle(ctx.Request.Context(), &article); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, article)
}

func (h *ArticleHandler) GetAll(ctx *gin.Context) {
	articles, err := h.articleSvc.GetArticles(ctx.Request.Context())
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, articles)
}

func (h *ArticleHandler) GetByID(ctx *gin.Context) {
	id := ctx.Param("id")

	article, err := h.articleSvc.GetArticleByID(ctx.Request.Context(), id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if article == nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "article not found"})
		return
	}

	ctx.JSON(http.StatusOK, article)
}
