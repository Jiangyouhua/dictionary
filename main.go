package main

import (
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("handle", handle)
	http.Handle("/", http.FileServer(http.Dir("./")))
	if err := http.ListenAndServe(":80", nil); err != nil {
		log.Fatal(err)
	}
}

func handle(w http.ResponseWriter, r *http.Request) {

}
