package atlas_test

/*
import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/kmesiab/equilibria/lambdas/lib/ai"
	"github.com/kmesiab/equilibria/lambdas/lib/atlas"
	"github.com/kmesiab/equilibria/lambdas/lib/config"
	"github.com/kmesiab/equilibria/lambdas/lib/log"
	"github.com/kmesiab/equilibria/lambdas/models"
)

func setupTestDB(t *testing.T) (*mongo.Client, *mongo.Collection) {
	ctx := context.Background()
	cfg := config.Get()
	clientOptions := atlas.GetClientOptions(cfg)
	client, err := atlas.GetMongoClient(ctx, clientOptions)
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}
	db := client.Database(cfg.AtlasDBName)
	coll := db.Collection("testcoll")
	return client, coll
}

func TestCreateOrUpdateDocument(t *testing.T) {
	client, coll := setupTestDB(t)
	defer func(client *mongo.Client, ctx context.Context) {
		err := client.Disconnect(ctx)
		if err != nil {
			log.New("Unrecoverable fatal error closing MongoDB client: %v", err).
				Log()
		}
	}(client, context.Background())

	aiService := &ai.OpenAICompletionService{RemoveEmojis: false}

	repo := atlas.NewMongoRepository(client, "testdb", "test coll", aiService)

	message := &models.Message{
		ID:         1,
		FromUserID: 1,
		Body:       "This is a test message",
		CreatedAt:  time.Now(),
	}

	err := repo.CreateOrUpdateDocument(context.Background(), message)
	assert.NoError(t, err)

	var result atlas.AtlasDocument
	err = coll.FindOne(context.Background(), bson.M{"message_id": message.ID}).Decode(&result)
	assert.NoError(t, err)
	assert.Equal(t, message.ID, result.MessageID)
}

func TestFindByID(t *testing.T) {
	client, coll := setupTestDB(t)
	defer func(client *mongo.Client, ctx context.Context) {
		err := client.Disconnect(ctx)
		if err != nil {
			log.New("Unrecoverable fatal error closing MongoDB client: %v", err).
				Log()
		}
	}(client, context.Background())

	aiService := &ai.OpenAICompletionService{RemoveEmojis: false}
	repo := atlas.NewMongoRepository(client, "testdb", "test coll", aiService)

	expectedDoc := atlas.AtlasDocument{
		UserID:    1,
		MessageID: 1,
		CreatedAt: time.Now(),
		Vectors:   []float32{0.1, 0.2, 0.3},
	}
	_, err := coll.InsertOne(context.Background(), expectedDoc)
	assert.NoError(t, err)

	result, err := repo.FindByID(context.Background(), 1)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, expectedDoc.MessageID, result.MessageID)
}

func TestDelete(t *testing.T) {
	client, coll := setupTestDB(t)
	defer func(client *mongo.Client, ctx context.Context) {
		err := client.Disconnect(ctx)
		if err != nil {
			log.New("Unrecoverable fatal error closing MongoDB client: %v", err).
				Log()
		}
	}(client, context.Background())

	aiService := &ai.OpenAICompletionService{RemoveEmojis: false}
	repo := atlas.NewMongoRepository(client, "testdb", "testcoll", aiService)

	doc := atlas.AtlasDocument{
		UserID:    1,
		MessageID: 1,
		CreatedAt: time.Now(),
		Vectors:   []float32{0.1, 0.2, 0.3},
	}
	_, err := coll.InsertOne(context.Background(), doc)
	assert.NoError(t, err)

	err = repo.Delete(context.Background(), 1)
	assert.NoError(t, err)

	err = coll.FindOne(context.Background(), bson.M{"message_id": 1}).Err()
	assert.Error(t, err)
	assert.Equal(t, mongo.ErrNoDocuments, err)
}

*/
