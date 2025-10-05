package collections

import (
	"github.com/fabianpoels/fabianpoels-api-go/config"
	"go.mongodb.org/mongo-driver/mongo"
)

func GetUserCollection(client mongo.Client) *mongo.Collection {
	return getCollection("users", client)
}

func getCollection(name string, client mongo.Client) *mongo.Collection {
	return client.Database(config.GetConfig().GetString("database")).Collection(name)
}
