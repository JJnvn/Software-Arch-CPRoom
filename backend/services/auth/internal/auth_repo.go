package internal

import (
	"errors"

	"github.com/JJnvn/Software-Arch-CPRoom/backend/services/auth/model"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type UserRepository struct {
	users map[string]*model.User // key: email
}

func NewUserRepository() *UserRepository {
	return &UserRepository{users: make(map[string]*model.User)}
}

func (r *UserRepository) Register(name, email, password string) (*model.User, error) {
	if _, exists := r.users[email]; exists {
		return nil, errors.New("email already registered")
	}
	hashed, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	user := &model.User{
		ID:       uuid.New().String(),
		Name:     name,
		Email:    email,
		Password: string(hashed),
	}
	r.users[email] = user
	return user, nil
}

func (r *UserRepository) GetByEmail(email string) (*model.User, error) {
	user, exists := r.users[email]
	if !exists {
		return nil, errors.New("user not found")
	}
	return user, nil
}
