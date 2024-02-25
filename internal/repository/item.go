package repository

import (
	"database/sql"
	"errors"
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

func (r *TodoItemRepostirory) GetAll(todoListId int) ([]domain.TodoItem, error) {
	var todoItems []domain.TodoItem

	query := fmt.Sprintf(`SELECT ti.* FROM %s ti INNER JOIN %s li ON li.item_id = ti.id 
	WHERE li.list_id = $1`, todoItemsTable, listsItemsTable)
	err := r.db.Select(&todoItems, query, todoListId)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}

	return todoItems, nil
}

func (r *TodoItemRepostirory) GetById(userId int, todoItemId int) (domain.TodoItem, error) {
	var todoItem domain.TodoItem

	query := fmt.Sprintf(`SELECT ti.* FROM %s ti INNER JOIN %s li ON li.item_id = ti.id INNER JOIN
	%s ul ON ul.list_id = li.list_id WHERE ul.user_id = $1`, todoItemsTable, listsItemsTable, usersListsTable)
	err := r.db.Get(&todoItem, query, userId)
	if err != nil {
		logrus.Error(err)
		if errors.Is(err, sql.ErrNoRows) {
			return todoItem, errors.New("not found")
		}
		return todoItem, err
	}

	return todoItem, nil
}
