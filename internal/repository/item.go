package repository

import (
	"fmt"

	"github.com/IvanMeln1k/go-todo-app/internal/domain"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
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
	tx, err := r.db.Begin()
	if err != nil {
		logrus.Error(err)
		return 0, err
	}

	query := fmt.Sprintf(`INSERT INTO %s (title, description, done) VALUES 
	($1, $2, $3) RETURNING id`, todoItemsTable)
	row := tx.QueryRow(query, todoItem.Title, todoItem.Description, todoItem.Done)
	if err = row.Scan(&todoItem.Id); err != nil {
		logrus.Error(err)
		tx.Rollback()
		return 0, err
	}

	query = fmt.Sprintf(`INSERT INTO %s (list_id, item_id) VALUES
	($1, $2)`, listsItemsTable)
	_, err = tx.Exec(query, todoListId, todoItem.Id)
	if err != nil {
		logrus.Error(err)
		tx.Rollback()
		return 0, err
	}

	tx.Commit()
	return todoItem.Id, nil
}