// Package service encapsulates
package service

import (
	"context"

	"github.com/AndiVS/myapp3.0/internal/broker"
	"github.com/AndiVS/myapp3.0/internal/model"
	"github.com/AndiVS/myapp3.0/internal/repository"
	"github.com/google/uuid"
)

// Cats interface for mocks
type Cats interface {
	AddCat(c context.Context, rec *model.Cat) (uuid.UUID, error)
	GetCat(c context.Context, id uuid.UUID) (*model.Cat, error)
	GetAllCat(c context.Context) ([]*model.Cat, error)
	UpdateCat(c context.Context, rec *model.Cat) error
	DeleteCat(c context.Context, id uuid.UUID) error
}

// CatService struct for rep
type CatService struct {
	Rep    repository.Cats
	CatMap map[string]*model.Cat
	Broker broker.Broker
}

// NewServiceCat used for setting services
func NewServiceCat(Rep, Broker interface{}, dbType, brokerType string) Cats {
	serviceCat := CatService{}

	switch dbType {
	case "mongodb":
		serviceCat.Rep = Rep.(*repository.Mongo)
	case "postgres":
		serviceCat.Rep = Rep.(*repository.Postgres)
	}

	catMap := make(map[string]*model.Cat)
	catsSlice, _ := serviceCat.Rep.SelectAllCat(context.Background())
	for _, cat := range catsSlice {
		catMap[cat.ID.String()] = cat
	}
	serviceCat.CatMap = catMap

	switch brokerType {
	case "redis":
		Redis := Broker.(*broker.Redis)
		serviceCat.Broker = Redis
	case "kafka":
		Kafka := Broker.(*broker.Kafka)
		serviceCat.Broker = Kafka
	case "rabbit":
		Rabbit := Broker.(*broker.RabbitMQ)
		serviceCat.Broker = Rabbit
	}

	go serviceCat.Broker.ConsumeEvents(Rep)

	return &serviceCat
}

// AddCat record about cat
func (s *CatService) AddCat(c context.Context, cat *model.Cat) (uuid.UUID, error) {
	s.CatMap[cat.ID.String()] = cat
	s.Broker.ProduceEvent("cat", "Insert", *cat)
	return s.Rep.InsertCat(c, cat)
}

// GetCat provides cat
func (s *CatService) GetCat(c context.Context, id uuid.UUID) (*model.Cat, error) {
	val, ok := s.CatMap[id.String()]
	if ok {
		return val, nil
	}
	return s.Rep.SelectCat(c, id)
}

// GetAllCat provides all cats
func (s *CatService) GetAllCat(c context.Context) ([]*model.Cat, error) {
	return s.Rep.SelectAllCat(c)
}

// UpdateCat updating record about cat
func (s *CatService) UpdateCat(c context.Context, cat *model.Cat) error {
	s.CatMap[cat.ID.String()] = cat
	s.Broker.ProduceEvent("cat", "Update", *cat)
	return s.Rep.UpdateCat(c, cat)
}

// DeleteCat record about cat
func (s *CatService) DeleteCat(c context.Context, id uuid.UUID) error {
	delete(s.CatMap, id.String())
	s.Broker.ProduceEvent("cat", "Delete", &model.Cat{ID: id})
	return s.Rep.DeleteCat(c, id)
}
