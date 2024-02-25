package repository

import (
	"github.com/IvanMeln1k/go-todo-app/internal/domain"
	"github.com/jmoiron/sqlx"
)

type TodoItemRepostirory struct {
	db *sqlx.DB
}

func NewTodoItemRepository(db *sqlx.DB) *TodoItemRepostirory {
	return &TodoItemRepostirory{
		db: db,
	}
}

func (r *TodoItemRepostirory) Create(todoListId int, todoItem domain.TodoItem) (int, error) {
	// tx, err := r.db.Begin()
	// if err != nil {
	// 	logrus.Error(err)
	// 	return 0, err
	// }

	return 0, nil
}