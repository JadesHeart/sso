package app

import (
	"log/slog"
	grpcapp "sso/internal/app/grpc"
	"sso/internal/services/auth"
	postgres "sso/internal/storage/postgres"
	"time"
)

type StorageApp struct {
	Storage *postgres.Storage
}

type App struct {
	GRPCServer *grpcapp.App
	StorageApp *StorageApp
}

func New(log *slog.Logger, grpcPort int, storagePath string, tokenTTL time.Duration) *App {
	storage, err := postgres.New(storagePath)
	if err != nil {
		panic(err)
	}

	authService := auth.New(log, storage, storage, storage, tokenTTL)

	grpcApp := grpcapp.New(log, *authService, grpcPort)
	return &App{
		GRPCServer: grpcApp,
		StorageApp: &StorageApp{
			Storage: storage,
		},
	}
}
