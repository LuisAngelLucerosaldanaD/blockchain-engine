package interceptor

import (
	"context"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

var routes = map[string]string{
	"/mine_proto.mineBlockServicesBlocks/GenerateBlockGenesis": "/mine_proto.mineBlockServicesBlocks/GenerateBlockGenesis",
	"/blocks_proto.blockServicesBlocks/GetBlock":               "/blocks_proto.blockServicesBlocks/GetBlock",
	"/blocks_proto.blockServicesBlocks/GetBlockByID":           "/blocks_proto.blockServicesBlocks/GetBlockByID",
}

// AuthInterceptor is a server interceptor for authentication and authorization
type AuthInterceptor struct {
	accessibleRoles map[string][]int
}

// NewAuthInterceptor returns a new auth interceptor
func NewAuthInterceptor() *AuthInterceptor {
	return &AuthInterceptor{}
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
	//TODO revisar manejo de rutas
	if _, ok := routes[method]; ok {
		return nil
	}

	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return status.Errorf(codes.Unauthenticated, "metadata is not provided")
	}

	values := md["authorization"]
	if len(values) == 0 {
		return status.Errorf(codes.Unauthenticated, "authorization token is not provided")
	}

	accessToken := values[0]
	_, err := Verify(accessToken)
	if err != nil {
		return status.Errorf(codes.Unauthenticated, "access token is invalid: %v", err)
	}
	/*
		for _, role := range accessibleRoles {
			if role == claims.Role {
				return nil
			}
		}*/

	return nil
}
