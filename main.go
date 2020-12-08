package main

import (
	"fmt"
	"log"
	"net/http"
	server "short_url/server"
)

// func hello(w http.ResponseWriter, req *http.Request) {
// 	fmt.Fprintf(w, "hello World\n")
// 	// fmt.Fprintf(w, printDog())
// }

// func printDog() string {
// 	var result string
// 	result =
// 		"    `.-::::::-.`    \n" +
// 			"    `.-::::::-.`    \n" +
// 			"    `.-::::::-.`    \n" +
// 			".:-::::::::::::::-:.\n" +
// 			"`_:::    ::    :::_`\n" +
// 			" .:( ^   :: ^   ):. \n" +
// 			" `:::   (..)   :::. \n" +
// 			" `:::::::UU:::::::` \n" +
// 			" .::::::::::::::::. \n" +
// 			" O::::::::::::::::O \n" +
// 			" -::::::::::::::::- \n" +
// 			" `::::::::::::::::` \n" +
// 			"  .::::::::::::::.  \n" +
// 			"    oO:::::::Oo     \n"
// 	return result
// }

// func headers(w http.ResponseWriter, req *http.Request) {
// 	for name, headers := range req.Header {
// 		for _, h := range headers {
// 			fmt.Fprintf(w, "%v: %v\n", name, h)
// 		}
// 	}
// }

func main() {
	fmt.Println("Server start")
	http.HandleFunc("/", server.Handler)

	err := http.ListenAndServe(":5500", nil)
	if err != nil {
		log.Fatal("ListenAdnServe", err)
	} else {
		log.Println("Listen 5050")
	}
}
