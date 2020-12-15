/**
 * This file handle the monoDB connection.
 *
 * @author Boris
 * @version 2020-12-09
 *
 */

package mongodb

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

//MongoClient is the variable of mongo.client
var MongoClient *mongo.Client

//InitRun function initialize the database
func InitRun() {
	MongoClient = InitMongoDB()
}

//InitMongoDB function start to connect to mongo Database
func InitMongoDB() *mongo.Client {
	fmt.Println("into InitMongoDB")
	//Build connection
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(
		"mongodb+srv://root:root@cluster0.qfx1p.mongodb.net/short-url?retryWrites=true&w=majority",
	))
	if err != nil {
		fmt.Println("connect error!")
		log.Fatal(err)
	}
	//check connection timeout
	// if err = client.Ping(ctx, readpref.Primary()); err != nil {
	// 	panic(err)
	// }
	return client
}
