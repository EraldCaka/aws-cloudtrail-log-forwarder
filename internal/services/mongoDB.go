package services

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoService struct {
	client   *mongo.Client
	database *mongo.Database
}

type Source struct {
	ID               string `json:"id"`
	SourceType       string `json:"sourceType"`
	Credentials      Credentials
	LogFetchInterval int    `json:"logFetchInterval"`
	CallbackURL      string `json:"callbackUrl"`
	S3Bucket         string `json:"s3Bucket"`
	S3Prefix         string `json:"s3Prefix"`
}

type Credentials struct {
	AccessKeyId     string `json:"accessKeyId"`
	SecretAccessKey string `json:"secretAccessKey"`
	Region          string `json:"region"`
}

func NewMongoService(uri, dbName string) (*MongoService, error) {
	clientOptions := options.Client().ApplyURI(uri)

	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MongoDB: %v", err)
	}

	err = client.Ping(context.Background(), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to ping MongoDB: %v", err)
	}

	db := client.Database(dbName)
	return &MongoService{
		client:   client,
		database: db,
	}, nil
}

func (m *MongoService) InsertSource(source Source) error {
	collection := m.database.Collection("sources")
	_, err := collection.InsertOne(context.Background(), source)
	if err != nil {
		return fmt.Errorf("failed to insert source into MongoDB: %v", err)
	}
	return nil
}

func (m *MongoService) GetSources() ([]Source, error) {
	var sources []Source
	collection := m.database.Collection("sources")
	cursor, err := collection.Find(context.Background(), bson.D{})
	if err != nil {
		return nil, fmt.Errorf("failed to fetch sources: %v", err)
	}

	for cursor.Next(context.Background()) {
		var source Source
		if err := cursor.Decode(&source); err != nil {
			return nil, fmt.Errorf("failed to decode source: %v", err)
		}
		sources = append(sources, source)
	}

	return sources, nil
}

func (m *MongoService) RemoveSource(id string) error {
	collection := m.database.Collection("sources")
	_, err := collection.DeleteOne(context.Background(), bson.M{"id": id})
	if err != nil {
		return fmt.Errorf("failed to remove source: %v", err)
	}
	return nil
}

func (m *MongoService) InsertDocument(collectionName string, document interface{}) (interface{}, error) {
	collection := m.database.Collection(collectionName)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	result, err := collection.InsertOne(ctx, document)
	if err != nil {
		return nil, fmt.Errorf("failed to insert document: %v", err)
	}
	return result.InsertedID, nil
}
