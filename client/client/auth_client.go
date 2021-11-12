package client

import (
	"context"

	"github.com/AndiVS/myapp3.0/protocol"
	"google.golang.org/grpc"

	"time"
)

// AuthClient is a client to call authentication RPC
type AuthClient struct {
	service  protocol.AuthServiceClient
	username string
	password string
}

// NewAuthClient returns a new auth client
func NewAuthClient(cc *grpc.ClientConn, username, password string) *AuthClient {
	service := protocol.NewAuthServiceClient(cc)
	return &AuthClient{service, username, password}
}

// SignIn login user and returns the access token
func (client *AuthClient) SignIn() (string, string, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Second)
	defer cancel()

	req := &protocol.SignInRequest{
		Username: client.username,
		Password: client.password,
	}

	res, err := client.service.SignIn(ctx, req)
	if err != nil {
		return "", "", err
	}

	return res.GetAccessToken(), res.GetRefreshToken(), nil
}
