/**
 * This main file excute the web server
 *
 * @author Boris
 * @version 2020-12-06
 *
 */

package main

import (
	"fmt"
	"log"
	"net/http"
	controller "short_url/controller"
	mongoDB "short_url/module"

	"github.com/gorilla/mux"
)

func RootEndpoint(response http.ResponseWriter, request *http.Request) {
	response.WriteHeader(200)
	response.Write([]byte("Hello World"))
}

func main() {
	fmt.Println("Server start")
	router := mux.NewRouter()

	// router.HandleFunc("/", RootEndpoint).Methods("GET")
	// // Init Database
	mongoDB.InitRun()

	//Show HomePage
	router.HandleFunc("/", controller.Index).Methods("GET")

	//Creat Short URL
	router.HandleFunc("/POST/url_resource", controller.CreateURL).Methods("POST")

	//Redirect Original URL
	router.HandleFunc("/{id}", controller.Redirect).Methods("GET")
	log.Fatal(http.ListenAndServe(":8000", router))
}
