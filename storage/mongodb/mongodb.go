package mongodb

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/Dimitriy14/shortener/encoder"
	"github.com/Dimitriy14/shortener/models"
)

type Repository struct {
	Uri        string
	DB         string
	Collection string
}

func NewMongoRepository(uri, db, collection string) *Repository {
	return &Repository{Uri: uri, DB: db, Collection: collection}
}

func NewMongoClient(uri string) (*mongo.Client, error) {
	client, err := mongo.NewClient(options.Client().ApplyURI(uri))

	if err != nil {
		return nil, err
	}

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	if err = client.Connect(ctx); err != nil {
		return nil, err
	}

	return client, nil
}

func (r Repository) SaveUrl(url string) (string, error) {
	client, err := NewMongoClient(r.Uri)

	if err != nil {
		return "", err
	}

	collection := client.Database(r.DB).Collection(r.Collection)

	id, err := collection.EstimatedDocumentCount(context.TODO())

	if err != nil {
		return "", err
	}

	id++

	genUrl := encoder.Encode(id)

	insertResult, err := collection.InsertOne(context.TODO(), models.Shortener{ID: id, OriginalURL: url, GeneratedURL: genUrl, Visited: false, Count: 0})

	if err != nil {
		return "", err
	}

	log.Println("Insert result: ", insertResult.InsertedID)

	return genUrl, nil
}

func (r Repository) GetURL(code string) (string, error) {
	client, err := NewMongoClient(r.Uri)

	if err != nil {
		return "", err
	}

	filter := bson.D{{"generatedurl", code}}

	update := bson.D{
		{"$inc", bson.D{
			{"count", 1},
		}},
		{"$set", bson.D{
			{"visited", true},
		}},
	}

	var shortener models.Shortener

	if err := client.Database(r.DB).Collection(r.Collection).FindOneAndUpdate(context.TODO(), filter, update).Decode(&shortener); err != nil {
		return "", err
	}

	return shortener.OriginalURL, nil
}
