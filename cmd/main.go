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
)

func main() {

	config := config.New("../files/conf/")
	log := logger.New(config.Logger, "../files/log/")

	connectionDB := connection.New(config.DataBase, log)

	memstorage := memory.New(connectionDB.StorageDb, log)
	app := app.New(log, memstorage)

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
