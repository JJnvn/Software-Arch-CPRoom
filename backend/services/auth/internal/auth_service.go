package internal

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"os"
	"strings"
	"time"

	"github.com/JJnvn/Software-Arch-CPRoom/backend/services/auth/models"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/oauth2"
)

type AuthService struct {
	repo     *AuthRepository
	oauthCfg *oauth2.Config
}

func NewAuthService(
	repo *AuthRepository,
	cfg *oauth2.Config,
) *AuthService {
	return &AuthService{
		repo:     repo,
		oauthCfg: cfg,
	}
}

func (s *AuthService) Register(name, email, password, role string) error {
	// hash password
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user := &models.User{
		Name:     name,
		Email:    email,
		Password: string(hashed),
		Role:     strings.ToUpper(role),
	}

	return s.repo.CreateUser(user)
}

func (s *AuthService) Login(email, password string) (*models.User, string, error) {
	user, err := s.repo.FindByEmail(email)
	if err != nil {
		return nil, "", errors.New("invalid email or password")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, "", errors.New("invalid email or password")
	}

	token, err := s.GenerateJWT(user.Email, user.Role)
	if err != nil {
		return nil, "", err
	}

	return user, token, nil
}

func (s *AuthService) HandleGitHubCallback(code string) (*models.User, error) {
	ctx := context.Background()
	token, err := s.oauthCfg.Exchange(ctx, code)
	if err != nil {
		return nil, err
	}

	client := s.oauthCfg.Client(ctx, token)
	resp, err := client.Get("https://api.github.com/user")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	data, _ := io.ReadAll(resp.Body)
	var ghUser struct {
		Login string `json:"login"`
		Email string `json:"email"`
		Name  string `json:"name"`
	}
	if err := json.Unmarshal(data, &ghUser); err != nil {
		return nil, err
	}

	email := ghUser.Email
	if email == "" {
		email = ghUser.Login + "@github.com"
	}

	user, _ := s.repo.FindByEmail(email)
	if user != nil {
		return user, nil
	}

	newUser := &models.User{
		Name:  ghUser.Name,
		Email: email,
		Role:  models.USER,
	}
	if err := s.repo.CreateUser(newUser); err != nil {
		return nil, err
	}

	return newUser, nil
}

func (s *AuthService) GetByEmail(email string) (*models.User, error) {
	return s.repo.FindByEmail(email)
}

func (s *AuthService) GenerateJWT(email, role string) (string, error) {
	secret := os.Getenv("JWT_SECRET")

	claims := jwt.MapClaims{
		"email": email,
		"role":  role,
		"exp":   time.Now().Add(7 * 24 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

func (s *AuthService) ParseJWT(tokenString string) (jwt.MapClaims, error) {
	secret := os.Getenv("JWT_SECRET")

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if err != nil || !token.Valid {
		return nil, errors.New("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("invalid claims")
	}

	return claims, nil
}
