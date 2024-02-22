package service

import (
	"crypto/sha1"

	"github.com/IvanMeln1k/go-todo-app/internal/domain"
	"github.com/IvanMeln1k/go-todo-app/internal/repository"
)

const salt string = "ghu835mgd823"

type AuthService struct {
	repo repository.Authorization
}

func NewAuthService(repo repository.Repository) *AuthService {
	return &AuthService{repo: repo}
}

func (s *AuthService) CreateUser(user domain.User) (int, error) {
	user.Password = s.hashPassword(user.Password);
	return s.repo.CreateUser(user);
}

func (s *AuthService) hashPassword(password string) string {
	hash := sha1.New();
	hash.Write([]byte(password)) 
	return string(hash.Sum([]byte(salt)))
}