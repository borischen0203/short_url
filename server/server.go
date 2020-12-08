package server

import (
	"fmt"
	"html/template"
	"net/http"
)

type pageData struct {
	Title     string
	Short_url string
	Long_url  string
}

func Handler(writer http.ResponseWriter, request *http.Request) {

	fmt.Println("into Server Handler")
	if request != nil {
		fmt.Println("!= nil What is :", request.Body)
	} else {
		fmt.Println("== nil What is :", request.Body)
	}
	// fmt.Fprintf(writer, "Hello, World")
	template := template.Must(template.ParseFiles("view/test.html"))
	// template := template.Must(template.ParseGlob("view/*.html"))
	data := new(pageData)
	data.Title = "Long URL 2 Short URL"
	data.Long_url = "This is long_url"
	data.Short_url = "This is short_url"
	template.Execute(writer, data)
}

// func main() {
// 	//set visit url
// 	http.HandleFunc("/home", handler)

// 	/*
// 	 * Start a web server.
// 	 * if fail, log.Fatal show error in output
// 	 */
// 	log.Fatal(http.ListenAndServe(":8080", nil))
// }
