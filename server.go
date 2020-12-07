package main

import (
	"html/template"
	"net/http"
)

type IndexData struct {
	Title   string
	Content string
}

func handler(writer http.ResponseWriter, request *http.Request) {
	// fmt.Fprintf(writer, "Hello, World")
	template := template.Must(template.ParseFiles("./index.html"))
	data := new(IndexData)
	data.Title = "Home page"
	data.Content = "This is conetent!"
	template.Execute(writer, data)
}

func main() {

	http.HandleFunc("/", handler)

	// Start a web server.
	http.ListenAndServe(":8080", nil)
}
