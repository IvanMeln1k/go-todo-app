package service

import (
	"github.com/IvanMeln1k/go-todo-app/internal/domain"
	"github.com/IvanMeln1k/go-todo-app/internal/repository"
)

type Authorization interface {
	CreateUser(user domain.User) (int, error)
	GenerateToken(username, password string) (string, error)
	ParseToken(tokenString string) (int, error)
}

type TodoList interface {

}

type TodoItem interface {

}

type Service struct {
	Authorization
	TodoList
	TodoItem
}

func NewService(repos *repository.Repository) *Service {
	return &Service{
		Authorization: NewAuthService(repos.Authorization),
	}
}