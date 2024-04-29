package cmd

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/SashaMelva/auth_by_token/internal/app"
	"github.com/SashaMelva/auth_by_token/internal/config"
	"github.com/SashaMelva/auth_by_token/internal/logger"
	"github.com/SashaMelva/auth_by_token/server/http"
	"github.com/SashaMelva/auth_by_token/storage/connection"
	"github.com/SashaMelva/auth_by_token/storage/memory"
)

func main() {

	config := config.New("../files/conf/")
	log := logger.New(config.Logger, "../files/log/")

	clientMongo := connection.New(config.DataBase, log)

	memstorage := memory.New(clientMongo, log, config.DataBase)
	app := app.New(log, memstorage, config.SecretJWT)

	httpServer := http.NewServer(log, app, config.HttpServer)

	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	go func() {
		<-ctx.Done()

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()

		if err := httpServer.Stop(ctx); err != nil {
			log.Error("failed to stop http server: " + err.Error())
		}
	}()

	log.Info("Services is running...")
	log.Debug("Debug mode enabled")

	if err := httpServer.Start(ctx); err != nil {
		log.Error("failed to start http server: " + err.Error())
		cancel()
		os.Exit(1)
	}
}
