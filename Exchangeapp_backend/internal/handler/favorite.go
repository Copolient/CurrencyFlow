package handler

import (
	"exchangeapp/internal/service"
	"log"
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
	userID, ok := getUserID(c)
	if !ok {
		return
	}

	var req favoriteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.favoriteSvc.AddFavorite(userID, req.FromCurrency, req.ToCurrency); err != nil {
		log.Printf("AddFavorite error: %v", err)
		genericError(c, http.StatusInternalServerError, "failed to add favorite")
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "added to favorites"})
}

func (h *FavoriteHandler) GetByUser(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		return
	}

	favorites, err := h.favoriteSvc.GetFavorites(userID)
	if err != nil {
		log.Printf("GetFavorites error: %v", err)
		genericError(c, http.StatusInternalServerError, "failed to fetch favorites")
		return
	}

	c.JSON(http.StatusOK, favorites)
}

func (h *FavoriteHandler) Remove(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		return
	}

	from := c.Query("from")
	to := c.Query("to")
	if from == "" || to == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "from and to query params are required"})
		return
	}

	fromCode, valid := sanitizeCurrencyCode(from)
	if !valid {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid from currency code"})
		return
	}
	toCode, valid := sanitizeCurrencyCode(to)
	if !valid {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid to currency code"})
		return
	}

	removed, err := h.favoriteSvc.RemoveFavorite(userID, fromCode, toCode)
	if err != nil {
		log.Printf("RemoveFavorite error: %v", err)
		genericError(c, http.StatusInternalServerError, "failed to remove favorite")
		return
	}

	if !removed {
		c.JSON(http.StatusNotFound, gin.H{"error": "favorite not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "removed from favorites"})
}

func (h *FavoriteHandler) Check(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		return
	}

	from := c.Query("from")
	to := c.Query("to")
	if from == "" || to == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "from and to query params are required"})
		return
	}

	fromCode, valid := sanitizeCurrencyCode(from)
	if !valid {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid from currency code"})
		return
	}
	toCode, valid := sanitizeCurrencyCode(to)
	if !valid {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid to currency code"})
		return
	}

	favorites, err := h.favoriteSvc.GetFavorites(userID)
	if err != nil {
		log.Printf("GetFavorites error: %v", err)
		genericError(c, http.StatusInternalServerError, "failed to check favorite")
		return
	}

	favorited := false
	for _, fav := range favorites {
		if fav.FromCurrency == fromCode && fav.ToCurrency == toCode {
			favorited = true
			break
		}
	}

	c.JSON(http.StatusOK, gin.H{"favorited": favorited})
}
