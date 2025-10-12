package internal

import (
	"context"
	"encoding/json"
	"errors"
	"io"

	"github.com/JJnvn/Software-Arch-CPRoom/backend/services/auth/model"
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

func (s *AuthService) Register(name, email, password string) error {
	// hash password
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user := &model.User{
		Name:     name,
		Email:    email,
		Password: string(hashed),
	}

	return s.repo.CreateUser(user)
}

func (s *AuthService) Login(email, password string) (*model.User, error) {
	user, err := s.repo.FindByEmail(email)
	if err != nil {
		return nil, errors.New("invalid email or password")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, errors.New("invalid email or password")
	}

	return user, nil
}

func (s *AuthService) HandleGitHubCallback(code string) (*model.User, error) {
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

	newUser := &model.User{
		Name:  ghUser.Name,
		Email: email,
	}
	if err := s.repo.CreateUser(newUser); err != nil {
		return nil, err
	}

	return newUser, nil
}

func (s *AuthService) GetByEmail(email string) (*model.User, error) {
	return s.repo.FindByEmail(email)
}
