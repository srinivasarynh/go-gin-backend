package repository

import "go-gin-backend/internal/models"

type UserRepositoryInterface interface {
Create(user *models.User) error
    GetByID(id uint) (*models.User, error)
    GetByEmail(email string) (*models.User, error)
    GetByUsername(username string) (*models.User, error)
    Update(user *models.User) error
    Delete(id uint) error
    List(limit, offset int) ([]*models.User, int64, error)
}
