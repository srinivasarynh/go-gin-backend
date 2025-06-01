package handlers


import (
    "net/http"
    "strconv"
    "go-gin-backend/internal/models"
    "go-gin-backend/internal/services"
    "go-gin-backend/internal/utils"
    "go-gin-backend/pkg/response"

    "github.com/gin-gonic/gin"
)

type UserHandler struct {
    userService *services.UserService
}

func NewUserHandler(userService *services.UserService) *UserHandler {
    return &UserHandler{userService: userService}
}

func (h *UserHandler) GetUser(c *gin.Context) {
    idParam := c.Param("id")
    id, err := strconv.ParseUint(idParam, 10, 32)
    if err != nil {
        response.Error(c, http.StatusBadRequest, "Invalid user ID")
        return
    }

    user, err := h.userService.GetUserByID(uint(id))
    if err != nil {
        response.Error(c, http.StatusNotFound, "User not found")
        return
    }

    response.Success(c, http.StatusOK, "User retrieved successfully", user)
}

func (h *UserHandler) UpdateUser(c *gin.Context) {
    idParam := c.Param("id")
    id, err := strconv.ParseUint(idParam, 10, 32)
    if err != nil {
        response.Error(c, http.StatusBadRequest, "Invalid user ID")
        return
    }

    // Check if user is updating their own profile or is admin
    currentUserID, exists := c.Get("user_id")
    if !exists || currentUserID.(uint) != uint(id) {
        response.Error(c, http.StatusForbidden, "Cannot update other user's profile")
        return
    }

    var req models.UpdateUserRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        response.Error(c, http.StatusBadRequest, "Invalid request body")
        return
    }

    if err := utils.ValidateStruct(&req); err != nil {
        response.ValidationError(c, err)
        return
    }

    user, err := h.userService.UpdateUser(uint(id), &req)
    if err != nil {
        response.Error(c, http.StatusBadRequest, err.Error())
        return
    }

    response.Success(c, http.StatusOK, "User updated successfully", user)
}

func (h *UserHandler) DeleteUser(c *gin.Context) {
    idParam := c.Param("id")
    id, err := strconv.ParseUint(idParam, 10, 32)
    if err != nil {
        response.Error(c, http.StatusBadRequest, "Invalid user ID")
        return
    }

    // Check if user is deleting their own profile or is admin
    currentUserID, exists := c.Get("user_id")
    if !exists || currentUserID.(uint) != uint(id) {
        response.Error(c, http.StatusForbidden, "Cannot delete other user's profile")
        return
    }

    if err := h.userService.DeleteUser(uint(id)); err != nil {
        response.Error(c, http.StatusInternalServerError, "Failed to delete user")
        return
    }

    response.Success(c, http.StatusOK, "User deleted successfully", nil)
}

func (h *UserHandler) ListUsers(c *gin.Context) {
    limitStr := c.DefaultQuery("limit", "10")
    offsetStr := c.DefaultQuery("offset", "0")

    limit, err := strconv.Atoi(limitStr)
    if err != nil || limit <= 0 {
        limit = 10
    }
    if limit > 100 {
        limit = 100 // Max limit
    }

    offset, err := strconv.Atoi(offsetStr)
    if err != nil || offset < 0 {
        offset = 0
    }

    users, total, err := h.userService.ListUsers(limit, offset)
    if err != nil {
        response.Error(c, http.StatusInternalServerError, "Failed to list users")
        return
    }

    result := map[string]interface{}{
        "users":  users,
        "total":  total,
        "limit":  limit,
        "offset": offset,
    }

    response.Success(c, http.StatusOK, "Users retrieved successfully", result)
}
