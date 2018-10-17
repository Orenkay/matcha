package cache

import (
	"github.com/go-redis/redis"
	"github.com/orenkay/matcha/internal/crypto"

	jwt "github.com/dgrijalva/jwt-go"
)

type AuthTokenService struct {
	redis *redis.Client
}

func NewAuthTokenService(r *redis.Client) *AuthTokenService {
	return &AuthTokenService{
		redis: r,
	}
}

func (s *AuthTokenService) Add(tokenString string, token *jwt.Token) error {
	return s.redis.HSet("tokens", tokenString, "").Err()
}

func (s *AuthTokenService) Delete(tokenString string) error {
	return s.redis.HDel("tokens", tokenString).Err()
}

func (s *AuthTokenService) Token(tokenString string) (*jwt.Token, error) {
	exists, err := s.redis.HExists("tokens", tokenString).Result()
	{
		if err != nil {
			return nil, err
		}
		if !exists {
			return nil, nil
		}
	}

	token, err := crypto.DecodeJWT(tokenString)

	// Here we'll check if token has expired
	// If token has expired we remove it from the store
	if !token.Valid {
		if ve, ok := err.(jwt.ValidationError); ok {
			if ve.Errors&jwt.ValidationErrorExpired != 0 {
				return nil, s.Delete(tokenString)
			}
		}
	}
	return token, nil
}
