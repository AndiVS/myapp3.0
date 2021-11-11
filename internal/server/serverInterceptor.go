package server

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"log"
	"myapp3.0/internal/service"
)

// AuthInterceptor is a server interceptor for authentication and authorization
type AuthInterceptor struct {
	access            *service.JWTManager
	refresh           *service.JWTManager
	accessibleForUser map[string]bool
}

// NewAuthInterceptor returns a new auth interceptor
func NewAuthInterceptor(access, refresh *service.JWTManager, accessibleForUser map[string]bool) *AuthInterceptor {
	return &AuthInterceptor{access: access, refresh: refresh, accessibleForUser: accessibleForUser}
}

// Unary returns a server interceptor function to authenticate and authorize unary RPC
func (interceptor *AuthInterceptor) Unary() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		log.Println("--> unary interceptor: ", info.FullMethod)

		err := interceptor.authorize(ctx, info.FullMethod)
		if err != nil {
			return nil, err
		}

		return handler(ctx, req)
	}
}

// Stream returns a server interceptor function to authenticate and authorize stream RPC
func (interceptor *AuthInterceptor) Stream() grpc.StreamServerInterceptor {
	return func(
		srv interface{},
		stream grpc.ServerStream,
		info *grpc.StreamServerInfo,
		handler grpc.StreamHandler,
	) error {
		log.Println("--> stream interceptor: ", info.FullMethod)

		err := interceptor.authorize(stream.Context(), info.FullMethod)
		if err != nil {
			return err
		}

		return handler(srv, stream)
	}
}

func (interceptor *AuthInterceptor) authorize(ctx context.Context, method string) error {

	if method == "/proto.AuthService/SignIn" || method == "/proto.AuthService/SignUp" {
		return nil
	}

	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return status.Errorf(codes.Unauthenticated, "metadata is not provided")
	}

	values := md["access"]
	if len(values) == 0 {
		return status.Errorf(codes.Unauthenticated, "authorization token is not provided")
	}
	accessToken := values[0]
	//claims, err := interceptor.access.Verify(accessToken)
	_, err := interceptor.access.Verify(accessToken)
	if err != nil {
		return status.Errorf(codes.Unauthenticated, "access token is invalid: %v", err)
	}

	/*values = md["refresh"]
		if len(values) == 0 {
			return status.Errorf(codes.Unauthenticated, "authorization token is not provided")
		}
		refreshToken := values[0]
		_, err = interceptor.access.Verify(refreshToken)

		accessibleForUser:= interceptor.accessibleForUser[method]
	   //  ctx.Set("IsAdmin", true)
		if !claims.IsAdmin && !accessibleForUser {
			return status.Error(codes.PermissionDenied, "no permission to access this RPC")
		}
	*/

	return nil
}

/*
func (interceptor *AuthInterceptor) refreshToke(ctx context.Context, method string) error {
	accessibleForUser:= interceptor.accessibleForUser[method]

	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return status.Errorf(codes.Unauthenticated, "metadata is not provided")
	}

	values := md["access"]
	if len(values) == 0 {
		return status.Errorf(codes.Unauthenticated, "authorization token is not provided")
	}
	accessToken := values[0]
	claims, err := interceptor.access.Verify(accessToken)
	if err != nil {
		return status.Errorf(codes.Unauthenticated, "access token is invalid: %v", err)
	}
	values = md["refresh"]
	if len(values) == 0 {
		return status.Errorf(codes.Unauthenticated, "authorization token is not provided")
	}



	if !claims.IsAdmin && !accessibleForUser {
		return status.Error(codes.PermissionDenied, "no permission to access this RPC")
	}

	return nil
}
*/
