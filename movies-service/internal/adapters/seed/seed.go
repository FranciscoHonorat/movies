package seed

import (
	"context"
	"encoding/json"
	"os"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type MoveiSeed struct {
	ID    int32  `json:"id" bson:"_id"`
	Title string `json:"title" bson:"title"`
	Year  string `json:"year" bson:"year"`
}

func Seed(ctx context.Context, collection *mongo.Collection, filePath string) error {
	count, err := collection.CountDocuments(ctx, bson.D{})
	if err != nil {
		return err
	}
	if count > 0 {
		return nil
	}

	data, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	var movies []MoveiSeed
	err = json.Unmarshal(data, &movies)
	if err != nil {
		return err
	}
	docs := []interface{}{}
	for _, movie := range movies {
		docs = append(docs, movie)
	}

	_, err = collection.InsertMany(ctx, docs)
	if err != nil {
		return err
	}
	return nil
}
