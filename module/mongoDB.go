package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

/**
 *	This constructor contains short URL, Original URL and Custom as a string, string and boolean.
 *
 *
 */
type MatchingURL struct {
	Short    string `bson: "short"`
	Original string `bson: "original"`
	Custom   bool   `bson: "custom"`
}

func Insert(collection *mongo.Collection, doc *MatchingURL) {
	// collection := c.Database("url_database").Collection("url_table")
	_, err := collection.InsertOne(context.TODO(), &doc)
	if err != nil {
		log.Fatal(err)
	}
	// fmt.Println(insertResult)
}

func SearchURL(collection *mongo.Collection, s string) {
	// err := collection.FindOne(context.TODO(), fliter).Decode(&amp ; s)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// return false
	// fmt.Println(result)
}

func main() {
	fmt.Println("Starting the MongoDB")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(
		"mongodb+srv://root:root@cluster0.qfx1p.mongodb.net/short-url?retryWrites=true&w=majority",
	))
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(ctx)

	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		log.Fatal(err)
	}
	database, err := client.ListDatabaseNames(ctx, bson.M{})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(database)
	data := MatchingURL{
		Short:    "google.com",
		Original: "short.com",
		Custom:   false,
	}
	// urlDatabase := client.Database("url_database")
	// urlCollection := urlDatabase.Collection("url_table")
	collection := client.Database("url_database").Collection("url_table")
	urlTableResult, err := collection.InsertOne(ctx, &data)
	if err != nil {
		log.Fatal(err)
	}
	// Insert(collection, &data)
	fmt.Println(urlTableResult.InsertedID)

	// if err != nil {
	// 	log.Fatal(err)
	// }
	// return false
	// fmt.Println(result)

	// err = collection.Aggregate(ctx, {$filter : { short : { $eq : "google.com" }})
	// if err != nil {
	// 	log.Fatal(err)
	// }
	fmt.Println("end")

	// 	{}
	// })

}
