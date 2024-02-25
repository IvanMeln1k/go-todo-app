package service

import (
	"github.com/IvanMeln1k/go-todo-app/internal/domain"
	"github.com/IvanMeln1k/go-todo-app/internal/repository"
)

type TodoItemService struct {
	repo repository.TodoItem
	listRepo repository.TodoList
}

func NewTodoItemService(repo repository.TodoItem, listRepo repository.TodoList) *TodoItemService {
	return &TodoItemService{
		repo: repo,
		listRepo: listRepo,
	}
}

func (s *TodoItemService) Create(userId int, todoListId int, todoItem domain.TodoItem) (int, error) {
	todoList, err := s.listRepo.GetById(userId, todoListId)
	if err != nil {
		return 0, err
	}

	return s.repo.Create(todoList.Id, todoItem)
}

func (s *TodoItemService) GetAll(userId int, todoListId int) ([]domain.TodoItem, error) {
	todoList, err := s.listRepo.GetById(userId, todoListId)
	if err != nil {
		return nil, err
	}
	
	return s.repo.GetAll(todoList.Id)
}

func (s *TodoItemService) GetById(userId int, todoItemId int) (domain.TodoItem, error) {
	return s.repo.GetById(userId, todoItemId)
}