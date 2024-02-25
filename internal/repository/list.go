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

type TodoListRepository struct {
	db *sqlx.DB
}

func NewTodoListRepository(db *sqlx.DB) *TodoListRepository {
	return &TodoListRepository{
		db: db,
	}
}

func (r *TodoListRepository) Create(userId int, todoList domain.TodoList) (int, error) {
	tx, err := r.db.Begin()
	if err != nil {
		logrus.Error(err)
		return 0, err
	}

	var id int
	query := fmt.Sprintf(`INSERT INTO %s (title, description)
	 VALUES ($1, $2) RETURNING id`, todoListsTable)
	row := tx.QueryRow(query, todoList.Title, todoList.Description)
	if err := row.Scan(&id); err != nil {
		logrus.Error(err)
		tx.Rollback()
		return 0, err
	}

	query = fmt.Sprintf(`INSERT INTO %s (user_id, list_id) VALUES ($1, $2)`, usersListsTable)
	_, err = tx.Exec(query, userId, id)
	if err != nil {
		logrus.Error(err)
		tx.Rollback()
		return 0, err
	}

	tx.Commit()
	return id, nil
}

func (r *TodoListRepository) GetAll(userId int) ([]domain.TodoList, error) {
	var todoLists []domain.TodoList

	query := fmt.Sprintf(`SELECT tl.* FROM %s tl INNER JOIN
	 %s ul ON ul.list_id = tl.id WHERE ul.user_id = $1`, todoListsTable, usersListsTable)
	err := r.db.Select(&todoLists, query, userId)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}

	return todoLists, nil
}

func (r *TodoListRepository) GetById(userId int, todoListId int) (domain.TodoList, error) {
	var todoList domain.TodoList

	query := fmt.Sprintf(`SELECT tl.* FROM %s tl INNER JOIN
	%s ul ON ul.list_id = tl.id WHERE ul.user_id = $1 AND tl.id = $2`, todoListsTable, usersListsTable)
	err := r.db.Get(&todoList, query, userId, todoListId)
	if err != nil {
		logrus.Error(err)
		if errors.Is(err, sql.ErrNoRows) {
			return todoList, errors.New("not found")
		}
		return todoList, err
	}

	return todoList, nil
}

func (r *TodoListRepository) Delete(userId int, todoListId int) error {
	query := fmt.Sprintf(`DELETE FROM %s tl USING %s ul WHERE ul.list_id = tl.id AND
	ul.user_id = $1 AND tl.id = $2 RETURNING tl.id`, todoListsTable, usersListsTable)
	row := r.db.QueryRow(query, userId, todoListId)

	var id int
	err := row.Scan(&id)
	if err != nil {
		logrus.Error(err)
		if errors.Is(err, sql.ErrNoRows) {
			return errors.New("not found")
		}
		return err
	}

	return nil
}

func (r *TodoListRepository) Update(userId int, todoListId int, updateTodoList domain.UpdateTodoList) (domain.TodoList, error) {
	var valueNames = make([]string, 0)
	var values = make([]interface{}, 0)
	var argId = 1

	addValue := func(name string, value interface{}) {
		valueNames = append(valueNames, fmt.Sprintf("%s = $%d", name, argId))
		values = append(values, value)
		argId++
	}

	if updateTodoList.Title != nil {
		addValue("title", *updateTodoList.Title)
	}
	if updateTodoList.Description != nil {
		addValue("description", *updateTodoList.Description)
	}

	setQuery := strings.Join(valueNames, ", ")
	query := fmt.Sprintf(`UPDATE %s tl SET %s FROM %s ul WHERE ul.list_id = tl.id AND ul.user_id = $%d
	AND tl.id = $%d RETURNING tl.*`, todoListsTable, setQuery, usersListsTable, argId, argId+1)
	values = append(values, userId, todoListId)

	var todoList domain.TodoList
	err := r.db.Get(&todoList, query, values...)
	if err != nil {
		logrus.Error(err)
		if errors.Is(err, sql.ErrNoRows) {
			return todoList, errors.New("not found")
		}
		return todoList, err
	}

	return todoList, err
}
