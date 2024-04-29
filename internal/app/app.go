package app

import (
	"errors"

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

func (a *App) GetTokens(userGUID string) (*model.Tokens, error) {
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

	tokens := model.Tokens{
		AccessToken:  accessToken,
		RefreshToken: refToken,
	}

	a.storage.SaveTokens(&model.TokenModel{
		UserGUID:     userGUID,
		AccessToken:  accessToken,
		RefreshToken: refToken,
	})

	return &tokens, nil
}

func (a *App) RefreshToken() (*model.Tokens, error) {
	return nil, nil
}
