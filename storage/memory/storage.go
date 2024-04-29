package memory

import (
	"context"
	"sync"

	"github.com/SashaMelva/auth_by_token/internal/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
)

type Storage struct {
	Logger       *zap.SugaredLogger
	Ctx          context.Context
	ClientMongo  *mongo.Client
	DataBaseName string
	sync.RWMutex
}

func New(client *mongo.Client, log *zap.SugaredLogger, conf *config.ConfigDB) *Storage {
	return &Storage{
		Logger:       log,
		ClientMongo:  client,
		DataBaseName: conf.NameDB,
	}
}

// func (s *Storage) GetCollection() {
// 	collection := s.ClientMongo.Database(s.DataBaseName).Collection("test")
// }
