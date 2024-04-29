package hendler

import (
	"context"
	"encoding/json"
	"net/http"
	"sync"
	"time"

	"github.com/SashaMelva/auth_by_token/internal/app"
	"go.uber.org/zap"
)

type Service struct {
	Logger zap.SugaredLogger
	app    app.App
	sync.RWMutex
}

type ErrorResponseBody struct {
	Status  int
	Message []byte
}

func NewService(log *zap.SugaredLogger, application *app.App, timeout time.Duration) *Service {
	return &Service{
		Logger: *log,
		app:    *application,
	}
}

func returnError(errorResponse *ErrorResponseBody, w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(errorResponse.Status)
	w.Write(errorResponse.Message)
}

func (s *Service) AuthHendler(w http.ResponseWriter, req *http.Request) {
	s.Logger.Debug("Path responce ", req.Method, req.URL)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	if req.Method == http.MethodGet {
		args := req.URL.Query()
		userGUID := args.Get("GUID")
		s.getToken(userGUID, w, req, ctx)
	}
}

func (s *Service) RefreshHendler(w http.ResponseWriter, req *http.Request) {
	s.Logger.Debug("Path responce ", req.Method, req.URL)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	if req.Method == http.MethodGet {
		s.refreshToken(w, req, ctx)
	}
}

func (s *Service) getToken(userGUID string, w http.ResponseWriter, req *http.Request, ctx context.Context) {
	tokens, err := s.app.GetTokens(userGUID, ctx)

	if err != nil {
		returnError(&ErrorResponseBody{
			Status:  http.StatusBadRequest,
			Message: []byte(err.Error()),
		}, w)
		return
	}

	json, err := json.Marshal(tokens)

	if err != nil {
		returnError(&ErrorResponseBody{
			Status:  http.StatusInternalServerError,
			Message: []byte(err.Error()),
		}, w)
		return
	}

	s.Logger.Info("OK")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(json)
}

func (s *Service) refreshToken(w http.ResponseWriter, req *http.Request, ctx context.Context) {
	tokens, err := s.app.RefreshToken()

	if err != nil {
		returnError(&ErrorResponseBody{
			Status:  http.StatusBadRequest,
			Message: []byte(err.Error()),
		}, w)
		return
	}

	json, err := json.Marshal(tokens)

	if err != nil {
		returnError(&ErrorResponseBody{
			Status:  http.StatusInternalServerError,
			Message: []byte(err.Error()),
		}, w)
		return
	}

	s.Logger.Info("OK")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(json)
}
