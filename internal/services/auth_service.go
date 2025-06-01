package services

import (
	"fmt"
	"go-gin-backend/internal/models"
	"go-gin-backend/internal/repository"
	"go-gin-backend/internal/utils"
)

type AuthService struct {
	userRepo repository.UserRepositoryInterface
	jwtSecret string
}

func NewAuthService(userRepo repository.UserRepository, jwtSecret string) *AuthService {
	return &AuthService{
		userRepo: &userRepo,
		jwtSecret: jwtSecret,
	}
}

func (s *AuthService) Register(req *models.CreateUserRequest) (*models.UserResponse, error) {
	if _, err := s.userRepo.GetByEmail(req.Email); err == nil {
		return nil, fmt.Errorf("user with email already exists")
	}

	if _, err := s.userRepo.GetByUsername(req.Username); err == nil {
		return nil, fmt.Errorf("user with username already exists")
	}

	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	user := &models.User{
		Email: req.Email,
		Username: req.Username,
		Password: hashedPassword,
		FirstName: req.FirstName,
		LastName: req.LastName,
		IsActive: true,
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return user.ToResponse(), nil
}

func (s *AuthService) Login(req *models.LoginRequest) (*models.LoginResponse, error) {
	user, err := s.userRepo.GetByEmail(req.Email)
	if err != nil {
		return nil, fmt.Errorf("invalid credentials")
	}

	if !user.IsActive {
		return nil, fmt.Errorf("account is deactivated")
	}

	if !utils.CheckPassword(req.Password, user.Password) {
		return nil, fmt.Errorf("invalid credentials")
	}

	claims := &models.JWTClaims{
		UserID: user.ID,
		Username: user.Username,
		Email: user.Email,
	}

	token, expiresAt, err := utils.GenerateJWT(claims, s.jwtSecret)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	return &models.LoginResponse{
		Token: token,
		ExpiresAt: expiresAt,
		User: user.ToResponse(),
	}, nil
}

func (s *AuthService) ValidateToken(tokenString string) (*models.JWTClaims, error) {
	return utils.ValidateJWT(tokenString, s.jwtSecret)
}
