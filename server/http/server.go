package http

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/SashaMelva/auth_by_token/internal/app"
	"github.com/SashaMelva/auth_by_token/internal/config"
	"github.com/SashaMelva/auth_by_token/server/hendler"
	"github.com/rs/cors"
	"go.uber.org/zap"
)

type Server struct {
	HttpServer *http.Server
}

func NewServer(log *zap.SugaredLogger, app *app.App, config *config.ConfigHttpServer) *Server {
	log.Info("URL api " + config.Host + ":" + config.Port)
	log.Debug("URL api running " + config.Host + ":" + config.Port)
	timeout := config.Timeout * time.Second

	mux := http.NewServeMux()
	h := hendler.NewService(log, app, timeout)

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Debug("Test path working")
		fmt.Fprintf(w, "Hello World!")
	})

	mux.HandleFunc("/auth", h.AuthHendler)
	mux.HandleFunc("/refresh", h.RefreshHendler)

	cors := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{
			http.MethodPost,
			http.MethodGet,
			http.MethodDelete,
			http.MethodPut,
		},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: false,
	})
	handler := cors.Handler(mux)
	return &Server{
		&http.Server{
			Addr:         config.Host + ":" + config.Port,
			Handler:      handler,
			ReadTimeout:  timeout,
			WriteTimeout: timeout,
		},
	}
}

func (s *Server) Start(ctx context.Context) error {
	err := s.HttpServer.ListenAndServe()
	<-ctx.Done()
	return err
}

func (s *Server) Stop(ctx context.Context) error {
	return s.HttpServer.Shutdown(ctx)
}
