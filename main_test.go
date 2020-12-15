/**
 * This main_test file test the function of this shorten URL server
 *
 * @author Boris
 * @version 2020-12-06
 *
 */
package main

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"net/url"
	controller "short_url/controller"

	mongoDB "short_url/module"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

func Router() *mux.Router {
	router := mux.NewRouter()
	// router.HandleFunc("/", RootEndpoint).Methods("GET")
	mongoDB.InitRun()
	controller.Init()
	router.HandleFunc("/", controller.Index).Methods("GET")
	router.HandleFunc("/POST/url_resource", controller.CreateURL).Methods("POST")
	router.HandleFunc("/{id}", controller.Redirect).Methods("GET")
	return router
}

//Test show index page
func TestIndex_1(t *testing.T) {
	request, _ := http.NewRequest("GET", "/", nil)
	response := httptest.NewRecorder()
	Router().ServeHTTP(response, request)
	assert.Equal(t, 200, response.Code, "expected Stats Code: 200")
}

//Test input non-valid word in URL
func TestIndex_2(t *testing.T) {
	request, _ := http.NewRequest("GET", "/home", nil)
	response := httptest.NewRecorder()
	Router().ServeHTTP(response, request)
	assert.Equal(t, 404, response.Code, "expected Stats Code: 404")
}

//Test user input a non-exist short URL
func TestRedirect_1(t *testing.T) {
	request, _ := http.NewRequest("GET", "/12345", nil) //12345 does not exist in the database
	response := httptest.NewRecorder()
	Router().ServeHTTP(response, request)
	assert.Equal(t, 404, response.Code, "expected Stats Code: 404")
}

//Test user input an exist short URL
func TestRedirect_2(t *testing.T) {
	request, _ := http.NewRequest("GET", "/wE3MrgM", nil) //wE3MrgM exist in the database
	response := httptest.NewRecorder()
	Router().ServeHTTP(response, request)
	resp := response.Result()
	originalURL, _ := resp.Location()
	assert.Equal(t, 301, response.Code, "expected Stats Code: 301")
	assert.Equal(t, originalURL.String(), "https://www.google.com/", "expected Stats Code: 301")
}

//Test CreatURL function
func TestCreateURL(t *testing.T) {
	data := url.Values{}
	data.Add("originalURL", "https://www.google.com/")
	data.Add("alias", "")
	request, _ := http.NewRequest("POST", "/POST/url_resource", bytes.NewBufferString(data.Encode()))
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded; param=value")
	response := httptest.NewRecorder()
	Router().ServeHTTP(response, request)
	assert.Equal(t, 200, response.Code, "expected Stats Code: 200")
}

// Test creat a short URL that Long URL doest not exist in database
func TestCreateWithoutAlias_1(t *testing.T) {
	mongoDB.InitRun()
	var request controller.RequestData
	request.OriginalURL = "https://www.netflix.com/" //non-exist Long URL
	request.Alias = ""
	responseData := controller.CreateWithoutAlias(request)
	assert.Equal(t, responseData.OriginalURL, request.OriginalURL, "error in CreateWithoutAlias function")
	assert.Equal(t, false, responseData.IsAlias, "error in CreateWithoutAlias function")
}

//Test create a short URL that Long URL exist in database
func TestCreateWithoutAlias_2(t *testing.T) {
	mongoDB.InitRun()
	var request controller.RequestData
	request.OriginalURL = "https://www.google.com/" // exist  Long URL
	request.Alias = ""
	responseData := controller.CreateWithoutAlias(request)
	expected := "http://localhost:8000/wE3MrgM" // short URL of exist long URL
	assert.Equal(t, expected, responseData.ShortURL, "error in CreateWithoutAlias function")
	assert.Equal(t, responseData.IsAlias, false, "error in CreateWithoutAlias function")

	request.OriginalURL = "https://www.spotify.com/" // exist  Long URL
	request.Alias = ""
	responseData = controller.CreateWithoutAlias(request)
	expected = "http://localhost:8000/xVDk18B" // short URL of exist long URL
	assert.Equal(t, expected, responseData.ShortURL, "error in CreateWithoutAlias function")
	assert.Equal(t, false, responseData.IsAlias, "error in CreateWithoutAlias function")
}

//Test input duplicate  slashes and spaces in the long URL
func TestPrefixSlash(t *testing.T) {
	actual1 := controller.PrefixSlash("https://www.google.com/")        //normal
	actual2 := controller.PrefixSlash("https://www.google.com")         //No slash
	actual3 := controller.PrefixSlash("https://www.google.com////////") //duplicate slash
	actual4 := controller.PrefixSlash("https://www.google.com        ") //duplicate space in the end
	actual5 := controller.PrefixSlash("        https://www.google.com") //duplicated in the head
	expeected := "https://www.google.com/"
	assert.Equal(t, expeected, actual1, "Error in PrefixSlash function")
	assert.Equal(t, expeected, actual2, "Error in PrefixSlash function")
	assert.Equal(t, expeected, actual3, "Error in PrefixSlash function")
	assert.Equal(t, expeected, actual4, "Error in PrefixSlash function")
	assert.Equal(t, expeected, actual5, "Error in PrefixSlash function")
}

//Test create short URL with non-exist long URL and non exist custom alias
func TestCreateWithAlias_1(t *testing.T) {
	mongoDB.InitRun()
	var request controller.RequestData
	request.OriginalURL = "https://www.netflix.com/" // non-exist  Long URL
	request.Alias = "watchNetflix"                   //non-exist custom alias
	responseData := controller.CreateWithAlias(request)
	expected := "http://localhost:8000/watchNetflix" // short URL with custom alias
	assert.Equal(t, expected, responseData.ShortURL, "error in CreateWithAlias function")
	assert.Equal(t, true, responseData.IsAlias, "error in CreateWithAlias function")

	request.OriginalURL = "https://www.spotify.com/" // non-exist  Long URL
	request.Alias = "listenSpotify"                  //non-exist custom alias
	responseData = controller.CreateWithAlias(request)
	expected = "http://localhost:8000/listenSpotify" // short URL with custom alias
	assert.Equal(t, expected, responseData.ShortURL, "error in CreateWithAlias function")
	assert.Equal(t, true, responseData.IsAlias, "error in CreateWithAlias function")
}

//Test create short URL with non-exist long URL and exist custom alias
func TestCreateWithAlias_2(t *testing.T) {
	mongoDB.InitRun()
	var request controller.RequestData
	request.OriginalURL = "https://www.birmingham.ac.uk/" // non-exist  Long URL
	request.Alias = "watchNetflix"                        //exist custom alias
	responseData := controller.CreateWithAlias(request)
	expected := controller.ResponseData{} // the expected should get an empty means the alias is not available
	assert.Equal(t, expected, responseData, "error in CreateWithAlias function")
}

//Test create short URL with exist long URL and non-exist custom alias
func TestCreateWithAlias_3(t *testing.T) {
	mongoDB.InitRun()
	var request controller.RequestData
	request.OriginalURL = "https://www.spotify.com/" // non-exist  Long URL
	request.Alias = "MySpotify"                      //exist custom alias
	responseData := controller.CreateWithAlias(request)
	expected := "http://localhost:8000/MySpotify" // the expected should get an empty means the alias is not available
	assert.Equal(t, expected, responseData.ShortURL, "error in CreateWithAlias function")
	assert.Equal(t, true, responseData.IsAlias, "error in CreateWithAlias function")
}
