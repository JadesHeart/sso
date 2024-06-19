package interceptors

import (
	"context"
	ssov1 "github.com/JadesHeart/protos/gen/go/sso"
	"google.golang.org/grpc"
	"sso/internal/lib/validator/auth"
)

func UnaryServerInterceptorValidate(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (interface{}, error) {
	switch r := req.(type) {

	case *ssov1.LoginRequest:
		if err := auth.ValidateLoginReq(r.GetEmail(), r.GetPassword(), r.GetAppId()); err != nil {
			return nil, err
		}
	case *ssov1.RegisterRequest:
		if err := auth.ValidateRegisterReq(r.GetEmail(), r.GetPassword()); err != nil {
			return nil, err
		}

	}

	return handler(ctx, req)
}
