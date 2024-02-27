package service

import (
	"context"
	"crypto/sha1"
	"errors"
	"fmt"
	"math/rand"
	"time"

	"github.com/IvanMeln1k/go-todo-app/internal/domain"
	"github.com/IvanMeln1k/go-todo-app/internal/repository"
	"github.com/dgrijalva/jwt-go"
	"github.com/sirupsen/logrus"
)

const (
	salt       string = "ghu835mgd823"
	signingkey        = "lgk;bfsdtrg"
	tokenTTL          = time.Hour * 12
	sessionTTL        = 30 * time.Hour * 24
)

type AuthService struct {
	repo repository.Authorization
}

func NewAuthService(repo repository.Authorization) *AuthService {
	return &AuthService{repo: repo}
}

var (
	ErrUsernameAlreadyInUse      = errors.New("username already in use")
	ErrInvalidUsernameOrPassowrd = errors.New("invalid username or password")
	ErrUserNotFound              = errors.New("user not found")
	ErrCreateUser                = errors.New("error to create user")
	ErrGetUser                   = errors.New("error to get user")
	ErrInternal                  = errors.New("internal error")
	ErrTokenExpired              = errors.New("token expired")
	ErrInvalidTokenSignature     = errors.New("invalid token signature")
	ErrInvalidSession            = errors.New("invalid session")
	ErrSessionExpired            = errors.New("session expired")
)

func (s *AuthService) CreateUser(user domain.User) (int, error) {
	user.Password = s.hashPassword(user.Password)
	userId, err := s.repo.CreateUser(user)
	if err != nil {
		if errors.Is(err, repository.ErrCreateUser) {
			return 0, ErrUsernameAlreadyInUse
		} else {
			return 0, ErrCreateUser
		}
	}
	return userId, nil
}

func (s *AuthService) hashPassword(password string) string {
	hash := sha1.New()
	hash.Write([]byte(password))
	return fmt.Sprintf("%x", hash.Sum([]byte(salt)))
}

type StandardClaimsWithUserId struct {
	jwt.StandardClaims
	UserId int `json:"user_id"`
}

func (s *AuthService) generateJWT(userId int) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &StandardClaimsWithUserId{
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(tokenTTL).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
		userId,
	})

	return token.SignedString([]byte(signingkey))
}

func (s *AuthService) generateRefreshToken() (string, error) {
	b := make([]byte, 32)

	src := rand.NewSource(time.Now().Unix())
	r := rand.New(src)

	if _, err := r.Read(b); err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", b), nil
}

func (s *AuthService) SignIn(ctx context.Context, username, password string) (Tokens, error) {
	user, err := s.repo.GetUser(username, s.hashPassword(password))
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			return Tokens{}, ErrInvalidUsernameOrPassowrd
		}
		return Tokens{}, ErrInternal
	}

	accessToken, err := s.generateJWT(user.Id)
	if err != nil {
		return Tokens{}, ErrInternal
	}

	refreshToken, err := s.generateRefreshToken()
	if err != nil {
		return Tokens{}, ErrInternal
	}

	cntSessions, err := s.repo.GetCntSessions(ctx, user.Id)
	if err != nil {
		return Tokens{}, ErrInternal
	}

	if cntSessions >= 5 {
		sessions, err := s.repo.GetAllSessions(ctx, user.Id)
		if err != nil {
			return Tokens{}, ErrInternal
		}
		domain.SortSessionsByTime(&sessions)
		cntSessions = len(sessions)
		for i := 0; i < len(sessions); i++ {
			if cntSessions < 5 && sessions[i].ExpiresAt.Unix() > time.Now().Unix() {
				break
			}
			err = s.repo.DeleteUserSession(ctx, user.Id, sessions[i].Id)
			if err != nil {
				return Tokens{}, ErrInternal
			}
			cntSessions--
		}
	}

	err = s.repo.CreateSession(ctx, domain.Session{
		UserId:    user.Id,
		Id:        refreshToken,
		ExpiresAt: time.Now().Add(sessionTTL),
	})
	if err != nil {
		return Tokens{}, ErrInternal
	}
	return Tokens{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (s *AuthService) Refresh(ctx context.Context, refreshToken string) (Tokens, error) {
	session, err := s.repo.GetSession(ctx, refreshToken)
	if err != nil {
		if errors.Is(err, repository.ErrSessionExpired) {
			return Tokens{}, ErrSessionExpired
		}
		return Tokens{}, ErrInternal
	}

	accessToken, err := s.generateJWT(session.UserId)
	if err != nil {
		return Tokens{}, ErrInternal
	}

	err = s.repo.DeleteUserSession(ctx, session.UserId, refreshToken)
	if err != nil {
		return Tokens{}, ErrInternal
	}

	refreshToken, err = s.generateRefreshToken()
	if err != nil {
		return Tokens{}, ErrInternal
	}

	err = s.repo.CreateSession(ctx, domain.Session{
		UserId:    session.UserId,
		Id:        refreshToken,
		ExpiresAt: time.Now().Add(sessionTTL),
	})
	if err != nil {
		return Tokens{}, ErrInternal
	}

	return Tokens{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (s *AuthService) ParseToken(tokenString string) (int, error) {
	token, err := jwt.ParseWithClaims(tokenString, &StandardClaimsWithUserId{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(signingkey), nil
	})

	claims, ok := token.Claims.(*StandardClaimsWithUserId)
	if !ok {
		return 0, ErrInvalidTokenSignature
	}

	if err != nil {
		logrus.Error(err)
		if claims.ExpiresAt != 0 && claims.ExpiresAt < time.Now().Unix() {
			return 0, ErrTokenExpired
		}
		return 0, ErrInvalidTokenSignature
	}

	return claims.UserId, nil
}
