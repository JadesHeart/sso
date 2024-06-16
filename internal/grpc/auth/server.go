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

type Auth interface {
	Login(ctx context.Context, email string, password string, appID int) (token string, err error)
	RegisterNewUser(ctx context.Context, email string, password string) (userID int64, err error)
	IsAdmin(ctx context.Context, userID int64) (bool, error)
}

// ServerAPI структура обрабатывающая входящие запросы
// TODO: получше описать что она делает
type serverAPI struct {
	ssov1.UnimplementedAuthServer
	auth Auth
}

// Register регистрирует обработчик, чтобы он обрабатывал запросы, которые поступают в gRPC сервер
func Register(grpc *grpc.Server, auth Auth) {
	ssov1.RegisterAuthServer(grpc, &serverAPI{auth: auth})
}

// Login ручка метода Login, валидирует данные и передаёт их в сервисный слой, для логирования юзера
func (s *serverAPI) Login(ctx context.Context, req *ssov1.LoginRequest) (*ssov1.LoginResponse, error) {
	// достаём валидатор из контекста
	v, ok := ctx.Value(validatorKey).(*validator.Validate)
	if !ok {
		return nil, status.Error(codes.Internal, "внутренняя ошибка")
	}
	// Проверяем email на валидность
	if err := v.Var(req.GetEmail(), "required,email"); err != nil {
		return nil, status.Error(codes.InvalidArgument, "email не валидна")
	}
	// Проверяем AppId на валидность
	if req.GetAppId() == emptyValue {
		return nil, status.Error(codes.InvalidArgument, "app_id не валиден")
	}
	// вызываем Login из сервисного слоя
	token, err := s.auth.Login(ctx, req.GetEmail(), req.GetPassword(), int(req.GetAppId()))
	if err != nil {
		// TODO: добавить обработку некоторых ошибок после реализации сервистного слоя
		return nil, status.Error(codes.Internal, "внутренняя ошибка")
	}

	return &ssov1.LoginResponse{Token: token}, nil
}

// Register ручка метода Register, валидирует данные и передаёт их в сервисный слой, для дальнейшей регистрации юзера
func (s *serverAPI) Register(ctx context.Context, req *ssov1.RegisterRequest) (*ssov1.RegisterResponse, error) {
	// достаём валидатор из контекста
	v, ok := ctx.Value(validatorKey).(*validator.Validate)
	if !ok {
		return nil, status.Error(codes.Internal, "внутренняя ошибка")
	}
	// Проверяем email на валидность
	if err := v.Var(req.GetEmail(), "required,email"); err != nil {
		return nil, status.Error(codes.InvalidArgument, "email не валидна")
	}
	// вызываем RegisterNewUser из сервисного слоя
	userID, err := s.auth.RegisterNewUser(ctx, req.GetEmail(), req.GetPassword())
	if err != nil {
		// TODO: добавить обработку некоторых после реализации сервистного слоя
		return nil, status.Error(codes.Internal, "внутренея ошибка")
	}

	return &ssov1.RegisterResponse{UserId: userID}, nil
}

// IsAdmin ручка метода IsAdmin, валидирует данные и передаёт их в сервисный слой, чтобы понять админ ли юзер или нет
func (s *serverAPI) IsAdmin(ctx context.Context, req *ssov1.IsAdminRequest) (*ssov1.IsAdminResponse, error) {
	// вызываем IsAdmin из сервисного слоя
	isAdmin, err := s.auth.IsAdmin(ctx, req.GetUserId())
	if err != nil {
		// TODO: добавить обработку некоторых после реализации сервистного слоя
		return nil, status.Error(codes.Internal, "внутренея ошибка")
	}
	return &ssov1.IsAdminResponse{IsAdmin: isAdmin}, nil
}
