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
)

type AuthService struct {
	repo repository.Authorization
}

func NewAuthService(repo repository.Authorization) *AuthService {
	return &AuthService{repo: repo}
}

func (s *AuthService) CreateUser(user domain.User) (int, error) {
	user.Password = s.hashPassword(user.Password)
	return s.repo.CreateUser(user)
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
		return Tokens{}, err
	}

	accessToken, err := s.generateJWT(user.Id)
	if err != nil {
		logrus.Error(err)
		return Tokens{}, err
	}

	refreshToken, err := s.generateRefreshToken()
	if err != nil {
		logrus.Error(err)
		return Tokens{}, err
	}

	err = s.repo.CreateRefreshToken(ctx, user.Id, refreshToken)
	if err != nil {
		logrus.Error(err)
	}

	return Tokens{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (s *AuthService) Refresh(ctx context.Context, refreshToken string) (Tokens, error) {
	newRefreshToken, err := s.generateRefreshToken()
	if err != nil {
		logrus.Error(err)
		return Tokens{}, err
	}

	userId, err := s.repo.Refresh(ctx, refreshToken, newRefreshToken)
	if err != nil {
		return Tokens{}, err
	}

	accessToken, err := s.generateJWT(userId)
	if err != nil {
		logrus.Error(err)
		return Tokens{}, err
	}

	return Tokens{
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken,
	}, nil
}

func (s *AuthService) ParseToken(tokenString string) (int, error) {
	var claims StandardClaimsWithUserId
	_, err := jwt.ParseWithClaims(tokenString, &claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(signingkey), nil
	})

	if claims.ExpiresAt < time.Now().Unix() {
		logrus.Error("token is expired")
		return 0, errors.New("token is expired")
	}

	if err != nil {
		logrus.Error(err)
		return 0, err
	}

	if claims.UserId == 0 {
		return 0, errors.New("invalid token signature")
	}

	return claims.UserId, nil
}
