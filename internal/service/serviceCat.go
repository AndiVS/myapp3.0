// Package service encapsulates
package service

import (
	"context"

	"github.com/AndiVS/myapp3.0/internal/model"
	"github.com/AndiVS/myapp3.0/internal/redis/consumer"
	"github.com/AndiVS/myapp3.0/internal/redis/producer"
	"github.com/AndiVS/myapp3.0/internal/repository"
	"github.com/google/uuid"

	"reflect"
)

// Cats interface for mocks
type Cats interface {
	AddCat(c context.Context, rec *model.Cat) (uuid.UUID, error)
	GetCat(c context.Context, id uuid.UUID) (*model.Cat, error)
	GetAllCat(c context.Context) ([]*model.Cat, error)
	UpdateCat(c context.Context, rec *model.Cat) error
	DeleteCat(c context.Context, id uuid.UUID) error
}

// ServiceCat struct for rep
type ServiceCat struct {
	Rep    repository.Cats
	CatMap map[string]*model.Cat
	Redis  *model.Redis
}

// NewServiceCat used for setting services
func NewServiceCat(Rep interface{}, Redis *model.Redis) Cats {
	serviceCat := ServiceCat{}

	var mongo *repository.Mongo
	var postgres *repository.Postgres

	switch reflect.TypeOf(Rep) {
	case reflect.TypeOf(mongo):
		serviceCat.Rep = Rep.(*repository.Mongo)
	case reflect.TypeOf(postgres):
		serviceCat.Rep = Rep.(*repository.Postgres)
	}

	catMap := make(map[string]*model.Cat)
	catsSlice, _ := serviceCat.Rep.SelectAllCat(context.Background())
	for _, cat := range catsSlice {
		catMap[cat.ID.String()] = cat
	}
	serviceCat.CatMap = catMap

	go consumer.ConsumeEvents(Redis, catMap)

	serviceCat.Redis = Redis

	return &serviceCat
}

// AddCat record about cat
func (s *ServiceCat) AddCat(c context.Context, cat *model.Cat) (uuid.UUID, error) {
	s.CatMap[cat.ID.String()] = cat
	producer.GenerateEvent("cat", "Insert", cat, s.Redis.Client, s.Redis.StreamName)
	return s.Rep.InsertCat(c, cat)
}

// GetCat provides cat
func (s *ServiceCat) GetCat(c context.Context, id uuid.UUID) (*model.Cat, error) {
	val, ok := s.CatMap[id.String()]
	if ok {
		return val, nil
	} else {
		return s.Rep.SelectCat(c, id)
	}
}

// GetAllCat provides all cats
func (s *ServiceCat) GetAllCat(c context.Context) ([]*model.Cat, error) {
	return s.Rep.SelectAllCat(c)
}

// UpdateCat updating record about cat
func (s *ServiceCat) UpdateCat(c context.Context, cat *model.Cat) error {
	s.CatMap[cat.ID.String()] = cat
	producer.GenerateEvent("cat", "Update", cat, s.Redis.Client, s.Redis.StreamName)
	return s.Rep.UpdateCat(c, cat)
}

// DeleteCat record about cat
func (s *ServiceCat) DeleteCat(c context.Context, id uuid.UUID) error {
	delete(s.CatMap, id.String())
	producer.GenerateEvent("cat", "Delete", model.Cat{ID: id}, s.Redis.Client, s.Redis.StreamName)
	return s.Rep.DeleteCat(c, id)
}
