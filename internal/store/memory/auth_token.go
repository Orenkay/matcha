package memory

import (
	"sync"

	jwt "github.com/dgrijalva/jwt-go"
)

type AuthTokenService struct {
	tokens map[string]*jwt.Token
	mu     sync.Mutex
}

func NewAuthTokenService() *AuthTokenService {
	return &AuthTokenService{
		tokens: make(map[string]*jwt.Token),
	}
}

func (s *AuthTokenService) Add(tokenString string, token *jwt.Token) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.tokens[tokenString] = token
}

func (s *AuthTokenService) Delete(tokenString string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.tokens, tokenString)
}

func (s *AuthTokenService) Exists(token string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	_, ok := s.tokens[token]
	return ok
}
