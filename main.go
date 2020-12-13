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
	server "short_url/controller"

	mongoDB "short_url/module"

	"github.com/gorilla/mux"
)

func main() {
	fmt.Println("Server start")
	router := mux.NewRouter()

	//Init Database
	mongoDB.InitRun()

	//Show HomePage
	router.HandleFunc("/", server.Index)

	//Creat Short URL
	// router.HandleFunc("/create", server.CreateURL).Methods("POST")
	router.HandleFunc("/POST/url_resource", server.CreateURL).Methods("POST")

	//Redirect Original URL
	router.HandleFunc("/{id}", server.Redirect).Methods("GET")

	log.Fatal(http.ListenAndServe(":8000", router))
}
