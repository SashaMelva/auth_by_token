package hendler

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"sync"
	"time"

	"github.com/SashaMelva/auth_by_token/internal/app"
	"github.com/SashaMelva/auth_by_token/storage/model"
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

	if req.Method == http.MethodPost {
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
	var tokens model.TokenModel

	body, err := io.ReadAll(req.Body)

	if err != nil {

		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
	} else {
		err = json.Unmarshal(body, &tokens)
		if err != nil {

			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Content-Type", "application/json; charset=utf-8")
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
		}
	}

	newTokens, err := s.app.RefreshToken(&tokens, ctx)

	if err != nil {
		returnError(&ErrorResponseBody{
			Status:  http.StatusBadRequest,
			Message: []byte(err.Error()),
		}, w)
		return
	}

	json, err := json.Marshal(newTokens)

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
