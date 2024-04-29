package app

import (
	"context"
	"errors"
	"strings"

	"github.com/SashaMelva/auth_by_token/internal/config"
	"github.com/SashaMelva/auth_by_token/internal/pkg"
	"github.com/SashaMelva/auth_by_token/storage/memory"
	"github.com/SashaMelva/auth_by_token/storage/model"
	"github.com/beevik/guid"
	"go.uber.org/zap"
)

type App struct {
	storage *memory.Storage
	Logger  *zap.SugaredLogger
	Tokens  *config.Tokens
}

func New(logger *zap.SugaredLogger, storage *memory.Storage, conf *config.Tokens) *App {
	return &App{
		storage: storage,
		Logger:  logger,
		Tokens:  conf,
	}
}

func (a *App) GetTokens(userGUID string, ctx context.Context) (*model.Tokens, error) {
	if !guid.IsGuid(userGUID) {
		a.Logger.Error("Не валидный GUID")
		return nil, errors.New("Не валидный GUID")
	}

	token, err := a.storage.GetTokenByUser(userGUID, ctx)

	if err != nil {
		a.Logger.Error(err)
		return nil, err
	}
	if token != nil {
		a.Logger.Debug(token)
		return nil, errors.New("В базе уже существует токен для данного пользователя")
	}

	accessToken, err := pkg.GenerateAccssesToken(userGUID, a.Tokens.SecretJWT)

	if err != nil {
		a.Logger.Error(err)
		return nil, err
	}

	refToken, err := pkg.GenerateRefreshToken()

	if err != nil {
		a.Logger.Error(err)
		return nil, err
	}

	tokens := model.Tokens{
		AccessToken:  accessToken,
		RefreshToken: refToken,
	}

	err = a.storage.SaveTokens(&model.RefreshToken{
		UserGUID:     userGUID,
		RefreshToken: refToken,
	}, ctx)

	if err != nil {
		a.Logger.Error(err)
		return nil, err
	}

	return &tokens, nil
}

func (a *App) RefreshToken(tokens *model.TokenModel, ctx context.Context) (*model.Tokens, error) {
	splitAccess := strings.Split(tokens.AccessToken, " ")
	if len(splitAccess) != 2 {
		return nil, errors.New("Неверный формат Access токена. Пример: Bearer wweqw4eq5e5eq")
	}
	if splitAccess[0] != "Bearer" || splitAccess[1] == "" {
		return nil, errors.New("Access токен не валидный")
	}

	userGUID, err := pkg.ParseAccessToken(splitAccess[1], a.Tokens.SecretJWT)

	if err != nil {
		a.Logger.Error(err)
		return nil, err
	}

	if !guid.IsGuid(userGUID) {
		a.Logger.Error("Не валидный GUID")
		return nil, errors.New("Не валидный GUID")
	}

	accessToken, err := pkg.GenerateAccssesToken(userGUID, a.Tokens.SecretJWT)

	if err != nil {
		a.Logger.Error(err)
		return nil, err
	}

	refToken, err := pkg.GenerateRefreshToken()

	if err != nil {
		a.Logger.Error(err)
		return nil, err
	}

	newTokens := model.Tokens{
		AccessToken:  accessToken,
		RefreshToken: refToken,
	}

	err = a.storage.UpdateTokenByUser(&model.RefreshToken{
		RefreshToken: refToken,
		UserGUID:     userGUID,
	}, ctx)

	if err != nil {
		a.Logger.Error(err)
		return nil, err
	}

	return &newTokens, nil
}
