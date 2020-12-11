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

//MongoClient .....
var MongoClient *mongo.Client

//InitRun function
func InitRun() {
	MongoClient = InitMongoDB()
}

//InitMongoDB function start the mongo Database
func InitMongoDB() *mongo.Client {
	fmt.Println("into InitMongoDB")
	//Build connection
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(
		"mongodb+srv://root:root@cluster0.qfx1p.mongodb.net/short-url?retryWrites=true&w=majority",
	))
	if err != nil {
		log.Fatal(err)
	}

	//check connection timeout
	err = client.Ping(ctx, readpref.Primary())
	defer cancel()
	if err != nil {
		log.Fatal(err)
	}
	return client
}

//InsertData function insert URL data into collection
func InsertData(collection *mongo.Collection, doc *MatchingURL) {
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

//RequestData data is from the request
// type RequestData struct {
// 	OriginalURL string `bson:"originalURL,omitempty"`
// 	CustomAlias string `bson:"customAlias,omitempty"`
// }
type RequestData struct {
	OriginalURL string `json:"originalURL,omitempty"`
	CustomAlias string `json:"customAlias,omitempty"`
	ShortURL    string `json:"shortURL,omitempty"`
}

//ResponseData data is from the request
type ResponseData struct {
	OriginalURL string `json:"originalURL,omitempty"`
	ShortURL    string `json:"shortURL,omitempty"`
}

func main() {
	fmt.Println("Starting the MongoDB")
	InitRun()
	collection := MongoClient.Database("url_database").Collection("url_table")
	// reqData := RequestData{
	// 	OriginalURL: "https://google.com",
	// }

	var findOne ResponseData
	collection.FindOne(context.Background(), bson.M{"OriginalURL": "https://www.linkedin.com"}).Decode(&findOne)
	// if findOne == nil {
	// 	collection.InsertOne(context.TODO(), bson.M{"OriginalURL": "https://www.linkedin.com", "ShortURL": "http://localhost:8000/778899"})
	// 	fmt.Println("insert successful")
	// }
	fmt.Println(findOne == ResponseData{})
	if (findOne != ResponseData{}) {
		fmt.Println("ok")
	}

	// if findOne["OriginalURL"] == nil {
	// 	fmt.Println("nil", findOne["OriginalURL"])
	// } else {
	// 	fmt.Println(findOne["OriginalURL"])
	// }

	// req, err := collection.Find(context.TODO(), bson.M{})

	// err := collection.FindOne(context.Background(), bson.M{"OriginalURL": "https://www.youtube.com/"}).Decode(&findOne)
	// err := collection.FindOne(context.Background(), bson.M{"ShortURL": bson.M{"OriginalURL": "https://www.youtube.com"}}).Decode(&findOne)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// fmt.Println(findOne == nil)
	// fmt.Println(findOne)
	// var find RequestData
	// cur, _ := collection.Find(context.Background(), bson.M{"OriginalURL": bson.M{"ShortURL": "https://www.youtube.com"}})
	// cur.Decode(&find)
	// fmt.Println(cur)

	// var find RequestData
	// // collection := MongoClient.Database("url_database").Collection("url_table")
	// // req, err = collection.CountDocuments(context.TODO(), bson.M{"OriginalURL": "https://google.com"})
	// findCursor, err := collection.Find(context.Background(), bson.M{"OriginalURL": "https://www.google.com/"}).Decode(&find)
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// findCursor.Decode(&find)
	// fmt.Println(find)

	// fmt.Println(find.OriginalURL)

	// var findLogic RequestData
	// logicFilter := bson.M{
	// 	"$and": bson.A{
	// 		bson.M{"custom": bson.M{"$gt": false}},
	// 		bson.M{"originalURL": bson.M{"$gt": "https://google.com/"}},
	// 	},
	// }
	// findLogicRes, err := collection.Find(context.Background(), logicFilter)
	// err = findLogicRes.Decode(&findLogic)
	// if err != nil {
	// 	fmt.Println("This ", err)
	// }
	// fmt.Println(findLogic)

}

// func main() {
// 	fmt.Println("Starting the MongoDB")
// 	InitRun()
// 	database, err := MongoClient.ListDatabaseNames(context.TODO(), bson.M{})
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	fmt.Println(database)
// 	collection := MongoClient.Database("url_database").Collection("url_table")
// 	fmt.Println("This collection", collection)

// 	// if err != nil {
// 	// 	log.Fatal(err)
// 	// }
// 	// fmt.Println(database)
// 	// ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
// 	// defer cancel()
// 	// client, err := mongo.Connect(ctx, options.Client().ApplyURI(
// 	// 	"mongodb+srv://root:root@cluster0.qfx1p.mongodb.net/short-url?retryWrites=true&w=majority",
// 	// ))
// 	// if err != nil {
// 	// 	log.Fatal(err)
// 	// }
// 	// defer client.Disconnect(ctx)

// 	// err = client.Ping(ctx, readpref.Primary())
// 	// if err != nil {
// 	// 	log.Fatal(err)
// 	// }
// 	// database, err := client.ListDatabaseNames(ctx, bson.M{})
// 	// if err != nil {
// 	// 	log.Fatal(err)
// 	// }
// 	// fmt.Println(database)
// 	// data := MatchingURL{
// 	// 	Short:    "google.com",
// 	// 	Original: "short.com",
// 	// 	Custom:   false,
// 	// }
// 	// urlDatabase := client.Database("url_database")
// 	// urlCollection := urlDatabase.Collection("url_table")
// 	// collection := client.Database("url_database").Collection("url_table")
// 	// urlTableResult, err := collection.InsertOne(ctx, &data)
// 	// if err != nil {
// 	// 	log.Fatal(err)
// 	// }
// 	// Insert(collection, &data)
// 	// fmt.Println(urlTableResult.InsertedID)

// 	// if err != nil {
// 	// 	log.Fatal(err)
// 	// }
// 	// return false
// 	// fmt.Println(result)

// 	// err = collection.Aggregate(ctx, {$filter : { short : { $eq : "google.com" }})
// 	// if err != nil {
// 	// 	log.Fatal(err)
// 	// }
// 	fmt.Println("end")

// 	// 	{}
// 	// })

// }
