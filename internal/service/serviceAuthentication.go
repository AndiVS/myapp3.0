package service

import (
	"context"
	"github.com/labstack/echo/v4"
	"myapp3.0/internal/model"
	"myapp3.0/internal/repository"
	"myapp3.0/protocol"
	"net/http"
	"reflect"
)

type Authentication interface {
	SignUp(c context.Context, user *model.User) error
	SignIn(c context.Context, user *model.User) (string, string, error)
}

type AuthenticationService struct {
	Rep repository.Users
	protocol.UnimplementedUserServiceServer
	Access   *JWTManager
	Refresh  *JWTManager
	HashSalt string
}

func NewServiceAuthentication(repositories interface{}, access, refresh *JWTManager, hashSalt string) Authentication {
	var mongo *repository.Mongo
	var postgres *repository.Postgres

	switch reflect.TypeOf(repositories) {
	case reflect.TypeOf(mongo):
		return &AuthenticationService{Rep: repositories.(*repository.Mongo), Access: access, Refresh: refresh, HashSalt: hashSalt}
	case reflect.TypeOf(postgres):
		return &AuthenticationService{Rep: repositories.(*repository.Postgres), Access: access, Refresh: refresh, HashSalt: hashSalt}
	}

	return nil
}

// SignUp record about cat
func (s *AuthenticationService) SignUp(c context.Context, user *model.User) error {
	user.Password = PasswordGenerator(user.Password, s.HashSalt)

	_, err := s.Rep.SelectUser(c, user.Username)
	if err != nil {
		return s.Rep.InsertUser(c, user)
	}

	return echo.NewHTTPError(http.StatusInternalServerError, "UNABLE TO INSERT ")
}

// SignIn generate token
func (s *AuthenticationService) SignIn(c context.Context, user *model.User) (string, string, error) {
	user.Password = PasswordGenerator(user.Password, s.HashSalt)

	user1, err := s.Rep.SelectUser(c, user.Username)
	if err != nil {
		return "", "", err
	}
	if !PasswordChek(user1.Password, user.Password) {
		return "", "", err
	}

	return GenerateTokens(user, s.Access, s.Refresh)
}
