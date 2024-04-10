package mongo

import (
	"context"
	"time"

	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"

	"github.com/juliocesar1235/golang-hex/shortener"
)

type mongoRepository struct {
	client		*mongo.Client
	database	string
	timeout		time.Duration
}

func newMongoClient(mongoURL string, mongoTimeout int) (*mongo.Client, error) {
	secs := time.Second
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(mongoTimeout)*secs)

	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURL))
	if err != nil {
		return nil, err
	}
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		return nil, err
	}
	return client, nil
}

func NewMongoRepository(mongoURL, mongoDB string, mongoTimeout int) (shortener.RedirectRepository, error) {
	mongoRepo := &mongoRepository{
		timeout: time.Duration(mongoTimeout) * time.Second,
		database: mongoDB,
	}
	client, err := newMongoClient(mongoURL, mongoTimeout)
	if err != nil {
		return nil, errors.Wrap(err, "repository.NewMongoRepo")
	}
	mongoRepo.client = client
	return mongoRepo, nil
}

func (mRepository *mongoRepository) Find(code string) (*shortener.Redirect, error) {
	ctx, cancel := context.WithTimeout(context.Background(), mRepository.timeout)
	defer cancel()
	redirect := &shortener.Redirect{}
	collection := mRepository.client.Database(mRepository.database).Collection("redirects")
	filter := bson.M{"code": code}
	err := collection.FindOne(ctx,filter).Decode(&redirect)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.Wrap(shortener.ErrRedirectNotFound, "repository.Redirect.Find")
		}
		return nil, errors.Wrap(err, "repository.Redirect.Find")
	}
	return redirect, nil
}

func (mRepository *mongoRepository) Store(redirect *shortener.Redirect) error {
	ctx, cancel := context.WithTimeout(context.Background(), mRepository.timeout)
	defer cancel()
	collection := mRepository.client.Database(mRepository.database).Collection("redirects")
	_, err := collection.InsertOne(
		ctx,
		bson.M{
			"code": 			redirect.Code,
			"url":				redirect.URL,
			"created_at":	redirect.CreatedAt,
		},
	)

	if err != nil {
		return errors.Wrap(err, "repository.Redirect.Store")
	}
	return nil
}