package database

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func goDotEnvVariable(key string) string {
	// load .env file
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	return os.Getenv(key)
}

// MongoInstance contains the Mongo client and database objects
type MongoInstance struct {
	Client *mongo.Client
	Db     *mongo.Database
}

var Mg MongoInstance

func init() {
	// Database settings
	dbName := goDotEnvVariable("DB_NAME")
	mongoURI := goDotEnvVariable("DB_URI")

	client, err := mongo.NewClient(options.Client().ApplyURI(mongoURI))

	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	err = client.Connect(ctx)
	db := client.Database(dbName)

	if err != nil {
		log.Fatal(err)
	}

	Mg = MongoInstance{
		Client: client,
		Db:     db,
	}
	fmt.Println("Connected to MongoDB!")
}
