package main

import (
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
	"time"
)

type MockJWT struct {
	mock.Mock
}

func (m *MockJWT) NewToken(user MockUser, app MockApp, duration time.Duration) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["uid"] = user.ID
	claims["email"] = user.Email
	claims["exp"] = time.Now().Add(duration).Unix()
	claims["app_id"] = app.ID

	tokenString, err := token.SignedString([]byte(app.Secret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func TestNewToken(t *testing.T) {
	t.Run("ValidToken", func(t *testing.T) {
		user := NewMockUser(1, "test@example.com")
		app := NewMockApp(1, "supersecret")
		duration := time.Hour

		mockJWT := new(MockJWT)

		mockJWT.On("NewToken", user, app, duration).Return("mock_token_string", nil)

		tokenString, err := mockJWT.NewToken(user, app, duration)

		assert.NoError(t, err, "Unexpected error from NewToken")
		assert.Equal(t, "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhcHBfaWQiOjEsImVtYWlsIjoidGVzdEBleGFtcGxlLmNvbSIsImV4cCI6MTcxODgzOTAxMywidWlkIjoxfQ", tokenString, "Token string mismatch")
	})

	t.Run("InvalidSecret", func(t *testing.T) {
		user := NewMockUser(1, "test@example.com")
		app := NewMockApp(1, "какой-то секрет") // Неверный секрет
		duration := time.Hour

		mockJWT := new(MockJWT)

		mockJWT.On("NewToken", user, app, duration).Return("", errors.New("signature validation failed"))

		tokenString, err := mockJWT.NewToken(user, app, duration)

		assert.Error(t, err, "Expected error from NewToken")
		assert.Empty(t, tokenString, "Expected empty token string")
	})

	t.Run("ExpiredToken", func(t *testing.T) {
		user := NewMockUser(1, "test@example.com")
		app := NewMockApp(1, "supersecret")
		duration := -time.Hour // Прошедшее время (истекший токен)

		mockJWT := new(MockJWT)

		mockJWT.On("NewToken", user, app, duration).Return("expired_token_string", nil)

		tokenString, err := mockJWT.NewToken(user, app, duration)

		assert.NoError(t, err, "Unexpected error from NewToken")
		assert.NotEmpty(t, tokenString, "Expected non-empty token string")

		expiredToken, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return []byte(app.Secret), nil
		})

		assert.NoError(t, err, "Unexpected error parsing token")
		assert.False(t, expiredToken.Valid, "Expected expired token")
	})

	t.Run("InvalidUserOrApp", func(t *testing.T) {
		invalidUser := NewMockUser(0, "")     // Недопустимый пользователь
		invalidApp := NewMockApp(0, "secret") // Недопустимое приложение
		duration := time.Hour

		mockJWT := new(MockJWT)

		mockJWT.On("NewToken", invalidUser, invalidApp, duration).Return("", errors.New("invalid user or app"))

		tokenString, err := mockJWT.NewToken(invalidUser, invalidApp, duration)

		assert.Error(t, err, "Expected error from NewToken")
		assert.Empty(t, tokenString, "Expected empty token string")
	})
}
