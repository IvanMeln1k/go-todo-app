package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/IvanMeln1k/go-todo-app/internal/domain"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
)

type AuthRepository struct {
	db  *sqlx.DB
	rdb *redis.Client
}

func NewAuthRepository(db *sqlx.DB, rdb *redis.Client) *AuthRepository {
	return &AuthRepository{
		db:  db,
		rdb: rdb,
	}
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
		return 0, errors.New("internal server error")
	}

	return id, nil
}

func (r *AuthRepository) GetUser(username, password string) (domain.User, error) {
	var user domain.User

	query := fmt.Sprintf(`SELECT * FROM %s WHERE username = $1 AND password_hash = $2`, usersTable)
	err := r.db.Get(&user, query, username, password)

	if err != nil {
		logrus.Error(err)
		if err == sql.ErrNoRows {
			return user, errors.New("user not found")
		}
		return user, err
	}
	return user, nil
}

func (r *AuthRepository) CreateRefreshToken(ctx context.Context, userId int, refreshToken string) error {
	_, err := r.rdb.ZAdd(ctx, fmt.Sprintf("sessionsuid%d", userId),
		redis.Z{Score: 0, Member: refreshToken}).Result()
	if err != nil {
		logrus.Error(err)
		return err
	}

	ssid := fmt.Sprintf("sessions:%s", refreshToken)
	_, err = r.rdb.HSet(ctx, ssid, map[string]interface{}{
		"userId": userId,
	}).Result()
	if err != nil {
		logrus.Error(err)
		return err
	}
	_, err = r.rdb.Expire(ctx, ssid, time.Hour*24*30).Result()
	if err != nil {
		logrus.Error(err)
		return err
	}

	return nil
}

func (r *AuthRepository) Refresh(ctx context.Context, refreshToken string, newRefreshToken string) (int, error) {
	ssid := fmt.Sprintf("sessions:%s", refreshToken)
	vals, err := r.rdb.HGetAll(ctx, ssid).Result()
	if err != nil {
		logrus.Error(err)
		return 0, err
	}

	_, err = r.rdb.Del(ctx, ssid).Result()
	if err != nil {
		logrus.Error(err)
		return 0, err
	}

	userIdStr, ok := vals["userId"]
	if !ok {
		logrus.Error("Session hasn't userId")
		return 0, errors.New("sessions hasn't userId")
	}
	userId, err := strconv.Atoi(userIdStr)
	if err != nil {
		logrus.Error("UserId isn't integer value")
		return 0, errors.New("userId isn't integer value")
	}

	ssuid := fmt.Sprintf("sessionsuid%d", userId)
	_, err = r.rdb.ZRem(ctx, ssuid, refreshToken).Result()
	if err != nil {
		logrus.Error(err)
		return 0, err
	}

	err = r.CreateRefreshToken(ctx, userId, newRefreshToken)
	if err != nil {
		return 0, err
	}

	return userId, nil
}
