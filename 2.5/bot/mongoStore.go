package bot

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type StoreInterface interface {
	Save(u Update) error
	FindSubsBySchedule(schedule string) ([]*Subscription, error)
	SaveLocation(u *Subscription) error
	SaveSchedule(uuid string, schedule string) error
	GetSubscriptionByChatId(id int) *Subscription
	GetSubscription(uuid string) (*Subscription, error)
}

type mongoDbStore struct {
	dbName                 string
	client                 *mongo.Client
	subscriptionCollection string
	updateCollection       string
}

func NewMongoDbStore(connectionString string, dbName string) (*mongoDbStore, error) {
	ctx := context.Background()
	client, err := mongo.Connect(
		ctx, options.Client().ApplyURI(connectionString),
		options.Client().SetServerSelectionTimeout(time.Second*3),
	)

	if err != nil {
		return nil, err
	}

	// Error if cannot connect.
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		return nil, err
	}

	return &mongoDbStore{
		client:                 client,
		dbName:                 dbName,
		updateCollection:       "update",
		subscriptionCollection: "subscription",
	}, nil
}

func (s *mongoDbStore) Save(u Update) error {
	collection := s.client.Database(s.dbName).Collection(s.updateCollection)
	_, err := collection.InsertOne(context.Background(), u)
	return err
}

func (s *mongoDbStore) FindSubsBySchedule(schedule string) ([]*Subscription, error) {
	collection := s.client.Database(s.dbName).Collection(s.subscriptionCollection)
	filter := bson.M{"schedule": schedule}

	var results []*Subscription

	cur, err := collection.Find(context.Background(), filter)
	if err != nil {
		return nil, err

	} else {
		for cur.Next(context.Background()) {
			var result *Subscription
			err := cur.Decode(&result)
			if err != nil {
				return nil, err
			} else {
				results = append(results, result)
			}
		}
	}

	return results, nil
}

func (s *mongoDbStore) SaveLocation(u *Subscription) error {
	collection := s.client.Database(s.dbName).Collection(s.subscriptionCollection)
	filter := bson.M{"id": u.Id}

	// Specify new values to replace the existing ones
	update := bson.M{
		"$set": bson.M{
			"chat_id":    u.ChatId,
			"latitude":   u.Latitude,
			"longitude":  u.Longitude,
			"updated_at": time.Now(),
		},
	}

	// Find one document matching the filter and update it
	result := collection.FindOneAndUpdate(context.Background(), filter, update)
	if result.Err() == mongo.ErrNoDocuments {
		// subscription not found, create it
		_, err := collection.InsertOne(context.Background(), u)
		if err != nil {
			return err
		}
	}

	return nil
}
func (s *mongoDbStore) SaveSchedule(uuid string, schedule string) error {
	collection := s.client.Database(s.dbName).Collection(s.subscriptionCollection)
	filter := bson.M{"id": uuid}

	// Specify new values to replace the existing ones
	update := bson.M{
		"$set": bson.M{
			"schedule":   schedule,
			"updated_at": time.Now(),
		},
	}

	// Find one document matching the filter and update it
	result := collection.FindOneAndUpdate(context.Background(), filter, update)
	if result.Err() == mongo.ErrNoDocuments {
		return result.Err()
	}

	return nil
}

func (s *mongoDbStore) GetSubscriptionByChatId(id int) *Subscription {
	collection := s.client.Database(s.dbName).Collection(s.subscriptionCollection)
	var result Subscription
	filter := bson.D{{Key: "chatid", Value: id}}
	err := collection.FindOne(context.Background(), filter).Decode(&result)
	if err != nil {
		// return nil and error if there's an issue fetching User
		return nil
	}
	return &result
}

func (s *mongoDbStore) GetSubscription(uuid string) (*Subscription, error) {
	collection := s.client.Database(s.dbName).Collection(s.subscriptionCollection)
	var result Subscription
	filter := bson.D{{Key:"id", Value:uuid}}
	err := collection.FindOne(context.Background(), filter).Decode(&result)
	if err != nil {
		// return nil and error if there's an issue fetching User
		return nil, fmt.Errorf("error retrieving subscription: %v", err)
	}
	return &result, nil
}

func (s *mongoDbStore) GetClient() *mongo.Client {
	return s.client
}
