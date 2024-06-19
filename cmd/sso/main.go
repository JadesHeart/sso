package main

import (
	"log/slog"
	"os"
	"os/signal"
	"sso/internal/app"
	"sso/internal/config"
	"syscall"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	cfg := config.MustLoad()

	log := loadLogger(cfg.Env)

	log.Info("Запуск приложения на порту:", slog.Any("grpc", cfg.GRPC.Port))

	application := app.New(log, cfg.GRPC.Port, cfg.StoragePath, cfg.TokenTTL)

	go application.GRPCServer.MustRun()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	sign := <-stop

	log.Info("Остановка работы приложения", slog.String("сигнал ОС", sign.String()))

	// TODO: сделать такой же grace full shot down для базы данных, в виде приложения
	application.GRPCServer.Stop()

	log.Info("Работа приложения была завершена")

	// TODO: запустить gRPC-сервер приложения
}

func loadLogger(loggerLevel string) *slog.Logger {
	var log *slog.Logger
	switch loggerLevel {
	case envLocal:
		log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envDev:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envProd:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	}

	return log
}
