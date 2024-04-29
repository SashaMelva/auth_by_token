package memory

import (
	"context"
	"sync"

	"github.com/SashaMelva/auth_by_token/internal/config"
	"github.com/SashaMelva/auth_by_token/storage/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
)

type Storage struct {
	Logger       *zap.SugaredLogger
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

func (s *Storage) SaveTokens(tokens *model.RefreshToken, ctx context.Context) error {
	s.Logger.Debug(tokens)
	collection := s.ClientMongo.Database(s.DataBaseName).Collection("test")
	res, err := collection.InsertOne(ctx, tokens)

	if err != nil {
		return err
	}

	s.Logger.Debug(res)
	return nil
}
func (s *Storage) GetTokenByUser(userGUID string, ctx context.Context) (*model.RefreshToken, error) {
	var token model.RefreshToken
	filter := bson.D{{"userguid", userGUID}}

	collection := s.ClientMongo.Database(s.DataBaseName).Collection("test")
	err := collection.FindOne(ctx, filter).Decode(&token)
	s.Logger.Debug(err)
	if err != nil {
		return nil, nil
	}

	return &token, nil
}

func (s *Storage) UpdateTokenByUser(reefToken model.RefreshToken, ctx context.Context) error {
	filter := bson.D{{"userguid", reefToken.UserGUID}}

	update := bson.D{
		{"$inc", bson.D{
			{"refreshtoken", reefToken.RefreshToken},
		}},
	}

	collection := s.ClientMongo.Database(s.DataBaseName).Collection("test")
	res, err := collection.UpdateOne(ctx, filter, update)

	if err != nil {
		return err
	}

	s.Logger.Debug(res)
	return nil
}
