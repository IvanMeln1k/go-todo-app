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

var (
	ErrUsernameAlreadyInUse    = errors.New("username already in use")
	ErrCreateUser              = errors.New("error to write data")
	ErrGetUser                 = errors.New("error to get data")
	ErrUserNotFound            = errors.New("user not found")
	ErrSessionExpiredOrInvalid = errors.New("session expired or invalid")
	ErrInternal                = errors.New("internal error")
)

func (r *AuthRepository) CreateUser(user domain.User) (int, error) {
	var id int

	query := fmt.Sprintf(`INSERT INTO %s (name, username, password_hash)
	 VALUES ($1, $2, $3) RETURNING id`, usersTable)
	row := r.db.QueryRow(query, user.Name, user.Username, user.Password)

	if err := row.Scan(&id); err != nil {
		logrus.Error(err)
		if pqErr, ok := err.(*pq.Error); ok {
			if pqErr.Code == "23505" {
				return 0, ErrUsernameAlreadyInUse
			}
		}
		return 0, ErrCreateUser
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
			return user, ErrUserNotFound
		}
		return user, ErrGetUser
	}
	return user, nil
}

func (r *AuthRepository) getSessionKey(refreshToken string) string {
	return fmt.Sprintf("sessions:%s", refreshToken)
}

func (r *AuthRepository) getUserSessionsKey(userId int) string {
	return fmt.Sprintf("userSessions:%d", userId)
}

func (r *AuthRepository) CreateSession(ctx context.Context, session domain.Session) error {
	pipe := r.rdb.Pipeline()

	sessionKey := r.getSessionKey(session.Id)
	userSessionKey := r.getUserSessionsKey(session.UserId)

	_, err := pipe.ZAdd(ctx, userSessionKey, redis.Z{
		Score:  0,
		Member: session.Id,
	}).Result()
	if err != nil {
		pipe.Discard()
		logrus.Error(err)
		return ErrInternal
	}

	_, err = pipe.HSet(ctx, sessionKey, map[string]interface{}{
		"userId": session.UserId,
	}).Result()
	if err != nil {
		pipe.Discard()
		logrus.Error(err)
		return ErrInternal
	}

	_, err = pipe.ExpireAt(ctx, sessionKey, session.ExpiresAt).Result()
	if err != nil {
		pipe.Discard()
		logrus.Error(err)
		return ErrInternal
	}

	_, err = pipe.Exec(ctx)
	if err != nil {
		logrus.Error(err)
		return ErrInternal
	}

	return nil
}

func (r *AuthRepository) bindSession(dict map[string]string) (domain.Session, error) {
	userIdStr, ok := dict["userId"]
	if !ok {
		return domain.Session{}, errors.New("bind error")
	}
	userId, err := strconv.Atoi(userIdStr)
	if err != nil {
		return domain.Session{}, errors.New("bind error")
	}
	return domain.Session{
		UserId: userId,
	}, nil
}

func (r *AuthRepository) GetSession(ctx context.Context, refreshToken string) (domain.Session, error) {
	sessionKey := r.getSessionKey(refreshToken)
	rez, err := r.rdb.HGetAll(ctx, sessionKey).Result()
	if err != nil {
		return domain.Session{}, nil
	}

	expireDuration, err := r.rdb.TTL(ctx, sessionKey).Result()
	if err != nil {
		return domain.Session{}, nil
	}
	expiresAt := time.Now().Add(expireDuration)
	if expiresAt.Unix() <= 0 {
		return domain.Session{}, ErrSessionExpiredOrInvalid
	}

	session, err := r.bindSession(rez)
	if err != nil {
		r.deleteSession(ctx, refreshToken)
		return domain.Session{}, ErrSessionExpiredOrInvalid
	}
	session.ExpiresAt = expiresAt
	session.Id = refreshToken

	return session, nil
}

func (r *AuthRepository) deleteSession(ctx context.Context, refreshToken string) error {
	_, err := r.rdb.Del(ctx, r.getSessionKey(refreshToken)).Result()
	return err
}

func (r *AuthRepository) DeleteUserSession(ctx context.Context, userId int, refreshToken string) error {
	pipe := r.rdb.Pipeline()

	sessionKey := r.getSessionKey(refreshToken)
	userSessionKey := r.getUserSessionsKey(userId)

	_, err := pipe.Del(ctx, sessionKey).Result()
	if err != nil {
		pipe.Discard()
		logrus.Error(err)
		return ErrInternal
	}

	_, err = pipe.ZRem(ctx, userSessionKey, refreshToken).Result()
	if err != nil {
		pipe.Discard()
		logrus.Error(err)
		return ErrInternal
	}

	_, err = pipe.Exec(ctx)
	if err != nil {
		logrus.Error(err)
		return ErrInternal
	}

	return nil
}

func (r *AuthRepository) DeleteAllUserSessions(ctx context.Context, userId int) error {
	pipe := r.rdb.Pipeline()

	userSessionKey := r.getUserSessionsKey(userId)

	sessions, err := r.GetAllSessions(ctx, userId)
	if err != nil {
		pipe.Discard()
		logrus.Error(err)
		return ErrInternal
	}
	var tokens []string
	for i := 0; i < len(sessions); i++ {
		tokens = append(tokens, sessions[i].Id)
	}
	var sessionKeys []string
	for i := 0; i < len(tokens); i++ {
		sessionKeys = append(sessionKeys, r.getSessionKey(tokens[i]))
	}

	_, err = pipe.ZRem(ctx, userSessionKey, tokens).Result()
	if err != nil {
		pipe.Discard()
		logrus.Error(err)
		return ErrInternal
	}
	_, err = pipe.Del(ctx, sessionKeys...).Result()
	if err != nil {
		pipe.Discard()
		logrus.Error(err)
		return ErrInternal
	}

	_, err = pipe.Exec(ctx)
	if err != nil {
		logrus.Error(err)
		return ErrInternal
	}

	return nil
}

func (r *AuthRepository) GetCntSessions(ctx context.Context, userId int) (int, error) {
	cnt, err := r.rdb.ZCard(ctx, r.getUserSessionsKey(userId)).Result()
	if err != nil {
		return 0, ErrInternal
	}
	return int(cnt), nil
}

func (r *AuthRepository) GetAllSessions(ctx context.Context, userId int) ([]domain.Session, error) {
	cnt, err := r.GetCntSessions(ctx, userId)
	if err != nil {
		logrus.Error(err)
		return nil, ErrInternal
	}

	refreshTokens, err := r.rdb.ZRange(ctx, r.getUserSessionsKey(userId), 0, int64(cnt)).Result()
	if err != nil {
		logrus.Error(err)
		return nil, ErrInternal
	}

	var sessions []domain.Session

	for i := 0; i < len(refreshTokens); i++ {
		session, err := r.GetSession(ctx, refreshTokens[i])
		if err != nil {
			if errors.Is(err, ErrSessionExpiredOrInvalid) {
				r.rdb.ZRem(ctx, r.getUserSessionsKey(userId), refreshTokens[i])
				continue
			}
			logrus.Error(err)
			return nil, ErrInternal
		}
		sessions = append(sessions, session)
	}

	return sessions, nil
}
