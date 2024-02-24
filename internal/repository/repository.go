package repository

import (
	"github.com/IvanMeln1k/go-todo-app/internal/domain"
	"github.com/jmoiron/sqlx"
)

const (
	usersTable = "users"
	todoListsTable = "todo_lists"
	usersListsTable = "users_lists"
	todoItemsTable = "todo_items"
	listsItemsTable = "lists_items"
)

type Authorization interface {
	CreateUser(user domain.User) (int, error)
	GetUser(username, password string) (domain.User, error) 
}

type TodoList interface {
	Create(userId int, list domain.TodoList) (int, error)
	GetAll(userId int) ([]domain.TodoList, error)
	GetById(userId int, todoListId int) (domain.TodoList, error)
	Delete(userId int, todoListId int) error
	Update(userId int, todoListId int, updateTodoList domain.UpdateTodoList) (domain.TodoList, error)
}

type TodoItem interface {

}

type Repository struct {
	Authorization
	TodoList
	TodoItem
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		Authorization: NewAuthRepository(db),
		TodoList: NewTodoListRepository(db),
	}
}