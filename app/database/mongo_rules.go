package database

import (
	"context"
	"errors"
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
	if rule.Ip == "" && rule.Path == "" {
		return errors.New("ip and path can't be empty")
	} else if rule.Time <= 0 {
		return errors.New("time must be greater than 0")
	} else if rule.MaxRequests < 0 {
		return errors.New("max_requests must be greater or equal to 0")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := mdb.Database("challenge_meli").Collection("rules")

	var existingRule model.Rule
	err := collection.FindOne(ctx, bson.M{"ip": rule.Ip, "path": rule.Path}).Decode(&existingRule)
	if err != mongo.ErrNoDocuments {
		if err != nil {
			return err
		}
		return errors.New("a rule with the same ip and path already exists")
	}

	_, err = collection.InsertOne(ctx, rule)
	if err != nil {
		return err
	}

	return nil
}

func GetRulesByIpOrPath(ip, path string) ([]model.Rule, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := mdb.Database("challenge_meli").Collection("rules")
	cursor, err := collection.Find(ctx, bson.M{"$or": []bson.M{{"ip": ip}, {"path": path}}})
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
