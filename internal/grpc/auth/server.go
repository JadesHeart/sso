package auth

import (
	"context"
	ssov1 "github.com/JadesHeart/protos/gen/go/sso"
	"github.com/go-playground/validator"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	validatorKey = "validatorKey" // TODO: вынести эту константу куда-то там
	emptyValue   = 0
)

// ServerAPI структура обрабатывающая входящие запросы
// TODO: получше описать что она делает
type serverAPI struct {
	ssov1.UnimplementedAuthServer
}

// Register регистрирует обработчик, чтобы он обрабатывал запросы, которые поступают в gRPC сервер
func Register(grpc *grpc.Server) {
	ssov1.RegisterAuthServer(grpc, &serverAPI{})
}

// Login todo:описать метод
func (s *serverAPI) Login(ctx context.Context, req *ssov1.LoginRequest) (*ssov1.LoginResponse, error) {
	// достаём валидатор из контекста
	v, ok := ctx.Value(validatorKey).(*validator.Validate)
	if !ok {
		return nil, status.Error(codes.Internal, "не удалось получить validator")
	}
	// Проверяем email на валидность
	if err := v.Var(req.GetEmail(), "required,email"); err != nil {
		return nil, status.Error(codes.InvalidArgument, "email не валидна")
	}
	// Проверяем AppId на валидность
	if req.GetAppId() == emptyValue {
		return nil, status.Error(codes.InvalidArgument, "app_id не валиден")
	}

	return &ssov1.LoginResponse{Token: "token123456"}, nil
}

// Register todo:описать метод
func (s *serverAPI) Register(ctx context.Context, req *ssov1.RegisterRequest) (*ssov1.RegisterResponse, error) {
	panic("implement me ") // пока что ставим заглужки, чтобы нормально имплементировать методы для удобного дебага
}

// IsAdmin todo:описать метод
func (s *serverAPI) IsAdmin(ctx context.Context, req *ssov1.IsAdminRequest) (*ssov1.IsAdminResponse, error) {
	panic("implement me ") // пока что ставим заглужки, чтобы нормально имплементировать методы для удобного дебага
}
