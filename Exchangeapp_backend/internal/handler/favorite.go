package handler

import (
	"exchangeapp/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type FavoriteHandler struct {
	favoriteSvc *service.FavoriteService
}

func NewFavoriteHandler(favoriteSvc *service.FavoriteService) *FavoriteHandler {
	return &FavoriteHandler{favoriteSvc: favoriteSvc}
}

type favoriteRequest struct {
	FromCurrency string `json:"fromCurrency" binding:"required,len=3"`
	ToCurrency   string `json:"toCurrency" binding:"required,len=3"`
}

func (h *FavoriteHandler) Add(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}

	var req favoriteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.favoriteSvc.AddFavorite(userID.(uint), req.FromCurrency, req.ToCurrency); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "added to favorites"})
}

func (h *FavoriteHandler) GetByUser(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}

	favorites, err := h.favoriteSvc.GetFavorites(userID.(uint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, favorites)
}

func (h *FavoriteHandler) Remove(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}

	from := c.Query("from")
	to := c.Query("to")
	if from == "" || to == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "from and to query params are required"})
		return
	}

	if err := h.favoriteSvc.RemoveFavorite(userID.(uint), from, to); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "removed from favorites"})
}

func (h *FavoriteHandler) Check(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}

	from := c.Query("from")
	to := c.Query("to")
	if from == "" || to == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "from and to query params are required"})
		return
	}

	favorites, err := h.favoriteSvc.GetFavorites(userID.(uint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	favorited := false
	for _, fav := range favorites {
		if fav.FromCurrency == from && fav.ToCurrency == to {
			favorited = true
			break
		}
	}

	c.JSON(http.StatusOK, gin.H{"favorited": favorited})
}
