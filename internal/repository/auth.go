package repository

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/IvanMeln1k/go-todo-app/internal/domain"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"github.com/sirupsen/logrus"
)

type AuthRepository struct {
	db *sqlx.DB
}

func NewAuthRepository(db *sqlx.DB) *AuthRepository {
	return &AuthRepository{db: db}
}

func (r *AuthRepository) CreateUser(user domain.User) (int, error) {
	var id int

	query := fmt.Sprintf(`INSERT INTO %s (name, username, password_hash)
	 VALUES ($1, $2, $3) RETURNING id`, usersTable)
	row := r.db.QueryRow(query, user.Name, user.Username, user.Password)

	if err := row.Scan(&id); err != nil {
		logrus.Error(err)
		if pqErr, ok := err.(*pq.Error); ok {
			if pqErr.Code == "23505" {
				return 0, errors.New("username already in use")
			}
		}
		return 0, errors.New("internal server error");
	}

	return id, nil
}

func (r *AuthRepository) GetUser(username, password string) (domain.User, error) {
	var user domain.User

	query := fmt.Sprintf(`SELECT * FROM %s WHERE username = $1 AND password_hash = $2`, usersTable)
	err := r.db.Get(&user, query, username, password)

	if err != nil {
		if err == sql.ErrNoRows {
			return user, errors.New("user not found")
		}
		return user, err
	}
	return user, nil;
}