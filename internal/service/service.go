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

// Get provides cat
func (serv *Service) Get(c context.Context, id string) (*model.Record, error) {
	return serv.Rep.Select(c, id)
}

// GetAll provides all cats
func (serv *Service) GetAll(c context.Context) ([]*model.Record, error) {
	return serv.Rep.SelectAll(c)
}

// Add record about cat
func (serv *Service) Add(c context.Context, rec *model.Record) error {
	return serv.Rep.Insert(c, rec)
}

// Update updating record about cat
func (serv *Service) Update(c context.Context, rec *model.Record) error {
	return serv.Rep.Update(c, rec)
}

// Delete record about cat
func (serv *Service) Delete(c context.Context, id string) error {
	return serv.Rep.Delete(c, id)
}
