package cache

import (
	"encoding/json"
	"strconv"

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

func (s *AuthTokenService) push(tokenString string, userID int64) error {
	var tokens []string
	target := strconv.FormatInt(userID, 10)
	b, err := s.redis.HGet("users_tokens", target).Bytes()
	{
		if err != nil && err != redis.Nil {
			return err
		}
		if err != redis.Nil {
			if err := json.Unmarshal(b, &tokens); err != nil {
				return err
			}
		}
	}
	if err := s.redis.HSet("tokens", tokenString, target).Err(); err != nil {
		return err
	}
	tokens = append(tokens, tokenString)
	b, err = json.Marshal(tokens)
	{
		if err != nil {
			return err
		}
	}
	return s.redis.HSet("users_tokens", target, string(b)).Err()
}

func (s *AuthTokenService) delete(tokenString string, userID int64) error {
	target := strconv.FormatInt(userID, 10)
	b, err := s.redis.HGet("users_tokens", target).Bytes()
	{
		if err != nil {
			return err
		}
	}
	var tokens []string
	if err := json.Unmarshal(b, &tokens); err != nil {
		return err
	}
	for i, t := range tokens {
		if t == tokenString {
			tokens = append(tokens[:i], tokens[i+1:]...)
			break
		}
	}
	if err := s.redis.HDel("tokens", tokenString).Err(); err != nil {
		return err
	}
	return s.redis.HSet("users_tokens", target, tokens).Err()
}

func (s *AuthTokenService) Add(tokenString string, userID int64) error {
	return s.push(tokenString, userID)
}

func (s *AuthTokenService) Delete(tokenString string) error {
	uid, err := s.redis.HGet("tokens", tokenString).Result()
	{
		if err != nil {
			if err == redis.Nil {
				return nil
			}
			return err
		}
	}
	userID, err := strconv.Atoi(uid)
	{
		if err != nil {
			return err
		}
	}
	return s.delete(tokenString, int64(userID))
}

func (s *AuthTokenService) DeleteByUserID(userID int64) error {
	target := strconv.FormatInt(userID, 10)
	b, err := s.redis.HGet("users_tokens", target).Bytes()
	{
		if err != nil {
			return err
		}
	}
	var tokens []string
	if err := json.Unmarshal(b, &tokens); err != nil {
		return err
	}
	for _, t := range tokens {
		if err := s.redis.HDel("tokens", t).Err(); err != nil {
			return err
		}
	}
	return s.redis.HDel("users_tokens", target).Err()
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
