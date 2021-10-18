// Package service encapsulates
package service

import (
	"context"
	"reflect"

	"myapp3.0/internal/model"
	"myapp3.0/internal/repository"
)

// Service struct for rep
type Service struct {
	Rep repository.Cats
}

func NewService(Rep interface{}) *Service {
	var mongo *repository.RepositoryMongo
	var postgres *repository.RepositoryPostgres

	switch reflect.TypeOf(Rep) {
	case reflect.TypeOf(mongo):
		return &Service{Rep: Rep.(*repository.RepositoryMongo)}
	case reflect.TypeOf(postgres):
		return &Service{Rep: Rep.(*repository.RepositoryPostgres)}
	}
	return nil
}

// AddC record about cat
func (serv *Service) AddC(c context.Context, rec *model.Record) error {
	return serv.Rep.InsertC(c, rec)
}

// GetC provides cat
func (serv *Service) GetC(c context.Context, id string) (*model.Record, error) {
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
func (serv *Service) DeleteC(c context.Context, id string) error {
	return serv.Rep.DeleteC(c, id)
}
