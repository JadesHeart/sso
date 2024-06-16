package grpcapp

import (
	"fmt"
	"google.golang.org/grpc"
	"log/slog"
	"net"
	authgrpc "sso/internal/grpc/auth"
)

const (
	opRun  = "grpcApp.Run"
	opStop = "grpcApp.Stop"
)

// App TODO: описать структуру
type App struct {
	log        *slog.Logger
	gRPCServer *grpc.Server
	port       int
}

// New конструктор сервера нового gRPC приложения
func New(log *slog.Logger, port int) *App {
	gRPCServer := grpc.NewServer()
	authgrpc.Register(gRPCServer)
	return &App{
		log:        log,
		gRPCServer: gRPCServer,
		port:       port,
	}
}

// MustRun в случае НЕзапуска сервера паникует
func (a *App) MustRun() {
	if err := a.Run(); err != nil {
		panic(err)
	}
}

// Run запускает сервер
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

// Stop останавливает gRPC сервер
func (a *App) Stop() {
	a.log.With(slog.String("op", opStop)).
		Info("Остановка gRPC сервера", slog.Int("port", a.port))

	a.gRPCServer.GracefulStop() // прекращает приём новых запросов и ждёт пока завершаться старые
}
