package auth

import (
	"time"

	"github.com/emday4prez/fs-go/internal/domain"
	"github.com/golang-jwt/jwt/v5"
)

type Service struct {
	secretKey []byte
}

func NewService(secret string) *Service {
	return &Service{
		secretKey: []byte(secret),
	}
}

type CustomClaims struct {
	UserID string `json:"userID"`
	jwt.RegisteredClaims
}

func (s *Service) GenerateJWT(user *domain.User) (string, error) {
	claims := CustomClaims{
		UserID: user.ID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),

			Issuer: "go-file-server",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedString, err := token.SignedString(s.secretKey)
	if err != nil {
		return "", err
	}

	return signedString, nil
}
