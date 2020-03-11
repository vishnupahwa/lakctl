package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

func main() {
	firstarg := "World!"
	if len(os.Args) > 1 {
		firstarg = os.Args[1]
	}
	log.Println("Running testapp with argument " + firstarg)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		_, _ = fmt.Fprint(w, "Hello "+firstarg)
	})
	log.Fatal(http.ListenAndServe(":9999", nil))
}
