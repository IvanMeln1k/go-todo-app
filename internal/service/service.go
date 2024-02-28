package service

import (
	"context"

	"github.com/IvanMeln1k/go-todo-app/internal/domain"
	"github.com/IvanMeln1k/go-todo-app/internal/repository"
)

type Tokens struct {
	AccessToken  string
	RefreshToken string
}

type Authorization interface {
	CreateUser(user domain.User) (int, error)
	SignIn(ctx context.Context, username, password string) (Tokens, error)
	Refresh(ctx context.Context, refreshToken string) (Tokens, error)
	Logout(ctx context.Context, refreshToken string) error
	LogoutAll(ctx context.Context, refreshToken string) error
	ParseToken(tokenString string) (int, error)
}

type TodoList interface {
	Create(userId int, todoList domain.TodoList) (int, error)
	GetAll(userId int) ([]domain.TodoList, error)
	GetById(userId int, todoListId int) (domain.TodoList, error)
	Delete(userId int, todoListId int) error
	Update(userId int, todoListId int, updateTodoList domain.UpdateTodoList) (domain.TodoList, error)
}

type TodoItem interface {
	Create(userId int, todoListId int, todoItem domain.TodoItem) (int, error)
	GetAll(userId int, todoListId int) ([]domain.TodoItem, error)
	GetById(userId int, todoItemId int) (domain.TodoItem, error)
	Delete(userId int, todoItemId int) error
	Update(userId int, todoItemId int, updateTodoItem domain.UpdateTodoItem) (domain.TodoItem, error)
}

type Service struct {
	Authorization
	TodoList
	TodoItem
}

func NewService(repos *repository.Repository) *Service {
	return &Service{
		Authorization: NewAuthService(repos.Authorization),
		TodoList:      NewTodoListService(repos.TodoList),
		TodoItem:      NewTodoItemService(repos.TodoItem, repos.TodoList),
	}
}
