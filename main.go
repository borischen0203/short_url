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
	"html/template"
	"log"
	"net/http"
	server "short_url/controller"
	mongoDB "short_url/module"
	"time"

	"github.com/gorilla/mux"
	"github.com/speps/go-hashids"

	"go.mongodb.org/mongo-driver/bson"
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

type ResponseData struct {
	OriginalURL string `json:"originalURL,omitempty"`
	ShortURL    string `json:"shortURL,omitempty"`
	ID          string `json:"id,omitempty"`
}

func RootEndPoint(w http.ResponseWriter, r *http.Request) {
	fmt.Println("into RootEndPoint")
	params := mux.Vars(r)
	fmt.Println(params["id"])

	collection := mongoDB.MongoClient.Database("url_database").Collection("url_table")
	//Find the Short URL
	var response ResponseData
	collection.FindOne(context.Background(), bson.M{"ID": params["id"]}).Decode(&response)
	// if err != nil {
	if (response == ResponseData{}) {
		w.WriteHeader(404)
		// w.Write([]byte("<h1>404</h1>"))
		// w.WriteHeader(http.StatusNotFound)
		// fmt.Fprintf(w, "Page Not found")
		template := template.Must(template.ParseGlob("view/errPage.html"))
		template.Execute(w, nil)
		fmt.Println("Page Not found")
		return
	}
	http.Redirect(w, r, response.OriginalURL, 301)
	fmt.Println("send request successful")
	// http.Redirect(w, r, "https://www.google.com/", 301)
}

func main() {
	fmt.Println("Server start")
	router := mux.NewRouter()
	mongoDB.InitRun()
	// database, err := mongoDB.MongoClient.ListDatabaseNames(context.TODO(), bson.M{})
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// fmt.Println(database)
	// collection := mongoDB.MongoClient.Database("url_database").Collection("url_table")
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
	router.HandleFunc("/", server.Index)
	// r.HandleFunc("/submission",).Methods("POST")
	// router.HandleFunc("/{name}", controller.JumpURL)

	router.HandleFunc("/create", server.CreateURL).Methods("POST")
	router.HandleFunc("/{id}", server.Redirect).Methods("GET") //into shortURL ->redirection

	// router.HandleFunc("/create", server.CreateEndPoint).Methods("POST")
	router.HandleFunc("/expand", ExpandEndPoint).Methods("GET")

	// err := http.ListenAndServe(":8000", router)
	// if err != nil {
	// 	log.Fatal("ListenAdnServe", err)
	// } else {
	// 	log.Println("Listen 8000")
	// }
	log.Fatal(http.ListenAndServe(":8000", router))
}
