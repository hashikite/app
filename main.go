package main

import (
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "index.html")
	})

	http.HandleFunc("/kite.gif", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "kite.gif")
	})

	log.Println("Serving hashy kites on http://0.0.0.0:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
