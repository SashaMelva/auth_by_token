package app

import (
	"github.com/SashaMelva/auth_by_token/storage/memory"
	"github.com/SashaMelva/auth_by_token/storage/model"
	"go.uber.org/zap"
)

type App struct {
	storage *memory.Storage
	Logger  *zap.SugaredLogger
}

func New(logger *zap.SugaredLogger, storage *memory.Storage, host *string) *App {
	return &App{
		storage: storage,
		Logger:  logger,
	}
}

func (a *App) GetTokens(userGUID string) (*model.Tokens, error) {
	return nil, nil
}

func (a *App) RefreshToken() (*model.Tokens, error) {
	return nil, nil
}
