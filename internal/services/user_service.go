package services

import (
	"fmt"
	"go-gin-backend/internal/models"
	"go-gin-backend/internal/repository"
)

type UserService struct {
	userRepo repository.UserRepositoryInterface
}

func NewUserService(userRepo repository.UserRepositoryInterface) *UserService {
	return &UserService{
		userRepo: userRepo,
	}
}

func (s *UserService) GetUserByID(id uint) (*models.UserResponse, error) {
	user, err := s.userRepo.GetByID(id)
	if err != nil {
		return nil, err
	}
	return user.ToResponse(), nil
}


func (s *UserService) UpdateUser(id uint, req *models.UpdateUserRequest) (*models.UserResponse, error) {
    user, err := s.userRepo.GetByID(id)
    if err != nil {
        return nil, err
    }

    if req.FirstName != nil {
        user.FirstName = *req.FirstName
    }
    if req.LastName != nil {
        user.LastName = *req.LastName
    }
    if req.Email != nil {
        existingUser, err := s.userRepo.GetByEmail(*req.Email)
        if err == nil && existingUser.ID != user.ID {
            return nil, fmt.Errorf("email already taken")
        }
        user.Email = *req.Email
    }

    if err := s.userRepo.Update(user); err != nil {
        return nil, err
    }

    return user.ToResponse(), nil
}

func (s *UserService) DeleteUser(id uint) error {
    return s.userRepo.Delete(id)
}

func (s *UserService) ListUsers(limit, offset int) ([]*models.UserResponse, int64, error) {
    users, total, err := s.userRepo.List(limit, offset)
    if err != nil {
        return nil, 0, err
    }

    userResponses := make([]*models.UserResponse, len(users))
    for i, user := range users {
        userResponses[i] = user.ToResponse()
    }

    return userResponses, total, nil
}
