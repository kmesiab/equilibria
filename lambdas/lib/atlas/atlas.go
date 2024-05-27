package atlas

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var globalMongoClient *mongo.Client

type AtlasDocument struct {
	UserID    int64     `json:"user_id" bson:"user_id"`
	MessageID int64     `json:"message_id" bson:"message_id"`
	CreatedAt time.Time `json:"date_created" bson:"date_created"`
	Vectors   []float32 `json:"vectors" bson:"vectors"`
}

/*
func MakeConnectionString(config *config.Config) string {
	return fmt.Sprintf(
		"mongodb+srv://%s:%s@%s/?retryWrites=true&w=majority&appName=%s",
		config.AtlasDBUser,
		config.AtlasDBPassword,
		config.AtlasDBHost,
		config.AtlasDBName,
	)
}
*/

func GetMongoClient(ctx context.Context, options *options.ClientOptions) (*mongo.Client, error) {
	var err error

	if globalMongoClient != nil {

		return globalMongoClient, nil
	}

	globalMongoClient, err = mongo.Connect(ctx, options)

	if err != nil {

		return nil, err
	}

	err = globalMongoClient.Ping(ctx, readpref.Primary())

	if err != nil {

		return nil, err
	}

	return globalMongoClient, nil

}

/*
func GetClientOptions(config *config.Config) *options.ClientOptions {

	return options.Client().ApplyURI(MakeConnectionString(config))
}


*/
