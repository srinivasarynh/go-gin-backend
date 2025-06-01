package models

import (
	"time"
	"gorm.io/gorm"
)

type User struct {
    ID        uint           `json:"id" gorm:"primaryKey"`
    Email     string         `json:"email" gorm:"uniqueIndex;not null" validate:"required,email"`
    Username  string         `json:"username" gorm:"uniqueIndex;not null" validate:"required,min=3,max=20"`
    Password  string         `json:"-" gorm:"not null" validate:"required,min=8"`
    FirstName string         `json:"first_name" validate:"required,min=2,max=50"`
    LastName  string         `json:"last_name" validate:"required,min=2,max=50"`
    IsActive  bool           `json:"is_active" gorm:"default:true"`
    CreatedAt time.Time      `json:"created_at"`
    UpdatedAt time.Time      `json:"updated_at"`
    DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}


type CreateUserRequest struct {
    Email     string `json:"email" validate:"required,email"`
    Username  string `json:"username" validate:"required,min=3,max=20"`
    Password  string `json:"password" validate:"required,min=8"`
    FirstName string `json:"first_name" validate:"required,min=2,max=50"`
    LastName  string `json:"last_name" validate:"required,min=2,max=50"`
}

type UpdateUserRequest struct {
    FirstName *string `json:"first_name,omitempty" validate:"omitempty,min=2,max=50"`
    LastName  *string `json:"last_name,omitempty" validate:"omitempty,min=2,max=50"`
    Email     *string `json:"email,omitempty" validate:"omitempty,email"`
}

type UserResponse struct {
    ID        uint      `json:"id"`
    Email     string    `json:"email"`
    Username  string    `json:"username"`
    FirstName string    `json:"first_name"`
    LastName  string    `json:"last_name"`
    IsActive  bool      `json:"is_active"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}

func (u *User) ToResponse() *UserResponse {
    return &UserResponse{
        ID:        u.ID,
        Email:     u.Email,
        Username:  u.Username,
        FirstName: u.FirstName,
        LastName:  u.LastName,
        IsActive:  u.IsActive,
        CreatedAt: u.CreatedAt,
        UpdatedAt: u.UpdatedAt,
    }
}
