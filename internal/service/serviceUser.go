// Package service encapsulates
package service

import (
	"context"

	"myapp3.0/internal/model"
	"myapp3.0/internal/repository"

	"reflect"
)

// Users interface for mocks
type Users interface {
	GetUser(c context.Context, username string) (*model.User, error)
	GetAllUser(c context.Context) ([]*model.User, error)
	UpdateUser(c context.Context, username string, isAdmin bool) error
	DeleteUser(c context.Context, username string) error
}

// ServiceUser for generating token
type ServiceUser struct {
	Rep repository.Users
}

// NewServiceUser  for setting new authorizer
func NewServiceUser(repositories interface{}) Users {
	var mongo *repository.Mongo
	var postgres *repository.Postgres

	switch reflect.TypeOf(repositories) {
	case reflect.TypeOf(mongo):
		return &ServiceUser{Rep: repositories.(*repository.Mongo)}
	case reflect.TypeOf(postgres):
		return &ServiceUser{Rep: repositories.(*repository.Postgres)}
	}

	return nil
}

// GetUser provides user
func (s *ServiceUser) GetUser(c context.Context, username string) (*model.User, error) {
	return s.Rep.SelectUser(c, username)
}

// GetAllUser provides all cats
func (s *ServiceUser) GetAllUser(c context.Context) ([]*model.User, error) {
	return s.Rep.SelectAllUser(c)
}

// UpdateUser updating record about cat
func (s *ServiceUser) UpdateUser(c context.Context, username string, isAdmin bool) error {
	return s.Rep.UpdateUser(c, username, isAdmin)
}

// DeleteUser record about cat
func (s *ServiceUser) DeleteUser(c context.Context, username string) error {
	return s.Rep.DeleteUser(c, username)
}
