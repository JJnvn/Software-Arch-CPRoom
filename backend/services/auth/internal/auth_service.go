package internal

import (
	"errors"
	"time"

	"github.com/JJnvn/Software-Arch-CPRoom/backend/services/auth/model"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var JWTSecret = []byte("super-secret-key") // in production, load from env

type AuthService struct {
	repo *UserRepository
}

func NewAuthService(repo *UserRepository) *AuthService {
	return &AuthService{repo: repo}
}

func (s *AuthService) Register(name, email, password string) (*model.User, error) {
	return s.repo.Register(name, email, password)
}

func (s *AuthService) Login(email, password string) (string, error) {
	user, err := s.repo.GetByEmail(email)
	if err != nil {
		return "", err
	}
	if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)) != nil {
		return "", errors.New("invalid credentials")
	}

	claims := jwt.MapClaims{
		"user_id": user.ID,
		"email":   user.Email,
		"exp":     time.Now().Add(time.Hour * 1).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(JWTSecret)
}

func (s *AuthService) ValidateToken(tokenStr string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		return JWTSecret, nil
	})
	if err != nil || !token.Valid {
		return nil, errors.New("invalid token")
	}
	return token.Claims.(jwt.MapClaims), nil
}
