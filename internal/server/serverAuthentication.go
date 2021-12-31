package server

import (
	"context"

	"github.com/AndiVS/myapp3.0/internal/model"
	"github.com/AndiVS/myapp3.0/internal/service"
	"github.com/AndiVS/myapp3.0/protocol"
)

// AuthenticationServer for grcp
type AuthenticationServer struct {
	Service service.Authentication
	*protocol.UnimplementedAuthServiceServer
}

// NewAuthenticationServer returns a new auth server
func NewAuthenticationServer(Service service.Authentication) protocol.AuthServiceServer {
	return &AuthenticationServer{Service: Service}
}

// SignUp User about cat
func (s *AuthenticationServer) SignUp(ctx context.Context, req *protocol.SignUpRequest) (*protocol.SignUpResponse, error) {
	user := model.User{Username: req.Username, Password: req.Password}

	err := s.Service.SignUp(ctx, &user)
	res := &protocol.SignUpResponse{Err: err.Error()}

	return res, nil
}

// SignIn is a unary RPC to login user
func (s *AuthenticationServer) SignIn(ctx context.Context, req *protocol.SignInRequest) (*protocol.SignInResponse, error) {
	user := model.User{Username: req.Username, Password: req.Password}
	accessToken, refreshToken, err := s.Service.SignIn(ctx, &user)
	if err != nil {
		return nil, err
	}

	res := &protocol.SignInResponse{AccessToken: accessToken, RefreshToken: refreshToken}
	return res, nil
}
