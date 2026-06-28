package handler

import (
	"exchangeapp/internal/repository"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type UserProfileHandler struct {
	userRepo repository.UserRepository
}

func NewUserProfileHandler(userRepo repository.UserRepository) *UserProfileHandler {
	return &UserProfileHandler{userRepo: userRepo}
}

type updateProfileRequest struct {
	Avatar string `json:"avatar"`
	Bio    string `json:"bio" binding:"max=500"`
}

func (h *UserProfileHandler) GetProfile(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}

	user, err := h.userRepo.FindByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":             user.ID,
		"username":       user.Username,
		"avatar":         user.Avatar,
		"bio":            user.Bio,
		"followersCount": user.FollowersCount,
		"followingCount": user.FollowingCount,
		"createdAt":      user.CreatedAt,
	})
}

func (h *UserProfileHandler) UpdateProfile(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		return
	}

	var req updateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.userRepo.FindByID(userID)
	if err != nil {
		log.Printf("FindByID error: %v", err)
		genericError(c, http.StatusNotFound, "user not found")
		return
	}

	user.Avatar = req.Avatar
	user.Bio = req.Bio

	if err := h.userRepo.Update(user); err != nil {
		log.Printf("Update user error: %v", err)
		genericError(c, http.StatusInternalServerError, "failed to update profile")
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "profile updated"})
}
