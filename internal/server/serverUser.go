package server

import (
	"context"
	"myapp3.0/internal/model"
	"myapp3.0/internal/service"
	"myapp3.0/protocol"
)

type UserServer struct {
	Service service.Users
	protocol.UnimplementedUserServiceServer
}

//NewCatServer
func NewUserServer(Service service.Users) *UserServer {
	return &UserServer{Service: Service}
}

// SearchUser provides cats
func (s *UserServer) SearchUser(ctx context.Context, in *protocol.SearchUserRequest) (*protocol.SearchUserResponse, error) {
	//ctx.Get("sdf")
	r, err := s.Service.GetUser(ctx, in.Username)
	if err != nil {
		return nil, err
	}

	rec := protocol.User{Username: r.Username, Password: r.Password, IsAdmin: r.IsAdmin}

	return &protocol.SearchUserResponse{User: &rec}, nil
}

// GetAllUser provides all cats
func (s *UserServer) GetAllUser(ctx context.Context, in *protocol.GetAllUserRequest) (*protocol.GetAllUserResponse, error) {
	var rec []*model.User

	rec, err := s.Service.GetAllUser(ctx)

	if err != nil {
		return nil, err
	}
	var respsl []*protocol.User

	for i := 0; i < len(rec); i++ {
		respsl[i] = &protocol.User{Username: rec[i].Username, Password: rec[i].Password, IsAdmin: rec[i].IsAdmin}
	}

	return &protocol.GetAllUserResponse{User: respsl}, nil
}

// UpdateUser updating User about cat
func (s *UserServer) UpdateUser(ctx context.Context, in *protocol.UpdateUserRequest) (*protocol.UpdateUserResponse, error) {
	err := s.Service.UpdateUser(ctx, in.User.Username, in.User.IsAdmin)
	if err != nil {
		return nil, err
	}

	return &protocol.UpdateUserResponse{Err: "sec"}, nil
}

// DeleteUser User about cat
func (s *UserServer) DeleteUser(ctx context.Context, in *protocol.DeleteUserRequest) (*protocol.DeleteUserResponse, error) {
	err := s.Service.DeleteUser(ctx, in.Username)
	if err != nil {
		return nil, err
	}

	return &protocol.DeleteUserResponse{Err: " sec "}, nil
}
