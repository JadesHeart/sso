package tests

import (
	ssov1 "github.com/JadesHeart/protos/gen/go/sso"
	"github.com/brianvoe/gofakeit"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"sso/internal/lib/validator"
	"sso/tests/suite"
	"testing"
	"time"
)

const (
	emptyAppID = 0
	appID      = 1
	appSecret  = "test-secret"

	passDefaultLen = 10
)

const deltaSecond = 1

func Test_Register_Login_Happy_Path(t *testing.T) {
	ctx, st := suite.New(t)

	email := gofakeit.Email()

	password := randomFakePassword(false, false)

	regResponse, err := st.AuthClient.Register(ctx, &ssov1.RegisterRequest{
		Email:    email,
		Password: password,
	})
	require.NoError(t, err)
	assert.NotEmpty(t, regResponse.GetUserId())

	loginResponse, err := st.AuthClient.Login(ctx, &ssov1.LoginRequest{
		Email:    email,
		Password: password,
		AppId:    appID,
	})
	require.NoError(t, err)

	loginTime := time.Now()

	token := loginResponse.GetToken()
	require.NotEmpty(t, token)

	parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return []byte(appSecret), nil
	})
	require.NoError(t, err)

	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	assert.True(t, ok)

	assert.Equal(t, regResponse.GetUserId(), int64(claims["uid"].(float64)))
	assert.Equal(t, email, claims["email"].(string))
	assert.Equal(t, appID, int(claims["app_id"].(float64)))

	assert.InDelta(t, loginTime.Add(st.Cfg.TokenTTL).Unix(), claims["exp"].(float64), deltaSecond)
}

func Test_Register_Invalid_Data(t *testing.T) {
	ctx, st := suite.New(t)

	tests := []struct {
		name          string
		email         string
		password      string
		expectedError string
	}{
		{
			name:          "Пустая почта",
			email:         " ",
			password:      randomFakePassword(false, false),
			expectedError: validator.EmailInvalid,
		},
		{
			name:          "Пустой пароль",
			email:         gofakeit.Email(),
			password:      " ",
			expectedError: validator.PasswordInvalid,
		},
		{
			name:          "Пароль со специальными символами",
			email:         gofakeit.Email(),
			password:      randomFakePassword(true, false),
			expectedError: validator.PasswordInvalid,
		},
		{
			name:          "Пароль с пробелами",
			email:         gofakeit.Email(),
			password:      randomFakePassword(false, true),
			expectedError: validator.PasswordInvalid,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := st.AuthClient.Register(ctx,
				&ssov1.RegisterRequest{
					Email:    tt.email,
					Password: tt.password,
				})
			require.Error(t, err)
			require.Contains(t, err.Error(), tt.expectedError)
		})
	}

}

func Test_Login_Invalid_Data(t *testing.T) {
	ctx, st := suite.New(t)

	tests := []struct {
		name          string
		email         string
		password      string
		appId         int32
		expectedError string
	}{
		{
			name:          "Пустая почта",
			email:         " ",
			password:      randomFakePassword(false, false),
			appId:         appID,
			expectedError: validator.EmailInvalid,
		},
		{
			name:          "Попытка входа несуществующего юзера",
			email:         gofakeit.Email(),
			password:      randomFakePassword(false, false),
			appId:         appID,
			expectedError: "не верный логин или пароль",
		},
		{
			name:          "Пустой пароль",
			email:         gofakeit.Email(),
			password:      " ",
			appId:         appID,
			expectedError: validator.PasswordInvalid,
		},
		{
			name:          "Пароль со специальными символами",
			email:         gofakeit.Email(),
			password:      randomFakePassword(true, false),
			appId:         appID,
			expectedError: validator.PasswordInvalid,
		},
		{
			name:          "Пароль с пробелами",
			email:         gofakeit.Email(),
			password:      randomFakePassword(false, true),
			appId:         appID,
			expectedError: validator.PasswordInvalid,
		},
		{
			name:          "Несуществующий app_id",
			email:         gofakeit.Email(),
			password:      randomFakePassword(false, false),
			appId:         emptyAppID,
			expectedError: validator.AppIdInvalid,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := st.AuthClient.Register(ctx, &ssov1.RegisterRequest{
				Email:    gofakeit.Email(),
				Password: randomFakePassword(false, false),
			})
			require.NoError(t, err)

			_, err = st.AuthClient.Login(ctx, &ssov1.LoginRequest{
				Email:    tt.email,
				Password: tt.password,
				AppId:    tt.appId,
			})
			require.Error(t, err)
			require.Contains(t, err.Error(), tt.expectedError)
		})
	}

}

func randomFakePassword(special, space bool) string {
	return gofakeit.Password(true, true, true, special, space, passDefaultLen)
}
