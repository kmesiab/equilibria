package atlas

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/kmesiab/equilibria/lambdas/lib/ai"
	"github.com/kmesiab/equilibria/lambdas/models"
)

type MongoRepository struct {
	client     *mongo.Client
	collection *mongo.Collection
	aiService  *ai.OpenAICompletionService // Embedding service
}

// NewMongoRepository initializes a new instance of MongoRepository
func NewMongoRepository(client *mongo.Client, dbName string, collName string, aiService *ai.OpenAICompletionService) *MongoRepository {
	return &MongoRepository{
		client:     client,
		collection: client.Database(dbName).Collection(collName),
		aiService:  aiService,
	}
}

// CreateOrUpdateDocument handles the insertion or update of documents including vector embeddings
func (r *MongoRepository) CreateOrUpdateDocument(ctx context.Context, message *models.Message) error {
	// Generate embeddings from the message body
	vectors, err := r.aiService.GetEmbeddings(message.Body)
	if err != nil {
		log.Printf("Error generating embeddings: %v", err)
		return err
	}

	// Define the MongoDB document
	doc := AtlasDocument{
		UserID:    message.FromUserID,
		MessageID: message.ID,
		CreatedAt: message.CreatedAt,
		Vectors:   vectors,
	}

	// Upsert option to insert or update the document
	opts := options.Update().SetUpsert(true)
	filter := bson.M{"message_id": message.ID}
	update := bson.M{"$set": doc}

	_, err = r.collection.UpdateOne(ctx, filter, update, opts)
	if err != nil {
		return fmt.Errorf("failed to upsert document: %v", err)
	}
	return nil
}

// FindByID retrieves a document by its MessageID from the MongoDB collection.
func (r *MongoRepository) FindByID(ctx context.Context, messageID int64) (*AtlasDocument, error) {
	var document AtlasDocument
	filter := bson.M{"message_id": messageID}
	err := r.collection.FindOne(ctx, filter).Decode(&document)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve document: %v", err)
	}
	return &document, nil
}

// Delete removes a document by its MessageID from the MongoDB collection.
func (r *MongoRepository) Delete(ctx context.Context, messageID int64) error {
	filter := bson.M{"message_id": messageID}
	_, err := r.collection.DeleteOne(ctx, filter)
	if err != nil {
		return fmt.Errorf("failed to delete document: %v", err)
	}
	return nil
}

// VectorSearch performs a basic vector search within a given tolerance.
// Note: This is a placeholder and may require a more sophisticated approach or use of a vector database.
func (r *MongoRepository) VectorSearch(ctx context.Context, vector []float32, tolerance float32) ([]AtlasDocument, error) {
	var documents []AtlasDocument
	// MongoDB does not natively support vector similarity search, so this would be a placeholder.
	// You may need to use a specialized service or database for efficient vector search.
	log.Println("Vector search is not implemented, as MongoDB does not support this natively.")
	return documents, fmt.Errorf("vector search not implemented")
}
