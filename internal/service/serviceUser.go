// Package service encapsulates
package service

import (
	"context"

	"github.com/AndiVS/myapp3.0/internal/model"
	"github.com/AndiVS/myapp3.0/internal/repository"

	"reflect"
)

// Users interface for mocks
type Users interface {
	GetUser(c context.Context, username string) (*model.User, error)
	GetAllUser(c context.Context) ([]*model.User, error)
	UpdateUser(c context.Context, username string, isAdmin bool) error
	DeleteUser(c context.Context, username string) error
}

// UserService for generating token
type UserService struct {
	Rep repository.Users
}

// NewServiceUser  for setting new authorizer
func NewServiceUser(repositories interface{}) Users {
	var mongo *repository.Mongo
	var postgres *repository.Postgres

	switch reflect.TypeOf(repositories) {
	case reflect.TypeOf(mongo):
		return &UserService{Rep: repositories.(*repository.Mongo)}
	case reflect.TypeOf(postgres):
		return &UserService{Rep: repositories.(*repository.Postgres)}
	}

	return nil
}

// GetUser provides user
func (s *UserService) GetUser(c context.Context, username string) (*model.User, error) {
	return s.Rep.SelectUser(c, username)
}

// GetAllUser provides all cats
func (s *UserService) GetAllUser(c context.Context) ([]*model.User, error) {
	return s.Rep.SelectAllUser(c)
}

// UpdateUser updating record about cat
func (s *UserService) UpdateUser(c context.Context, username string, isAdmin bool) error {
	return s.Rep.UpdateUser(c, username, isAdmin)
}

// DeleteUser record about cat
func (s *UserService) DeleteUser(c context.Context, username string) error {
	return s.Rep.DeleteUser(c, username)
}
