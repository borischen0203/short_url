package controller

import (
	"fmt"
	"html/template"
	"net/http"
	"regexp"
)

type pageData struct {
	Title string

	Alias   string
	LongURL string
}

const (
	AliasURL = "https://BorisLongToShort.com/"
)

var tpl *template.Template

func Init() {
	tpl = template.Must(template.ParseGlob("view/*.html"))
}

var count int

/**
 * This function show the index page
 *
 * @ writer
 * @ request
 *
 */
func Index(writer http.ResponseWriter, request *http.Request) {
	Init()
	tpl.ExecuteTemplate(writer, "app.html", nil)

}

// func Handler(writer http.ResponseWriter, request *http.Request) {
// 	fmt.Println("into Server Handler", count)
// 	//Init()
// 	tpl.ExecuteTemplate(writer, "app.html", nil)
// 	// if request != nil {
// 	// 	fmt.Println("!= nil What is :", request.Body)
// 	// } else {
// 	// 	fmt.Println("== nil What is :", request.Body)
// 	// }
// 	// fmt.Fprintf(writer, "Hello, World")
// 	//template := template.Must(template.ParseFiles("view/test.html"))
// 	// template := template.Must(template.ParseGlob("view/*.html"))
// 	// data := new(pageData)
// 	// data.Title = "Long URL 2 Short URL"
// 	// data.Long_url = "This is long_url"
// 	// data.Short_url = "This is short_url"
// 	//template.Execute(writer, nil)
// }

func HandleURL(writer http.ResponseWriter, request *http.Request) {
	fmt.Println("into URL Handler")
	if request.Method != "POST" {
		http.Redirect(writer, request, "/", http.StatusSeeOther)
		return
	}

	LongURL := request.FormValue("url")
	fmt.Println("Long ", LongURL)
	Alias := request.FormValue("alias")
	fmt.Println("Alias", Alias)
	ShortURL := AliasURL + Alias
	fmt.Println("short", ShortURL)

	result := struct {
		Original string
		Shorten  string
	}{
		Original: LongURL,
		Shorten:  ShortURL,
	}
	fmt.Println(result)

	tpl.ExecuteTemplate(writer, "submission.html", result)

}

func isValidAlias(s string) bool {
	result, _ := regexp.MatchString("^[a-zA-Z0-9]+$", s)
	return result
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
