package context

import (
	"context"
	"github.com/go-playground/validator"
	"google.golang.org/grpc"
)

const (
	validatorKey = "validatorKey"
)

// UnaryServerInterceptor middleware для добавления значений в context
func UnaryServerInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	// Добавление значений в контекст
	newCtx := WithValidator(ctx)

	// Вызов следующего обработчика с новым контекстом
	return handler(newCtx, req)
}

// WithValidator добавляем в контекст validator, чтобы использовать его в ручках
func WithValidator(ctx context.Context) context.Context {
	v := validator.New()
	return context.WithValue(ctx, validatorKey, v)
}
