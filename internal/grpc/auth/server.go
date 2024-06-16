package auth

import (
	"context"
	ssov1 "github.com/JadesHeart/protos/gen/go/sso"
	"google.golang.org/grpc"
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
	panic("implement me ") // пока что ставим заглужки, чтобы нормально имплементировать методы для удобного дебага
}

// Register todo:описать метод
func (s *serverAPI) Register(ctx context.Context, req *ssov1.RegisterRequest) (*ssov1.RegisterResponse, error) {
	panic("implement me ") // пока что ставим заглужки, чтобы нормально имплементировать методы для удобного дебага
}

// IsAdmin todo:описать метод
func (s *serverAPI) IsAdmin(ctx context.Context, req *ssov1.IsAdminRequest) (*ssov1.IsAdminResponse, error) {
	panic("implement me ") // пока что ставим заглужки, чтобы нормально имплементировать методы для удобного дебага
}
