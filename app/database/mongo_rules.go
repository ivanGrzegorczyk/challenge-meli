package database

import (
	"context"
	"os"
	"time"

	"github.com/ivanGrzegorczyk/challenge_meli/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var mdb *mongo.Client

func init() {
	var err error
	mdb, err = newMongoClient()
	if err != nil {
		panic(err)
	}
}

func newMongoClient() (*mongo.Client, error) {
	mongoHost := os.Getenv("MONGO_HOST")
	if mongoHost == "" {
		mongoHost = "mongodb://localhost"
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoHost+":27017"))
	if err != nil {
		return nil, err
	}

	return client, nil
}

func InsertRule(rule model.Rule) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := mdb.Database("challenge_meli").Collection("rules")
	_, err := collection.InsertOne(ctx, rule)
	if err != nil {
		return err
	}

	return nil
}

func GetRulesByIp(ip string) ([]model.Rule, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := mdb.Database("challenge_meli").Collection("rules")
	cursor, err := collection.Find(ctx, bson.M{"ip": ip})
	if err != nil {
		return nil, err
	}

	var rules []model.Rule
	err = cursor.All(ctx, &rules)
	if err != nil {
		return nil, err
	}

	return rules, nil
}

func GetRulesByPath(path string) ([]model.Rule, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := mdb.Database("challenge_meli").Collection("rules")
	cursor, err := collection.Find(ctx, bson.M{"path": path})
	if err != nil {
		return nil, err
	}

	var rules []model.Rule
	err = cursor.All(ctx, &rules)
	if err != nil {
		return nil, err
	}

	return rules, nil
}
