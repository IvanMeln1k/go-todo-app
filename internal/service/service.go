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
	Delete(userId int, todoItemId int) (error)
}

type Service struct {
	Authorization
	TodoList
	TodoItem
}

func NewService(repos *repository.Repository) *Service {
	return &Service{
		Authorization: NewAuthService(repos.Authorization),
		TodoList: NewTodoListService(repos.TodoList),
		TodoItem: NewTodoItemService(repos.TodoItem, repos.TodoList),
	}
}