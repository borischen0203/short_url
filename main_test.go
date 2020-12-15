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

// func TestRootEndpoint(t *testing.T) {
// 	request, _ := http.NewRequest("GET", "/", nil)
// 	response := httptest.NewRecorder()
// 	Router().ServeHTTP(response, request)
// 	// fmt.Println(response.Code)
// 	assert.Equal(t, 200, response.Code, "expected Stats Code: 200")
// 	// fmt.Println(response.Body)
// 	// assert.Equal(t, "Hello World", response.Body.String(), "Incorrect body")
// }

func TestIndex_1(t *testing.T) {
	request, _ := http.NewRequest("GET", "/", nil)
	response := httptest.NewRecorder()
	Router().ServeHTTP(response, request)
	assert.Equal(t, 200, response.Code, "expected Stats Code: 200")
}

func TestIndex_2(t *testing.T) {
	request, _ := http.NewRequest("GET", "/home", nil)
	response := httptest.NewRecorder()
	Router().ServeHTTP(response, request)
	assert.Equal(t, 404, response.Code, "expected Stats Code: 404")
}

//TestRedirect_1 function test user input a non-exist short URL
func TestRedirect_1(t *testing.T) {
	request, _ := http.NewRequest("GET", "/12345", nil) //12345 does not exist in the database
	response := httptest.NewRecorder()
	Router().ServeHTTP(response, request)
	assert.Equal(t, 404, response.Code, "expected Stats Code: 404")
}

//TestRedirect_2 function test user input an exist short URL
func TestRedirect_2(t *testing.T) {
	request, _ := http.NewRequest("GET", "/wE3MrgM", nil) //wE3MrgM exist in the database
	response := httptest.NewRecorder()
	Router().ServeHTTP(response, request)
	resp := response.Result()
	originalURL, _ := resp.Location()
	assert.Equal(t, 301, response.Code, "expected Stats Code: 301")
	assert.Equal(t, originalURL.String(), "https://www.google.com/", "expected Stats Code: 301")
}

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
