/**
 * This main file excute the web server.
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

func main() {
	fmt.Println("Server start")
	router := mux.NewRouter()

	//Init Database
	mongoDB.InitRun()

	//Show HomePage
	router.HandleFunc("/", controller.Index).Methods("GET")

	//Creat Short URL
	router.HandleFunc("/POST/url_resource", controller.CreateURL).Methods("POST")

	//Redirect Original URL
	router.HandleFunc("/{id}", controller.Redirect).Methods("GET")

	//Server listen at port 8000
	log.Fatal(http.ListenAndServe(":8000", router))
}
