package service

import (
	"crypto/sha1"
	"errors"
	"fmt"
	"time"

	"github.com/IvanMeln1k/go-todo-app/internal/domain"
	"github.com/IvanMeln1k/go-todo-app/internal/repository"
	"github.com/dgrijalva/jwt-go"
	"github.com/sirupsen/logrus"
)

const ( 
	salt string = "ghu835mgd823"
	signingkey = "lgk;bfsdtrg"
	tokenTTL = time.Hour * 12
)

type AuthService struct {
	repo repository.Authorization
}

func NewAuthService(repo repository.Authorization) *AuthService {
	return &AuthService{repo: repo}
}

func (s *AuthService) CreateUser(user domain.User) (int, error) {
	user.Password = s.hashPassword(user.Password);
	return s.repo.CreateUser(user);
}

func (s *AuthService) hashPassword(password string) string {
	hash := sha1.New();
	hash.Write([]byte(password)) 
	return fmt.Sprintf("%x", hash.Sum([]byte(salt)))
}

type StandardClaimsWithUserId struct {
	jwt.StandardClaims
	UserId int `json:"user_id"`
}

func (s *AuthService) GenerateToken(username, password string) (string, error) {
	user, err := s.repo.GetUser(username, s.hashPassword(password))
	if err != nil {
		return "", err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &StandardClaimsWithUserId{
		jwt.StandardClaims {
			ExpiresAt: time.Now().UTC().Add(tokenTTL).Unix(),
			IssuedAt: time.Now().UTC().Unix(),	
		},
		user.Id, 
	})

	return token.SignedString([]byte(signingkey))
}

func (s *AuthService) ParseToken(tokenString string) (int, error) {
	var claims StandardClaimsWithUserId
	_, err := jwt.ParseWithClaims(tokenString, &claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(signingkey), nil
	})

	if err != nil {
		logrus.Error(err)
		return 0, err
	}

	if claims.UserId == 0 {
		return 0, errors.New("invalid token signature")
	}

	return claims.UserId, nil
}