package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/IvanMeln1k/go-todo-app/internal/domain"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
)

type TodoItemRepository struct {
	db *sqlx.DB
}

func NewTodoItemRepository(db *sqlx.DB) *TodoItemRepository {
	return &TodoItemRepository{
		db: db,
	}
}

func (r *TodoItemRepository) Create(todoListId int, todoItem domain.TodoItem) (int, error) {
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

func (r *TodoItemRepository) GetAll(todoListId int) ([]domain.TodoItem, error) {
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

func (r *TodoItemRepository) GetById(userId int, todoItemId int) (domain.TodoItem, error) {
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

func (r *TodoItemRepository) Delete(userId int, todoItemId int) (error) {
	query := fmt.Sprintf(`DELETE FROM %s ti USING %s li, %s ul WHERE li.item_id = ti.id AND
	li.list_id = ul.list_id AND ul.user_id = $1 AND ti.id = $2 RETURNING ti.id`,
	todoItemsTable, listsItemsTable, usersListsTable)
	var id int
	fmt.Println(query)
	row := r.db.QueryRow(query, userId, todoItemId)
	if err := row.Scan(&id); err != nil {
		logrus.Error(err)
		if errors.Is(err, sql.ErrNoRows) {
			return errors.New("not found")
		}
		return err
	}
	
	return nil
}

func (r *TodoItemRepository) Update(userId int, todoItemId int, updateTodoItem domain.UpdateTodoItem) (domain.TodoItem, error) {
	var values = make([]interface{}, 0)
	var names = make([]string, 0)
	var argId = 1

	appendArg := func(name string, value interface{}) {
		names = append(names, fmt.Sprintf("%s = $%d", name, argId))
		values = append(values, value)
		argId++
	}

	if updateTodoItem.Title != nil {
		appendArg("title", *updateTodoItem.Title)
	}

	if updateTodoItem.Description != nil {
		appendArg("description", *updateTodoItem.Description)
	}

	if updateTodoItem.Done != nil {
		appendArg("done", *updateTodoItem.Done)
	}

	setQuery := strings.Join(names, ", ")
	values = append(values, userId, todoItemId)

	query := fmt.Sprintf(`UPDATE %s ti SET %s FROM %s li, %s ul WHERE ti.id = li.item_id AND
	ul.list_id = li.list_id AND ul.user_id = $%d AND ti.id = $%d RETURNING ti.*`, todoItemsTable, setQuery,
	listsItemsTable, usersListsTable, argId, argId + 1)

	var todoItem domain.TodoItem
	err := r.db.Get(&todoItem, query, values...)
	if err != nil {
		logrus.Error(err)
		if errors.Is(err, sql.ErrNoRows) {
			return todoItem, errors.New("not found")
		}
		return todoItem, err
	}

	return todoItem, nil
}