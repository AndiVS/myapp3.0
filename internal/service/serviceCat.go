// Package service encapsulates
package service

import (
	"context"

	"myapp3.0/internal/model"
	"myapp3.0/internal/repository"
)

// Records contains usecase logic for cats
type Records interface {
	SelectAllC(ctx context.Context) ([]*model.Record, error)
	SelectC(ctx context.Context, id string) (*model.Record, error)
	InsertC(ctx context.Context, rec *model.Record) error
	DeleteC(ctx context.Context, id string) error
	UpdateC(ctx context.Context, rec *model.Record) error
}

// Service struct for rep
type Service struct {
	//Rep interface{}
	//Rep *repository.RepositoryPostgres
	Rep *repository.RepositoryMongo
}

// New function for customization service
func New(Rep *repository.RepositoryMongo) *Service {
	return &Service{Rep: Rep}
}

// GetC provides cat
func (serv *Service) GetC(c context.Context, id string) (*model.Record, error) {
	return serv.Rep.SelectC(c, id)
}

// GetAllC provides all cats
func (serv *Service) GetAllC(c context.Context) ([]*model.Record, error) {
	return serv.Rep.SelectAllC(c)
}

// AddC record about cat
func (serv *Service) AddC(c context.Context, rec *model.Record) error {
	return serv.Rep.InsertC(c, rec)
}

// UpdateC updating record about cat
func (serv *Service) UpdateC(c context.Context, rec *model.Record) error {
	return serv.Rep.UpdateC(c, rec)
}

// DeleteC record about cat
func (serv *Service) DeleteC(c context.Context, id string) error {
	return serv.Rep.DeleteC(c, id)
}
