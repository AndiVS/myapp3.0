// Package service encapsulates
package service

import (
	"context"

	"myapp3.0/internal/model"
)

// GetU provides cat
func (serv *Service) GetU(c context.Context, username string) (*model.User, error) {
	return serv.Rep.SelectU(c, username)
}

// GetAllU provides all cats
func (serv *Service) GetAllU(c context.Context) ([]*model.User, error) {
	return serv.Rep.SelectAllU(c)
}

// AddU record about cat
func (serv *Service) AddU(c context.Context, rec *model.User) error {
	return serv.Rep.InsertU(c, rec)
}

// UpdateU updating record about cat
func (serv *Service) UpdateU(c context.Context, username string, isAdmin bool) error {
	return serv.Rep.UpdateU(c, username, isAdmin)
}

// DeleteU record about cat
func (serv *Service) DeleteU(c context.Context, username string) error {
	return serv.Rep.DeleteU(c, username)
}
