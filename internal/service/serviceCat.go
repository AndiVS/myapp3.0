// Package service encapsulates
package service

import (
	"context"

	"github.com/google/uuid"
	"myapp3.0/internal/model"
	"myapp3.0/internal/repository"

	"reflect"
)

// Cats interface for mocks
type Cats interface {
	AddC(c context.Context, rec *model.Record) (uuid.UUID, error)
	GetC(c context.Context, id uuid.UUID) (*model.Record, error)
	GetAllC(c context.Context) ([]*model.Record, error)
	UpdateC(c context.Context, rec *model.Record) error
	DeleteC(c context.Context, id uuid.UUID) error
}

// Service struct for rep
type Service struct {
	Rep repository.Cats
}

// NewService used for setting services
func NewService(Rep interface{}) Cats {
	var mongo *repository.Mongo
	var postgres *repository.Postgres

	switch reflect.TypeOf(Rep) {
	case reflect.TypeOf(mongo):
		return &Service{Rep: Rep.(*repository.Mongo)}
	case reflect.TypeOf(postgres):
		return &Service{Rep: Rep.(*repository.Postgres)}
	}
	return nil
}

// AddC record about cat
func (serv *Service) AddC(c context.Context, rec *model.Record) (uuid.UUID, error) {
	return serv.Rep.InsertC(c, rec)
}

// GetC provides cat
func (serv *Service) GetC(c context.Context, id uuid.UUID) (*model.Record, error) {
	return serv.Rep.SelectC(c, id)
}

// GetAllC provides all cats
func (serv *Service) GetAllC(c context.Context) ([]*model.Record, error) {
	return serv.Rep.SelectAllC(c)
}

// UpdateC updating record about cat
func (serv *Service) UpdateC(c context.Context, rec *model.Record) error {
	return serv.Rep.UpdateC(c, rec)
}

// DeleteC record about cat
func (serv *Service) DeleteC(c context.Context, id uuid.UUID) error {
	return serv.Rep.DeleteC(c, id)
}
