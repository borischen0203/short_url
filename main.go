/**
 * This main file excute the web server
 *
 * @author Boris
 * @version 2020-12-06
 *
 */

package main

import (
	"encoding/json"
	"fmt"
	"log"
	d "module/mongodb"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/speps/go-hashids"
	"go.mongodb.org/mongo-driver/mongo"
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

//ProduceUniqueID return a unique ID
func ProduceUniqueID() string {
	h, _ := hashids.NewWithData(hashids.NewData())
	now := time.Now()
	fmt.Println(now)
	ID, _ := h.Encode([]int{int(now.Unix())})
	return ID
}

func createURL(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, World")
}

type MyUrl struct {
	ID       string `json:"id,omitempty`
	LongURL  string `json:"longURL,omitempty"`
	ShortURL string `json:"shortURL,omitempty"`
}

// var bucket *gocb.Bucket
// var bucketName string
var bucket *mongo.Collection

func ExpandEndPoint(w http.ResponseWriter, r *http.Request) {

}

//CreateEndPoint do
func CreateEndPoint(w http.ResponseWriter, r *http.Request) {
	fmt.Println("into createEndpoint function")
	var url MyUrl
	json.NewDecoder(r.Body).Decode(&url)
	url.ID = ProduceUniqueID()
	url.ShortURL = "http://localhost:8000/" + url.ID
	json.NewEncoder(w).Encode(url)
	/**
	 * write a query to check dose longUrl() exist
	 *
	 *
	 */
	// err =
	//  if err!=nil{
	// 	w.WriterHeader(401)
	// 	w.Write([]byte(err.Error()))
	// 	return
	// }

}

func RootEndPoint(w http.ResponseWriter, r *http.Request) {
	fmt.Println("into EndPoint")
	// params :=mux.Vars(r)
	/**
	 * write a query to find the longURL
	 *
	 *
	 *
	 */
	AlongURL := "https://google.com"
	http.Redirect(w, r, AlongURL, 301)

}

func main() {
	fmt.Println("Server start")
	router := mux.NewRouter()
	d.InitRun()
	// DB.InitRun()
	// database, err := DB.MongoClient.ListDatabaseNames(context.TODO(), bson.M{})
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// fmt.Println(database)
	// collection := DB.MongoClient.Database("url_database").Collection("url_table")
	// fmt.Println("This collection", collection)

	// fmt.Println("Starting the MongoDB")
	// ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	// defer cancel()
	// client, err := mongo.Connect(ctx, options.Client().ApplyURI(
	// 	"mongodb+srv://root:root@cluster0.qfx1p.mongodb.net/short-url?retryWrites=true&w=majority",
	// ))
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// bucket := client.Database("url_database").Collection("url_table")
	// fmt.Println(bucket)
	// http.HandleFunc("/", server.Handler)
	// http.HandleFunc("/", controller.Index)
	// r.HandleFunc("/submission",).Methods("POST")
	// router.HandleFunc("/{name}", controller.JumpURL)
	// http.HandleFunc("/jump", controller.JumpURL)
	// router.HandleFunc("/long2short", getURL).Methods("GET")
	router.HandleFunc("/create", CreateEndPoint).Methods("POST")
	router.HandleFunc("/expand", ExpandEndPoint).Methods("GET")
	router.HandleFunc("/{id}", RootEndPoint).Methods("GET") //into shortURL ->redirection

	// http.HandleFunc("/submission", controller.HandleURL)

	// err := http.ListenAndServe(":8000", router)
	// if err != nil {
	// 	log.Fatal("ListenAdnServe", err)
	// } else {
	// 	log.Println("Listen 8000")
	// }
	log.Fatal(http.ListenAndServe(":8000", router))
}
