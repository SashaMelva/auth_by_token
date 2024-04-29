package connection

import (
	"context"

	"github.com/SashaMelva/auth_by_token/internal/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
)

func New(conf *config.ConfigDB, log *zap.SugaredLogger) *mongo.Client {
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://" + conf.Host + ":" + conf.Port))
	if err != nil {
		log.Fatal(err)
	}

	err = client.Connect(context.TODO())
	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
	}

	log.Debug("Create connection mongo")

	return client
}
