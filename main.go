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
)

func main() {
	fmt.Println("Server start")
	// http.HandleFunc("/", server.Handler)
	http.HandleFunc("/", controller.Index)
	http.HandleFunc("/submission", controller.HandleURL)

	err := http.ListenAndServe(":5500", nil)
	if err != nil {
		log.Fatal("ListenAdnServe", err)
	} else {
		log.Println("Listen 5050")
	}
}
