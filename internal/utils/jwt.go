package utils

import (
	"fmt"
	"go-gin-backend/internal/models"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/golang-jwt/jwt/v5"
)

func GenerateJWT(claims *models.JWTClaims, secret string) (string, time.Time, error) {
	expirationTime := time.Now().Add(24 * time.Hour)

	jwtClaims := jwt.MapClaims{
		"user_id": claims.UserID,
		"username": claims.Username,
		"email": claims.Email,
		"exp": expirationTime.Unix(),
		"iat": time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwtClaims)
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", time.Time{}, fmt.Errorf("failed to sign token: %w", err)
	}
	return tokenString, expirationTime, nil
}

func ValidateJWT(tokenString, secret string) (*models.JWTClaims, error) {
    token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
        if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
            return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
        }
        return []byte(secret), nil
    })

    if err != nil {
        return nil, fmt.Errorf("failed to parse token: %w", err)
    }

    if !token.Valid {
        return nil, fmt.Errorf("invalid token")
    }

    claims, ok := token.Claims.(jwt.MapClaims)
    if !ok {
        return nil, fmt.Errorf("failed to parse claims")
    }

    userID, ok := claims["user_id"].(float64)
    if !ok {
        return nil, fmt.Errorf("invalid user_id claim")
    }

    username, ok := claims["username"].(string)
    if !ok {
        return nil, fmt.Errorf("invalid username claim")
    }

    email, ok := claims["email"].(string)
    if !ok {
        return nil, fmt.Errorf("invalid email claim")
    }

    return &models.JWTClaims{
        UserID:   uint(userID),
        Username: username,
        Email:    email,
    }, nil
}
