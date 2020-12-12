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

type pageData struct {
	Title   string
	Alias   string
	LongURL string
}

// type RequestData struct {
// 	OriginalURL string `json:"originalURL,omitempty"`
// 	CustomAlias string `json:"customAlias,omitempty"`
// 	ShortURL    string `json:"shortURL,omitempty"`
// }

type MyUrl struct {
	ID       string `json:"id,omitempty`
	LongURL  string `json:"longURL,omitempty"`
	ShortURL string `json:"shortURL,omitempty"`
}

var tpl *template.Template

//RequestData data is from the request
// type RequestData struct {
// 	OriginalURL string `json:"originalURL,omitempty"`
// 	// CustomAlias string `json:"customAlias,omitempty"`
// }

//ResponseData data is from the request
// type ResponseData struct {
// 	OriginalURL string `json:"originalURL,omitempty"`
// 	ShortURL    string `json:"shortURL,omitempty"`
// }

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
}
type ResponseData struct {
	OriginalURL string `json:"originalURL,omitempty"`
	ShortURL    string `json:"shortURL,omitempty"`
	ID          string `json:"id,omitempty"`
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
	fmt.Println("into createEndpoint function")
	str := (r.FormValue("originalURL"))
	request := strings.TrimSpace(str)
	// var request RequestData
	collection := mongoDB.MongoClient.Database("url_database").Collection("url_table")
	var response ResponseData
	//check the longURL exist or not
	// var findOne bson.M
	collection.FindOne(context.Background(), bson.M{"OriginalURL": request}).Decode(&response)
	// fmt.Println(findOne["ShortURL"])
	if (response == ResponseData{}) { //if the long URL does not exist, create new one
		response.ID = ProduceUniqueID()
		response.ShortURL = "http://localhost:8000/" + response.ID
		response.OriginalURL = request
		collection.InsertOne(context.TODO(), bson.M{
			"ID":          response.ID,
			"OriginalURL": response.OriginalURL,
			"ShortURL":    response.ShortURL})
		fmt.Println("Insert successful")
	} else {
		fmt.Println("The Original URL exists")
	}
	// err =
	//  if err!=nil{
	// 	w.WriteHeader(401)
	// 	w.Write([]byte(err.Error()))
	// 	return
	// }
	tpl.ExecuteTemplate(w, "submission.html", response)
	fmt.Println("Send respond successful")

}

//RootEndPoint Function redirect the link to long URL
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
		// w.Write([]byte(err.Error()))
		return
	}
	http.Redirect(w, r, response.OriginalURL, 301)
	fmt.Println("send request successful")
}

//ProduceUniqueID return a unique ID
func ProduceUniqueID() string {
	h, _ := hashids.NewWithData(hashids.NewData())
	now := time.Now()
	fmt.Println(now)
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
