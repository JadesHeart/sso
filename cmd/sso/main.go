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
	cfg := config.MustLoad() // загрузка конфига

	log := loadLogger(cfg.Env) // инициализация логера

	log.Info("Start app on port:", slog.Any("grpc", cfg.GRPC.Port))

	application := app.New(log, cfg.GRPC.Port, cfg.StoragePath, cfg.TokenTTL) // инициализация нового auth gRPC сервера

	//запуск auth gRPC сервера, если не выходит, то паникуем;
	//запускаем его внутри отдельной горутины, для grace full shut down
	go application.GRPCServer.MustRun()

	// сделано так, чтобы ОС дала сигнал нашей программе, что пора завершаться и она сделает это сама
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT) // функция Notify при получении сигнала запишет его в канал

	// ждём пока что-то придёт в канал, а значит блокируем код дальше;
	// в это время горутина обрабатывает gRPC запросы
	sign := <-stop

	log.Info("Остановка работы приложения", slog.String("сигнал ОС", sign.String()))

	// TODO: сделать такой же grace full shot down для базы данных, в виде приложения
	application.GRPCServer.Stop()

	log.Info("Работа приложения была завершена")
	// TODO: запустить gRPC-сервер приложения
}

// loadLogger инициализирует логер в зависимости от переменной пришедшей в аргументе
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
