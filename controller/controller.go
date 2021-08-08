/**
 * This controller file handle input URL, and generate short URL
 *
 * @author Boris
 * @version 2020-12-09
 *
 */

package controller

import (
	"context"
	"fmt"
	"html/template"
	"net/http"
	"os"
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
var Host string
var collection string

// NotAvailable ....
type NotAvailable struct {
	Title string
}

//RequestData constructor creates two parts: OriginalURL and Alias. which are a string and a string, respectively.
type RequestData struct {
	OriginalURL string `json:"originalURL,omitempty"`
	Alias       string `json:"alias,omitempty"`
}

//ResponseData constructor creates four parts: OriginalURL, ShortURL, ID and IsAlias. which are a string, a string,a string and boolean, respectively.
type ResponseData struct {
	OriginalURL string `json:"originalURL,omitempty"`
	ShortURL    string `json:"shortURL,omitempty"`
	ID          string `json:"id,omitempty"`
	IsAlias     bool   `json:"isAlias,omitempty"`
}

//Init function initialize the html file
func Init() {
	tpl = template.Must(template.ParseGlob("view/*.html"))
}

// type HostNumber struct {
// 	PortNumber string
// }

//Index function show the home page of server.
func Index(writer http.ResponseWriter, request *http.Request) {
	fmt.Println("Show home page")
	if Host = os.Getenv("HOST"); Host == "" {
		Host = "http://localhost:8000/"
	}
	if collection = os.Getenv("Collection"); collection == "" {
		collection = "url_table"
	}
	data := struct {
		HostDomain string
	}{
		HostDomain: Host,
	}

	Init()
	// var Host HostNumber
	// hostNumber = request.Host
	// Host.PortNumber = request.Host
	// tpl.ExecuteTemplate(writer, "app.html", Host)

	err := tpl.ExecuteTemplate(writer, "index.html", data)
	if err != nil {
		writer.WriteHeader(400)
		writer.Write([]byte("error load index.html"))
	}
}

//PrefixSlash function handle input Long URL with duplicate skahes and spaces.
func PrefixSlash(inputURL string) string {
	url := strings.TrimSpace(inputURL)
	hasSlash := strings.HasSuffix(url, "/")
	url = path.Clean(url)
	if hasSlash && !strings.HasSuffix(url, "/") {
		url += "/"
	} else if strings.HasSuffix(url, ".com") {
		url += "/"
	}
	url = strings.Replace(url, "https:/", "https://", 1)
	url = strings.Replace(url, "http:/", "http://", 1)
	return url
}

//CreateURL function create a new short URL and excute the template to show users.
func CreateURL(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Into createURL function")

	//Get the data from request
	var request RequestData
	str := r.PostFormValue("originalURL")
	// fmt.Println("The Host Number is" + hostNumber)
	request.OriginalURL = PrefixSlash(str)
	request.Alias = r.PostFormValue("alias")
	fmt.Println("User Input URL:", request.OriginalURL)
	fmt.Println("User Input Alias:", request.Alias)

	var res NotAvailable
	//Check URL domain, forbidden input http://localhost:8000
	forbiddenInput := Host
	if strings.Contains(request.OriginalURL, forbiddenInput) {
		res.Title = "URL domain banned"
		tpl.ExecuteTemplate(w, "notAvailable.html", res)
		return
	}
	if request.Alias == "" { //No Custom Alias input
		fmt.Println("Into No alias process")
		requestData := CreateWithoutAlias(request)
		tpl.ExecuteTemplate(w, "create.html", requestData)
	} else { // With Custom Alias
		fmt.Println("Into alias process")
		requestData := CreateWithAlias(request)
		if (requestData == ResponseData{}) { // Custom alias is not available
			res.Title = "Alias is not available"
			tpl.ExecuteTemplate(w, "notAvailable.html", res)
		} else { //Custom alias is available
			tpl.ExecuteTemplate(w, "create.html", requestData)
		}
	}
	fmt.Println("Send CreateURL respond successful")
}

//CreateWithAlias handle the condition that user input custom alias
func CreateWithAlias(request RequestData) ResponseData {
	/**
	 *	The function handle four situations:
	 *  non-exist long URL and non-exist custom Alias -> creat new short URK with alias
	 *  non-exist long URL and exist custom Alias     -> Custom alias is not available
	 *  exist long URL and non-exist custom Alias     -> creat new short URK with alias
	 * 	exist long URL and exist custom Alias         -> get previous short URL
	 */
	var response ResponseData
	collection := mongoDB.MongoClient.Database("url_database").Collection(collection)
	collection.FindOne(context.Background(), bson.M{"ID": request.Alias}).Decode(&response)
	if (response == ResponseData{}) { //The custom alias is available
		response.ID = request.Alias
		response.ShortURL = Host + response.ID
		response.OriginalURL = request.OriginalURL
		collection.InsertOne(context.TODO(), bson.M{
			"ID":          response.ID,
			"OriginalURL": response.OriginalURL,
			"ShortURL":    response.ShortURL,
			"IsAlias":     true})
		fmt.Println("Insert custom alias to DB successful")
	} else if response.OriginalURL != request.OriginalURL { //The custom alias is not available
		response = ResponseData{}
		fmt.Println("The alias is not available")
	}
	fmt.Println("Create short URL with custom alias successful")
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
	collection := mongoDB.MongoClient.Database("url_database").Collection(collection)
	collection.FindOne(context.Background(), bson.M{"OriginalURL": request.OriginalURL, "IsAlias": false}).Decode(&response)

	if (response == ResponseData{}) {
		response.ID = ProduceUniqueID()
		response.ShortURL = Host + response.ID
		response.OriginalURL = request.OriginalURL
		collection.InsertOne(context.TODO(), bson.M{
			"ID":          response.ID,
			"OriginalURL": response.OriginalURL,
			"ShortURL":    response.ShortURL,
			"IsAlias":     false})
		fmt.Println("Insert successful")
	}
	fmt.Println("Create short URL without custom alias successful")
	return response
}

//Redirect Function redirect the short URL to long URL
func Redirect(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Into Redirect")
	params := mux.Vars(r)
	collection := mongoDB.MongoClient.Database("url_database").Collection(collection)
	var response ResponseData
	collection.FindOne(context.Background(), bson.M{"ID": params["id"]}).Decode(&response)
	if (response != ResponseData{}) { //Find the Short URL
		http.Redirect(w, r, response.OriginalURL, 301)
		fmt.Println("Redirect successful")
	} else { //Can not found the page, short URL does not exist
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

//ProduceUniqueID return a unique ID by time
func ProduceUniqueID() string {
	h, _ := hashids.NewWithData(hashids.NewData())
	now := time.Now()
	ID, _ := h.Encode([]int{int(now.Unix())})
	return ID
}

//IsValidAlias function
// func IsValidAlias(s string) bool {
// 	result, _ := regexp.MatchString("^[a-zA-Z0-9]+$", s)
// 	return result
// }
