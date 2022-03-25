package standupnotesservice

import (
	"context"
	"encoding/json"
	"log"
	"os"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type StandupNote struct {
	Id          int64 `json:"_id"`
	Date        string
	Yesterday   string
	Today       string
	Impediments string
	GoBacks     string
}

func GetDBCollection() *mongo.Collection {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	uri := os.Getenv("MONGODB_URI")
	if uri == "" {
		log.Fatal("You must set your 'MONGODB_URI' environmental variable. See\n\t https://docs.mongodb.com/drivers/go/current/usage-examples/#environment-variable")
	}
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()

	return client.Database("standup-notes").Collection("standup-notes")
}

func GetAllNotes() []StandupNote {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	uri := os.Getenv("MONGODB_URI")
	if uri == "" {
		log.Fatal("You must set your 'MONGODB_URI' environmental variable. See\n\t https://docs.mongodb.com/drivers/go/current/usage-examples/#environment-variable")
	}
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()

	coll := client.Database("standup-notes").Collection("standup-notes")
	var notes []StandupNote

	cursor, err := coll.Find(context.TODO(), bson.D{})
	if err != nil {
		log.Fatal(err)
	}

	if err = cursor.All(context.TODO(), &notes); err != nil {
		log.Fatal(err)
	}

	return notes
}

func GetNote(date string) (StandupNote, error) {
	var note StandupNote
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	uri := os.Getenv("MONGODB_URI")
	if uri == "" {
		log.Fatal("You must set your 'MONGODB_URI' environmental variable. See\n\t https://docs.mongodb.com/drivers/go/current/usage-examples/#environment-variable")
	}
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	if err != nil {
		return note, err
	}
	defer func() {
		if err := client.Disconnect(context.TODO()); err != nil {
			return
		}
	}()
	coll := client.Database("standup-notes").Collection("standup-notes")
	var result bson.M
	filter := bson.M{"Date": date}
	err = coll.FindOne(context.TODO(), filter).Decode(&result)
	if err == mongo.ErrNoDocuments {
		return note, err
	}
	if err != nil {
		return note, err
	}

	jsonData, err := json.MarshalIndent(result, "", "  ")

	if err != nil {
		return note, err
	}

	json.Unmarshal(jsonData, &note)
	return note, nil
}
