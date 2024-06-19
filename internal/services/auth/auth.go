package auth

import (
	"context"
	"errors"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"log/slog"
	"sso/internal/domain/models"
	"sso/internal/lib/jwt"
	"sso/internal/lib/sl"
	"sso/internal/storage"
	"time"
)

const (
	regUserPath = "auth.RegisterNewUser"
	LoginPath   = "auth.Login"
	isAdminPath = "auth.IsAdmin"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrInvalidAppId       = errors.New("invalid app id")
	UserAlreadyExist      = errors.New("user already exist")
)

type Auth struct {
	log         *slog.Logger
	usrSaver    UserSaver
	usrProvider UserProvider
	appProvider AppProvider
	tokenTTL    time.Duration
}

type UserSaver interface {
	SaveUser(
		ctx context.Context,
		email string,
		passHash []byte,
	) (uid int64, err error)
}

type UserProvider interface {
	User(ctx context.Context, email string) (models.User, error)
	IsAdmin(ctx context.Context, userID int64) (bool, error)
}

type AppProvider interface {
	App(ctx context.Context, appID int) (models.App, error)
}

func New(
	log *slog.Logger,
	userSaver UserSaver,
	userProvider UserProvider,
	appProvider AppProvider,
	tokenTTL time.Duration,
) *Auth {
	return &Auth{
		log:         log,
		usrSaver:    userSaver,
		usrProvider: userProvider,
		appProvider: appProvider,
		tokenTTL:    tokenTTL,
	}
}

func (a *Auth) Login(
	ctx context.Context,
	email string,
	password string,
	appID int,
) (string, error) {
	log := a.log.With(
		slog.Any("op", LoginPath),
	)

	user, err := a.usrProvider.User(ctx, email)
	if err != nil {
		if errors.Is(err, storage.ErrUserNotFound) {
			log.Warn("пользователь не найден", sl.Err(err))

			return "", fmt.Errorf("%s: %w", LoginPath, ErrInvalidCredentials)
		}
		log.Error("не удалось получить юзера", sl.Err(err))

		return "", fmt.Errorf("%s: %w", LoginPath, err)
	}

	if err := bcrypt.CompareHashAndPassword(user.PassHash, []byte(password)); err != nil {
		log.Info("неверные учётные данные", sl.Err(err))

		return "", fmt.Errorf("%s: %w", LoginPath, ErrInvalidCredentials)
	}

	app, err := a.appProvider.App(ctx, appID)
	if err != nil {
		return "", fmt.Errorf("%s: %w", LoginPath, err)
	}

	token, err := jwt.NewToken(user, app, a.tokenTTL)
	if err != nil {
		log.Error("не удалось сгенерировать токен", sl.Err(err))

		return "", fmt.Errorf("%s: %w", LoginPath, err)
	}

	return token, nil
}

func (a *Auth) RegisterNewUser(
	ctx context.Context,
	email string,
	password string,
) (int64, error) {
	log := a.log.With(
		slog.Any("op", regUserPath),
	)

	passHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Error("не удалось сгенерировать хэш от пароля", sl.Err(err))

		return 0, fmt.Errorf("%s: %w", regUserPath, err)
	}

	id, err := a.usrSaver.SaveUser(ctx, email, passHash)
	if err != nil {
		if errors.Is(err, storage.ErrUserExists) {
			log.Warn("пользователь уже существует", sl.Err(err))

			return 0, fmt.Errorf("%s: %w", regUserPath, UserAlreadyExist)
		}

		log.Error("не удалось сохранить юзера", sl.Err(err))

		return 0, fmt.Errorf("%s: %w", regUserPath, err)
	}

	return id, nil
}

func (a *Auth) IsAdmin(ctx context.Context, userID int) (bool, error) {
	log := a.log.With(
		slog.Any("op", isAdminPath),
	)

	isAdmin, err := a.usrProvider.IsAdmin(ctx, int64(userID))
	if err != nil {
		if errors.Is(err, storage.ErrAppNotFound) {
			log.Warn("пользователь не найден", sl.Err(err))

			return false, fmt.Errorf("%s: %w", isAdminPath, ErrInvalidAppId)
		}
		return false, fmt.Errorf("%s: %w", isAdminPath, err)
	}

	log.Warn("проверяем админ ли юзер: ", slog.Bool("is_admin", isAdmin), userID)

	return isAdmin, nil
}
