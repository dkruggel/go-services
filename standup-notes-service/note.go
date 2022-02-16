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

func GetNote(date string) StandupNote {
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

	var result bson.M
	filter := bson.M{"Date": date}
	err = coll.FindOne(context.TODO(), filter).Decode(&result)
	if err == mongo.ErrNoDocuments {
		//return fmt.Sprintf("No document was found with the date %s\n", date)
		var note StandupNote
		return note
	}
	if err != nil {
		panic(err)
	}

	jsonData, err := json.MarshalIndent(result, "", "  ")

	if err != nil {
		panic(err)
	}

	var note StandupNote
	json.Unmarshal(jsonData, &note)
	return note

	// return fmt.Sprintf("Date: %s\nYesterday: %s\nToday: %s\nImpediments: %s\nGo Backs: %s\n", note.Date, note.Yesterday, note.Today, note.Impediments, note.GoBacks)
}
