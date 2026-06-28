package handler

import (
	"exchangeapp/internal/service"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type PostHandler struct {
	postSvc *service.PostService
}

func NewPostHandler(postSvc *service.PostService) *PostHandler {
	return &PostHandler{postSvc: postSvc}
}

type postRequest struct {
	Content  string `json:"content" binding:"required,max=2000"`
	Currency string `json:"currency"`
}

func (h *PostHandler) Create(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		return
	}

	var req postRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.postSvc.CreatePost(userID, req.Content, req.Currency); err != nil {
		log.Printf("CreatePost error: %v", err)
		genericError(c, http.StatusInternalServerError, "failed to create post")
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "post created"})
}

func (h *PostHandler) GetAll(c *gin.Context) {
	feedType := c.DefaultQuery("type", "latest")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "20"))

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 50 {
		pageSize = 20
	}

	uid := getUserIDOptional(c)

	posts, err := h.postSvc.GetPosts(feedType, uid, page, pageSize)
	if err != nil {
		log.Printf("GetPosts error: %v", err)
		genericError(c, http.StatusInternalServerError, "failed to fetch posts")
		return
	}

	c.JSON(http.StatusOK, posts)
}

func (h *PostHandler) Like(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		return
	}

	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid post id"})
		return
	}

	liked, err := h.postSvc.LikePost(uint(id), userID)
	if err != nil {
		log.Printf("LikePost error: %v", err)
		genericError(c, http.StatusInternalServerError, "failed to like post")
		return
	}

	if !liked {
		c.JSON(http.StatusOK, gin.H{"message": "already liked"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "liked"})
}
