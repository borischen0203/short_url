package controller

import (
	"context"
	"fmt"
	"html/template"
	"net/http"
	"path"
	mongoDB "short_url/module"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/speps/go-hashids"
	"go.mongodb.org/mongo-driver/bson"
)

var tpl *template.Template
var hostNumber string

//RequestData ...
type RequestData struct {
	OriginalURL string `json:"originalURL,omitempty"`
	Alias       string `json:"alias,omitempty"`
}

//ResponseData  ...
type ResponseData struct {
	OriginalURL string `json:"originalURL,omitempty"`
	ShortURL    string `json:"shortURL,omitempty"`
	ID          string `json:"id,omitempty"`
	IsAlias     bool   `json:"isAlias,omitempty"`
}

// type HostNumber struct {
// 	PortNumber string
// }

//Init function ...
func Init() {
	tpl = template.Must(template.ParseGlob("view/*.html"))
}

//Index function ...
func Index(writer http.ResponseWriter, request *http.Request) {
	fmt.Println("show home page")
	Init()
	// var Host HostNumber
	// hostNumber = request.Host
	// Host.PortNumber = request.Host
	// tpl.ExecuteTemplate(writer, "app.html", Host)
	err := tpl.ExecuteTemplate(writer, "index.html", nil)
	if err != nil {
		writer.WriteHeader(400)
		writer.Write([]byte("error load index.html"))
	}
	//  else {

	// 	writer.WriteHeader(200)
	// 	// var response ResponseData
	// 	// response.ID = "1234567"
	// 	// response.OriginalURL = "google.com"
	// 	// response.ShortURL = "google.com/1234567"
	// 	// response.IsAlias = false
	// 	// json.NewEncoder(writer).Encode(response)

	// }
}

//PrefixSlash function ...
func PrefixSlash(inputURL string) string {
	url := strings.TrimSpace(inputURL)
	// hasSlash := strings.HasSuffix(url, "/")
	url = path.Clean(url)
	if !strings.HasSuffix(url, "/") {
		url += "/"
	}
	//  else if strings.HasSuffix(url, ".com") {
	// 	url += "/"
	// }
	url = strings.Replace(url, "https:/", "https://", 1)
	url = strings.Replace(url, "http:/", "http://", 1)
	return url
}

//CreateURL function create a new short URL
func CreateURL(w http.ResponseWriter, r *http.Request) {
	fmt.Println("into createURL function")

	var request RequestData
	str := r.PostFormValue("originalURL")
	// fmt.Println("The Host Number is" + hostNumber)
	request.OriginalURL = PrefixSlash(str)
	request.Alias = r.PostFormValue("alias")
	fmt.Println("User Input URL:", request.OriginalURL)
	fmt.Println("User Input Alias:", request.Alias)

	if request.Alias == "" { //No Custom Alias input
		fmt.Println("Into No alias process")
		requestData := CreateWithoutAlias(request)
		tpl.ExecuteTemplate(w, "create.html", requestData)
	} else { // With Custom Alias
		fmt.Println("Into alias process")
		requestData := CreateWithAlias(request)
		if (requestData == ResponseData{}) {
			tpl.ExecuteTemplate(w, "notAvailable.html", nil)
		} else {
			tpl.ExecuteTemplate(w, "create.html", requestData)
		}
	}
	fmt.Println("Send respond successful")
}

//CreateWithAlias handle the condition that user input custom alias
func CreateWithAlias(request RequestData) ResponseData {
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
		// tpl.ExecuteTemplate(w, "create.html", response)
		fmt.Println("Insert custom alias to DB successful")
	} else if response.OriginalURL != request.OriginalURL {
		// tpl.ExecuteTemplate(w, "notAvailable.html", nil)
		response = ResponseData{}
	}
	// else {
	// tpl.ExecuteTemplate(w, "create.html", response)
	// fmt.Println("Show Alias result page")
	// }
	fmt.Println("Create short URL with custom alias  successful")
	return response
}

//CreateWithoutAlias handle the condition that user dose not entry alias
func CreateWithoutAlias(request RequestData) ResponseData {
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
	// tpl.ExecuteTemplate(w, "create.html", response)
	fmt.Println("Create short URL without custom alias  successful")
	return response
}

//Redirect Function redirect the link to long URL
func Redirect(w http.ResponseWriter, r *http.Request) {
	fmt.Println("into Redirect")
	params := mux.Vars(r)
	collection := mongoDB.MongoClient.Database("url_database").Collection("url_table")
	//Find the Short URL
	var response ResponseData
	collection.FindOne(context.Background(), bson.M{"ID": params["id"]}).Decode(&response)
	if (response != ResponseData{}) {
		http.Redirect(w, r, response.OriginalURL, 301)
		fmt.Println("Redirect successful")
	} else {
		notFound(w, r)
	}
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

// func IsValidAlias(s string) bool {
// 	result, _ := regexp.MatchString("^[a-zA-Z0-9]+$", s)
// 	return result
// }
