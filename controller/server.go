package controller

import (
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"regexp"
	mongoDB "short_url/module"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/speps/go-hashids"
	"go.mongodb.org/mongo-driver/bson"
)

var tpl *template.Template
var hostNumber string

type HostNumber struct {
	PortNumber string
}

func Init() {
	tpl = template.Must(template.ParseGlob("view/*.html"))
}

/**
 * This function show the index page
 *
 * @ writer
 * @ request
 *
 */
func Index(writer http.ResponseWriter, request *http.Request) {
	fmt.Println("show home page")
	Init()
	// var Host HostNumber
	// hostNumber = request.Host
	// Host.PortNumber = request.Host
	// tpl.ExecuteTemplate(writer, "app.html", Host)
	tpl.ExecuteTemplate(writer, "app.html", nil)
}

// func Handler(writer http.ResponseWriter, request *http.Request) {
// 	fmt.Println("into Server Handler", count)
// 	//Init()
// 	tpl.ExecuteTemplate(writer, "app.html", nil)
// 	// if request != nil {
// 	// 	fmt.Println("!= nil What is :", request.Body)
// 	// } else {
// 	// 	fmt.Println("== nil What is :", request.Body)
// 	// }
// 	// fmt.Fprintf(writer, "Hello, World")
// 	//template := template.Must(template.ParseFiles("view/test.html"))
// 	// template := template.Must(template.ParseGlob("view/*.html"))
// 	// data := new(pageData)
// 	// data.Title = "Long URL 2 Short URL"
// 	// data.Long_url = "This is long_url"
// 	// data.Short_url = "This is short_url"
// 	//template.Execute(writer, nil)
// }

type RequestData struct {
	OriginalURL string `json:"originalURL,omitempty"`
	Alias       string `json:"alias,omitempty"`
}

type ResponseData struct {
	OriginalURL string `json:"originalURL,omitempty"`
	ShortURL    string `json:"shortURL,omitempty"`
	ID          string `json:"id,omitempty"`
	IsAlias     bool   `json:"isAlias,omitempty"`
}

//CreateEndPoint do
func CreateEndPoint(w http.ResponseWriter, r *http.Request) {
	fmt.Println("into createEndpoint function")
	var request RequestData
	json.NewDecoder(r.Body).Decode(&request)
	// url.ID = ProduceUniqueID()
	// url.ShortURL = "http://localhost:8000/" + url.ID

	collection := mongoDB.MongoClient.Database("url_database").Collection("url_table")
	var response ResponseData
	//check the longURL exist or not
	// var findOne bson.M
	collection.FindOne(context.Background(), bson.M{"OriginalURL": request.OriginalURL}).Decode(&response)
	// fmt.Println(findOne["ShortURL"])
	if (response == ResponseData{}) { //if the long URL does not exist, create new one
		response.ID = ProduceUniqueID()
		response.ShortURL = "http://localhost:8000/" + response.ID
		response.OriginalURL = request.OriginalURL
		collection.InsertOne(context.TODO(), bson.M{
			"ID":          response.ID,
			"OriginalURL": response.OriginalURL,
			"ShortURL":    response.ShortURL})
		fmt.Println("insert successful")
	} else {
		fmt.Println("The Original URL exist")
	}
	// err =
	//  if err!=nil{
	// 	w.WriteHeader(401)
	// 	w.Write([]byte(err.Error()))
	// 	return
	// }
	json.NewEncoder(w).Encode(response)
	fmt.Println("Send respond successful")

}

//CreateURL function create a new short URL
func CreateURL(w http.ResponseWriter, r *http.Request) {
	fmt.Println("into createURL function")
	var request RequestData
	str := r.FormValue("originalURL")
	fmt.Println("The Host Number is" + hostNumber)
	request.OriginalURL = strings.TrimSpace(str)
	request.Alias = r.FormValue("alias")
	// collection := mongoDB.MongoClient.Database("url_database").Collection("url_table")
	fmt.Println("User Input URL:", request.OriginalURL)
	fmt.Println("User Input Alias:", request.Alias)

	// var response ResponseData
	// collection.FindOne(context.Background(), bson.M{"OriginalURL": request.OriginalURL}).Decode(&response)

	if request.Alias == "" { //No Custom Alias input
		CreateWithoutAlias(w, request)
	} else { // With Custom Alias
		CreateWithAlias(w, request)
	}
	fmt.Println("Send respond successful")
}

func PrefixHead(url string) {

}

//CreateWithAlias handle the condition that user input custom alias
func CreateWithAlias(w http.ResponseWriter, request RequestData) {
	/**
	 *
	 *
	 */
	var response ResponseData
	collection := mongoDB.MongoClient.Database("url_database").Collection("url_table")
	collection.FindOne(context.Background(), bson.M{"ID": request.Alias}).Decode(&response)

	if (response == ResponseData{}) { //The custom alias does not exist
		response.ID = request.Alias
		response.ShortURL = "http://localhost:8000/" + response.ID
		response.OriginalURL = request.OriginalURL
		collection.InsertOne(context.TODO(), bson.M{
			"ID":          response.ID,
			"OriginalURL": response.OriginalURL,
			"ShortURL":    response.ShortURL,
			"IsAlias":     true})
		tpl.ExecuteTemplate(w, "create.html", response)
		fmt.Println("Insert custom alias to DB successful")
	} else if response.OriginalURL != request.OriginalURL {
		tpl.ExecuteTemplate(w, "notAvailable.html", nil)
	} else {
		tpl.ExecuteTemplate(w, "create.html", response)
		fmt.Println("Show Alias result page")
	}
	fmt.Println("Create short URL with custom alias  successful")

}

//CreateWithoutAlias handle the condition that user dose not entry alias
func CreateWithoutAlias(w http.ResponseWriter, request RequestData) {
	/**
	 *  The request is the result of query with Original URL
	 *	if the response is empty, then create a new short URL
	 *
	 * 	Or the response is not empty, but the Original URL had a custom alias URL, it should
	 *  also create a new short URL
	 */
	var response ResponseData
	collection := mongoDB.MongoClient.Database("url_database").Collection("url_table")
	collection.FindOne(context.Background(), bson.M{"OriginalURL": request.OriginalURL, "IsAlias": false}).Decode(&response)
	if (response == ResponseData{}) {
		response.ID = ProduceUniqueID()
		response.ShortURL = "http://localhost:8000/" + response.ID
		response.OriginalURL = request.OriginalURL
		collection.InsertOne(context.TODO(), bson.M{
			"ID":          response.ID,
			"OriginalURL": response.OriginalURL,
			"ShortURL":    response.ShortURL,
			"IsAlias":     false})
		fmt.Println("Insert successful")

	}
	tpl.ExecuteTemplate(w, "create.html", response)
	fmt.Println("Create short URL without custom alias  successful")
}

//Redirect Function redirect the link to long URL
func Redirect(w http.ResponseWriter, r *http.Request) {
	fmt.Println("into Redirect")
	params := mux.Vars(r)
	fmt.Println("Now, ID is:", params["id"])
	if params["id"] == "" {
		return
	}
	collection := mongoDB.MongoClient.Database("url_database").Collection("url_table")
	//Find the Short URL
	var response ResponseData
	collection.FindOne(context.Background(), bson.M{"ID": params["id"]}).Decode(&response)
	// if err != nil {
	if (response == ResponseData{}) {
		notFound(w, r)
	}
	http.Redirect(w, r, response.OriginalURL, 301)
	fmt.Println("send request successful")
}

func notFound(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	template, err := template.ParseFiles("view/errPage.html")
	if err != nil {
		fmt.Println("Not found page open error:", err)
	}
	template.Execute(w, nil)
}

//ProduceUniqueID return a unique ID
func ProduceUniqueID() string {
	h, _ := hashids.NewWithData(hashids.NewData())
	now := time.Now()
	ID, _ := h.Encode([]int{int(now.Unix())})
	return ID
}

// func HandleURL(writer http.ResponseWriter, request *http.Request) {
// 	fmt.Println("into URL Handler")
// 	if request.Method != "POST" {
// 		http.Redirect(writer, request, "/", http.StatusSeeOther)
// 		return
// 	}

// 	LongURL := request.FormValue("url")
// 	fmt.Println("Long ", LongURL)
// 	Alias := request.FormValue("alias")
// 	fmt.Println("Alias", Alias)
// 	ShortURL := AliasURL + Alias
// 	fmt.Println("short", ShortURL)

// 	result := struct {
// 		Original string
// 		Shorten  string
// 	}{
// 		Original: LongURL,
// 		Shorten:  ShortURL,
// 	}
// 	fmt.Println(result)

// 	tpl.ExecuteTemplate(writer, "submission.html", result)

// }

func IsValidAlias(s string) bool {
	result, _ := regexp.MatchString("^[a-zA-Z0-9]+$", s)
	return result
}

// func main() {
// 	//set visit url
// 	http.HandleFunc("/home", handler)

// 	/*
// 	 * Start a web server.
// 	 * if fail, log.Fatal show error in output
// 	 */
// 	log.Fatal(http.ListenAndServe(":8080", nil))
// }
