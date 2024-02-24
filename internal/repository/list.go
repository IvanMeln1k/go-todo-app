package repository

import (
	"database/sql"
	"errors"
	"fmt"

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
	tx, err := r.db.Begin();
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
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("not found")
		}
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