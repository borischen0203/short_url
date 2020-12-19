/**
 * This main file excute the web server.
 *
 * @author Boris
 * @version 2020-12-06
 *
 */

package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	controller "short_url/controller"
	mongoDB "short_url/module"

	"github.com/go-redis/redis/v8"
	"github.com/gorilla/mux"
)

// var client *redis.Client

var ctx = context.Background()

func ExampleClient() {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	err := rdb.Set(ctx, "Hello", "world", 0).Err()
	if err != nil {
		panic(err)
	}

	val, err := rdb.Get(ctx, "Hello").Result()
	if err != nil {
		panic(err)
	}
	fmt.Println("Hello", val)

	val2, err := rdb.Get(ctx, "go").Result()
	if err == redis.Nil {
		fmt.Println("key2 does not exist")
	} else if err != nil {
		panic(err)
	} else {
		fmt.Println("go", val2)
	}
	// Output: key value
	// key2 does not exist
}

func main() {
	fmt.Println("Server start")
	ExampleClient()
	router := mux.NewRouter()

	//Init Database
	mongoDB.InitRun()

	//Show HomePage
	router.HandleFunc("/", controller.Index).Methods("GET")

	//Creat Short URL
	router.HandleFunc("/POST/url_resource", controller.CreateURL).Methods("POST")

	//Redirect Original URL
	router.HandleFunc("/{id}", controller.Redirect).Methods("GET")

	//Server listen at port 8000
	log.Fatal(http.ListenAndServe(":8000", router))
}
