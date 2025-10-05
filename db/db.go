package db

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/fabianpoels/fabianpoels-api-go/collections"
	"github.com/fabianpoels/fabianpoels-api-go/config"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var mongoClient *mongo.Client

func GetDbClient() *mongo.Client {
	if mongoClient == nil {
		DbConnect()
	}
	return mongoClient
}

func DbConnect() {
	// mongoDB config
	username := os.Getenv("MONGODB_USER")
	password := os.Getenv("MONGODB_PASSW")
	host := os.Getenv("MONGODB_HOST")
	port := os.Getenv("MONGODB_PORT")
	mongoUrl := fmt.Sprintf("mongodb://%s:%s", host, port)
	if username != "" && password != "" {
		mongoUrl = fmt.Sprintf("mongodb://%s:%s@%s:%s", username, password, host, port)
	}
	clientOptions := options.Client().ApplyURI(mongoUrl)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	log.Printf("Connecting to db with uri: %s", mongoUrl)

	// Init connection
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal("⛒ Connection Failed to Database")
		log.Fatal(err)
	}

	// Check the connection
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal("⛒ Connection Failed to Database")
		log.Fatal(err)
	}

	log.Println("Connected to database: " + config.GetConfig().GetString("database"))

	mongoClient = client
}

func CreateIndexes() {
	client := GetDbClient()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// (re) create indexes
	userIndexes(ctx, client)
}

func userIndexes(ctx context.Context, client *mongo.Client) {
	userEmailIndex := mongo.IndexModel{
		Keys:    bson.D{{Key: "email", Value: 1}},
		Options: options.Index().SetUnique(true),
	}
	name, err := collections.GetUserCollection(*client).Indexes().CreateOne(ctx, userEmailIndex)
	if err != nil {
		log.Fatal("⛒ Error creating User email index")
		log.Fatal(err)
	}
	log.Println("Created User index: " + name)
}
