package service

import (
	"github.com/IvanMeln1k/go-todo-app/internal/domain"
	"github.com/IvanMeln1k/go-todo-app/internal/repository"
)

type TodoListService struct {
	repo repository.TodoList
}

func NewTodoListService(repo repository.TodoList) *TodoListService {
	return &TodoListService{
		repo: repo,
	}
}

func (s *TodoListService) Create(userId int, todoList domain.TodoList) (int, error) {
	return s.repo.Create(userId, todoList)
}

func (s *TodoListService) GetAll(userId int) ([]domain.TodoList, error) {
	return s.repo.GetAll(userId)
}

func (s *TodoListService) GetById(userId int, todoListId int) (domain.TodoList, error) {
	return s.repo.GetById(userId, todoListId)
}

func (s *TodoListService) Delete(userId int, todoListId int) error {
	return s.repo.Delete(userId, todoListId)
}

func (s *TodoListService) Update(userId int, todoListId int, updateTodoList domain.UpdateTodoList) (domain.TodoList, error) {
	return s.repo.Update(userId, todoListId, updateTodoList)
}
