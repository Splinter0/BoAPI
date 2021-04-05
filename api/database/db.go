package database

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/splinter0/api/models"
)

func DBInstance() *mongo.Client {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}
	URI := os.Getenv("MONGO_URI")
	client, err := mongo.NewClient(options.Client().ApplyURI(URI))
	if err != nil {
		log.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	defer cancel()
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connected to MongoDB!")

	return client
}

var Client *mongo.Client = DBInstance()

func OpenCollection(client *mongo.Client, collectionName string) *mongo.Collection {
	return client.Database("test").Collection(collectionName)
}

func AddUser(user models.User) bool {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	_, err := OpenCollection(Client, "users").InsertOne(ctx, user)
	defer cancel()
	if err != nil {
		fmt.Println(err)
		return false
	}
	return true
}

func FindUser(username string) models.User {
	var u models.User
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	tasks := OpenCollection(Client, "users")
	tasks.FindOne(ctx, bson.M{"username": username}).Decode(&u)
	defer cancel()
	return u
}

func AddUserToken(username, token string) {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	clients := OpenCollection(Client, "users")
	clients.UpdateOne(
		ctx,
		bson.M{"username": username},
		bson.D{
			{"$set", bson.D{{"last", time.Now()}, {"token", token}}},
		},
	)
	defer cancel()
}

func GetUserToken(username string) string {
	user := FindUser(username)
	return user.Token
}
