// package main

// import (
// 	"fmt"
// 	g "gopher"
// )

// func main() {
// 	fmt.Println("hello World")
// 	g.ShowDog()
// }
package main

import (
	"fmt"
	"net/http"
)

func hello(w http.ResponseWriter, req *http.Request) {
	// fmt.Fprintf(w, "hello World\n")
	fmt.Fprintf(w, printDog())
}

func printDog() string {
	var result string
	result =
		"    `.-::::::-.`    \n" +
			"    `.-::::::-.`    \n" +
			"    `.-::::::-.`    \n" +
			".:-::::::::::::::-:.\n" +
			"`_:::    ::    :::_`\n" +
			" .:( ^   :: ^   ):. \n" +
			" `:::   (..)   :::. \n" +
			" `:::::::UU:::::::` \n" +
			" .::::::::::::::::. \n" +
			" O::::::::::::::::O \n" +
			" -::::::::::::::::- \n" +
			" `::::::::::::::::` \n" +
			"  .::::::::::::::.  \n" +
			"    oO:::::::Oo     \n"
	return result
}

func headers(w http.ResponseWriter, req *http.Request) {
	for name, headers := range req.Header {
		for _, h := range headers {
			fmt.Fprintf(w, "%v: %v\n", name, h)
		}
	}
}

// func main() {
// 	// http.HandleFunc("/hello", hello)
// 	// http.HandleFunc("/headers", headers)

// 	// http.ListenAndServe(":5001", nil)
// }
