/**
 * This main file excute the web server
 *
 * @author Boris
 * @version 2020-12-06
 *
 */

package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type URLdata struct {
	Long  string `json:"long"`
	Short string `json:"short"`
	Alias string `json:alias`
}

var urldatas []URLdata

func getURL(w http.ResponseWriter, r *http.Request) {
	fmt.Println("into getURL function")
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(urldatas)
	param := mux.Vars(r) //Get params
	fmt.Print(param["long"])
}

func createURL(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, World")
}

type MyUrl struct {
	ID       string `json:"id,omitempty`
	LongUrl  string `json:"longUrl,omitempty"`
	ShortUrl string `json:"shortUrl,omitempty"`
}

// var bucket *gocb.Bucket
// var bucketName string

func ExpandEndPint(w http.ResponseWriter, r *http.Request) {

}

func CreateEndPint(w http.ResponseWriter, r *http.Request) {

}

func RootEndPint(w http.ResponseWriter, r *http.Request) {

}

func main() {
	fmt.Println("Server start")
	router := mux.NewRouter()

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
	// http.HandleFunc("/", server.Handler)
	// http.HandleFunc("/", controller.Index)
	// r.HandleFunc("/submission",).Methods("POST")
	// router.HandleFunc("/{name}", controller.JumpURL)
	// http.HandleFunc("/jump", controller.JumpURL)

	//mock data
	urldatas = append(urldatas, URLdata{
		Long:  " https://google.com",
		Alias: "borisgood",
	})

	urldatas = append(urldatas, URLdata{
		Long:  " https://youtube.com",
		Alias: "borisbad",
	})

	router.HandleFunc("/long2short", getURL).Methods("GET")

	router.HandleFunc("/create", CreateEndPint).Methods("PUT")
	router.HandleFunc("/expand", ExpandEndPint).Methods("GET")
	router.HandleFunc("/{id}", RootEndPint).Methods("GET") //into shortURL ->redirection

	// http.HandleFunc("/submission", controller.HandleURL)

	// err := http.ListenAndServe(":8000", router)
	// if err != nil {
	// 	log.Fatal("ListenAdnServe", err)
	// } else {
	// 	log.Println("Listen 8000")
	// }
	log.Fatal(http.ListenAndServe(":8000", router))
}
