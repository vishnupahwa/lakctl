package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	log.Println("Running testapp")
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		_, _ = fmt.Fprint(w, "Hello World!")
	})
	log.Fatal(http.ListenAndServe(":9999", nil))
}
