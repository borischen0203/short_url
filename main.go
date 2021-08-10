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
	"os"

	controller "github.com/borischen0203/short_url/controller"
	mongoDB "github.com/borischen0203/short_url/module"

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
	router.HandleFunc("/api/url_resource", controller.CreateURL).Methods("POST")

	//Redirect Original URL
	router.HandleFunc("/{id}", controller.Redirect).Methods("GET")

	//Server listen at port 8000
	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
		log.Printf("Defaulting to port %s", port)
	}

	log.Printf("Listening on port %s", port)
	if err := http.ListenAndServe(":"+port, router); err != nil {
		log.Fatal(err)
	}
}
