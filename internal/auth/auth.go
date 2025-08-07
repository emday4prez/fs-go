package auth

import (
	"errors"
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

func (s *Service) ValidateJWT(tokenString string) (*CustomClaims, error) {
	// parse the token with custom claims and secret key
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		// provides the key for validation.
		return s.secretKey, nil
	})

	if err != nil {
		return nil, err // this could be due to an expired token or invalid signature
	}

	// check if the token is valid and get the claims
	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}
