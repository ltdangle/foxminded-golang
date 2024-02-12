package bot

import (
	"context"
	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// Test storage flow.
func Test_mongoDbStore(t *testing.T) {
	// go test -short ./...
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}

	store, err := NewMongoDbStore("mongodb://localhost:27017", "fox")
	if err != nil {
		t.Fatalf(err.Error())
	}

	// Clear collection.
	client := store.GetClient()
	collection := client.Database("fox").Collection("subscription")
	err = collection.Drop(context.Background())
	if err != nil {
		t.Fatalf(err.Error())
	}

	saveSubscription(t, store, "10:10:10")
	saveSubscription(t, store, "10:10:10")

	subs, err := store.FindSubsBySchedule("10:10:10")
	if err != nil {
		t.Fatalf(err.Error())
	}

	assert.Equal(t, 2, len(subs))
}

func saveSubscription(t *testing.T, store *mongoDbStore, schedule string) {
	// Create new subscription via geolocation.
	chatId := rand.Int()

	sub := NewSubscription()
	sub.ChatId = chatId
	sub.FirsName = "Test user"
	sub.Latitude = rand.Float64()
	sub.Longitude = rand.Float64()
	sub.UpdatedAt = time.Now()

	err := store.SaveLocation(sub)
	if err != nil {
		t.Fatalf(err.Error())
	}

	// Update schedule.
	err = store.SaveSchedule(sub.Id, schedule)
	if err != nil {
		t.Fatalf(err.Error())
	}

	// Retrieve subscription by chat id.
	rSub := store.GetSubscriptionByChatId(chatId)
	if err != nil {
		t.Fatalf(err.Error())
	}

	// Saves subscription is returned.
	assert.Equal(t, sub.Id, rSub.Id)
	// Schedule matches.
	assert.Equal(t, schedule, rSub.Schedule)
}
