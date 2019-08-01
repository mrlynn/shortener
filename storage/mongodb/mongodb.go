// Package mongodb is an example MongoDB client implementation written in Go
// This package extends a Repository metaphor which wraps access methods for
// the MongoDB database and collections that are used in the example.
package mongodb

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/mrlynn/shortener/encoder"
	"github.com/mrlynn/shortener/models"
)

// Repository is a struct used to store details of your MongoDB Database,
// Collection and a pointer to a MongoDB Client object
type Repository struct {
	DB         string
	Collection string
	Client     *mongo.Client
}

// NewMongoRepository takes a Database, a Collection name,
// a pointer to a MongoDB Client object and returns a pointer to a Repository
func NewMongoRepository(db, collection string, client *mongo.Client) *Repository {
	return &Repository{DB: db, Collection: collection, Client: client}
}

// NewMongoClient takes a single parameter uri and returns either an nil, and error,
// or a newly created MongoDB Client object
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

// SaveUrl is a function that has one receiver of type Repository,
// takes a single url string parameter and generates a shortened url,
// returning that after it is saved in the database.
// A call to SaveUrl will result in a document to be created which contains
// the details of a new generated url.
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

// GetURL is a function that has one receiver of type Repository, accepts
// a single parameter of type string that is an encoded URL and returns
// details from the encoded URL. A call to this function results in an
// update of the document stored in the database such that the visit count
// is incremented by one.
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
