// Package service encapsulates
package service

import (
	"context"

	"myapp3.0/internal/model"
	"myapp3.0/internal/repository"
)

// Service struct for rep
type Service struct {
	Rep *repository.Repository
}

// New function for customization service
func New(Rep *repository.Repository) *Service {
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
