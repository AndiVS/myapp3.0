package server

import (
	"context"
	"github.com/google/uuid"
	"myapp3.0/internal/model"
	"myapp3.0/internal/service"
	"myapp3.0/protocol"
)

type CatServer struct {
	Service service.Cats
	protocol.UnimplementedCatServiceServer
}

//NewCatServer as
func NewCatServer(Service service.Cats) *CatServer {
	return &CatServer{Service: Service}
}

// CreateCat Cat about cat
func (s *CatServer) CreateCat(ctx context.Context, in *protocol.CreateCatRequest) (*protocol.CreateCatResponse, error) {
	rec := model.Cat{Name: in.Cat.Name, Type: in.Cat.Type}
	id, err := s.Service.AddCat(ctx, &rec)
	if err != nil {
		return nil, err
	}
	return &protocol.CreateCatResponse{Id: "id.String() " + id.String()}, nil
}

// SearchCat provides cat
func (s *CatServer) SearchCat(ctx context.Context, in *protocol.SearchCatRequest) (*protocol.SearchCatResponse, error) {
	_id, err1 := uuid.Parse(in.Id)
	if err1 != nil {
		return &protocol.SearchCatResponse{}, err1
	}

	r, err := s.Service.GetCat(ctx, _id)
	if err != nil {
		return nil, err
	}
	rec := protocol.Cat{Name: r.Name, Id: r.ID.String(), Type: r.Type}

	return &protocol.SearchCatResponse{Cat: &rec}, nil
}

// GetAllCat provides all cats
func (s *CatServer) GetAllCat(ctx context.Context, in *protocol.GetAllCatRequest) (*protocol.GetAllCatResponse, error) {
	var rec []*model.Cat

	rec, err := s.Service.GetAllCat(ctx)

	if err != nil {
		return nil, err
	}
	var protorecsl []*protocol.Cat

	for i := 0; i < len(rec); i++ {
		protorecsl[i] = &protocol.Cat{Id: rec[i].ID.String(), Name: rec[i].Name, Type: rec[i].Type}
	}

	return &protocol.GetAllCatResponse{Cats: protorecsl}, nil
}

// UpdateCat updating Cat about cat
func (s *CatServer) UpdateCat(ctx context.Context, in *protocol.UpdateCatRequest) (*protocol.UpdateCatResponse, error) {
	_id, err1 := uuid.Parse(in.Cat.Id)
	if err1 != nil {
		return &protocol.UpdateCatResponse{}, err1
	}
	rec := model.Cat{ID: _id, Name: in.Cat.Name, Type: in.Cat.Type}

	err := s.Service.UpdateCat(ctx, &rec)
	if err != nil {
		return nil, err
	}

	return &protocol.UpdateCatResponse{Err: "sec"}, nil
}

// DeleteCat Cat about cat
func (s *CatServer) DeleteCat(ctx context.Context, in *protocol.DeleteCatRequest) (*protocol.DeleteCatResponse, error) {
	_id, err1 := uuid.Parse(in.Id)
	if err1 != nil {
		return &protocol.DeleteCatResponse{}, err1
	}

	err := s.Service.DeleteCat(ctx, _id)
	if err != nil {
		return nil, err
	}

	return &protocol.DeleteCatResponse{Err: " sec "}, nil
}
