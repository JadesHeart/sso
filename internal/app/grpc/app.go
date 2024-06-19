package grpcapp

import (
	"fmt"
	"google.golang.org/grpc"
	"log/slog"
	"net"
	authgrpc "sso/internal/grpc/auth"
	"sso/internal/interceptors"
	"sso/internal/services/auth"
)

const (
	opRun  = "grpcApp.Run"
	opStop = "grpcApp.Stop"
)

type App struct {
	log        *slog.Logger
	gRPCServer *grpc.Server
	port       int
}

func New(
	log *slog.Logger,
	authService auth.Auth,
	port int,

) *App {
	gRPCServer := grpc.NewServer(
		grpc.UnaryInterceptor(interceptors.UnaryServerInterceptorValidate),
	)
	authgrpc.Register(gRPCServer, &authService)
	return &App{
		log:        log,
		gRPCServer: gRPCServer,
		port:       port,
	}
}

func (a *App) MustRun() {
	if err := a.Run(); err != nil {
		panic(err)
	}
}

func (a *App) Run() error {
	log := a.log.With(slog.String("op", opRun), slog.Int("port", a.port))

	l, err := net.Listen("tcp", fmt.Sprintf(":%d", a.port))
	if err != nil {
		return fmt.Errorf("%s: %w", opRun, err)
	}

	log.Info("gRPC сервер запущен: ", slog.String("address: ", l.Addr().String()))

	if err := a.gRPCServer.Serve(l); err != nil {
		return fmt.Errorf("%s: %w", opRun, err)

	}

	return nil
}

func (a *App) Stop() {
	a.log.With(slog.String("op", opStop)).
		Info("Остановка gRPC сервера", slog.Int("port", a.port))

	a.gRPCServer.GracefulStop() // прекращает приём новых запросов и ждёт пока завершаться старые
}
