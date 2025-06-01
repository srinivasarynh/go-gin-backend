package handlers

import (
	"go-gin-backend/internal/models"
	"go-gin-backend/internal/services"
	"go-gin-backend/internal/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authService *services.AuthService
}

func NewAuthHandler(authService *services.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req models.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "invalid request body")
		return
	}
	
	if err := utils.ValidateStruct(&req); err != nil {
		response.ValidationError(c, err)
		return
	}

	user, err := h.authService.Register(&req)
	if err != nil {
		response.Error(c, http.StatusConflict, err.Error())
		return
	}

	response.Success(c, http.StatusCreated, "user registered successfully", user)
}

func (h *AuthHandler) Login(c *gin.Context) {
    var req models.LoginRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        response.Error(c, http.StatusBadRequest, "Invalid request body")
        return
    }

    if err := utils.ValidateStruct(&req); err != nil {
        response.ValidationError(c, err)
        return
    }

    loginResponse, err := h.authService.Login(&req)
    if err != nil {
        response.Error(c, http.StatusUnauthorized, err.Error())
        return
    }

    response.Success(c, http.StatusOK, "Login successful", loginResponse)
}

func (h *AuthHandler) Me(c *gin.Context) {
    userID, exists := c.Get("user_id")
    if !exists {
        response.Error(c, http.StatusUnauthorized, "User not authenticated")
        return
    }

    username, _ := c.Get("username")
    email, _ := c.Get("email")

    userInfo := map[string]interface{}{
        "user_id":  userID,
        "username": username,
        "email":    email,
    }

    response.Success(c, http.StatusOK, "User info retrieved", userInfo)
}
