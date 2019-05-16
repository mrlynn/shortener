package mongodb

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/Dimitriy14/shortener/encoder"
	"github.com/Dimitriy14/shortener/models"
)

type Repository struct {
	DB         string
	Collection string
	Client     *mongo.Client
}

func NewMongoRepository(db, collection string, client *mongo.Client) *Repository {
	return &Repository{DB: db, Collection: collection, Client: client}
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
	collection := r.Client.Database(r.DB).Collection(r.Collection)

	id, err := collection.EstimatedDocumentCount(context.TODO())

	if err != nil {
		return "", err
	}

	id++

	encoded := encoder.Encode(id)

	genUrl := fmt.Sprintf("http://localhost:8080/go/%s", encoded)

	insertResult, err := collection.InsertOne(context.TODO(), models.Shortener{ID: id, OriginalURL: url, GeneratedURL: genUrl, Visited: false, Count: 0})

	if err != nil {
		return "", err
	}

	log.Println("Insert result: ", insertResult.InsertedID)

	return genUrl, nil
}

func (r Repository) GetURL(code string) (string, error) {
	genUrl := fmt.Sprintf("http://localhost:8080/go/%s", code)
	filter := bson.D{{"generatedurl", genUrl}}

	update := bson.D{
		{"$inc", bson.D{
			{"count", 1},
		}},
		{"$set", bson.D{
			{"visited", true},
		}},
	}

	var shortener models.Shortener

	if err := r.Client.Database(r.DB).Collection(r.Collection).FindOneAndUpdate(context.TODO(), filter, update).Decode(&shortener); err != nil {
		return "", err
	}

	return shortener.OriginalURL, nil
}

func (r Repository) GetInfo() ([]models.Shortener, error) {
	cursor, err := r.Client.Database(r.DB).Collection(r.Collection).Find(context.TODO(), bson.D{})

	if err != nil {
		return nil, err
	}

	var shorteners []models.Shortener

	for cursor.Next(context.TODO()) {
		var shortener models.Shortener

		if err = cursor.Decode(&shortener); err != nil {
			return nil, err
		}

		shorteners = append(shorteners, shortener)
	}

	return shorteners, nil
}
