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
	"strings"

	// "log"
	"net/http"
	controller "short_url/controller"
	// mongoDB "short_url/module"
	// "github.com/gorilla/mux"
)

func Add(value1 int, value2 int) int {
	return value1 + value2
}

func RootEndpoint(response http.ResponseWriter, request *http.Request) {
	response.WriteHeader(200)
	response.Write([]byte("Hello World"))
}

func main() {
	fmt.Println("Server start")
	// router := mux.NewRouter()
	// router.HandleFunc("/", RootEndpoint).Methods("GET")
	str := "https://www.google.com/                     "
	fmt.Println(controller.PrefixSlash(str))
	p := "https://www.google.com//"
	p = strings.TrimSuffix(p, "/")
	fmt.Println(p)
	// // Init Database
	// mongoDB.InitRun()

	// //Show HomePage
	// router.HandleFunc("/", controller.Index).Methods("GET")

	// //Creat Short URL
	// router.HandleFunc("/POST/url_resource", controller.CreateURL).Methods("POST")

	// //Redirect Original URL
	// router.HandleFunc("/{id}", controller.Redirect).Methods("GET")

	// log.Fatal(http.ListenAndServe(":8000", router))
}
