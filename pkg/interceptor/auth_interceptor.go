package interceptor

import (
	"context"

	"github.com/juanmachuca95/ahorcado_go/pkg/servicejwt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type AuthInterceptor struct {
	accessiblesRoles map[string][]string
	servJWT          servicejwt.AuthJWT
}

func NewAuthInterceptor() *AuthInterceptor {
	sJWT := servicejwt.NewServiceJWT()
	return &AuthInterceptor{
		accessiblesRoles: accessibleRoles(),
		servJWT:          sJWT,
	}
}

func (a *AuthInterceptor) UnaryInterceptor() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {

		//log.Println("--> unary interceptor: ", info.FullMethod)
		if err := a.authorize(ctx, info.FullMethod); err != nil {
			return nil, err
		}
		return handler(ctx, req)
	}
}

func (a *AuthInterceptor) StreamInterceptor() grpc.StreamServerInterceptor {
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {

		//log.Println("--> stream interceptor: ", info.FullMethod)
		if err := a.authorize(ss.Context(), info.FullMethod); err != nil {
			return err
		}
		return handler(srv, ss)
	}
}

func (interceptor *AuthInterceptor) authorize(ctx context.Context, method string) error {
	accessibleRoles, ok := interceptor.accessiblesRoles[method]
	if !ok {
		// everyone can access
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
	claims, err := interceptor.servJWT.ValidateToken(accessToken)
	if err != nil {
		return status.Errorf(codes.Unauthenticated, "access token is invalid: %v", err)
	}

	for _, role := range accessibleRoles {
		if role == claims.Role {
			return nil
		}
	}

	return status.Error(codes.PermissionDenied, "no permission to access this RPC")
}

func accessibleRoles() map[string][]string {
	//const gameServicePath = "/protos.Ahorcado/"

	return map[string][]string{
		//gameServicePath + "GetRandomGame": {"player"},
	}
}
