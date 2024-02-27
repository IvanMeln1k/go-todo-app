package repository

import (
	"context"

	"github.com/IvanMeln1k/go-todo-app/internal/domain"
	"github.com/jmoiron/sqlx"
	"github.com/redis/go-redis/v9"
)

const (
	usersTable      = "users"
	todoListsTable  = "todo_lists"
	usersListsTable = "users_lists"
	todoItemsTable  = "todo_items"
	listsItemsTable = "lists_items"
)

type Authorization interface {
	CreateUser(user domain.User) (int, error)
	GetUser(username, password string) (domain.User, error)
	CreateRefreshToken(ctx context.Context, userId int, refreshToken string) error
	Refresh(ctx context.Context, refreshToken string, newRefreshToken string) (int, error)
}

type TodoList interface {
	Create(userId int, list domain.TodoList) (int, error)
	GetAll(userId int) ([]domain.TodoList, error)
	GetById(userId int, todoListId int) (domain.TodoList, error)
	Delete(userId int, todoListId int) error
	Update(userId int, todoListId int, updateTodoList domain.UpdateTodoList) (domain.TodoList, error)
}

type TodoItem interface {
	Create(todoListId int, todoItem domain.TodoItem) (int, error)
	GetAll(todoListId int) ([]domain.TodoItem, error)
	GetById(userId int, todoItemId int) (domain.TodoItem, error)
	Delete(userId int, todoItemId int) error
	Update(userId int, todoItemId int, updateTodoItem domain.UpdateTodoItem) (domain.TodoItem, error)
}

type Repository struct {
	Authorization
	TodoList
	TodoItem
}

func NewRepository(db *sqlx.DB, rdb *redis.Client) *Repository {
	return &Repository{
		Authorization: NewAuthRepository(db, rdb),
		TodoList:      NewTodoListRepository(db),
		TodoItem:      NewTodoItemRepository(db),
	}
}
